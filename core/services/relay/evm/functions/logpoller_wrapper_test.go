package functions

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type subscriber struct {
	updates       sync.WaitGroup
	expectedCalls int
}

const (
	CoordinatorContractV100 = "Functions Coordinator v1.0.0"
	CoordinatorContractV200 = "Functions Coordinator v2.0.0"
	OracleRequestV200       = `[{"constant":true,"inputs":[{"indexed":true,"internalType":"bytes32","name":"requestId","type":"bytes32"},{"indexed":true,"internalType":"address","name":"requestingContract","type":"address"},{"indexed":false,"internalType":"address","name":"requestInitiator","type":"address"},{"indexed":false,"internalType":"uint64","name":"subscriptionId","type":"uint64"},{"indexed":false,"internalType":"address","name":"subscriptionOwner","type":"address"},{"indexed":false,"internalType":"bytes","name":"data","type":"bytes"},{"indexed":false,"internalType":"uint16","name":"dataVersion","type":"uint16"},{"indexed":false,"internalType":"bytes32","name":"flags","type":"bytes32"},{"indexed":false,"internalType":"uint64","name":"callbackGasLimit","type":"uint64"},{"components":[{"internalType":"bytes32","name":"requestId","type":"bytes32"},{"internalType":"address","name":"coordinator","type":"address"},{"internalType":"uint96","name":"estimatedTotalCostJuels","type":"uint96"},{"internalType":"address","name":"client","type":"address"},{"internalType":"uint64","name":"subscriptionId","type":"uint64"},{"internalType":"uint32","name":"callbackGasLimit","type":"uint32"},{"internalType":"uint72","name":"adminFee","type":"uint72"},{"internalType":"uint72","name":"donFee","type":"uint72"},{"internalType":"uint40","name":"gasOverheadBeforeCallback","type":"uint40"},{"internalType":"uint40","name":"gasOverheadAfterCallback","type":"uint40"},{"internalType":"uint32","name":"timeoutTimestamp","type":"uint32"},{"internalType":"uint72","name":"operationFee","type":"uint72"}],"indexed":false,"internalType":"structFunctionsResponse.Commitment","name":"commitment","type":"tuple"}],"name":"OracleRequest","type":"function"}]`
)

var routerAddressBytes []byte
var routerAddressHex common.Address
var coordinatorAddressBytes []byte
var coordinatorAddressHex common.Address

func (s *subscriber) UpdateRoutes(activeCoordinator common.Address, proposedCoordinator common.Address) error {
	if s.expectedCalls == 0 {
		panic("unexpected call to UpdateRoutes")
	}
	if activeCoordinator == (common.Address{}) {
		panic("activeCoordinator should not be zero")
	}
	s.expectedCalls--
	s.updates.Done()
	return nil
}

func newSubscriber(expectedCalls int) *subscriber {
	sub := &subscriber{expectedCalls: expectedCalls}
	sub.updates.Add(expectedCalls)
	return sub
}

func addr(lastByte string) ([]byte, error) {
	contractAddr, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000" + lastByte)
	if err != nil {
		return []byte{}, err
	}
	return contractAddr, nil
}

func setUp(t *testing.T, updateFrequencySec uint32) (*lpmocks.LogPoller, types.LogPollerWrapper, *evmclimocks.Client) {
	lggr := logger.TestLogger(t)
	client := evmclimocks.NewClient(t)
	lp := lpmocks.NewLogPoller(t)
	config := config.PluginConfig{
		ContractUpdateCheckFrequencySec: updateFrequencySec,
		ContractVersion:                 1,
	}
	routerAddressBytes, err := addr("01")
	require.NoError(t, err)
	routerAddressHex = common.BytesToAddress(routerAddressBytes)
	coordinatorAddressBytes, err = addr("02")
	require.NoError(t, err)
	coordinatorAddressHex = common.BytesToAddress(coordinatorAddressBytes)
	lpWrapper, err := NewLogPollerWrapper(routerAddressHex, config, client, lp, lggr)
	require.NoError(t, err)

	return lp, lpWrapper, client
}

func getMockedRequestLogV1(t *testing.T) logpoller.Log {
	// NOTE: Change this to be a more readable log generation
	data, err := hex.DecodeString("000000000000000000000000c113ba31b0080f940ca5812bbccc1e038ea9efb40000000000000000000000000000000000000000000000000000000000000001000000000000000000000000c113ba31b0080f940ca5812bbccc1e038ea9efb4000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001117082cd81744eb9504dc37f53a86db7e3fb24929b8e7507b097d501ab5b315fb20e0000000000000000000000001b4f2b0e6363097f413c249910d5bc632993ed08000000000000000000000000000000000000000000000000015bcf880382c000000000000000000000000000665785a800593e8fa915208c1ce62f6e57fd75ba0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000001117000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004f588000000000000000000000000000000000000000000000000000000000000c350000000000000000000000000000000000000000000000000000000000000021c00000000000000000000000000000000000000000000000000000000000008866c636f64654c6f636174696f6ec258200000000000000000000000000000000000000000000000000000000000000000686c616e6775616765c25820000000000000000000000000000000000000000000000000000000000000000066736f757263657907d0633836366665643238326533313137636466303836633934396662613133643834666331376131656335353934656361643034353133646632326137623538356333363763633132326236373138306334383737303435616235383033373463353066313862346564386132346131323437383532363731623030633035663237373163663036363632333663333236393939323139363866323833346438626462616266306661643165313237613837643237363936323831643965656539326134646263316337356137316136656333613135356438633230616661643064623432383362613433353736303734653035633433633561653061656466643332323838346536613231386466323430323630316436356437316131303061633065376563643037663565646364633535643562373932646130626632353665623038363139336463376431333965613764373965653531653831356465333834386565643363366330353837393265366461333434363738626436373239346636643639656564356132663836323835343965616530323235323835346232666361333635646265623032383433386537326234383465383864316136646563373933633739656265353834666465363465663831383363313365386231623735663037636532303963393138633532643637613735343862653236366433663964316439656132613162303166633838376231316162383739663164333861373833303563373031316533643938346130393863663634383931316536653065383038396365306130363230393136663134323935343036336630376239343931326435666331393366303138633764616135363136323562313966376463323036663930353365623234643036323234616164326338623430646162663631656166666635326234653831373239353837333830313561643730663739316663643864333739343035353737393563383937363164636665333639373938373437353439633234643530646464303563623337613465613863353162306530313032363738643433653766306563353039653434633564343764353335626261363831303936383264643864653439326532363633646336653133653532383539663664336565306533633430336236366362653338643236366137356163373639363863613465653331396166363965373431333137393162653630376537353832373430366164653038306335623239653665343262386563386137373761663865383166336234616337626263666531643066616633393338613664353061316561633835643933643234343066313863333037356237306433626134663930323836396439383937663266636562626262366263646439333436633336633663643838626434336265306562333134323562343665613765386338336638386230363933343836383666366134313839623535666132666431396634326264333730313634616339356530303635656461663130373761633131366632393930303833616631333839636661666336613433323439376531363437393762633738616633366335613435366136646661326636626430626639326136613930366130653930313130626266323265613066333163663364353132663466303331653236343330633831663935656431323362323938356266623830623161396432646337306232356264613961386261303839323833666166663634383661316231646235613938353564346237363966623835663531353063393935306462303964373536326537353133633234653531636163366634366634633231636234373561613937363166666466626434656138613531626465613432383037313466363538393630656336643139656539373237626339316635313665346466306665346264613762623035343161393462326334396636323938616132396337656130646662653635346632306437663164323239633066303262356535326137363031376237306439383232643533383166623966613166393361353861376338383632326631326462643363623937323363626132313639633337643538303939336333663666393065323039336331336130363132323334303064393731363031656262313631343332613966666333373033396562663537326364326566666635636562323539346236346462336261616431633734663532653938343938353964383363313238353465376263393764363432363464653931343735386333386438383739343132333937653263643534653431366234373962363331623830626633306266653062366239353564393066356362303435346361373531303963393938366330636536316165356566376534653433353036313432633633646235363862383634353139623463306636366137633161376661336538666431323231376666336665383164663830643138386232646334343833356132663332323733666133353139633531343764643233353763326161346336326461386238353232306535386130333565373662633133316634623734376632663731643263663933376431303832356138316533623963323136663962316134646431663239383463656635656363656265353530363662363061373263363063323864303336653766386635323131343735386638326366323330646636363930636364617267739f64617267316461726732ff6f736563726574734c6f636174696f6ec2582000000000000000000000000000000000000000000000000000000000000000016773656372657473430102030000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	topic0, err := hex.DecodeString("bf50768ccf13bd0110ca6d53a9c4f1f3271abdd4c24a56878863ed25b20598ff")
	require.NoError(t, err)
	// Create a random requestID
	topic1 := make([]byte, 32)
	_, err = rand.Read(topic1)
	require.NoError(t, err)
	topic2, err := hex.DecodeString("000000000000000000000000665785a800593e8fa915208c1ce62f6e57fd75ba")
	require.NoError(t, err)
	return logpoller.Log{
		Topics: [][]byte{topic0, topic1, topic2},
		Data:   data,
	}
}

func getMockedRequestLogV2(t *testing.T) logpoller.Log {
	// NOTE: Change this to be a more readable log generation
	data, err := hex.DecodeString("0000000000000000000000007e5f4552091a69125d5dfcb7b8c2659029395bdf00000000000000000000000000000000000000000000000000000000000000010000000000000000000000007e5f4552091a69125d5dfcb7b8c2659029395bdf000000000000000000000000000000000000000000000000000000000000026000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000157c0e93aa2d812e96736b4fa41a2e4688e841879a42eead6f06eef541bea9ea289b000000000000000000000000b9816fc57977d5a786e654c7cf76767be63b966e00000000000000000000000000000000000000000000000003fa7025a86e0db80000000000000000000000005cf7f96627f3c9903763d128a1cc5d97556a6b990000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000157c000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000ecd8fae906aaaa0000000000000000000000000000000000000000000000000000000000019a280000000000000000000000000000000000000000000000000000000000016ef6000000000000000000000000000000000000000000000000000000004996041c00000000000000000000000000000000000000000000000000ecd8fae906aaaa00000000000000000000000000000000000000000000000000000000000000796c636f64654c6f636174696f6ec258200000000000000000000000000000000000000000000000000000000000000000686c616e6775616765c25820000000000000000000000000000000000000000000000000000000000000000066736f757263657572657475726e202768656c6c6f20776f726c64273b00000000000000")
	require.NoError(t, err)
	topic0, err := hex.DecodeString("718684b6c135c1277575a7b5c7365bc9587d5ebfd899230d5fa11360f6143bfb")
	require.NoError(t, err)
	// Create a random requestID
	topic1 := make([]byte, 32)
	_, err = rand.Read(topic1)
	require.NoError(t, err)
	topic2, err := hex.DecodeString("000000000000000000000000665785a800593e8fa915208c1ce62f6e57fd75ba")
	require.NoError(t, err)
	return logpoller.Log{
		Topics: [][]byte{topic0, topic1, topic2},
		Data:   data,
	}
}

func TestLogPollerWrapper_SingleSubscriberEmptyEvents_CoordinatorV1(t *testing.T) {
	lp, lpWrapper, client := setUp(t, 100_000) // check only once
	lp.On("LatestBlock").Return(logpoller.LogPollerBlock{BlockNumber: int64(100)}, nil)

	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{}, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getContractById
		To:   &routerAddressHex,
		Data: []uint8{0xa9, 0xc9, 0xa9, 0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(coordinatorAddressBytes, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getProposedContractById
		To:   &routerAddressHex,
		Data: []uint8{0x6a, 0x22, 0x15, 0xde, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(addr("00"))
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	typeAndVersionResponse, err := encodeTypeAndVersion(CoordinatorContractV100)
	require.NoError(t, err)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
		To:   &coordinatorAddressHex,
		Data: hexutil.MustDecode("0x181f5a77"),
	}, mock.Anything).Return(typeAndVersionResponse, nil)
	subscriber := newSubscriber(1)
	lpWrapper.SubscribeToUpdates("mock_subscriber", subscriber)

	servicetest.Run(t, lpWrapper)
	subscriber.updates.Wait()
	reqs, resps, err := lpWrapper.LatestEvents()
	require.NoError(t, err)
	require.Equal(t, 0, len(reqs))
	require.Equal(t, 0, len(resps))
}

func TestLogPollerWrapper_SingleSubscriberEmptyEvents_CoordinatorV2(t *testing.T) {
	lp, lpWrapper, client := setUp(t, 100_000) // check only once
	lp.On("LatestBlock").Return(logpoller.LogPollerBlock{BlockNumber: int64(100)}, nil)

	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{}, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getContractById
		To:   &routerAddressHex,
		Data: []uint8{0xa9, 0xc9, 0xa9, 0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(coordinatorAddressBytes, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getProposedContractById
		To:   &routerAddressHex,
		Data: []uint8{0x6a, 0x22, 0x15, 0xde, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(addr("00"))
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	typeAndVersionResponse, err := encodeTypeAndVersion(CoordinatorContractV200)
	require.NoError(t, err)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
		To:   &coordinatorAddressHex,
		Data: hexutil.MustDecode("0x181f5a77"),
	}, mock.Anything).Return(typeAndVersionResponse, nil)
	subscriber := newSubscriber(1)
	lpWrapper.SubscribeToUpdates("mock_subscriber", subscriber)

	servicetest.Run(t, lpWrapper)
	subscriber.updates.Wait()
	reqs, resps, err := lpWrapper.LatestEvents()
	require.NoError(t, err)
	require.Equal(t, 0, len(reqs))
	require.Equal(t, 0, len(resps))
}

func TestLogPollerWrapper_ErrorOnZeroAddresses(t *testing.T) {
	lp, lpWrapper, client := setUp(t, 100_000) // check only once
	lp.On("LatestBlock").Return(logpoller.LogPollerBlock{BlockNumber: int64(100)}, nil)

	client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(addr("00"))

	servicetest.Run(t, lpWrapper)
	_, _, err := lpWrapper.LatestEvents()
	require.Error(t, err)
}

func TestLogPollerWrapper_LatestEvents_ReorgHandlingV1(t *testing.T) {
	lp, lpWrapper, client := setUp(t, 100_000)
	lp.On("LatestBlock").Return(logpoller.LogPollerBlock{BlockNumber: int64(100)}, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getContractById
		To:   &routerAddressHex,
		Data: []uint8{0xa9, 0xc9, 0xa9, 0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(coordinatorAddressBytes, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getProposedContractById
		To:   &routerAddressHex,
		Data: []uint8{0x6a, 0x22, 0x15, 0xde, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(addr("00"))
	typeAndVersionResponse, err := encodeTypeAndVersion(CoordinatorContractV100)
	require.NoError(t, err)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
		To:   &coordinatorAddressHex,
		Data: hexutil.MustDecode("0x181f5a77"),
	}, mock.Anything).Return(typeAndVersionResponse, nil)
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	subscriber := newSubscriber(1)
	lpWrapper.SubscribeToUpdates("mock_subscriber", subscriber)
	mockedLog := getMockedRequestLogV1(t)
	// All logPoller queries for responses return none
	lp.On("Logs", mock.Anything, mock.Anything, functions_coordinator_1_1_0.FunctionsCoordinator110OracleResponse{}.Topic(), mock.Anything).Return([]logpoller.Log{}, nil)
	// On the first logPoller query for requests, the request log appears
	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{mockedLog}, nil).Once()
	// On the 2nd query, the request log disappears
	lp.On("Logs", mock.Anything, mock.Anything, functions_coordinator_1_1_0.FunctionsCoordinator110OracleRequest{}.Topic(), mock.Anything).Return([]logpoller.Log{}, nil).Once()
	// On the 3rd query, the original request log appears again
	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{mockedLog}, nil).Once()

	servicetest.Run(t, lpWrapper)
	subscriber.updates.Wait()

	oracleRequests, _, err := lpWrapper.LatestEvents()
	require.NoError(t, err)
	assert.Equal(t, 1, len(oracleRequests))
	oracleRequests, _, err = lpWrapper.LatestEvents()
	require.NoError(t, err)
	assert.Equal(t, 0, len(oracleRequests))
	require.NoError(t, err)
	oracleRequests, _, err = lpWrapper.LatestEvents()
	require.NoError(t, err)
	assert.Equal(t, 0, len(oracleRequests))
}

func TestLogPollerWrapper_LatestEvents_ReorgHandlingV2(t *testing.T) {
	lp, lpWrapper, client := setUp(t, 100_000)
	lp.On("LatestBlock").Return(logpoller.LogPollerBlock{BlockNumber: int64(100)}, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getContractById
		To:   &routerAddressHex,
		Data: []uint8{0xa9, 0xc9, 0xa9, 0x18, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(coordinatorAddressBytes, nil)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // getProposedContractById
		To:   &routerAddressHex,
		Data: []uint8{0x6a, 0x22, 0x15, 0xde, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	}, mock.Anything).Return(addr("00"))
	typeAndVersionResponse, err := encodeTypeAndVersion(CoordinatorContractV200)
	require.NoError(t, err)
	client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
		To:   &coordinatorAddressHex,
		Data: hexutil.MustDecode("0x181f5a77"),
	}, mock.Anything).Return(typeAndVersionResponse, nil)
	lp.On("RegisterFilter", mock.Anything).Return(nil)
	subscriber := newSubscriber(1)
	lpWrapper.SubscribeToUpdates("mock_subscriber", subscriber)
	mockedLog := getMockedRequestLogV2(t)
	// All logPoller queries for responses return none
	lp.On("Logs", mock.Anything, mock.Anything, functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(), mock.Anything).Return([]logpoller.Log{}, nil)
	// On the first logPoller query for requests, the request log appears
	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{mockedLog}, nil).Once()
	// On the 2nd query, the request log disappears
	lp.On("Logs", mock.Anything, mock.Anything, functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(), mock.Anything).Return([]logpoller.Log{}, nil).Once()
	// On the 3rd query, the original request log appears again
	lp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]logpoller.Log{mockedLog}, nil).Once()

	servicetest.Run(t, lpWrapper)
	subscriber.updates.Wait()

	oracleRequests, _, err := lpWrapper.LatestEvents()
	require.NoError(t, err)
	assert.Equal(t, 1, len(oracleRequests))
	oracleRequests, _, err = lpWrapper.LatestEvents()
	require.NoError(t, err)
	assert.Equal(t, 0, len(oracleRequests))
	require.NoError(t, err)
	oracleRequests, _, err = lpWrapper.LatestEvents()
	require.NoError(t, err)
	assert.Equal(t, 0, len(oracleRequests))
}

func TestLogPollerWrapper_FilterPreviouslyDetectedEvents_TruncatesLogs(t *testing.T) {
	_, lpWrapper, _ := setUp(t, 100_000)

	inputLogs := make([]logpoller.Log, maxLogsToProcess+100)
	for i := 0; i < 1100; i++ {
		inputLogs[i] = getMockedRequestLogV1(t)
	}

	functionsLpWrapper := lpWrapper.(*logPollerWrapper)
	mockedDetectedEvents := detectedEvents{isPreviouslyDetected: make(map[[32]byte]struct{})}
	outputLogs := functionsLpWrapper.filterPreviouslyDetectedEvents(inputLogs, &mockedDetectedEvents, "request")

	assert.Equal(t, maxLogsToProcess, len(outputLogs))
	assert.Equal(t, 1000, len(mockedDetectedEvents.detectedEventsOrdered))
	assert.Equal(t, 1000, len(mockedDetectedEvents.isPreviouslyDetected))
}

func TestLogPollerWrapper_FilterPreviouslyDetectedEvents_SkipsInvalidLog(t *testing.T) {
	_, lpWrapper, _ := setUp(t, 100_000)
	inputLogs := []logpoller.Log{getMockedRequestLogV1(t)}
	inputLogs[0].Topics = [][]byte{[]byte("invalid topic")}
	mockedDetectedEvents := detectedEvents{isPreviouslyDetected: make(map[[32]byte]struct{})}

	functionsLpWrapper := lpWrapper.(*logPollerWrapper)
	outputLogs := functionsLpWrapper.filterPreviouslyDetectedEvents(inputLogs, &mockedDetectedEvents, "request")

	assert.Equal(t, 0, len(outputLogs))
	assert.Equal(t, 0, len(mockedDetectedEvents.detectedEventsOrdered))
	assert.Equal(t, 0, len(mockedDetectedEvents.isPreviouslyDetected))
}

func TestLogPollerWrapper_FilterPreviouslyDetectedEvents_FiltersPreviouslyDetectedEvent(t *testing.T) {
	_, lpWrapper, _ := setUp(t, 100_000)
	mockedRequestLog := getMockedRequestLogV1(t)
	inputLogs := []logpoller.Log{mockedRequestLog}
	var mockedRequestId [32]byte
	copy(mockedRequestId[:], mockedRequestLog.Topics[1])

	mockedDetectedEvents := detectedEvents{
		isPreviouslyDetected:  make(map[[32]byte]struct{}),
		detectedEventsOrdered: make([]detectedEvent, 1),
	}
	mockedDetectedEvents.isPreviouslyDetected[mockedRequestId] = struct{}{}
	mockedDetectedEvents.detectedEventsOrdered[0] = detectedEvent{
		requestId:    mockedRequestId,
		timeDetected: time.Now().Add(-time.Second * time.Duration(logPollerCacheDurationSecDefault+1)),
	}

	functionsLpWrapper := lpWrapper.(*logPollerWrapper)
	outputLogs := functionsLpWrapper.filterPreviouslyDetectedEvents(inputLogs, &mockedDetectedEvents, "request")

	assert.Equal(t, 0, len(outputLogs))
	// Ensure that expired events are removed from the cache
	assert.Equal(t, 0, len(mockedDetectedEvents.detectedEventsOrdered))
	assert.Equal(t, 0, len(mockedDetectedEvents.isPreviouslyDetected))
}

func encodeTypeAndVersion(typeAndVersion string) ([]byte, error) {
	stringAbiType, _ := abi.NewType("string", "string", nil)
	abiDec := abi.Arguments{
		{Type: stringAbiType},
	}
	nameEncoded, err := abiDec.Pack(typeAndVersion)
	if err != nil {
		return nil, err
	}
	return nameEncoded, nil
}
