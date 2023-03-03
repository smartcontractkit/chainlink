package median

type epochRound struct {
	Epoch uint32
	Round uint8
}

func (x epochRound) Less(y epochRound) bool {
	return x.Epoch < y.Epoch || (x.Epoch == y.Epoch && x.Round < y.Round)
}
