// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

interface IFunctionsRequest {
  struct Commitment {
    uint96 adminFee; // -----------┐
    address coordinator; // -------┘
    address client; // ------------┐
    uint64 subscriptionId; //      |
    uint32 callbackGasLimit; // ---┘
    uint96 estimatedTotalCostJuels; // TODO pack the following
    uint40 timeoutTimestamp;
    bytes32 requestId;
    uint80 donFee;
    uint40 gasOverheadBeforeCallback;
    uint40 gasOverheadAfterCallback;
  }
}
