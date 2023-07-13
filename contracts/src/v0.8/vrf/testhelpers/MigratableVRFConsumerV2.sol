// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/VRFCoordinatorV2Interface.sol";
import "./MigratableVRFConsumerBaseV2.sol";

contract MigratableVRFConsumerV2 is MigratableVRFConsumerBaseV2 {
  mapping(uint256 => uint256[]) public s_randomWords;
  uint256 public s_requestId;

  bytes4 private constant REQUEST_RANDOM_WORDS_SELECTOR = bytes4(keccak256("requestRandomWords(bytes32,uint64,uint16,uint32,uint32,bool)"));

  constructor(address vrfCoordinator, uint64 subId) MigratableVRFConsumerBaseV2(vrfCoordinator, subId) {}

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    require(requestId == s_requestId, "request ID is incorrect");
    s_randomWords[s_requestId] = randomWords;
  }

  function requestRandomness(
    bytes32 keyHash,
    uint16 minReqConfs,
    uint32 callbackGasLimit,
    uint32 numWords,
    bool nativePayment
  ) external returns (uint256) {
    bytes memory callData = abi.encodeWithSelector(REQUEST_RANDOM_WORDS_SELECTOR, keyHash, s_subId, minReqConfs, callbackGasLimit, numWords, nativePayment);
    // solhint-disable-next-line avoid-low-level-calls
    (bool success, bytes memory ret) = s_vrfCoordinator.call(callData);
    require(success, "request random words failed");
    s_requestId = uint256(bytes32(ret));
    return s_requestId;
  }
}
