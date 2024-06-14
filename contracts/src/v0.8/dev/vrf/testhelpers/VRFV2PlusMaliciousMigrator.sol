// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../interfaces/IVRFMigratableConsumerV2Plus.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../libraries/VRFV2PlusClient.sol";

contract VRFV2PlusMaliciousMigrator is IVRFMigratableConsumerV2Plus {
  IVRFCoordinatorV2Plus s_vrfCoordinator;

  constructor(address _vrfCoordinator) {
    s_vrfCoordinator = IVRFCoordinatorV2Plus(_vrfCoordinator);
  }

  /**
   * @inheritdoc IVRFMigratableConsumerV2Plus
   */
  function setCoordinator(address _vrfCoordinator) public override {
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
