// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import './UpkeepBase.sol';
import './UpkeepInterface.sol';

abstract contract UpkeepCompatible is UpkeepBase, UpkeepInterface {}
