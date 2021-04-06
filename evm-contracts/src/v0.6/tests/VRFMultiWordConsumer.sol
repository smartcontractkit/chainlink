// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

import "../interfaces/LinkTokenInterface.sol";
import "../VRFCoordinator.sol";
import "../VRFConsumerBase.sol";

contract VRFMultiWordConsumer is VRFConsumerBase {

    uint256 numRandomWords;
    bytes32[] public randomnessOutput;
    bytes32 public requestId;

    constructor(address _vrfCoordinator, address _link) public
        // solhint-disable-next-line no-empty-blocks
        VRFConsumerBase(_vrfCoordinator, _link) { /* empty */ }

    function fulfillRandomness(bytes32 _requestId, uint256 _randomness)
        internal override
    {
        for (uint256 i = 0; i < randomnessOutput.length; i++) {
            randomnessOutput[i] = keccak256(abi.encode(_randomness, i));
        }
        requestId = _requestId;
    }

    function testRequestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed, uint256 _numRandomWords)
        external returns (bytes32 _requestId)
    {
        randomnessOutput = new bytes32[](_numRandomWords);
        return requestRandomness(_keyHash, _fee, _seed);
    }
}
