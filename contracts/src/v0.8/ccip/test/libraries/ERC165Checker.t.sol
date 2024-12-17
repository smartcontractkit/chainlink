// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";

import {ERC165Checker} from "../../libraries/ERC165Checker.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";

import {Test} from "forge-std/Test.sol";

contract ERC165CheckerTest is Test {
  using ERC165Checker for address;

  address internal s_receiver;

  bytes4 public constant EXAMPLE_INTERFACE_ID = 0xdeadbeef;

  constructor() {
    s_receiver = address(new MaybeRevertMessageReceiver(false));
  }

  function test__supportsInterface() public view {
    assertTrue(s_receiver._supportsInterface(type(IAny2EVMMessageReceiver).interfaceId));
  }

  function test__getSupportedInterfaces() public view {
    bytes4[] memory interfaceIds = new bytes4[](1);

    interfaceIds[0] = type(IAny2EVMMessageReceiver).interfaceId;

    bool[] memory supportedIds = s_receiver._getSupportedInterfaces(interfaceIds);
    assertTrue(supportedIds[0]);
    assertEq(interfaceIds.length, supportedIds.length);
  }

  function test__supportsAllInterfaces() public view {
    bytes4[] memory interfaceIds = new bytes4[](1);

    interfaceIds[0] = type(IAny2EVMMessageReceiver).interfaceId;

    assertTrue(s_receiver._supportsAllInterfaces(interfaceIds));
  }

  function test__supportsAllInterfaces_notAllSupported() public view {
    bytes4[] memory interfaceIds = new bytes4[](2);

    interfaceIds[0] = type(IAny2EVMMessageReceiver).interfaceId;
    interfaceIds[1] = EXAMPLE_INTERFACE_ID;

    assertFalse(s_receiver._supportsAllInterfaces(interfaceIds));
  }

  function test__supportsAllInterfaces_notSupportsERC165() public view {
    bytes4[] memory interfaceIds = new bytes4[](2);

    interfaceIds[0] = type(IAny2EVMMessageReceiver).interfaceId;
    interfaceIds[1] = EXAMPLE_INTERFACE_ID;

    // An address that does not support ERC165
    address randomAddress = address(0xdead);

    assertFalse(randomAddress._supportsAllInterfaces(interfaceIds));
  }

  function test_supportsInterface_RevertWhen_NotEnoughGasForSupportsInterface() public {
    vm.expectRevert(ERC165Checker.NotEnoughGasForSupportsInterfaceCall.selector);

    // Library calls cannot be called with gas limit overrides, so a public function must be exposed
    // instead which can proxy the call to the library.

    // The gas limit was chosen so that after overhead, <30k would remain to trigger the error.
    this.invokeERC165Checker{gas: 35_000}();
  }

  // Meant to test the call with a manual gas limit override
  function invokeERC165Checker() external view {
    s_receiver._supportsInterface(EXAMPLE_INTERFACE_ID);
  }
}
