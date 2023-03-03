package player_idx

func init() {
	if int(MaxPlayer) < 0 {
		panic("idx must be an unsigned type")
	}
	if one+MaxPlayer != 0 {
		panic("maxPlayer must be twos complement representation of -1")
	}
}

func (pi *PlayerIdx) mustBeNonZero() {
	if err := pi.NonZero(); err != nil {
		panic(err)
	}
}

var one = Int(1)
