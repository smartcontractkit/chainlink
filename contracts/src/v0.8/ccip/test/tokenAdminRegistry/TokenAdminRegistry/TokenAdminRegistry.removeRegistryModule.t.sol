// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistrySetup.t.sol";

contract TokenAdminRegistry_removeRegistryModule is TokenAdminRegistrySetup {
  function test_removeRegistryModule() public {
    address newModule = makeAddr("newModule");

    s_tokenAdminRegistry.addRegistryModule(newModule);

    assertTrue(s_tokenAdminRegistry.isRegistryModule(newModule));

    vm.expectEmit();
    emit TokenAdminRegistry.RegistryModuleRemoved(newModule);

    s_tokenAdminRegistry.removeRegistryModule(newModule);

    assertFalse(s_tokenAdminRegistry.isRegistryModule(newModule));

    // Assert the event is not emitted if the module is already removed.
    vm.recordLogs();
    s_tokenAdminRegistry.removeRegistryModule(newModule);

    vm.assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_RevertWhen_removeRegistryModule_OnlyOwner() public {
    address newModule = makeAddr("newModule");
    vm.stopPrank();

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_tokenAdminRegistry.removeRegistryModule(newModule);
  }
}
