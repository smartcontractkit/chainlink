// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_parseRemoteDecimals is TokenPoolSetup {
  function test_parseRemoteDecimals() public view {
    uint256 remoteDecimals = 12;
    bytes memory encodedDecimals = abi.encode(remoteDecimals);

    assertEq(s_tokenPool.parseRemoteDecimals(encodedDecimals), remoteDecimals);

    assertEq(s_tokenPool.parseRemoteDecimals(s_tokenPool.encodeLocalDecimals()), s_tokenPool.getTokenDecimals());
  }

  function test_parseRemoteDecimals_NoDecimalsDefaultsToLocalDecimals() public view {
    assertEq(s_tokenPool.parseRemoteDecimals(""), s_tokenPool.getTokenDecimals());
  }

  function test_RevertWhen_parseRemoteDecimalsWhen_InvalidRemoteChainDecimals_DigitTooLarge() public {
    bytes memory encodedDecimals = abi.encode(uint256(256));

    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidRemoteChainDecimals.selector, encodedDecimals));

    s_tokenPool.parseRemoteDecimals(encodedDecimals);
  }

  function test_RevertWhen_parseRemoteDecimalsWhen_InvalidRemoteChainDecimals_WrongType() public {
    bytes memory encodedDecimals = abi.encode(uint256(256), "wrong type");

    vm.expectRevert(abi.encodeWithSelector(TokenPool.InvalidRemoteChainDecimals.selector, encodedDecimals));

    s_tokenPool.parseRemoteDecimals(encodedDecimals);
  }
}
