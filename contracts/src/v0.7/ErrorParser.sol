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
   * @return string decoded revert message.
   */
  function getRevertMessage(bytes memory result) internal pure returns (string memory) {
    if (result.length < 68) return "call failed silently";
    assembly {
      result := add(result, 0x04)
    }
    return abi.decode(result, (string));
  }
}
