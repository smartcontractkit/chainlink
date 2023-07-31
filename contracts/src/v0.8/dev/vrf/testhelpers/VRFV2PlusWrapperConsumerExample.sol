// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../VRFV2PlusWrapperConsumerBase.sol";
import "../../../shared/access/ConfirmedOwner.sol";

contract VRFV2PlusWrapperConsumerExample is VRFV2PlusWrapperConsumerBase, ConfirmedOwner {
  event WrappedRequestFulfilled(uint256 requestId, uint256[] randomWords, uint256 payment);
  event WrapperRequestMade(uint256 indexed requestId, uint256 paid);

  struct RequestStatus {
    uint256 paid;
    bool fulfilled;
    uint256[] randomWords;
    bool native;
  }

  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */ public s_requests;

  constructor(
    address _link,
    address _vrfV2Wrapper
  ) ConfirmedOwner(msg.sender) VRFV2PlusWrapperConsumerBase(_link, _vrfV2Wrapper) {}

  function makeRequest(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external onlyOwner returns (uint256 requestId) {
    requestId = requestRandomness(_callbackGasLimit, _requestConfirmations, _numWords);
    uint256 paid = VRF_V2_PLUS_WRAPPER.calculateRequestPrice(_callbackGasLimit);
    s_requests[requestId] = RequestStatus({paid: paid, randomWords: new uint256[](0), fulfilled: false, native: false});
    emit WrapperRequestMade(requestId, paid);
    return requestId;
  }

  function makeRequestNative(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external onlyOwner returns (uint256 requestId) {
    requestId = requestRandomnessPayInNative(_callbackGasLimit, _requestConfirmations, _numWords);
    uint256 paid = VRF_V2_PLUS_WRAPPER.calculateRequestPriceNative(_callbackGasLimit);
    s_requests[requestId] = RequestStatus({paid: paid, randomWords: new uint256[](0), fulfilled: false, native: true});
    emit WrapperRequestMade(requestId, paid);
    return requestId;
  }

  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    require(s_requests[_requestId].paid > 0, "request not found");
    s_requests[_requestId].fulfilled = true;
    s_requests[_requestId].randomWords = _randomWords;
    emit WrappedRequestFulfilled(_requestId, _randomWords, s_requests[_requestId].paid);
  }

  function getRequestStatus(
    uint256 _requestId
  ) external view returns (uint256 paid, bool fulfilled, uint256[] memory randomWords) {
    require(s_requests[_requestId].paid > 0, "request not found");
    RequestStatus memory request = s_requests[_requestId];
    return (request.paid, request.fulfilled, request.randomWords);
  }

  /// @notice withdrawLink withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in juels
  function withdrawLink(uint256 amount) external onlyOwner {
    LINK.transfer(owner(), amount);
  }

  /// @notice withdrawNative withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in wei
  function withdrawNative(uint256 amount) external onlyOwner {
    (bool success, ) = payable(owner()).call{value: amount}("");
    require(success, "withdrawNative failed");
  }
}
