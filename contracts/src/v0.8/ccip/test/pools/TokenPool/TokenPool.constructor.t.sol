// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from
  "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";

contract TokenPool_constructor is TokenPoolSetup {
  function test_constructor() public view {
    assertEq(address(s_token), address(s_tokenPool.getToken()));
    assertEq(address(s_mockRMNRemote), s_tokenPool.getRmnProxy());
    assertFalse(s_tokenPool.getAllowListEnabled());
    assertEq(address(s_sourceRouter), s_tokenPool.getRouter());
    assertEq(DEFAULT_TOKEN_DECIMALS, s_tokenPool.getTokenDecimals());
  }

  function test_constructor_DecimalCallFails() public {
    uint8 decimals = 255;

    vm.mockCallRevert(address(s_token), abi.encodeWithSelector(IERC20Metadata.decimals.selector), "decimals fails");

    s_tokenPool =
      new TokenPoolHelper(s_token, decimals, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter));

    assertEq(s_tokenPool.getTokenDecimals(), decimals);
  }

  // Reverts

  function test_RevertWhen_constructorWhen_ZeroAddressNotAllowed() public {
    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);

    s_tokenPool = new TokenPoolHelper(
      IERC20(address(0)), DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );
  }

  function test_RevertWhen_constructorWhen_InvalidDecimalArgs() public {
    uint8 invalidDecimals = DEFAULT_TOKEN_DECIMALS + 1;

    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.InvalidDecimalArgs.selector, invalidDecimals, DEFAULT_TOKEN_DECIMALS)
    );

    s_tokenPool =
      new TokenPoolHelper(s_token, invalidDecimals, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter));
  }
}
