// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {BaseTest} from "../BaseTest.t.sol";
import {SortedSetValidationUtil} from "../../../shared/util/SortedSetValidationUtil.sol";

contract SortedSetValidationUtilBaseTest is BaseTest {
  uint256 constant OFFSET = 5;

  modifier _ensureSetLength(uint256 subsetLength, uint256 supersetLength) {
    vm.assume(subsetLength > 0 && supersetLength > 0 && subsetLength <= supersetLength);
    _;
  }

  function _createSets(
    uint256 subsetLength,
    uint256 supersetLength
  ) internal pure returns (bytes32[] memory subset, bytes32[] memory superset) {
    subset = new bytes32[](subsetLength);
    superset = new bytes32[](supersetLength);
  }

  function _convertArrayToSortedSet(bytes32[] memory arr, uint256 offSet) internal pure {
    for (uint256 i = 1; i < arr.length; ++i) {
      arr[i] = bytes32(uint256(arr[i - 1]) + offSet);
    }
  }

  function _convertToUnsortedSet(bytes32[] memory arr, uint256 ptr1, uint256 ptr2) internal pure {
    // Swap two elements to make it unsorted
    (arr[ptr1], arr[ptr2]) = (arr[ptr2], arr[ptr1]);
  }

  function _convertArrayToSubset(bytes32[] memory subset, bytes32[] memory superset) internal pure {
    for (uint256 i; i < subset.length; ++i) {
      subset[i] = superset[i];
    }
  }

  function _makeInvalidSubset(bytes32[] memory subset, bytes32[] memory superset, uint256 ptr) internal pure {
    _convertArrayToSubset(subset, superset);
    subset[ptr] = bytes32(uint256(subset[ptr]) + 1);
  }

  function _convertArrayToHaveDuplicates(bytes32[] memory arr, uint256 ptr1, uint256 ptr2) internal pure {
    arr[ptr2] = arr[ptr1];
  }
}

contract SortedSetValidationUtil_CheckIsValidUniqueSubsetTest is SortedSetValidationUtilBaseTest {
  // Successes.

  function test__checkIsValidUniqueSubset_ValidSubset_Success() public pure {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _convertArrayToSubset(subset, superset);

    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  // Reverts.

  function test__checkIsValidUniqueSubset_EmptySubset_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(0, 5);
    _convertArrayToSortedSet(superset, OFFSET);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.EmptySet.selector));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_EmptySuperset_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 0);
    _convertArrayToSortedSet(subset, OFFSET);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.EmptySet.selector));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_NotASubset_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _makeInvalidSubset(subset, superset, 1);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASubset.selector, subset, superset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_UnsortedSubset_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _convertToUnsortedSet(subset, 1, 2);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, subset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_UnsortedSuperset_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _convertArrayToSubset(subset, superset);
    _convertToUnsortedSet(superset, 1, 2);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, superset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_HasDuplicates_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _convertArrayToSubset(subset, superset);
    _convertArrayToHaveDuplicates(subset, 1, 2);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, subset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_SubsetLargerThanSuperset_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(6, 5);
    _convertArrayToSortedSet(subset, OFFSET);
    _convertArrayToSortedSet(superset, OFFSET);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASubset.selector, subset, superset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_SubsetEqualsSuperset_NoRevert() public pure {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(5, 5);
    _convertArrayToSortedSet(subset, OFFSET);
    _convertArrayToSortedSet(superset, OFFSET);

    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_SingleElementSubset() public pure {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(1, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _convertArrayToSubset(subset, superset);

    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_SingleElementSubsetAndSuperset_Equal() public pure {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(1, 1);
    _convertArrayToSortedSet(subset, OFFSET);
    _convertArrayToSortedSet(superset, OFFSET);

    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_SingleElementSubsetAndSuperset_NotEqual_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(1, 1);
    _convertArrayToSortedSet(subset, OFFSET);
    superset[0] = bytes32(uint256(subset[0]) + 10); // Different value

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASubset.selector, subset, superset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }

  function test__checkIsValidUniqueSubset_SupersetHasDuplicates_Reverts() public {
    (bytes32[] memory subset, bytes32[] memory superset) = _createSets(3, 5);
    _convertArrayToSortedSet(superset, OFFSET);
    _convertArrayToSubset(subset, superset);
    _convertArrayToHaveDuplicates(superset, 1, 2);

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, superset));
    SortedSetValidationUtil._checkIsValidUniqueSubset(subset, superset);
  }
}
