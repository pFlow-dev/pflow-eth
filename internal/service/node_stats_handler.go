package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (s Server) NodeStatsHandler(w http.ResponseWriter, _ *http.Request) {
	// Prepare the SQL query to select the latest stats from the materialized view
	query := `SELECT sync_data FROM node_sync_data_view LIMIT 1;`

	// Variable to hold the JSON data
	var syncData json.RawMessage

	// Execute the query
	err := s.NodeDb.QueryRow(query).Scan(&syncData)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No data available", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	_, err = w.Write(syncData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
