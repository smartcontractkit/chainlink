// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";

import {ERC165CheckerReverting} from "../../libraries/ERC165CheckerReverting.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";

import {Test} from "forge-std/Test.sol";

contract ERC165CheckerReverting_supportsInterfaceReverting is Test {
  using ERC165CheckerReverting for address;

  address internal s_receiver;

  bytes4 internal constant EXAMPLE_INTERFACE_ID = 0xdeadbeef;

  error InsufficientGasForStaticCall();

  constructor() {
    s_receiver = address(new MaybeRevertMessageReceiver(false));
  }

  function test__supportsInterfaceReverting() public view {
    assertTrue(s_receiver._supportsInterfaceReverting(type(IAny2EVMMessageReceiver).interfaceId));
  }

  // Reverts

  function test__supportsInterfaceReverting_RevertWhen_NotEnoughGasForSupportsInterface() public {
    vm.expectRevert(InsufficientGasForStaticCall.selector);

    // Library calls cannot be called with gas limit overrides, so a public function must be exposed
    // instead which can proxy the call to the library.

    // The gas limit was chosen so that after overhead, <31k would remain to trigger the error.
    this.invokeERC165Checker{gas: 33_000}();
  }

  // Meant to test the call with a manual gas limit override
  function invokeERC165Checker() external view {
    s_receiver._supportsInterfaceReverting(EXAMPLE_INTERFACE_ID);
  }
}
