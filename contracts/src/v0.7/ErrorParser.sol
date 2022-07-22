// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

/**
 * @dev library for parsing hex-encoded errors external calls.
 */
library ErrorParser {

  /**
   * @notice extracts revert message from .call() result.
   *
   * @dev this extracts revert message by discarding the first 4 bytes of the
   * encoded result, the signature for Error(msg), and decode the rest as string.
   *
   * @param result encoded bytes returned from .call()
   */
  function revertWithMessage(bytes memory result) internal pure {
    if (result.length == 0) revert();
    assembly {
      revert(add(32, result),  mload(result))
    }
  }
}
