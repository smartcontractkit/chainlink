package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	evmConfigMocks "github.com/smartcontractkit/chainlink-evm/pkg/config/mocks"
	evmMonMocks "github.com/smartcontractkit/chainlink-evm/pkg/monitor/mocks"
	legacyEvmORMMocks "github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	evmtxmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"
	bridgeORMMocks "github.com/smartcontractkit/chainlink/v2/core/bridges/mocks"
	coremocks "github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	chainlinkMocks "github.com/smartcontractkit/chainlink/v2/core/services/chainlink/mocks"
	feedsMocks "github.com/smartcontractkit/chainlink/v2/core/services/feeds/mocks"
	jobORMMocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	keystoreMocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	pipelineMocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
	webhookmocks "github.com/smartcontractkit/chainlink/v2/core/services/webhook/mocks"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
	authProviderMocks "github.com/smartcontractkit/chainlink/v2/core/sessions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/web/auth"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
	"github.com/smartcontractkit/chainlink/v2/core/web/schema"
)

type mocks struct {
	bridgeORM            *bridgeORMMocks.ORM
	evmORM               *evmtest.TestConfigs
	jobORM               *jobORMMocks.ORM
	authProvider         *authProviderMocks.AuthenticationProvider
	pipelineORM          *pipelineMocks.ORM
	feedsSvc             *feedsMocks.Service
	cfg                  *chainlinkMocks.GeneralConfig
	scfg                 *evmConfigMocks.ChainScopedConfig
	ocr                  *keystoreMocks.OCR
	ocr2                 *keystoreMocks.OCR2
	csa                  *keystoreMocks.CSA
	keystore             *keystoreMocks.Master
	ethKs                *keystoreMocks.Eth
	p2p                  *keystoreMocks.P2P
	vrf                  *keystoreMocks.VRF
	solana               *keystoreMocks.Solana
	aptos                *keystoreMocks.Aptos
	cosmos               *keystoreMocks.Cosmos
	starknet             *keystoreMocks.StarkNet
	tron                 *keystoreMocks.Tron
	ton                  *keystoreMocks.TON
	chain                *legacyEvmORMMocks.Chain
	legacyEVMChains      *legacyEvmORMMocks.LegacyChainContainer
	relayerChainInterops *chainlinkMocks.FakeRelayerChainInteroperators
	ethClient            *clienttest.Client
	eIMgr                *webhookmocks.ExternalInitiatorManager
	balM                 *evmMonMocks.BalanceMonitor
	txmStore             *evmtxmgrmocks.EvmTxStore
	auditLogger          *audit.AuditLoggerService
}

// gqlTestFramework is a framework wrapper containing the objects needed to run
// a GQL test.
type gqlTestFramework struct {
	t *testing.T

	// The mocked chainlink.Application
	App *coremocks.Application

	// The root GQL schema
	RootSchema *graphql.Schema

	Mocks *mocks
}

// setupFramework sets up the framework for all GQL testing
func setupFramework(t *testing.T) *gqlTestFramework {
	t.Helper()

	var (
		app        = coremocks.NewApplication(t)
		rootSchema = graphql.MustParseSchema(
			schema.MustGetRootSchema(),
			&Resolver{App: app},
		)
	)

	// Setup mocks
	// Note - If you add a new mock make sure you assert it's expectation below.
	m := &mocks{
		bridgeORM:            bridgeORMMocks.NewORM(t),
		evmORM:               evmtest.NewTestConfigs(),
		jobORM:               jobORMMocks.NewORM(t),
		feedsSvc:             feedsMocks.NewService(t),
		authProvider:         authProviderMocks.NewAuthenticationProvider(t),
		pipelineORM:          pipelineMocks.NewORM(t),
		cfg:                  chainlinkMocks.NewGeneralConfig(t),
		scfg:                 evmConfigMocks.NewChainScopedConfig(t),
		ocr:                  keystoreMocks.NewOCR(t),
		ocr2:                 keystoreMocks.NewOCR2(t),
		csa:                  keystoreMocks.NewCSA(t),
		keystore:             keystoreMocks.NewMaster(t),
		ethKs:                keystoreMocks.NewEth(t),
		p2p:                  keystoreMocks.NewP2P(t),
		vrf:                  keystoreMocks.NewVRF(t),
		solana:               keystoreMocks.NewSolana(t),
		aptos:                keystoreMocks.NewAptos(t),
		cosmos:               keystoreMocks.NewCosmos(t),
		starknet:             keystoreMocks.NewStarkNet(t),
		tron:                 keystoreMocks.NewTron(t),
		ton:                  keystoreMocks.NewTON(t),
		chain:                legacyEvmORMMocks.NewChain(t),
		legacyEVMChains:      legacyEvmORMMocks.NewLegacyChainContainer(t),
		relayerChainInterops: &chainlinkMocks.FakeRelayerChainInteroperators{},
		ethClient:            clienttest.NewClient(t),
		eIMgr:                webhookmocks.NewExternalInitiatorManager(t),
		balM:                 evmMonMocks.NewBalanceMonitor(t),
		txmStore:             evmtxmgrmocks.NewEvmTxStore(t),
		auditLogger:          &audit.AuditLoggerService{},
	}

	lggr := logger.TestLogger(t)
	app.Mock.On("GetAuditLogger", mock.Anything, mock.Anything).Return(audit.NoopLogger).Maybe()
	app.Mock.On("GetLogger").Return(lggr).Maybe()

	f := &gqlTestFramework{
		t:          t,
		App:        app,
		RootSchema: rootSchema,
		Mocks:      m,
	}

	return f
}

// Timestamp returns a static timestamp.
//
// Use this in tests by interpolating it into the result string. If you don't
// want to interpolate you can instead use the formatted output of
// `2021-01-01T00:00:00Z`
func (f *gqlTestFramework) Timestamp() time.Time {
	f.t.Helper()

	return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
}

// withAuthenticatedUser injects a session into the request context
func (f *gqlTestFramework) withAuthenticatedUser(ctx context.Context) context.Context {
	user := clsessions.User{Email: "gqltester@chain.link", Role: clsessions.UserRoleAdmin}

	return auth.WithGQLAuthenticatedSession(ctx, user, "gqltesterSession")
}

// GQLTestCase represents a single GQL request test.
type GQLTestCase struct {
	name          string
	authenticated bool
	before        func(context.Context, *gqlTestFramework)
	query         string
	variables     map[string]interface{}
	result        string
	errors        []*gqlerrors.QueryError
}

// RunGQLTests runs a set of GQL tests cases
func RunGQLTests(t *testing.T, testCases []GQLTestCase) {
	t.Helper()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := setupFramework(t)
			ctx := loader.InjectDataloader(testutils.Context(t), f.App)

			if tc.authenticated {
				ctx = f.withAuthenticatedUser(ctx)
			}

			if tc.before != nil {
				tc.before(ctx, f)
			}

			// This does not print out the correct stack trace as the `RunTest`
			// function does not call t.Helper(). It insteads displays the file
			// and line location of the `gqltesting` package.
			//
			// This would need to be fixed upstream.
			gqltesting.RunTest(t, &gqltesting.Test{
				Context:        ctx,
				Schema:         f.RootSchema,
				Query:          tc.query,
				Variables:      tc.variables,
				ExpectedResult: tc.result,
				ExpectedErrors: tc.errors,
			})
		})
	}
}

// unauthorizedTestCase generates an unauthorized test case from another test
// case.
//
// The paths will be the query/mutation definition name
func unauthorizedTestCase(tc GQLTestCase, paths ...interface{}) GQLTestCase {
	tc.name = "not authorized"
	tc.authenticated = false
	tc.result = "null"
	tc.errors = []*gqlerrors.QueryError{
		{
			ResolverError: unauthorizedError{},
			Path:          paths,
			Message:       "Unauthorized",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		},
	}

	return tc
}
