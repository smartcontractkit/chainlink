// SPDX-License-Identifier: MIT
pragma solidity 0.8.13;

import {KeeperRegistryExecutableInterface} from "../KeeperRegistry.sol";
import "../ConfirmedOwner.sol";

/**
 * @notice This contract serves as a wrapper around a keeper registry's checkUpkeep function.
 */
contract KeeperRegistryCheckUpkeepGasUsageWrapper is ConfirmedOwner {
  KeeperRegistryExecutableInterface private immutable i_keeperRegistry;

  /**
   * @param keeperRegistry address of a keeper registry
   */
  constructor(KeeperRegistryExecutableInterface keeperRegistry) ConfirmedOwner(msg.sender) {
    i_keeperRegistry = keeperRegistry;
  }

  /**
   * @return the keeper registry
   */
  function getKeeperRegistry() external view returns (KeeperRegistryExecutableInterface) {
    return i_keeperRegistry;
  }

  /**
   * @notice This function is called by monitoring service to estimate how much gas checkUpkeep functions will consume.
   * @param id identifier of the upkeep to check
   */
  function measureCheckGas(uint256 id)
    external
    returns (
      bool,
      bytes memory,
      uint256
    )
  {
    (, , , , address lastKeeper, , , ) = i_keeperRegistry.getUpkeep(id);
    (, , address[] memory keepers) = i_keeperRegistry.getState();

    uint256 index = block.number % keepers.length;
    address from = keepers[index];
    if (from == lastKeeper) {
      from = keepers[(index + 1) % keepers.length];
    }

    uint256 startGas = gasleft();
    try i_keeperRegistry.checkUpkeep(id, from) returns (
      bytes memory performData,
      uint256 maxLinkPayment,
      uint256 gasLimit,
      uint256 adjustedGasWei,
      uint256 linkEth
    ) {
      uint256 gasUsed = startGas - gasleft();
      return (true, performData, gasUsed);
    } catch {
      uint256 gasUsed = startGas - gasleft();
      return (false, "", gasUsed);
    }
  }
}
