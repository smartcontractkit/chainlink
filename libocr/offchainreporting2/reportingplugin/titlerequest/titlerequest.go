package titlerequest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2titlerequest"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Transactions that are buried confirmationDepth blocks deep in the chain are
// considered confirmed.
const confirmationDepth = 5

// Scan logs at most 2000 blocks old.
const lookback = 2000

// We truncate titles longer than this.
const maxTitleLen = 500

// Duration for which a report is considered pending
const pendingDuration = 2 * time.Minute

var _ types.ReportingPluginFactory = (*TitleRequestPluginFactory)(nil)

type TitleRequestPluginFactory struct {
	Client   *ethclient.Client
	Contract *ocr2titlerequest.OCR2TitleRequest
}

func (fac *TitleRequestPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	return &TitleRequestPlugin{
			config.F,
			fac.Client,
			fac.Contract,
			map[[32]byte]bool{},
			map[[32]byte]time.Time{},
		}, types.ReportingPluginInfo{
			"Title Request ReportingPlugin",
			false,
			types.ReportingPluginLimits{
				// queries are empty
				0,
				// observations are at most 32 (request id) + 32 (title offset) + 32
				// (title len) + maxTitleLen chars, let's generously round to 1000
				1_000,
				// reports follow the same format as observations
				1_000,
			},
		}, nil
}

var _ types.ReportingPlugin = (*TitleRequestPlugin)(nil)

type TitleRequestPlugin struct {
	F        int
	client   *ethclient.Client
	contract *ocr2titlerequest.OCR2TitleRequest
	// Keep track of fulfilled requestIDs to avoid duplicate fulfillments. For a
	// production implementation, you'd want to store this persistently.
	fulfilled map[[32]byte]bool
	// Keep track of pending requestIDs to avoid duplicate fulfillments.
	pending map[[32]byte]time.Time
}

func (trp *TitleRequestPlugin) Query(context.Context, types.ReportTimestamp) (types.Query, error) {
	// We don't use a query for this reporting plugin, so we can just leave it empty here
	return types.Query{}, nil
}

func (trp *TitleRequestPlugin) Observation(ctx context.Context, _ types.ReportTimestamp, _ types.Query) (types.Observation, error) {
	// This is where we start to do interesting things. We look for the oldest
	// unfulfilled non-pending request. If we don't find such a request, we send
	// an empty Observation. Otherwise, we send an observation in the format of
	// the report, i.e. an abi-encoded tuple of (requestID, webpage title).

	blocknumber, err := trp.client.BlockNumber(ctx)
	if err != nil {
		return types.Observation{}, err
	}
	confirmedBlocknumer := blocknumber - confirmationDepth

	if err := trp.updateFulfilled(ctx, blocknumber); err != nil {
		return types.Observation{}, err
	}

	var request *ocr2titlerequest.OCR2TitleRequestTitleRequest
	{
		// Find oldest open request. A production implementation would likely do
		// this in a separate background go routine and not scan overlapping
		// block ranges over and over.
		it, err := trp.contract.FilterTitleRequest(&bind.FilterOpts{
			Start:   blocknumber - lookback,
			End:     &confirmedBlocknumer,
			Context: ctx,
		})
		if err != nil {
			return types.Observation{}, err
		}
		defer it.Close()

		for it.Next() {
			if trp.fulfilled[it.Event.RequestId] || trp.isPending(it.Event.RequestId) {
				continue
			}

			request = it.Event
			break
		}
	}

	if request == nil {
		return types.Observation{}, nil
	}

	report, err := encodeReport(request.RequestId, title(ctx, request.Url))
	if err != nil {
		return types.Observation{}, err
	}

	return types.Observation(report), nil
}

func (trp *TitleRequestPlugin) Report(_ context.Context, _ types.ReportTimestamp, _ types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	// Not the most efficient implementation, but it gets the job done. Find
	// any observation that has been sent by at least F+1 nodes, use it as
	// report.
	for _, ao1 := range aos {
		if _, _, err := decodeReport(types.Report(ao1.Observation)); err != nil {
			continue
		}

		voteCount := 0
		for _, ao2 := range aos {
			if bytes.Equal(ao1.Observation, ao2.Observation) {
				voteCount++
			}
			if voteCount > trp.F {
				// At least F+1 oracles "voted" for the same report. Since we
				// assume that at most F oracles are faulty/malicious, this
				// implies that at least one honest oracle voted for this
				// report.
				return true, types.Report(ao1.Observation), nil
			}
		}
	}

	return false, nil, nil
}

func (trp *TitleRequestPlugin) ShouldAcceptFinalizedReport(_ context.Context, _ types.ReportTimestamp, report types.Report) (bool, error) {
	requestID, _, err := decodeReport(report)
	if err != nil {
		return false, nil
	}

	// Avoid duplicate reports
	if trp.isPending(requestID) {
		return false, nil
	}

	// Mark as pending and accept to transmit
	trp.setPending(requestID)
	return true, nil
}

func (trp *TitleRequestPlugin) ShouldTransmitAcceptedReport(ctx context.Context, _ types.ReportTimestamp, report types.Report) (bool, error) {
	requestID, _, err := decodeReport(report)
	if err != nil {
		return false, nil
	}

	blocknumber, err := trp.client.BlockNumber(ctx)
	if err != nil {
		return false, err
	}

	if err := trp.updateFulfilled(ctx, blocknumber); err != nil {
		return false, err
	}

	// Don't broadcast if the request is already fulfilled
	return !trp.fulfilled[requestID], nil
}

func (trp *TitleRequestPlugin) Close() error {
	// No background go-routines or other resources are held by
	// TitleRequestPlugin, no need to do anything here.
	return nil
}

// Updates the set of fulfilled requests. A production implementation would
// likely do this in a separate background go routine and not scan overlapping
// block ranges over and over.
func (trp *TitleRequestPlugin) updateFulfilled(ctx context.Context, blocknumber uint64) error {
	it, err := trp.contract.FilterTitleFulfillment(&bind.FilterOpts{
		Start:   blocknumber - lookback,
		End:     &blocknumber,
		Context: ctx,
	})
	if err != nil {
		return err
	}
	defer it.Close()

	for it.Next() {
		trp.fulfilled[it.Event.RequestId] = true
	}

	return nil
}

func (trp *TitleRequestPlugin) isPending(requestID [32]byte) bool {
	if t, ok := trp.pending[requestID]; ok && !t.Add(pendingDuration).Before(time.Now()) {
		return true
	}
	return false
}

func (trp *TitleRequestPlugin) setPending(requestID [32]byte) {
	trp.pending[requestID] = time.Now()
}

// Returns title (or an empty string), truncated to maxTitleLen chars. There are
// many things wrong with this function, it's just a quick and dirty prototype!
func title(ctx context.Context, url string) string {
	// This is bad practice, we're GETting an untrusted URL
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ""
	}

	httpResp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return ""
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return ""
	}

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return ""
	}
	bodyString := string(bodyBytes)

	// Obviously not the right way to parse HTML, but hey, at least we don't use
	// regexps ;-)
	titleOpen := strings.Index(bodyString, "<title>")
	titleClose := strings.Index(bodyString, "</title>")

	if !(titleOpen+len("<title>") < titleClose && 0 <= titleOpen && 0 <= titleClose) {
		return ""
	}

	result := bodyString[titleOpen+len("<title>") : titleClose]
	if len(result) > maxTitleLen {
		result = result[:maxTitleLen]
	}
	return result
}

func encodeReport(requestID [32]byte, title string) (types.Report, error) {
	return makeReportArgs().Pack(requestID, title)
}

func decodeReport(report types.Report) (requestID [32]byte, title string, err error) {
	unpacked, err := makeReportArgs().Unpack(report)
	if err != nil {
		return [32]byte{}, "", err
	}
	var ok bool
	requestID, ok = unpacked[0].([32]byte)
	if !ok {
		return [32]byte{}, "", fmt.Errorf("cast to big.Int failed, got %T", unpacked[0])
	}
	title, ok = unpacked[1].(string)
	if !ok {
		return [32]byte{}, "", fmt.Errorf("cast to string failed, got %T", unpacked[1])
	}
	return requestID, title, nil
}

func makeReportArgs() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "requestID", Type: mustNewType("bytes32")},
		{Name: "title", Type: mustNewType("string")},
	})
}
