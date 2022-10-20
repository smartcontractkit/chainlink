// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";
import "../AuthorizedReceiver.sol";
import "../../ConfirmedOwner.sol";

/**
 * @title OCR2DR oracle contract
 */
contract OCR2DROracle is OCR2DROracleInterface, AuthorizedReceiver, ConfirmedOwner {
  event OracleRequest(bytes32 requestId, bytes data);
  event OracleResponse(bytes32 requestId);

  error EmptyRequestData();
  error InvalidRequestID();
  error LowGasForConsumer();

  struct Commitment {
    address client;
    uint256 subscriptionId;
  }

  uint256 private constant MINIMUM_CONSUMER_GAS_LIMIT = 400000;

  bytes private s_donPublicKey;
  uint256 private s_nonce;
  mapping(bytes32 => Commitment) private s_commitments;

  constructor(address owner, bytes memory donPublicKey) ConfirmedOwner(owner) {
    s_donPublicKey = donPublicKey;
  }

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure virtual returns (string memory) {
    return "OCR2DROracle 0.0.0";
  }

  /// @inheritdoc OCR2DROracleInterface
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_donPublicKey;
  }

  /// @inheritdoc OCR2DROracleInterface
  function sendRequest(uint256 subscriptionId, bytes calldata data) external override returns (bytes32) {
    if (data.length == 0) {
      revert EmptyRequestData();
    }
    s_nonce++;
    bytes32 requestId = keccak256(abi.encodePacked(msg.sender, s_nonce));
    s_commitments[requestId] = Commitment(msg.sender, subscriptionId);
    emit OracleRequest(requestId, data);
    return requestId;
  }

  /// @inheritdoc OCR2DROracleInterface
  function fulfillRequest(
    bytes32 requestId,
    bytes calldata response,
    bytes calldata err
  ) external override validateRequestId(requestId) validateAuthorizedSender {
    OCR2DRClientInterface client = OCR2DRClientInterface(s_commitments[requestId].client);
    emit OracleResponse(requestId);
    if (gasleft() < MINIMUM_CONSUMER_GAS_LIMIT) {
      revert LowGasForConsumer();
    }
    client.handleOracleFulfillment(requestId, response, err);
    delete s_commitments[requestId];
  }

  modifier validateRequestId(bytes32 requestId) {
    if (s_commitments[requestId].client == address(0)) {
      revert InvalidRequestID();
    }
    _;
  }

  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return isAuthorizedSender(msg.sender) || owner() == msg.sender;
  }
}
