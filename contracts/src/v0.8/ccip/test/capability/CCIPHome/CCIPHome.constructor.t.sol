// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_constructor is CCIPHomeTestSetup {
  function test_constructor() public {
    CCIPHome ccipHome = new CCIPHome(CAPABILITIES_REGISTRY);

    assertEq(address(ccipHome.getCapabilityRegistry()), CAPABILITIES_REGISTRY);
  }

  function test_RevertWhen_constructor_CapabilitiesRegistryAddressZero() public {
    vm.expectRevert(CCIPHome.ZeroAddressNotAllowed.selector);
    new CCIPHome(address(0));
  }
}
