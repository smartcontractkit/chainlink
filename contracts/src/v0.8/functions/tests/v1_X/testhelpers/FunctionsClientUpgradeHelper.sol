// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRequest} from "../../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsClient} from "../../../dev/v1_X/FunctionsClient.sol";
import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";

contract FunctionsClientUpgradeHelper is FunctionsClient, ConfirmedOwner {
  using FunctionsRequest for FunctionsRequest.Request;

  constructor(address router) FunctionsClient(router) ConfirmedOwner(msg.sender) {}

  event ResponseReceived(bytes32 indexed requestId, bytes result, bytes err);

  /**
   * @notice Send a simple request
   *
   * @param donId DON ID
   * @param source JavaScript source code
   * @param secrets Encrypted secrets payload
   * @param args List of arguments accessible from within the source code
   * @param subscriptionId Funtions billing subscription ID
   * @param callbackGasLimit Maximum amount of gas used to call the client contract's `handleOracleFulfillment` function
   * @return Functions request ID
   */
  function sendRequest(
    bytes32 donId,
    string calldata source,
    bytes calldata secrets,
    string[] calldata args,
    bytes[] memory bytesArgs,
    uint64 subscriptionId,
    uint32 callbackGasLimit
  ) public onlyOwner returns (bytes32) {
    FunctionsRequest.Request memory req;
    req._initializeRequestForInlineJavaScript(source);
    if (secrets.length > 0) req._addSecretsReference(secrets);
    if (args.length > 0) req._setArgs(args);
    if (bytesArgs.length > 0) req._setBytesArgs(bytesArgs);

    return _sendRequest(FunctionsRequest._encodeCBOR(req), subscriptionId, callbackGasLimit, donId);
  }

  function sendRequestBytes(
    bytes memory data,
    uint64 subscriptionId,
    uint32 callbackGasLimit,
    bytes32 donId
  ) public returns (bytes32 requestId) {
    return _sendRequest(data, subscriptionId, callbackGasLimit, donId);
  }

  /**
   * @notice Same as sendRequest but for DONHosted secrets
   */
  function sendRequestWithDONHostedSecrets(
    bytes32 donId,
    string calldata source,
    uint8 slotId,
    uint64 slotVersion,
    string[] calldata args,
    uint64 subscriptionId,
    uint32 callbackGasLimit
  ) public onlyOwner returns (bytes32) {
    FunctionsRequest.Request memory req;
    req._initializeRequestForInlineJavaScript(source);
    req._addDONHostedSecrets(slotId, slotVersion);

    if (args.length > 0) req._setArgs(args);

    return _sendRequest(FunctionsRequest._encodeCBOR(req), subscriptionId, callbackGasLimit, donId);
  }

  // @notice Sends a Chainlink Functions request
  // @param data The CBOR encoded bytes data for a Functions request
  // @param subscriptionId The subscription ID that will be charged to service the request
  // @param callbackGasLimit the amount of gas that will be available for the fulfillment callback
  // @return requestId The generated request ID for this request
  function _sendRequestToProposed(
    bytes memory data,
    uint64 subscriptionId,
    uint32 callbackGasLimit,
    bytes32 donId
  ) internal returns (bytes32) {
    bytes32 requestId = i_functionsRouter.sendRequestToProposed(
      subscriptionId,
      data,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      donId
    );
    emit RequestSent(requestId);
    return requestId;
  }

  /**
   * @notice Send a simple request to the proposed contract
   *
   * @param donId DON ID
   * @param source JavaScript source code
   * @param secrets Encrypted secrets payload
   * @param args List of arguments accessible from within the source code
   * @param subscriptionId Funtions billing subscription ID
   * @param callbackGasLimit Maximum amount of gas used to call the client contract's `handleOracleFulfillment` function
   * @return Functions request ID
   */
  function sendRequestToProposed(
    bytes32 donId,
    string calldata source,
    bytes calldata secrets,
    string[] calldata args,
    bytes[] memory bytesArgs,
    uint64 subscriptionId,
    uint32 callbackGasLimit
  ) public onlyOwner returns (bytes32) {
    FunctionsRequest.Request memory req;
    req._initializeRequestForInlineJavaScript(source);
    if (secrets.length > 0) req._addSecretsReference(secrets);
    if (args.length > 0) req._setArgs(args);
    if (bytesArgs.length > 0) req._setBytesArgs(bytesArgs);

    return _sendRequestToProposed(FunctionsRequest._encodeCBOR(req), subscriptionId, callbackGasLimit, donId);
  }

  /**
   * @notice Same as sendRequestToProposed but for DONHosted secrets
   */
  function sendRequestToProposedWithDONHostedSecrets(
    bytes32 donId,
    string calldata source,
    uint8 slotId,
    uint64 slotVersion,
    string[] calldata args,
    uint64 subscriptionId,
    uint32 callbackGasLimit
  ) public onlyOwner returns (bytes32) {
    FunctionsRequest.Request memory req;
    req._initializeRequestForInlineJavaScript(source);
    req._addDONHostedSecrets(slotId, slotVersion);

    if (args.length > 0) req._setArgs(args);

    return _sendRequestToProposed(FunctionsRequest._encodeCBOR(req), subscriptionId, callbackGasLimit, donId);
  }

  /**
   * @notice Callback that is invoked once the DON has resolved the request or hit an error
   *
   * @param requestId The request ID, returned by sendRequest()
   * @param response Aggregated response from the user code
   * @param err Aggregated error from the user code or from the execution pipeline
   * Either response or error parameter will be set, but never both
   */
  function _fulfillRequest(bytes32 requestId, bytes memory response, bytes memory err) internal override {
    emit ResponseReceived(requestId, response, err);
  }
}
