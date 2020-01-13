pragma solidity 0.5.0;

import "../interfaces/LinkTokenInterface.sol";
import "../dev/VRFCoordinator.sol";

contract VRFConsumer {

  LinkTokenInterface LINK;
  address vrfCoordinator;
  uint256 public randomnessOutput;
  bytes32 public requestId;

  constructor(address _vrfCoordinator, address _link) public {
    vrfCoordinator = _vrfCoordinator;
    LINK = LinkTokenInterface(_link);
  }

  function requestRandomness(bytes32 _sAId, uint256 _fee, uint256 _seed) external returns (uint256) {
    LINK.transferAndCall(vrfCoordinator, _fee, abi.encode(_sAId, _seed));
  }

  function fulfillRandomness(bytes32 _requestId, uint256 _seed) external {
    randomnessOutput = _seed;
    requestId =_requestId;
  }
}
