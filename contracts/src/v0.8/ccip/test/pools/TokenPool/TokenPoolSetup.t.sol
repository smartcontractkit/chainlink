// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {BurnMintERC677} from "../../../../shared/token/ERC677/BurnMintERC677.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";
import {RouterSetup} from "../../router/Router/RouterSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPoolSetup is RouterSetup {
  IERC20 internal s_token;
  TokenPoolHelper internal s_tokenPool;

  function setUp() public virtual override {
    RouterSetup.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);

    s_tokenPool = new TokenPoolHelper(s_token, new address[](0), address(s_mockRMN), address(s_sourceRouter));
  }
}
