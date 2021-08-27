// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../SimpleWriteAccessController.sol";

/* ./dev dependencies - to be re/moved after audit */
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";
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
  // Config for L1 -> L2 `createRetryableTicket` call
  struct GasConfiguration {
    uint256 maxSubmissionCost;
    uint256 maxGasPrice;
    uint256 gasCostL2;
    uint256 gasLimitL2;
    address refundableAddr;
  }

  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address constant private FLAG_ARBITRUM_SEQ_OFFLINE = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-seq-offline")) - 1)));
  // Encode underlying Flags call/s
  bytes constant private CALL_RAISE_FLAG = abi.encodeWithSelector(FlagsInterface.raiseFlag.selector, FLAG_ARBITRUM_SEQ_OFFLINE);
  bytes constant private CALL_LOWER_FLAG = abi.encodeWithSelector(FlagsInterface.lowerFlag.selector, FLAG_ARBITRUM_SEQ_OFFLINE);
  int256 constant private ANSWER_SEQ_OFFLINE = 1;

  address immutable public CROSS_DOMAIN_MESSENGER;
  address immutable public L2_CROSS_DOMAIN_FORWARDER;
  address immutable public L2_FLAGS;

  AccessControllerInterface private s_gasConfigAccessController;
  GasConfiguration private s_gasConfig;

  /**
   * @notice emitted when a new gas configuration is set
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maxGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddr address where gas excess on L2 will be sent
   */
  event GasConfigurationSet(
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address indexed refundableAddr
  );

  /**
   * @notice emitted when a new gas access-control contract is set
   * @param previous the address prior to the current setting
   * @param current the address of the new access-control contract
   */
  event GasAccessControllerSet(
    address indexed previous,
    address indexed current
  );

  /**
   * @param crossDomainMessengerAddr address the xDomain bridge messenger (Arbitrum Inbox L1) contract address
   * @param l2CrossDomainForwarderAddr the L2 Forwarder contract address
   * @param l2FlagsAddr the L2 Flags contract address
   * @param gasConfigAccessControllerAddr address of the access controller for managing gas price on Arbitrum
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maxGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param gasLimitL2 gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundableAddr address where gas excess on L2 will be sent
   */
  constructor(
    address crossDomainMessengerAddr,
    address l2CrossDomainForwarderAddr,
    address l2FlagsAddr,
    address gasConfigAccessControllerAddr,
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address refundableAddr
  ) {
    require(crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
    require(l2CrossDomainForwarderAddr != address(0), "Invalid L2 xDomain Forwarder address");
    require(l2FlagsAddr != address(0), "Invalid Flags contract address");
    CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
    L2_CROSS_DOMAIN_FORWARDER = l2CrossDomainForwarderAddr;
    L2_FLAGS = l2FlagsAddr;
    // additional configuration
    s_gasConfigAccessController = AccessControllerInterface(gasConfigAccessControllerAddr);
    _setGasConfiguration(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddr);
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

  /// @return gas config AccessControllerInterface contract address
  function gasConfigAccessController()
    external
    view
    virtual
    returns (address)
  {
    return address(s_gasConfigAccessController);
  }

  /// @return stored GasConfiguration
  function gasConfig()
    external
    view
    virtual
    returns (GasConfiguration memory)
  {
    return s_gasConfig;
  }

  /// @notice makes this contract payable as it need funds to pay for L2 transactions fees on L1.
  receive() external payable {}

  /**
   * @notice withdraws all funds availbale in this contract to the msg.sender
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
   * @notice withdraws all funds availbale in this contract to the address specified
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
   * @notice sets gas config AccessControllerInterface contract
   * @dev only owner can call this
   * @param accessController new AccessControllerInterface contract address
   */
  function setGasAccessController(
    address accessController
  )
    external
    onlyOwner()
  {
    _setGasAccessController(accessController);
  }

  /**
   * @notice sets Arbitrum gas configuration
   * @dev access control provided by s_gasConfigAccessController
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maxGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param gasLimitL2 gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundableAddr address where gas excess on L2 will be sent
   */
  function setGasConfiguration(
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address refundableAddr
  )
    external
  {
    require(s_gasConfigAccessController.hasAccess(msg.sender, msg.data), "Access required to set config");
    _setGasConfiguration(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddr);
  }

  /**
   * @notice validate method updates the state of an L2 Flag in case of change on the Arbitrum Sequencer.
   * A one answer considers the service as offline.
   * In case the previous answer is the same as the current it does not trigger any tx on L2. In other case,
   * a retryable ticket is created on the Arbitrum L1 Inbox contract. The tx gas fee can be paid from this
   * contract providing a value, or the same address on L2.
   * @dev access control provided internally by SimpleWriteAccessController
   * @param previousAnswer previous aggregator answer
   * @param currentAnswer new aggregator answer
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
    // NOTICE: if gasCostL2 is zero the payment is processed on L2 so the L2 address needs to be funded, as it will
    // paying the fee. We also ignore the returned msg number, that can be queried via the InboxMessageDelivered event.
    IInbox(CROSS_DOMAIN_MESSENGER).createRetryableTicket{value: s_gasConfig.gasCostL2}(
      L2_CROSS_DOMAIN_FORWARDER,
      0, // L2 call value
      // NOTICE: maxSubmissionCost info will possibly become available on L1 after the London fork. At that time this
      // contract could start querying/calculating it directly so we wouldn't need to configure it statically. On L2 this
      // info is available via `ArbRetryableTx.getSubmissionPrice`.
      s_gasConfig.maxSubmissionCost, // Max submission cost of sending data length
      s_gasConfig.refundableAddr, // excessFeeRefundAddress
      s_gasConfig.refundableAddr, // callValueRefundAddress
      s_gasConfig.gasLimitL2,
      s_gasConfig.maxGasPrice,
      message
    );
    // return success
    return true;
  }

  function _setGasConfiguration(
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address refundableAddr
  )
    internal
  {
    // L2 will pay the fee if gasCostL2 is zero
    if (gasCostL2 > 0) {
      uint256 minGasCostValue = maxSubmissionCost + gasLimitL2 * maxGasPrice;
      require(gasCostL2 >= minGasCostValue, "Gas cost provided is too low");
    }
    s_gasConfig = GasConfiguration(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddr);
    emit GasConfigurationSet(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddr);
  }

  function _setGasAccessController(
    address accessController
  )
    internal
  {
    address previousAccessController = address(s_gasConfigAccessController);
    if (accessController != previousAccessController) {
      s_gasConfigAccessController = AccessControllerInterface(accessController);
      emit GasAccessControllerSet(previousAccessController, accessController);
    }
  }
}
