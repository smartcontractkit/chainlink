// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

/*
 * @title Common
 * @author Michael Fletcher
 * @notice Common functions and structs
 */
library Common {
  // @notice The asset struct to hold the address of an asset and amount
  struct Asset {
    address assetAddress;
    uint256 amount;
  }

  // @notice Struct to hold the address and its associated weight
  struct AddressAndWeight {
    address addr;
    uint64 weight;
  }

  /**
   * @notice Checks if an array of AddressAndWeight has duplicate addresses
   * @param recipients The array of AddressAndWeight to check
   * @return bool True if there are duplicates, false otherwise
   */
  function _hasDuplicateAddresses(address[] memory recipients) internal pure returns (bool) {
    for (uint256 i = 0; i < recipients.length; ) {
      for (uint256 j = i + 1; j < recipients.length; ) {
        if (recipients[i] == recipients[j]) {
          return true;
        }
        unchecked {
          ++j;
        }
      }
      unchecked {
        ++i;
      }
    }
    return false;
  }

  /**
   * @notice Checks if an array of AddressAndWeight has duplicate addresses
   * @param recipients The array of AddressAndWeight to check
   * @return bool True if there are duplicates, false otherwise
   */
  function _hasDuplicateAddresses(Common.AddressAndWeight[] memory recipients) internal pure returns (bool) {
    for (uint256 i = 0; i < recipients.length; ) {
      for (uint256 j = i + 1; j < recipients.length; ) {
        if (recipients[i].addr == recipients[j].addr) {
          return true;
        }
        unchecked {
          ++j;
        }
      }
      unchecked {
        ++i;
      }
    }
    return false;
  }

  /**
   * @notice sorts a list of addresses numerically
   * @param arr The array of addresses to sort
   * @param left the start index
   * @param right the end index
   */
  function _quickSort(address[] memory arr, int256 left, int256 right) internal pure {
    int256 i = left;
    int256 j = right;
    if (i == j) return;
    address pivot = arr[uint256(left + (right - left) / 2)];
    while (i <= j) {
      while (uint160(arr[uint256(i)]) < uint160(pivot)) i++;
      while (uint160(pivot) < uint160(arr[uint256(j)])) j--;
      if (i <= j) {
        (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
        i++;
        j--;
      }
    }
    if (left < j) _quickSort(arr, left, j);
    if (i < right) _quickSort(arr, i, right);
  }
}
