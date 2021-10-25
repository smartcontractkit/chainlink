package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/web/loader"
	"github.com/smartcontractkit/chainlink/core/web/schema"
)

// gqlTestFramework is a framework wrapper containing the objects needed to run
// a GQL test.
type gqlTestFramework struct {
	t *testing.T

	// The mocked chainlink.Application
	App *mocks.Application

	// The root GQL schema
	RootSchema *graphql.Schema

	// Contains the context with an injected dataloader
	Ctx context.Context
}

// setupFramework sets up the framework for all GQL testing
func setupFramework(t *testing.T) *gqlTestFramework {
	t.Helper()

	var (
		app        = &mocks.Application{}
		rootSchema = graphql.MustParseSchema(
			schema.MustGetRootSchema(),
			&Resolver{App: app},
		)
		ctx = loader.InjectDataloader(context.Background(), app)
	)

	t.Cleanup(func() {
		app.AssertExpectations(t)
	})

	return &gqlTestFramework{
		t:          t,
		App:        app,
		RootSchema: rootSchema,
		Ctx:        ctx,
	}
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
