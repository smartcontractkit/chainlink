// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./iAggregator.sol";
import "./iAggregatorV3.sol";

interface iAggregatorV2V3 is iAggregator, iAggregatorV3 {}
