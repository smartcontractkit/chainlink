// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {USDCTokenPool} from "../../pools/USDC/USDCTokenPool.sol";

contract USDCTokenPoolHelper is USDCTokenPool {
  constructor(
    USDCConfig memory config,
    IBurnMintERC20 token,
    address[] memory allowlist,
    address armProxy,
    uint32 localDomainIdentifier
  ) USDCTokenPool(config, token, allowlist, armProxy, localDomainIdentifier) {}

  function validateMessage(bytes memory usdcMessage, SourceTokenDataPayload memory sourceTokenData) external view {
    return _validateMessage(usdcMessage, sourceTokenData);
  }
}
