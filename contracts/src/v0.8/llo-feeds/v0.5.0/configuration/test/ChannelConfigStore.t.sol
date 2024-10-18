// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {IChannelConfigStore} from "../interfaces/IChannelConfigStore.sol";
import {Test} from "forge-std/Test.sol";
import {ChannelConfigStore} from "../ChannelConfigStore.sol";
import {ExposedChannelConfigStore} from "./mocks/ExposedChannelConfigStore.sol";

/**
 * @title ChannelConfigStoreTest
 * @author samsondav
 * @notice Base class for ChannelConfigStore tests
 */
contract ChannelConfigStoreTest is Test {
  ExposedChannelConfigStore public channelConfigStore;
  event NewChannelDefinition(uint256 indexed donId, uint32 version, string url, bytes32 sha);

  function setUp() public virtual {
    channelConfigStore = new ExposedChannelConfigStore();
  }

  function testTypeAndVersion() public view {
    assertEq(channelConfigStore.typeAndVersion(), "ChannelConfigStore 0.0.1");
  }

  function testSupportsInterface() public view {
    assertTrue(channelConfigStore.supportsInterface(type(IChannelConfigStore).interfaceId));
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");

    vm.startPrank(address(2));
    channelConfigStore.setChannelDefinitions(42, "url", keccak256("sha"));
  }

  function testSetChannelDefinitions() public {
    vm.expectEmit();
    emit NewChannelDefinition(42, 1, "url", keccak256("sha"));
    channelConfigStore.setChannelDefinitions(42, "url", keccak256("sha"));

    vm.expectEmit();
    emit NewChannelDefinition(42, 2, "url2", keccak256("sha2"));
    channelConfigStore.setChannelDefinitions(42, "url2", keccak256("sha2"));

    assertEq(channelConfigStore.exposedReadChannelDefinitionStates(42), uint32(2));
  }
}
