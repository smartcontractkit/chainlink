// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {AggregatorValidatorInterface} from "../../../shared/interfaces/AggregatorValidatorInterface.sol";
import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {ZKSyncSequencerUptimeFeedInterface} from "./../interfaces/ZKSyncSequencerUptimeFeedInterface.sol";

import {SimpleWriteAccessController} from "../../../shared/access/SimpleWriteAccessController.sol";

import {IBridgehub, L2TransactionRequestDirect} from "@zksync/contracts/l1-contracts/contracts/bridgehub/IBridgehub.sol";

/**
 * @title ZKSyncValidator - makes cross chain call to update the Sequencer Uptime Feed on L2
 */
contract ZKSyncValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  int256 private constant ANSWER_SEQ_OFFLINE = 1;
  uint32 private s_gasLimit;

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L2_UPTIME_FEED_ADDR;

  /**
   * @notice emitted when gas cost to spend on L2 is updated
   * @param gasLimit updated gas cost
   */
  event GasLimitUpdated(uint32 gasLimit);

  /**
   * @param l1CrossDomainMessengerAddress address the Bridgehub contract address
   * @param l2UptimeFeedAddr the address of the OptimismSequencerUptimeFeed contract address
   * @param gasLimit the gasLimit to use for sending a message from L1 to L2
   */
  constructor(address l1CrossDomainMessengerAddress, address l2UptimeFeedAddr, uint32 gasLimit) {
    // solhint-disable-next-line gas-custom-errors
    require(l1CrossDomainMessengerAddress != address(0), "Invalid xDomain Messenger address");
    // solhint-disable-next-line gas-custom-errors
    require(l2UptimeFeedAddr != address(0), "Invalid OptimismSequencerUptimeFeed contract address");
    L1_CROSS_DOMAIN_MESSENGER_ADDRESS = l1CrossDomainMessengerAddress;
    L2_UPTIME_FEED_ADDR = l2UptimeFeedAddr;
    s_gasLimit = gasLimit;
  }

  /**
   * @notice versions:
   *
   * - ZKSyncValidator 0.1.0: initial release
   * - ZKSyncValidator 1.0.0: change target of L2 sequencer status update
   *   - now calls `updateStatus` on an L2 OptimismSequencerUptimeFeed contract instead of
   *     directly calling the Flags contract
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ZKSyncValidator 1.0.0";
  }

  /**
   * @notice sets the new gas cost to spend when sending cross chain message
   * @param gasLimit the updated gas cost
   */
  function setGasLimit(uint32 gasLimit) external onlyOwner {
    s_gasLimit = gasLimit;
    emit GasLimitUpdated(gasLimit);
  }

  /**
   * @notice fetches the gas cost of sending a cross chain message
   */
  function getGasLimit() external view returns (uint32) {
    return s_gasLimit;
  }

  /**
   * @notice validate method sends an xDomain L2 tx to update Uptime Feed contract on L2.
   * @dev A message is sent using the Bridgehub. This method is accessed controlled.
   * @param currentAnswer new aggregator answer - value of 1 considers the sequencer offline.
   */
  function validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) external override checkAccess returns (bool) {
    // Encode the OptimismSequencerUptimeFeed call
    bytes4 selector = ZKSyncSequencerUptimeFeedInterface.updateStatus.selector;
    bool status = currentAnswer == ANSWER_SEQ_OFFLINE;
    uint64 timestamp = uint64(block.timestamp);
    // Encode `status` and `timestamp`
    bytes memory message = abi.encodeWithSelector(selector, status, timestamp);
    bytes[] memory emptyBytes;

    // Create the L2 transaction request
    L2TransactionRequestDirect memory l2TransactionRequestDirect = L2TransactionRequestDirect({
      chainId: 1, // TODO look at docs
      mintValue: 0,
      l2Contract: L2_UPTIME_FEED_ADDR,
      l2Value: 0, // TODO confirm
      l2Calldata: message,
      l2GasLimit: s_gasLimit,
      l2GasPerPubdataByteLimit: 0, // TODO confirm
      factoryDeps: emptyBytes,
      refundRecipient: msg.sender
    });

    // Make the xDomain call
    IBridgehub(L1_CROSS_DOMAIN_MESSENGER_ADDRESS).requestL2TransactionDirect(l2TransactionRequestDirect);

    return true;
  }
}
