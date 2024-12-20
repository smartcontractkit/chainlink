// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPoolWithAllowListSetup is TokenPoolSetup {
  address[] internal s_allowedSenders;

  function setUp() public virtual override {
    TokenPoolSetup.setUp();

    s_allowedSenders.push(STRANGER);
    s_allowedSenders.push(OWNER);

    s_tokenPool = new TokenPoolHelper(
      s_token, DEFAULT_TOKEN_DECIMALS, s_allowedSenders, address(s_mockRMNRemote), address(s_sourceRouter)
    );
  }
}
