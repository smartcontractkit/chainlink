// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveNodeOperatorsTest is BaseTest {
  event NodeOperatorRemoved(uint256 nodeOperatorId);

  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 0;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 1;

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_RevertWhen_CalledByNonOwner() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    uint256[] memory nodeOperatorsToRemove = new uint256[](2);
    nodeOperatorsToRemove[1] = 1;
    s_capabilityRegistry.removeNodeOperators(nodeOperatorsToRemove);
  }

  function test_RemovesNodeOperator() public {
    changePrank(ADMIN);

    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit NodeOperatorRemoved(TEST_NODE_OPERATOR_ONE_ID);
    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit NodeOperatorRemoved(TEST_NODE_OPERATOR_TWO_ID);
    uint256[] memory nodeOperatorsToRemove = new uint256[](2);
    nodeOperatorsToRemove[1] = 1;
    s_capabilityRegistry.removeNodeOperators(nodeOperatorsToRemove);

    CapabilityRegistry.NodeOperator memory nodeOperatorOne = s_capabilityRegistry.getNodeOperator(0);
    assertEq(nodeOperatorOne.admin, address(0));
    assertEq(nodeOperatorOne.name, "");

    CapabilityRegistry.NodeOperator memory nodeOperatorTwo = s_capabilityRegistry.getNodeOperator(1);
    assertEq(nodeOperatorTwo.admin, address(0));
    assertEq(nodeOperatorTwo.name, "");
  }
}
