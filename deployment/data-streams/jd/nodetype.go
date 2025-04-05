package jd

type NodeType int

const (
	NodeTypeOracle NodeType = iota
	NodeTypeBootstrap
)

func (nt NodeType) String() string {
	switch nt {
	case NodeTypeOracle:
		return "oracle"
	case NodeTypeBootstrap:
		return "bootstrap"
	default:
		return "unknown"
	}
}
