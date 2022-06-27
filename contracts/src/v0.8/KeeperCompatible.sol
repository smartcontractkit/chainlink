// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./KeeperBase.sol";
import "./interfaces/iKeeperCompatible.sol";

abstract contract KeeperCompatible is KeeperBase, iKeeperCompatible {}
