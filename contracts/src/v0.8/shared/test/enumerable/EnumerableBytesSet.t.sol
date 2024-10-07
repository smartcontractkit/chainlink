// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {EnumerableBytesSet} from "../../dev/enumerable/EnumerableBytesSet.sol";

import {Test} from "../../../vendor/forge-std/src/Test.sol";

contract EnumerableBytesSetTest is Test {
  function _assertBytesArrayEq(bytes[] memory a, bytes[] memory b) internal {
    assertEq(a.length, b.length);
    for (uint256 i = 0; i < a.length; i++) {
      assertEq(a[i], b[i]);
    }
  }
}

contract EnumerableBytesSetTest_Add is EnumerableBytesSetTest {
  using EnumerableBytesSet for EnumerableBytesSet.BytesSet;

  EnumerableBytesSet.BytesSet private s_set;

  function test_add_SingleValue() public {
    bytes memory value = "value";
    bytes[] memory expected = new bytes[](1);
    expected[0] = value;

    assertFalse(s_set.contains(value));
    assertTrue(s_set.add(value));
    assertEq(s_set.length(), 1);
    assertEq(s_set.at(0), value);
    assertTrue(s_set.contains(value));
    _assertBytesArrayEq(s_set.values(), expected);
  }

  function test_add_AlreadyExistingValue() public {
    bytes memory value = "value";
    bytes[] memory expected = new bytes[](1);
    expected[0] = value;

    assertTrue(s_set.add(value));
    assertFalse(s_set.add(value));
    assertEq(s_set.length(), 1);
    assertEq(s_set.at(0), value);
    assertTrue(s_set.contains(value));
    _assertBytesArrayEq(s_set.values(), expected);
  }

  function test_add_MultipleUniqueValues() public {
    bytes memory value1 = "value1";
    bytes memory value2 = "value2";
    bytes[] memory expected = new bytes[](2);
    expected[0] = value1;
    expected[1] = value2;

    assertTrue(s_set.add(value1));
    assertTrue(s_set.add(value2));
    assertEq(s_set.length(), 2);
    assertTrue(s_set.contains(value1));
    assertTrue(s_set.contains(value2));
    assertEq(s_set.at(0), value1);
    assertEq(s_set.at(1), value2);
    _assertBytesArrayEq(s_set.values(), expected);
  }

  function testFuzz_add(bytes[2] memory values) public {
    bytes[] memory expected = new bytes[](values.length);

    for (uint256 i = 0; i < values.length; ++i) {
      // Ensure uniqueness
      expected[i] = bytes.concat(values[i], abi.encodePacked(i));
      s_set.add(expected[i]);

      assertEq(s_set.at(i), expected[i]);
      assertTrue(s_set.contains(expected[i]));
    }

    assertEq(s_set.length(), values.length);
    _assertBytesArrayEq(s_set.values(), expected);
  }
}

contract EnumerableBytesSet_Remove is EnumerableBytesSetTest {
  using EnumerableBytesSet for EnumerableBytesSet.BytesSet;

  EnumerableBytesSet.BytesSet private s_set;

  function setUp() public {
    s_set.add("value1");
    s_set.add("value2");
  }

  function test_remove_SingleExistingValue() public {
    bytes memory value = "value1";
    bytes[] memory expected = new bytes[](1);
    expected[0] = "value2";

    assertTrue(s_set.remove(value));
    assertEq(s_set.length(), 1);
    assertFalse(s_set.contains(value));
    assertEq(s_set.at(0), "value2");
    _assertBytesArrayEq(s_set.values(), expected);
  }

  function test_remove_MultipleExistingValues() public {
    bytes memory value1 = "value1";
    bytes memory value2 = "value2";
    bytes[] memory expected = new bytes[](0);

    vm.expectRevert();
    assertEq(s_set.at(0), "");
    vm.expectRevert();
    assertEq(s_set.at(1), "");

    assertTrue(s_set.remove(value1));
    assertTrue(s_set.remove(value2));
    assertEq(s_set.length(), 0);
    assertFalse(s_set.contains(value1));
    assertFalse(s_set.contains(value2));
    _assertBytesArrayEq(s_set.values(), expected);
  }

  function test_remove_SingleNonExistingValue() public {
    bytes memory value = "value3";
    bytes[] memory expected = new bytes[](2);
    expected[0] = "value1";
    expected[1] = "value2";

    assertFalse(s_set.remove(value));
    assertEq(s_set.length(), 2);
    assertFalse(s_set.contains(value));
    assertEq(s_set.at(0), "value1");
    assertEq(s_set.at(1), "value2");
    _assertBytesArrayEq(s_set.values(), expected);
  }
}
