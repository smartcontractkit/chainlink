// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

// GetClusterNodes returns information about all the nodes participating in the cluster.
func (cl *Client) GetClusterNodes(ctx context.Context) (out []*GetClusterNodesResult, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getClusterNodes", nil)
	return
}

type GetClusterNodesResult struct {
	// Node public key.
	Pubkey solana.PublicKey `json:"pubkey"`

	// TODO: "" or nil ?

	// Gossip network address for the node.
	Gossip *string `json:"gossip,omitempty"`

	// TPU network address for the node.
	TPU *string `json:"tpu,omitempty"`

	// JSON RPC network address for the node, or empty if the JSON RPC service is not enabled.
	RPC *string `json:"rpc,omitempty"`

	// The software version of the node, or empty if the version information is not available.
	Version *string `json:"version,omitempty"`

	// TODO: what type is this?
	// The unique identifier of the node's feature set.
	FeatureSet uint32 `json:"featureSet,omitempty"`

	// The shred version the node has been configured to use.
	ShredVersion uint16 `json:"shredVersion,omitempty"`
}
