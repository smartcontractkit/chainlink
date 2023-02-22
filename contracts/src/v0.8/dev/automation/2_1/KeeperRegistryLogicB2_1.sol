// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./KeeperRegistryBase2_1.sol";

contract KeeperRegistryLogicB2_1 is KeeperRegistryBase2_1 {
  /**
   * @dev see KeeperRegistry master contract for constructor description
   */
  constructor(
    Mode mode,
    address link,
    address linkNativeFeed,
    address fastGasFeed,
    address forwarderFactory
  ) KeeperRegistryBase2_1(mode, link, linkNativeFeed, fastGasFeed, forwarderFactory) {}
}
