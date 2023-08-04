// SPDX-License-Identifier: MIT
// Example of a single consumer contract which owns the subscription.
pragma solidity ^0.8.0;

import "../../../shared/interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";

/// @notice This contract is used for testing only and should not be used for production.
contract VRFV2PlusSingleConsumerExample is VRFConsumerBaseV2Plus {
  IVRFCoordinatorV2Plus COORDINATOR;
  LinkTokenInterface LINKTOKEN;

  struct RequestConfig {
    uint256 subId;
    uint32 callbackGasLimit;
    uint16 requestConfirmations;
    uint32 numWords;
    bytes32 keyHash;
    bool nativePayment;
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
    bytes32 keyHash,
    bool nativePayment
  ) VRFConsumerBaseV2Plus(vrfCoordinator) {
    COORDINATOR = IVRFCoordinatorV2Plus(vrfCoordinator);
    LINKTOKEN = LinkTokenInterface(link);
    s_owner = msg.sender;
    s_requestConfig = RequestConfig({
      subId: 0, // Unset initially
      callbackGasLimit: callbackGasLimit,
      requestConfirmations: requestConfirmations,
      numWords: numWords,
      keyHash: keyHash,
      nativePayment: nativePayment
    });
    subscribe();
  }

  function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
    require(requestId == s_requestId, "request ID is incorrect");
    s_randomWords = randomWords;
  }

  // Assumes the subscription is funded sufficiently.
  function requestRandomWords() external onlyOwner {
    RequestConfig memory rc = s_requestConfig;
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: rc.keyHash,
      subId: rc.subId,
      requestConfirmations: rc.requestConfirmations,
      callbackGasLimit: rc.callbackGasLimit,
      numWords: rc.numWords,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: rc.nativePayment}))
    });
    // Will revert if subscription is not set and funded.
    s_requestId = COORDINATOR.requestRandomWords(req);
  }

  // Assumes this contract owns link
  // This method is analogous to VRFv1, except the amount
  // should be selected based on the keyHash (each keyHash functions like a "gas lane"
  // with different link costs).
  function fundAndRequestRandomWords(uint256 amount) external onlyOwner {
    RequestConfig memory rc = s_requestConfig;
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_requestConfig.subId));
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: rc.keyHash,
      subId: rc.subId,
      requestConfirmations: rc.requestConfirmations,
      callbackGasLimit: rc.callbackGasLimit,
      numWords: rc.numWords,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: rc.nativePayment}))
    });
    // Will revert if subscription is not set and funded.
    s_requestId = COORDINATOR.requestRandomWords(req);
  }

  // Assumes this contract owns link
  function topUpSubscription(uint256 amount) external onlyOwner {
    LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_requestConfig.subId));
  }

  function withdraw(uint256 amount, address to) external onlyOwner {
    LINKTOKEN.transfer(to, amount);
  }

  function unsubscribe(address to) external onlyOwner {
    // Returns funds to this address
    COORDINATOR.cancelSubscription(s_requestConfig.subId, to);
    s_requestConfig.subId = 0;
  }

  // Keep this separate in case the contract want to unsubscribe and then
  // resubscribe.
  function subscribe() public onlyOwner {
    // Create a subscription, current subId
    address[] memory consumers = new address[](1);
    consumers[0] = address(this);
    s_requestConfig.subId = COORDINATOR.createSubscription();
    COORDINATOR.addConsumer(s_requestConfig.subId, consumers[0]);
  }
}
