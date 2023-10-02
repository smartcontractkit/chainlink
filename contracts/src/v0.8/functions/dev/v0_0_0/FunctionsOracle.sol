// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IFunctionsOracle, IFunctionsBillingRegistry} from "./interfaces/IFunctionsOracle.sol";
import {OCR2BaseUpgradeable} from "./ocr/OCR2BaseUpgradeable.sol";
import {AuthorizedOriginReceiverUpgradeable} from "./accessControl/AuthorizedOriginReceiverUpgradeable.sol";
import {Initializable} from "../../../vendor/openzeppelin-contracts-upgradeable/v4.8.1/proxy/utils/Initializable.sol";

/**
 * @title Functions Oracle contract
 * @notice Contract that nodes of a Decentralized Oracle Network (DON) interact with
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
contract FunctionsOracle is Initializable, IFunctionsOracle, OCR2BaseUpgradeable, AuthorizedOriginReceiverUpgradeable {
  event OracleRequest(
    bytes32 indexed requestId,
    address requestingContract,
    address requestInitiator,
    uint64 subscriptionId,
    address subscriptionOwner,
    bytes data
  );
  event OracleResponse(bytes32 indexed requestId);
  event UserCallbackError(bytes32 indexed requestId, string reason);
  event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData);
  event InvalidRequestID(bytes32 indexed requestId);
  event ResponseTransmitted(bytes32 indexed requestId, address transmitter);

  error EmptyRequestData();
  error InconsistentReportData();
  error EmptyPublicKey();
  error EmptyBillingRegistry();
  error UnauthorizedPublicKeyChange();

  bytes private s_donPublicKey;
  IFunctionsBillingRegistry private s_registry;
  mapping(address => bytes) private s_nodePublicKeys;

  bytes private s_thresholdPublicKey;

  /**
   * @dev Initializes the contract.
   */
  function initialize() public initializer {
    __OCR2Base_initialize(true);
    __AuthorizedOriginReceiver_initialize(true);
  }

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "FunctionsOracle 0.0.0";
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function getRegistry() external view override returns (address) {
    return address(s_registry);
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function setRegistry(address registryAddress) external override onlyOwner {
    if (registryAddress == address(0)) {
      revert EmptyBillingRegistry();
    }
    s_registry = IFunctionsBillingRegistry(registryAddress);
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function getThresholdPublicKey() external view override returns (bytes memory) {
    return s_thresholdPublicKey;
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function setThresholdPublicKey(bytes calldata thresholdPublicKey) external override onlyOwner {
    if (thresholdPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_thresholdPublicKey = thresholdPublicKey;
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_donPublicKey;
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function setDONPublicKey(bytes calldata donPublicKey) external override onlyOwner {
    if (donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_donPublicKey = donPublicKey;
  }

  /**
   * @dev check if node is in current transmitter list
   */
  function _isTransmitter(address node) internal view returns (bool) {
    address[] memory nodes = this.transmitters();
    for (uint256 i = 0; i < nodes.length; i++) {
      if (nodes[i] == node) {
        return true;
      }
    }
    return false;
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function setNodePublicKey(address node, bytes calldata publicKey) external override {
    // Owner can set anything. Transmitters can set only their own key.
    if (!(msg.sender == owner() || (_isTransmitter(msg.sender) && msg.sender == node))) {
      revert UnauthorizedPublicKeyChange();
    }
    s_nodePublicKeys[node] = publicKey;
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function deleteNodePublicKey(address node) external override {
    // Owner can delete anything. Others can delete only their own key.
    if (!(msg.sender == owner() || msg.sender == node)) {
      revert UnauthorizedPublicKeyChange();
    }
    delete s_nodePublicKeys[node];
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function getAllNodePublicKeys() external view override returns (address[] memory, bytes[] memory) {
    address[] memory nodes = this.transmitters();
    bytes[] memory keys = new bytes[](nodes.length);
    for (uint256 i = 0; i < nodes.length; i++) {
      keys[i] = s_nodePublicKeys[nodes[i]];
    }
    return (nodes, keys);
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function getRequiredFee(
    bytes calldata /* data */,
    IFunctionsBillingRegistry.RequestBilling memory /* billing */
  ) public pure override returns (uint96) {
    // NOTE: Optionally, compute additional fee split between nodes of the DON here
    // e.g. 0.1 LINK * s_transmitters.length
    return 0;
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    uint256 gasPrice
  ) external view override registryIsSet returns (uint96) {
    IFunctionsBillingRegistry.RequestBilling memory billing = IFunctionsBillingRegistry.RequestBilling(
      subscriptionId,
      msg.sender,
      gasLimit,
      gasPrice
    );
    uint96 donFee = getRequiredFee(data, billing);
    uint96 registryFee = s_registry.getRequiredFee(data, billing);
    return s_registry.estimateCost(gasLimit, gasPrice, donFee, registryFee);
  }

  /**
   * @inheritdoc IFunctionsOracle
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit
  ) external override registryIsSet validateAuthorizedSender returns (bytes32) {
    if (data.length == 0) {
      revert EmptyRequestData();
    }
    bytes32 requestId = s_registry.startBilling(
      data,
      IFunctionsBillingRegistry.RequestBilling(subscriptionId, msg.sender, gasLimit, tx.gasprice)
    );
    emit OracleRequest(
      requestId,
      msg.sender,
      tx.origin,
      subscriptionId,
      s_registry.getSubscriptionOwner(subscriptionId),
      data
    );
    return requestId;
  }

  function _beforeSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _afterSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _validateReport(
    bytes32 /* configDigest */,
    uint40 /* epochAndRound */,
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
    if (requestIds.length == 0 || requestIds.length != results.length || requestIds.length != errors.length) {
      revert ReportInvalid();
    }

    uint256 reportValidationGasShare = (initialGas - gasleft()) / requestIds.length;

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
      returns (IFunctionsBillingRegistry.FulfillResult result) {
        if (result == IFunctionsBillingRegistry.FulfillResult.USER_SUCCESS) {
          emit OracleResponse(requestIds[i]);
          emit ResponseTransmitted(requestIds[i], transmitter);
        } else if (result == IFunctionsBillingRegistry.FulfillResult.USER_ERROR) {
          emit UserCallbackError(requestIds[i], "error in callback");
          emit ResponseTransmitted(requestIds[i], transmitter);
        } else if (result == IFunctionsBillingRegistry.FulfillResult.INVALID_REQUEST_ID) {
          emit InvalidRequestID(requestIds[i]);
        }
      } catch (bytes memory reason) {
        emit UserCallbackRawError(requestIds[i], reason);
        emit ResponseTransmitted(requestIds[i], transmitter);
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

  /**
   * @dev This empty reserved space is put in place to allow future versions to add new
   * variables without shifting down storage in the inheritance chain.
   * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
   */
  uint256[48] private __gap;
}
