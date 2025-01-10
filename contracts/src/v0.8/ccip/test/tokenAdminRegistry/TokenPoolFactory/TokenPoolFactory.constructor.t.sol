// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITokenAdminRegistry} from "../../../interfaces/ITokenAdminRegistry.sol";

import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenPoolFactory} from "../../../tokenAdminRegistry/TokenPoolFactory/TokenPoolFactory.sol";

import {Create2} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";
import {TokenPoolFactorySetup} from "./TokenPoolFactorySetup.t.sol";

contract TokenPoolFactory_constructor is TokenPoolFactorySetup {
  using Create2 for bytes32;

  function test_RevertWhen_constructor() public {
    // Revert cause the tokenAdminRegistry is address(0)
    vm.expectRevert(TokenPoolFactory.InvalidZeroAddress.selector);
    new TokenPoolFactory(ITokenAdminRegistry(address(0)), RegistryModuleOwnerCustom(address(0)), address(0), address(0));

    new TokenPoolFactory(
      ITokenAdminRegistry(address(0xdeadbeef)),
      RegistryModuleOwnerCustom(address(0xdeadbeef)),
      address(0xdeadbeef),
      address(0xdeadbeef)
    );
  }
}
