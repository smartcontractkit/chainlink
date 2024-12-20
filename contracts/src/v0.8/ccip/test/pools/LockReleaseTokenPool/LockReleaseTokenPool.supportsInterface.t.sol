// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IPoolV1} from "../../../interfaces/IPool.sol";

import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

contract LockReleaseTokenPool_supportsInterface is LockReleaseTokenPoolSetup {
  function test_SupportsInterface() public view {
    assertTrue(s_lockReleaseTokenPool.supportsInterface(type(IPoolV1).interfaceId));
    assertTrue(s_lockReleaseTokenPool.supportsInterface(type(IERC165).interfaceId));
  }
}
