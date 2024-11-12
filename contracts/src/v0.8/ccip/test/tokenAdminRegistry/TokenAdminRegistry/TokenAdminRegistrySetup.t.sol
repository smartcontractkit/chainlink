// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenSetup} from "../../TokenSetup.t.sol";

contract TokenAdminRegistrySetup is TokenSetup {
  address internal s_registryModule = makeAddr("registryModule");

  function setUp() public virtual override {
    TokenSetup.setUp();

    s_tokenAdminRegistry.addRegistryModule(s_registryModule);
  }
}
