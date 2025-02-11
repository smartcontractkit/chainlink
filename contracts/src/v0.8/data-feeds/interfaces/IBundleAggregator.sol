// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IBundleBaseAggregator} from "./IBundleBaseAggregator.sol";
import {ICommonAggregator} from "./ICommonAggregator.sol";

interface IBundleAggregator is IBundleBaseAggregator, ICommonAggregator {}
