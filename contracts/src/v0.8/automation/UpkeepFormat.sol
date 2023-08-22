// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

/**
 * @dev this struct is only maintained for backwards compatibility with MigratableKeeperRegistryInterface
 * it should be deprecated in the future in favor of MigratableKeeperRegistryInterfaceV2
 */
enum UpkeepFormat {
  V1,
  V2,
  V3
}
