// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../../automation/interfaces/KeeperCompatibleInterface.sol";
import "../../automation/interfaces/v1_2/KeeperRegistryInterface1_2.sol";
import "../../shared/access/ConfirmedOwner.sol";

error NoKeeperNodes();
error InsufficientInterval();

/**
 * @notice A canary upkeep which requires a different keeper to service its upkeep at an interval. This makes sure that
 * all keepers are in a healthy state.
 */
contract CanaryUpkeep1_2 is KeeperCompatibleInterface, ConfirmedOwner {
  uint256 private s_keeperIndex;
  uint256 private s_interval;
  uint256 private s_timestamp;
  KeeperRegistryExecutableInterface private immutable i_keeperRegistry;

  /**
   * @param keeperRegistry address of a keeper registry
   */
  constructor(KeeperRegistryExecutableInterface keeperRegistry, uint256 interval) ConfirmedOwner(msg.sender) {
    i_keeperRegistry = keeperRegistry;
    s_timestamp = block.timestamp;
    s_interval = interval;
    s_keeperIndex = 0;
  }

  /**
   * @return the current keeper index
   */
  function getKeeperIndex() external view returns (uint256) {
    return s_keeperIndex;
  }

  /**
   * @return the current timestamp
   */
  function getTimestamp() external view returns (uint256) {
    return s_timestamp;
  }

  /**
   * @return the current interval
   */
  function getInterval() external view returns (uint256) {
    return s_interval;
  }

  /**
   * @return the keeper registry
   */
  function getKeeperRegistry() external view returns (KeeperRegistryExecutableInterface) {
    return i_keeperRegistry;
  }

  /**
   * @notice updates the interval
   * @param interval the new interval
   */
  function setInterval(uint256 interval) external onlyOwner {
    s_interval = interval;
  }

  /**
   * @notice returns true if keeper array is not empty and sufficient time has passed
   */
  function checkUpkeep(bytes calldata /* checkData */) external view override returns (bool, bytes memory) {
    bool upkeepNeeded = block.timestamp >= s_interval + s_timestamp;
    return (upkeepNeeded, bytes(""));
  }

  /**
   * @notice checks keepers array limit, timestamp limit, and requires transaction origin must be the anticipated keeper.
   * If all checks pass, update the keeper index and timestamp. Otherwise, revert this transaction.
   */
  function performUpkeep(bytes calldata /* performData */) external override {
    (State memory _s, Config memory _c, address[] memory keepers) = i_keeperRegistry.getState();
    if (keepers.length == 0) {
      revert NoKeeperNodes();
    }
    if (block.timestamp < s_interval + s_timestamp) {
      revert InsufficientInterval();
    }
    // if keepers array is shortened, this statement will make sure keeper index is always valid
    if (s_keeperIndex >= keepers.length) {
      s_keeperIndex = 0;
    }

    require(tx.origin == keepers[s_keeperIndex], "transaction origin is not the anticipated keeper.");
    s_keeperIndex = (s_keeperIndex + 1) % keepers.length;
    s_timestamp = block.timestamp;
  }
}
