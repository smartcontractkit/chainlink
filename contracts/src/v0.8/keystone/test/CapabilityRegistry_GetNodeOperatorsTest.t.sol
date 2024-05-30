// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetNodeOperatorsTest is BaseTest {
  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 1;
  uint32 private constant TEST_NODE_OPERATOR_TWO_ID = 2;

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_CorrectlyFetchesNodeOperators() public view {
    CapabilityRegistry.NodeOperator[] memory nodeOperators = s_capabilityRegistry.getNodeOperators();
    assertEq(nodeOperators.length, 2);

    assertEq(nodeOperators[0].admin, NODE_OPERATOR_ONE_ADMIN);
    assertEq(nodeOperators[0].name, NODE_OPERATOR_ONE_NAME);

    assertEq(nodeOperators[1].admin, NODE_OPERATOR_TWO_ADMIN);
    assertEq(nodeOperators[1].name, NODE_OPERATOR_TWO_NAME);
  }

  function test_DoesNotIncludeRemovedNodeOperators() public {
    changePrank(ADMIN);
    uint32[] memory nodeOperatorsToRemove = new uint32[](1);
    nodeOperatorsToRemove[0] = 2;
    s_capabilityRegistry.removeNodeOperators(nodeOperatorsToRemove);

    CapabilityRegistry.NodeOperator[] memory nodeOperators = s_capabilityRegistry.getNodeOperators();
    assertEq(nodeOperators.length, 1);

    assertEq(nodeOperators[0].admin, NODE_OPERATOR_ONE_ADMIN);
    assertEq(nodeOperators[0].name, NODE_OPERATOR_ONE_NAME);
  }
}
