// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPool_constructor is TokenPoolSetup {
  function test_immutableFields_Success() public view {
    assertEq(address(s_token), address(s_tokenPool.getToken()));
    assertEq(address(s_mockRMN), s_tokenPool.getRmnProxy());
    assertEq(false, s_tokenPool.getAllowListEnabled());
    assertEq(address(s_sourceRouter), s_tokenPool.getRouter());
  }

  // Reverts
  function test_ZeroAddressNotAllowed_Revert() public {
    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);

    s_tokenPool = new TokenPoolHelper(IERC20(address(0)), new address[](0), address(s_mockRMN), address(s_sourceRouter));
  }
}
