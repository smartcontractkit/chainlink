// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {ByteUtil} from "../ByteUtil.sol";

contract ByteUtilTest is Test {
  using ByteUtil for bytes;

  bytes internal constant B_512 =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000";
  bytes internal constant B_128 = hex"ffffffffffffffffffffffffffffffff";
  bytes internal constant B_16 = hex"ffff";
  bytes internal constant B_EMPTY = new bytes(0);

  bytes4 internal constant MALFORMED_ERROR_SELECTOR = bytes4(keccak256("MalformedData()"));

  function test_readUint256Max() public {
    //read the first 32 bytes
    uint256 result = B_512.readUint256(0);

    //the result should be the max value of a uint256
    assertEq(result, type(uint256).max);
  }

  function test_readUint192Max() public {
    //read the first 24 bytes
    uint256 result = B_512.readUint192(0);

    //the result should be the max value of a uint192
    assertEq(result, type(uint192).max);
  }

  function test_readUint32Max() public {
    //read the first 4 bytes
    uint256 result = B_512.readUint32(0);

    //the result should be the max value of a uint32
    assertEq(result, type(uint32).max);
  }

  function test_readUint256Min() public {
    //read the second 32 bytes
    uint256 result = B_512.readUint256(32);

    //the result should be the min value of a uint256
    assertEq(result, type(uint256).min);
  }

  function test_readUint192Min() public {
    //read the second 24 bytes
    uint256 result = B_512.readUint192(32);

    //the result should be the min value of a uint192
    assertEq(result, type(uint192).min);
  }

  function test_readUint32Min() public {
    //read the second 4 bytes
    uint256 result = B_512.readUint32(32);

    //the result should be the min value of a uint32
    assertEq(result, type(uint32).min);
  }

  function test_readUint256MultiWord() public {
    //read the first 32 bytes
    uint256 result = B_512.readUint256(31);

    //the result should be the last byte from the first word (ff), and 31 bytes from the second word (0000) (0xFF...0000)
    assertEq(result, type(uint256).max << 248);
  }

  function test_readUint192MultiWord() public {
    //read the first 24 bytes
    uint256 result = B_512.readUint192(31);

    //the result should be the last byte from the first word (ff), and 23 bytes from the second word (0000) (0xFF...0000)
    assertEq(result, type(uint192).max << 184);
  }

  function test_readUint32MultiWord() public {
    //read the first 4 bytes
    uint256 result = B_512.readUint32(31);

    //the result should be the last byte from the first word (ff), and 3 bytes from the second word (0000) (0xFF...0000)
    assertEq(result, type(uint32).max << 24);
  }

  function test_readUint256WithNotEnoughBytes() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //try and read 32 bytes from a 16 byte number
    B_128.readUint256(0);
  }

  function test_readUint192WithNotEnoughBytes() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //try and read 24 bytes from a 16 byte number
    B_128.readUint192(0);
  }

  function test_readUint32WithNotEnoughBytes() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //try and read 4 bytes from a 2 byte number
    B_16.readUint32(0);
  }

  function test_readUint256WithEmptyArray() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //read 20 bytes from an empty array
    B_EMPTY.readUint256(0);
  }

  function test_readUint192WithEmptyArray() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //read 20 bytes from an empty array
    B_EMPTY.readUint192(0);
  }

  function test_readUint32WithEmptyArray() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //read 20 bytes from an empty array
    B_EMPTY.readUint32(0);
  }

  function test_readAddress() public {
    //read the first 20 bytes
    address result = B_512.readAddress(0);

    //the result should be the max value of a uint256
    assertEq(result, address(type(uint160).max));
  }

  function test_readZeroAddress() public {
    //read the first 32 bytes after the first word
    address result = B_512.readAddress(32);

    //the result should be 0x00...0
    assertEq(result, address(type(uint160).min));
  }

  function test_readAddressMultiWord() public {
    //read the first 20 bytes after byte 13
    address result = B_512.readAddress(13);

    //the result should be the value last 19 bytes of the first word (ffff..) and the first byte of the second word (00) (0xFFFF..00)
    assertEq(result, address(type(uint160).max << 8));
  }

  function test_readAddressWithNotEnoughBytes() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //read 20 bytes from a 16 byte array
    B_128.readAddress(0);
  }

  function test_readAddressWithEmptyArray() public {
    //should revert if there's not enough bytes
    vm.expectRevert(MALFORMED_ERROR_SELECTOR);

    //read the first 20 bytes of an empty array
    B_EMPTY.readAddress(0);
  }
}
