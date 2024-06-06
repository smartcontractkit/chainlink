// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {EnumerableMap} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableMap.sol";

// TODO: the lib can be replaced with OZ v5.1 post-upgrade, which has AddressToAddressMap and AddressToBytes32Map
library EnumerableMapAddresses {
  using EnumerableMap for EnumerableMap.UintToAddressMap;
  using EnumerableMap for EnumerableMap.Bytes32ToBytes32Map;

  struct AddressToAddressMap {
    EnumerableMap.UintToAddressMap _inner;
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function set(AddressToAddressMap storage map, address key, address value) internal returns (bool) {
    return map._inner.set(uint256(uint160(key)), value);
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function remove(AddressToAddressMap storage map, address key) internal returns (bool) {
    return map._inner.remove(uint256(uint160(key)));
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function contains(AddressToAddressMap storage map, address key) internal view returns (bool) {
    return map._inner.contains(uint256(uint160(key)));
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function length(AddressToAddressMap storage map) internal view returns (uint256) {
    return map._inner.length();
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function at(AddressToAddressMap storage map, uint256 index) internal view returns (address, address) {
    (uint256 key, address value) = map._inner.at(index);
    return (address(uint160(key)), value);
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function tryGet(AddressToAddressMap storage map, address key) internal view returns (bool, address) {
    return map._inner.tryGet(uint256(uint160(key)));
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function get(AddressToAddressMap storage map, address key) internal view returns (address) {
    return map._inner.get(uint256(uint160(key)));
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function get(
    AddressToAddressMap storage map,
    address key,
    string memory errorMessage
  ) internal view returns (address) {
    return map._inner.get(uint256(uint160(key)), errorMessage);
  }

  // AddressToBytes32Map

  struct AddressToBytes32Map {
    EnumerableMap.Bytes32ToBytes32Map _inner;
  }

  /**
   * @dev Adds a key-value pair to a map, or updates the value for an existing
   * key. O(1).
   *
   * Returns true if the key was added to the map, that is if it was not
   * already present.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function set(AddressToBytes32Map storage map, address key, bytes32 value) internal returns (bool) {
    return map._inner.set(bytes32(uint256(uint160(key))), value);
  }

  /**
   * @dev Removes a value from a map. O(1).
   *
   * Returns true if the key was removed from the map, that is if it was present.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function remove(AddressToBytes32Map storage map, address key) internal returns (bool) {
    return map._inner.remove(bytes32(uint256(uint160(key))));
  }

  /**
   * @dev Returns true if the key is in the map. O(1).
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function contains(AddressToBytes32Map storage map, address key) internal view returns (bool) {
    return map._inner.contains(bytes32(uint256(uint160(key))));
  }

  /**
   * @dev Returns the number of elements in the map. O(1).
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function length(AddressToBytes32Map storage map) internal view returns (uint256) {
    return map._inner.length();
  }

  /**
   * @dev Returns the element stored at position `index` in the map. O(1).
   * Note that there are no guarantees on the ordering of values inside the
   * array, and it may change when more values are added or removed.
   *
   * Requirements:
   *
   * - `index` must be strictly less than {length}.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function at(AddressToBytes32Map storage map, uint256 index) internal view returns (address, bytes32) {
    (bytes32 key, bytes32 value) = map._inner.at(index);
    return (address(uint160(uint256(key))), value);
  }

  /**
   * @dev Tries to returns the value associated with `key`. O(1).
   * Does not revert if `key` is not in the map.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function tryGet(AddressToBytes32Map storage map, address key) internal view returns (bool, bytes32) {
    (bool success, bytes32 value) = map._inner.tryGet(bytes32(uint256(uint160(key))));
    return (success, value);
  }

  /**
   * @dev Returns the value associated with `key`. O(1).
   *
   * Requirements:
   *
   * - `key` must be in the map.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function get(AddressToBytes32Map storage map, address key) internal view returns (bytes32) {
    return map._inner.get(bytes32(uint256(uint160(key))));
  }
}
