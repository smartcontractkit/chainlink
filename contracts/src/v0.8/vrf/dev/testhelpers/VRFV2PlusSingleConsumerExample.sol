// SPDX-License-Identifier: MIT
// Example of a single consumer contract which owns the subscription.
pragma solidity ^0.8.0;

import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {VRFConsumerBaseV2Plus} from "../VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

/// @notice This contract is used for testing only and should not be used for production.
contract VRFV2PlusSingleConsumerExample is VRFConsumerBaseV2Plus {
  // solhint-disable-next-line chainlink-solidity/prefix-storage-variables-with-s-underscore
  LinkTokenInterface internal LINKTOKEN;

  // solhint-disable-next-line gas-struct-packing
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
  address internal s_owner;

  constructor(
    address vrfCoordinator,
    address link,
    uint32 callbackGasLimit,
    uint16 requestConfirmations,
    uint32 numWords,
    bytes32 keyHash,
    bool nativePayment
  ) VRFConsumerBaseV2Plus(vrfCoordinator) {
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

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal override {
    // solhint-disable-next-line gas-custom-errors
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
    s_requestId = s_vrfCoordinator.requestRandomWords(req);
  }

  // Assumes this contract owns link
  // This method is analogous to VRFv1, except the amount
  // should be selected based on the keyHash (each keyHash functions like a "gas lane"
  // with different link costs).
  function fundAndRequestRandomWords(uint256 amount) external onlyOwner {
    RequestConfig memory rc = s_requestConfig;
    LINKTOKEN.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_requestConfig.subId));
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: rc.keyHash,
      subId: rc.subId,
      requestConfirmations: rc.requestConfirmations,
      callbackGasLimit: rc.callbackGasLimit,
      numWords: rc.numWords,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: rc.nativePayment}))
    });
    // Will revert if subscription is not set and funded.
    s_requestId = s_vrfCoordinator.requestRandomWords(req);
  }

  // Assumes this contract owns link
  function topUpSubscription(uint256 amount) external onlyOwner {
    LINKTOKEN.transferAndCall(address(s_vrfCoordinator), amount, abi.encode(s_requestConfig.subId));
  }

  function withdraw(uint256 amount, address to) external onlyOwner {
    LINKTOKEN.transfer(to, amount);
  }

  function unsubscribe(address to) external onlyOwner {
    // Returns funds to this address
    s_vrfCoordinator.cancelSubscription(s_requestConfig.subId, to);
    s_requestConfig.subId = 0;
  }

  // Keep this separate in case the contract want to unsubscribe and then
  // resubscribe.
  function subscribe() public onlyOwner {
    // Create a subscription, current subId
    address[] memory consumers = new address[](1);
    consumers[0] = address(this);
    s_requestConfig.subId = s_vrfCoordinator.createSubscription();
    s_vrfCoordinator.addConsumer(s_requestConfig.subId, consumers[0]);
  }
}
