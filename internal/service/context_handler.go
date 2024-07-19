package service

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pflow-dev/pflow-eth/metamodel"
	"math/big"
	"net/http"
)

type ModelApiResponse struct {
	Nonce   int64       `json:"nonce"`
	Context interface{} `json:"context,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (s Service) ContextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues := r.URL.Query()
	addr := queryValues.Get("addr")
	nonce := queryValues.Get("nonce")
	if nonce == "" {
		nonce = "0"
	}

	response := ModelApiResponse{}
	nonceInt, _ := new(big.Int).SetString(nonce, 10)
	response.Nonce = nonceInt.Int64()

	if addr == "" {
		response.Error = "Invalid address"
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	address := common.HexToAddress(addr)
	metamodelInstance, err := metamodel.NewMetamodel(address, s.Client)
	if err != nil {
		response.Error = "Failed to instantiate a Metamodel contract"
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	modelData, err := metamodelInstance.Context(nil)
	if err != nil {
		response.Error = "Failed to call model method"
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Context = modelData

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
