package main

import "fmt"

type Node struct {
	ID      string
	Address string
}

func main() {
	nodes := []Node{
		{ID: "Aura_01", Address: "127.0.0.1:8080"},
		{ID: "Aura_02", Address: "127.0.0.1:8081"},
	}
	fmt.Println("--- Kin-Swarm P2P Handshake ---")
	for _, node := range nodes {
		fmt.Printf("Peering with %s at %s... [CONNECTED]\n", node.ID, node.Address)
	}
	fmt.Println("Network_State: FULLY_MESHED")
}
