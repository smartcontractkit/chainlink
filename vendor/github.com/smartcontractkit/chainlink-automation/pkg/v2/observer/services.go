package observer

type SimpleService struct {
	F func() error
	C func()
}

func (sw *SimpleService) Do() error {
	return sw.F()
}

func (sw *SimpleService) Stop() {
	sw.C()
}
