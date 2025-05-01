package config

type UpdateAptosChainConfig struct {
	ChainSelector uint64
	UpdateCCIP    bool
	UpdateOffRamp bool
	UpdateOnRamp  bool
	UpdateRouter  bool
}
