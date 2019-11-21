// package stub implements a completely in-memory "network", as a place holder
// for the eventual networking system the threshold-signature scheme will use.
package stub

import (
	board "chainlink/core/services/signatures/dkg/bulletin_board"
)

type BoardMap map[board.BoardKey]board.BoardValue

// network represents the bird's eye view of all the boards combined. This is
// how the nodes communicate, in this stub board. More realistic implementations
// will replace this with some kind of network communication.
var network = make(map[board.PublicKey]*stubBoards)

type stubBoards struct {
	board         BoardMap
	subscriptions []board.Subscription
	listeners     []board.BoardUpdateHandler
}

var _ board.Boards = &stubBoards{} // assert interface implementation

type stubBoard struct {
	publicKey board.PublicKey
}

var _ board.BulletinBoard = &stubBoard{} // assert interface implementation

type ownStubBoard struct {
	publicKey board.PublicKey
}

var _ board.OwnBulletinBoard = &ownStubBoard{} // assert interface implementation

func Boards(key board.PublicKey) board.Boards {
	rv := &stubBoards{}
	network[key] = rv
	rv.board = make(BoardMap)
	return rv
}

// MakeOwnBoard implements board.Boards.MakeOwnBoard on stubBoards
func (b *stubBoards) MakeOwnBoard(key board.SecretKey) (board.OwnBulletinBoard, error) {
	return &ownStubBoard{publicKey: key.PublicKey()}, nil
}

// MakeBoard implements board.Boards.MakeBoard on stubBoards
func (b *stubBoards) MakeBoard(key board.PublicKey) (board.BulletinBoard, error) {
	return &stubBoard{publicKey: key}, nil
}

// Get implements board.BulletinBoard.Get on stubBoard
func (b *stubBoard) Get(key board.BoardKey) board.BoardValue {
	return network[b.publicKey].board[key]
}

// Subscribe implements board.BulletinBoard.Subscribe on stubBoard
func (b *stubBoard) Subscribe(key board.Subscription, handler board.SubscriptionHandler) {
	board := network[b.publicKey]
	board.subscriptions = append(board.subscriptions, key)
	board.listeners = append(board.listeners, handler)
}

// Publish implements board.OwnBulletinBoard.Publish on ownStubBoard
func (b *ownStubBoard) Publish(key board.BoardKey, value board.BoardValue) {
	board := network[b.publicKey]
	board.board[key] = value
	for subidx, subscription := range board.subscriptions {
		if subscription.Match(key) {
			board.listeners[subidx](key, value)
		}
	}
}
