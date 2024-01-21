package node

import (
	"GoBlockchain/internal/blockchain"
	"GoBlockchain/internal/proto/pb"
	"context"
	"fmt"
)

type Server struct {
	pb.UnimplementedNodeServiceServer
	Blockchain *blockchain.Blockchain
}

func (s *Server) GetNodes(ctx context.Context, req *pb.GetNodesReq) (*pb.NodeList, error) {
	return &pb.NodeList{Nodes: s.Blockchain.Nodes}, nil
}

// broadcast block
// broadcast transaction
// broadcast node
// broadcast chain
// consensus
// Peer discovery
// fork resolution

// Func to query hardcoded nodes for their nodelist
// Func to add new node to nodelist

func AddNode(blockchain *blockchain.Blockchain, address string) {
	blockchain.Nodes = append(blockchain.Nodes, address)
	fmt.Println("Node successfully added!")
}

// func validateNodes() bool {}

// func broadcastNodes(blockchain *pb.Blockchain) {
// 	for _, node := range blockchain.Nodes {
// 		url := fmt.Sprintf("http://%s/getNodes", node)
// 		response, err :=  http.Get(url)
// 		struct NodeRes{
// 			Nodes string `json:"nodes"`
// 		}
// 		json.NewEncoder(w).Encode(response)
// 		res.Body.nodes
// 	}
// 	// for each node in node list, get nodes and calculate differences.
// 	// If there are new nodes in our list for them, add the new node and re broadcast
// 	// If there are new nodes for us, add the new nodes and re broadcast
// }
