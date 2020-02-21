pragma solidity ^0.6.0;

import "./vendor/SafeMath.sol";
import "./dev/SignedSafeMath.sol";

library Median {
  using SignedSafeMath for int256;

  int256 constant intMax = 57896044618658097711785492504343953926634992332820282019728792003956564819967;

  /**
   * @notice Returns the sorted middle, or the average of the two middle indexed items if the
   * array has an even number of elements.
   * @dev The list passed as an argument isn't modified.
   * @dev This algorithm has expected runtime O(n), but for adversarially chosen inputs
   * the runtime is O(n^2).
   * @param list The list of elements to compare
   */
  function calculate(int256[] memory list)
    internal
    pure
    returns (int256)
  {
    require(0 < list.length, "list must not be empty");
    if (list.length <= 9) {
      return shortList(list);
    } else {
      return longList(copy(list));
    }
  }

  /**
   * @notice See documentation for function calculate.
   * @dev The list passed as an argument may be permuted.
   */
  function calculateInplace(int256[] memory list)
    internal
    pure
    returns (int256)
  {
    require(0 < list.length, "list must not be empty");
    if (list.length <= 9) {
      return shortList(list);
    } else {
      return longList(list);
    }
  }

  /**
   * @notice Optimized median computation for lists of length at most 9
   * @dev Assumes that 0 < list.len <= 9
   * @dev Does not modify list
   */
  function shortList(int256[] memory list) private pure returns (int256) {
    // Uses an optimal sorting network (https://en.wikipedia.org/wiki/Sorting_network)
    // for lists of length 9. Network layout is taken from https://stackoverflow.com/a/46801450

    uint256 len = list.length;
    int256 x0 = list[0];
    if (len == 1) {return x0;}
    // --- end of subnetwork for lists of length <= 1
    int256 x1 = list[1];
    if (x0 > x1) {(x0, x1) = (x1, x0);}
    if (len == 2) {return SignedSafeMath.avg(x0, x1);}
    // --- end of subnetwork for lists of length <= 2
    int256 x2 = list[2];
    if (x1 > x2) {(x1, x2) = (x2, x1);}
    if (x0 > x1) {(x0, x1) = (x1, x0);}
    if (len == 3) {return x1;}
    // --- end of subnetwork for lists of length <= 3
    int256 x3 = list[3];
    int256 x4 = 4 < len ? list[4] : intMax;
    int256 x5 = 5 < len ? list[5] : intMax;
    int256 x6 = 6 < len ? list[6] : intMax;
    int256 x7 = 7 < len ? list[7] : intMax;
    int256 x8 = 8 < len ? list[8] : intMax;
    if (x3 > x4) {(x3, x4) = (x4, x3);}
    if (x6 > x7) {(x6, x7) = (x7, x6);}
    if (x4 > x5) {(x4, x5) = (x5, x4);}
    if (x7 > x8) {(x7, x8) = (x8, x7);}
    if (x3 > x4) {(x3, x4) = (x4, x3);}
    if (x6 > x7) {(x6, x7) = (x7, x6);}
    if (x0 > x3) {(x0, x3) = (x3, x0);}
    if (x3 > x6) {(x3, x6) = (x6, x3);}
    if (x0 > x3) {(x0, x3) = (x3, x0);}
    if (x1 > x4) {(x1, x4) = (x4, x1);}
    if (x4 > x7) {(x4, x7) = (x7, x4);}
    if (x1 > x4) {(x1, x4) = (x4, x1);}
    if (x5 > x8) {(x5, x8) = (x8, x5);}
    if (x2 > x5) {(x2, x5) = (x5, x2);}
    if (x2 > x4) {(x2, x4) = (x4, x2);}
    if (x4 > x6) {(x4, x6) = (x6, x4);}
    if (x2 > x4) {(x2, x4) = (x4, x2);}
    if (x1 > x3) {(x1, x3) = (x3, x1);}
    if (x2 > x3) {(x2, x3) = (x3, x2);}
    // Since we don't care about fully sorting list, but only want the median,
    // some unnecessary comparators have been commented out below.
    // if (x5 > x8) {(x5, x8) = (x8, x5);}
    // if (x5 > x7) {(x5, x7) = (x7, x5);}
    // if (x5 > x6) {(x5, x6) = (x6, x5);}
    if (len == 4) {return SignedSafeMath.avg(x1, x2);}
    if (len == 5) {return x2;}
    if (len == 6) {return SignedSafeMath.avg(x2, x3);}
    if (len == 7) {return x3;}
    if (len == 8) {return SignedSafeMath.avg(x3, x4);}
    if (len == 9) {return x4;}
    revert("list.length > 9");
  }

  /**
   * @notice Median computation for lists of any length
   */
  function longList(int256[] memory list)
    private
    pure
    returns (int256)
  {
    uint256 len = list.length;
    uint256 middleIndex = len / 2;
    if (len % 2 == 0) {
      int256 median1;
      int256 median2;
      (median1, median2) = quickselectTwo(list, 0, len - 1, middleIndex - 1, middleIndex);
      return SignedSafeMath.avg(median1, median2);
    } else {
      return quickselect(list, 0, len - 1, middleIndex);
    }
  }

  /**
   * @notice Selects the k-th ranked element from list, looking only at indices between lo and hi
   * (inclusive). Modifies list in-place.
   */
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

  /**
   * @notice Selects the k1-th and k2-th ranked elements from list, looking only at indices between
   * lo and hi (inclusive). Modifies list in-place.
   */
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

  /**
   * @notice Partitions list in-place using Hoare's partitioning scheme.
   * Only elements of list between indices lo and hi (inclusive) will be modified.
   * Returns an index i, such that:
   * - lo <= i < hi
   * - forall j in [lo, i]. list[j] <= list[i]
   * - forall j in [i, hi]. list[i] <= list[j]
   */
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

  /**
   * @notice Makes an in-memory copy of the array passed in
   * @param list Reference to the array to be copied
   */
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
