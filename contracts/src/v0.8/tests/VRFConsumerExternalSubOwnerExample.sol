pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../dev/VRFConsumerBaseV2.sol";

contract VRFConsumerExternalSubOwnerExample is VRFConsumerBaseV2 {

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
            subId: 0, // Initially unset
            callbackGasLimit: callbackGasLimit,
            requestConfirmations: requestConfirmations,
            numWords: numWords,
            jobID: bytes32(0) // Initially unset
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

    function setSubscriptionID(
        uint64 subId
    )
        public
    {
        s_requestConfig.subId = subId;
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
