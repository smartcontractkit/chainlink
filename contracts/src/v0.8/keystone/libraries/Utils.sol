// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// solhint-disable gas-custom-errors
library Utils {
  // solhint-disable avoid-low-level-calls, chainlink-solidity/explicit-returns
  function _splitSignature(bytes memory sig) internal pure returns (bytes32 r, bytes32 s, uint8 v) {
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
  function _splitReport(
    bytes memory rawReport
  ) internal pure returns (bytes32 workflowId, bytes32 workflowExecutionId) {
    require(rawReport.length > 64, "invalid report length");
    assembly {
      // skip first 32 bytes, contains length of the report
      workflowId := mload(add(rawReport, 32))
      workflowExecutionId := mload(add(rawReport, 64))
    }
  }
}
