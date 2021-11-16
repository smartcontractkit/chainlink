// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../dev/VRFConsumerBaseV2.sol";

contract VRFExternalSubOwnerExample is VRFConsumerBaseV2 {
  VRFCoordinatorV2Interface COORDINATOR;
  LinkTokenInterface LINKTOKEN;

  struct RequestConfig {
    uint64 subId;
    uint32 callbackGasLimit;
    uint16 requestConfirmations;
    uint32 numWords;
    bytes32 keyHash;
  }
  RequestConfig public s_requestConfig;
  uint256[] public s_randomWords;
  uint256 public s_requestId;
  address s_owner;

  constructor(
    address vrfCoordinator,
    address link,
    uint32 callbackGasLimit,
    uint16 requestConfirmations,
    uint32 numWords,
    bytes32 keyHash
  ) VRFConsumerBaseV2(vrfCoordinator) {
    COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(link);
    s_requestConfig = RequestConfig({
      subId: 0, // Initially unset
      callbackGasLimit: callbackGasLimit,
      requestConfirmations: requestConfirmations,
      numWords: numWords,
      keyHash: keyHash
    });
    s_owner = msg.sender;
  }

  function fulfillRandomWords(
    uint256, /* requestId */
    uint256[] memory randomWords
  ) internal override {
    s_randomWords = randomWords;
  }

  function requestRandomWords() external onlyOwner {
    RequestConfig memory rc = s_requestConfig;
    // Will revert if subscription is not set and funded.
    s_requestId = COORDINATOR.requestRandomWords(
      rc.keyHash,
      rc.subId,
      rc.requestConfirmations,
      rc.callbackGasLimit,
      rc.numWords
    );
  }

  function setSubscriptionID(uint64 subId) external onlyOwner {
    s_requestConfig.subId = subId;
  }

  modifier onlyOwner() {
    require(msg.sender == s_owner);
    _;
  }
}
