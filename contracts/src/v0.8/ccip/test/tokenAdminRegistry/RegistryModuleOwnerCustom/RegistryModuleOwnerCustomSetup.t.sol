// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {BurnMintERC677Helper} from "../../helpers/BurnMintERC677Helper.sol";

import {Test} from "forge-std/Test.sol";

contract RegistryModuleOwnerCustomSetup is Test {
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;

  RegistryModuleOwnerCustom internal s_registryModuleOwnerCustom;
  TokenAdminRegistry internal s_tokenAdminRegistry;
  address internal s_token;

  function setUp() public virtual {
    vm.startPrank(OWNER);

    s_tokenAdminRegistry = new TokenAdminRegistry();
    s_token = address(new BurnMintERC677Helper("Test", "TST"));
    s_registryModuleOwnerCustom = new RegistryModuleOwnerCustom(address(s_tokenAdminRegistry));
    s_tokenAdminRegistry.addRegistryModule(address(s_registryModuleOwnerCustom));
  }
}
