package changeset

type TokenSymbol string

const (
	LinkSymbol   TokenSymbol = "LINK"
	WethSymbol   TokenSymbol = "WETH"
	USDCSymbol   TokenSymbol = "USDC"
	USDCName     string      = "USD Coin"
	LinkDecimals             = 18
	WethDecimals             = 18
	UsdcDecimals             = 6
)
