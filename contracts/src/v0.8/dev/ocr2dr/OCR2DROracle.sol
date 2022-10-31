// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";
import "../../ConfirmedOwner.sol";
import "./OCR2Base.sol";
import "./OCR2DRBillingAbstract.sol";

/**
 * @title OCR2DR oracle contract
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
contract OCR2DROracle is OCR2DROracleInterface, OCR2Base, OCR2DRBillingAbstract {
  event OracleRequest(bytes32 requestId, bytes data, RequestBilling billing);
  event OracleResponse(bytes32 requestId);
  event UserCallbackError(bytes32 requestId, string reason);
  event UserCallbackRawError(bytes32 requestId, bytes lowLevelData);

  error EmptyRequestData();
  error InvalidRequestID();
  error InconsistentReportData();
  error EmptyPublicKey();

  struct Commitment {
    address client;
    RequestBilling billing;
  }

  bytes private s_donPublicKey;
  uint256 private s_nonce;
  mapping(bytes32 => Commitment) private s_commitments;

  constructor() OCR2Base(true) {}

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "OCR2DROracle 0.0.0";
  }

  /// @inheritdoc OCR2DROracleInterface
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_donPublicKey;
  }

  /// @inheritdoc OCR2DROracleInterface
  function setDONPublicKey(bytes calldata donPublicKey) external override onlyOwner {
    if (donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_donPublicKey = donPublicKey;
  }

  /// @inheritdoc OCR2DROracleInterface
  function sendRequest(bytes calldata data, RequestBilling calldata billing) external override returns (bytes32) {
    if (data.length == 0) {
      revert EmptyRequestData();
    }
    _preRequestBilling(data, billing);
    s_nonce++;
    bytes32 requestId = keccak256(abi.encodePacked(msg.sender, s_nonce));
    s_commitments[requestId] = Commitment(msg.sender, billing);
    emit OracleRequest(requestId, data, billing);
    return requestId;
  }

  function fulfillRequest(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint32 gasLimit
  ) internal validateRequestId(requestId) returns (uint32) {
    uint32 initialGas = gasleft(); // This line must come first
    OCR2DRClientInterface client = OCR2DRClientInterface(s_commitments[requestId].client);
    delete s_commitments[requestId];
    try client.handleOracleFulfillment.gas(gasLimit)(requestId, response, err) {
      emit OracleResponse(requestId);
    } catch Error(string memory reason) {
      emit UserCallbackError(requestId, reason);
    } catch (bytes memory lowLevelData) {
      emit UserCallbackRawError(requestId, lowLevelData);
    }
    return initialGas - gasleft();
  }

  function _beforeSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _afterSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _validateReport(
    bytes32, /* configDigest */
    uint40, /* epochAndRound */
    bytes memory report
  ) internal override {
    bytes32[] memory requestIds;
    bytes[] memory results;
    bytes[] memory errors;
    (requestIds, results, errors) = abi.decode(report, (bytes32[], bytes[], bytes[]));
    if (requestIds.length != results.length && requestIds.length != errors.length) {
      revert InconsistentReportData();
    }
  }

  function _report(
    uint32 initialGas,
    address transmitter,
    address[maxNumOracles] memory signers,
    bytes calldata report
  ) internal override {
    bytes32[] memory requestIds;
    bytes[] memory results;
    bytes[] memory errors;
    (requestIds, results, errors) = abi.decode(report, (bytes32[], bytes[], bytes[]));

    for (uint256 i = 0; i < requestIds.length; i++) {
      commitment = s_commitments[requestIds[i]];
      uint32 gasUsed = fulfillRequest(requestIds[i], results[i], errors[i], billingConfig.gasLimit);
      _postFulfillBilling(commitment.client, commitment.billing, transmitter, signers, initialGas, gasUsed);
    }
  }

  modifier validateRequestId(bytes32 requestId) {
    if (s_commitments[requestId].client == address(0)) {
      revert InvalidRequestID();
    }
    _;
  }
}
