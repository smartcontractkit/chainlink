// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/transmissions11/solmate/blob/97bdb2003b70382996a79a406813f76417b1cf90/src/utils/LibString.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

/// @notice Library for converting numbers into strings and other string operations.
/// @author Solady (https://github.com/vectorized/solady/blob/main/src/utils/LibString.sol)
/// @author Modified from Solmate (https://github.com/transmissions11/solmate/blob/main/src/utils/LibString.sol)
///
/// Note:
/// For performance and bytecode compactness, most of the string operations are restricted to
/// byte strings (7-bit ASCII), except where otherwise specified.
/// Usage of byte string operations on charsets with runes spanning two or more bytes
/// can lead to undefined behavior.
library LibString {
  /// @dev Returns a string from a small bytes32 string.
  /// `s` must be null-terminated, or behavior will be undefined.
  function fromSmallString(bytes32 s) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      result := mload(0x40)
      let n := 0
      for {

      } byte(n, s) {
        n := add(n, 1)
      } {

      } // Scan for '\0'.
      mstore(result, n)
      let o := add(result, 0x20)
      mstore(o, s)
      mstore(add(o, n), 0)
      mstore(0x40, add(result, 0x40))
    }
  }

  /// @dev Returns the string as a normalized null-terminated small string.
  function toSmallString(string memory s) internal pure returns (bytes32 result) {
    /// @solidity memory-safe-assembly
    assembly {
      result := mload(s)
      if iszero(lt(result, 33)) {
        mstore(0x00, 0xec92f9a3) // `TooBigForSmallString()`.
        revert(0x1c, 0x04)
      }
      result := shl(shl(3, sub(32, result)), mload(add(s, result)))
    }
  }
}
