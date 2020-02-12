pragma solidity ^0.6.0;

import "./vendor/SafeMath.sol";
import "./vendor/SignedSafeMath.sol";

library Median {
  using SignedSafeMath for int256;

  /// @notice Returns the sorted middle, or the average of the two middle indexed items if the
  /// array has an even number of elements.
  /// @dev The list passed as an argument isn't modified.
  /// @dev This algorithm has expected runtime O(n), but for adversarially chosen inputs
  /// the runtime is O(n^2).
  /// @param list The list of elements to compare
  function calculate(int256[] memory list)
    internal
    pure
    returns (int256)
  {
    return calculateInplace(copy(list));
  }

  /// @notice See documentation for function calculate.
  /// @dev The list passed as an argument *is* modified.
  function calculateInplace(int256[] memory list)
    internal
    pure
    returns (int256)
  {
    uint256 len = list.length;
    require(0 < len, "input must not be empty");
    uint256 middleIndex = len / 2;
    if (len % 2 == 0) {
      int256 median1;
      int256 median2;
      (median1, median2) = quickselectTwo(list, 0, len - 1, middleIndex - 1, middleIndex);
      return safeAvg(median1, median2);
    } else {
      return quickselect(list, 0, len - 1, middleIndex);
    }
  }

  /// @notice Selects the k-th ranked element from list, looking only at indices between lo and hi
  /// (inclusive). Modifies list in-place.
  function quickselect(int256[] memory list, uint lo, uint hi, uint k)
    private
    pure
    returns (int256)
  {
    require(lo <= k);
    require(k <= hi);
    while (lo < hi) {
      uint pivotIndex = partition(list, lo, hi);
      if (k <= pivotIndex) {
        // since pivotIndex < (original hi passed to partition),
        // termination is guaranteed in this case
        hi = pivotIndex;
      } else {
        // since (original lo passed to partition) <= pivotIndex,
        // termination is guaranteed in this case
        lo = pivotIndex + 1;
      }
    }
    return list[lo];
  }

  /// @notice Selects the k1-th and k2-th ranked elements from list, looking only at indices between
  /// lo and hi (inclusive). Modifies list in-place.
  function quickselectTwo(int256[] memory list, uint lo, uint hi, uint k1, uint k2)
    internal
    pure
    returns (int256, int256)
  {
    require(k1 < k2);
    require(lo <= k1 && k1 <= hi);
    require(lo <= k2 && k2 <= hi);

    while (true) {
      uint pivotIdx = partition(list, lo, hi);
      if (k2 <= pivotIdx) {
        hi = pivotIdx;
      } else if (pivotIdx < k1) {
        lo = pivotIdx + 1;
      } else {
        assert(k1 <= pivotIdx && pivotIdx < k2);
        int256 r1 = quickselect(list, lo, pivotIdx, k1);
        int256 r2 = quickselect(list, pivotIdx + 1, hi, k2);
        return (r1, r2);
      }
    }
  }

  /// @notice Partitions list in-place using Hoare's partitioning scheme.
  /// Only elements of list between indices lo and hi (inclusive) will be modified.
  /// Returns an index i, such that:
  /// - lo <= i < hi
  /// - forall j in [lo, i]. list[j] <= list[i]
  /// - forall j in [i, hi]. list[i] <= list[j]
  function partition(int256[] memory list, uint lo, uint hi)
    private
    pure
    returns (uint256)
  {
    // We don't care about overflow of the addition, because it would require a list
    // larger than any feasible computer's memory.
    int256 pivot = list[(lo + hi) / 2];
    lo -= 1; // this can underflow. that's intentional.
    hi += 1;
    while (true) {
      do {
        lo += 1;
      } while (list[lo] < pivot);
      do {
        hi -= 1;
      } while (list[hi] > pivot);
      if (lo < hi) {
        (list[lo], list[hi]) = (list[hi], list[lo]);
      } else {
        // Let orig_lo and orig_hi be the original values of lo and hi passed to partition.
        // Then, hi < orig_hi, because hi decreases *strictly* monotonically
        // in each loop iteration and
        // - either list[orig_hi] > pivot, in which case the first loop iteration
        //   will achieve hi < orig_hi;
        // - or list[orig_hi] <= pivot, in which case at least two loop iterations are
        //   needed:
        //   - lo will have to stop at least once in the interval
        //     [orig_lo, (orig_lo + orig_hi)/2]
        //   - (orig_lo + orig_hi)/2 < orig_hi
        return hi;
      }
    }
  }

  /// @notice Computes average of a and b using SignedSafeMath
  function safeAvg(int256 a, int256 b)
    private
    pure
    returns (int256)
  {
    int256 remainder = (a % 2 + b % 2) / 2;
    return (a / 2).add(b / 2).add(remainder);
  }


  /// @notice Makes an in-memory copy of the array passed in
  /// @param list Reference to the array to be copied
  function copy(int256[] memory list)
    private
    pure
    returns(int256[] memory)
  {
    int256[] memory list2 = new int256[](list.length);
    for (uint256 i = 0; i < list.length; i++) {
      list2[i] = list[i];
    }
    return list2;
  }
}
