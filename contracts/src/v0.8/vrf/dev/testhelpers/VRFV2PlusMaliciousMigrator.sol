// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IVRFMigratableConsumerV2Plus} from "../interfaces/IVRFMigratableConsumerV2Plus.sol";
import {IVRFCoordinatorV2Plus} from "../interfaces/IVRFCoordinatorV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

contract VRFV2PlusMaliciousMigrator is IVRFMigratableConsumerV2Plus {
  IVRFCoordinatorV2Plus internal s_vrfCoordinator;

  constructor(address _vrfCoordinator) {
    s_vrfCoordinator = IVRFCoordinatorV2Plus(_vrfCoordinator);
  }

  /**
   * @inheritdoc IVRFMigratableConsumerV2Plus
   */
  function setCoordinator(address /* _vrfCoordinator */) public override {
    // try to re-enter, should revert
    // args don't really matter
    s_vrfCoordinator.requestRandomWords(
      VRFV2PlusClient.RandomWordsRequest({
        keyHash: bytes32(0),
        subId: 0,
        requestConfirmations: 0,
        callbackGasLimit: 0,
        numWords: 0,
        extraArgs: ""
      })
    );
  }
}
