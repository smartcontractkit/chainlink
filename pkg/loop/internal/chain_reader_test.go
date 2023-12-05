package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

func TestVersionedBytesFunctions(t *testing.T) {
	const unsupportedVer = 25913
	t.Run("EncodeVersionedBytes unsupported type", func(t *testing.T) {
		expected := types.ErrInvalidType
		invalidData := make(chan int)

		_, err := encodeVersionedBytes(invalidData, JSONEncodingVersion2)
		if err == nil || !errors.Is(err, expected) {
			t.Errorf("expected error: %s, but got: %v", expected, err)
		}
	})

	t.Run("EncodeVersionedBytes unsupported encoding version", func(t *testing.T) {
		expected := fmt.Errorf("unsupported encoding version %d for data map[key:value]", unsupportedVer)
		data := map[string]interface{}{
			"key": "value",
		}

		_, err := encodeVersionedBytes(data, unsupportedVer)
		if err == nil || err.Error() != expected.Error() {
			t.Errorf("expected error: %s, but got: %v", expected, err)
		}
	})

	t.Run("DecodeVersionedBytes", func(t *testing.T) {
		var decodedData map[string]interface{}
		expected := fmt.Errorf("unsupported encoding version %d for versionedData [97 98 99 100 102]", unsupportedVer)
		versionedBytes := &pb.VersionedBytes{
			Version: unsupportedVer, // Unsupported version
			Data:    []byte("abcdf"),
		}

		err := decodeVersionedBytes(&decodedData, versionedBytes)
		if err == nil || err.Error() != expected.Error() {
			t.Errorf("expected error: %s, but got: %v", expected, err)
		}
	})
}

func TestChainReaderClient(t *testing.T) {
	RunChainReaderInterfaceTests(t, &interfaceTester{})

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	es := &errorServer{}
	pb.RegisterChainReaderServer(s, es)

	chSrv := make(chan error)
	go func() {
		chSrv <- s.Serve(lis)
	}()
	defer func() {
		s.Stop()
		err := <-chSrv
		assert.NoError(t, err)
	}()
	conn := connFromLis(t, lis)
	client := &chainReaderClient{grpc: pb.NewChainReaderClient(conn)}
	ctx := context.Background()

	type testCase struct {
		errType error
		errMsg  string
	}

	testCases := []testCase{
		{types.ErrChainReaderConfigMissing, ""},
		{types.ErrInvalidConfig, "ChainReader config doesn't match abi"},
		{types.ErrInvalidType, "wrong length"},
	}

	for _, tc := range testCases {
		_, ok := tc.errType.(error)
		require.True(t, ok)
		if tc.errMsg != "" {
			es.err = fmt.Errorf("%w: %s", tc.errType, tc.errMsg)
		} else {
			es.err = tc.errType
		}

		t.Run("GetLatestValue unwraps errors from server "+tc.errType.Error(), func(t *testing.T) {
			err := client.GetLatestValue(ctx, types.BoundContract{}, "method", "anything", "anything")
			assert.ErrorIs(t, tc.errType, err) // Note: our custom error type must be first arg, otherwise gRPC's Is() is called instead of ours
		})
	}

	// make sure that errors come from client directly
	es.err = nil
	invalidTypeErr := types.ErrInvalidType
	t.Run("GetLatestValue returns error if type cannot be encoded in the wire format", func(t *testing.T) {
		err := client.GetLatestValue(ctx, types.BoundContract{}, "method", &cannotEncode{}, &TestStruct{})
		assert.ErrorIs(t, err, &invalidTypeErr)
	})
}

type interfaceTester struct {
	lis    *bufconn.Listener
	server *grpc.Server
	conn   *grpc.ClientConn
	fs     *fakeCodecServer
}

var encoder = makeEncoder()

func makeEncoder() cbor.EncMode {
	opts := cbor.CoreDetEncOptions()
	opts.Sort = cbor.SortCanonical
	e, _ := opts.EncMode()
	return e
}

func (it *interfaceTester) SetLatestValue(ctx context.Context, t *testing.T, testStruct *TestStruct) types.BoundContract {
	it.fs.SetLatestValue(testStruct)
	return types.BoundContract{}
}

func (it *interfaceTester) GetPrimitiveContract(ctx context.Context, t *testing.T) types.BoundContract {
	return types.BoundContract{}
}

func (it *interfaceTester) GetSliceContract(ctx context.Context, t *testing.T) types.BoundContract {
	return types.BoundContract{}
}

func (it *interfaceTester) GetAccountBytes(_ int) []byte {
	return []byte{1, 2, 3}
}

func (it *interfaceTester) Setup(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	it.lis = lis
	it.fs = &fakeCodecServer{lock: &sync.Mutex{}}
	it.server = grpc.NewServer()
	pb.RegisterChainReaderServer(it.server, &chainReaderServer{impl: it.fs})
	go func() {
		if err := it.server.Serve(lis); err != nil {
			t.Error(err)
		}
	}()

	srvCh := make(chan error)
	go func() {
		srvCh <- it.server.Serve(lis)
	}()

	t.Cleanup(func() {
		if it.server != nil {
			it.server.Stop()
			err := <-srvCh
			assert.NoError(t, err)
		}
		if it.conn != nil {
			assert.NoError(t, it.conn.Close())
		}

		it.lis = nil
		it.server = nil
		it.conn = nil
	})
}

func (it *interfaceTester) Name() string {
	return "relay client"
}

func (it *interfaceTester) GetChainReader(t *testing.T) types.ChainReader {
	if it.conn == nil {
		it.conn = connFromLis(t, it.lis)
	}

	return &chainReaderClient{grpc: pb.NewChainReaderClient(it.conn)}
}

type fakeCodecServer struct {
	lastItem any
	latest   []TestStruct
	lock     *sync.Mutex
}

func (f *fakeCodecServer) SetLatestValue(ts *TestStruct) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.latest = append(f.latest, *ts)
}

func (f *fakeCodecServer) GetLatestValue(ctx context.Context, _ types.BoundContract, method string, params, returnVal any) error {
	if method != MethodTakingLatestParamsReturningTestStruct {
		return errors.New("unknown method " + method)
	}

	f.lock.Lock()
	defer f.lock.Unlock()
	lp := params.(*map[string]interface{})
	i := (*lp)["I"].(uint64)
	return mapstructure.Decode(f.latest[i-1], returnVal)
}

type errorServer struct {
	err error
	pb.UnimplementedChainReaderServer
}

func (e *errorServer) GetLatestValue(context.Context, *pb.GetLatestValueRequest) (*pb.GetLatestValueReply, error) {
	return nil, e.err
}

func connFromLis(t *testing.T, lis *bufconn.Listener) *grpc.ClientConn {
	conn, err := grpc.Dial("bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	require.NoError(t, err)
	return conn
}

type cannotEncode struct{}

func (*cannotEncode) MarshalBinary() ([]byte, error) {
	return nil, errors.New("nope")
}

func (*cannotEncode) UnmarshalBinary() error {
	return errors.New("nope")
}

func (*cannotEncode) MarshalText() ([]byte, error) {
	return nil, errors.New("nope")
}

func (*cannotEncode) UnmarshalText() error {
	return errors.New("nope")
}
