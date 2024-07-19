package service

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	json "github.com/gibson042/canonicaljson-go"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"math/big"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(response)
}

func (s *Service) FaucetHandler(w http.ResponseWriter, r *http.Request) {
	if s.getNetwork(r.Host) != "hardhat" {
		respondWithError(w, http.StatusBadRequest, "Faucet is only available on hardhat network")
		return
	}

	addr := r.URL.Query().Get("addr")
	if !common.IsHexAddress(addr) {
		respondWithError(w, http.StatusBadRequest, "Invalid address")
		return
	}

	client, err := ethclient.Dial(config.Endpoint)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to Ethereum network")
		return
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load private key")
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Error casting public key to ECDSA")
		return
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(r.Context(), fromAddress)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get nonce")
		return
	}

	value := big.NewInt(1e18) // 1 ETH
	gasLimit := uint64(21000) // ETH transfer gas limit
	gasPrice, err := client.SuggestGasPrice(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to suggest gas price")
		return
	}

	toAddress := common.HexToAddress(addr)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get network ID")
		return
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign transaction")
		return
	}

	err = client.SendTransaction(r.Context(), signedTx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send transaction")
		return
	}

	txHash := signedTx.Hash().Hex()
	response := map[string]string{"txHash": txHash}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
