// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {RegistryModuleOwnerCustomSetup} from "./RegistryModuleOwnerCustomSetup.t.sol";

contract RegistryModuleOwnerCustom_constructor is RegistryModuleOwnerCustomSetup {
  function test_constructor_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(RegistryModuleOwnerCustom.AddressZero.selector));

    new RegistryModuleOwnerCustom(address(0));
  }
}
