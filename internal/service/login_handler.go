package service

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"net/http"
	"time"
)

type jsonResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func loginRespondWithJSON(w http.ResponseWriter, statusCode int, response jsonResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func verifySig(from, sigHex string, msg []byte) bool {
	sig := hexutil.MustDecode(sigHex)
	msg = accounts.TextHash(msg)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27
	}
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	return from == recoveredAddr.Hex()
}

func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	signature := r.URL.Query().Get("signature")
	if signature == "" {
		loginRespondWithJSON(w, http.StatusBadRequest, jsonResponse{Success: false, Error: "Missing signature"})
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		loginRespondWithJSON(w, http.StatusBadRequest, jsonResponse{Success: false, Error: "Missing address"})
		return
	}

	session := ""

	timeStamp := r.URL.Query().Get("ts")
	t, tsErr := time.Parse("2006-01-02T15:04:05.000Z", timeStamp)

	loginWindow := 15 * time.Second

	if tsErr != nil || t.Before(time.Now().Add(0-loginWindow)) || t.After(time.Now().Add(loginWindow)) {
		loginRespondWithJSON(w, http.StatusUnauthorized, jsonResponse{Success: false, Error: "Invalid timestamp"})
		return
	}

	if !verifySig(address, signature, []byte("authenticate "+timeStamp)) {
		loginRespondWithJSON(w, http.StatusUnauthorized, jsonResponse{Success: false, Error: "Invalid signature"})
		return
	}

	loginRespondWithJSON(w, http.StatusOK, jsonResponse{Success: true, Data: map[string]string{"session": session}})
}
