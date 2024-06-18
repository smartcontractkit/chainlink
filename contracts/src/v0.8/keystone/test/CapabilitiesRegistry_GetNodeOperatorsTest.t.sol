// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_GetNodeOperatorsTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_CapabilitiesRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_CorrectlyFetchesNodeOperators() public view {
    CapabilitiesRegistry.NodeOperator[] memory nodeOperators = s_CapabilitiesRegistry.getNodeOperators();
    assertEq(nodeOperators.length, 3);

    assertEq(nodeOperators[0].admin, NODE_OPERATOR_ONE_ADMIN);
    assertEq(nodeOperators[0].name, NODE_OPERATOR_ONE_NAME);

    assertEq(nodeOperators[1].admin, NODE_OPERATOR_TWO_ADMIN);
    assertEq(nodeOperators[1].name, NODE_OPERATOR_TWO_NAME);

    assertEq(nodeOperators[2].admin, NODE_OPERATOR_THREE);
    assertEq(nodeOperators[2].name, NODE_OPERATOR_THREE_NAME);
  }

  function test_DoesNotIncludeRemovedNodeOperators() public {
    changePrank(ADMIN);
    uint32[] memory nodeOperatorsToRemove = new uint32[](1);
    nodeOperatorsToRemove[0] = 2;
    s_CapabilitiesRegistry.removeNodeOperators(nodeOperatorsToRemove);

    CapabilitiesRegistry.NodeOperator[] memory nodeOperators = s_CapabilitiesRegistry.getNodeOperators();
    assertEq(nodeOperators.length, 2);

    assertEq(nodeOperators[0].admin, NODE_OPERATOR_ONE_ADMIN);
    assertEq(nodeOperators[0].name, NODE_OPERATOR_ONE_NAME);
    assertEq(nodeOperators[1].admin, NODE_OPERATOR_THREE);
    assertEq(nodeOperators[1].name, NODE_OPERATOR_THREE_NAME);
  }
}
