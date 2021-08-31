// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../SimpleWriteAccessController.sol";

/* ./dev dependencies - to be re/moved after audit */
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";
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
  // Config for L1 -> L2 `createRetryableTicket` call
  struct GasConfiguration {
    uint256 maxSubmissionCost;
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

  AccessControllerInterface private s_gasConfigAccessController;
  GasConfiguration private s_gasConfig;

  /**
   * @notice emitted when a new gas configuration is set
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param gasPriceBid maximum gas price to pay on L2
   * @param gasCostL2Value value to send to L2 to cover cost (submission + gas)
   * @param refundAddr address where gas excess on L2 will be sent
   */
  event GasConfigurationSet(
    uint256 maxSubmissionCost,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address indexed refundAddr
  );

  /**
   * @notice emitted when a new ETH withdrawal from L2 was requested
   * @param id unique id of the published retryable transaction (keccak256(requestID, uint(0))
   * @param amount of funds to withdraw
   */
  event WithdrawalFromL2Requested(
    uint256 indexed id,
    uint256 amount
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
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param gasPriceBid maximum gas price to pay on L2
   * @param gasCostL2Value value to send to L2 to cover cost (submission + gas)
   * @param maxGas gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundAddr address where gas excess on L2 will be sent
   */
  constructor(
    address crossDomainMessengerAddr,
    address l2CrossDomainForwarderAddr,
    address l2FlagsAddr,
    address gasConfigAccessControllerAddr,
    uint256 maxSubmissionCost,
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
    // additional configuration
    s_gasConfigAccessController = AccessControllerInterface(gasConfigAccessControllerAddr);
    _setGasConfiguration(maxSubmissionCost, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
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
   * @param maxGas gas limit for immediate L2 execution attempt.
   * @param gasPriceBid maximum gas price to pay on L2
   * @param refundAddr address where gas excess on L2 will be sent
   * @return id unique id of the published retryable transaction (keccak256(requestID, uint(0))
   */
  function withdrawFundsFromL2(
    uint256 amount,
    uint256 maxSubmissionCost,
    uint256 maxGas,
    uint256 gasPriceBid,
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

    id = IInbox(CROSS_DOMAIN_MESSENGER).createRetryableTicketNoRefundAliasRewrite(
        ARBSYS_ADDR,
        amount,
        maxSubmissionCost,
        refundAddr,
        refundAddr,
        maxGas,
        gasPriceBid,
        l1ToL2Calldata
    );
    emit WithdrawalFromL2Requested(id, amount);
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
   * @param maxSubmissionCost maximum submission cost willing to pay on L2
   * @param gasPriceBid maximum gas price to pay on L2
   * @param gasCostL2Value value to send to L2 to cover cost (submission + gas)
   * @param maxGas gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundAddr address where gas excess on L2 will be sent
   */
  function setGasConfiguration(
    uint256 maxSubmissionCost,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address refundAddr
  )
    external
  {
    require(s_gasConfigAccessController.hasAccess(msg.sender, msg.data), "Access required to set config");
    _setGasConfiguration(maxSubmissionCost, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
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
    // NOTICE: if gasCostL2Value is zero the payment is processed on L2 so the L2 address needs to be funded, as it will
    // paying the fee. We also ignore the returned msg number, that can be queried via the InboxMessageDelivered event.
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

  function _setGasConfiguration(
    uint256 maxSubmissionCost,
    uint256 gasPriceBid,
    uint256 gasCostL2Value,
    uint256 maxGas,
    address refundAddr
  )
    internal
  {
    // L2 will pay the fee if gasCostL2Value is zero
    if (gasCostL2Value > 0) {
      uint256 minGasCostL2Value = maxSubmissionCost + maxGas * gasPriceBid;
      require(gasCostL2Value >= minGasCostL2Value, "Gas cost provided is too low");
    }
    s_gasConfig = GasConfiguration(maxSubmissionCost, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
    emit GasConfigurationSet(maxSubmissionCost, gasPriceBid, gasCostL2Value, maxGas, refundAddr);
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
