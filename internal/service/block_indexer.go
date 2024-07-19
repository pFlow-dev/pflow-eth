package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"os"
	"time"
)

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

	notify_channel := "node_sync_channel"
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
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			s.Event("block_indexer_stopped", map[string]interface{}{
				"pid":      os.Getpid(),
				"networks": networks,
				"duration": time.Since(start).String(),
			})
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
			for schema, enabled := range networks {
				if enabled {
					s.syncNodeWithBlockchain(schema)
				}
			}
		}
	}
}
