// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {ILiquidityContainer} from "../../interfaces/ILiquidityContainer.sol";

import {LiquidityManager} from "../../LiquidityManager.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract LiquidityManagerHelper is LiquidityManager {
  constructor(
    IERC20 token,
    uint64 localChainSelector,
    ILiquidityContainer localLiquidityContainer,
    uint256 targetTokens,
    address finance
  ) LiquidityManager(token, localChainSelector, localLiquidityContainer, targetTokens, finance) {}

  function report(bytes calldata rep, uint64 ocrSeqNum) external {
    _report(rep, ocrSeqNum);
  }
}
