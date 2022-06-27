package pipeline

func (t *ETHTxTask) SetGasLimitGwei(v uint32) {
	t.gasLimitGwei = &v
}
