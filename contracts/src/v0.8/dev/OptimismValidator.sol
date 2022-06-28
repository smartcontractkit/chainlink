// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../SimpleWriteAccessController.sol";

/* ./dev dependencies - to be moved from ./dev after audit */
import "./interfaces/FlagsInterface.sol";
import "./interfaces/ForwarderInterface.sol";
import "./interfaces/OptimismSequencerUptimeFeedInterface.sol";
import "@eth-optimism/contracts/L1/messaging/IL1CrossDomainMessenger.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";

/**
 * @title OptimismValidator - makes xDomain L2 Flags contract call (using L2 xDomain Forwarder contract)
 * @notice Allows to raise and lower Flags on the Optimism L2 network through L1 bridge
 *  - The internal AccessController controls the access of the validate method
 */
contract OptimismValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  int256 private constant ANSWER_SEQ_OFFLINE = 1;
  uint32 private s_gasLimit;

  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  address public immutable L2_UPTIME_FEED_ADDR;

  /**
   * @notice emitted when a new ETH withdrawal from L2 was requested
   * @param amount of funds to withdraw
   */
  event L2WithdrawalRequested(uint256 amount);

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
   * @notice makes this contract payable
   * @dev receives funds:
   *  - to use them (if configured) to pay for L2 execution on L1
   *  - when withdrawing funds from L2 xDomain alias address (pay for L2 execution on L2)
   */
  receive() external payable {}

  /**
   * @notice withdraws all funds available in this contract to the msg.sender
   * @dev only owner can call this
   */
  function withdrawFunds() external onlyOwner {
    address payable recipient = payable(msg.sender);
    uint256 amount = address(this).balance;
    Address.sendValue(recipient, amount);
  }

  /**
   * @notice withdraws all funds available in this contract to the address specified
   * @dev only owner can call this
   * @param recipient address where to send the funds
   */
  function withdrawFundsTo(address payable recipient) external onlyOwner {
    uint256 amount = address(this).balance;
    Address.sendValue(recipient, amount);
  }

  /**
   * @notice validate method sends an xDomain L2 tx to update Flags contract, in case of change from `previousAnswer`.
   * @dev A message is created on the Optimism L1 Inbox contract. This method is accessed controlled.
   * @param previousAnswer previous aggregator answer
   * @param currentAnswer new aggregator answer - value of 1 considers the service offline.
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
