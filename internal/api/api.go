package api

import (
	"GoBlockchain/internal/blockchain"
	"GoBlockchain/internal/blockchain/pb"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func Start() {
	fmt.Println("Starting Go Blockchain...")
	bc := blockchain.NewBlockchain()
	nodeId := uuid.New()
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/mineBlock", withBlockchain(handleMineBlock, bc, nodeId))
	// http.HandleFunc("/getChain", handleGetChain)
	// http.HandleFunc("/isChainValid", handleIsChainValid)
	// http.HandleFunc("/addTransaction", handleAddTransaction)
	// http.HandleFunc("/connectNode", handleConnectNode)
	// http.HandleFunc("/replaceChain", handleReplaceChain)
	http.ListenAndServe(":3003", nil)
}

func withBlockchain(next http.HandlerFunc, bc *pb.Blockchain, nodeId uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add 'bc' to the request context
		ctx := context.WithValue(r.Context(), "blockchain", bc)
		ctx = context.WithValue(ctx, "nodeAddress", nodeId)
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
	bc := r.Context().Value("blockchain").(*pb.Blockchain)
	nodeId := r.Context().Value("nodeId").(uuid.UUID)
	prevBlock := blockchain.GetPreviousBlock(bc)
	prevProof := prevBlock.Proof
	previousHash := blockchain.HashBlock(prevBlock)
	newProof := blockchain.ProofOfWork(prevProof)
	blockchain.AddTransaction(bc, "The Network", nodeId.String(), 1)
	block := blockchain.CreateBlock(bc, newProof, previousHash)
	response := struct {
		Message   string `json:"message"`
		Index     int64  `json:"index"`
		Timestamp int64  `json:"timestamp"`
		PrevHash  []byte `json:"prevHash"`
	}{
		Message:   "Block mined",
		Index:     block.Index,
		Timestamp: block.Timestamp,
		PrevHash:  block.PrevHash,
	}
	json.NewEncoder(w).Encode(response)
}
