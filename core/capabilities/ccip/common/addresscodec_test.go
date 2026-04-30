package common_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	sel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

type testAddressCodec struct{}

func (testAddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return fmt.Sprintf("address-%x", addr), nil
}

func (testAddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return []byte(addr), nil
}

func (testAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	return []byte{oracleID}, nil
}

func (testAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return fmt.Sprintf("transmitter-%x", addr), nil
}

func TestAddressCodec_RegisterCodec(t *testing.T) {
	t.Parallel()

	addressCodec := ccipcommon.NewAddressCodec(nil)
	require.False(t, addressCodec.HasCodec(sel.FamilyEVM))

	addressCodec.RegisterCodec(sel.FamilyEVM, testAddressCodec{})
	require.True(t, addressCodec.HasCodec(sel.FamilyEVM))

	selector := ccipocr3common.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)

	address, err := addressCodec.AddressBytesToString([]byte{0x01}, selector)
	require.NoError(t, err)
	require.Equal(t, "address-01", address)

	transmitter, err := addressCodec.TransmitterBytesToString([]byte{0x02}, selector)
	require.NoError(t, err)
	require.Equal(t, "transmitter-02", transmitter)

	addressBytes, err := addressCodec.AddressStringToBytes("hello", selector)
	require.NoError(t, err)
	require.Equal(t, ccipocr3common.UnknownAddress("hello"), addressBytes)

	oracleAddress, err := addressCodec.OracleIDAsAddressBytes(7, selector)
	require.NoError(t, err)
	require.Equal(t, []byte{7}, oracleAddress)
}

func TestAddressCodec_ConcurrentAccessThroughCopies(t *testing.T) {
	t.Parallel()

	selector := ccipocr3common.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)
	addressCodec := ccipcommon.NewAddressCodec(map[string]ccipcommon.ChainSpecificAddressCodec{
		sel.FamilyEVM: testAddressCodec{},
	})
	addressCodecCopy := addressCodec

	const (
		readers    = 8
		iterations = 1000
	)

	start := make(chan struct{})
	errCh := make(chan error, readers)
	var wg sync.WaitGroup

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			for i := 0; i < iterations; i++ {
				if !addressCodecCopy.HasCodec(sel.FamilyEVM) {
					errCh <- errors.New("expected EVM codec to be registered")
					return
				}

				address, err := addressCodecCopy.AddressBytesToString([]byte{0x01}, selector)
				if err != nil {
					errCh <- err
					return
				}
				if address != "address-01" {
					errCh <- fmt.Errorf("unexpected address: %s", address)
					return
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start

		for i := 0; i < iterations; i++ {
			addressCodec.RegisterCodec(sel.FamilyEVM, testAddressCodec{})
		}
	}()

	close(start)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err)
	}
}
