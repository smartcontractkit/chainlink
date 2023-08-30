// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

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
  function hasDuplicates(Common.AddressAndWeight[] memory recipients) internal pure returns (bool) {
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
    * @notice Checks if an array of addresses has duplicate addresses
    * @param addressList The array of addresses to check
    * @return bool True if there are duplicates, false otherwise
    */
  function hasDuplicates(address[] memory addressList) internal pure returns (bool) {
    for (uint256 i; i < addressList.length; ) {
      for (uint256 j = i + 1; j < addressList.length; ) {
        if (addressList[i] == addressList[j]) {
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

  function contains(address[] memory addressList, address addr) internal pure returns (bool, uint256) {
    uint256 i;
    for (; i < addressList.length; ) {
      if (addressList[i] == addr) {
        return (true, i);
      }
      unchecked {
        ++i;
      }
    }
    return (false, 0);
  }
}
