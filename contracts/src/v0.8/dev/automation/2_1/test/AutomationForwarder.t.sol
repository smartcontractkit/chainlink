// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import 'forge-std/Test.sol';
import '../AutomationForwarder.sol';
import {IAutomationRegistryConsumer} from "../interfaces/IAutomationRegistryConsumer.sol";
import '../mocks/MockCustomerTarget.sol';
import '../KeeperRegistryBase2_1.sol';
import '../mocks/MockKeeperRegistry2_1.sol';

contract AutomationForwarderTest is Test {

    AutomationForwarder public forwarder;
    IAutomationRegistryConsumer public default_registry;
    MockCustomerTarget public default_target;

    uint256 constant GAS = 5_000;

     function setUp() public {
        default_registry = IAutomationRegistryConsumer(new MockKeeperRegistry2_1());
        default_target = new MockCustomerTarget();
        vm.prank(address(default_registry));
        forwarder = new AutomationForwarder(1, address(default_target));
    }

    function getSelector(string memory _func) public pure returns (bytes memory) {
        bytes4 selector = bytes4(keccak256(bytes(_func)));
        bytes memory selectorBytes = abi.encodePacked(selector);
        return selectorBytes;
    }
  

    function test_forward() public {
        bytes memory selector = getSelector('performUpkeep()');
        vm.prank(address(default_registry));
        bool val = forwarder.forward(GAS, selector);
        assertEq(val, true);
    }

    function test_forward_WrongFunctionSelector() public {
        bytes memory selector = getSelector('performUpkeep(int num)');
        vm.prank(address(default_registry));
        bool val = forwarder.forward(GAS, selector);
        assertFalse(val);
    }

    function test_forward_NotFromRegistry() public {
        bytes memory selector = getSelector('performUpkeep()');
        vm.expectRevert(AutomationForwarder.NotAuthorized.selector);
        bool val = forwarder.forward(GAS, selector);
    }

    function test_updateRegistry() public {
        address newRegistry = address(1);
        vm.prank(address(default_registry));
        forwarder.updateRegistry(address(newRegistry));
        IAutomationRegistryConsumer newReg = forwarder.getRegistry();
        assertEq(address(newReg), newRegistry);
    }

    function test_updateRegistry_NotFromRegistry() public {
        address newRegistry = address(1);
        vm.expectRevert(AutomationForwarder.NotAuthorized.selector);
        forwarder.updateRegistry(address(newRegistry));
    }

    function test_getters() public {
        bytes memory reg = abi.encodePacked(address(forwarder.getRegistry()));
        assertEq(reg, abi.encodePacked(default_registry), 'getRegistry()');
        bytes memory targ = abi.encodePacked(address(forwarder.getTarget()));
        assertEq(targ, abi.encodePacked(default_target), 'getTarget()');
        uint256 id = forwarder.getUpkeepID();
        assertEq(id, 1, 'getUpkeepId()');
    }


}