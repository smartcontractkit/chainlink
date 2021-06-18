pragma solidity ^0.8.0;

interface VRFCoordinatorV2Interface {

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint32  subId,   // A data structure for billing
        uint16  minimumRequestConfirmations,
        uint32  callbackGasLimit,
        uint16 numWords  // Desired number of random words
    )
    external
    returns (uint256 requestId);

    function createSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    returns (uint32 subId);

    function fundSubscription(
        uint32 subId,
        uint256 amount
    )
    external;
}
