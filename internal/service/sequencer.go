package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pflow-dev/pflow-eth/config"
	"log"
	"os"
	"time"
)

func (s *Server) getHightestSequence() (highestSequence int64) {
	err := s.NodeDb.QueryRow("SELECT MAX(sequence) FROM transaction_states").Scan(&highestSequence)
	if err != nil {
		log.Printf("Error finding highest sequence: %v", err)
		return -1
	}
	return highestSequence
}

func (s *Server) insertSequence(nextSeq int64) (int64, error) {
	sqlStatement := `
    INSERT INTO transaction_states (sequence, transaction_hash, state)
    SELECT sequence, transaction_hash, NULL
    FROM (
        SELECT sequence, transaction_hash
        FROM transaction_log_view
        WHERE sequence = {{}} 
    ) AS next_sequence
    `

	result, err := s.NodeDb.Exec(sqlStatement, nextSeq)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Server) SyncNodeWithBlockchain() bool {
	_, err := s.NodeDb.Exec("REFRESH MATERIALIZED VIEW node_sync_data_view")
	if err != nil {
		log.Printf("Error refreshing materialized view node_sync_data_view: %v", err)
		return false
	} else {
		log.Println("Materialized view node_sync_data_view refreshed successfully.")
		return true
	}
}

func (s *Server) ConfirmSentTransactions() (ok bool) {
	var jsonString string
	err := s.NodeDb.QueryRow("SELECT confirm_transactions();").Scan(&jsonString)
	if err != nil {
		log.Printf("Error calling remove_confirmed_transactions: %v", err)
		return
	}

	// Unmarshal the JSON string into a map
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		log.Printf("Error unmarshaling JSON from confirm_transactions: %v", err)
		return false
	}

	s.Event("confirm_transactions", result)
	return true
}

func (s *Server) ImportNextTransaction() (ok bool) {
	highestSequence := s.getHightestSequence()
	if highestSequence < 0 {
		return false
	}
	nextSeq := highestSequence + 1
	rows, err := s.insertSequence(nextSeq)
	if err != nil {
		log.Printf("Error inserting sequence %d: %v", nextSeq, err)
		return false
	}
	if rows > 0 {
		s.Event("sequence", map[string]interface{}{
			"sequence": nextSeq,
		})
		log.Printf("Inserted sequence %d", nextSeq)
		return true
	}
	return false
}

func (s *Server) ComputeTransactionStates() (ok bool) {
	// REVIEW: test this function
	sqlStatement := `SELECT sequence, transaction_hash FROM transaction_states WHERE state IS NULL order by sequence asc`

	// Execute the query
	rows, err := s.NodeDb.Query(sqlStatement)
	if err != nil {
		log.Printf("Error querying null state transactions: %v", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var sequence int64
		var transactionHash string

		err := rows.Scan(&sequence, &transactionHash)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return false
		}
		log.Printf("Found null state transaction: sequence %d, hash %s", sequence, transactionHash)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return false
	}

	return true
}

func (s *Server) getNotificationListener() *pq.Listener {
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
	err = listener.Listen("node_sync_channel")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	return listener
}

func (s *Server) runSequencer(ctx context.Context) {
	s.Event("sequencer_started", map[string]interface{}{
		"pid": os.Getpid(),
	})
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	start := time.Now()

	listener := s.getNotificationListener()

	for {
		select {
		case <-ctx.Done():
			s.Event("sequencer_stopped", map[string]interface{}{
				"pid":      os.Getpid(),
				"duration": time.Since(start).String(),
			})
			os.Exit(0)
		case notification := <-listener.Notify:
			fmt.Printf("Received notification on channel %s: %v\n", notification.Channel, notification.Extra)
			s.ConfirmSentTransactions()
			for {
				if !s.ImportNextTransaction() {
					break
				}
			}
			s.ComputeTransactionStates()
		case <-ticker.C:
			s.SyncNodeWithBlockchain()
		}
	}
}
