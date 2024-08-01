package db

import (
	"time"
)

type Node struct {
	ID        int32
	Name      string
	ChainID   string `db:"starknet_chain_id"`
	URL       string
	APIKey    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
