// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../SimpleWriteAccessController.sol";

import "./interfaces/OptimismSequencerUptimeFeedInterface.sol";
import "@eth-optimism/contracts/L1/messaging/IL1CrossDomainMessenger.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";

/**
 * @title OptimismValidator - makes cross chain call to update the Sequencer Uptime Feed on L2
 */
contract OptimismValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  int256 private constant ANSWER_SEQ_OFFLINE = 1;
  uint32 private s_gasLimit;

  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  address public immutable L2_UPTIME_FEED_ADDR;

  /**
   * @notice emitted when gas cost to spend on L2 is updated
   * @param gasLimit updated gas cost
   */
  event GasLimitUpdated(uint32 gasLimit);

  /**
   * @param l1CrossDomainMessengerAddress address the L1CrossDomainMessenger contract address
   * @param l2UptimeFeedAddr the address of the OptimismSequencerUptimeFeed contract address
   * @param gasLimit the gasLimit to use for sending a message from L1 to L2
   */
  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    uint32 gasLimit
  ) {
    require(l1CrossDomainMessengerAddress != address(0), "Invalid xDomain Messenger address");
    require(l2UptimeFeedAddr != address(0), "Invalid OptimismSequencerUptimeFeed contract address");
    L1_CROSS_DOMAIN_MESSENGER_ADDRESS = l1CrossDomainMessengerAddress;
    L2_UPTIME_FEED_ADDR = l2UptimeFeedAddr;
    s_gasLimit = gasLimit;
  }

  /**
   * @notice versions:
   *
   * - OptimismValidator 0.1.0: initial release
   * - OptimismValidator 1.0.0: change target of L2 sequencer status update
   *   - now calls `updateStatus` on an L2 OptimismSequencerUptimeFeed contract instead of
   *     directly calling the Flags contract
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismValidator 1.0.0";
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
   * @dev A message is sent using the L1CrossDomainMessenger. This method is accessed controlled.
   * @param previousAnswer previous aggregator answer
   * @param currentAnswer new aggregator answer - value of 1 considers the sequencer offline.
   */
  function validate(
    uint256, /* previousRoundId */
    int256 previousAnswer,
    uint256, /* currentRoundId */
    int256 currentAnswer
  ) external override checkAccess returns (bool) {
    // Encode the OptimismSequencerUptimeFeed call
    bytes4 selector = OptimismSequencerUptimeFeedInterface.updateStatus.selector;
    bool status = currentAnswer == ANSWER_SEQ_OFFLINE;
    uint64 timestamp = uint64(block.timestamp);
    // Encode `status` and `timestamp`
    bytes memory message = abi.encodeWithSelector(selector, status, timestamp);
    // Make the xDomain call
    IL1CrossDomainMessenger(L1_CROSS_DOMAIN_MESSENGER_ADDRESS).sendMessage(
      L2_UPTIME_FEED_ADDR, // target
      message,
      s_gasLimit
    );
    // return success
    return true;
  }
}
