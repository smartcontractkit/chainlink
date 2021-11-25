// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../SimpleWriteAccessController.sol";

/* ./dev dependencies - to be moved from ./dev after audit */
import "./interfaces/ArbitrumSequencerUptimeFeedInterface.sol";
import "./interfaces/FlagsInterface.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import "./vendor/arb-os/e8d9696f21/contracts/arbos/builtin/ArbSys.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";

/**
 * @title ArbitrumValidator - makes xDomain L2 Flags contract call (using L2 xDomain Forwarder contract)
 * @notice Allows to raise and lower Flags on the Arbitrum L2 network through L1 bridge
 *  - The internal AccessController controls the access of the validate method
 *  - Gas configuration is controlled by a configurable external SimpleWriteAccessController
 *  - Funds on the contract are managed by the owner
 */
contract ArbitrumValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  enum PaymentStrategy {
    L1,
    L2
  }
  // Config for L1 -> L2 Arbitrum retryable ticket message
  struct GasConfig {
    uint256 maxGas;
    uint256 gasPriceBid;
    address gasPriceL1FeedAddr;
  }

  /// @dev Precompiled contract that exists in every Arbitrum chain at address(100). Exposes a variety of system-level functionality.
  address constant ARBSYS_ADDR = address(0x0000000000000000000000000000000000000064);

  int256 private constant ANSWER_SEQ_OFFLINE = 1;

  address public immutable CROSS_DOMAIN_MESSENGER;
  address public immutable L2_SEQ_STATUS_RECORDER;
  // L2 xDomain alias address of this contract
  address public immutable L2_ALIAS = AddressAliasHelper.applyL1ToL2Alias(address(this));

  PaymentStrategy private s_paymentStrategy;
  GasConfig private s_gasConfig;
  AccessControllerInterface private s_configAC;

  /**
   * @notice emitted when a new payment strategy is set
   * @param paymentStrategy strategy describing how the contract pays for xDomain calls
   */
  event PaymentStrategySet(PaymentStrategy indexed paymentStrategy);

  /**
   * @notice emitted when a new gas configuration is set
   * @param maxGas gas limit for immediate L2 execution attempt.
   * @param gasPriceBid maximum L2 gas price to pay
   * @param gasPriceL1FeedAddr address of the L1 gas price feed (used to approximate Arbitrum retryable ticket submission cost)
   */
  event GasConfigSet(uint256 maxGas, uint256 gasPriceBid, address indexed gasPriceL1FeedAddr);

  /**
   * @notice emitted when a new gas access-control contract is set
   * @param previous the address prior to the current setting
   * @param current the address of the new access-control contract
   */
  event ConfigACSet(address indexed previous, address indexed current);

  /**
   * @notice emitted when a new ETH withdrawal from L2 was requested
   * @param id unique id of the published retryable transaction (keccak256(requestID, uint(0))
   * @param amount of funds to withdraw
   */
  event L2WithdrawalRequested(uint256 indexed id, uint256 amount, address indexed refundAddr);

  /**
   * @param crossDomainMessengerAddr address the xDomain bridge messenger (Arbitrum Inbox L1) contract address
   * @param l2ArbitrumSequencerUptimeFeedAddr the L2 Flags contract address
   * @param configACAddr address of the access controller for managing gas price on Arbitrum
   * @param maxGas gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param gasPriceBid maximum L2 gas price to pay
   * @param gasPriceL1FeedAddr address of the L1 gas price feed (used to approximate Arbitrum retryable ticket submission cost)
   * @param paymentStrategy strategy describing how the contract pays for xDomain calls
   */
  constructor(
    address crossDomainMessengerAddr,
    address l2ArbitrumSequencerUptimeFeedAddr,
    address configACAddr,
    uint256 maxGas,
    uint256 gasPriceBid,
    address gasPriceL1FeedAddr,
    PaymentStrategy paymentStrategy
  ) {
    require(crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
    require(l2ArbitrumSequencerUptimeFeedAddr != address(0), "Invalid ArbitrumSequencerUptimeFeed contract address");
    CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
    L2_SEQ_STATUS_RECORDER = l2ArbitrumSequencerUptimeFeedAddr;
    // Additional L2 payment configuration
    _setConfigAC(configACAddr);
    _setGasConfig(maxGas, gasPriceBid, gasPriceL1FeedAddr);
    _setPaymentStrategy(paymentStrategy);
  }

  /**
   * @notice versions:
   *
   * - ArbitrumValidator 0.1.0: initial release
   * - ArbitrumValidator 0.2.0: critical Arbitrum network update
   *   - xDomain `msg.sender` backwards incompatible change (now an alias address)
   *   - new `withdrawFundsFromL2` fn that withdraws from L2 xDomain alias address
   *   - approximation of `maxSubmissionCost` using a L1 gas price feed
   * - ArbitrumValidator 1.0.0: change target of L2 sequencer status update
   *   - now calls `updateStatus` on an L2 ArbitrumSequencerUptimeFeed contract instead of
   *     directly calling the Flags contract
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ArbitrumValidator 1.0.0";
  }

  /// @return stored PaymentStrategy
  function paymentStrategy() external view virtual returns (PaymentStrategy) {
    return s_paymentStrategy;
  }

  /// @return stored GasConfig
  function gasConfig() external view virtual returns (GasConfig memory) {
    return s_gasConfig;
  }

  /// @return config AccessControllerInterface contract address
  function configAC() external view virtual returns (address) {
    return address(s_configAC);
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
   * @notice withdraws funds from L2 xDomain alias address (representing this L1 contract)
   * @dev only owner can call this
   * @param amount of funds to withdraws
   * @param refundAddr address where gas excess on L2 will be sent
   *   WARNING: `refundAddr` is not aliased! Make sure you can recover the refunded funds on L2.
   * @return id unique id of the published retryable transaction (keccak256(requestID, uint(0))
   */
  function withdrawFundsFromL2(uint256 amount, address refundAddr) external onlyOwner returns (uint256 id) {
    // Build an xDomain message to trigger the ArbSys precompile, which will create a L2 -> L1 tx transferring `amount`
    bytes memory message = abi.encodeWithSelector(ArbSys.withdrawEth.selector, address(this));
    // Make the xDomain call
    // NOTICE: We approximate the max submission cost of sending a retryable tx with specific calldata length.
    uint256 maxSubmissionCost = _approximateMaxSubmissionCost(message.length);
    uint256 maxGas = 120_000; // static `maxGas` for L2 -> L1 transfer
    uint256 gasPriceBid = s_gasConfig.gasPriceBid;
    uint256 l1PaymentValue = s_paymentStrategy == PaymentStrategy.L1
      ? _maxRetryableTicketCost(maxSubmissionCost, maxGas, gasPriceBid)
      : 0;
    // NOTICE: In the case of PaymentStrategy.L2 the L2 xDomain alias address needs to be funded, as it will be paying the fee.
    id = IInbox(CROSS_DOMAIN_MESSENGER).createRetryableTicketNoRefundAliasRewrite{value: l1PaymentValue}(
      ARBSYS_ADDR, // target
      amount, // L2 call value (requested)
      maxSubmissionCost,
      refundAddr, // excessFeeRefundAddress
      refundAddr, // callValueRefundAddress
      maxGas,
      gasPriceBid,
      message
    );
    emit L2WithdrawalRequested(id, amount, refundAddr);
  }

  /**
   * @notice sets config AccessControllerInterface contract
   * @dev only owner can call this
   * @param accessController new AccessControllerInterface contract address
   */
  function setConfigAC(address accessController) external onlyOwner {
    _setConfigAC(accessController);
  }

  /**
   * @notice sets Arbitrum gas configuration
   * @dev access control provided by `configAC`
   * @param maxGas gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param gasPriceBid maximum L2 gas price to pay
   * @param gasPriceL1FeedAddr address of the L1 gas price feed (used to approximate Arbitrum retryable ticket submission cost)
   */
  function setGasConfig(
    uint256 maxGas,
    uint256 gasPriceBid,
    address gasPriceL1FeedAddr
  ) external onlyOwnerOrConfigAccess {
    _setGasConfig(maxGas, gasPriceBid, gasPriceL1FeedAddr);
  }

  /**
   * @notice sets the payment strategy
   * @dev access control provided by `configAC`
   * @param paymentStrategy strategy describing how the contract pays for xDomain calls
   */
  function setPaymentStrategy(PaymentStrategy paymentStrategy) external onlyOwnerOrConfigAccess {
    _setPaymentStrategy(paymentStrategy);
  }

  /**
   * @notice validate method sends an xDomain L2 tx to update Flags contract, in case of change from `previousAnswer`.
   * @dev A retryable ticket is created on the Arbitrum L1 Inbox contract. The tx gas fee can be paid from this
   *   contract providing a value, or if no L1 value is sent with the xDomain message the gas will be paid by
   *   the L2 xDomain alias account (generated from `address(this)`). This method is accessed controlled.
   * @param previousAnswer previous aggregator answer
   * @param currentAnswer new aggregator answer - value of 1 considers the service offline.
   */
  function validate(
    uint256, /* previousRoundId */
    int256 previousAnswer,
    uint256, /* currentRoundId */
    int256 currentAnswer
  ) external override checkAccess returns (bool) {
    // Avoids resending to L2 the same tx on every call
    if (previousAnswer == currentAnswer) {
      return true;
    }

    // Excess gas on L2 will be sent to the L2 xDomain alias address of this contract
    address refundAddr = L2_ALIAS;
    // Encode the ArbitrumSequencerUptimeFeed call
    bytes4 selector = ArbitrumSequencerUptimeFeedInterface.updateStatus.selector;
    bool status = currentAnswer == ANSWER_SEQ_OFFLINE;
    uint64 timestamp = uint64(block.timestamp);
    // Encode `status` and `timestamp`
    bytes memory message = abi.encodeWithSelector(selector, status, timestamp);
    // Make the xDomain call
    // NOTICE: We approximate the max submission cost of sending a retryable tx with specific calldata length.
    uint256 maxSubmissionCost = _approximateMaxSubmissionCost(message.length);
    uint256 maxGas = s_gasConfig.maxGas;
    uint256 gasPriceBid = s_gasConfig.gasPriceBid;
    uint256 l1PaymentValue = s_paymentStrategy == PaymentStrategy.L1
      ? _maxRetryableTicketCost(maxSubmissionCost, maxGas, gasPriceBid)
      : 0;
    // NOTICE: In the case of PaymentStrategy.L2 the L2 xDomain alias address needs to be funded, as it will be paying the fee.
    // We also ignore the returned msg number, that can be queried via the `InboxMessageDelivered` event.
    IInbox(CROSS_DOMAIN_MESSENGER).createRetryableTicketNoRefundAliasRewrite{value: l1PaymentValue}(
      L2_SEQ_STATUS_RECORDER, // target
      0, // L2 call value
      maxSubmissionCost,
      refundAddr, // excessFeeRefundAddress
      refundAddr, // callValueRefundAddress
      maxGas,
      gasPriceBid,
      message
    );
    // return success
    return true;
  }

  /// @notice internal method that stores the payment strategy
  function _setPaymentStrategy(PaymentStrategy paymentStrategy) internal {
    s_paymentStrategy = paymentStrategy;
    emit PaymentStrategySet(paymentStrategy);
  }

  /// @notice internal method that stores the gas configuration
  function _setGasConfig(
    uint256 maxGas,
    uint256 gasPriceBid,
    address gasPriceL1FeedAddr
  ) internal {
    require(maxGas > 0, "Max gas is zero");
    require(gasPriceBid > 0, "Gas price bid is zero");
    require(gasPriceL1FeedAddr != address(0), "Gas price Aggregator is zero address");
    s_gasConfig = GasConfig(maxGas, gasPriceBid, gasPriceL1FeedAddr);
    emit GasConfigSet(maxGas, gasPriceBid, gasPriceL1FeedAddr);
  }

  /// @notice Internal method that stores the configuration access controller
  function _setConfigAC(address accessController) internal {
    address previousAccessController = address(s_configAC);
    if (accessController != previousAccessController) {
      s_configAC = AccessControllerInterface(accessController);
      emit ConfigACSet(previousAccessController, accessController);
    }
  }

  /**
   * @notice Internal method that approximates the `maxSubmissionCost` (using the L1 gas price feed)
   * @dev On L2 this info is available via `ArbRetryableTx.getSubmissionPrice`.
   * @param calldataSizeInBytes xDomain message size in bytes
   */
  function _approximateMaxSubmissionCost(uint256 calldataSizeInBytes) internal view returns (uint256) {
    (, int256 l1GasPriceInWei, , , ) = AggregatorV3Interface(s_gasConfig.gasPriceL1FeedAddr).latestRoundData();
    uint256 l1GasPriceEstimate = uint256(l1GasPriceInWei) * 3; // add 200% buffer (price volatility error margin)
    return (l1GasPriceEstimate * calldataSizeInBytes) / 256 + l1GasPriceEstimate;
  }

  /// @notice Internal helper method that calculates the total cost of the xDomain retryable ticket call
  function _maxRetryableTicketCost(
    uint256 maxSubmissionCost,
    uint256 maxGas,
    uint256 gasPriceBid
  ) internal pure returns (uint256) {
    return maxSubmissionCost + maxGas * gasPriceBid;
  }

  /// @dev reverts if the caller does not have access to change the configuration
  modifier onlyOwnerOrConfigAccess() {
    require(
      msg.sender == owner() || (address(s_configAC) != address(0) && s_configAC.hasAccess(msg.sender, msg.data)),
      "No access"
    );
    _;
  }
}
