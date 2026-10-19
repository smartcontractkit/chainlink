package solana

import (
	"encoding/binary"
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	df1_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/store"
)

func TestParseDataID(t *testing.T) {
	t.Parallel()

	id, err := ParseDataID("0x0003aabbccddeeff0001020300000000")
	require.NoError(t, err)
	require.Equal(t, [16]byte{0x00, 0x03, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x01, 0x02, 0x03}, id)

	noPrefix, err := ParseDataID("0003aabbccddeeff0001020300000000")
	require.NoError(t, err)
	require.Equal(t, id, noPrefix)

	for _, bad := range []string{"0x0003", "0x0003aabbccddeeff000102030000000000", "0xzz03aabbccddeeff0001020300000000", ""} {
		_, err := ParseDataID(bad)
		require.Error(t, err, bad)
	}

	_, err = parseDataIDs([]string{"0x0003aabbccddeeff0001020300000000", "0003aabbccddeeff0001020300000000"})
	require.ErrorContains(t, err, "duplicate")
	_, err = parseDataIDs([]string{"0x00000000000000000000000000000000"})
	require.ErrorContains(t, err, "zero")
}

func TestDataIDLayout(t *testing.T) {
	t.Parallel()

	//            feed id         index svr path instance reserved
	id := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10}
	require.True(t, isSvrDataID(id))
	require.Equal(t, [2]byte{0xa, 0xb}, svrPathID(id))
	require.Equal(t, [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0xc, 0xd, 0xe, 0xf, 0x10}, nonSvrDataID(id))
	require.Equal(t, [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, indexedFeedID(id))
	require.False(t, isSvrDataID(nonSvrDataID(id)))
}

func TestSettingsAccount(t *testing.T) {
	t.Parallel()

	store := storeAddresses{
		ProgramID: solana.MustPublicKeyFromBase58("9qCsUGTdKCanybkroJuPzAhuRU8P2SZN4NCexevSZAwu"),
		State:     solana.MustPublicKeyFromBase58("11111111111111111111111111111114"),
	}
	nonSvr := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	svr := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1}

	feed, err := store.feed(nonSvr)
	require.NoError(t, err)
	settings, err := store.settings(nonSvr)
	require.NoError(t, err)
	require.Equal(t, feed, settings, "a non-SVR data ID's settings live in its feed")

	config, err := store.dataIDConfig(svr)
	require.NoError(t, err)
	settings, err = store.settings(svr)
	require.NoError(t, err)
	require.Equal(t, config, settings, "an SVR path's settings live in its data ID config")

	metas, err := store.settingsMetas([][16]byte{nonSvr, svr})
	require.NoError(t, err)
	require.Equal(t, []*solana.AccountMeta{solana.Meta(feed).WRITE(), solana.Meta(config).WRITE()}, metas)
}

func TestInt128(t *testing.T) {
	t.Parallel()

	for _, s := range []string{"0", "1", "-1", "325012000000000000000", "-325012000000000000000", maxInt128.String(), minInt128.String()} {
		v, err := parseAnswer(s)
		require.NoError(t, err)
		require.Equal(t, s, toInt128(v).BigInt().String())

		le := make([]byte, 16)
		i := toInt128(v)
		binary.LittleEndian.PutUint64(le[:8], i.Lo)
		binary.LittleEndian.PutUint64(le[8:], i.Hi)
		require.Equal(t, s, fromLittleEndianInt128(le).String())
	}

	_, err := parseAnswer(new(big.Int).Add(maxInt128, big.NewInt(1)).String())
	require.Error(t, err)
	_, err = parseAnswer(new(big.Int).Sub(minInt128, big.NewInt(1)).String())
	require.Error(t, err)
	_, err = parseAnswer("1.5")
	require.Error(t, err)
}

func TestRescaleAnswer(t *testing.T) {
	t.Parallel()

	answer, _ := new(big.Int).SetString("3250129999999999999999", 10) // 3250.129999... at 18 decimals
	require.Equal(t, "325012999999", rescaleAnswer(answer, 8).String())
	require.Equal(t, "-325012999999", rescaleAnswer(new(big.Int).Neg(answer), 8).String(), "truncates toward zero")
	require.Equal(t, answer.String(), rescaleAnswer(answer, 18).String())
	require.Equal(t, answer.String()+"00", rescaleAnswer(answer, 20).String())
}

func TestCheckCutoverValue(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	storeAnswer, _ := new(big.Int).SetString("3250120000000000000000", 10) // 3250.12 at 18 decimals
	observedAt := uint32(now.Unix() - 30)

	tests := []struct {
		name         string
		df1Answer    int64
		observedAt   uint32
		maxStaleness time.Duration
		maxBps       uint64
		wantErr      string
	}{
		{name: "equal", df1Answer: 325012000000, observedAt: observedAt, maxStaleness: time.Minute},
		{name: "within tolerance", df1Answer: 325300000000, observedAt: observedAt, maxStaleness: time.Minute, maxBps: 10},
		{name: "outside tolerance", df1Answer: 325400000000, observedAt: observedAt, maxStaleness: time.Minute, maxBps: 10, wantErr: "deviates"},
		{name: "stale", df1Answer: 325012000000, observedAt: observedAt, maxStaleness: 10 * time.Second, wantErr: "old"},
		{name: "zero DF1 answer", df1Answer: 0, observedAt: observedAt, maxStaleness: time.Minute, maxBps: 10_000, wantErr: "deviates"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkCutoverValue(storeAnswer, tc.observedAt, 8, big.NewInt(tc.df1Answer), now, tc.maxStaleness, tc.maxBps)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestDF1FeedLatest(t *testing.T) {
	t.Parallel()

	liveLength := uint32(3)
	data := make([]byte, df1HeaderSize+int(liveLength)*df1TransmissionSize)
	write := func(i int, timestamp uint32, answer int64) {
		tr := data[df1HeaderSize+i*df1TransmissionSize:]
		binary.LittleEndian.PutUint32(tr[8:12], timestamp)
		v := toInt128(big.NewInt(answer))
		binary.LittleEndian.PutUint64(tr[16:24], v.Lo)
		binary.LittleEndian.PutUint64(tr[24:32], v.Hi)
	}
	write(0, 4, 400)
	write(1, 2, 200)
	write(2, 3, -300)

	f := &df1Feed{Header: &df1_store.Transmissions{LatestRoundId: 4, LiveLength: liveLength, LiveCursor: 1}, data: data}
	latest, err := f.latest()
	require.NoError(t, err)
	require.Equal(t, uint32(4), latest.Timestamp)
	require.Equal(t, int64(400), latest.Answer.Int64())

	f.Header.LiveCursor = 0 // wrapped around
	latest, err = f.latest()
	require.NoError(t, err)
	require.Equal(t, int64(-300), latest.Answer.Int64())

	f.Header.LatestRoundId = 0
	_, err = f.latest()
	require.ErrorContains(t, err, "no round")
}

func TestSortedFeedAdmins(t *testing.T) {
	t.Parallel()

	a, b := solana.MustPublicKeyFromBase58("11111111111111111111111111111114"), solana.MustPublicKeyFromBase58("11111111111111111111111111111113")
	sorted, err := sortedFeedAdmins([]solana.PublicKey{a, b})
	require.NoError(t, err)
	require.Equal(t, []solana.PublicKey{b, a}, sorted)

	_, err = sortedFeedAdmins([]solana.PublicKey{a, a})
	require.ErrorContains(t, err, "duplicate")
	_, err = sortedFeedAdmins([]solana.PublicKey{{}})
	require.ErrorContains(t, err, "zero")
}

func TestWorkflowMetadata(t *testing.T) {
	t.Parallel()

	storeProgram := solana.MustPublicKeyFromBase58("9qCsUGTdKCanybkroJuPzAhuRU8P2SZN4NCexevSZAwu")
	forwarder := Sender{
		ProgramID: solana.MustPublicKeyFromBase58("whV7Q5pi17hPPyaPksToDw1nMx6Lh8qmNWKFaLRQ4wz"),
		StateID:   solana.MustPublicKeyFromBase58("11111111111111111111111111111114"),
	}
	metadata, err := workflowMetadata(storeProgram, []StoreWorkflow{{Forwarder: forwarder, Owner: [20]uint8{19: 11}, Name: [10]uint8{9: 10}}})
	require.NoError(t, err)
	require.Len(t, metadata, 1)

	sender, _, err := solana.FindProgramAddress([][]byte{[]byte("forwarder"), forwarder.StateID.Bytes(), storeProgram.Bytes()}, forwarder.ProgramID)
	require.NoError(t, err)
	require.Equal(t, sender, metadata[0].AllowedSender)
	require.Equal(t, [20]uint8{19: 11}, metadata[0].AllowedWorkflowOwner)
	require.Equal(t, [10]uint8{9: 10}, metadata[0].AllowedWorkflowName)

	_, err = workflowMetadata(storeProgram, []StoreWorkflow{{}})
	require.ErrorContains(t, err, "forwarder program and state are required")
	_, err = workflowMetadata(storeProgram, make([]StoreWorkflow, maxWorkflowMetadatas+1))
	require.ErrorContains(t, err, "at most")
}

func TestChunk(t *testing.T) {
	t.Parallel()

	require.Nil(t, chunk([]int{}, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, chunk([]int{1, 2, 3, 4, 5}, 2))
	require.Equal(t, [][]int{{1, 2}}, chunk([]int{1, 2}, 2))
}
