// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./interfaces/ArbitrumInboxInterface.sol";
import "./interfaces/AggregatorValidatorInterface.sol";
import "./interfaces/FlagsInterface.sol";
import "../v0.6/SimpleWriteAccessController.sol";

contract ArbitrumValidator is SimpleWriteAccessController, AggregatorValidatorInterface {

  bytes4 constant private RAISE_SELECTOR = FlagsInterface.raiseFlag.selector;
  bytes4 constant private LOWER_SELECTOR = FlagsInterface.lowerFlags.selector;

  address private s_flagsAddress;
  // Follows: https://eips.ethereum.org/EIPS/eip-1967
  address private s_arbitrumFlag = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-offline")) - 1)));

  ArbitrumInboxInterface private s_arbitrumInbox;
  SimpleWriteAccessController private s_gasConfigAccessController;

  struct GasConfiguration {
    uint256 maximumSubmissionCost;
    uint32 maximumGasPrice;
    uint256 gasCostL2;
    address refundableAddress;
  }
  GasConfiguration private s_gasConfig;
  uint32 constant private s_L2GasLimit = 30000000;
  uint32 constant private s_maxSubmissionCostIncreaseRatio = 13;

  /**
   * @param aggregatorAddress default aggregator with access to validate
   * @param inboxAddress address of the Arbitrum Inbox L1 contract
   * @param flagAddress address of the Chainlink L2 Flags contract
   * @param gasConfigAccessController address of the access controller for managing gas price on Arbitrum
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maximumGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  constructor(
    address aggregatorAddress,
    address inboxAddress,
    address flagAddress,
    address gasConfigAccessController,
    uint256 maxSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  ) {
    s_arbitrumInbox = ArbitrumInboxInterface(inboxAddress);
    s_gasConfigAccessController = SimpleWriteAccessController(gasConfigAccessController);
    s_flagsAddress = flagAddress;
    _setGasConfiguration(maxSubmissionCost, maximumGasPrice, gasCostL2, refundableAddress);

    SimpleWriteAccessController(address(this)).addAccess(aggregatorAddress);
  }
  
  fallback() external payable {}

  function withdrawFunds() 
    external 
    onlyOwner() 
  {
    address payable to = payable(msg.sender);
    to.transfer(address(this).balance);
  }

  function withdrawFundsTo(
    address payable to
  ) 
    external
    onlyOwner() 
  {
    to.transfer(address(this).balance);
  }

  function setGasConfiguration(
    uint256 maxSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  )
    external
  {
    require(s_gasConfigAccessController.hasAccess(msg.sender, msg.data), "Only gas configuration admin can call");
    _setGasConfiguration(maxSubmissionCost, maximumGasPrice, gasCostL2, refundableAddress);
  }

  function validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) 
    external 
    override
    checkAccess() 
    returns (bool) 
  {
    bytes memory data = currentAnswer == 1 ? abi.encodeWithSelector(RAISE_SELECTOR, s_arbitrumFlag) : abi.encodeWithSelector(LOWER_SELECTOR, [s_arbitrumFlag]);

    s_arbitrumInbox.createRetryableTicket{value: s_gasConfig.gasCostL2}(
      s_flagsAddress, 
      0, // L2 call value
      s_gasConfig.maximumSubmissionCost, // Max submission cost of sending data length
      s_gasConfig.refundableAddress, // excessFeeRefundAddress
      s_gasConfig.refundableAddress, // callValueRefundAddress
      s_L2GasLimit,
      s_gasConfig.maximumGasPrice, 
      data
    );
    return true;
  }
  
  event GasConfigurationSet(
    uint256 maximumSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  );
  
  function _setGasConfiguration(
    uint256 _maximumSubmissionCost,
    uint32 _maximumGasPrice,
    uint256 _gasCostL2,
    address _refundableAddress
  ) internal {
    uint256 minGasCostValue = _maximumSubmissionCost * s_maxSubmissionCostIncreaseRatio + s_L2GasLimit * _maximumGasPrice;
    require(_gasCostL2 >= minGasCostValue, "The gas cost value provided is too low to cover the L2 transactions");
    s_gasConfig = GasConfiguration(_maximumSubmissionCost, _maximumGasPrice, _gasCostL2, _refundableAddress);
    emit GasConfigurationSet(_maximumSubmissionCost, _maximumGasPrice, _gasCostL2, _refundableAddress);
  }
}