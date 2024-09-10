package spec

import "reflect"

type SmartContract struct {
	Name      string
	Functions []Function
}

type Function struct {
	Name        string
	Description string
	Inputs      []Arg
	Outputs     []Arg
	ReadOnly    bool
}

type Arg struct {
	Type reflect.Type
	Name string
}

type OnChainProductSpec struct {
	Description    string
	SmartContracts []SmartContract
}
