// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_UpdateNodeOperatorTest is BaseTest {
    event NodeOperatorUpdated(uint256 nodeOperatorId, address indexed admin, string name);

    uint256 private constant TEST_NODE_OPERATOR_ID = 0;
    address private constant NEW_NODE_OPERATOR_ADMIN = address(3);
    string private constant NEW_NODE_OPERATOR_NAME = "new-node-operator";

    function setUp() public override {
        BaseTest.setUp();
        changePrank(ADMIN);
        s_capabilityRegistry.addNodeOperator(NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);
    }

    function test_RevertWhen_CalledByNonAdminAndNonOwner() public {
        changePrank(STRANGER);
        vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
        s_capabilityRegistry.updateNodeOperator(TEST_NODE_OPERATOR_ID, NEW_NODE_OPERATOR_ADMIN, NEW_NODE_OPERATOR_NAME);
    }

    function test_RevertWhen_NodeOperatorAdminIsZeroAddress() public {
        changePrank(ADMIN);
        vm.expectRevert(CapabilityRegistry.InvalidNodeOperatorAdmin.selector);
        s_capabilityRegistry.updateNodeOperator(TEST_NODE_OPERATOR_ID, address(0), NEW_NODE_OPERATOR_NAME);
    }

    function test_RevertWhen_NodeOperatorIsNotUpdated() public {
        changePrank(ADMIN);
        vm.expectRevert(CapabilityRegistry.InvalidNodeOperatorUpdate.selector);
        s_capabilityRegistry.updateNodeOperator(TEST_NODE_OPERATOR_ID, NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);
    }

    function test_RevertWhen_NodeOperatorDoesNotExist() public {
        changePrank(ADMIN);
        vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NonExistentNodeOperator.selector, 1));
        s_capabilityRegistry.updateNodeOperator(1, NEW_NODE_OPERATOR_ADMIN, NEW_NODE_OPERATOR_NAME);
    }

    function test_UpdatesNodeOperator() public {
        changePrank(ADMIN);
        vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
        emit NodeOperatorUpdated(TEST_NODE_OPERATOR_ID, NEW_NODE_OPERATOR_ADMIN, NEW_NODE_OPERATOR_NAME);
        s_capabilityRegistry.updateNodeOperator(TEST_NODE_OPERATOR_ID, NEW_NODE_OPERATOR_ADMIN, NEW_NODE_OPERATOR_NAME);
        CapabilityRegistry.NodeOperator memory nodeOperator = s_capabilityRegistry.getNodeOperator(0);
        assertEq(nodeOperator.admin, NEW_NODE_OPERATOR_ADMIN);
        assertEq(nodeOperator.name, NEW_NODE_OPERATOR_NAME);
    }
}
