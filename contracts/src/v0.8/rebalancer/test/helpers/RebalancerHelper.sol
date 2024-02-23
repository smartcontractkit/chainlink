// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {ILiquidityContainer} from "../../interfaces/ILiquidityContainer.sol";

import {Rebalancer} from "../../Rebalancer.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract RebalancerHelper is Rebalancer {
  constructor(
    IERC20 token,
    uint64 localChainSelector,
    ILiquidityContainer localLiquidityContainer
  ) Rebalancer(token, localChainSelector, localLiquidityContainer) {}

  function report(bytes calldata rep, uint64 ocrSeqNum) external {
    _report(rep, ocrSeqNum);
  }
}
