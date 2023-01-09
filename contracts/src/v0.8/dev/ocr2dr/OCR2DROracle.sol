// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DROracleInterface.sol";
import "../ocr2/OCR2Base.sol";
import "./AuthorizedOriginReceiver.sol";

/**
 * @title OCR2DR oracle contract
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
contract OCR2DROracle is OCR2DROracleInterface, OCR2Base, AuthorizedOriginReceiver {
  event OracleRequest(bytes32 indexed requestId, uint64 subscriptionId, bytes data);
  event OracleResponse(bytes32 indexed requestId);
  event UserCallbackError(bytes32 indexed requestId, string reason);
  event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData);

  error EmptyRequestData();
  error InconsistentReportData();
  error EmptyPublicKey();
  error EmptyBillingRegistry();
  error InvalidRequestID();

  bytes private s_donPublicKey;
  OCR2DRRegistryInterface private s_registry;

  constructor() OCR2Base(true) {}

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "OCR2DROracle 0.0.0";
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function getRegistry() external view override returns (address) {
    return address(s_registry);
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function setRegistry(address registryAddress) external override onlyOwner {
    if (registryAddress == address(0)) {
      revert EmptyBillingRegistry();
    }
    s_registry = OCR2DRRegistryInterface(registryAddress);
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_donPublicKey;
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function setDONPublicKey(bytes calldata donPublicKey) external override onlyOwner {
    if (donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_donPublicKey = donPublicKey;
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function getRequiredFee(
    bytes calldata, /* data */
    OCR2DRRegistryInterface.RequestBilling memory /* billing */
  ) public pure override returns (uint96) {
    // NOTE: Optionally, compute additional fee split between oracles here
    // e.g. 0.1 LINK * s_transmitters.length
    return 0;
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    uint256 gasPrice
  ) external view override registryIsSet returns (uint96) {
    OCR2DRRegistryInterface.RequestBilling memory billing = OCR2DRRegistryInterface.RequestBilling(
      subscriptionId,
      msg.sender,
      gasLimit,
      gasPrice
    );
    uint96 requiredFee = getRequiredFee(data, billing);
    uint96 registryFee = getRequiredFee(data, billing);
    return s_registry.estimateCost(gasLimit, gasPrice, requiredFee, registryFee);
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    uint256 gasPrice
  ) external override registryIsSet validateAuthorizedSender returns (bytes32) {
    if (data.length == 0) {
      revert EmptyRequestData();
    }
    bytes32 requestId = s_registry.startBilling(
      data,
      OCR2DRRegistryInterface.RequestBilling(subscriptionId, msg.sender, gasLimit, gasPrice)
    );
    emit OracleRequest(requestId, subscriptionId, data);
    return requestId;
  }

  function _beforeSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _afterSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _validateReport(
    bytes32, /* configDigest */
    uint40, /* epochAndRound */
    bytes memory /* report */
  ) internal pure override returns (bool) {
    // validate within _report to save gas
    return true;
  }

  function _report(
    uint256 initialGas,
    address transmitter,
    uint8 signerCount,
    address[maxNumOracles] memory signers,
    bytes calldata report
  ) internal override registryIsSet {
    bytes32[] memory requestIds;
    bytes[] memory results;
    bytes[] memory errors;
    (requestIds, results, errors) = abi.decode(report, (bytes32[], bytes[], bytes[]));
    if (requestIds.length != results.length && requestIds.length != errors.length) {
      revert ReportInvalid();
    }

    uint256 reportValidationGasShare = (initialGas - gasleft()) / signerCount;

    for (uint256 i = 0; i < requestIds.length; i++) {
      try
        s_registry.fulfillAndBill(
          requestIds[i],
          results[i],
          errors[i],
          transmitter,
          signers,
          signerCount,
          reportValidationGasShare,
          gasleft()
        )
      returns (bool success) {
        if (success) {
          emit OracleResponse(requestIds[i]);
        } else {
          emit UserCallbackError(requestIds[i], "error in callback");
        }
      } catch (bytes memory reason) {
        emit UserCallbackRawError(requestIds[i], reason);
      }
    }
  }

  /**
   * @dev Reverts if the the registry is not set
   */
  modifier registryIsSet() {
    if (address(s_registry) == address(0)) {
      revert EmptyBillingRegistry();
    }
    _;
  }

  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return msg.sender == owner();
  }
}
