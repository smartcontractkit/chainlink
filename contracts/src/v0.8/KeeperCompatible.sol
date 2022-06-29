// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./KeeperBase.sol";
import "./interfaces/IKeeperCompatible.sol";

abstract contract KeeperCompatible is KeeperBase, IKeeperCompatible {}
