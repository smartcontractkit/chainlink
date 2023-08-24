// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface VRFV2PlusWrapperInterface {
  /**
   * @return the request ID of the most recent VRF V2 request made by this wrapper. This should only
   * be relied option within the same transaction that the request was made.
   */
  function lastRequestId() external view returns (uint256);

  /**
   * @notice Requests randomness from the VRF V2 wrapper, paying in native token.
   *
   * @param _callbackGasLimit is the gas limit for the request.
   * @param _requestConfirmations number of request confirmations to wait before serving a request.
   * @param _numWords is the number of words to request.
   */
  function requestRandomWordsInNative(
    uint32 _callbackGasLimit,
    uint16 _requestConfirmations,
    uint32 _numWords
  ) external payable returns (uint256 requestId);

  /**
   * @notice Returns the price registry being used by this wrapper.
  */
  function getPriceRegistry() external view returns (address);
}
