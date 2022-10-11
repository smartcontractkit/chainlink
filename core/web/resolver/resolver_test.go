package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/stretchr/testify/mock"

	bridgeORMMocks "github.com/smartcontractkit/chainlink/core/bridges/mocks"
	evmConfigMocks "github.com/smartcontractkit/chainlink/core/chains/evm/config/mocks"
	evmORMMocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	txmgrMocks "github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	configMocks "github.com/smartcontractkit/chainlink/core/config/mocks"
	coremocks "github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	feedsMocks "github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	jobORMMocks "github.com/smartcontractkit/chainlink/core/services/job/mocks"
	keystoreMocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	pipelineMocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"
	clsessions "github.com/smartcontractkit/chainlink/core/sessions"
	sessionsMocks "github.com/smartcontractkit/chainlink/core/sessions/mocks"
	"github.com/smartcontractkit/chainlink/core/web/auth"
	"github.com/smartcontractkit/chainlink/core/web/loader"
	"github.com/smartcontractkit/chainlink/core/web/schema"
)

type mocks struct {
	bridgeORM   *bridgeORMMocks.ORM
	evmORM      *evmtest.MockORM
	jobORM      *jobORMMocks.ORM
	sessionsORM *sessionsMocks.ORM
	pipelineORM *pipelineMocks.ORM
	feedsSvc    *feedsMocks.Service
	cfg         *configMocks.GeneralConfig
	scfg        *evmConfigMocks.ChainScopedConfig
	ocr         *keystoreMocks.OCR
	ocr2        *keystoreMocks.OCR2
	csa         *keystoreMocks.CSA
	keystore    *keystoreMocks.Master
	ethKs       *keystoreMocks.Eth
	p2p         *keystoreMocks.P2P
	vrf         *keystoreMocks.VRF
	solana      *keystoreMocks.Solana
	chain       *evmORMMocks.Chain
	chainSet    *evmORMMocks.ChainSet
	ethClient   *evmORMMocks.Client
	eIMgr       *webhookmocks.ExternalInitiatorManager
	balM        *evmORMMocks.BalanceMonitor
	txmORM      *txmgrMocks.ORM
	auditLogger *audit.AuditLoggerService
}

// gqlTestFramework is a framework wrapper containing the objects needed to run
// a GQL test.
type gqlTestFramework struct {
	t *testing.T

	// The mocked chainlf.Mocks.chainSetink.Application
	App *coremocks.Application

	// The root GQL schema
	RootSchema *graphql.Schema

	// Contains the context with an injected dataloader
	Ctx context.Context

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
		ctx = loader.InjectDataloader(testutils.Context(t), app)
	)

	// Setup mocks
	// Note - If you add a new mock make sure you assert it's expectation below.
	m := &mocks{
		bridgeORM:   bridgeORMMocks.NewORM(t),
		evmORM:      evmtest.NewMockORM(nil, nil),
		jobORM:      jobORMMocks.NewORM(t),
		feedsSvc:    feedsMocks.NewService(t),
		sessionsORM: sessionsMocks.NewORM(t),
		pipelineORM: pipelineMocks.NewORM(t),
		cfg:         configMocks.NewGeneralConfig(t),
		scfg:        evmConfigMocks.NewChainScopedConfig(t),
		ocr:         keystoreMocks.NewOCR(t),
		ocr2:        keystoreMocks.NewOCR2(t),
		csa:         keystoreMocks.NewCSA(t),
		keystore:    keystoreMocks.NewMaster(t),
		ethKs:       keystoreMocks.NewEth(t),
		p2p:         keystoreMocks.NewP2P(t),
		vrf:         keystoreMocks.NewVRF(t),
		solana:      keystoreMocks.NewSolana(t),
		chain:       evmORMMocks.NewChain(t),
		chainSet:    evmORMMocks.NewChainSet(t),
		ethClient:   evmORMMocks.NewClient(t),
		eIMgr:       webhookmocks.NewExternalInitiatorManager(t),
		balM:        evmORMMocks.NewBalanceMonitor(t),
		txmORM:      txmgrMocks.NewORM(t),
		auditLogger: &audit.AuditLoggerService{},
	}

	app.Mock.On("GetAuditLogger", mock.Anything, mock.Anything).Return(audit.NoopLogger).Maybe()

	f := &gqlTestFramework{
		t:          t,
		App:        app,
		RootSchema: rootSchema,
		Ctx:        ctx,
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

// injectAuthenticatedUser injects a session into the request context
func (f *gqlTestFramework) injectAuthenticatedUser() {
	f.t.Helper()

	user := clsessions.User{Email: "gqltester@chain.link", Role: clsessions.UserRoleAdmin}

	f.Ctx = auth.SetGQLAuthenticatedSession(f.Ctx, user, "gqltesterSession")
}

// GQLTestCase represents a single GQL request test.
type GQLTestCase struct {
	name          string
	authenticated bool
	before        func(*gqlTestFramework)
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

			var (
				f = setupFramework(t)
			)

			if tc.authenticated {
				f.injectAuthenticatedUser()
			}

			if tc.before != nil {
				tc.before(f)
			}

			// This does not print out the correct stack trace as the `RunTest`
			// function does not call t.Helper(). It insteads displays the file
			// and line location of the `gqltesting` package.
			//
			// This would need to be fixed upstream.
			gqltesting.RunTest(t, &gqltesting.Test{
				Context:        f.Ctx,
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
