pragma solidity ^0.8.0;

abstract contract VRFConsumerBaseV2 {
  error OnlyCoordinatorCanFulfill(address have, address want);
  address immutable private vrfCoordinator;

  constructor(
    address _vrfCoordinator
  )
  {
      vrfCoordinator = _vrfCoordinator;
  }

  function fulfillRandomWords(
    uint256 requestId,
    uint256[] memory randomWords
  )
    internal virtual;

  function rawFulfillRandomWords(
    uint256 requestId,
    uint256[] memory randomWords
  )
    external
  {
    if (msg.sender != vrfCoordinator) {
      revert OnlyCoordinatorCanFulfill(msg.sender, vrfCoordinator);
    }
    fulfillRandomWords(requestId, randomWords);
  }
}
