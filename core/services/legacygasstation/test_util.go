package legacygasstation

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	TokenName    string = "tokenName"
	TokenVersion string = "1"
)

var (
	FromAddress      common.Address = common.HexToAddress("0x8b94A8792dcbb482F2A49569AcE5E7c29fF5c93d")
	TargetAddress    common.Address = common.HexToAddress("0x61aF15229af1CEd7Ca4a80f623F85a8a91420C04")
	ForwarderAddress common.Address = common.HexToAddress("0x16dCd7F98B0a62e09e86Ba95201334Da9C718Da1")
	ReceiverAddress  common.Address = common.HexToAddress("0xCA1767E3f243874b7d5c0a231e5F1cEA3659A59e")
	OfframpAddress   common.Address = common.HexToAddress("0x93d8156d7F0271cCe4bcE07501fdD0534D335932")
)

type TestLegacyGaslessTx struct {
	ID                 string
	Forwarder          common.Address
	From               common.Address
	Target             common.Address
	Receiver           common.Address
	Nonce              *big.Int
	Amount             *big.Int
	SourceChainID      uint64
	DestinationChainID uint64
	ValidUntilTime     *big.Int
	Signature          []byte
	Status             types.Status
	FailureReason      *string
	TokenName          string
	TokenVersion       string
	CCIPMessageID      *common.Hash
	EthTxID            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func LegacyGaslessTx(t *testing.T,
	testTx TestLegacyGaslessTx,
) types.LegacyGaslessTx {
	id := uuid.New().String()
	if testTx.ID != "" {
		id = testTx.ID
	}
	forwarder := ForwarderAddress
	if !utils.IsEmptyAddress(testTx.Forwarder) {
		forwarder = testTx.Forwarder
	}
	from := FromAddress
	if !utils.IsEmptyAddress(testTx.From) {
		from = testTx.From
	}
	target := TargetAddress
	if !utils.IsEmptyAddress(testTx.Target) {
		target = testTx.Target
	}
	receiver := ReceiverAddress
	if !utils.IsEmptyAddress(testTx.Receiver) {
		receiver = testTx.Receiver
	}
	nonce, ok := new(big.Int).SetString("41348595284370659424721411682442273173799930917222432033995375917025790037177", 10)
	require.True(t, ok)
	if testTx.Nonce != nil {
		nonce = testTx.Nonce
	}
	amount, ok := new(big.Int).SetString("10000000000000000000000", 10)
	require.True(t, ok)
	if testTx.Amount != nil {
		amount = testTx.Amount
	}
	validUntilTime := big.NewInt(1682576193)
	if testTx.ValidUntilTime != nil {
		validUntilTime = testTx.ValidUntilTime
	}
	sourceChainID := uint64(1337)
	if testTx.SourceChainID != 0 {
		sourceChainID = testTx.SourceChainID
	}
	destChainID := uint64(1000)
	if testTx.DestinationChainID != 0 {
		destChainID = testTx.DestinationChainID
	}
	status := types.Submitted
	if testTx.Status != status {
		status = testTx.Status
	}
	tokenName := TokenName
	if testTx.TokenName != "" {
		tokenName = TokenName
	}
	tokenVersion := TokenVersion
	if testTx.TokenVersion != "" {
		tokenVersion = TokenVersion
	}
	signature, err := base64.StdEncoding.DecodeString("a9VQaaVBf5W2O/rppOutrjsoq9Sk7m+aVoBuT/2o2ykT3hzzHQmtDmELLr/noQeUPqHdSWDPh1xL540G/FNm+xs=")
	require.NoError(t, err)
	if testTx.Signature != nil {
		signature = testTx.Signature[:]
	}
	tx := types.LegacyGaslessTx{
		ID:                 id,
		Forwarder:          forwarder,
		From:               from,
		Target:             target,
		Receiver:           receiver,
		Nonce:              utils.NewBig(nonce),
		Amount:             utils.NewBig(amount),
		SourceChainID:      sourceChainID,
		DestinationChainID: destChainID,
		ValidUntilTime:     utils.NewBig(validUntilTime),
		Signature:          signature,
		Status:             status,
		FailureReason:      testTx.FailureReason,
		TokenName:          tokenName,
		TokenVersion:       tokenVersion,
		CCIPMessageID:      testTx.CCIPMessageID,
		EthTxID:            testTx.EthTxID,
	}
	return tx
}

func AssertTxEquals(t *testing.T, expected, actual types.LegacyGaslessTx) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Forwarder, actual.Forwarder)
	require.Equal(t, expected.Target, actual.Target)
	require.Equal(t, expected.Receiver, actual.Receiver)
	require.True(t, expected.Nonce.Cmp(actual.Nonce) == 0)
	require.True(t, expected.Amount.Cmp(actual.Amount) == 0)
	require.True(t, expected.ValidUntilTime.Cmp(actual.ValidUntilTime) == 0)
	require.Equal(t, expected.SourceChainID, actual.SourceChainID)
	require.Equal(t, expected.DestinationChainID, actual.DestinationChainID)
	require.True(t, bytes.Equal(expected.Signature, actual.Signature))
	require.Equal(t, expected.Status, actual.Status)
	require.Equal(t, expected.TokenName, actual.TokenName)
	require.Equal(t, expected.TokenVersion, actual.TokenVersion)
	require.Equal(t, expected.EthTxID, actual.EthTxID)
	if expected.FailureReason == nil {
		require.Nil(t, actual.FailureReason)
	} else {
		require.Equal(t, *expected.FailureReason, *actual.FailureReason)
	}
	if expected.CCIPMessageID == nil {
		require.Nil(t, actual.CCIPMessageID)
	} else {
		require.True(t, bytes.Equal(expected.CCIPMessageID.Bytes(), actual.CCIPMessageID.Bytes()))
	}
}

type TestStatusUpdateServer struct {
	Server *http.Server
	Port   uint16
}

func NewUnstartedStatusUpdateServer(t *testing.T) TestStatusUpdateServer {
	router := gin.Default()
	router.POST("/return_success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"jsonrpc": "2.0",
		})
	})
	router.POST("/return_not_found", func(c *gin.Context) {
		c.AbortWithStatus(404)
	})
	router.POST("/return_error", func(c *gin.Context) {
		c.AbortWithStatus(404)
		c.JSON(http.StatusOK, gin.H{
			"jsonrpc": "2.0",
		})
	})

	port := GetFreePort(t)
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     router,
		ReadTimeout: 30 * time.Second,
	}

	return TestStatusUpdateServer{
		Server: server,
		Port:   port,
	}
}

func (s *TestStatusUpdateServer) Start() error {
	return s.Server.ListenAndServe()
}

func (s *TestStatusUpdateServer) Stop() error {
	return s.Server.Close()
}

func GetFreePort(t *testing.T) uint16 {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer func() { assert.NoError(t, l.Close()) }()

	return uint16(l.Addr().(*net.TCPAddr).Port)
}
