// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IPoolV1} from "../../../../interfaces/IPool.sol";
import {USDCTokenPoolSetup} from "./USDCTokenPoolSetup.t.sol";

import {IERC165} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

contract USDCTokenPool_supportsInterface is USDCTokenPoolSetup {
  function test_SupportsInterface() public view {
    assertTrue(s_usdcTokenPool.supportsInterface(type(IPoolV1).interfaceId));
    assertTrue(s_usdcTokenPool.supportsInterface(type(IERC165).interfaceId));
  }
}
