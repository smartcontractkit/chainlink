package protocol

type EpochRound struct {
	epoch uint32
	round uint8
}

func (x EpochRound) Less(y EpochRound) bool {
	return x.epoch < y.epoch || (x.epoch == y.epoch && x.round < y.round)
}
