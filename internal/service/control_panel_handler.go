package service

import (
	"encoding/json"
	"net/http"
)

// Define a structure for the JSON response
type ActionResponse struct {
	Action string `json:"action"`
	Result string `json:"result"`
}

func (s *Service) ResetDb() bool {
	_, err := s.NodeDb.Exec("SELECT truncate_node_data()")
	return err == nil
}

func (s *Service) InsertBlockNumber(blockNumber uint64) bool {
	stmt, err := s.NodeDb.Prepare("INSERT INTO block_numbers (block_number) VALUES ($1)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(blockNumber)
	return err == nil
}

func (s *Service) ControlPanelHandler(w http.ResponseWriter, r *http.Request) {
	if s.getNetwork(r.Host) != "hardhat" {
		respondWithError(w, http.StatusBadRequest, "ControlPanel is only available on hardhat network")
		return
	}
	action := r.URL.Query().Get("cmd")
	var result string

	switch action {
	case "reset_db":
		if s.ResetDb() {
			result = "Reset All Tables"
		} else {
			result = "Failed to reset database"
		}
	case "init_block_numbers":
		if s.InsertBlockNumber(0) {
			result = "Init Block 0"
		} else {
			result = "Failed to init block numbers"
		}
	case "sync":
		if s.syncNodeWithBlockchain(s.getNetwork(r.Host)) {
			result = "Synced with blockchain"
		} else {
			result = "Failed to sync with blockchain"
		}
	default:
		// If the action is invalid, immediately return a JSON response indicating the error
		respondWithJSON(w, http.StatusBadRequest, ActionResponse{Action: action, Result: "Invalid action"})
		return
	}

	// If the action is valid and has been processed, return a success JSON response
	respondWithJSON(w, http.StatusOK, ActionResponse{Action: action, Result: result})
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, response ActionResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
