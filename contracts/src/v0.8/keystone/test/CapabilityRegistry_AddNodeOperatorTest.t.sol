// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddNodeOperatorTest is BaseTest {
    event NodeOperatorAdded(uint256 nodeOperatorId, address indexed admin, string name);

    function test_RevertWhen_CalledByNonAdmin() public {
        changePrank(STRANGER);
        vm.expectRevert("Only callable by owner");
        s_capabilityRegistry.addNodeOperator(NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);
    }

    function test_RevertWhen_NodeOperatorAdminAddressZero() public {
        changePrank(ADMIN);
        vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
        s_capabilityRegistry.addNodeOperator(address(0), NODE_OPERATOR_ONE_NAME);
    }

    function test_AddNodeOperator() public {
        changePrank(ADMIN);

        vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
        emit NodeOperatorAdded(0, NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);
        s_capabilityRegistry.addNodeOperator(NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);

        CapabilityRegistry.NodeOperator memory nodeOperator = s_capabilityRegistry.getNodeOperator(0);

        assertEq(nodeOperator.id, 0);
        assertEq(nodeOperator.admin, NODE_OPERATOR_ONE_ADMIN);
        assertEq(nodeOperator.name, NODE_OPERATOR_ONE_NAME);
    }
}
