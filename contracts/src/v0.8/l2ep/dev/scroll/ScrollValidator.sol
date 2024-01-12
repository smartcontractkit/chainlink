// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ScrollSequencerUptimeFeedInterface} from "../interfaces/ScrollSequencerUptimeFeedInterface.sol";

import {GasLimitValidator} from "../GasLimitValidator.sol";

import {IL1ScrollMessenger} from "@scroll-tech/contracts/L1/IL1ScrollMessenger.sol";

/// @title ScrollValidator - makes cross chain call to update the Sequencer Uptime Feed on L2
contract ScrollValidator is GasLimitValidator {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ScrollValidator 1.0.0";

  /// @param l1CrossDomainMessengerAddress address the L1CrossDomainMessenger contract address
  /// @param l2UptimeFeedAddr the address of the ScrollSequencerUptimeFeed contract address
  /// @param gasLimit the gasLimit to use for sending a message from L1 to L2
  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    uint32 gasLimit
  ) GasLimitValidator(l1CrossDomainMessengerAddress, l2UptimeFeedAddr, gasLimit) {}

  /// @notice validate method sends an xDomain L2 tx to update Uptime Feed contract on L2.
  /// @dev A message is sent using the L1CrossDomainMessenger. This method is accessed controlled.
  /// @param currentAnswer new aggregator answer - value of 1 considers the sequencer offline.
  function validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) external override checkAccess returns (bool) {
    // Make the xDomain call
    IL1ScrollMessenger(L1_CROSS_DOMAIN_MESSENGER_ADDRESS).sendMessage(
      L2_UPTIME_FEED_ADDR,
      0,
      abi.encodeWithSelector(
        ScrollSequencerUptimeFeedInterface.updateStatus.selector,
        currentAnswer == ANSWER_SEQ_OFFLINE,
        uint64(block.timestamp)
      ),
      s_gasLimit
    );

    // return success
    return true;
  }
}
