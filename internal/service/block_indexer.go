package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/lib/pq"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"github.com/pflow-dev/pflow-eth/metamodel"
	"os"
	"strings"
	"time"
)

const notify_channel = "node_sync_channel"

func convertToWebSocketURL(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http://", "ws://", 1)
	} else if strings.HasPrefix(url, "https://") {
		return strings.Replace(url, "https://", "wss://", 1)
	}
	return url
}

func (s *Service) subscribeToBlockchainEvents(sink chan *metamodel.MetamodelSignaledEvent) (event.Subscription, error) {
	address := common.HexToAddress(config.Address)
	client, dialErr := ethclient.Dial(convertToWebSocketURL(config.Endpoint))
	if dialErr != nil {
		return nil, dialErr
	}
	m, err := metamodel.NewMetamodel(address, client)
	if err != nil {
		panic(err)
	}
	opts := &bind.WatchOpts{Context: context.Background()}
	return m.WatchSignaledEvent(opts, sink, nil, nil, nil)
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

func (s *Service) SignalEvent(evt *metamodel.MetamodelSignaledEvent) {
	next := make(map[string]interface{})
	next["sequence"] = evt.Sequence.Int64()
	next["role"] = evt.Role
	next["action"] = evt.ActionId
	next["scalar"] = evt.Scalar.Int64()

	next["address"] = evt.Raw.Address
	next["block"] = evt.Raw.BlockNumber
	next["tx"] = evt.Raw.TxHash
	s.Event("signaled_event", next)
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
		fmt.Printf("event_subscription disabled: %s\n", subError)
	} else {
		fmt.Printf("event_subscription: enabled\n")
	}
	start := time.Now()
	highestSeq := int64(-1)

	syncByBlocks := func() {
		for schema, enabled := range networks {
			if enabled {
				s.syncNodeWithBlockchain(schema)
			}
		}
	}
	refreshDataView := func() {
		for schema, enabled := range networks {
			if enabled {
				s.refreshDataView(schema)
			}
		}
	}

	insertBlock := func(address common.Address, blockNumber uint64) {
		for schema, enabled := range networks {
			if enabled {
				// FIXME: also check if the address belongs to the schema
				s.setSearchPath(schema)
				if s.InsertBlockNumber(blockNumber) {
					s.addressHeights[address] = blockNumber
					refreshDataView()
				}
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
				s.confirmSentTransactions(schema)
			} else {
				panic("notification from unknown schema: " + schema)
			}
		case <-ticker.C:
			onChainSequence := s.getContractSequence()
			if highestSeq < 0 {
				highestSeq = s.getHightestSequence()
			}
			if highestSeq != onChainSequence {
				syncByBlocks()
			}
		case evt := <-eventChannel:
			if highestSeq < 0 { // initialize from db
				highestSeq = s.getHightestSequence()
			}
			insertBlock(evt.Raw.Address, evt.Raw.BlockNumber)
			s.SignalEvent(evt)
			if highestSeq+1 == evt.Sequence.Int64() {
				highestSeq = evt.Sequence.Int64()
			} else {
				highestSeq = -1
				syncByBlocks()
			}
		}
	}
}
