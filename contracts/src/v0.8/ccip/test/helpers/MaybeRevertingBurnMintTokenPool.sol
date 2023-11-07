// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";
import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

contract MaybeRevertingBurnMintTokenPool is BurnMintTokenPool {
  bytes public s_revertReason = "";
  bytes public s_sourceTokenData = "";

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address armProxy
  ) BurnMintTokenPool(token, allowlist, armProxy) {}

  function setShouldRevert(bytes calldata revertReason) external {
    s_revertReason = revertReason;
  }

  function setSourceTokenData(bytes calldata sourceTokenData) external {
    s_sourceTokenData = sourceTokenData;
  }

  function lockOrBurn(
    address originalSender,
    bytes calldata,
    uint256 amount,
    uint64,
    bytes calldata
  ) external virtual override onlyOnRamp checkAllowList(originalSender) whenHealthy returns (bytes memory) {
    bytes memory revertReason = s_revertReason;
    if (revertReason.length != 0) {
      assembly {
        revert(add(32, revertReason), mload(revertReason))
      }
    }
    _consumeOnRampRateLimit(amount);
    IBurnMintERC20(address(i_token)).burn(amount);
    emit Burned(msg.sender, amount);
    return s_sourceTokenData;
  }

  /// @notice Reverts depending on the value of `s_revertReason`
  function releaseOrMint(
    bytes memory,
    address receiver,
    uint256 amount,
    uint64,
    bytes memory
  ) external virtual override whenHealthy onlyOffRamp {
    bytes memory revertReason = s_revertReason;
    if (revertReason.length != 0) {
      assembly {
        revert(add(32, revertReason), mload(revertReason))
      }
    }
    _consumeOffRampRateLimit(amount);
    IBurnMintERC20(address(i_token)).mint(receiver, amount);
    emit Minted(msg.sender, receiver, amount);
  }
}
