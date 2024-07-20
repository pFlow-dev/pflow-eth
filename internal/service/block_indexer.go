package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/lib/pq"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"github.com/pflow-dev/pflow-eth/metamodel"
	"math/big"
	"os"
	"time"
)

const notify_channel = "node_sync_channel"

func (s *Service) subscribeToBlockchainEvents(sink chan *metamodel.MetamodelSignaledEvent) (event.Subscription, error) {
	address := common.HexToAddress(config.Address)
	m, err := metamodel.NewMetamodel(address, s.Client)
	if err != nil {
		panic(err)
	}
	opts := &bind.WatchOpts{Context: context.Background()}
	scalars := []*big.Int{}
	scalars = append(scalars, big.NewInt(0))
	scalars = append(scalars, big.NewInt(2))
	scalars = append(scalars, big.NewInt(1))

	// REVIEW: can we listen to all events?
	return m.WatchSignaledEvent(opts, sink, []uint8{0, 1}, []uint8{0, 1, 2, 3}, scalars)
}

func (s *Service) getNotificationListener() *pq.Listener {
	_, err := sql.Open("postgres", config.DbConn)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(config.DbConn, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen(notify_channel)
	if err != nil {
		panic(err)
	}

	fmt.Println("psql listener: " + notify_channel)
	return listener
}

func (s *Service) runBlockIndexer(ctx context.Context, networks map[string]bool) {
	s.Event("block_indexer_started", map[string]interface{}{
		"pid":      os.Getpid(),
		"networks": networks,
	})
	listener := s.getNotificationListener()
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	eventChannel := make(chan *metamodel.MetamodelSignaledEvent)

	eventSub, subError := s.subscribeToBlockchainEvents(eventChannel)
	if subError != nil {
		fmt.Printf("event subscription disabled: %s", subError)
	}
	start := time.Now()

	syncBlockchainNetworks := func() {
		for schema, enabled := range networks {
			if enabled {
				s.syncNodeWithBlockchain(schema)
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			s.Event("block_indexer_stopped", map[string]interface{}{
				"pid":      os.Getpid(),
				"networks": networks,
				"duration": time.Since(start).String(),
			})
			eventSub.Unsubscribe()
			_ = listener.UnlistenAll()
			os.Exit(0)
		case notification := <-listener.Notify:
			if schema := notification.Extra; networks[schema] == true {
				fmt.Println("refresh notification from schema: " + schema)
				s.confirmSentTransactions(schema)
				for {
					if !s.importNextTransaction(schema) {
						break
					}
				}
				s.computeTransactionStates(schema)
			} else {
				fmt.Println("notification from unknown schema: " + schema)
			}
		case <-ticker.C:
			syncBlockchainNetworks()
		case evt := <-eventChannel:
			fmt.Printf("Received event: %v\n", evt)
			syncBlockchainNetworks()
		}
	}
}
