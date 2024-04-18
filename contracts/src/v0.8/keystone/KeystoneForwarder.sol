// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {Report} from "./libraries/Report.sol";

contract KeystoneForwarder is ConfirmedOwner, TypeAndVersionInterface {
  event ReportForwarded(bytes32 indexed workflowExecutionId, address indexed transmitter, bool success);

  error ReentrantCall();

  /// @notice This error is returned when the data with report is invalid.
  /// This can happen if the data is shorter than SELECTOR_LENGTH + REPORT_LENGTH.
  /// @param data the data that was received
  error InvalidData(bytes data);

  /// @notice This error is returned when the report signature length is invalid.
  /// @param signature the signature that was received
  error InvalidSignature(bytes signature);

  uint256 private constant SELECTOR_LENGTH = 4;
  uint256 private constant REPORT_LENGTH = 64;

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
    if (data.length < SELECTOR_LENGTH + REPORT_LENGTH) {
      revert InvalidData(data);
    }

    // data is an encoded call with the selector prefixed: (bytes4 selector, bytes report, ...)
    // we are able to partially decode just the first param, since we don't know the rest
    bytes memory rawReport = abi.decode(data[4:], (bytes));

    // TODO: we probably need some type of f value config?

    // Report will always be over 64 bytes here because of the check above
    bytes32 hash = keccak256(rawReport);

    // validate signatures
    for (uint256 i = 0; i < signatures.length; i++) {
      // TODO: is libocr-style multiple bytes32 arrays more optimal?
      (bytes32 r, bytes32 s, uint8 v) = _splitSignature(signatures[i]);
      // address signer = ecrecover(hash, v, r, s);
      ecrecover(hash, v, r, s);
      // TODO: we need to store oracle cluster similar to aggregator then, to validate valid signer list
    }

    (, /* bytes32 workflowId */ bytes32 workflowExecutionId) = Report.getMetadata(rawReport);

    // report was already processed
    if (s_reports[workflowExecutionId] != address(0)) {
      return false;
    }

    // solhint-disable-next-line avoid-low-level-calls
    (bool success /* bytes memory result */, ) = targetAddress.call(data);

    s_reports[workflowExecutionId] = msg.sender;

    emit ReportForwarded(workflowExecutionId, msg.sender, success);
    return true;
  }

  // get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(bytes32 workflowExecutionId) external view returns (address) {
    return s_reports[workflowExecutionId];
  }

  function _splitSignature(bytes memory sig) internal pure returns (bytes32 r, bytes32 s, uint8 v) {
    if (sig.length != 65) {
      revert InvalidSignature(sig);
    }

    assembly {
      /*
      First 32 bytes stores the length of the signature

      add(sig, 32) = pointer of sig + 32
      effectively, skips first 32 bytes of signature

      mload(p) loads next 32 bytes starting at the memory address p into memory
      */

      // first 32 bytes, after the length prefix
      r := mload(add(sig, 32))
      // second 32 bytes
      s := mload(add(sig, 64))
      // final byte (first byte of the next 32 bytes)
      v := byte(0, mload(add(sig, 96)))
    }

    return (r, s, v);
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
