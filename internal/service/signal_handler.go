package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	json "github.com/gibson042/canonicaljson-go"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pflow-dev/pflow-eth/metamodel"
)

// Example private key (DO NOT USE IN PRODUCTION)
const privateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
const contractAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
const senderAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

var chainId = big.NewInt(31337)

type EventLog struct {
	Event []string `json:"event"`
	Data  string   `json:"data"`
}

type SignalResponse struct {
	TransactionHash string     `json:"transaction_hash"`
	Sender          string     `json:"sender"`
	Contract        string     `json:"contract"`
	EventLog        []EventLog `json:"event_log"`
}

func signTransaction(from common.Address, tx *types.Transaction) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	signer := types.LatestSignerForChainID(chainId)

	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		fmt.Printf("transaction_type: %T\n", tx)
		fmt.Printf("transaction: %v\n", tx)
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}
	return signedTx, nil
}

func (s *Service) Signal(ctx context.Context, action uint8, scalar *big.Int) (*SignalResponse, error) {
	contractAddress := common.HexToAddress(contractAddress)
	metamodelInstance, err := metamodel.NewMetamodel(contractAddress, s.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a Metamodel contract: %v", err)
	}

	txn, err := metamodelInstance.Signal(&bind.TransactOpts{
		From:    common.HexToAddress(senderAddress),
		Nonce:   nil,
		Signer:  signTransaction,
		Context: ctx,
	}, action, scalar)
	if err != nil {
		return nil, fmt.Errorf("failed to call signal error: %v action: %v scalar: %v", err, action, scalar)
	}

	_, err = bind.WaitMined(ctx, s.Client, txn)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	receipt, err := s.Client.TransactionReceipt(ctx, txn.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed")
	}

	res := &SignalResponse{
		TransactionHash: txn.Hash().Hex(),
		Sender:          senderAddress,
	}

	for _, log := range receipt.Logs {
		eventLog := EventLog{
			Data: "0x" + strings.Trim(fmt.Sprintf("%x", log.Data), "\""),
		}

		// Extract the event topics
		for _, topic := range log.Topics {
			eventLog.Event = append(eventLog.Event, topic.Hex())
		}

		// Append the eventLog to the response
		res.EventLog = append(res.EventLog, eventLog)
	}

	return res, nil
}

func (s *Service) SignalMany(ctx context.Context, actions []uint8, scalars []*big.Int, nonce *big.Int) (*SignalResponse, error) {
	var res = &SignalResponse{}

	if len(actions) != len(scalars) {
		return nil, fmt.Errorf("actions and scalars arrays must be of the same length")
	}

	contractAddress := common.HexToAddress(contractAddress)
	metamodelInstance, err := metamodel.NewMetamodel(contractAddress, s.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a Metamodel contract: %v", err)
	}

	_ = nonce // REVIEW: can we use this nonce?
	txn, err := metamodelInstance.SignalMany(&bind.TransactOpts{
		From:    common.HexToAddress(senderAddress),
		Nonce:   nil,
		Signer:  signTransaction,
		Context: ctx,
	}, actions, scalars)

	if err != nil {
		return nil, fmt.Errorf("failed to call model method: %v", err)
	}

	receipt, err := s.Client.TransactionReceipt(ctx, txn.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	res.TransactionHash = txn.Hash().Hex()
	res.Sender = senderAddress
	res.Contract = contractAddress.Hex()

	for _, log := range receipt.Logs {
		eventLog := EventLog{
			Data: "0x" + strings.Trim(fmt.Sprintf("%x", log.Data), "\""),
		}

		for _, topic := range log.Topics {
			eventLog.Event = append(eventLog.Event, topic.Hex())
		}

		res.EventLog = append(res.EventLog, eventLog)
	}

	return res, nil
}

type SignalApiResponse struct {
	Response *SignalResponse `json:"response,omitempty"`
	Nonce    int64           `json:"nonce"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Nonce int64  `json:"nonce,omitempty"`
}

func (s *Service) SignalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues := r.URL.Query()
	actionsStr := queryValues.Get("action") // "1,2,3"
	scalarsStr := queryValues.Get("scalar") // "100,200,300"
	nonceStr := queryValues.Get("nonce")    // "123"

	if nonceStr == "" {
		nonceStr = "0"
	}

	actionsStrSlice := strings.Split(actionsStr, ",")
	scalarsStrSlice := strings.Split(scalarsStr, ",")

	var actions []uint8
	var scalars []*big.Int

	if len(actionsStrSlice) != len(scalarsStrSlice) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Actions and scalars arrays must be of the same length", Nonce: 0})
		return
	}

	for _, a := range actionsStrSlice {
		action, err := strconv.ParseUint(a, 10, 8)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid action value", Nonce: 0})
			return
		}
		actions = append(actions, uint8(action))
	}

	for _, s := range scalarsStrSlice {
		scalar, ok := big.NewInt(0).SetString(s, 10)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid scalar value", Nonce: 0})
			return
		}
		scalars = append(scalars, scalar)
	}

	nonce, _ := big.NewInt(0).SetString(nonceStr, 10)

	res, err := s.SignalMany(r.Context(), actions, scalars, nonce)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error(), Nonce: nonce.Int64()})
		return
	}

	response := SignalApiResponse{
		Nonce:    nonce.Int64(),
		Response: res,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
