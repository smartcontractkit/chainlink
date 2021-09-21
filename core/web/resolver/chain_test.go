package resolver

import (
	"testing"

	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/stretchr/testify/mock"

	evmORMMocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_Chains(t *testing.T) {
	var (
		f      = setupFramework(t)
		evmORM = &evmORMMocks.ORM{}

		chainID = *utils.NewBigI(1)
		nodeID  = int32(200)
	)

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			evmORM,
		)
	})

	f.App.On("EVMORM").Return(evmORM)
	evmORM.On("Chains", PageDefaultOffset, PageDefaultLimit).Return([]types.Chain{
		{
			ID:        chainID,
			Enabled:   true,
			CreatedAt: f.Timestamp(),
		},
	}, 1, nil)
	evmORM.On("GetNodesByChainIDs", []utils.Big{chainID}).
		Return([]types.Node{
			{
				ID:         nodeID,
				EVMChainID: chainID,
			},
		}, nil)

	gqltesting.RunTest(t, &gqltesting.Test{
		Context: f.Ctx,
		Schema:  f.RootSchema,
		Query: `
			{
				chains {
					id
					enabled
					createdAt
					nodes {
						id
					}
				}
			}
		`,
		ExpectedResult: `
			{
				"chains": [{
					"id": "1",
					"enabled": true,
					"createdAt": "2021-01-01T00:00:00Z",
					"nodes": [{
						"id": "200"
					}]
				}]
			}
		`,
	})
}

func Test_Chain(t *testing.T) {
	var (
		f      = setupFramework(t)
		evmORM = &evmORMMocks.ORM{}

		chainID = *utils.NewBigI(1)
		nodeID  = int32(200)
	)

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			evmORM,
		)
	})

	f.App.On("EVMORM").Return(evmORM)
	evmORM.On("Chain", chainID).Return(types.Chain{
		ID:        chainID,
		Enabled:   true,
		CreatedAt: f.Timestamp(),
	}, nil)
	evmORM.On("GetNodesByChainIDs", []utils.Big{chainID}).
		Return([]types.Node{
			{
				ID:         nodeID,
				EVMChainID: chainID,
			},
		}, nil)

	gqltesting.RunTest(t, &gqltesting.Test{
		Context: f.Ctx,
		Schema:  f.RootSchema,
		Query: `
			{
				chain(id: "1") {
					id
					enabled
					createdAt
					nodes {
						id
					}
				}
			}
		`,
		ExpectedResult: `
			{
				"chain": {
					"id": "1",
					"enabled": true,
					"createdAt": "2021-01-01T00:00:00Z",
					"nodes": [{
						"id": "200"
					}]
				}
			}
		`,
	})
}
