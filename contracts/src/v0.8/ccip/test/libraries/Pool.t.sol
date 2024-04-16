// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Pool} from "../../libraries/Pool.sol";
import {BaseTest} from "../BaseTest.t.sol";
import "forge-std/console.sol";

contract Pool__generatePoolReturnDataV1 is BaseTest {
  function test__generatePoolReturnDataV1_Success() public view {
    bytes memory remotePoolAddress = abi.encode(address(this));
    bytes memory destPoolData = abi.encode(address(this));

    bytes memory generatedReturnData = Pool._generatePoolReturnDataV1(remotePoolAddress, destPoolData);

    Pool.PoolReturnDataV1 memory poolReturnDataV1 = Pool._decodePoolReturnDataV1(generatedReturnData);

    assertEq(poolReturnDataV1.destPoolAddress, remotePoolAddress);
    assertEq(poolReturnDataV1.destPoolData, destPoolData);
  }

  function test_fuzz__generatePoolReturnDataV1_Success(
    bytes memory destPoolData,
    bytes memory remotePoolAddress
  ) public pure {
    bytes memory generatedReturnData = Pool._generatePoolReturnDataV1(remotePoolAddress, destPoolData);

    Pool.PoolReturnDataV1 memory poolReturnDataV1 = Pool._decodePoolReturnDataV1(generatedReturnData);

    assertEq(poolReturnDataV1.destPoolAddress, remotePoolAddress);
    assertEq(poolReturnDataV1.destPoolData, destPoolData);
  }
}

contract Pool__decodePoolReturnDataV1 is BaseTest {
  function test__decodePoolReturnDataV1_InvalidTag_Revert() public {
    bytes memory remotePoolAddress = abi.encode(address(this));
    bytes memory destPoolData = abi.encode(address(this));

    bytes memory generatedReturnData = Pool._generatePoolReturnDataV1(remotePoolAddress, destPoolData);

    generatedReturnData[0] = 0x00;

    vm.expectRevert(abi.encodeWithSelector(Pool.InvalidTag.selector, bytes4(generatedReturnData)));
    Pool._decodePoolReturnDataV1(generatedReturnData);
  }
}

contract Pool__removeFirstFourBytes is BaseTest {
  function test_fuzz__removeFirstFourBytes_Success(bytes calldata data) public {
    if (data.length < 4) {
      vm.expectRevert(abi.encodeWithSelector(Pool.MalformedPoolReturnData.selector, data));
      Pool._removeFirstFourBytes(data);
      return;
    }

    bytes memory result = Pool._removeFirstFourBytes(data);
    assertEq(result, data[4:]);
  }

  function test__removeFirstFourBytes_EmptyData_Success() public {
    bytes memory input = abi.encodePacked(bytes4(0x00112233));

    assertEq(input.length, 4);
    bytes memory result = Pool._removeFirstFourBytes(input);

    assertEq(result.length, 0);
    assertEq(bytes(""), result);
  }

  function test__removeFirstFourBytes_MalformedPoolReturnData_Revert() public {
    bytes memory input = abi.encodePacked(bytes1(0x84));

    vm.expectRevert(abi.encodeWithSelector(Pool.MalformedPoolReturnData.selector, input));
    Pool._removeFirstFourBytes(input);
  }
}
