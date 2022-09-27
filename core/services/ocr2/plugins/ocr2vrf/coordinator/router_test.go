package coordinator

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/coordinator/mocks"
)

var nilOpts *bind.CallOpts

func TestRouter_SProvingKeyHash(t *testing.T) {
	beacon := mocks.NewVRFBeaconInterface(t)
	keyHash := [32]byte{1}
	router := vrfRouter{
		beacon: beacon,
	}
	beacon.On("SProvingKeyHash", mock.Anything).Return(keyHash, nil).Once()

	result, err := router.SProvingKeyHash(nilOpts)
	assert.NoError(t, err)
	assert.Equal(t, keyHash, result)
}

func TestRouter_SKeyID(t *testing.T) {
	beacon := mocks.NewVRFBeaconInterface(t)
	keyID := [32]byte{2}
	router := vrfRouter{
		beacon: beacon,
	}
	beacon.On("SKeyID", mock.Anything).Return(keyID, nil).Once()

	result, err := router.SKeyID(nilOpts)
	assert.NoError(t, err)
	assert.Equal(t, keyID, result)
}

func TestRouter_IBeaconPeriodBlocks(t *testing.T) {
	coordinator := mocks.NewVRFCoordinatorInterface(t)
	periodBlocks := big.NewInt(3)
	router := vrfRouter{
		coordinator: coordinator,
	}
	coordinator.On("IBeaconPeriodBlocks", mock.Anything).Return(periodBlocks, nil).Once()

	result, err := router.IBeaconPeriodBlocks(nilOpts)
	assert.NoError(t, err)
	assert.Equal(t, periodBlocks, result)
}

func TestRouter_GetConfirmationDelays(t *testing.T) {
	coordinator := mocks.NewVRFCoordinatorInterface(t)
	confDelays := [8]*big.Int{big.NewInt(4)}
	router := vrfRouter{
		coordinator: coordinator,
	}
	coordinator.On("GetConfirmationDelays", mock.Anything).Return(confDelays, nil).Once()

	result, err := router.GetConfirmationDelays(nilOpts)
	assert.NoError(t, err)
	assert.Equal(t, confDelays, result)
}

func TestRouter_ParseLog(t *testing.T) {
	t.Parallel()

	t.Run("parse beacon log", func(t *testing.T) {
		addr := newAddress(t)
		log := types.Log{
			Address: addr,
		}
		parsedLog := vrf_beacon.VRFBeaconNewTransmission{}
		beacon := mocks.NewVRFBeaconInterface(t)
		router := vrfRouter{
			beacon: beacon,
		}
		beacon.On("Address").Return(addr).Once()
		beacon.On("ParseLog", log).Return(parsedLog, nil).Once()

		result, err := router.ParseLog(log)
		assert.NoError(t, err)
		assert.Equal(t, result, parsedLog)
	})

	t.Run("parse coordinator log", func(t *testing.T) {
		addr := newAddress(t)
		log := types.Log{
			Address: addr,
		}
		parsedLog := vrf_coordinator.VRFCoordinatorRandomnessRequested{}
		beacon := mocks.NewVRFBeaconInterface(t)
		coordinator := mocks.NewVRFCoordinatorInterface(t)
		router := vrfRouter{
			beacon:      beacon,
			coordinator: coordinator,
		}
		beacon.On("Address").Return(newAddress(t)).Once()
		coordinator.On("Address").Return(addr).Once()
		coordinator.On("ParseLog", log).Return(parsedLog, nil).Once()

		result, err := router.ParseLog(log)
		assert.NoError(t, err)
		assert.Equal(t, result, parsedLog)
	})

	t.Run("parse log unexpected log", func(t *testing.T) {
		log := types.Log{
			Address: newAddress(t),
		}
		beacon := mocks.NewVRFBeaconInterface(t)
		coordinator := mocks.NewVRFCoordinatorInterface(t)
		router := vrfRouter{
			beacon:      beacon,
			coordinator: coordinator,
		}
		beacon.On("Address").Return(newAddress(t)).Once()
		coordinator.On("Address").Return(newAddress(t)).Once()

		result, err := router.ParseLog(log)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse log")
	})
}
