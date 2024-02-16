// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IForwarder} from "./interfaces/IForwarder.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";

// solhint-disable custom-errors, no-unused-vars
contract KeystoneForwarder is IForwarder, ConfirmedOwner, TypeAndVersionInterface {
  error ReentrantCall();

  struct HotVars {
    bool reentrancyGuard; // guard against reentrancy
  }

  HotVars internal s_hotVars; // Mixture of config and state, commonly accessed

  mapping(bytes32 => address) internal s_reports;

  constructor() ConfirmedOwner(msg.sender) {}

  // solhint-disable avoid-low-level-calls, chainlink-solidity/explicit-returns
  function splitSignature(bytes memory sig) public pure returns (bytes32 r, bytes32 s, uint8 v) {
    require(sig.length == 65, "invalid signature length");

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

    // implicitly return (r, s, v)
  }

  // solhint-disable avoid-low-level-calls, chainlink-solidity/explicit-returns
  function splitReport(bytes memory rawReport) public pure returns (bytes32 workflowId, bytes32 workflowExecutionId) {
    require(rawReport.length > 64, "invalid report length");
    assembly {
      workflowId := mload(add(rawReport, 4))
      workflowExecutionId := mload(add(rawReport, 36)) // 4 + 32
    }
  }

  // send a report to targetAddress
  function report(
    address targetAddress,
    bytes calldata data,
    bytes[] calldata signatures
  ) external nonReentrant returns (bool) {
    // data is an encoded call with the selector prefixed: (bytes4 selector, bytes report, ...)
    // we are able to partially decode just the first param, since we don't know the rest
    bytes memory rawReport = abi.decode(data[4:], (bytes));

    // TODO: we probably need some type of f value config?

    bytes32 hash = keccak256(rawReport);

    // validate signatures
    for (uint256 i = 0; i < signatures.length; i++) {
      // TODO: is libocr-style multiple bytes32 arrays more optimal?
      (bytes32 r, bytes32 s, uint8 v) = splitSignature(signatures[i]);
      address signer = ecrecover(hash, v, r, s);
      // TODO: we need to store oracle cluster similar to aggregator then, to validate valid signer list
    }

    (bytes32 workflowId, bytes32 workflowExecutionId) = splitReport(rawReport);

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
