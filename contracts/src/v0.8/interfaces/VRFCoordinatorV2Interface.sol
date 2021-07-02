// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface VRFCoordinatorV2Interface {

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint64  subId,   // A data structure for billing
        uint64  minimumRequestConfirmations,
        uint64  callbackGasLimit,
        uint64  numWords  // Desired number of random words
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
}
