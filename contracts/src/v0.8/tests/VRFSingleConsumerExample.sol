// SPDX-License-Identifier: MIT
// Example of a single consumer contract which owns the subscription.
pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../dev/VRFConsumerBaseV2.sol";

contract VRFSingleConsumerExample is VRFConsumerBaseV2 {

    VRFCoordinatorV2Interface COORDINATOR;
    LinkTokenInterface LINKTOKEN;

    struct RequestConfig {
        uint64 subId;
        uint32 callbackGasLimit;
        uint16 requestConfirmations;
        uint32 numWords;
        bytes32 jobID;
    }
    RequestConfig s_requestConfig;
    uint256[] s_randomWords;
    uint256 s_requestId;

    constructor(
        address vrfCoordinator,
        address link,
        uint32 callbackGasLimit,
        uint16 requestConfirmations,
        uint32 numWords
    )
    VRFConsumerBaseV2(vrfCoordinator)
    {
        COORDINATOR = VRFCoordinatorV2Interface(vrfCoordinator);
        LINKTOKEN = LinkTokenInterface(link);
        s_requestConfig = RequestConfig({
            subId: 0, // Unset
            callbackGasLimit: callbackGasLimit,
            requestConfirmations: requestConfirmations,
            numWords: numWords,
            jobID: bytes32(0)
        });
    }

    function fulfillRandomWords(
        uint256 requestId,
        uint256[] memory randomWords
    )
        internal
        override
    {
        s_randomWords = randomWords;
    }

    function requestRandomWords()
        external
    {
        RequestConfig memory rc = s_requestConfig;
        // Will revert if subscription is not set and funded.
        s_requestId = COORDINATOR.requestRandomWords(
            rc.jobID,
            rc.subId,
            rc.requestConfirmations,
            rc.callbackGasLimit,
            rc.numWords);
    }

    // Assumes this contract owns link
    function topUpSubscription(
        uint256 amount
    )
        external
    {
        LINKTOKEN.transferAndCall(
            address(COORDINATOR),
            amount,
            abi.encode(s_requestConfig.subId));
    }

    function unsubscribe()
        external
    {
        // Returns funds to this address
        COORDINATOR.cancelSubscription(s_requestConfig.subId, address(this));
        s_requestConfig.subId = 0;
    }

    function subscribe()
        external
    {
        address[] memory consumers = new address[](1);
        consumers[0] = address(this);
        s_requestConfig.subId = COORDINATOR.createSubscription(consumers);
    }

    // Allows users to dynamically select jobs with particular
    // gas price ceilings.
    function setJobID(
        bytes32 jobID
    )
        public
    {
        s_requestConfig.jobID = jobID;
    }
}
