package txm

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	"github.com/smartcontractkit/chainlink-framework/chains"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/storage"
)

func TestOrchestratorLifecycle(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	address1 := testutils.NewAddress()
	keystore := keystest.Addresses{address1}
	client := newMockClient(t)
	var nonce uint64
	client.On("PendingNonceAt", mock.Anything, address1).Return(nonce, nil).Once()
	txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
	ab := newMockAttemptBuilder(t)
	oab := newMockOrchestratorAttemptBuilder[common.Hash, chains.Head[common.Hash]](t)
	oab.On("Start", mock.Anything).Return(nil).Once()
	oab.On("Close").Return(nil).Once()
	config := Config{BlockTime: 1 * time.Minute, RetryBlockThreshold: 5}

	txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
	o := NewTxmOrchestrator(lggr, testutils.FixtureChainID, txm, txStore, nil, keystore, oab)
	servicetest.Run(t, o)
}

func TestOrchestratorAddressUpdate(t *testing.T) {
	t.Parallel()

	lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
	address1 := testutils.NewAddress()
	address2 := testutils.NewAddress()
	keystore := &keystest.Addresses{address1}
	client := newMockClient(t)
	var nonce uint64
	txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
	ab := newMockAttemptBuilder(t)
	oab := newMockOrchestratorAttemptBuilder[common.Hash, chains.Head[common.Hash]](t)
	oab.On("Start", mock.Anything).Return(nil).Once()
	oab.On("Close").Return(nil).Once()
	config := Config{BlockTime: 1 * time.Minute, RetryBlockThreshold: 5}

	txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
	o := NewTxmOrchestrator(lggr, testutils.FixtureChainID, txm, txStore, nil, keystore, oab)
	client.On("PendingNonceAt", mock.Anything, address1).Return(nonce, nil).Once()
	servicetest.Run(t, o)
	tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("Set initial nonce for address: %v to %d", address1, nonce))

	client.On("PendingNonceAt", mock.Anything, address2).Return(nonce+1, nil).Once()
	*keystore = keystest.Addresses{address2}
	tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("Added the following address to InMemoryStoreManager: %s", address2))
	tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("Removed the following address from InMemoryStoreManager: %s", address1))
	tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("Restarted the following addresses: %v", *keystore))
}
