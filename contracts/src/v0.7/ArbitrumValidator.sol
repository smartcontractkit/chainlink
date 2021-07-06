// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0; // Could we use 0.8.0

import "./interfaces/ArbitrumInboxInterface.sol";
import "./interfaces/AggregatorValidatorInterface.sol";
import "../v0.6/SimpleWriteAccessController.sol";

contract ArbitrumValidator is SimpleWriteAccessController, AggregatorValidatorInterface {
  address s_arbitrumInbox;
  address s_flags;
  // Follows: https://eips.ethereum.org/EIPS/eip-1967
  address s_arbitrumFlag = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-offline")) - 1)));
  address s_gasConfigAccessController;
  address s_refundableAddress;

  struct GasConfiguration {
    uint32 maximumGasPrice;
    uint256 gasCostL2;
  }
  GasConfiguration internal s_gasConfig;

  /**
   * @param inboxAddress address of the Arbitrum Inbox L1 contract
   * @param flagAddress address of the Chainlink L2 Flags contract
   * @param gasConfigAccessController address of the access controller for managing gas price on Arbitrum
   * @param maximumGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  constructor(
    address inboxAddress, 
    address flagAddress,
    address gasConfigAccessController,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  ) 
    public
  {
    s_arbitrumInbox = inboxAddress;
    s_flags = flagAddress;
    s_gasConfigAccessController = gasConfigAccessController;
    s_refundableAddress = refundableAddress;
    // TODO: Is it possible to give default access to the aggregator?
    // addAccess(aggregatorAddress);
    _setGasConfigurationInternal(maximumGasPrice, gasCostL2);
  }
  
  // Accept ETH funds
  fallback() external payable {}


  function setGasConfiguration(
    uint32 maximumGasPrice,
    uint256 gasCostL2
  )
    external
    override
  {
    SimpleWriteAccessController access = SimpleWriteAccessController(s_gasConfigAccessController);
    require(access.hasAccess(msg.sender, msg.data), "Only billing admin can call");
    _setGasConfigurationInternal(maximumGasPrice, gasCostL2);
  }
  
  function setRefundableAddress(
    address refundableAddress
  ) 
    external 
    override
    checkAccess()
  {
    s_refundableAddress = refundableAddress;
  }

  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) 
    external 
    override
    checkAccess() 
    returns (bool) 
  {
    bytes memory data = currentAnswer == 1 ? abi.encodeWithSignature("raiseFlag(address)", s_arbitrumFlag) : abi.encodeWithSignature("lowerFlags(address[])", [s_arbitrumFlag]);
    IInbox arbitrumInbox = IInbox(s_arbitrumInbox);
    // Validator should be funded in L1 and send some value to pay for the L2 gas
    // uint256 minL2Cost = maxSubmissionCost + (s_gasConfig.maximumGasPrice*gasLimit);
    arbitrumInbox.createRetryableTicket{value: s_gasConfig.gasCostL2}(
      s_flags, 
      0, // L2 call value
      13700320797, // Max submission cost of sending data length
      s_refundableAddress, // excessFeeRefundAddress
      s_refundableAddress, // callValueRefundAddress
      30000000, // L2 gas limit
      s_gasConfig.maximumGasPrice, 
      data
    );
    return true;
  }
  
  event GasConfigurationSet(
    uint32 maximumGasPrice,
    uint256 gasCostL2
  );

  
  function _setGasConfigurationInternal(
    uint32 _maximumGasPrice,
    uint256 _gasCostL2
  ) internal {
      s_gasConfig = GasConfiguration(_maximumGasPrice, _gasCostL2);
      emit GasConfigurationSet(_maximumGasPrice, _gasCostL2);
  }
}