// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";

contract VRFConsumerV2 {
    uint256[] public s_randomWords;
    uint256 public s_requestId;
    VRFCoordinatorV2Interface COORDINATOR;
    LinkTokenInterface LINKTOKEN;
    uint64 public s_subId;
    uint256 public s_gasAvailable;

    constructor(address vrfCoordinator, address link)
    {
        COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
        LINKTOKEN = LinkTokenInterface(link);
    }

    function fulfillRandomWords(
        uint256 requestId,
        uint256[] memory randomWords
    )
    external
    {
        s_gasAvailable = gasleft();
        s_randomWords = randomWords;
        s_requestId = requestId;
    }

    function testCreateSubscriptionAndFund(
        uint256 amount
    )
    external
    {
        if (s_subId == 0) {
            address[] memory consumers = new address[](1);
            consumers[0] = address(this);
            s_subId = COORDINATOR.createSubscription(consumers);
        }
        // Approve the link transfer.
        LINKTOKEN.approve(address(COORDINATOR), amount);
        // Transfer link to the coordinator.
        COORDINATOR.fundSubscription(s_subId, amount);
    }

    function testRequestRandomness(
        bytes32 keyHash,
        uint64 subId,
        uint64 minReqConfs,
        uint64 callbackGasLimit,
        uint64 numWords)
    external
    returns (uint256)
    {
        return COORDINATOR.requestRandomWords(keyHash, subId, minReqConfs, callbackGasLimit, numWords);
    }
}
