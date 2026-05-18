// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {VRFV2PlusWrapperConsumerBase} from "../VRFV2PlusWrapperConsumerBase.sol";
import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

contract VRFV2PlusWrapperConsumerExample is VRFV2PlusWrapperConsumerBase, ConfirmedOwner {
  event WrappedRequestFulfilled(uint256 requestId, uint256[] randomWords, uint256 payment);
  event WrapperRequestMade(uint256 indexed requestId, uint256 paid);

  // solhint-disable-next-line gas-struct-packing
  struct RequestStatus {
    uint256 paid;
    bool fulfilled;
    uint256[] randomWords;
    bool native;
  }

  mapping(uint256 => RequestStatus) /* requestId */ /* requestStatus */ public s_requests;

  constructor(address _vrfV2Wrapper) ConfirmedOwner(msg.sender) VRFV2PlusWrapperConsumerBase(_vrfV2Wrapper) {}

  function makeRequest(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external onlyOwner returns (uint256 requestId) {
    bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}));
    uint256 paid;
    (requestId, paid) = requestRandomness(_callbackGasLimit, _requestConfirmations, _numWords, extraArgs);
    s_requests[requestId] = RequestStatus({paid: paid, randomWords: new uint256[](0), fulfilled: false, native: false});
    emit WrapperRequestMade(requestId, paid);
    return requestId;
  }

  function makeRequestNative(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external onlyOwner returns (uint256 requestId) {
    bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}));
    uint256 paid;
    (requestId, paid) = requestRandomnessPayInNative(_callbackGasLimit, _requestConfirmations, _numWords, extraArgs);
    s_requests[requestId] = RequestStatus({paid: paid, randomWords: new uint256[](0), fulfilled: false, native: true});
    emit WrapperRequestMade(requestId, paid);
    return requestId;
  }

  // solhint-disable-next-line chainlink-solidity/prefix-internal-functions-with-underscore
  function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords) internal override {
    // solhint-disable-next-line gas-custom-errors
    require(s_requests[_requestId].paid > 0, "request not found");
    s_requests[_requestId].fulfilled = true;
    s_requests[_requestId].randomWords = _randomWords;
    emit WrappedRequestFulfilled(_requestId, _randomWords, s_requests[_requestId].paid);
  }

  function getRequestStatus(
    uint256 _requestId
  ) external view returns (uint256 paid, bool fulfilled, uint256[] memory randomWords) {
    // solhint-disable-next-line gas-custom-errors
    require(s_requests[_requestId].paid > 0, "request not found");
    RequestStatus memory request = s_requests[_requestId];
    return (request.paid, request.fulfilled, request.randomWords);
  }

  /// @notice withdrawLink withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in juels
  function withdrawLink(uint256 amount) external onlyOwner {
    i_linkToken.transfer(owner(), amount);
  }

  /// @notice withdrawNative withdraws the amount specified in amount to the owner
  /// @param amount the amount to withdraw, in wei
  function withdrawNative(uint256 amount) external onlyOwner {
    (bool success, ) = payable(owner()).call{value: amount}("");
    // solhint-disable-next-line gas-custom-errors
    require(success, "withdrawNative failed");
  }

  event Received(address, uint256);

  receive() external payable {
    emit Received(msg.sender, msg.value);
  }
}
