// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {DualAggregator} from "../../DualAggregator.sol";

import {AccessControllerInterface} from "../../../shared/interfaces/AccessControllerInterface.sol";
import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";

contract DualAggregatorHelper is DualAggregator {
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_,
    address secondaryProxy_,
    uint32 cutoffTime_,
    uint32 maxSyncIterations_
  )
    DualAggregator(
      link,
      minAnswer_,
      maxAnswer_,
      billingAccessController,
      requesterAccessController,
      decimals_,
      description_,
      secondaryProxy_,
      cutoffTime_,
      maxSyncIterations_
    )
  {}

  function getSyncPrimaryRound() public view returns (uint32 roundId) {
    return _getSyncPrimaryRound();
  }

  function configDigestFromConfigData(
    uint256 chainId,
    address contractAddress,
    uint64 configCount,
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) public pure returns (bytes32) {
    return _configDigestFromConfigData(
      chainId,
      contractAddress,
      configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function totalLinkDue() public view returns (uint256 linkDue) {
    return _totalLinkDue();
  }

  function getHotVars() public view returns (HotVars memory) {
    return s_hotVars;
  }
}
