// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ILinkToken.sol";
import "../interfaces/IVRFCoordinatorV2.sol";
import "../VRFConsumerBaseV2.sol";
import "../ConfirmedOwner.sol";

/**
 * @title The VRFLoadTestExternalSubOwner contract.
 * @notice Allows making many VRF V2 randomness requests in a single transaction for load testing.
 */
contract VRFLoadTestExternalSubOwner is VRFConsumerBaseV2, ConfirmedOwner {
  IVRFCoordinatorV2 public immutable COORDINATOR;
  ILinkToken public immutable LINK;

  uint256 public s_responseCount;

  constructor(address _vrfCoordinator, address _link) VRFConsumerBaseV2(_vrfCoordinator) ConfirmedOwner(msg.sender) {
    COORDINATOR = IVRFCoordinatorV2(_vrfCoordinator);
    LINK = ILinkToken(_link);
  }

  function fulfillRandomWords(uint256, uint256[] memory) internal override {
    s_responseCount++;
  }

  function requestRandomWords(
    uint64 _subId,
    uint16 _requestConfirmations,
    bytes32 _keyHash,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      COORDINATOR.requestRandomWords(_keyHash, _subId, _requestConfirmations, 50_000, 1);
    }
  }
}
