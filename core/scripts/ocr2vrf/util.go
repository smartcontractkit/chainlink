package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	"github.com/smartcontractkit/ocr2vrf/gethwrappers/vrf"
)

func DeployDKG(e helpers.Environment) (blockhashStoreAddress common.Address) {
	_, tx, _, err := dkg.DeployDKG(e.Owner, e.Ec)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func AddClientToDKG(e helpers.Environment, dkgAddress string, keyId string, clientAddress string) {
	var keyIdBytes [32]byte
	copy(keyIdBytes[:], keyId)
	fmt.Printf("Encoded Key ID: 0x%x \n", keyIdBytes)

	dkg, err := dkg.NewDKG(common.HexToAddress(dkgAddress), e.Ec)
	helpers.PanicErr(err)

	tx, err := dkg.AddClient(e.Owner, keyIdBytes, common.HexToAddress(clientAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func RemoveClientFromDKG(e helpers.Environment, dkgAddress string, keyId string, clientAddress string) {
	var keyIdBytes [32]byte
	copy(keyIdBytes[:], keyId)
	fmt.Printf("Key ID in bytes: 0x%x \n", keyIdBytes)

	dkg, err := dkg.NewDKG(common.HexToAddress(dkgAddress), e.Ec)
	helpers.PanicErr(err)

	tx, err := dkg.RemoveClient(e.Owner, keyIdBytes, common.HexToAddress(clientAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func SetDKGConfig(e helpers.Environment, dkgAddress string, signers string, transmitters string, f uint8, onchainConfig string, offchainConfig string, offchainConfigVersion uint64) {
	dkg := newDKG(common.HexToAddress(dkgAddress), e.Ec)
	decodedOnchainConfig, err := hex.DecodeString(onchainConfig)
	helpers.PanicErr(err)
	decodedOffchainConfig, err := hex.DecodeString(offchainConfig)
	helpers.PanicErr(err)
	dkg.SetConfig(e.Owner, helpers.ParseAddressSlice(signers), helpers.ParseAddressSlice(transmitters), f, decodedOnchainConfig, offchainConfigVersion, decodedOffchainConfig)
}

func SetVRFConfig(e helpers.Environment, vrfAddress string, signers string, transmitters string, f uint8, onchainConfig string, offchainConfig string, offchainConfigVersion uint64) {
	vrf := newVRF(common.HexToAddress(vrfAddress), e.Ec)
	decodedOnchainConfig, err := hex.DecodeString(onchainConfig)
	helpers.PanicErr(err)
	decodedOffchainConfig, err := hex.DecodeString(offchainConfig)
	helpers.PanicErr(err)
	vrf.SetConfig(e.Owner, helpers.ParseAddressSlice(signers), helpers.ParseAddressSlice(transmitters), f, decodedOnchainConfig, offchainConfigVersion, decodedOffchainConfig)
}

func DeployVRF(e helpers.Environment, dkgAddress string, keyId string) (blockhashStoreAddress common.Address) {
	var keyIdBytes [32]byte
	copy(keyIdBytes[:], keyId)
	fmt.Printf("Key ID in bytes: 0x%x \n", keyIdBytes)

	_, tx, _, err := vrf.DeployVRF(e.Owner, e.Ec, common.HexToAddress(dkgAddress), keyIdBytes)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func newVRF(addr common.Address, client *ethclient.Client) *vrf.VRF {
	vrf, err := vrf.NewVRF(addr, client)
	helpers.PanicErr(err)
	return vrf
}

func newDKG(addr common.Address, client *ethclient.Client) *dkg.DKG {
	dkg, err := dkg.NewDKG(addr, client)
	helpers.PanicErr(err)
	return dkg
}
