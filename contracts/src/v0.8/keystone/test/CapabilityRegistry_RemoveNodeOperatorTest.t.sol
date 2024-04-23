// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveNodeOperatorTest is BaseTest {
    event NodeOperatorRemoved(uint256 nodeOperatorId);

    uint256 private constant TEST_NODE_OPERATOR_ID = 0;

    function setUp() public override {
        BaseTest.setUp();
        changePrank(ADMIN);
        s_capabilityRegistry.addNodeOperator(NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);
    }

    function test_RevertWhen_CalledByNonAdminAndNonOwner() public {
        changePrank(STRANGER);
        vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
        s_capabilityRegistry.removeNodeOperator(TEST_NODE_OPERATOR_ID);
    }

    function test_RevertWhen_NodeOperatorDoesNotExist() public {
        changePrank(ADMIN);
        vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NonExistentNodeOperator.selector, 1));
        s_capabilityRegistry.removeNodeOperator(1);
    }

    function test_RemovesNodeOperator() public {
        changePrank(ADMIN);

        vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
        emit NodeOperatorRemoved(TEST_NODE_OPERATOR_ID);
        s_capabilityRegistry.removeNodeOperator(TEST_NODE_OPERATOR_ID);
        CapabilityRegistry.NodeOperator memory nodeOperator = s_capabilityRegistry.getNodeOperator(0);
        assertEq(nodeOperator.admin, address(0));
        assertEq(nodeOperator.name, "");
    }
}
