// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/AggregatorV3Interface.sol";

/**
 * Network: Fantom Testnet
 * Base: LINK/USD
 * Base Address: 0x6d5689Ad4C1806D1BA0c70Ab95ebe0Da6B204fC5
 * Quote: FTM/USD
 * Quote Address: 0xe04676B9A9A2973BCb0D1478b5E1E9098BBB7f3D
 * Decimals: 18
 *
 * Network: AVAX Testnet
 * Base: LINK/USD
 * Base Address: 0x34C4c526902d88a3Aa98DB8a9b802603EB1E3470
 * Quote: AVAX/USD
 * Quote Address: 0x5498BB86BC934c8D34FDA08E81D444153d0D06aD
 * Decimals: 18
 *
 * Chainlink Data Feeds can be used in combination to derive denominated price pairs in other
 * currencies.
 *
 * If you require a denomination other than what is provided, you can use two data feeds to derive
 * the pair that you need.
 *
 * For example, if you needed a LINK / FTM price, you could take the LINK / USD feed and the
 * FTM / USD feed and derive LINK / FTM using division.
 * (LINK/USD)/(FTM/USD) = LINK/FTM
 */
contract DerivedPriceFeed is AggregatorV3Interface {
  uint256 public constant override version = 0;

  AggregatorV3Interface public immutable BASE;
  AggregatorV3Interface public immutable QUOTE;
  uint8 public immutable DECIMALS;

  constructor(
    address _base,
    address _quote,
    uint8 _decimals
  ) {
    require(_decimals > uint8(0) && _decimals <= uint8(18), "Invalid _decimals");
    DECIMALS = _decimals;
    BASE = AggregatorV3Interface(_base);
    QUOTE = AggregatorV3Interface(_quote);
  }

  function decimals() external view override returns (uint8) {
    return DECIMALS;
  }

  function getRoundData(uint80)
    external
    pure
    override
    returns (
      uint80,
      int256,
      uint256,
      uint256,
      uint80
    )
  {
    revert("not implemented - use latestRoundData()");
  }

  function description() external pure override returns (string memory) {
    return "DerivedPriceFeed.sol";
  }

  function latestRoundData()
    external
    view
    override
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return (uint80(0), getDerivedPrice(), block.timestamp, block.timestamp, uint80(0));
  }

  // https://docs.chain.link/docs/get-the-latest-price/#getting-a-different-price-denomination
  function getDerivedPrice() internal view returns (int256) {
    (, int256 basePrice, , , ) = BASE.latestRoundData();
    uint8 baseDecimals = BASE.decimals();
    basePrice = scalePrice(basePrice, baseDecimals, DECIMALS);

    (, int256 quotePrice, , , ) = QUOTE.latestRoundData();
    uint8 quoteDecimals = QUOTE.decimals();
    quotePrice = scalePrice(quotePrice, quoteDecimals, DECIMALS);

    return (basePrice * int256(10**uint256(DECIMALS))) / quotePrice;
  }

  function scalePrice(
    int256 _price,
    uint8 _priceDecimals,
    uint8 _decimals
  ) internal pure returns (int256) {
    if (_priceDecimals < _decimals) {
      return _price * int256(10**uint256(_decimals - _priceDecimals));
    } else if (_priceDecimals > _decimals) {
      return _price / int256(10**uint256(_priceDecimals - _decimals));
    }
    return _price;
  }
}
