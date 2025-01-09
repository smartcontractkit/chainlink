// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IVerifierProxy} from "./interfaces/IVerifierProxy.sol";
import {IVerifier} from "./interfaces/IVerifier.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IVerifierFeeManager} from "./interfaces/IVerifierFeeManager.sol";
import {Common} from "../libraries/Common.sol";

/**
 * The verifier proxy contract is the gateway for all report verification requests
 * on a chain.  It is responsible for taking in a verification request and routing
 * it to the correct verifier contract.
 */
contract VerifierProxy is IVerifierProxy, ConfirmedOwner, ITypeAndVersion {
  /// @notice This event is emitted whenever a new verifier contract is set
  /// @param oldConfigDigest The config digest that was previously the latest config
  /// digest of the verifier contract at the verifier address.
  /// @param oldConfigDigest The latest config digest of the verifier contract
  /// at the verifier address.
  /// @param verifierAddress The address of the verifier contract that verifies reports for
  /// a given digest
  event VerifierSet(bytes32 oldConfigDigest, bytes32 newConfigDigest, address verifierAddress);

  /// @notice This event is emitted whenever a new verifier contract is initialized
  /// @param verifierAddress The address of the verifier contract that verifies reports
  event VerifierInitialized(address verifierAddress);

  /// @notice This event is emitted whenever a verifier is unset
  /// @param configDigest The config digest that was unset
  /// @param verifierAddress The Verifier contract address unset
  event VerifierUnset(bytes32 configDigest, address verifierAddress);

  /// @notice This event is emitted when a new access controller is set
  /// @param oldAccessController The old access controller address
  /// @param newAccessController The new access controller address
  event AccessControllerSet(address oldAccessController, address newAccessController);

  /// @notice This event is emitted when a new fee manager is set
  /// @param oldFeeManager The old fee manager address
  /// @param newFeeManager The new fee manager address
  event FeeManagerSet(address oldFeeManager, address newFeeManager);

  /// @notice This error is thrown whenever an address tries
  /// to exeecute a transaction that it is not authorized to do so
  error AccessForbidden();

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown when trying to set a verifier address
  /// for a digest that has already been initialized
  /// @param configDigest The digest for the verifier that has
  /// already been set
  /// @param verifier The address of the verifier the digest was set for
  error ConfigDigestAlreadySet(bytes32 configDigest, address verifier);

  /// @notice This error is thrown when trying to set a verifier address that has already been initialized
  error VerifierAlreadyInitialized(address verifier);

  /// @notice This error is thrown when the verifier at an address does
  /// not conform to the verifier interface
  error VerifierInvalid();

  /// @notice This error is thrown when the fee manager at an address does
  /// not conform to the fee manager interface
  error FeeManagerInvalid();

  /// @notice This error is thrown whenever a verifier is not found
  /// @param configDigest The digest for which a verifier is not found
  error VerifierNotFound(bytes32 configDigest);

  /// @notice This error is thrown whenever billing fails.
  error BadVerification();

  /// @notice Mapping of authorized verifiers
  mapping(address => bool) private s_initializedVerifiers;

  /// @notice Mapping between config digests and verifiers
  mapping(bytes32 => address) private s_verifiersByConfig;

  /// @notice The contract to control addresses that are allowed to verify reports
  AccessControllerInterface public s_accessController;

  /// @notice The contract to control fees for report verification
  IVerifierFeeManager public s_feeManager;

  constructor(AccessControllerInterface accessController) ConfirmedOwner(msg.sender) {
    s_accessController = accessController;
  }

  modifier checkAccess() {
    AccessControllerInterface ac = s_accessController;
    if (address(ac) != address(0) && !ac.hasAccess(msg.sender, msg.data)) revert AccessForbidden();
    _;
  }

  modifier onlyInitializedVerifier() {
    if (!s_initializedVerifiers[msg.sender]) revert AccessForbidden();
    _;
  }

  modifier onlyValidVerifier(address verifierAddress) {
    if (verifierAddress == address(0)) revert ZeroAddress();
    if (!IERC165(verifierAddress).supportsInterface(IVerifier.verify.selector)) revert VerifierInvalid();
    _;
  }

  modifier onlyUnsetConfigDigest(bytes32 configDigest) {
    address configDigestVerifier = s_verifiersByConfig[configDigest];
    if (configDigestVerifier != address(0)) revert ConfigDigestAlreadySet(configDigest, configDigestVerifier);
    _;
  }

  /// @inheritdoc ITypeAndVersion
  function typeAndVersion() external pure override returns (string memory) {
    return "VerifierProxy 2.0.0";
  }

  /// @inheritdoc IVerifierProxy
  function verify(
    bytes calldata payload,
    bytes calldata parameterPayload
  ) external payable checkAccess returns (bytes memory) {
    IVerifierFeeManager feeManager = s_feeManager;

    // Bill the verifier
    if (address(feeManager) != address(0)) {
      feeManager.processFee{value: msg.value}(payload, parameterPayload, msg.sender);
    }

    return _verify(payload);
  }

  /// @inheritdoc IVerifierProxy
  function verifyBulk(
    bytes[] calldata payloads,
    bytes calldata parameterPayload
  ) external payable checkAccess returns (bytes[] memory verifiedReports) {
    IVerifierFeeManager feeManager = s_feeManager;

    // Bill the verifier
    if (address(feeManager) != address(0)) {
      feeManager.processFeeBulk{value: msg.value}(payloads, parameterPayload, msg.sender);
    }

    //verify the reports
    verifiedReports = new bytes[](payloads.length);
    for (uint256 i; i < payloads.length; ++i) {
      verifiedReports[i] = _verify(payloads[i]);
    }

    return verifiedReports;
  }

  function _verify(bytes calldata payload) internal returns (bytes memory verifiedReport) {
    // First 32 bytes of the signed report is the config digest
    bytes32 configDigest = bytes32(payload);
    address verifierAddress = s_verifiersByConfig[configDigest];
    if (verifierAddress == address(0)) revert VerifierNotFound(configDigest);

    return IVerifier(verifierAddress).verify(payload, msg.sender);
  }

  /// @inheritdoc IVerifierProxy
  function initializeVerifier(address verifierAddress) external override onlyOwner onlyValidVerifier(verifierAddress) {
    if (s_initializedVerifiers[verifierAddress]) revert VerifierAlreadyInitialized(verifierAddress);

    s_initializedVerifiers[verifierAddress] = true;
    emit VerifierInitialized(verifierAddress);
  }

  /// @inheritdoc IVerifierProxy
  function setVerifier(
    bytes32 currentConfigDigest,
    bytes32 newConfigDigest,
    Common.AddressAndWeight[] calldata addressesAndWeights
  ) external override onlyUnsetConfigDigest(newConfigDigest) onlyInitializedVerifier {
    s_verifiersByConfig[newConfigDigest] = msg.sender;

    // Empty recipients array will be ignored and must be set off chain
    if (addressesAndWeights.length > 0) {
      if (address(s_feeManager) == address(0)) {
        revert ZeroAddress();
      }

      s_feeManager.setFeeRecipients(newConfigDigest, addressesAndWeights);
    }

    emit VerifierSet(currentConfigDigest, newConfigDigest, msg.sender);
  }

  /// @inheritdoc IVerifierProxy
  function unsetVerifier(bytes32 configDigest) external override onlyOwner {
    address verifierAddress = s_verifiersByConfig[configDigest];
    if (verifierAddress == address(0)) revert VerifierNotFound(configDigest);
    delete s_verifiersByConfig[configDigest];
    emit VerifierUnset(configDigest, verifierAddress);
  }

  /// @inheritdoc IVerifierProxy
  function getVerifier(bytes32 configDigest) external view override returns (address) {
    return s_verifiersByConfig[configDigest];
  }

  /// @inheritdoc IVerifierProxy
  function setAccessController(AccessControllerInterface accessController) external onlyOwner {
    address oldAccessController = address(s_accessController);
    s_accessController = accessController;
    emit AccessControllerSet(oldAccessController, address(accessController));
  }

  /// @inheritdoc IVerifierProxy
  function setFeeManager(IVerifierFeeManager feeManager) external onlyOwner {
    if (address(feeManager) == address(0)) revert ZeroAddress();

    if (
      !IERC165(feeManager).supportsInterface(IVerifierFeeManager.processFee.selector) ||
      !IERC165(feeManager).supportsInterface(IVerifierFeeManager.processFeeBulk.selector)
    ) revert FeeManagerInvalid();

    address oldFeeManager = address(s_feeManager);
    s_feeManager = IVerifierFeeManager(feeManager);
    emit FeeManagerSet(oldFeeManager, address(feeManager));
  }
}
