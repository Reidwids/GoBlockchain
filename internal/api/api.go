package api

import (
	"GoBlockchain/internal/blockchain"
	"GoBlockchain/internal/blockchain/pb"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Start() {
	fmt.Println("Starting Go Blockchain...")
	bc := blockchain.NewBlockchain()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/mineBlock", withBlockchain(handleMineBlock, bc))
	// http.HandleFunc("/getChain", handleGetChain)
	// http.HandleFunc("/isChainValid", handleIsChainValid)
	// http.HandleFunc("/addTransaction", handleAddTransaction)
	// http.HandleFunc("/connectNode", handleConnectNode)
	// http.HandleFunc("/replaceChain", handleReplaceChain)
	http.ListenAndServe(":3003", nil)
}

func withBlockchain(next http.HandlerFunc, bc *pb.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add 'bc' to the request context
		ctx := context.WithValue(r.Context(), "blockchain", bc)
		next(w, r.WithContext(ctx))
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Message: "Golang Blockchain",
		Version: "v1.3.0",
	}
	json.NewEncoder(w).Encode(response)
}

func handleMineBlock(w http.ResponseWriter, r *http.Request) {
	// bc := r.Context().Value("blockchain").(*pb.Blockchain)
	// previousBlock :=
}
