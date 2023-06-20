// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "forge-std/Test.sol";
import {AutomationForwarder} from "../../AutomationForwarder.sol";
import {IAutomationRegistryConsumer} from "../../interfaces/IAutomationRegistryConsumer.sol";
import {UpkeepMock} from "../../mocks/UpkeepMock.sol";
import {MockKeeperRegistry2_1} from "../../mocks/MockKeeperRegistry2_1.sol";

// in contracts directory, run
// forge test --match-path src/v0.8/dev/automation/2_1/test/AutomationForwarder/AutomationForwarder.t.sol

contract AutomationForwarderSetUp is Test {
  AutomationForwarder internal forwarder;
  IAutomationRegistryConsumer internal default_registry;
  UpkeepMock internal default_target;
  address internal OWNER;
  address internal constant STRANGER = address(999);
  uint256 constant GAS = 1e18;

  function setUp() public {
    default_registry = IAutomationRegistryConsumer(new MockKeeperRegistry2_1());
    default_target = new UpkeepMock();
    default_target.setCanCheck(true);
    default_target.setCanPerform(true);
    vm.startPrank(address(default_registry));
    forwarder = new AutomationForwarder(1, address(default_target));
    OWNER = address(default_registry);
  }

  function getSelector(string memory _func, bytes memory myData) public pure returns (bytes memory) {
    bytes4 functionSignature = bytes4(keccak256(bytes(_func)));
    return abi.encodeWithSelector(functionSignature, myData);
  }
}

contract AutomationForwarderTest_forward is AutomationForwarderSetUp {
  function testBasicSuccess() public {
    bytes memory selector = getSelector("performUpkeep(bytes)", "performDataHere");
    bool val = forwarder.forward(GAS, selector);
    assertEq(val, true);
    bool performed = default_target.performed();
    assertEq(performed, true);
  }

  function testWrongFunctionSelectorSuccess() public {
    bytes memory selector = getSelector("performUpkeep(bytes calldata data)", "");
    bool val = forwarder.forward(GAS, selector);
    assertFalse(val);
    assertFalse(default_target.performed());
  }

  function testNotAuthorizedReverts() public {
    bytes memory selector = getSelector("performUpkeep(bytes)", "");
    changePrank(STRANGER);
    vm.expectRevert(AutomationForwarder.NotAuthorized.selector);
    bool val = forwarder.forward(GAS, selector);
    assertFalse(val);
    assertFalse(default_target.performed());
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
