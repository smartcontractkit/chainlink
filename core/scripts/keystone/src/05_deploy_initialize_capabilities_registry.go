package src

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
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

	aptosTargetDonPeers = []peer{
		{
			PeerID: "12D3KooWNBr1AD3vD3dzSLgg1tK56qyJoenDx7EYNnZpbr1g4jD6",
			Signer: "a41f9a561ff2266d94240996a76f9c2b3b7d8184",
		},
		{
			PeerID: "12D3KooWRRgWiZGw5GYsPa62CkwFNKJb5u4hWo4DinnvjG6GE6Nj",
			Signer: "e4f3c7204776530fb7833db6f9dbfdb8bd0ec96892965324a71c20d6776f67f0",
		},
		{
			PeerID: "12D3KooWKwzgUHw5YbqUsYUVt3yiLSJcqc8ANofUtqHX6qTm7ox2",
			Signer: "4071ea00e2e2c76b3406018ba9f66bf6b9aee3a6762e62ac823b1ee91ba7d7b0",
		},
		{
			PeerID: "12D3KooWBRux5o2bw1j3SQwEHzCspjkt7Xe3Y3agRUuab2SUnExj",
			Signer: "6f5180c7d276876dbe413bf9b0efff7301d1367f39f4bac64180090cab70989b",
		},
		{
			PeerID: "12D3KooWFqvDaMSDGa6eMSTF9en6G2c3ZbGLmaA5Xs3AgxVBPb8B",
			Signer: "dbce9a6df8a04d54e52a109d01ee9b5d32873b1d2436cf7b7fae61fd6eca46f8",
		},
	}
)

type deployAndInitializeCapabilitiesRegistryCommand struct{}

func NewDeployAndInitializeCapabilitiesRegistryCommand() *deployAndInitializeCapabilitiesRegistryCommand {
	return &deployAndInitializeCapabilitiesRegistryCommand{}
}

func (c *deployAndInitializeCapabilitiesRegistryCommand) Name() string {
	return "deploy-and-initialize-capabilities-registry"
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

func peerToNode(nopID uint32, p peer) (kcr.CapabilitiesRegistryNodeParams, error) {
	peerIDB, err := peerIDToB(p.PeerID)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert peerID: %w", err)
	}

	sig := strings.TrimPrefix(p.Signer, "0x")
	signerB, err := hex.DecodeString(sig)
	if err != nil {
		return kcr.CapabilitiesRegistryNodeParams{}, fmt.Errorf("failed to convert signer: %w", err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return kcr.CapabilitiesRegistryNodeParams{
		NodeOperatorId: nopID,
		P2pId:          peerIDB,
		Signer:         sigb,
	}, nil
}

func newCapabilityConfig() *capabilitiespb.CapabilityConfig {
	return &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
	}
}

// Run expects the following environment variables to be set:
//
//  1. Deploys the CapabilitiesRegistry contract
//  2. Configures it with a hardcode DON setup, as used by our staging environment.
func (c *deployAndInitializeCapabilitiesRegistryCommand) Run(args []string) {
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

	var reg *kcr.CapabilitiesRegistry
	if *capabilityRegistryAddress == "" {
		reg = deployCapabilitiesRegistry(env)
	} else {
		addr := common.HexToAddress(*capabilityRegistryAddress)
		r, innerErr := kcr.NewCapabilitiesRegistry(addr, env.Ec)
		if err != nil {
			panic(innerErr)
		}

		reg = r
	}

	streamsTrigger := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "streams-trigger",
		Version:        "1.0.0",
		CapabilityType: uint8(0), // trigger
	}
	sid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, streamsTrigger.LabelledName, streamsTrigger.Version)
	if err != nil {
		panic(err)
	}

	writeChain := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_ethereum-testnet-sepolia",
		Version:        "1.0.0",
		CapabilityType: uint8(3), // target
	}
	wid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, writeChain.LabelledName, writeChain.Version)
	if err != nil {
		log.Printf("failed to call GetHashedCapabilityId: %s", err)
	}

	aptosWriteChain := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "write_aptos",
		Version:        "1.0.0",
		CapabilityType: uint8(3), // target
	}
	awid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, aptosWriteChain.LabelledName, aptosWriteChain.Version)
	if err != nil {
		log.Printf("failed to call GetHashedCapabilityId: %s", err)
	}

	ocr := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "offchain_reporting",
		Version:        "1.0.0",
		CapabilityType: uint8(2), // consensus
	}
	ocrid, err := reg.GetHashedCapabilityId(&bind.CallOpts{}, ocr.LabelledName, ocr.Version)
	if err != nil {
		log.Printf("failed to call GetHashedCapabilityId: %s", err)
	}

	tx, err := reg.AddCapabilities(env.Owner, []kcr.CapabilitiesRegistryCapability{
		streamsTrigger,
		writeChain,
		aptosWriteChain,
		ocr,
	})
	if err != nil {
		log.Printf("failed to call AddCapabilities: %s", err)
	}

	helpers.ConfirmTXMined(ctx, env.Ec, tx, env.ChainID)

	tx, err = reg.AddNodeOperators(env.Owner, []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: env.Owner.From,
			Name:  "STAGING_NODE_OPERATOR",
		},
	})
	if err != nil {
		log.Printf("failed to AddNodeOperators: %s", err)
	}

	receipt := helpers.ConfirmTXMined(ctx, env.Ec, tx, env.ChainID)

	recLog, err := reg.ParseNodeOperatorAdded(*receipt.Logs[0])
	if err != nil {
		panic(err)
	}

	nopID := recLog.NodeOperatorId
	nodes := []kcr.CapabilitiesRegistryNodeParams{}
	for _, wfPeer := range workflowDonPeers {
		n, innerErr := peerToNode(nopID, wfPeer)
		if innerErr != nil {
			panic(innerErr)
		}

		n.HashedCapabilityIds = [][32]byte{ocrid}
		nodes = append(nodes, n)
	}

	for _, triggerPeer := range triggerDonPeers {
		n, innerErr := peerToNode(nopID, triggerPeer)
		if innerErr != nil {
			panic(innerErr)
		}

		n.HashedCapabilityIds = [][32]byte{sid}
		nodes = append(nodes, n)
	}

	for _, targetPeer := range targetDonPeers {
		n, innerErr := peerToNode(nopID, targetPeer)
		if innerErr != nil {
			panic(innerErr)
		}

		n.HashedCapabilityIds = [][32]byte{wid}
		nodes = append(nodes, n)
	}

	for _, targetPeer := range aptosTargetDonPeers {
		n, innerErr := peerToNode(nopID, targetPeer)
		if innerErr != nil {
			panic(innerErr)
		}

		n.HashedCapabilityIds = [][32]byte{awid}
		nodes = append(nodes, n)
	}

	tx, err = reg.AddNodes(env.Owner, nodes)
	if err != nil {
		log.Printf("failed to AddNodes: %s", err)
	}

	helpers.ConfirmTXMined(ctx, env.Ec, tx, env.ChainID)

	// workflow DON
	ps, err := peers(workflowDonPeers)
	if err != nil {
		panic(err)
	}

	cc := newCapabilityConfig()
	ccb, err := proto.Marshal(cc)
	if err != nil {
		panic(err)
	}

	cfgs := []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: ocrid,
			Config:       ccb,
		},
	}
	_, err = reg.AddDON(env.Owner, ps, cfgs, true, true, 2)
	if err != nil {
		log.Printf("workflowDON: failed to AddDON: %s", err)
	}

	// trigger DON
	ps, err = peers(triggerDonPeers)
	if err != nil {
		panic(err)
	}

	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				// F + 1
				MinResponsesToAggregate: uint32(1) + 1,
			},
		},
	}
	configb, err := proto.Marshal(config)
	if err != nil {
		panic(err)
	}
	cfgs = []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: sid,
			Config:       configb,
		},
	}
	_, err = reg.AddDON(env.Owner, ps, cfgs, true, false, 1)
	if err != nil {
		log.Printf("triggerDON: failed to AddDON: %s", err)
	}

	// target DON
	ps, err = peers(targetDonPeers)
	if err != nil {
		panic(err)
	}

	targetCapabilityConfig := newCapabilityConfig()
	targetCapabilityConfig.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteTargetConfig{
		RemoteTargetConfig: &capabilitiespb.RemoteTargetConfig{
			RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
		},
	}

	remoteTargetConfigBytes, err := proto.Marshal(targetCapabilityConfig)
	if err != nil {
		panic(err)
	}

	cfgs = []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: wid,
			Config:       remoteTargetConfigBytes,
		},
	}
	_, err = reg.AddDON(env.Owner, ps, cfgs, true, false, 1)
	if err != nil {
		log.Printf("targetDON: failed to AddDON: %s", err)
	}

	// Aptos target DON
	ps, err = peers(aptosTargetDonPeers)
	if err != nil {
		panic(err)
	}

	targetCapabilityConfig = newCapabilityConfig()
	targetCapabilityConfig.RemoteConfig = &capabilitiespb.CapabilityConfig_RemoteTargetConfig{
		RemoteTargetConfig: &capabilitiespb.RemoteTargetConfig{
			RequestHashExcludedAttributes: []string{"signed_report.Signatures"},
		},
	}

	remoteTargetConfigBytes, err = proto.Marshal(targetCapabilityConfig)
	if err != nil {
		panic(err)
	}

	cfgs = []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: awid,
			Config:       remoteTargetConfigBytes,
		},
	}
	_, err = reg.AddDON(env.Owner, ps, cfgs, true, false, 1)
	if err != nil {
		log.Printf("targetDON: failed to AddDON: %s", err)
	}
}

func deployCapabilitiesRegistry(env helpers.Environment) *kcr.CapabilitiesRegistry {
	_, tx, contract, err := kcr.DeployCapabilitiesRegistry(env.Owner, env.Ec)
	if err != nil {
		panic(err)
	}

	addr := helpers.ConfirmContractDeployed(context.Background(), env.Ec, tx, env.ChainID)
	fmt.Printf("CapabilitiesRegistry address: %s", addr)
	return contract
}
