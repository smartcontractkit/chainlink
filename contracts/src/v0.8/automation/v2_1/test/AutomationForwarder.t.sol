// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {IAutomationRegistryConsumer} from "../../interfaces/IAutomationRegistryConsumer.sol";
import {IAutomationForwarder} from "../../interfaces/IAutomationForwarder.sol";
import {AutomationForwarder} from "../AutomationForwarder.sol";
import {AutomationForwarderLogic} from "../AutomationForwarderLogic.sol";
import {MockKeeperRegistry2_1} from "../../mocks/MockKeeperRegistry2_1.sol";
import {UpkeepCounter} from "../../testhelpers/UpkeepCounter.sol";
import {BaseTest} from "./BaseTest.t.sol";

// in contracts directory, run
// forge test --match-path src/v0.8/automation/v2_1/test/AutomationForwarder.t.sol

contract AutomationForwarderSetUp is BaseTest {
  IAutomationForwarder internal forwarder;
  AutomationForwarderLogic internal logicContract;
  IAutomationRegistryConsumer internal default_registry;
  UpkeepCounter internal default_target;
  uint256 constant GAS = 1e18;

  function setUp() public override {
    // BaseTest.setUp() not called since we want calls to iniatiate from default_registry, not from some predefined owner
    default_registry = IAutomationRegistryConsumer(new MockKeeperRegistry2_1());
    default_target = new UpkeepCounter(10000, 1);
    vm.startPrank(address(default_registry));
    logicContract = new AutomationForwarderLogic();
    forwarder = IAutomationForwarder(
      address(new AutomationForwarder(address(default_target), address(default_registry), address(logicContract)))
    );
    // OWNER not necessary?
    OWNER = address(default_registry);
  }

  function getSelector(string memory _func, bytes memory myData) public pure returns (bytes memory) {
    bytes4 functionSignature = bytes4(keccak256(bytes(_func)));
    return abi.encodeWithSelector(functionSignature, myData);
  }
}

contract AutomationForwarder_forward is AutomationForwarderSetUp {
  function testBasicSuccess() public {
    uint256 prevCount = default_target.counter();
    bytes memory selector = getSelector("performUpkeep(bytes)", "performDataHere");
    (bool val, uint256 gasUsed) = forwarder.forward(GAS, selector);
    assertEq(val, true);
    assertGt(gasUsed, 0);
    uint256 newCount = default_target.counter();
    assertEq(newCount, prevCount + 1);
  }

  function testWrongFunctionSelectorSuccess() public {
    uint256 prevCount = default_target.counter();
    bytes memory selector = getSelector("performUpkeep(bytes calldata data)", "");
    (bool val, uint256 gasUsed) = forwarder.forward(GAS, selector);
    assertFalse(val);
    assertGt(gasUsed, 0);
    uint256 newCount = default_target.counter();
    assertEq(newCount, prevCount);
  }

  function testNotAuthorizedReverts() public {
    uint256 prevCount = default_target.counter();
    bytes memory selector = getSelector("performUpkeep(bytes)", "");
    changePrank(STRANGER);
    vm.expectRevert();
    (bool val, uint256 gasUsed) = forwarder.forward(GAS, selector);
    assertFalse(val);
    assertEq(gasUsed, 0);
    uint256 newCount = default_target.counter();
    assertEq(newCount, prevCount);
  }
}

contract AutomationForwarder_updateRegistry is AutomationForwarderSetUp {
  function testBasicSuccess() public {
    address newRegistry = address(1);
    forwarder.updateRegistry(address(newRegistry));
    IAutomationRegistryConsumer newReg = forwarder.getRegistry();
    assertEq(address(newReg), newRegistry);
  }

  function testNotFromRegistryNotAuthorizedReverts() public {
    address newRegistry = address(1);
    changePrank(STRANGER);
    vm.expectRevert();
    forwarder.updateRegistry(address(newRegistry));
  }
}
