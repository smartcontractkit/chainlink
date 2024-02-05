// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {IAutomationForwarder} from "../interfaces/IAutomationForwarder.sol";
import {AutomationForwarder} from "../AutomationForwarder.sol";
import {AutomationForwarderLogic} from "../AutomationForwarderLogic.sol";
import "forge-std/Test.sol";

// in contracts directory, run
// forge test --match-path src/v0.8/automation/test/AutomationForwarder.t.sol

contract Target {
  function handler() external pure {}

  function handlerRevert() external pure {
    revert("revert");
  }
}

contract AutomationForwarderTestSetUp is Test {
  address internal constant REGISTRY = 0x3e19ef5Aaa2606655f5A677A97E085cf3811067c;
  address internal constant STRANGER = 0x618fae5d04963B2CEf533F247Eb2C46Bf1801D3b;

  IAutomationForwarder internal forwarder;
  address internal TARGET;

  function setUp() public {
    TARGET = address(new Target());
    AutomationForwarderLogic logicContract = new AutomationForwarderLogic();
    forwarder = IAutomationForwarder(address(new AutomationForwarder(TARGET, REGISTRY, address(logicContract))));
  }
}

contract AutomationForwarderTest_constructor is AutomationForwarderTestSetUp {
  function testInitialValues() external {
    assertEq(address(forwarder.getRegistry()), REGISTRY);
    assertEq(forwarder.getTarget(), TARGET);
  }

  function testTypeAndVersion() external {
    assertEq(forwarder.typeAndVersion(), "AutomationForwarder 1.0.0");
  }
}

contract AutomationForwarderTest_forward is AutomationForwarderTestSetUp {
  function testOnlyCallableByTheRegistry() external {
    vm.prank(REGISTRY);
    forwarder.forward(100000, abi.encodeWithSelector(Target.handler.selector));
    vm.prank(STRANGER);
    vm.expectRevert();
    forwarder.forward(100000, abi.encodeWithSelector(Target.handler.selector));
  }

  function testReturnsSuccessValueAndGasUsed() external {
    vm.startPrank(REGISTRY);
    (bool success, uint256 gasUsed) = forwarder.forward(100000, abi.encodeWithSelector(Target.handler.selector));
    assertTrue(success);
    assertGt(gasUsed, 0);
    (success, gasUsed) = forwarder.forward(100000, abi.encodeWithSelector(Target.handlerRevert.selector));
    assertFalse(success);
    assertGt(gasUsed, 0);
  }
}

contract AutomationForwarderTest_updateRegistry is AutomationForwarderTestSetUp {
  function testOnlyCallableByTheActiveRegistry() external {
    address newRegistry = address(1);
    vm.startPrank(REGISTRY);
    forwarder.updateRegistry(newRegistry);
    assertEq(address(forwarder.getRegistry()), newRegistry);
    vm.expectRevert();
    forwarder.updateRegistry(REGISTRY);
  }
}
