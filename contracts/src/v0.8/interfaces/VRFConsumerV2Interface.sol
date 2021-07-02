// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface VRFConsumerV2Interface {
    function fulfillRandomWords(
        uint256 requestId,
        uint256[] memory randomWords
    )
    external;
}
