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
    uint256 deviationPercentageThreshold;
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
  event FlaggingComparisonUpdated(
    address indexed aggregator,
    address indexed compoundOracle,
    uint8 decimals,
    uint256 deviationPercentageThreshold
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
    uint256 compoundPercentageDeviationThreshold
  ) 
    public 
    onlyOwner() 
  {
    CompoundAssetDetails memory compDetails = s_comparisons[aggregator];
    compDetails.symbol = compoundSymbol;
    compDetails.decimals = compoundDecimals;
    compDetails.deviationPercentageThreshold = compoundPercentageDeviationThreshold;
    s_comparisons[aggregator] = compDetails;
  }

  function check(address[] memory aggregators) public view returns (address[] memory) {
    address[] memory invalidAggregators = new address[](aggregators.length);
    uint256 invalidCount = 0;
    for (uint256 i = 0; i < aggregators.length; i++) {
      address aggregator = aggregators[i];
      if (!isValid(aggregator)) {
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

  function isValid(address aggregator) private view returns (bool valid) {
    // get aggregator price
    // get the open oracle price
    // convert the open oracle decimals to aggregator decimals
    // check if the difference if within the threshold
      // return true if is / false if not
  }
}
