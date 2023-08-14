// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {EnumerableMapAddresses} from "../../enumerable/EnumerableMapAddresses.sol";

contract EnumerableMapAddressesHelper {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  EnumerableMapAddresses.AddressToAddressMap internal s_map;

  function set(address key, address value) external returns (bool) {
    return s_map.set(key, value);
  }

  function remove(address key) external returns (bool) {
    return s_map.remove(key);
  }

  function contains(address key) external view returns (bool) {
    return s_map.contains(key);
  }

  function length() external view returns (uint256) {
    return s_map.length();
  }

  function at(uint256 index) external view returns (address, address) {
    return s_map.at(index);
  }

  function tryGet(address key) external view returns (bool, address) {
    return s_map.tryGet(key);
  }

  function get(address key) external view returns (address) {
    return s_map.get(key);
  }

  function get(address key, string memory errorMessage) external view returns (address) {
    return s_map.get(key, errorMessage);
  }
}
