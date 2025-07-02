package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/por_mock_ocr3plugin/chainsupport"
	"github.com/smartcontractkit/por_mock_ocr3plugin/contractconfig"
	"github.com/smartcontractkit/por_mock_ocr3plugin/myname"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const network = "sepolia"

func main() {
	if err := fund(); err != nil {
		panic(err)
	}
	fmt.Printf("https://sepolia.etherscan.io/address/%s\n", contractconfig.DestinationAddress().Hex())
	fmt.Println("My destination address is ", contractconfig.DestinationAddress().Hex())
}

func fund() error {
	client := chainsupport.EthClient(chainsupport.InfuraUrl(network))
	ctx := context.Background()

	deployerPrivateKey := contractconfig.GodPrivateKey()
	opts, err := bind.NewKeyedTransactorWithChainID(
		&deployerPrivateKey,
		big.NewInt(int64(chainsupport.NetworkToChainID[network])),
	)
	fmt.Printf("Funder address: %s\n", opts.From.Hex())
	if err != nil {
		panic(fmt.Sprintf("bind.NewKeyedTransactorWithChainID: %v", err))
	}

	nonce, err := client.PendingNonceAt(ctx, opts.From)
	if err != nil {
		return err
	}

	to := contractconfig.TransmitterAddress(0)

	fmt.Println("Name of the developer whose account we are funding:", myname.Name)
	fmt.Println("Funded OCR Transmitter Address", to.Hex())

	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:    big.NewInt(int64(chainsupport.NetworkToChainID[network])),
		Nonce:      nonce,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(100e9),
		Gas:        100_000,
		To:         &to,
		Value:      big.NewInt(5e16),
		Data:       nil,
		AccessList: nil,

		V: nil, R: nil, S: nil,
	})

	signedTx, err := opts.Signer(opts.From, tx)
	if err != nil {
		return err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}

	return nil
}
