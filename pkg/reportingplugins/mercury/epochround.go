package mercury

type EpochRound struct {
	Epoch uint32
	Round uint8
}

func (x EpochRound) Less(y EpochRound) bool {
	return x.Epoch < y.Epoch || (x.Epoch == y.Epoch && x.Round < y.Round)
}
