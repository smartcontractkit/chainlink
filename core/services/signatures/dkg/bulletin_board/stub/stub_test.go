package stub

import (
	"chainlink/core/services/signatures/cryptotest"
	board "chainlink/core/services/signatures/dkg/bulletin_board"
	"crypto/cipher"
	"testing"

	"github.com/stretchr/testify/require"
)

var numParticipants = 5

type Participant struct {
	key       board.SecretKey
	publicKey board.PublicKey
	boards    board.Boards
}

func makeParticipant(s cipher.Stream) Participant {
	rv := Participant{}
	rv.key = board.PickKey(s)
	rv.publicKey = rv.key.PublicKey()
	rv.boards = Boards(rv.publicKey)
	rv.boards.MakeOwnBoard(rv.key)
	return rv
}

func TestStub_Setup(t *testing.T) {
	randomStreamScalar := cryptotest.NewStream(t, 0)
	var participants []Participant
	for i := 0; i < numParticipants; i++ {
		participants = append(participants, makeParticipant(randomStreamScalar))
	}
	count := 0
	inc := func(_ board.Key)
	for i, participant := range participants {
		for j := 0; j < numParticipants; j++ {
			if j != i {
				oboard, err := participant.boards.MakeBoard(participants[j].publicKey)
				require.NoError(t, err)
				oboard.Subscribe(board.AllMessages, func)
			}
		}
	}
}
