// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";

import {ERC165Checker} from "../../libraries/ERC165Checker.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";

import {Test} from "forge-std/Test.sol";

contract ERC165CheckerTest is Test {
  using ERC165Checker for address;

  MaybeRevertMessageReceiver internal s_receiver;

  bytes4 public constant EXAMPLE_INTERFACE_ID = 0xdeadbeef;

  constructor() {
    s_receiver = new MaybeRevertMessageReceiver(false);
  }

  function test_supportsInterface() public view {
    assertTrue(s_receiver.supportsInterface(type(IAny2EVMMessageReceiver).interfaceId));
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
    address(s_receiver)._supportsInterface(EXAMPLE_INTERFACE_ID);
  }
}
