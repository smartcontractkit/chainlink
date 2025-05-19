// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ICapabilityConfiguration} from "../../../../keystone/interfaces/ICapabilityConfiguration.sol";

import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";
import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_supportsInterface is CCIPHomeTestSetup {
  function test_supportsInterface() public view {
    assertTrue(s_ccipHome.supportsInterface(type(IERC165).interfaceId));
    assertTrue(s_ccipHome.supportsInterface(type(ICapabilityConfiguration).interfaceId));
  }
}
