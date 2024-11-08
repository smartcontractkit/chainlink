// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../../keystone/interfaces/ICapabilityConfiguration.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";
import {CCIPHome} from "../../capability/CCIPHome.sol";
import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_constructor is CCIPHomeTestSetup {
  function test_constructor_success() public {
    CCIPHome ccipHome = new CCIPHome(CAPABILITIES_REGISTRY);

    assertEq(address(ccipHome.getCapabilityRegistry()), CAPABILITIES_REGISTRY);
  }

  function test_supportsInterface_success() public view {
    assertTrue(s_ccipHome.supportsInterface(type(IERC165).interfaceId));
    assertTrue(s_ccipHome.supportsInterface(type(ICapabilityConfiguration).interfaceId));
  }

  function test_getCapabilityConfiguration_success() public view {
    bytes memory config = s_ccipHome.getCapabilityConfiguration(DEFAULT_DON_ID);
    assertEq(config.length, 0);
  }

  function test_constructor_CapabilitiesRegistryAddressZero_reverts() public {
    vm.expectRevert(CCIPHome.ZeroAddressNotAllowed.selector);
    new CCIPHome(address(0));
  }
}
