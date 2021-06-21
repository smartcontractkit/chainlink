// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";

contract VRFConsumerV2 {

    uint256[] public randomWords;
    uint256 public requestId;
    VRFCoordinatorV2Interface COORDINATOR;
    LinkTokenInterface LINKTOKEN;
    uint64 public subId;
    uint256 public gasAvailable;

    constructor(address vrfCoordinator, address link)
    {
        COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
        LINKTOKEN = LinkTokenInterface(link);
    }

    function fulfillRandomWords(uint256 _requestId, uint256[] memory _randomWords)
    external
    {
        gasAvailable = gasleft();
        randomWords = _randomWords;
        requestId = _requestId;
    }

    function testCreateSubscriptionAndFund(uint256 amount) external {
        if (subId == 0) {
            address[] memory consumers = new address[](1);
            consumers[0] = address(this);
            subId = COORDINATOR.createSubscription(consumers);
        }
        // Approve the link transfer.
        LINKTOKEN.approve(address(COORDINATOR), amount);
        // Transfer link to the coordinator.
        COORDINATOR.fundSubscription(subId, amount);
    }

    function testRequestRandomness(bytes32 _keyHash, uint64 _subId, uint64 minReqConfs, uint64 callbackGasLimit, uint64 numWords)
    external
returns (uint256)
    {
        return COORDINATOR.requestRandomWords(_keyHash, _subId, minReqConfs, callbackGasLimit, numWords);
    }
}
