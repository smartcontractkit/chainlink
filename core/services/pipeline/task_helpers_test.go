package pipeline

func (t *ETHTxTask) SetSpecGasLimit(v uint32) {
	t.specGasLimit = &v
}
