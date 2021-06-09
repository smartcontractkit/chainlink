package csa_test

import (
	"context"
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/csa"
	"github.com/smartcontractkit/chainlink/core/services/csa/mocks"
	storeorm "github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestService struct {
	csa.Service
	orm *mocks.ORM
}

func Test_Service_CreateCSAKey(t *testing.T) {
	var (
		id  = uint(1)
		key = csa.CSAKey{}
		ctx = context.Background()
	)

	testCases := []struct {
		name     string
		beforeFn func(t *testing.T, ts *TestService)
		want     *csa.CSAKey
		wantErr  error
	}{
		{
			name: "success",
			beforeFn: func(t *testing.T, ts *TestService) {
				ts.orm.On("CountCSAKeys").
					Return(int64(0), nil).
					Once()
				ts.orm.On("CreateCSAKey", ctx, mock.IsType(&csa.CSAKey{})).
					Return(id, nil)
				ts.orm.On("GetCSAKey", ctx, id).
					Return(&key, nil)
			},
			want: &key,
		},
		{
			name: "success",
			beforeFn: func(t *testing.T, ts *TestService) {
				ts.orm.On("CountCSAKeys").
					Return(int64(1), nil).
					Once()
			},
			wantErr: errors.New("can only have 1 CSA key"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ts := setupTestService(t)

			if tc.beforeFn != nil {
				tc.beforeFn(t, ts)
			}

			actual, err := ts.CreateCSAKey()
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, actual)
			}
		})
	}
}

func Test_Service_ListCSAKeys(t *testing.T) {
	t.Parallel()

	ts := setupTestService(t)

	var (
		key  = csa.CSAKey{}
		keys = []csa.CSAKey{key}
	)

	ts.orm.On("ListCSAKeys", context.Background()).
		Return(keys, nil)

	actual, err := ts.ListCSAKeys()
	require.NoError(t, err)

	assert.Equal(t, keys, actual)
}

func Test_Service_CountCSAKeys(t *testing.T) {
	t.Parallel()

	ts := setupTestService(t)

	var (
		count = int64(1)
	)

	ts.orm.On("CountCSAKeys").
		Return(count, nil)

	actual, err := ts.CountCSAKeys()
	require.NoError(t, err)

	assert.Equal(t, count, actual)
}

func setupTestService(t *testing.T) *TestService {
	orm := &mocks.ORM{}
	cfg := &storeorm.Config{}
	cfg.SetKeystorePassword("passphrase")

	t.Cleanup(func() {
		orm.AssertExpectations(t)
	})

	return &TestService{
		orm:     orm,
		Service: csa.NewService(cfg, orm, utils.FastScryptParams),
	}
}
