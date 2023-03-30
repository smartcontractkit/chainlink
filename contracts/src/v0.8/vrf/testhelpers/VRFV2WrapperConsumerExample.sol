// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../VRFV2WrapperConsumerBase.sol";
import "../../ConfirmedOwner.sol";

contract VRFV2WrapperConsumerExample is VRFV2WrapperConsumerBase, ConfirmedOwner {
  event WrappedRequestFulfilled(uint256 requestId, uint256[] randomWords, uint256 payment);
  event WrapperRequestMade(uint256 indexed requestId, uint256 paid);

  struct RequestStatus {
    uint256 paid;
    bool fulfilled;
    uint256[] randomWords;
  }
  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */
    public s_requests;

  constructor(address _link, address _vrfV2Wrapper)
    ConfirmedOwner(msg.sender)
    VRFV2WrapperConsumerBase(_link, _vrfV2Wrapper)
  {}

  function makeRequest(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external onlyOwner returns (uint256 requestId) {
    requestId = requestRandomness(_callbackGasLimit, _requestConfirmations, _numWords);
    uint256 paid = VRF_V2_WRAPPER.calculateRequestPrice(_callbackGasLimit);
    s_requests[requestId] = RequestStatus({paid: paid, randomWords: new uint256[](0), fulfilled: false});
    emit WrapperRequestMade(requestId, paid);
    return requestId;
  }

  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    require(s_requests[_requestId].paid > 0, "request not found");
    s_requests[_requestId].fulfilled = true;
    s_requests[_requestId].randomWords = _randomWords;
    emit WrappedRequestFulfilled(_requestId, _randomWords, s_requests[_requestId].paid);
  }

  function getRequestStatus(uint256 _requestId)
    external
    view
    returns (
      uint256 paid,
      bool fulfilled,
      uint256[] memory randomWords
    )
  {
    require(s_requests[_requestId].paid > 0, "request not found");
    RequestStatus memory request = s_requests[_requestId];
    return (request.paid, request.fulfilled, request.randomWords);
  }
}
