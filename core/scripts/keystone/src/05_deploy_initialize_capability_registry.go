package src

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"google.golang.org/protobuf/proto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
)

type peer struct {
	PeerID string
	Signer string
}

var (
	workflowDonPeers = []peer{
		{
			PeerID: "12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N",
			Signer: "0x9639dCc7D0ca4468B5f684ef89F12F0B365c9F6d",
		},
		{
			PeerID: "12D3KooWG1AyvwmCpZ93J8pBQUE1SuzrjDXnT4BeouncHR3jWLCG",
			Signer: "0x8f0fAE64f5f75067833ed5deDC2804B62b21383d",
		},
		{
			PeerID: "12D3KooWGeUKZBRMbx27FUTgBwZa9Ap9Ym92mywwpuqkEtz8XWyv",
			Signer: "0xf09A863D920840c13277e76F43CFBdfB22b8FB7C",
		},
		{
			PeerID: "12D3KooW9zYWQv3STmDeNDidyzxsJSTxoCTLicafgfeEz9nhwhC4",
			Signer: "0x7eD90b519bC3054a575C464dBf39946b53Ff90EF",
		},
		{
			PeerID: "12D3KooWG1AeBnSJH2mdcDusXQVye2jqodZ6pftTH98HH6xvrE97",
			Signer: "0x8F572978673d711b2F061EB7d514BD46EAD6668A",
		},
		{
			PeerID: "12D3KooWBf3PrkhNoPEmp7iV291YnPuuTsgEDHTscLajxoDvwHGA",
			Signer: "0x21eF07Dfaf8f7C10CB0d53D18b641ee690541f9D",
		},
		{
			PeerID: "12D3KooWP3FrMTFXXRU2tBC8aYvEBgUX6qhcH9q2JZCUi9Wvc2GX",
			Signer: "0x7Fa21F6f716CFaF8f249564D72Ce727253186C89",
		},
	}
	triggerDonPeers = []peer{
		{
			PeerID: "12D3KooWBaiTbbRwwt2fbNifiL7Ew9tn3vds9AJE3Nf3eaVBX36m",
			Signer: "0x9CcE7293a4Cc2621b61193135A95928735e4795F",
		},
		{
			PeerID: "12D3KooWS7JSY9fzSfWgbCE1S3W2LNY6ZVpRuun74moVBkKj6utE",
			Signer: "0x3c775F20bCB2108C1A818741Ce332Bb5fe0dB925",
		},
		{
			PeerID: "12D3KooWMMTDXcWhpVnwrdAer1jnVARTmnr3RyT3v7Djg8ZuoBh9",
			Signer: "0x50314239e2CF05555ceeD53E7F47eB2A8Eab0dbB",
		},
		{
			PeerID: "12D3KooWGzVXsKxXsF4zLgxSDM8Gzx1ywq2pZef4PrHMKuVg4K3P",
			Signer: "0xd76A4f98898c3b9A72b244476d7337b50D54BCd8",
		},
		{
			PeerID: "12D3KooWSyjmmzjVtCzwN7bXzZQFmWiJRuVcKBerNjVgL7HdLJBW",
			Signer: "0x656A873f6895b8a03Fb112dE927d43FA54B2c92A",
		},
		{
			PeerID: "12D3KooWLGz9gzhrNsvyM6XnXS3JRkZoQdEzuAvysovnSChNK5ZK",
			Signer: "0x5d1e87d87bF2e0cD4Ea64F381a2dbF45e5f0a553",
		},
		{
			PeerID: "12D3KooWAvZnvknFAfSiUYjATyhzEJLTeKvAzpcLELHi4ogM3GET",
			Signer: "0x91d9b0062265514f012Eb8fABA59372fD9520f56",
		},
	}
	targetDonPeers = []peer{
		{
			PeerID: "12D3KooWJrthXtnPHw7xyHFAxo6NxifYTvc8igKYaA6wRRRqtsMb",
			Signer: "0x3F82750353Ea7a051ec9bA011BC628284f9a5327",
		},
		{
			PeerID: "12D3KooWFQekP9sGex4XhqEJav5EScjTpDVtDqJFg1JvrePBCEGJ",
			Signer: "0xc23545876A208AA0443B1b8d552c7be4FF4b53F0",
		},
		{
			PeerID: "12D3KooWFLEq4hYtdyKWwe47dXGEbSiHMZhmr5xLSJNhpfiEz8NF",
			Signer: "0x82601Fa43d8B1dC1d4eB640451aC86a7CDA37011",
		},
		{
			PeerID: "12D3KooWN2hztiXNNS1jMQTTvvPRYcarK1C7T3Mdqk4x4gwyo5WS",
			Signer: "0x1a684B3d8f917fe496b7B1A8b29EDDAED64F649f",
		},
	}
)

type deployAndInitializeCapabilityRegistryCommand struct{}

func NewDeployAndInitializeCapabilityRegistryCommand() *deployAndInitializeCapabilityRegistryCommand {
	return &deployAndInitializeCapabilityRegistryCommand{}
}

func (c *deployAndInitializeCapabilityRegistryCommand) Name() string {
	return "deploy-and-initialize-capability-registry"
}

func peerIDToB(peerID string) ([32]byte, error) {
	var peerIDB ragetypes.PeerID
	err := peerIDB.UnmarshalText([]byte(peerID))
	if err != nil {
		return [32]byte{}, err
	}

	return peerIDB, nil
}

func peers(ps []peer) ([][32]byte, error) {
	out := [][32]byte{}
	for _, p := range ps {
		b, err := peerIDToB(p.PeerID)
		if err != nil {
			return nil, err
		}

		out = append(out, b)
	}

	return out, nil
}

func peerToNode(nopID uint32, p peer) (kcr.CapabilityRegistryNodeInfo, error) {
	peerIDB, err := peerIDToB(p.PeerID)
	if err != nil {
		return kcr.CapabilityRegistryNodeInfo{}, fmt.Errorf("failed to convert peerID: %w", err)
	}

	sig := strings.TrimPrefix(p.Signer, "0x")
	signerB, err := hex.DecodeString(sig)
	if err != nil {
		return kcr.CapabilityRegistryNodeInfo{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return kcr.CapabilityRegistryNodeInfo{
		NodeOperatorId: nopID,
		P2pId:          [32]byte(peerIDB),
		Signer:         sigb,
	}, nil
}

// Run expects the follow environment variables to be set:
//
//  1. Deploys the CapabilityRegistry contract
//  2. Configures it with a hardcode DON setup, as used by our staging environment.
func (c *deployAndInitializeCapabilityRegistryCommand) Run(args []string) {
	ctx := context.Background()

	fs := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	// create flags for all of the env vars then set the env vars to normalize the interface
	// this is a bit of a hack but it's the easiest way to make this work
	ethUrl := fs.String("ethurl", "", "URL of the Ethereum node")
	chainID := fs.Int64("chainid", 11155111, "chain ID of the Ethereum network to deploy to")
	accountKey := fs.String("accountkey", "", "private key of the account to deploy from")
	capabilityRegistryAddress := fs.String("craddress", "", "address of the capability registry")

	err := fs.Parse(args)
	if err != nil ||
		*ethUrl == "" || ethUrl == nil ||
		*chainID == 0 || chainID == nil ||
		*accountKey == "" || accountKey == nil {
		fs.Usage()
		os.Exit(1)
	}

	os.Setenv("ETH_URL", *ethUrl)
	os.Setenv("ETH_CHAIN_ID", fmt.Sprintf("%d", *chainID))
	os.Setenv("ACCOUNT_KEY", *accountKey)

	env := helpers.SetupEnv(false)

	var reg *kcr.CapabilityRegistry
	if *capabilityRegistryAddress == "" {
		reg = deployCapabilityRegistry(env)
	} else {
		addr := common.HexToAddress(*capabilityRegistryAddress)
		r, err := kcr.NewCapabilityRegistry(addr, env.Ec)
		if err != nil {
			panic(err)
		}

		reg = r
	}

	streamsTrigger := kcr.CapabilityRegistryCapability{
		LabelledName:   "streams-trigger",
		Version:        "1.0.0",
		CapabilityType: uint8(0), // trigger
	}
	sid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, streamsTrigger.LabelledName, streamsTrigger.Version)
	if err != nil {
		panic(err)
	}

	writeChain := kcr.CapabilityRegistryCapability{
		LabelledName:   "write_ethereum-testnet-sepolia",
		Version:        "1.0.0",
		CapabilityType: uint8(3), // target
	}
	wid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChain.LabelledName, writeChain.Version)
	if err != nil {
		log.Printf("failed to call GetHashedCapabilityId: %w", err)
	}

	ocr := kcr.CapabilityRegistryCapability{
		LabelledName:   "offchain_reporting",
		Version:        "1.0.0",
		CapabilityType: uint8(3), // target
	}
	ocrid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, ocr.LabelledName, ocr.Version)
	if err != nil {
		log.Printf("failed to call GetHashedCapabilityId: %w", err)
	}

	tx, err := reg.AddCapabilities(env.Owner, []kcr.CapabilityRegistryCapability{
		streamsTrigger,
		writeChain,
		ocr,
	})
	if err != nil {
		log.Printf("failed to call AddCapabilities: %w", err)
	}

	helpers.ConfirmTXMined(ctx, env.Ec, tx, env.ChainID)

	tx, err = reg.AddNodeOperators(env.Owner, []kcr.CapabilityRegistryNodeOperator{
		{
			Admin: env.Owner.From,
			Name:  "STAGING_NODE_OPERATOR",
		},
	})
	if err != nil {
		log.Printf("failed to AddNodeOperators: %w", err)
	}

	receipt := helpers.ConfirmTXMined(ctx, env.Ec, tx, env.ChainID)

	recLog, err := reg.ParseNodeOperatorAdded(*receipt.Logs[0])
	if err != nil {
		panic(err)
	}

	nopID := recLog.NodeOperatorId
	nodes := []kcr.CapabilityRegistryNodeInfo{}
	for _, wfPeer := range workflowDonPeers {
		n, err := peerToNode(nopID, wfPeer)
		if err != nil {
			panic(err)
		}

		n.HashedCapabilityIds = [][32]byte{ocrid}
		nodes = append(nodes, n)
	}

	for _, triggerPeer := range triggerDonPeers {
		n, err := peerToNode(nopID, triggerPeer)
		if err != nil {
			panic(err)
		}

		n.HashedCapabilityIds = [][32]byte{sid}
		nodes = append(nodes, n)
	}

	for _, targetPeer := range targetDonPeers {
		n, err := peerToNode(nopID, targetPeer)
		if err != nil {
			panic(err)
		}

		n.HashedCapabilityIds = [][32]byte{wid}
		nodes = append(nodes, n)
	}

	tx, err = reg.AddNodes(env.Owner, nodes)
	if err != nil {
		log.Printf("failed to AddNodes: %w", err)
	}

	helpers.ConfirmTXMined(ctx, env.Ec, tx, env.ChainID)

	// workflow DON
	ps, err := peers(workflowDonPeers)
	if err != nil {
		panic(err)
	}

	cfgs := []kcr.CapabilityRegistryCapabilityConfiguration{
		{
			CapabilityId: ocrid,
		},
	}
	tx, err = reg.AddDON(env.Owner, ps, cfgs, false, true, 2)
	if err != nil {
		log.Printf("workflowDON: failed to AddDON: %w", err)
	}

	// trigger DON
	ps, err = peers(triggerDonPeers)
	if err != nil {
		panic(err)
	}

	config := &remotetypes.RemoteTriggerConfig{
		RegistrationRefreshMs: 20000,
		RegistrationExpiryMs:  60000,
		// F + 1
		MinResponsesToAggregate: uint32(1) + 1,
	}
	configb, err := proto.Marshal(config)
	if err != nil {
		panic(err)
	}
	cfgs = []kcr.CapabilityRegistryCapabilityConfiguration{
		{
			CapabilityId: sid,
			Config:       configb,
		},
	}
	tx, err = reg.AddDON(env.Owner, ps, cfgs, true, false, 1)
	if err != nil {
		log.Printf("triggerDON: failed to AddDON: %w", err)
	}

	// target DON
	ps, err = peers(targetDonPeers)
	if err != nil {
		panic(err)
	}

	cfgs = []kcr.CapabilityRegistryCapabilityConfiguration{
		{
			CapabilityId: wid,
		},
	}
	tx, err = reg.AddDON(env.Owner, ps, cfgs, true, false, 1)
	if err != nil {
		log.Printf("targetDON: failed to AddDON: %w", err)
	}
}

func deployCapabilityRegistry(env helpers.Environment) *kcr.CapabilityRegistry {
	_, tx, contract, err := kcr.DeployCapabilityRegistry(env.Owner, env.Ec)
	if err != nil {
		panic(err)
	}

	addr := helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
	fmt.Printf("CapabilityRegistry address: %s", addr)
	return contract
}
