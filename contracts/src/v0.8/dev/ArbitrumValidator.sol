// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../SimpleWriteAccessController.sol";

/* ./dev dependencies - to be re/moved after audit */
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import "./vendor/arb-os/e8d9696f21/contracts/arbos/builtin/ArbSys.sol";
import "./interfaces/FlagsInterface.sol";
import "./interfaces/ForwarderInterface.sol";

/**
 * @title ArbitrumValidator - makes xDomain L2 Flags contract call (using L2 xDomain Forwarder contract)
 * @notice Allows to raise and lower Flags on the Arbitrum L2 network through L1 bridge
 *  - The internal AccessController controls the access of the validate method
 *  - Gas configuration is controlled by a configurable external SimpleWriteAccessController
 *  - Funds on the contract are managed by the owner
 */
contract ArbitrumValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  // Config for L1 -> L2 Arbitrum retryable ticket message
  struct GasConfig {
    uint256 maxSubmissionCost;
    address gasFeedAddr;
    uint256 gasPriceBid;
    uint256 gasCostL2Value;
    uint256 maxGas;
    address refundAddr;
  }

  /// @dev Precompiled contract that exists in every Arbitrum chain at address(100). Exposes a variety of system-level functionality.
  address constant ARBSYS_ADDR = address(0x0000000000000000000000000000000000000064);

  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address constant private FLAG_ARBITRUM_SEQ_OFFLINE = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-seq-offline")) - 1)));
  // Encode underlying Flags call/s
  bytes constant private CALL_RAISE_FLAG = abi.encodeWithSelector(FlagsInterface.raiseFlag.selector, FLAG_ARBITRUM_SEQ_OFFLINE);
  bytes constant private CALL_LOWER_FLAG = abi.encodeWithSelector(FlagsInterface.lowerFlag.selector, FLAG_ARBITRUM_SEQ_OFFLINE);
  int256 constant private ANSWER_SEQ_OFFLINE = 1;

  address immutable public CROSS_DOMAIN_MESSENGER;
  address immutable public L2_CROSS_DOMAIN_FORWARDER;
  address immutable public L2_FLAGS;

  GasConfig private s_gasConfig;
  AccessControllerInterface private s_gasConfigAC;

  /**
   * @notice emitted when a new gas configuration is set
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param gasFeedAddr address of the L1 gas feed (used to approximate Arbitrum `maxSubmissionCost`, if `maxSubmissionCost == 0`)
   * @param gasPriceBid maximum gas price to pay on L2
   * @param gasCostL2Value value to send to L2 to cover cost (submission + gas)
   * @param refundAddr address where gas excess on L2 will be sent
   */
  event GasConfigSet(
    uint256 maxSubmissionCost,
    address indexed gasFeedAddr,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address indexed refundAddr
  );

  /**
   * @notice emitted when a new gas access-control contract is set
   * @param previous the address prior to the current setting
   * @param current the address of the new access-control contract
   */
  event GasConfigACSet(
    address indexed previous,
    address indexed current
  );

  /**
   * @notice emitted when a new ETH withdrawal from L2 was requested
   * @param id unique id of the published retryable transaction (keccak256(requestID, uint(0))
   * @param amount of funds to withdraw
   */
  event L2WithdrawalRequested(
    uint256 indexed id,
    uint256 amount
  );

  /**
   * @param crossDomainMessengerAddr address the xDomain bridge messenger (Arbitrum Inbox L1) contract address
   * @param l2CrossDomainForwarderAddr the L2 Forwarder contract address
   * @param l2FlagsAddr the L2 Flags contract address
   * @param gasConfigACAddr address of the access controller for managing gas price on Arbitrum
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param gasFeedAddr address of the L1 gas feed (used to approximate Arbitrum `maxSubmissionCost`, if `maxSubmissionCost == 0`)
   * @param gasPriceBid maximum gas price to pay on L2
   * @param gasCostL2Value value to send to L2 to cover cost (submission + gas)
   * @param maxGas gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundAddr address where gas excess on L2 will be sent
   */
  constructor(
    address crossDomainMessengerAddr,
    address l2CrossDomainForwarderAddr,
    address l2FlagsAddr,
    address gasConfigACAddr,
    uint256 maxSubmissionCost,
    address gasFeedAddr,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address refundAddr
  ) {
    require(crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
    require(l2CrossDomainForwarderAddr != address(0), "Invalid L2 xDomain Forwarder address");
    require(l2FlagsAddr != address(0), "Invalid Flags contract address");
    CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
    L2_CROSS_DOMAIN_FORWARDER = l2CrossDomainForwarderAddr;
    L2_FLAGS = l2FlagsAddr;
    // Additional L2 payment configuration
    s_gasConfigAC = AccessControllerInterface(gasConfigACAddr);
    _setGasConfig(maxSubmissionCost, gasFeedAddr, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
  }

  /**
   * @notice versions:
   *
   * - ArbitrumValidator 0.1.0: initial release
   * - ArbitrumValidator 0.2.0: critical Arbitrum network update, xDomain `msg.sender` backwards incompatible change
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion()
    external
    pure
    virtual
    override
    returns (
      string memory
    )
  {
    return "ArbitrumValidator 0.2.0";
  }

  /// @return stored GasConfig
  function gasConfig()
    external
    view
    virtual
    returns (GasConfig memory)
  {
    return s_gasConfig;
  }

  /// @return gas config AccessControllerInterface contract address
  function gasConfigAC()
    external
    view
    virtual
    returns (address)
  {
    return address(s_gasConfigAC);
  }

  /**
   * @notice makes this contract payable
   * @dev receives funds:
   *  - to use them (if configured) to pay for L2 execution on L1
   *  - when withdrawing funds from L2 xDomain alias address (pay for L2 execution on L2)
   */
  receive()
    external
    payable
  {}

  /**
   * @notice withdraws all funds available in this contract to the msg.sender
   * @dev only owner can call this
   */
  function withdrawFunds()
    external
    onlyOwner()
  {
    address payable to = payable(msg.sender);
    to.transfer(address(this).balance);
  }

  /**
   * @notice withdraws all funds available in this contract to the address specified
   * @dev only owner can call this
   * @param to address where to send the funds
   */
  function withdrawFundsTo(
    address payable to
  ) 
    external
    onlyOwner()
  {
    to.transfer(address(this).balance);
  }

  /**
   * @notice withdraws funds from L2 xDomain alias address (representing this L1 contract)
   * @dev only owner can call this
   * @param amount of funds to withdraws
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param refundAddr address where gas excess on L2 will be sent
   * @return id unique id of the published retryable transaction (keccak256(requestID, uint(0))
   */
  function withdrawFundsFromL2(
    uint256 amount,
    uint256 maxSubmissionCost,
    address refundAddr
  )
    external
    onlyOwner()
    returns (uint256 id)
  {
    // We want the L1 to L2 tx to trigger the Arbsys precompile
    // then create a L2 to L1 transaction transferring `amount`
    bytes memory l1ToL2Calldata = abi.encodeWithSelector(
        ArbSys.sendTxToL1.selector,
        address(this)
    );

    // If `refundAddr` not set, default to the L2 xDomain alias address
    refundAddr = _selectRefundAddr(refundAddr);
    // If `maxSubmissionCost` not set, approximate
    maxSubmissionCost = _selectMaxSubmissionCost(maxSubmissionCost, l1ToL2Calldata.length);

    id = IInbox(CROSS_DOMAIN_MESSENGER).createRetryableTicketNoRefundAliasRewrite(
        ARBSYS_ADDR,
        amount,
        maxSubmissionCost,
        refundAddr,
        refundAddr,
        100_000, // static `maxGas` for L2 -> L1 transfer
        s_gasConfig.gasPriceBid,
        l1ToL2Calldata
    );
    emit L2WithdrawalRequested(id, amount);
  }

  /**
   * @notice sets gas config AccessControllerInterface contract
   * @dev only owner can call this
   * @param accessController new AccessControllerInterface contract address
   */
  function setGasConfigAC(
    address accessController
  )
    external
    onlyOwner()
  {
    _setGasConfigAC(accessController);
  }

  /**
   * @notice sets Arbitrum gas configuration
   * @dev access control provided by s_gasConfigAC
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param gasFeedAddr address of the L1 gas feed (used to approximate Arbitrum `maxSubmissionCost`, if `maxSubmissionCost == 0`)
   * @param gasPriceBid maximum gas price to pay on L2
   * @param gasCostL2Value value to send to L2 to cover cost (submission + gas)
   * @param maxGas gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundAddr address where gas excess on L2 will be sent
   */
  function setGasConfig(
    uint256 maxSubmissionCost,
    address gasFeedAddr,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address refundAddr
  )
    external
  {
    require(s_gasConfigAC.hasAccess(msg.sender, msg.data), "No access");
    _setGasConfig(maxSubmissionCost, gasFeedAddr, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
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
    uint256 /* previousRoundId */,
    int256 previousAnswer,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) 
    external
    override
    checkAccess()
    returns (bool)
  {
    // Avoids resending to L2 the same tx on every call
    if (previousAnswer == currentAnswer) {
      return true;
    }

    // Encode the Forwarder call
    bytes4 selector = ForwarderInterface.forward.selector;
    address target = L2_FLAGS;
    // Choose and encode the underlying Flags call
    bytes memory data = currentAnswer == ANSWER_SEQ_OFFLINE ? CALL_RAISE_FLAG : CALL_LOWER_FLAG;
    bytes memory message = abi.encodeWithSelector(selector, target, data);
    // Make the xDomain call
    // NOTICE: if `gasCostL2Value` is zero the payment is processed on L2. In that case the L2 xDomain alias address needs to be funded,
    //   as it will be paying the fee. We also ignore the returned msg number, that can be queried via the `InboxMessageDelivered` event.
    IInbox(CROSS_DOMAIN_MESSENGER).createRetryableTicketNoRefundAliasRewrite{value: s_gasConfig.gasCostL2Value}(
      L2_CROSS_DOMAIN_FORWARDER,
      0, // L2 call value
      // NOTICE: maxSubmissionCost info will possibly become available on L1 after the London fork. At that time this
      // contract could start querying/calculating it directly so we wouldn't need to configure it statically. On L2 this
      // info is available via `ArbRetryableTx.getSubmissionPrice`.
      s_gasConfig.maxSubmissionCost, // Max submission cost of sending data length
      s_gasConfig.refundAddr, // excessFeeRefundAddress
      s_gasConfig.refundAddr, // callValueRefundAddress
      s_gasConfig.maxGas,
      s_gasConfig.gasPriceBid,
      message
    );
    // return success
    return true;
  }

  /// @notice internal method that stores the gas configuration
  function _setGasConfig(
    uint256 maxSubmissionCost,
    address gasFeedAddr,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address refundAddr
  )
    internal
  {
    // L2 xDomain alias addr pays the fee if `gasCostL2Value` is zero
    // Else, we check if configured L1 payment is sufficient (and not excessive)
    if (gasCostL2Value > 0) {
      uint256 minGasCostL2Value = _selectMaxSubmissionCost(maxSubmissionCost, CALL_RAISE_FLAG.length) + maxGas * gasPriceBid;
      uint256 maxGasCostL2Value = 2 * minGasCostL2Value;
      require(gasCostL2Value >= minGasCostL2Value, "gasCostL2Value < MIN");
      require(gasCostL2Value < maxGasCostL2Value, "gasCostL2Value >= MAX");
    }
    // If `refundAddr` not set, default to the L2 xDomain alias address
    refundAddr = _selectRefundAddr(refundAddr);

    s_gasConfig = GasConfig(maxSubmissionCost, gasFeedAddr, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
    emit GasConfigSet(maxSubmissionCost, gasFeedAddr, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
  }

  /// @notice internal method that stores the gas configuration access controller
  function _setGasConfigAC(
    address accessController
  )
    internal
  {
    address previousAccessController = address(s_gasConfigAC);
    if (accessController != previousAccessController) {
      s_gasConfigAC = AccessControllerInterface(accessController);
      emit GasConfigACSet(previousAccessController, accessController);
    }
  }

  /// @notice internal method that selects the `refundAddr` (either configured or L2 xDomain alias of `address(this)`)
  function _selectRefundAddr(
    address refundAddr
  )
    internal
    view
    returns (address)
  {
    return refundAddr == address(0x0) ? AddressAliasHelper.applyL1ToL2Alias(address(this)) : refundAddr;
  }

  /// @notice internal method that selects the `maxSubmissionCost` (either configured or approximated)
  function _selectMaxSubmissionCost(
    uint256 maxSubmissionCost,
    uint256 calldataSizeInBytes
  )
    internal
    view
    returns (uint256)
  {
    return maxSubmissionCost != 0 ? maxSubmissionCost : _approximateMaxSubmissionCost(calldataSizeInBytes);
  }

  /// @notice internal method that approximates the `maxSubmissionCost` (using the gas price feed)
  function _approximateMaxSubmissionCost(
    uint256 calldataSizeInBytes
  )
    internal
    view
    returns (uint256)
  {
    (,int256 l1GasPriceInWei,,,) = AggregatorV3Interface(s_gasConfig.gasFeedAddr).latestRoundData();
    uint256 l1GasPriceEstimate = uint256(l1GasPriceInWei) * 2; // add 100% buffer
    return (l1GasPriceEstimate / 256) * calldataSizeInBytes + l1GasPriceEstimate;
  }
}
