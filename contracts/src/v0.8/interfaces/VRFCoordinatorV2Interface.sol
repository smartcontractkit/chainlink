// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface VRFCoordinatorV2Interface {

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint64  subId,   // A data structure for billing
        uint16  minimumRequestConfirmations,
        uint32  callbackGasLimit,
        uint32  numWords,  // Desired number of random words
        uint32 consumerID
    )
    external
    returns (uint256 requestId);

    function createSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    returns (uint64 subId);

    function fundSubscription(
        uint64 subId,
        uint96 amount
    )
    external;

    function updateSubscription(
        uint64 subId,
        address[] memory consumers
    )
    external;
}
