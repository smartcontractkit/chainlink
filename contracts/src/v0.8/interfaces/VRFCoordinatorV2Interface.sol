pragma solidity ^0.8.0;

interface VRFCoordinatorV2Interface {

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint16  minimumRequestConfirmations,
        uint16  callbackGasLimit,
        uint256 subId,   // A data structure for billing
        uint256 numWords  // Desired number of random words
    )
    external
    returns (uint256 requestId);

    function createSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    returns (uint256 subId);

    function fundSubscription(
        uint256 subId,
        uint256 amount
    )
    external;
}
