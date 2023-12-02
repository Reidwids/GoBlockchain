package api

import (
	"GoBlockchain/internal/blockchain/pb"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func Start() {
	fmt.Println("Starting Go Blockchain...")
	http.HandleFunc("/", handleRoot)
	// http.HandleFunc("/download", handleDownload)
	http.ListenAndServe(":3003", nil)
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

func GetChainFromNode(node string) ([]*pb.Block, error) {
	url := fmt.Sprintf("http://%s/getChain", node)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error Connecting with node:", node, err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	// Decode the protobuf-encoded data
	var chainBuf pb.Chain
	err = proto.Unmarshal(body, &chainBuf)
	if err != nil {
		fmt.Println("Error decoding protobuf data:", err)
		return nil, err
	}

	receivedChain := chainBuf.Chain

	return receivedChain, nil
}
