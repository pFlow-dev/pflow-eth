package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"net/http"
)

// construct a uuid build string from the JS and CSS build numbers
func getBuild() string {
	hexString := config.JsBuild + "000000000001" + "0000" + config.CssBuild

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexString[0:8], hexString[8:12], hexString[12:16], hexString[16:20], hexString[20:32])
}

func (s Service) SyncStatsHandler(w http.ResponseWriter, r *http.Request) {
	s.setSearchPathForRequest(r)

	query := `SELECT sync_data FROM node_sync_data_view LIMIT 1;`

	// Variable to hold the JSON data
	var syncData json.RawMessage

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
	status := make(map[string]interface{})
	err = json.Unmarshal(syncData, &status)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}
	status["build"] = getBuild()
	err = json.NewEncoder(w).Encode(status)
}
