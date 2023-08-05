// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsSubscriptions} from "../interfaces/IFunctionsSubscriptions.sol";

// @title Library of types that are used for fulfillment of a Functions request
library FunctionsResponse {
  // Used to send request information from the Router to the Coordinator
  struct RequestMeta {
    address requestingContract; // The client contract that is sending the request
    bytes data; // CBOR encoded Chainlink Functions request data, use FunctionsRequest library to encode a request
    uint64 subscriptionId; // Identifier of the billing subscription that will be charged for the request
    uint16 dataVersion; // The version of the structure of the CBOR encoded request data
    bytes32 flags; // Per-subscription account flags
    uint32 callbackGasLimit; // The amount of gas that the callback to the consuming contract will be given
    uint72 adminFee; // Flat fee (in Juels of LINK) that will be paid to the Router Owner for operation of the network
    IFunctionsSubscriptions.Consumer consumer; // Details about the consumer making the request
    IFunctionsSubscriptions.Subscription subscription; // Details about the subscription that will be charged for the request
  }

  enum FulfillResult {
    USER_SUCCESS, // 0
    USER_ERROR, // 1
    INVALID_REQUEST_ID, // 2
    COST_EXCEEDS_COMMITMENT, // 3
    INSUFFICIENT_GAS_PROVIDED, // 4
    SUBSCRIPTION_BALANCE_INVARIANT_VIOLATION, // 5
    INVALID_COMMITMENT // 6
  }

  struct Commitment {
    uint72 adminFee; // -----------┐
    address coordinator; // -------┘
    address client; // ------------┐
    uint64 subscriptionId; //      |
    uint32 callbackGasLimit; // ---┘
    uint96 estimatedTotalCostJuels;
    uint40 timeoutTimestamp;
    bytes32 requestId;
    uint72 donFee;
    uint40 gasOverheadBeforeCallback;
    uint40 gasOverheadAfterCallback;
  }
}
