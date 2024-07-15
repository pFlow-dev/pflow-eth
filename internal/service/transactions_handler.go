package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type TransactionLog struct {
	Sequence        int64  `json:"sequence"`
	BlockNumber     int64  `json:"block_number"`
	Role            int64  `json:"role"`
	Action          int64  `json:"action"`
	Scalar          int64  `json:"scalar"`
	FromAddress     string `json:"from_address"`
	TransactionHash string `json:"transaction_hash"`
	TopicHash       string `json:"topic_hash"`
	Removed         bool   `json:"removed"`
}

const transactionsQuery = `
SELECT
  sequence,
  block_number,
  role,
  action,
  scalar,
  from_address,
  transaction_hash,
  topic_hash,
  removed
FROM
  transaction_logs_view
`

func (s Server) TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	seqStr := r.URL.Query().Get("seq")
	if seqStr != "" {
		endStr = seqStr
		startStr = seqStr
	}

	var start, end int
	var err error

	if startStr != "" {
		start, err = strconv.Atoi(startStr)
		if err != nil {
			http.Error(w, "Invalid start parameter", http.StatusBadRequest)
			return
		}
	} else {
		start = 0
	}

	if endStr != "" {
		end, err = strconv.Atoi(endStr)
		if err != nil {
			http.Error(w, "Invalid end parameter", http.StatusBadRequest)
			return
		}
	}

	if start > end && end != 0 {
		http.Error(w, "Start parameter must be less than or equal to end parameter", http.StatusBadRequest)
		return
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(transactionsQuery)

	queryParams := make([]interface{}, 0)
	conditions := make([]string, 0)

	if start > 0 {
		conditions = append(conditions, "sequence >= $1")
		queryParams = append(queryParams, start)
	}

	if end > 0 {
		if len(conditions) > 0 {
			conditions = append(conditions, " AND ")
		}
		conditions = append(conditions, "sequence <= $2")
		queryParams = append(queryParams, end)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, ""))
	}

	queryBuilder.WriteString(" ORDER BY sequence DESC")

	fmt.Printf("Query: %s\n", queryBuilder.String())
	fmt.Printf("Params: %v\n", queryParams)

	rows, err := s.NodeDb.Query(queryBuilder.String(), queryParams...)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []TransactionLog
	for rows.Next() {
		var tx TransactionLog
		err = rows.Scan(&tx.Sequence, &tx.BlockNumber, &tx.Role, &tx.Action, &tx.Scalar, &tx.FromAddress, &tx.TransactionHash, &tx.TopicHash, &tx.Removed)
		if err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, tx)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(transactions) == 0 {
		// Return empty array instead of null
		transactions = make([]TransactionLog, 0)
	}
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
