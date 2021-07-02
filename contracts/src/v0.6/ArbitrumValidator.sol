// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0; // Could we use 0.8.0

interface ArbitrumValidator {
  function setMaximumGasPrice(uint32) external;
  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) external returns (bool);
}

import "./interfaces/ArbitrumInboxInterface.sol";
import "./SimpleWriteAccessController.sol";

contract ArbitrumValidator is SimpleWriteAccessController {
  address s_arbitrumInbox;
  address s_flags;
  // Follows: https://eips.ethereum.org/EIPS/eip-1967
  address s_arbitrumFlag = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-offline")) - 1)));

  address s_billingAccessController;
  uint32 s_maximumGasPrice;

  /**
   * @param aggregatorAddress address of the Aggregator using validate
   * @param inboxAddress address of the Arbitrum Inbox L1 contract
   * @param flagAddress address of the Chainlink L2 Flags contract
   * @param billingAccessControllerAddress address of the access controller for managing gas price on Arbitrum
   * @param maximumGasPrice maximum gas price to pay on L2
   */
  constructor(
    address aggregatorAddress,
    address inboxAddress, 
    address flagAddress,
    address billingAccessControllerAddress,
    uint32 maximumGasPrice
  ) 
  {
    s_arbitrumInbox = inboxAddress;
    s_flagsAddress = flagAddress;
    s_maximumGasPrice = maximumGasPrice;
    s_billingAccessController = billingAccessControllerAddress;

    // Default access. Give access to aggregator
    addAccess(aggregatorAddress);
  }

  function setMaximumGasPrice(
    uint32 gasPrice
  )
    external
  {
    SimpleWriteAccessController access = SimpleWriteAccessController(s_billingAccessController);
    require(access.hasAccess(msg.sender, msg.data), "Only billing admin can call");
    _setMaximumGasPriceInternal(gasPrice);
  }


  function validate(
    uint256 previousRoundId,
    int256 previousAnswer,
    uint256 currentRoundId,
    int256 currentAnswer
  ) 
    external 
    hasAccess() 
    returns (bool) 
  {
    bytes data = currentAnswer == 1 ? abi.encodeWithSignature("raiseFlag(address)", arbitrumFlag) : abi.encodeWithSignature("lowerFlags(address[])", [arbitrumFlag]);
    IInbox arbitrumInbox = IInbox(s_arbitrumInbox);
    // Validator should be funded in L1 and send some value to pay for the L2 gas
    arbitrumInbox.sendL1FundedContractTransaction(
      30000000,
      s_maximumGasPrice,
      flagsAddress,
      data
    );
    return true;
  }

  event MaximumGasPriceSet(
    uint32 maximumGasPrice
  );

  function _setMaximumGasPriceInternal(
    uint256 _gasPrice
  )
    internal
  {
    s_maximumGasPrice = _gasPrice;
    emit MaximumGasPriceSet(_gasPrice);
  }
}