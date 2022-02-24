package main

import (
	"os"

	"example.com/m/PBFT/network"
)

func main() {

	nodeID := os.Args[1]
	server := network.NewServer(nodeID)
	server.Start()

}
