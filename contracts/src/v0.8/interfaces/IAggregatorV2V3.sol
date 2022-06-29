// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IAggregator.sol";
import "./IAggregatorV3.sol";

interface IAggregatorV2V3 is IAggregator, IAggregatorV3 {}
