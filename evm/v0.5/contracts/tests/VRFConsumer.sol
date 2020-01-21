pragma solidity 0.5.0;

import "../interfaces/LinkTokenInterface.sol";
import "../dev/VRFCoordinator.sol";
import "../dev/VRFConsumerBase.sol";

contract VRFConsumer is VRFConsumerBase {

  uint256 public randomnessOutput;
  bytes32 public requestId;

  constructor(address _vrfCoordinator, address _link)
    VRFConsumerBase(_vrfCoordinator, _link) public {

  }

  function fulfillRandomness(bytes32 _requestId, uint256 _randomness) external {
    randomnessOutput = _randomness;
    requestId =_requestId;
  }
}
