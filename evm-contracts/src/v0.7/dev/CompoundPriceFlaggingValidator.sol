// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./ConfirmedOwner.sol";
import "../vendor/SafeMathChainlink.sol";
import "../interfaces/FlagsInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/UniswapAnchoredView.sol";
import "./UpkeepCompatible.sol";

contract CompoundPriceFlaggingValidator is ConfirmedOwner, UpkeepCompatible {
  using SafeMathChainlink for uint256;

  struct CompoundAssetDetails {
    string symbol;
    uint8 decimals;
    // 1        = 100%
    // 10       = 10%
    // 20       = 5%
    // 50       = 2%
    // 100      = 1%
    // 200      = 0.5%
    // 500      = 0.2%
    // 1000     = 0.1%
    uint256 deviationThresholdDenominator;
  }

  FlagsInterface private s_flags;
  UniswapAnchoredView private s_openOracle;
  mapping(address => CompoundAssetDetails) private s_comparisons;

  event OpenOracleAddressUpdated(
    address indexed from,
    address indexed to
  );
  event FlagsAddressUpdated(
    address indexed from,
    address indexed to
  );
  event ComparisonUpdated(
    address indexed aggregator,
    string indexed symbol,
    uint8 decimals,
    uint256 deviationThresholdDenominator
  );
  
  constructor(address flagsAddress, address compoundOracleAddress)
    ConfirmedOwner(msg.sender)
  {
    setFlagsAddress(flagsAddress);
    setOpenOracleAddress(compoundOracleAddress);
  }

  function setOpenOracleAddress(address oracleAddress)
    public
    onlyOwner()
  {
    address previous = address(s_openOracle);
    if (previous != oracleAddress) {
      s_openOracle = UniswapAnchoredView(oracleAddress);
      emit OpenOracleAddressUpdated(previous, oracleAddress);
    }
  }

  /**
   * @notice Updates the flagging contract address for raising flags
   * @param flagsAddress sets the address of the flags contract
   */
  function setFlagsAddress(address flagsAddress)
    public
    onlyOwner()
  {
    address previous = address(s_flags);
    if (previous != flagsAddress) {
      s_flags = FlagsInterface(flagsAddress);
      emit FlagsAddressUpdated(previous, flagsAddress);
    }
  }

  function setThreshold(
    address aggregator,
    string calldata compoundSymbol,
    uint8 compoundDecimals,
    uint256 compoundDeviationThresholdDenominator
  ) 
    public 
    onlyOwner() 
  {
    CompoundAssetDetails memory compDetails = s_comparisons[aggregator];
    compDetails.symbol = compoundSymbol;
    compDetails.decimals = compoundDecimals;
    compDetails.deviationThresholdDenominator = compoundDeviationThresholdDenominator;
    s_comparisons[aggregator] = compDetails;
    emit ComparisonUpdated(
      aggregator,
      compoundSymbol,
      compoundDecimals,
      compoundDeviationThresholdDenominator
    );
  }

  function check(address[] memory aggregators) public view returns (address[] memory) {
    address[] memory invalidAggregators = new address[](aggregators.length);
    uint256 invalidCount = 0;
    for (uint256 i = 0; i < aggregators.length; i++) {
      address aggregator = aggregators[i];
      if (isInvalid(aggregator)) {
        invalidAggregators[invalidCount] = aggregator;
        invalidCount++;
      }
    }

    if (aggregators.length != invalidCount) {
      assembly {
        mstore(invalidAggregators, invalidCount)
      }
    }
    return invalidAggregators;
  }

  function update(address[] memory aggregators) public returns (address[] memory){
    address[] memory invalidAggregators = check(aggregators);
    s_flags.raiseFlags(invalidAggregators);
    return invalidAggregators;
  }

  function checkForUpkeep(bytes calldata data) external view override returns (bool, bytes memory) {
    address[] memory invalidAggregators = check(abi.decode(data, (address[])));
    bool needsUpkeep = (invalidAggregators.length > 0);
    return (needsUpkeep, abi.encode(invalidAggregators));
  }

  function performUpkeep(bytes calldata data) external override {
    update(abi.decode(data, (address[])));
  }


  function isInvalid(address aggregator) private view returns (bool invalid) {
    CompoundAssetDetails memory compDetails = s_comparisons[aggregator];
    // Get aggregator price & decimals
    AggregatorV3Interface priceFeed = AggregatorV3Interface(aggregator);
    (,int256 unsignedPrice,,,) = priceFeed.latestRoundData();
    uint256 aggregatorPrice = uint256(unsignedPrice);
    uint8 decimals = priceFeed.decimals();
    // Get compound price
    uint256 compPrice = s_openOracle.price(compDetails.symbol);

    // Convert prices so they match decimals
    if (decimals > compDetails.decimals) {
      uint8 diff = decimals - compDetails.decimals;
      uint256 multiplier = 10**uint256(diff);
      compPrice = compPrice * multiplier;
    }
    else if (decimals < compDetails.decimals) {
      uint8 diff = compDetails.decimals - decimals;
      uint256 multiplier = 10**uint256(diff);
      aggregatorPrice = aggregatorPrice * multiplier;
    }

    // Deviation amount threshold from the aggregator price
    uint256 deviationAmountThreshold = aggregatorPrice.div(compDetails.deviationThresholdDenominator);

    // Calculate deviation
    uint256 deviation;
    if (aggregatorPrice > compPrice) {
      deviation = aggregatorPrice.sub(compPrice);
    }
    else if (aggregatorPrice < compPrice) {
      deviation = compPrice.sub(aggregatorPrice);
    }
    invalid = (deviation >= deviationAmountThreshold);
  }
}
