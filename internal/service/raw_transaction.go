package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"net/http"
)

type TransactionRequest struct {
	RawTransaction string `json:"rawTransaction"`
}

func (s *Service) PublishRawTransaction(ctx context.Context, rawTx string) (string, error) {
	// Example: Connect to an Ethereum node
	client, err := ethclient.Dial(config.Endpoint)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Decode the transaction
	txBytes, err := hexutil.Decode(rawTx)
	if err != nil {
		return "", err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		return "", err
	}

	// Example: Send the transaction
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

// RawTransactionHandler handles the HTTP request for publishing a raw Ethereum transaction
func (s *Service) RawTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	txHash, err := s.PublishRawTransaction(r.Context(), req.RawTransaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish transaction: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"txHash": txHash}
	json.NewEncoder(w).Encode(response)
}
