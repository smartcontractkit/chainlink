// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";

import {ERC165CheckerReverting} from "../../../libraries/ERC165CheckerReverting.sol";
import {ERC165CheckerRevertingSetup} from "./ERC165CheckerRevertingSetup.t.sol";

contract ERC165CheckerReverting_supportsInterfaceReverting is ERC165CheckerRevertingSetup {
  using ERC165CheckerReverting for address;

  function test__supportsInterfaceReverting() public view {
    assertTrue(s_receiver._supportsInterfaceReverting(type(IAny2EVMMessageReceiver).interfaceId));
  }

  // Reverts

  function test__supportsInterfaceReverting_RevertWhen_NotEnoughGasForSupportsInterface() public {
    vm.expectRevert(NOT_ENOUGH_GAS_SIG);

    // Library calls cannot be called with gas limit overrides, so a public function must be exposed
    // instead which can proxy the call to the library.

    // The gas limit was chosen so that after overhead, <30k would remain to trigger the error.
    this.invokeERC165Checker{gas: 35_000}();
  }

  // Meant to test the call with a manual gas limit override
  function invokeERC165Checker() external view {
    s_receiver._supportsInterfaceReverting(EXAMPLE_INTERFACE_ID);
  }
}
