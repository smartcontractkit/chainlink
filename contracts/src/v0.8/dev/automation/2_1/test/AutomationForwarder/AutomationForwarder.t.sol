// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IAutomationRegistryConsumer} from "../../interfaces/IAutomationRegistryConsumer.sol";
import {AutomationForwarder} from "../../AutomationForwarder.sol";
import {MockKeeperRegistry2_1} from "../../mocks/MockKeeperRegistry2_1.sol";
import {UpkeepCounter} from "../../mocks/UpkeepCounter.sol";
import "forge-std/Test.sol";

// in contracts directory, run
// forge test --match-path src/v0.8/dev/automation/2_1/test/AutomationForwarder/AutomationForwarder.t.sol

contract AutomationForwarderSetUp is Test {
  AutomationForwarder internal forwarder;
  IAutomationRegistryConsumer internal default_registry;
  UpkeepCounter internal default_target;
  address internal OWNER;
  address internal constant STRANGER = address(999);
  uint256 constant GAS = 1e18;

  function setUp() public {
    default_registry = IAutomationRegistryConsumer(new MockKeeperRegistry2_1());
    default_target = new UpkeepCounter(10000, 1);
    vm.startPrank(address(default_registry));
    forwarder = new AutomationForwarder(1, address(default_target), address(default_registry));
    OWNER = address(default_registry);
  }

  function getSelector(string memory _func, bytes memory myData) public pure returns (bytes memory) {
    bytes4 functionSignature = bytes4(keccak256(bytes(_func)));
    return abi.encodeWithSelector(functionSignature, myData);
  }
}

contract AutomationForwarderTest_forward is AutomationForwarderSetUp {
  function testBasicSuccess() public {
    uint256 prevCount = default_target.counter();
    bytes memory selector = getSelector("performUpkeep(bytes)", "performDataHere");
    bool val = forwarder.forward(GAS, selector);
    assertEq(val, true);
    uint256 newCount = default_target.counter();
    assertEq(newCount, prevCount + 1);
  }

  function testWrongFunctionSelectorSuccess() public {
    uint256 prevCount = default_target.counter();
    bytes memory selector = getSelector("performUpkeep(bytes calldata data)", "");
    bool val = forwarder.forward(GAS, selector);
    assertFalse(val);
    uint256 newCount = default_target.counter();
    assertEq(newCount, prevCount);
  }

  function testNotAuthorizedReverts() public {
    uint256 prevCount = default_target.counter();
    bytes memory selector = getSelector("performUpkeep(bytes)", "");
    changePrank(STRANGER);
    vm.expectRevert(AutomationForwarder.NotAuthorized.selector);
    bool val = forwarder.forward(GAS, selector);
    assertFalse(val);
    uint256 newCount = default_target.counter();
    assertEq(newCount, prevCount);
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
