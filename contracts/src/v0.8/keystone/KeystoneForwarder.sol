// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IForwarder} from "./interfaces/IForwarder.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {Utils} from "./libraries/Utils.sol";

// solhint-disable gas-custom-errors, no-unused-vars
contract KeystoneForwarder is IForwarder, ConfirmedOwner, TypeAndVersionInterface {
  error ReentrantCall();

  struct HotVars {
    bool reentrancyGuard; // guard against reentrancy
  }

  HotVars internal s_hotVars; // Mixture of config and state, commonly accessed

  mapping(bytes32 => address) internal s_reports;

  constructor() ConfirmedOwner(msg.sender) {}

  // send a report to targetAddress
  function report(
    address targetAddress,
    bytes calldata data,
    bytes[] calldata signatures
  ) external nonReentrant returns (bool) {
    require(data.length > 4 + 64, "invalid data length");

    // data is an encoded call with the selector prefixed: (bytes4 selector, bytes report, ...)
    // we are able to partially decode just the first param, since we don't know the rest
    bytes memory rawReport = abi.decode(data[4:], (bytes));

    // TODO: we probably need some type of f value config?

    bytes32 hash = keccak256(rawReport);

    // validate signatures
    for (uint256 i = 0; i < signatures.length; i++) {
      // TODO: is libocr-style multiple bytes32 arrays more optimal?
      (bytes32 r, bytes32 s, uint8 v) = Utils._splitSignature(signatures[i]);
      address signer = ecrecover(hash, v, r, s);
      // TODO: we need to store oracle cluster similar to aggregator then, to validate valid signer list
    }

    (bytes32 workflowId, bytes32 workflowExecutionId) = Utils._splitReport(rawReport);

    // report was already processed
    if (s_reports[workflowExecutionId] != address(0)) {
      return false;
    }

    // solhint-disable-next-line avoid-low-level-calls
    (bool success, bytes memory result) = targetAddress.call(data);

    s_reports[workflowExecutionId] = msg.sender;
    return true;
  }

  // get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(bytes32 workflowExecutionId) external view returns (address) {
    return s_reports[workflowExecutionId];
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "KeystoneForwarder 1.0.0";
  }

  /**
   * @dev replicates Open Zeppelin's ReentrancyGuard but optimized to fit our storage
   */
  modifier nonReentrant() {
    if (s_hotVars.reentrancyGuard) revert ReentrantCall();
    s_hotVars.reentrancyGuard = true;
    _;
    s_hotVars.reentrancyGuard = false;
  }
}
