// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

import "./UniswapConfig.sol";
import "./UniswapLib.sol";

struct Observation {
  uint timestamp;
  uint acc;
}

contract UniswapAnchoredView is UniswapConfig {
  using FixedPoint for *;

  /// @notice The number of wei in 1 ETH
  uint public constant ethBaseUnit = 1e18;

  /// @notice A common scaling factor to maintain precision
  uint public constant expScale = 1e18;

  /// @notice The Open Oracle Reporter
  address public immutable reporter;

  /// @notice The highest ratio of the new price to the anchor price that will still trigger the price to be updated
  uint public immutable upperBoundAnchorRatio;

  /// @notice The lowest ratio of the new price to the anchor price that will still trigger the price to be updated
  uint public immutable lowerBoundAnchorRatio;

  /// @notice The minimum amount of time in seconds required for the old uniswap price accumulator to be replaced
  uint public immutable anchorPeriod;

  /// @notice Official prices by symbol hash
  mapping(bytes32 => uint) public prices;

  /// @notice Circuit breaker for using anchor price oracle directly, ignoring reporter
  bool public reporterInvalidated;

  /// @notice The old observation for each symbolHash
  mapping(bytes32 => Observation) public oldObservations;

  /// @notice The new observation for each symbolHash
  mapping(bytes32 => Observation) public newObservations;

  /// @notice The event emitted when new prices are posted but the stored price is not updated due to the anchor
  event PriceGuarded(
    string symbol,
    uint reporter,
    uint anchor
  );

  /// @notice The event emitted when the stored price is updated
  event PriceUpdated(
    string symbol,
    uint price
  );

  /// @notice The event emitted when anchor price is updated
  event AnchorPriceUpdated(
    string symbol,
    uint anchorPrice,
    uint oldTimestamp,
    uint newTimestamp
  );

  /// @notice The event emitted when the uniswap window changes
  event UniswapWindowUpdated(
    bytes32 indexed symbolHash,
    uint oldTimestamp,
    uint newTimestamp,
    uint oldPrice,
    uint newPrice
  );

  /// @notice The event emitted when reporter invalidates itself
  event ReporterInvalidated(
    address reporter
  );

  bytes32 constant ethHash = keccak256(abi.encodePacked("ETH"));
  bytes32 constant rotateHash = keccak256(abi.encodePacked("rotate"));

  /**
    * @notice Construct a uniswap anchored view for a set of token configurations
    * @dev Note that to avoid immature TWAPs, the system must run for at least a single anchorPeriod before using.
    * @param reporter_ The reporter whose prices are to be used
    * @param anchorToleranceMantissa_ The percentage tolerance that the reporter may deviate from the uniswap anchor
    * @param anchorPeriod_ The minimum amount of time required for the old uniswap price accumulator to be replaced
    * @param configs The static token configurations which define what prices are supported and how
    */
  constructor(
    address reporter_,
    uint anchorToleranceMantissa_,
    uint anchorPeriod_,
    TokenConfig[] memory configs
  )
    UniswapConfig(configs)
    public
  {
    reporter = reporter_;
    anchorPeriod = anchorPeriod_;

    // Allow the tolerance to be whatever the deployer chooses, but prevent under/overflow (and prices from being 0)
    upperBoundAnchorRatio = anchorToleranceMantissa_ > uint(-1) - 100e16 ? uint(-1) : 100e16 + anchorToleranceMantissa_;
    lowerBoundAnchorRatio = anchorToleranceMantissa_ < 100e16 ? 100e16 - anchorToleranceMantissa_ : 1;

    for (uint i = 0; i < configs.length; i++) {
      TokenConfig memory config = configs[i];
      require(config.baseUnit > 0, "baseUnit must be greater than zero");
      address uniswapMarket = config.uniswapMarket;
      if (config.priceSource == PriceSource.REPORTER) {
        require(uniswapMarket != address(0), "reported prices must have an anchor");
        bytes32 symbolHash = config.symbolHash;
        uint cumulativePrice = currentCumulativePrice(config);
        oldObservations[symbolHash].timestamp = block.timestamp;
        newObservations[symbolHash].timestamp = block.timestamp;
        oldObservations[symbolHash].acc = cumulativePrice;
        newObservations[symbolHash].acc = cumulativePrice;
        emit UniswapWindowUpdated(symbolHash, block.timestamp, block.timestamp, cumulativePrice, cumulativePrice);
      } else {
          require(uniswapMarket == address(0), "only reported prices utilize an anchor");
      }
    }
  }

  /**
    * @notice Get the official price for a symbol
    * @param symbol The symbol to fetch the price of
    * @return Price denominated in USD, with 6 decimals
    */
  function price(
    string memory symbol
  )
    external
    view
    returns (
      uint
    )
  {
    TokenConfig memory config = getTokenConfigBySymbol(symbol);
    return priceInternal(config);
  }

  function priceInternal(
    TokenConfig memory config
  )
    internal
    view
    returns (
      uint
    )
  {
    if (config.priceSource == PriceSource.REPORTER) return prices[config.symbolHash];
    if (config.priceSource == PriceSource.FIXED_USD) return config.fixedPrice;
    if (config.priceSource == PriceSource.FIXED_ETH) {
      uint usdPerEth = prices[ethHash];
      require(usdPerEth > 0, "ETH price not set, cannot convert to dollars");
      return mul(usdPerEth, config.fixedPrice) / ethBaseUnit;
    }
  }

  /**
    * @notice Get the underlying price of a cToken
    * @dev Implements the PriceOracle interface for Compound v2.
    * @param cToken The cToken address for price retrieval
    * @return Price denominated in USD, with 18 decimals, for the given cToken address
    */
  function getUnderlyingPrice(
    address cToken
  )
    external
    view
    returns (
      uint
    )
  {
    TokenConfig memory config = getTokenConfigByCToken(cToken);
    // Comptroller needs prices in the format: ${raw price} * 1e(36 - baseUnit)
    // Since the prices in this view have 6 decimals, we must scale them by 1e(36 - 6 - baseUnit)
    return mul(1e30, priceInternal(config)) / config.baseUnit;
  }

  /**
    * @notice Post open oracle reporter prices, and recalculate stored price by comparing to anchor
    * @dev We let anyone pay to post anything, but only prices from configured reporter will be stored in the view.
    * @param messages The messages to post to the oracle
    * @param signatures The signatures for the corresponding messages
    * @param symbols The symbols to compare to anchor for authoritative reading
    */
  function postPrices(
    bytes[] calldata messages,
    bytes[] calldata signatures,
    string[] calldata symbols
  )
    external
  {
    require(messages.length == signatures.length, "messages and signatures must be 1:1");

    // Save the prices
    for (uint i = 0; i < messages.length; i++) {
      priceData.put(messages[i], signatures[i]);
    }

    uint ethPrice = fetchEthPrice();

    // Try to update the view storage
    for (uint i = 0; i < symbols.length; i++) {
      postPriceInternal(symbols[i], ethPrice);
    }
  }

  function postPriceInternal(
    string memory symbol,
    uint ethPrice
  )
    internal
  {
    TokenConfig memory config = getTokenConfigBySymbol(symbol);
    require(config.priceSource == PriceSource.REPORTER, "only reporter prices get posted");

    bytes32 symbolHash = keccak256(abi.encodePacked(symbol));
    uint reporterPrice = priceData.getPrice(reporter, symbol);
    uint anchorPrice;
    if (symbolHash == ethHash) {
      anchorPrice = ethPrice;
    } else {
      anchorPrice = fetchAnchorPrice(symbol, config, ethPrice);
    }

    if (reporterInvalidated) {
      prices[symbolHash] = anchorPrice;
      emit PriceUpdated(symbol, anchorPrice);
    } else if (isWithinAnchor(reporterPrice, anchorPrice)) {
      prices[symbolHash] = reporterPrice;
      emit PriceUpdated(symbol, reporterPrice);
    } else {
      emit PriceGuarded(symbol, reporterPrice, anchorPrice);
    }
  }

  function isWithinAnchor(
    uint reporterPrice,
    uint anchorPrice
  )
    internal
    view
    returns (
      bool
    )
  {
    if (reporterPrice > 0) {
      uint anchorRatio = mul(anchorPrice, 100e16) / reporterPrice;
      return anchorRatio <= upperBoundAnchorRatio && anchorRatio >= lowerBoundAnchorRatio;
    }
    return false;
  }

  /**
    * @dev Fetches the current token/eth price accumulator from uniswap.
    */
  function currentCumulativePrice(
    TokenConfig memory config
  )
    internal
    view
    returns (
      uint
    )
  {
    (uint cumulativePrice0, uint cumulativePrice1,) = UniswapV2OracleLibrary.currentCumulativePrices(config.uniswapMarket);
    if (config.isUniswapReversed) {
      return cumulativePrice1;
    } else {
      return cumulativePrice0;
    }
  }

  /**
    * @dev Fetches the current eth/usd price from uniswap, with 6 decimals of precision.
    *  Conversion factor is 1e18 for eth/usdc market, since we decode uniswap price statically with 18 decimals.
    */
  function fetchEthPrice()
    internal
    returns (
      uint
    )
  {
      return fetchAnchorPrice("ETH", getTokenConfigBySymbolHash(ethHash), ethBaseUnit);
  }

  /**
    * @dev Fetches the current token/usd price from uniswap, with 6 decimals of precision.
    * @param conversionFactor 1e18 if seeking the ETH price, and a 6 decimal ETH-USDC price in the case of other assets
    */
  function fetchAnchorPrice(
    string memory symbol,
    TokenConfig memory config,
    uint conversionFactor
  )
    internal
    virtual
    returns (
      uint
    )
  {
    (uint nowCumulativePrice, uint oldCumulativePrice, uint oldTimestamp) = pokeWindowValues(config);

    // This should be impossible, but better safe than sorry
    require(block.timestamp > oldTimestamp, "now must come after before");
    uint timeElapsed = block.timestamp - oldTimestamp;

    // Calculate uniswap time-weighted average price
    // Underflow is a property of the accumulators: https://uniswap.org/audit.html#orgc9b3190
    FixedPoint.uq112x112 memory priceAverage = FixedPoint.uq112x112(uint224((nowCumulativePrice - oldCumulativePrice) / timeElapsed));
    uint rawUniswapPriceMantissa = priceAverage.decode112with18();
    uint unscaledPriceMantissa = mul(rawUniswapPriceMantissa, conversionFactor);
    uint anchorPrice;

    // Adjust rawUniswapPrice according to the units of the non-ETH asset
    // In the case of ETH, we would have to scale by 1e6 / USDC_UNITS, but since baseUnit2 is 1e6 (USDC), it cancels

    // In the case of non-ETH tokens
    // a. pokeWindowValues already handled uniswap reversed cases, so priceAverage will always be Token/ETH TWAP price.
    // b. conversionFactor = ETH price * 1e6
    // unscaledPriceMantissa = priceAverage(token/ETH TWAP price) * expScale * conversionFactor
    // so ->
    // anchorPrice = priceAverage * tokenBaseUnit / ethBaseUnit * ETH_price * 1e6
    //             = priceAverage * conversionFactor * tokenBaseUnit / ethBaseUnit
    //             = unscaledPriceMantissa / expScale * tokenBaseUnit / ethBaseUnit
    anchorPrice = mul(unscaledPriceMantissa, config.baseUnit) / ethBaseUnit / expScale;

    emit AnchorPriceUpdated(symbol, anchorPrice, oldTimestamp, block.timestamp);

    return anchorPrice;
  }

  /**
    * @dev Get time-weighted average prices for a token at the current timestamp.
    *  Update new and old observations of lagging window if period elapsed.
    */
  function pokeWindowValues(
    TokenConfig memory config
  )
    internal
    returns (
      uint,
      uint,
      uint
    )
  {
    bytes32 symbolHash = config.symbolHash;
    uint cumulativePrice = currentCumulativePrice(config);

    Observation memory newObservation = newObservations[symbolHash];

    // Update new and old observations if elapsed time is greater than or equal to anchor period
    uint timeElapsed = block.timestamp - newObservation.timestamp;
    if (timeElapsed >= anchorPeriod) {
      oldObservations[symbolHash].timestamp = newObservation.timestamp;
      oldObservations[symbolHash].acc = newObservation.acc;

      newObservations[symbolHash].timestamp = block.timestamp;
      newObservations[symbolHash].acc = cumulativePrice;
      emit UniswapWindowUpdated(config.symbolHash, newObservation.timestamp, block.timestamp, newObservation.acc, cumulativePrice);
    }
    return (cumulativePrice, oldObservations[symbolHash].acc, oldObservations[symbolHash].timestamp);
  }

  /**
    * @notice Invalidate the reporter, and fall back to using anchor directly in all cases
    * @dev Only the reporter may sign a message which allows it to invalidate itself.
    *  To be used in cases of emergency, if the reporter thinks their key may be compromised.
    * @param message The data that was presumably signed
    * @param signature The fingerprint of the data + private key
    */
  function invalidateReporter(
    bytes memory message,
    bytes memory signature
  )
    external
  {
      (string memory decodedMessage, ) = abi.decode(message, (string, address));
      require(keccak256(abi.encodePacked(decodedMessage)) == rotateHash, "invalid message must be 'rotate'");
      require(source(message, signature) == reporter, "invalidation message must come from the reporter");
      reporterInvalidated = true;
      emit ReporterInvalidated(reporter);
  }

  /**
    * @notice Recovers the source address which signed a message
    * @dev Comparing to a claimed address would add nothing,
    *  as the caller could simply perform the recover and claim that address.
    * @param message The data that was presumably signed
    * @param signature The fingerprint of the data + private key
    * @return The source address which signed the message, presumably
    */
  function source(
    bytes memory message,
    bytes memory signature
  )
    public
    pure
    returns (
      address
    )
  {
      (bytes32 r, bytes32 s, uint8 v) = abi.decode(signature, (bytes32, bytes32, uint8));
      bytes32 hash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", keccak256(message)));
      return ecrecover(hash, v, r, s);
  }

  /// @dev Overflow proof multiplication
  function mul(
    uint a,
    uint b
  )
    internal
    pure
    returns (
      uint
    )
  {
      if (a == 0) return 0;
      uint c = a * b;
      require(c / a == b, "multiplication overflow");
      return c;
  }
}
