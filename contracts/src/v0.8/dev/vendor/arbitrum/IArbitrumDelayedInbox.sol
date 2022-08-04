// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";

interface IArbitrumDelayedInbox is IInbox {
    function calculateRetryableSubmissionFee(uint256 dataLength, uint256 baseFee)
        external
        view
        returns (uint256);
}
