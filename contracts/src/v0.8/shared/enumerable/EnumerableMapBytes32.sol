// SPDX-License-Identifier: MIT
/* solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore */
pragma solidity ^0.8.0;

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/**
 * @dev Library for managing an enumerable variant of Solidity's
 * https://solidity.readthedocs.io/en/latest/types.html#mapping-types[`mapping`]
 * type.
 *
 * Maps have the following properties:
 *
 * - Entries are added, removed, and checked for existence in constant time
 * (O(1)).
 * - Entries are enumerated in O(n). No guarantees are made on the ordering.
 *
 * ```
 * contract Example {
 *     // Add the library methods
 *     using EnumerableMapBytes32 for EnumerableMapBytes32.Bytes32ToBytesMap;
 *
 *     // Declare a set state variable
 *     EnumerableMapBytes32.Bytes32ToBytesMap private myMap;
 * }
 * ```
 *
 * The following map types are supported:
 *
 * - `bytes32 -> bytes` (`Bytes32ToBytes`)
 *
 * [WARNING]
 * ====
 * Trying to delete such a structure from storage will likely result in data corruption, rendering the structure
 * unusable.
 * See https://github.com/ethereum/solidity/pull/11843[ethereum/solidity#11843] for more info.
 *
 * In order to clean up an EnumerableMapBytes32, you should remove all elements one by one.
 * ====
 */
library EnumerableMapBytes32 {
  using EnumerableSet for EnumerableSet.Bytes32Set;

  error NonexistentKeyError();

  struct Bytes32ToBytesMap {
    EnumerableSet.Bytes32Set _keys;
    mapping(bytes32 => bytes) _values;
  }

  /**
   * @dev Adds a key-value pair to a map, or updates the value for an existing
   * key. O(1).
   *
   * Returns true if the key was added to the map, that is if it was not
   * already present.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function set(Bytes32ToBytesMap storage map, bytes32 key, bytes memory value) internal returns (bool) {
    map._values[key] = value;
    return map._keys.add(key);
  }

  /**
   * @dev Removes a key-value pair from a map. O(1).
   *
   * Returns true if the key was removed from the map, that is if it was present.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function remove(Bytes32ToBytesMap storage map, bytes32 key) internal returns (bool) {
    delete map._values[key];
    return map._keys.remove(key);
  }

  /**
   * @dev Returns true if the key is in the map. O(1).
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function contains(Bytes32ToBytesMap storage map, bytes32 key) internal view returns (bool) {
    return map._keys.contains(key);
  }

  /**
   * @dev Returns the number of key-value pairs in the map. O(1).
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function length(Bytes32ToBytesMap storage map) internal view returns (uint256) {
    return map._keys.length();
  }

  /**
   * @dev Returns the key-value pair stored at position `index` in the map. O(1).
   *
   * Note that there are no guarantees on the ordering of entries inside the
   * array, and it may change when more entries are added or removed.
   *
   * Requirements:
   *
   * - `index` must be strictly less than {length}.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function at(Bytes32ToBytesMap storage map, uint256 index) internal view returns (bytes32, bytes memory) {
    bytes32 key = map._keys.at(index);
    return (key, map._values[key]);
  }

  /**
   * @dev Tries to returns the value associated with `key`. O(1).
   * Does not revert if `key` is not in the map.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function tryGet(Bytes32ToBytesMap storage map, bytes32 key) internal view returns (bool, bytes memory) {
    bytes memory value = map._values[key];
    if (value.length == 0) {
      return (contains(map, key), bytes(""));
    } else {
      return (true, value);
    }
  }

  /**
   * @dev Returns the value associated with `key`. O(1).
   *
   * Requirements:
   *
   * - `key` must be in the map.
   */
  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function get(Bytes32ToBytesMap storage map, bytes32 key) internal view returns (bytes memory) {
    bytes memory value = map._values[key];
    if (value.length == 0 && !contains(map, key)) {
      revert NonexistentKeyError();
    }
    return value;
  }
}
