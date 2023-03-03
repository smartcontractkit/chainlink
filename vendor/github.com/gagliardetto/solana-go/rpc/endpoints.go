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

// See more: https://docs.solana.com/cluster/rpc-endpoints

const (
	protocolHTTPS = "https://"
	protocolWSS   = "wss://"
)

type Cluster struct {
	Name string
	RPC  string
	WS   string
}

var (
	MainNetBeta = Cluster{
		Name: "mainnet-beta",
		RPC:  MainNetBeta_RPC,
		WS:   MainNetBeta_WS,
	}
	TestNet = Cluster{
		Name: "testnet",
		RPC:  TestNet_RPC,
		WS:   TestNet_WS,
	}
	DevNet = Cluster{
		Name: "devnet",
		RPC:  DevNet_RPC,
		WS:   DevNet_WS,
	}
	LocalNet = Cluster{
		Name: "localnet",
		RPC:  LocalNet_RPC,
		WS:   LocalNet_WS,
	}
)

const (
	hostDevNet           = "api.devnet.solana.com"
	hostTestNet          = "api.testnet.solana.com"
	hostMainNetBeta      = "api.mainnet-beta.solana.com"
	hostMainNetBetaSerum = "solana-api.projectserum.com"
)

const (
	DevNet_RPC           = protocolHTTPS + hostDevNet
	TestNet_RPC          = protocolHTTPS + hostTestNet
	MainNetBeta_RPC      = protocolHTTPS + hostMainNetBeta
	MainNetBetaSerum_RPC = protocolHTTPS + hostMainNetBetaSerum
	LocalNet_RPC         = "http://127.0.0.1:8899"
)

const (
	DevNet_WS           = protocolWSS + hostDevNet
	TestNet_WS          = protocolWSS + hostTestNet
	MainNetBeta_WS      = protocolWSS + hostMainNetBeta
	MainNetBetaSerum_WS = protocolWSS + hostMainNetBetaSerum
	LocalNet_WS         = "ws://127.0.0.1:8900"
)
