// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../VRFConsumerBaseV2.sol";

/**
 * @title The VRFLoadTestExternalSubOwner contract.
 * @notice Allows making many VRF V2 randomness requests in a single transaction for load testing.
 */
contract VRFLoadTestExternalSubOwner is VRFConsumerBaseV2 {
  VRFCoordinatorV2Interface public immutable COORDINATOR;
  LinkTokenInterface public immutable LINK;

  uint256 public s_responseCount;
  address s_owner;

  constructor(address _vrfCoordinator, address _link) VRFConsumerBaseV2(_vrfCoordinator) {
    COORDINATOR = VRFCoordinatorV2Interface(_vrfCoordinator);
    LINK = LinkTokenInterface(_link);
    s_owner = msg.sender;
  }

  function fulfillRandomWords(uint256, uint256[] memory) internal override {
    s_responseCount++;
  }

  function requestRandomWords(
    uint64 _subId,
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords,
    bytes32 _keyHash,
    uint16 _requestCount
  ) external onlyOwner {
    for (uint16 i = 0; i < _requestCount; i++) {
      COORDINATOR.requestRandomWords(
        _keyHash,
        _subId,
        _requestConfirmations,
        _callbackGasLimit,
        _numWords
      );
    }
  }

  function transferOwnership(address newOwner) external onlyOwner {
    s_owner = newOwner;
  }

  modifier onlyOwner() {
    require(msg.sender == s_owner);
    _;
  }
}
