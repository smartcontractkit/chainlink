// SPDX-License-Identifier: MIT
/* solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore */
pragma solidity ^0.8.0;

/// Library for managing sets of bytes. Reuses OpenZeppelin's EnumerableSet library logic but for bytes.
library EnumerableBytesSet {
  struct BytesSet {
    bytes[] _values;
    mapping(bytes value => uint256) _positions;
  }

  /// @dev Adds a value to a set. O(1).
  /// @param set The set to add the value to.
  /// @param value The value to add.
  /// @return True if the value was added to the set, false if the value was already in the set.
  function add(BytesSet storage set, bytes memory value) internal returns (bool) {
    return _add(set, value);
  }

  function _add(BytesSet storage set, bytes memory value) private returns (bool) {
    if (!_contains(set, value)) {
      set._values.push(value);
      // The value is stored at length-1, but we add 1 to all indexes
      // and use 0 as a sentinel value
      set._positions[value] = set._values.length;
      return true;
    } else {
      return false;
    }
  }

  /// @dev Removes a value from a set. O(1).
  /// @param set The set to remove the value from.
  /// @param value The value to remove.
  /// @return True if the value was removed from the set, false if the value was not in the set.
  function remove(BytesSet storage set, bytes memory value) internal returns (bool) {
    return _remove(set, value);
  }

  function _remove(BytesSet storage set, bytes memory value) private returns (bool) {
    // We cache the value's position to prevent multiple reads from the same storage slot
    uint256 position = set._positions[value];

    if (position != 0) {
      // Equivalent to contains(set, value)
      // To delete an element from the _values array in O(1), we swap the element to delete with the last one in
      // the array, and then remove the last element (sometimes called as 'swap and pop').
      // This modifies the order of the array, as noted in {at}.

      uint256 valueIndex = position - 1;
      uint256 lastIndex = set._values.length - 1;

      if (valueIndex != lastIndex) {
        bytes memory lastValue = set._values[lastIndex];

        // Move the lastValue to the index where the value to delete is
        set._values[valueIndex] = lastValue;
        // Update the tracked position of the lastValue (that was just moved)
        set._positions[lastValue] = position;
      }

      // Delete the slot where the moved value was stored
      set._values.pop();

      // Delete the tracked position for the deleted slot
      delete set._positions[value];

      return true;
    } else {
      return false;
    }
  }

  /// @dev Checks if a value is in a set. O(1).
  /// @param set The set to check the value in.
  /// @param value The value to check.
  /// @return True if the value is in the set, false otherwise.
  function contains(BytesSet storage set, bytes memory value) internal view returns (bool) {
    return _contains(set, value);
  }

  function _contains(BytesSet storage set, bytes memory value) private view returns (bool) {
    return set._positions[value] != 0;
  }

  /// @dev Returns the number of values in the set. O(1).
  /// @param set The set to count values in.
  /// @return The number of values in the set.
  function length(BytesSet storage set) internal view returns (uint256) {
    return _length(set);
  }

  function _length(BytesSet storage set) private view returns (uint256) {
    return set._values.length;
  }

  /// @dev Returns the value stored at position `index` in the set. O(1).
  /// Note that there are no guarantees on the ordering of values inside the array, and it may change when more values
  /// are added or removed.
  /// @dev precondition - `index` must be strictly less than {length}.
  /// @param set The set to get the value from.
  /// @param index The index to get the value at.
  /// @return The value stored at the specified index.
  function at(BytesSet storage set, uint256 index) internal view returns (bytes memory) {
    return _at(set, index);
  }

  function _at(BytesSet storage set, uint256 index) private view returns (bytes memory) {
    return set._values[index];
  }

  /// @dev Returns the entire set in an array
  ///
  /// WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed to
  /// mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that this
  /// function has an unbounded cost, and using it as part of a state-changing function may render the function
  /// uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
  /// @param set The set to get the values from.
  ///
  /// @return An array containing all the values in the set.
  function values(BytesSet storage set) internal view returns (bytes[] memory) {
    bytes[] memory store = _values(set);
    bytes[] memory result;

    assembly ("memory-safe") {
      result := store
    }

    return result;
  }

  function _values(BytesSet storage set) private view returns (bytes[] memory) {
    return set._values;
  }
}
