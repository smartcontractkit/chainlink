// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../dev/VRFConsumerBaseV2.sol";

contract VRFConsumerV2 is VRFConsumerBaseV2 {
    uint256[] public s_randomWords;
    uint256 public s_requestId;
    VRFCoordinatorV2Interface COORDINATOR;
    LinkTokenInterface LINKTOKEN;
    uint64 public s_subId;
    uint256 public s_gasAvailable;

    constructor(
        address vrfCoordinator,
        address link
    )
        VRFConsumerBaseV2(vrfCoordinator)
    {
        COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
        LINKTOKEN = LinkTokenInterface(link);
    }

    function fulfillRandomWords(
        uint256 requestId,
        uint256[] memory randomWords
    )
        internal
        override
    {
        s_gasAvailable = gasleft();
        s_randomWords = randomWords;
        s_requestId = requestId;
    }

    function testCreateSubscriptionAndFund(
        uint96 amount
    )
        external
    {
        if (s_subId == 0) {
            s_subId = COORDINATOR.createSubscription();
            COORDINATOR.addConsumer(s_subId, address(this));
        }
        // Approve the link transfer.
        LINKTOKEN.transferAndCall(address(COORDINATOR), amount, abi.encode(s_subId));
    }

    function updateSubscription(
        address[] memory consumers
    )
        external
    {
        require(s_subId != 0, "subID not set");
        for (uint256 i = 0; i < consumers.length; i++)  {
            COORDINATOR.addConsumer(s_subId, consumers[i]);
        }
    }

    function testRequestRandomness(
        bytes32 keyHash,
        uint64 subId,
        uint16 minReqConfs,
        uint32 callbackGasLimit,
        uint32 numWords
    )
        external
        returns (uint256)
    {
        return COORDINATOR.requestRandomWords(keyHash, subId, minReqConfs, callbackGasLimit, numWords);
    }
}
