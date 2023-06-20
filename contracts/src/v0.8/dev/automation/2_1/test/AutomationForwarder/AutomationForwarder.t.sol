// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import 'forge-std/Test.sol';
import {AutomationForwarder} from '../../AutomationForwarder.sol';
import {IAutomationRegistryConsumer} from "../../interfaces/IAutomationRegistryConsumer.sol";
import {MockCustomerTarget} from '../../mocks/MockCustomerTarget.sol';
import {KeeperRegistryBase2_1} from '../../KeeperRegistryBase2_1.sol';
import {MockKeeperRegistry2_1} from '../../mocks/MockKeeperRegistry2_1.sol';

contract AutomationForwarderSetUp is Test {

    AutomationForwarder internal forwarder;
    IAutomationRegistryConsumer internal default_registry;
    MockCustomerTarget internal default_target;
    address internal OWNER;
    address internal constant STRANGER = address(999);

    uint256 constant GAS = 5_000;

     function setUp() public {
        default_registry = IAutomationRegistryConsumer(new MockKeeperRegistry2_1());
        default_target = new MockCustomerTarget();
        vm.startPrank(address(default_registry));
        forwarder = new AutomationForwarder(1, address(default_target));
        OWNER = address(default_registry);
    }

    function getSelector(string memory _func) public pure returns (bytes memory) {
        bytes4 selector = bytes4(keccak256(bytes(_func)));
        bytes memory selectorBytes = abi.encodePacked(selector);
        return selectorBytes;
    }

}

contract AutomationForwarderTest_forward is AutomationForwarderSetUp {
    function testBasicSuccess() public {
        bytes memory selector = getSelector('performUpkeep()');
        bool val = forwarder.forward(GAS, selector);
        assertEq(val, true);
    }

    function testWrongFunctionSelectorSuccess() public {
        bytes memory selector = getSelector('performUpkeep(int num)');
        bool val = forwarder.forward(GAS, selector);
        assertFalse(val);
    }

    function testNotFromRegistryNotAuthorizedReverts() public {
        bytes memory selector = getSelector('performUpkeep()');
        changePrank(STRANGER);
        vm.expectRevert(AutomationForwarder.NotAuthorized.selector);
        bool val = forwarder.forward(GAS, selector);
    }
}

contract AutomationForwarderTest_updateRegistry is AutomationForwarderSetUp {
    function testBasicSuccess() public {
        address newRegistry = address(1);
        forwarder.updateRegistry(address(newRegistry));
        IAutomationRegistryConsumer newReg = forwarder.getRegistry();
        assertEq(address(newReg), newRegistry);
    }

    function testNotFromRegistryNotAuthorizedReverts() public {
        address newRegistry = address(1);
        changePrank(STRANGER);
        vm.expectRevert(AutomationForwarder.NotAuthorized.selector);
        forwarder.updateRegistry(address(newRegistry));
    }

}
