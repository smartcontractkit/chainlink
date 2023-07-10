// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {ByteUtil} from "../internal/ByteUtil.sol";

contract ByteUtilTest is Test {

    using ByteUtil for bytes;

    bytes internal constant B_512 = hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000";
    bytes internal constant B_128 = hex"ffffffffffffffffffffffffffffffff";
    bytes internal constant B_EMPTY = new bytes(0);

    bytes4 internal constant MALFORMED_ERROR_SELECTOR = bytes4(keccak256("MalformedData()"));

    function test_readUint256Max() public {
        //read the first 32 bytes
        uint256 result = B_512.readUint256(0);

        //the result should be the max value of a uint256
        assertEq(result, type(uint256).max);
    }

    function test_readUint256Min() public {
        //read the second 32 bytes
        uint256 result = B_512.readUint256(32);

        //the result should be the min value of a uint256
        assertEq(result, type(uint256).min);
    }

    function test_readUint256MultiWord() public {
        //read the first 32 bytes
        uint256 result = B_512.readUint256(31);

        //the result should be the last byte from the first word (ff), and 31 bytes from the second word (0000) (0xFF...0000)
        assertEq(result, type(uint256).max << 248);
    }

    function test_readUint256WithNotEnoughBytes() public {
        //should revert if there's not enough bytes
        vm.expectRevert(MALFORMED_ERROR_SELECTOR);

        //try and read 32 bytes from a 16 byte number
        uint256 result = B_128.readUint256(0);
    }

    function test_readUint256WithEmptyArray() public {
        //should revert if there's not enough bytes
        vm.expectRevert(MALFORMED_ERROR_SELECTOR);

        //read 20 bytes from an empty array
        uint256 result = B_EMPTY.readUint256(0);
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
        address result = B_128.readAddress(0);
    }

    function test_readAddressWithEmptyArray() public {
        //should revert if there's not enough bytes
        vm.expectRevert(MALFORMED_ERROR_SELECTOR);

        //read the first 20 bytes of an empty array
        address result = B_EMPTY.readAddress(0);
    }

}
