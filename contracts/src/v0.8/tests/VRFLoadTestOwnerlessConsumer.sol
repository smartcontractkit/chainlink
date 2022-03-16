// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../VRFConsumerBase.sol";
import "../interfaces/ERC677ReceiverInterface.sol";

/**
 * @title The VRFLoadTestOwnerlessConsumer contract.
 * @notice Allows making many VRF V1 randomness requests in a single transaction for load testing.
 */
contract VRFLoadTestOwnerlessConsumer is VRFConsumerBase, ERC677ReceiverInterface {
  // The price of each VRF request in Juels. 1 LINK = 1e18 Juels.
  uint256 public immutable PRICE;

  uint256 public s_responseCount;

  constructor(
    address _vrfCoordinator,
    address _link,
    uint256 _price
  ) VRFConsumerBase(_vrfCoordinator, _link) {
    PRICE = _price;
  }

  function fulfillRandomness(bytes32, uint256) internal override {
    s_responseCount++;
  }

  /**
   * @dev Creates as many randomness requests as can be made with the funds transferred.
   * @param _amount The amount of LINK transferred to pay for these requests.
   * @param _data The data passed to transferAndCall on LinkToken. Must be an abi-encoded key hash.
   */
  function onTokenTransfer(
    address,
    uint256 _amount,
    bytes calldata _data
  ) external override {
    if (msg.sender != address(LINK)) {
      revert("only callable from LINK");
    }
    bytes32 keyHash = abi.decode(_data, (bytes32));

    uint256 spent = 0;
    while (spent + PRICE <= _amount) {
      requestRandomness(keyHash, PRICE);
      spent += PRICE;
    }
  }
}
