// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

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
  /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
  /*                        CUSTOM ERRORS                       */
  /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

  /// @dev The length of the output is too small to contain all the hex digits.
  error HexLengthInsufficient();

  /// @dev The length of the string is more than 32 bytes.
  error TooBigForSmallString();

  /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
  /*                         CONSTANTS                          */
  /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

  /// @dev The constant returned when the `search` is not found in the string.
  uint256 internal constant NOT_FOUND = type(uint256).max;

  /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
  /*                     DECIMAL OPERATIONS                     */
  /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

  /// @dev Returns the base 10 decimal representation of `value`.
  function toString(uint256 value) internal pure returns (string memory str) {
    /// @solidity memory-safe-assembly
    assembly {
      // The maximum value of a uint256 contains 78 digits (1 byte per digit), but
      // we allocate 0xa0 bytes to keep the free memory pointer 32-byte word aligned.
      // We will need 1 word for the trailing zeros padding, 1 word for the length,
      // and 3 words for a maximum of 78 digits.
      str := add(mload(0x40), 0x80)
      // Update the free memory pointer to allocate.
      mstore(0x40, add(str, 0x20))
      // Zeroize the slot after the string.
      mstore(str, 0)

      // Cache the end of the memory to calculate the length later.
      let end := str

      let w := not(0) // Tsk.
      // We write the string from rightmost digit to leftmost digit.
      // The following is essentially a do-while loop that also handles the zero case.
      for {
        let temp := value
      } 1 {

      } {
        str := add(str, w) // `sub(str, 1)`.
        // Write the character to the pointer.
        // The ASCII index of the '0' character is 48.
        mstore8(str, add(48, mod(temp, 10)))
        // Keep dividing `temp` until zero.
        temp := div(temp, 10)
        if iszero(temp) {
          break
        }
      }

      let length := sub(end, str)
      // Move the pointer 32 bytes leftwards to make room for the length.
      str := sub(str, 0x20)
      // Store the length.
      mstore(str, length)
    }
  }

  /// @dev Returns the base 10 decimal representation of `value`.
  function toString(int256 value) internal pure returns (string memory str) {
    if (value >= 0) {
      return toString(uint256(value));
    }
    unchecked {
      str = toString(uint256(-value));
    }
    /// @solidity memory-safe-assembly
    assembly {
      // We still have some spare memory space on the left,
      // as we have allocated 3 words (96 bytes) for up to 78 digits.
      let length := mload(str) // Load the string length.
      mstore(str, 0x2d) // Store the '-' character.
      str := sub(str, 1) // Move back the string pointer by a byte.
      mstore(str, add(length, 1)) // Update the string length.
    }
  }

  /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
  /*                   HEXADECIMAL OPERATIONS                   */
  /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

  /// @dev Returns the hexadecimal representation of `value`,
  /// left-padded to an input length of `length` bytes.
  /// The output is prefixed with "0x" encoded using 2 hexadecimal digits per byte,
  /// giving a total length of `length * 2 + 2` bytes.
  /// Reverts if `length` is too small for the output to contain all the digits.
  function toHexString(uint256 value, uint256 length) internal pure returns (string memory str) {
    str = toHexStringNoPrefix(value, length);
    /// @solidity memory-safe-assembly
    assembly {
      let strLength := add(mload(str), 2) // Compute the length.
      mstore(str, 0x3078) // Write the "0x" prefix.
      str := sub(str, 2) // Move the pointer.
      mstore(str, strLength) // Write the length.
    }
  }

  /// @dev Returns the hexadecimal representation of `value`,
  /// left-padded to an input length of `length` bytes.
  /// The output is prefixed with "0x" encoded using 2 hexadecimal digits per byte,
  /// giving a total length of `length * 2` bytes.
  /// Reverts if `length` is too small for the output to contain all the digits.
  function toHexStringNoPrefix(uint256 value, uint256 length) internal pure returns (string memory str) {
    /// @solidity memory-safe-assembly
    assembly {
      // We need 0x20 bytes for the trailing zeros padding, `length * 2` bytes
      // for the digits, 0x02 bytes for the prefix, and 0x20 bytes for the length.
      // We add 0x20 to the total and round down to a multiple of 0x20.
      // (0x20 + 0x20 + 0x02 + 0x20) = 0x62.
      str := add(mload(0x40), and(add(shl(1, length), 0x42), not(0x1f)))
      // Allocate the memory.
      mstore(0x40, add(str, 0x20))
      // Zeroize the slot after the string.
      mstore(str, 0)

      // Cache the end to calculate the length later.
      let end := str
      // Store "0123456789abcdef" in scratch space.
      mstore(0x0f, 0x30313233343536373839616263646566)

      let start := sub(str, add(length, length))
      let w := not(1) // Tsk.
      let temp := value
      // We write the string from rightmost digit to leftmost digit.
      // The following is essentially a do-while loop that also handles the zero case.
      for {

      } 1 {

      } {
        str := add(str, w) // `sub(str, 2)`.
        mstore8(add(str, 1), mload(and(temp, 15)))
        mstore8(str, mload(and(shr(4, temp), 15)))
        temp := shr(8, temp)
        if iszero(xor(str, start)) {
          break
        }
      }

      if temp {
        mstore(0x00, 0x2194895a) // `HexLengthInsufficient()`.
        revert(0x1c, 0x04)
      }

      // Compute the string's length.
      let strLength := sub(end, str)
      // Move the pointer and write the length.
      str := sub(str, 0x20)
      mstore(str, strLength)
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output is prefixed with "0x" and encoded using 2 hexadecimal digits per byte.
  /// As address are 20 bytes long, the output will left-padded to have
  /// a length of `20 * 2 + 2` bytes.
  function toHexString(uint256 value) internal pure returns (string memory str) {
    str = toHexStringNoPrefix(value);
    /// @solidity memory-safe-assembly
    assembly {
      let strLength := add(mload(str), 2) // Compute the length.
      mstore(str, 0x3078) // Write the "0x" prefix.
      str := sub(str, 2) // Move the pointer.
      mstore(str, strLength) // Write the length.
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output is prefixed with "0x".
  /// The output excludes leading "0" from the `toHexString` output.
  /// `0x00: "0x0", 0x01: "0x1", 0x12: "0x12", 0x123: "0x123"`.
  function toMinimalHexString(uint256 value) internal pure returns (string memory str) {
    str = toHexStringNoPrefix(value);
    /// @solidity memory-safe-assembly
    assembly {
      let o := eq(byte(0, mload(add(str, 0x20))), 0x30) // Whether leading zero is present.
      let strLength := add(mload(str), 2) // Compute the length.
      mstore(add(str, o), 0x3078) // Write the "0x" prefix, accounting for leading zero.
      str := sub(add(str, o), 2) // Move the pointer, accounting for leading zero.
      mstore(str, sub(strLength, o)) // Write the length, accounting for leading zero.
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output excludes leading "0" from the `toHexStringNoPrefix` output.
  /// `0x00: "0", 0x01: "1", 0x12: "12", 0x123: "123"`.
  function toMinimalHexStringNoPrefix(uint256 value) internal pure returns (string memory str) {
    str = toHexStringNoPrefix(value);
    /// @solidity memory-safe-assembly
    assembly {
      let o := eq(byte(0, mload(add(str, 0x20))), 0x30) // Whether leading zero is present.
      let strLength := mload(str) // Get the length.
      str := add(str, o) // Move the pointer, accounting for leading zero.
      mstore(str, sub(strLength, o)) // Write the length, accounting for leading zero.
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output is encoded using 2 hexadecimal digits per byte.
  /// As address are 20 bytes long, the output will left-padded to have
  /// a length of `20 * 2` bytes.
  function toHexStringNoPrefix(uint256 value) internal pure returns (string memory str) {
    /// @solidity memory-safe-assembly
    assembly {
      // We need 0x20 bytes for the trailing zeros padding, 0x20 bytes for the length,
      // 0x02 bytes for the prefix, and 0x40 bytes for the digits.
      // The next multiple of 0x20 above (0x20 + 0x20 + 0x02 + 0x40) is 0xa0.
      str := add(mload(0x40), 0x80)
      // Allocate the memory.
      mstore(0x40, add(str, 0x20))
      // Zeroize the slot after the string.
      mstore(str, 0)

      // Cache the end to calculate the length later.
      let end := str
      // Store "0123456789abcdef" in scratch space.
      mstore(0x0f, 0x30313233343536373839616263646566)

      let w := not(1) // Tsk.
      // We write the string from rightmost digit to leftmost digit.
      // The following is essentially a do-while loop that also handles the zero case.
      for {
        let temp := value
      } 1 {

      } {
        str := add(str, w) // `sub(str, 2)`.
        mstore8(add(str, 1), mload(and(temp, 15)))
        mstore8(str, mload(and(shr(4, temp), 15)))
        temp := shr(8, temp)
        if iszero(temp) {
          break
        }
      }

      // Compute the string's length.
      let strLength := sub(end, str)
      // Move the pointer and write the length.
      str := sub(str, 0x20)
      mstore(str, strLength)
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output is prefixed with "0x", encoded using 2 hexadecimal digits per byte,
  /// and the alphabets are capitalized conditionally according to
  /// https://eips.ethereum.org/EIPS/eip-55
  function toHexStringChecksummed(address value) internal pure returns (string memory str) {
    str = toHexString(value);
    /// @solidity memory-safe-assembly
    assembly {
      let mask := shl(6, div(not(0), 255)) // `0b010000000100000000 ...`
      let o := add(str, 0x22)
      let hashed := and(keccak256(o, 40), mul(34, mask)) // `0b10001000 ... `
      let t := shl(240, 136) // `0b10001000 << 240`
      for {
        let i := 0
      } 1 {

      } {
        mstore(add(i, i), mul(t, byte(i, hashed)))
        i := add(i, 1)
        if eq(i, 20) {
          break
        }
      }
      mstore(o, xor(mload(o), shr(1, and(mload(0x00), and(mload(o), mask)))))
      o := add(o, 0x20)
      mstore(o, xor(mload(o), shr(1, and(mload(0x20), and(mload(o), mask)))))
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output is prefixed with "0x" and encoded using 2 hexadecimal digits per byte.
  function toHexString(address value) internal pure returns (string memory str) {
    str = toHexStringNoPrefix(value);
    /// @solidity memory-safe-assembly
    assembly {
      let strLength := add(mload(str), 2) // Compute the length.
      mstore(str, 0x3078) // Write the "0x" prefix.
      str := sub(str, 2) // Move the pointer.
      mstore(str, strLength) // Write the length.
    }
  }

  /// @dev Returns the hexadecimal representation of `value`.
  /// The output is encoded using 2 hexadecimal digits per byte.
  function toHexStringNoPrefix(address value) internal pure returns (string memory str) {
    /// @solidity memory-safe-assembly
    assembly {
      str := mload(0x40)

      // Allocate the memory.
      // We need 0x20 bytes for the trailing zeros padding, 0x20 bytes for the length,
      // 0x02 bytes for the prefix, and 0x28 bytes for the digits.
      // The next multiple of 0x20 above (0x20 + 0x20 + 0x02 + 0x28) is 0x80.
      mstore(0x40, add(str, 0x80))

      // Store "0123456789abcdef" in scratch space.
      mstore(0x0f, 0x30313233343536373839616263646566)

      str := add(str, 2)
      mstore(str, 40)

      let o := add(str, 0x20)
      mstore(add(o, 40), 0)

      value := shl(96, value)

      // We write the string from rightmost digit to leftmost digit.
      // The following is essentially a do-while loop that also handles the zero case.
      for {
        let i := 0
      } 1 {

      } {
        let p := add(o, add(i, i))
        let temp := byte(i, value)
        mstore8(add(p, 1), mload(and(temp, 15)))
        mstore8(p, mload(shr(4, temp)))
        i := add(i, 1)
        if eq(i, 20) {
          break
        }
      }
    }
  }

  /// @dev Returns the hex encoded string from the raw bytes.
  /// The output is encoded using 2 hexadecimal digits per byte.
  function toHexString(bytes memory raw) internal pure returns (string memory str) {
    str = toHexStringNoPrefix(raw);
    /// @solidity memory-safe-assembly
    assembly {
      let strLength := add(mload(str), 2) // Compute the length.
      mstore(str, 0x3078) // Write the "0x" prefix.
      str := sub(str, 2) // Move the pointer.
      mstore(str, strLength) // Write the length.
    }
  }

  /// @dev Returns the hex encoded string from the raw bytes.
  /// The output is encoded using 2 hexadecimal digits per byte.
  function toHexStringNoPrefix(bytes memory raw) internal pure returns (string memory str) {
    /// @solidity memory-safe-assembly
    assembly {
      let length := mload(raw)
      str := add(mload(0x40), 2) // Skip 2 bytes for the optional prefix.
      mstore(str, add(length, length)) // Store the length of the output.

      // Store "0123456789abcdef" in scratch space.
      mstore(0x0f, 0x30313233343536373839616263646566)

      let o := add(str, 0x20)
      let end := add(raw, length)

      for {

      } iszero(eq(raw, end)) {

      } {
        raw := add(raw, 1)
        mstore8(add(o, 1), mload(and(mload(raw), 15)))
        mstore8(o, mload(and(shr(4, mload(raw)), 15)))
        o := add(o, 2)
      }
      mstore(o, 0) // Zeroize the slot after the string.
      mstore(0x40, add(o, 0x20)) // Allocate the memory.
    }
  }

  /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
  /*                   RUNE STRING OPERATIONS                   */
  /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

  /// @dev Returns the number of UTF characters in the string.
  function runeCount(string memory s) internal pure returns (uint256 result) {
    /// @solidity memory-safe-assembly
    assembly {
      if mload(s) {
        mstore(0x00, div(not(0), 255))
        mstore(0x20, 0x0202020202020202020202020202020202020202020202020303030304040506)
        let o := add(s, 0x20)
        let end := add(o, mload(s))
        for {
          result := 1
        } 1 {
          result := add(result, 1)
        } {
          o := add(o, byte(0, mload(shr(250, mload(o)))))
          if iszero(lt(o, end)) {
            break
          }
        }
      }
    }
  }

  /// @dev Returns if this string is a 7-bit ASCII string.
  /// (i.e. all characters codes are in [0..127])
  function is7BitASCII(string memory s) internal pure returns (bool result) {
    /// @solidity memory-safe-assembly
    assembly {
      let mask := shl(7, div(not(0), 255))
      result := 1
      let n := mload(s)
      if n {
        let o := add(s, 0x20)
        let end := add(o, n)
        let last := mload(end)
        mstore(end, 0)
        for {

        } 1 {

        } {
          if and(mask, mload(o)) {
            result := 0
            break
          }
          o := add(o, 0x20)
          if iszero(lt(o, end)) {
            break
          }
        }
        mstore(end, last)
      }
    }
  }

  /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
  /*                   BYTE STRING OPERATIONS                   */
  /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

  // For performance and bytecode compactness, byte string operations are restricted
  // to 7-bit ASCII strings. All offsets are byte offsets, not UTF character offsets.
  // Usage of byte string operations on charsets with runes spanning two or more bytes
  // can lead to undefined behavior.

  /// @dev Returns `subject` all occurrences of `search` replaced with `replacement`.
  function replace(
    string memory subject,
    string memory search,
    string memory replacement
  ) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let subjectLength := mload(subject)
      let searchLength := mload(search)
      let replacementLength := mload(replacement)

      subject := add(subject, 0x20)
      search := add(search, 0x20)
      replacement := add(replacement, 0x20)
      result := add(mload(0x40), 0x20)

      let subjectEnd := add(subject, subjectLength)
      if iszero(gt(searchLength, subjectLength)) {
        let subjectSearchEnd := add(sub(subjectEnd, searchLength), 1)
        let h := 0
        if iszero(lt(searchLength, 0x20)) {
          h := keccak256(search, searchLength)
        }
        let m := shl(3, sub(0x20, and(searchLength, 0x1f)))
        let s := mload(search)
        for {

        } 1 {

        } {
          let t := mload(subject)
          // Whether the first `searchLength % 32` bytes of
          // `subject` and `search` matches.
          if iszero(shr(m, xor(t, s))) {
            if h {
              if iszero(eq(keccak256(subject, searchLength), h)) {
                mstore(result, t)
                result := add(result, 1)
                subject := add(subject, 1)
                if iszero(lt(subject, subjectSearchEnd)) {
                  break
                }
                continue
              }
            }
            // Copy the `replacement` one word at a time.
            for {
              let o := 0
            } 1 {

            } {
              mstore(add(result, o), mload(add(replacement, o)))
              o := add(o, 0x20)
              if iszero(lt(o, replacementLength)) {
                break
              }
            }
            result := add(result, replacementLength)
            subject := add(subject, searchLength)
            if searchLength {
              if iszero(lt(subject, subjectSearchEnd)) {
                break
              }
              continue
            }
          }
          mstore(result, t)
          result := add(result, 1)
          subject := add(subject, 1)
          if iszero(lt(subject, subjectSearchEnd)) {
            break
          }
        }
      }

      let resultRemainder := result
      result := add(mload(0x40), 0x20)
      let k := add(sub(resultRemainder, result), sub(subjectEnd, subject))
      // Copy the rest of the string one word at a time.
      for {

      } lt(subject, subjectEnd) {

      } {
        mstore(resultRemainder, mload(subject))
        resultRemainder := add(resultRemainder, 0x20)
        subject := add(subject, 0x20)
      }
      result := sub(result, 0x20)
      let last := add(add(result, 0x20), k) // Zeroize the slot after the string.
      mstore(last, 0)
      mstore(0x40, add(last, 0x20)) // Allocate the memory.
      mstore(result, k) // Store the length.
    }
  }

  /// @dev Returns the byte index of the first location of `search` in `subject`,
  /// searching from left to right, starting from `from`.
  /// Returns `NOT_FOUND` (i.e. `type(uint256).max`) if the `search` is not found.
  function indexOf(string memory subject, string memory search, uint256 from) internal pure returns (uint256 result) {
    /// @solidity memory-safe-assembly
    assembly {
      for {
        let subjectLength := mload(subject)
      } 1 {

      } {
        if iszero(mload(search)) {
          if iszero(gt(from, subjectLength)) {
            result := from
            break
          }
          result := subjectLength
          break
        }
        let searchLength := mload(search)
        let subjectStart := add(subject, 0x20)

        result := not(0) // Initialize to `NOT_FOUND`.

        subject := add(subjectStart, from)
        let end := add(sub(add(subjectStart, subjectLength), searchLength), 1)

        let m := shl(3, sub(0x20, and(searchLength, 0x1f)))
        let s := mload(add(search, 0x20))

        if iszero(and(lt(subject, end), lt(from, subjectLength))) {
          break
        }

        if iszero(lt(searchLength, 0x20)) {
          for {
            let h := keccak256(add(search, 0x20), searchLength)
          } 1 {

          } {
            if iszero(shr(m, xor(mload(subject), s))) {
              if eq(keccak256(subject, searchLength), h) {
                result := sub(subject, subjectStart)
                break
              }
            }
            subject := add(subject, 1)
            if iszero(lt(subject, end)) {
              break
            }
          }
          break
        }
        for {

        } 1 {

        } {
          if iszero(shr(m, xor(mload(subject), s))) {
            result := sub(subject, subjectStart)
            break
          }
          subject := add(subject, 1)
          if iszero(lt(subject, end)) {
            break
          }
        }
        break
      }
    }
  }

  /// @dev Returns the byte index of the first location of `search` in `subject`,
  /// searching from left to right.
  /// Returns `NOT_FOUND` (i.e. `type(uint256).max`) if the `search` is not found.
  function indexOf(string memory subject, string memory search) internal pure returns (uint256 result) {
    result = indexOf(subject, search, 0);
  }

  /// @dev Returns the byte index of the first location of `search` in `subject`,
  /// searching from right to left, starting from `from`.
  /// Returns `NOT_FOUND` (i.e. `type(uint256).max`) if the `search` is not found.
  function lastIndexOf(
    string memory subject,
    string memory search,
    uint256 from
  ) internal pure returns (uint256 result) {
    /// @solidity memory-safe-assembly
    assembly {
      for {

      } 1 {

      } {
        result := not(0) // Initialize to `NOT_FOUND`.
        let searchLength := mload(search)
        if gt(searchLength, mload(subject)) {
          break
        }
        let w := result

        let fromMax := sub(mload(subject), searchLength)
        if iszero(gt(fromMax, from)) {
          from := fromMax
        }

        let end := add(add(subject, 0x20), w)
        subject := add(add(subject, 0x20), from)
        if iszero(gt(subject, end)) {
          break
        }
        // As this function is not too often used,
        // we shall simply use keccak256 for smaller bytecode size.
        for {
          let h := keccak256(add(search, 0x20), searchLength)
        } 1 {

        } {
          if eq(keccak256(subject, searchLength), h) {
            result := sub(subject, add(end, 1))
            break
          }
          subject := add(subject, w) // `sub(subject, 1)`.
          if iszero(gt(subject, end)) {
            break
          }
        }
        break
      }
    }
  }

  /// @dev Returns the byte index of the first location of `search` in `subject`,
  /// searching from right to left.
  /// Returns `NOT_FOUND` (i.e. `type(uint256).max`) if the `search` is not found.
  function lastIndexOf(string memory subject, string memory search) internal pure returns (uint256 result) {
    result = lastIndexOf(subject, search, uint256(int256(-1)));
  }

  /// @dev Returns true if `search` is found in `subject`, false otherwise.
  function contains(string memory subject, string memory search) internal pure returns (bool) {
    return indexOf(subject, search) != NOT_FOUND;
  }

  /// @dev Returns whether `subject` starts with `search`.
  function startsWith(string memory subject, string memory search) internal pure returns (bool result) {
    /// @solidity memory-safe-assembly
    assembly {
      let searchLength := mload(search)
      // Just using keccak256 directly is actually cheaper.
      // forgefmt: disable-next-item
      result := and(
        iszero(gt(searchLength, mload(subject))),
        eq(keccak256(add(subject, 0x20), searchLength), keccak256(add(search, 0x20), searchLength))
      )
    }
  }

  /// @dev Returns whether `subject` ends with `search`.
  function endsWith(string memory subject, string memory search) internal pure returns (bool result) {
    /// @solidity memory-safe-assembly
    assembly {
      let searchLength := mload(search)
      let subjectLength := mload(subject)
      // Whether `search` is not longer than `subject`.
      let withinRange := iszero(gt(searchLength, subjectLength))
      // Just using keccak256 directly is actually cheaper.
      // forgefmt: disable-next-item
      result := and(
        withinRange,
        eq(
          keccak256(
            // `subject + 0x20 + max(subjectLength - searchLength, 0)`.
            add(add(subject, 0x20), mul(withinRange, sub(subjectLength, searchLength))),
            searchLength
          ),
          keccak256(add(search, 0x20), searchLength)
        )
      )
    }
  }

  /// @dev Returns `subject` repeated `times`.
  function repeat(string memory subject, uint256 times) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let subjectLength := mload(subject)
      if iszero(or(iszero(times), iszero(subjectLength))) {
        subject := add(subject, 0x20)
        result := mload(0x40)
        let output := add(result, 0x20)
        for {

        } 1 {

        } {
          // Copy the `subject` one word at a time.
          for {
            let o := 0
          } 1 {

          } {
            mstore(add(output, o), mload(add(subject, o)))
            o := add(o, 0x20)
            if iszero(lt(o, subjectLength)) {
              break
            }
          }
          output := add(output, subjectLength)
          times := sub(times, 1)
          if iszero(times) {
            break
          }
        }
        mstore(output, 0) // Zeroize the slot after the string.
        let resultLength := sub(output, add(result, 0x20))
        mstore(result, resultLength) // Store the length.
        // Allocate the memory.
        mstore(0x40, add(result, add(resultLength, 0x20)))
      }
    }
  }

  /// @dev Returns a copy of `subject` sliced from `start` to `end` (exclusive).
  /// `start` and `end` are byte offsets.
  function slice(string memory subject, uint256 start, uint256 end) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let subjectLength := mload(subject)
      if iszero(gt(subjectLength, end)) {
        end := subjectLength
      }
      if iszero(gt(subjectLength, start)) {
        start := subjectLength
      }
      if lt(start, end) {
        result := mload(0x40)
        let resultLength := sub(end, start)
        mstore(result, resultLength)
        subject := add(subject, start)
        let w := not(0x1f)
        // Copy the `subject` one word at a time, backwards.
        for {
          let o := and(add(resultLength, 0x1f), w)
        } 1 {

        } {
          mstore(add(result, o), mload(add(subject, o)))
          o := add(o, w) // `sub(o, 0x20)`.
          if iszero(o) {
            break
          }
        }
        // Zeroize the slot after the string.
        mstore(add(add(result, 0x20), resultLength), 0)
        // Allocate memory for the length and the bytes,
        // rounded up to a multiple of 32.
        mstore(0x40, add(result, and(add(resultLength, 0x3f), w)))
      }
    }
  }

  /// @dev Returns a copy of `subject` sliced from `start` to the end of the string.
  /// `start` is a byte offset.
  function slice(string memory subject, uint256 start) internal pure returns (string memory result) {
    result = slice(subject, start, uint256(int256(-1)));
  }

  /// @dev Returns all the indices of `search` in `subject`.
  /// The indices are byte offsets.
  function indicesOf(string memory subject, string memory search) internal pure returns (uint256[] memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let subjectLength := mload(subject)
      let searchLength := mload(search)

      if iszero(gt(searchLength, subjectLength)) {
        subject := add(subject, 0x20)
        search := add(search, 0x20)
        result := add(mload(0x40), 0x20)

        let subjectStart := subject
        let subjectSearchEnd := add(sub(add(subject, subjectLength), searchLength), 1)
        let h := 0
        if iszero(lt(searchLength, 0x20)) {
          h := keccak256(search, searchLength)
        }
        let m := shl(3, sub(0x20, and(searchLength, 0x1f)))
        let s := mload(search)
        for {

        } 1 {

        } {
          let t := mload(subject)
          // Whether the first `searchLength % 32` bytes of
          // `subject` and `search` matches.
          if iszero(shr(m, xor(t, s))) {
            if h {
              if iszero(eq(keccak256(subject, searchLength), h)) {
                subject := add(subject, 1)
                if iszero(lt(subject, subjectSearchEnd)) {
                  break
                }
                continue
              }
            }
            // Append to `result`.
            mstore(result, sub(subject, subjectStart))
            result := add(result, 0x20)
            // Advance `subject` by `searchLength`.
            subject := add(subject, searchLength)
            if searchLength {
              if iszero(lt(subject, subjectSearchEnd)) {
                break
              }
              continue
            }
          }
          subject := add(subject, 1)
          if iszero(lt(subject, subjectSearchEnd)) {
            break
          }
        }
        let resultEnd := result
        // Assign `result` to the free memory pointer.
        result := mload(0x40)
        // Store the length of `result`.
        mstore(result, shr(5, sub(resultEnd, add(result, 0x20))))
        // Allocate memory for result.
        // We allocate one more word, so this array can be recycled for {split}.
        mstore(0x40, add(resultEnd, 0x20))
      }
    }
  }

  /// @dev Returns a arrays of strings based on the `delimiter` inside of the `subject` string.
  function split(string memory subject, string memory delimiter) internal pure returns (string[] memory result) {
    uint256[] memory indices = indicesOf(subject, delimiter);
    /// @solidity memory-safe-assembly
    assembly {
      let w := not(0x1f)
      let indexPtr := add(indices, 0x20)
      let indicesEnd := add(indexPtr, shl(5, add(mload(indices), 1)))
      mstore(add(indicesEnd, w), mload(subject))
      mstore(indices, add(mload(indices), 1))
      let prevIndex := 0
      for {

      } 1 {

      } {
        let index := mload(indexPtr)
        mstore(indexPtr, 0x60)
        if iszero(eq(index, prevIndex)) {
          let element := mload(0x40)
          let elementLength := sub(index, prevIndex)
          mstore(element, elementLength)
          // Copy the `subject` one word at a time, backwards.
          for {
            let o := and(add(elementLength, 0x1f), w)
          } 1 {

          } {
            mstore(add(element, o), mload(add(add(subject, prevIndex), o)))
            o := add(o, w) // `sub(o, 0x20)`.
            if iszero(o) {
              break
            }
          }
          // Zeroize the slot after the string.
          mstore(add(add(element, 0x20), elementLength), 0)
          // Allocate memory for the length and the bytes,
          // rounded up to a multiple of 32.
          mstore(0x40, add(element, and(add(elementLength, 0x3f), w)))
          // Store the `element` into the array.
          mstore(indexPtr, element)
        }
        prevIndex := add(index, mload(delimiter))
        indexPtr := add(indexPtr, 0x20)
        if iszero(lt(indexPtr, indicesEnd)) {
          break
        }
      }
      result := indices
      if iszero(mload(delimiter)) {
        result := add(indices, 0x20)
        mstore(result, sub(mload(indices), 2))
      }
    }
  }

  /// @dev Returns a concatenated string of `a` and `b`.
  /// Cheaper than `string.concat()` and does not de-align the free memory pointer.
  function concat(string memory a, string memory b) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let w := not(0x1f)
      result := mload(0x40)
      let aLength := mload(a)
      // Copy `a` one word at a time, backwards.
      for {
        let o := and(add(aLength, 0x20), w)
      } 1 {

      } {
        mstore(add(result, o), mload(add(a, o)))
        o := add(o, w) // `sub(o, 0x20)`.
        if iszero(o) {
          break
        }
      }
      let bLength := mload(b)
      let output := add(result, aLength)
      // Copy `b` one word at a time, backwards.
      for {
        let o := and(add(bLength, 0x20), w)
      } 1 {

      } {
        mstore(add(output, o), mload(add(b, o)))
        o := add(o, w) // `sub(o, 0x20)`.
        if iszero(o) {
          break
        }
      }
      let totalLength := add(aLength, bLength)
      let last := add(add(result, 0x20), totalLength)
      // Zeroize the slot after the string.
      mstore(last, 0)
      // Stores the length.
      mstore(result, totalLength)
      // Allocate memory for the length and the bytes,
      // rounded up to a multiple of 32.
      mstore(0x40, and(add(last, 0x1f), w))
    }
  }

  /// @dev Returns a copy of the string in either lowercase or UPPERCASE.
  /// WARNING! This function is only compatible with 7-bit ASCII strings.
  function toCase(string memory subject, bool toUpper) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let length := mload(subject)
      if length {
        result := add(mload(0x40), 0x20)
        subject := add(subject, 1)
        let flags := shl(add(70, shl(5, toUpper)), 0x3ffffff)
        let w := not(0)
        for {
          let o := length
        } 1 {

        } {
          o := add(o, w)
          let b := and(0xff, mload(add(subject, o)))
          mstore8(add(result, o), xor(b, and(shr(b, flags), 0x20)))
          if iszero(o) {
            break
          }
        }
        result := mload(0x40)
        mstore(result, length) // Store the length.
        let last := add(add(result, 0x20), length)
        mstore(last, 0) // Zeroize the slot after the string.
        mstore(0x40, add(last, 0x20)) // Allocate the memory.
      }
    }
  }

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

  /// @dev Returns the small string, with all bytes after the first null byte zeroized.
  function normalizeSmallString(bytes32 s) internal pure returns (bytes32 result) {
    /// @solidity memory-safe-assembly
    assembly {
      for {

      } byte(result, s) {
        result := add(result, 1)
      } {

      } // Scan for '\0'.
      mstore(0x00, s)
      mstore(result, 0x00)
      result := mload(0x00)
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

  /// @dev Returns a lowercased copy of the string.
  /// WARNING! This function is only compatible with 7-bit ASCII strings.
  function lower(string memory subject) internal pure returns (string memory result) {
    result = toCase(subject, false);
  }

  /// @dev Returns an UPPERCASED copy of the string.
  /// WARNING! This function is only compatible with 7-bit ASCII strings.
  function upper(string memory subject) internal pure returns (string memory result) {
    result = toCase(subject, true);
  }

  /// @dev Escapes the string to be used within HTML tags.
  function escapeHTML(string memory s) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let end := add(s, mload(s))
      result := add(mload(0x40), 0x20)
      // Store the bytes of the packed offsets and strides into the scratch space.
      // `packed = (stride << 5) | offset`. Max offset is 20. Max stride is 6.
      mstore(0x1f, 0x900094)
      mstore(0x08, 0xc0000000a6ab)
      // Store "&quot;&amp;&#39;&lt;&gt;" into the scratch space.
      mstore(0x00, shl(64, 0x2671756f743b26616d703b262333393b266c743b2667743b))
      for {

      } iszero(eq(s, end)) {

      } {
        s := add(s, 1)
        let c := and(mload(s), 0xff)
        // Not in `["\"","'","&","<",">"]`.
        if iszero(and(shl(c, 1), 0x500000c400000000)) {
          mstore8(result, c)
          result := add(result, 1)
          continue
        }
        let t := shr(248, mload(c))
        mstore(result, mload(and(t, 0x1f)))
        result := add(result, shr(5, t))
      }
      let last := result
      mstore(last, 0) // Zeroize the slot after the string.
      result := mload(0x40)
      mstore(result, sub(last, add(result, 0x20))) // Store the length.
      mstore(0x40, add(last, 0x20)) // Allocate the memory.
    }
  }

  /// @dev Escapes the string to be used within double-quotes in a JSON.
  /// If `addDoubleQuotes` is true, the result will be enclosed in double-quotes.
  function escapeJSON(string memory s, bool addDoubleQuotes) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      let end := add(s, mload(s))
      result := add(mload(0x40), 0x20)
      if addDoubleQuotes {
        mstore8(result, 34)
        result := add(1, result)
      }
      // Store "\\u0000" in scratch space.
      // Store "0123456789abcdef" in scratch space.
      // Also, store `{0x08:"b", 0x09:"t", 0x0a:"n", 0x0c:"f", 0x0d:"r"}`.
      // into the scratch space.
      mstore(0x15, 0x5c75303030303031323334353637383961626364656662746e006672)
      // Bitmask for detecting `["\"","\\"]`.
      let e := or(shl(0x22, 1), shl(0x5c, 1))
      for {

      } iszero(eq(s, end)) {

      } {
        s := add(s, 1)
        let c := and(mload(s), 0xff)
        if iszero(lt(c, 0x20)) {
          if iszero(and(shl(c, 1), e)) {
            // Not in `["\"","\\"]`.
            mstore8(result, c)
            result := add(result, 1)
            continue
          }
          mstore8(result, 0x5c) // "\\".
          mstore8(add(result, 1), c)
          result := add(result, 2)
          continue
        }
        if iszero(and(shl(c, 1), 0x3700)) {
          // Not in `["\b","\t","\n","\f","\d"]`.
          mstore8(0x1d, mload(shr(4, c))) // Hex value.
          mstore8(0x1e, mload(and(c, 15))) // Hex value.
          mstore(result, mload(0x19)) // "\\u00XX".
          result := add(result, 6)
          continue
        }
        mstore8(result, 0x5c) // "\\".
        mstore8(add(result, 1), mload(add(c, 8)))
        result := add(result, 2)
      }
      if addDoubleQuotes {
        mstore8(result, 34)
        result := add(1, result)
      }
      let last := result
      mstore(last, 0) // Zeroize the slot after the string.
      result := mload(0x40)
      mstore(result, sub(last, add(result, 0x20))) // Store the length.
      mstore(0x40, add(last, 0x20)) // Allocate the memory.
    }
  }

  /// @dev Escapes the string to be used within double-quotes in a JSON.
  function escapeJSON(string memory s) internal pure returns (string memory result) {
    result = escapeJSON(s, false);
  }

  /// @dev Returns whether `a` equals `b`.
  function eq(string memory a, string memory b) internal pure returns (bool result) {
    /// @solidity memory-safe-assembly
    assembly {
      result := eq(keccak256(add(a, 0x20), mload(a)), keccak256(add(b, 0x20), mload(b)))
    }
  }

  /// @dev Returns whether `a` equals `b`, where `b` is a null-terminated small string.
  function eqs(string memory a, bytes32 b) internal pure returns (bool result) {
    /// @solidity memory-safe-assembly
    assembly {
      // These should be evaluated on compile time, as far as possible.
      let m := not(shl(7, div(not(iszero(b)), 255))) // `0x7f7f ...`.
      let x := not(or(m, or(b, add(m, and(b, m)))))
      let r := shl(7, iszero(iszero(shr(128, x))))
      r := or(r, shl(6, iszero(iszero(shr(64, shr(r, x))))))
      r := or(r, shl(5, lt(0xffffffff, shr(r, x))))
      r := or(r, shl(4, lt(0xffff, shr(r, x))))
      r := or(r, shl(3, lt(0xff, shr(r, x))))
      // forgefmt: disable-next-item
      result := gt(
        eq(mload(a), add(iszero(x), xor(31, shr(3, r)))),
        xor(shr(add(8, r), b), shr(add(8, r), mload(add(a, 0x20))))
      )
    }
  }

  /// @dev Packs a single string with its length into a single word.
  /// Returns `bytes32(0)` if the length is zero or greater than 31.
  function packOne(string memory a) internal pure returns (bytes32 result) {
    /// @solidity memory-safe-assembly
    assembly {
      // We don't need to zero right pad the string,
      // since this is our own custom non-standard packing scheme.
      result := mul(
        // Load the length and the bytes.
        mload(add(a, 0x1f)),
        // `length != 0 && length < 32`. Abuses underflow.
        // Assumes that the length is valid and within the block gas limit.
        lt(sub(mload(a), 1), 0x1f)
      )
    }
  }

  /// @dev Unpacks a string packed using {packOne}.
  /// Returns the empty string if `packed` is `bytes32(0)`.
  /// If `packed` is not an output of {packOne}, the output behavior is undefined.
  function unpackOne(bytes32 packed) internal pure returns (string memory result) {
    /// @solidity memory-safe-assembly
    assembly {
      // Grab the free memory pointer.
      result := mload(0x40)
      // Allocate 2 words (1 for the length, 1 for the bytes).
      mstore(0x40, add(result, 0x40))
      // Zeroize the length slot.
      mstore(result, 0)
      // Store the length and bytes.
      mstore(add(result, 0x1f), packed)
      // Right pad with zeroes.
      mstore(add(add(result, 0x20), mload(result)), 0)
    }
  }

  /// @dev Packs two strings with their lengths into a single word.
  /// Returns `bytes32(0)` if combined length is zero or greater than 30.
  function packTwo(string memory a, string memory b) internal pure returns (bytes32 result) {
    /// @solidity memory-safe-assembly
    assembly {
      let aLength := mload(a)
      // We don't need to zero right pad the strings,
      // since this is our own custom non-standard packing scheme.
      result := mul(
        // Load the length and the bytes of `a` and `b`.
        or(shl(shl(3, sub(0x1f, aLength)), mload(add(a, aLength))), mload(sub(add(b, 0x1e), aLength))),
        // `totalLength != 0 && totalLength < 31`. Abuses underflow.
        // Assumes that the lengths are valid and within the block gas limit.
        lt(sub(add(aLength, mload(b)), 1), 0x1e)
      )
    }
  }

  /// @dev Unpacks strings packed using {packTwo}.
  /// Returns the empty strings if `packed` is `bytes32(0)`.
  /// If `packed` is not an output of {packTwo}, the output behavior is undefined.
  function unpackTwo(bytes32 packed) internal pure returns (string memory resultA, string memory resultB) {
    /// @solidity memory-safe-assembly
    assembly {
      // Grab the free memory pointer.
      resultA := mload(0x40)
      resultB := add(resultA, 0x40)
      // Allocate 2 words for each string (1 for the length, 1 for the byte). Total 4 words.
      mstore(0x40, add(resultB, 0x40))
      // Zeroize the length slots.
      mstore(resultA, 0)
      mstore(resultB, 0)
      // Store the lengths and bytes.
      mstore(add(resultA, 0x1f), packed)
      mstore(add(resultB, 0x1f), mload(add(add(resultA, 0x20), mload(resultA))))
      // Right pad with zeroes.
      mstore(add(add(resultA, 0x20), mload(resultA)), 0)
      mstore(add(add(resultB, 0x20), mload(resultB)), 0)
    }
  }

  /// @dev Directly returns `a` without copying.
  function directReturn(string memory a) internal pure {
    assembly {
      // Assumes that the string does not start from the scratch space.
      let retStart := sub(a, 0x20)
      let retSize := add(mload(a), 0x40)
      // Right pad with zeroes. Just in case the string is produced
      // by a method that doesn't zero right pad.
      mstore(add(retStart, retSize), 0)
      // Store the return offset.
      mstore(retStart, 0x20)
      // End the transaction, returning the string.
      return(retStart, retSize)
    }
  }
}
