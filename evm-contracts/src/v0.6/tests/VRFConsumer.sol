pragma solidity 0.6.2;

import "../interfaces/LinkTokenInterface.sol";
import "../VRFCoordinator.sol";
import "../VRFConsumerBase.sol";

contract VRFConsumer is VRFConsumerBase {

  uint256 public randomnessOutput;
  bytes32 public requestId;

  constructor(address _vrfCoordinator, address _link) public
    // solhint-disable-next-line no-empty-blocks
    VRFConsumerBase(_vrfCoordinator, _link) { /* empty */ }

  function fulfillRandomness(bytes32 _requestId, uint256 _randomness)
    external override
  {
    randomnessOutput = _randomness;
    requestId = _requestId;
  }
}
