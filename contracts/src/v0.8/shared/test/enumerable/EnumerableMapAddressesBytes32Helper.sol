// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {EnumerableMapAddresses} from "../../enumerable/EnumerableMapAddresses.sol";

contract EnumerableMapAddressesBytes32Helper {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;

  EnumerableMapAddresses.AddressToBytes32Map internal s_map;

  function set(address key, bytes32 value) external returns (bool) {
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

  function at(uint256 index) external view returns (address, bytes32) {
    return s_map.at(index);
  }

  function tryGet(address key) external view returns (bool, bytes32) {
    return s_map.tryGet(key);
  }

  function get(address key) external view returns (bytes32) {
    return s_map.get(key);
  }
}
