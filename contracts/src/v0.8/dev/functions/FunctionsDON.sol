// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/FunctionsDONInterface.sol";
import "../ocr2/OCR2Base.sol";
import "./AuthorizedOriginReceiver.sol";

/**
 * @title Functions Decentralized Oracle Network (DON) contract
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
contract FunctionsDON is FunctionsDONInterface, OCR2Base, AuthorizedOriginReceiver {
  event DONRequest(bytes32 indexed requestId, uint64 subscriptionId, bytes data);
  event DONResponse(bytes32 indexed requestId);
  event UserCallbackError(bytes32 indexed requestId, string reason);
  event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData);

  error EmptyRequestData();
  error InconsistentReportData();
  error EmptyPublicKey();
  error EmptyBillingRegistry();
  error InvalidRequestID();

  bytes private s_donPublicKey;
  FunctionsBillingRegistryInterface private s_registry;

  constructor() OCR2Base(true) {}

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "FunctionsDON 0.0.0";
  }

  /**
   * @inheritdoc FunctionsDONInterface
   */
  function getRegistry() external view override returns (address) {
    return address(s_registry);
  }

  /**
   * @inheritdoc FunctionsDONInterface
   */
  function setRegistry(address registryAddress) external override onlyOwner {
    if (registryAddress == address(0)) {
      revert EmptyBillingRegistry();
    }
    s_registry = FunctionsBillingRegistryInterface(registryAddress);
  }

  /**
   * @inheritdoc FunctionsDONInterface
   */
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_donPublicKey;
  }

  /**
   * @inheritdoc FunctionsDONInterface
   */
  function setDONPublicKey(bytes calldata donPublicKey) external override onlyOwner {
    if (donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_donPublicKey = donPublicKey;
  }

  /**
   * @inheritdoc FunctionsDONInterface
   */
  function getRequiredFee(
    bytes calldata, /* data */
    FunctionsBillingRegistryInterface.RequestBilling memory /* billing */
  ) public pure override returns (uint96) {
    // NOTE: Optionally, compute additional fee split between nodes of the DON here
    // e.g. 0.1 LINK * s_transmitters.length
    return 0;
  }

  /**
   * @inheritdoc FunctionsDONInterface
   */
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    uint256 gasPrice
  ) external view override registryIsSet returns (uint96) {
    FunctionsBillingRegistryInterface.RequestBilling memory billing = FunctionsBillingRegistryInterface.RequestBilling(
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
   * @inheritdoc FunctionsDONInterface
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
      FunctionsBillingRegistryInterface.RequestBilling(subscriptionId, msg.sender, gasLimit, gasPrice)
    );
    emit DONRequest(requestId, subscriptionId, data);
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
          emit DONResponse(requestIds[i]);
        } else {
          emit UserCallbackError(requestIds[i], "error in callback");
        }
      } catch (bytes memory reason) {
        emit UserCallbackRawError(requestIds[i], reason);
      }
    }
  }

  /**
   * @dev Reverts if the the billing registry is not set
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
