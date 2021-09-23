// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../SimpleWriteAccessController.sol";

/* ./dev dependencies - to be moved from ./dev after audit */
import "./interfaces/FlagsInterface.sol";
import "./interfaces/ForwarderInterface.sol";
import "./vendor/@eth-optimism/contracts/0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";

/**
 * @title OptimismValidator - makes xDomain L2 Flags contract call (using L2 xDomain Forwarder contract)
 * @notice Allows to raise and lower Flags on the Optimism L2 network through L1 bridge
 *  - The internal AccessController controls the access of the validate method
 */
contract OptimismValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address public constant FLAG_OPTIMISM_SEQ_OFFLINE =
    address(bytes20(bytes32(uint256(keccak256("chainlink.flags.optimism-seq-offline")) - 1)));
  // Encode underlying Flags call/s
  bytes private constant CALL_RAISE_FLAG =
    abi.encodeWithSelector(FlagsInterface.raiseFlag.selector, FLAG_OPTIMISM_SEQ_OFFLINE);
  bytes private constant CALL_LOWER_FLAG =
    abi.encodeWithSelector(FlagsInterface.lowerFlag.selector, FLAG_OPTIMISM_SEQ_OFFLINE);
  uint32 private constant CALL_GAS_LIMIT = 1_200_000;
  int256 private constant ANSWER_SEQ_OFFLINE = 1;

  address public immutable CROSS_DOMAIN_MESSENGER;
  address public immutable L2_CROSS_DOMAIN_FORWARDER;
  address public immutable L2_FLAGS;

  /**
   * @param crossDomainMessengerAddr address the xDomain bridge messenger (Optimism bridge L1) contract address
   * @param l2CrossDomainForwarderAddr the L2 Forwarder contract address
   * @param l2FlagsAddr the L2 Flags contract address
   */
  constructor(
    address crossDomainMessengerAddr,
    address l2CrossDomainForwarderAddr,
    address l2FlagsAddr
  ) {
    require(crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
    require(l2CrossDomainForwarderAddr != address(0), "Invalid L2 xDomain Forwarder address");
    require(l2FlagsAddr != address(0), "Invalid L2 Flags address");
    CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
    L2_CROSS_DOMAIN_FORWARDER = l2CrossDomainForwarderAddr;
    L2_FLAGS = l2FlagsAddr;
  }

  /**
   * @notice versions:
   *
   * - OptimismValidator 0.1.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismValidator 0.1.0";
  }

  /**
   * @notice validate method sends an xDomain L2 tx to update Flags contract, in case of change from `previousAnswer`.
   * @dev A message is sent via the Optimism CrossDomainMessenger L1 contract. The "payment" for L2 execution happens on L1,
   *   using the gas attached to this tx (some extra gas is burned by the Optimism bridge to avoid DoS attacks).
   *   This method is accessed controlled.
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
      return true; // noop
    }

    // Encode the Forwarder call
    bytes4 selector = ForwarderInterface.forward.selector;
    address target = L2_FLAGS;
    // Choose and encode the underlying Flags call
    bytes memory data = currentAnswer == ANSWER_SEQ_OFFLINE ? CALL_RAISE_FLAG : CALL_LOWER_FLAG;
    bytes memory message = abi.encodeWithSelector(selector, target, data);
    // Make the xDomain call
    iOVM_CrossDomainMessenger(CROSS_DOMAIN_MESSENGER).sendMessage(L2_CROSS_DOMAIN_FORWARDER, message, CALL_GAS_LIMIT);
    // return success
    return true;
  }
}
