// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {USDCTokenPool} from "../../pools/USDC/USDCTokenPool.sol";
import {ITokenMessenger} from "../../pools/USDC/ITokenMessenger.sol";

contract USDCTokenPoolHelper is USDCTokenPool {
  constructor(
    ITokenMessenger tokenMessenger,
    IBurnMintERC20 token,
    address[] memory allowlist,
    address armProxy,
    address router
  ) USDCTokenPool(tokenMessenger, token, allowlist, armProxy, router) {}

  function validateMessage(bytes memory usdcMessage, SourceTokenDataPayload memory sourceTokenData) external view {
    return _validateMessage(usdcMessage, sourceTokenData);
  }
}
