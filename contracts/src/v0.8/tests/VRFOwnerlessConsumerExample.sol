// SPDX-License-Identifier: MIT
// An example VRF V1 consumer contract that can be triggered using a transferAndCall from the link
// contract.
pragma solidity ^0.8.0;

import "../VRFConsumerBase.sol";
import "../interfaces/ERC677ReceiverInterface.sol";

contract VRFOwnerlessConsumerExample is VRFConsumerBase, ERC677ReceiverInterface {
  uint256 public s_randomnessOutput;
  bytes32 public s_requestId;

  error OnlyCallableFromLink();

  constructor(address _vrfCoordinator, address _link) VRFConsumerBase(_vrfCoordinator, _link) {
    /* empty */
  }

  function fulfillRandomness(bytes32 requestId, uint256 _randomness) internal override {
    require(requestId == s_requestId, "request ID is incorrect");
    s_randomnessOutput = _randomness;
  }

  /**
   * @dev Creates a new randomness request. This function can only be used by calling
   * transferAndCall on the LinkToken contract.
   * @param _amount The amount of LINK transferred to pay for this request.
   * @param _data The data passed to transferAndCall on LinkToken. Must be an abi-encoded key hash.
   */
  function onTokenTransfer(
    address, /* sender */
    uint256 _amount,
    bytes calldata _data
  ) external override {
    if (msg.sender != address(LINK)) {
      revert OnlyCallableFromLink();
    }

    bytes32 keyHash = abi.decode(_data, (bytes32));
    s_requestId = requestRandomness(keyHash, _amount);
  }
}
