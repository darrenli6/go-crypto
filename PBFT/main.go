package main

import (
	"os"

	"goland/go-crypto/PBFT/network"
)

func main() {

	nodeID := os.Args[1]
	server := network.NewServer(nodeID)
	server.Start()

}
