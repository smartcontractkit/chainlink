// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IBundleBaseAggregator} from "./IBundleBaseAggregator.sol";
import {ICommonAggregator} from "./ICommonAggregator.sol";

interface IBundleAggregatorProxy is IBundleBaseAggregator, ICommonAggregator {
  function proposedAggregator() external view returns (address);

  function confirmAggregator(
    address aggregatorAddress
  ) external;

  function aggregator() external view returns (address);
}
