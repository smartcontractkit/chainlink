// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IDestinationVerifier} from "./interfaces/IDestinationVerifier.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";
import {IAccessController} from "../../shared/interfaces/IAccessController.sol";
import {IDestinationVerifierProxy} from "./interfaces/IDestinationVerifierProxy.sol";
import {IDestinationVerifierProxyVerifier} from "./interfaces/IDestinationVerifierProxyVerifier.sol";
import {IDestinationVerifierFeeManager} from "./interfaces/IDestinationVerifierFeeManager.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/**
 * @title DestinationVerifier
 * @author Michael Fletcher
 * @notice This contract will be used to verify reports based on the oracle signatures. This is not the source verifier which required individual fee configurations, instead, this checks that a report has been signed by one of the configured oracles.
 */
contract DestinationVerifier is
  IDestinationVerifier,
  IDestinationVerifierProxyVerifier,
  ConfirmedOwner,
  ITypeAndVersion
{
  /// @notice The list of DON configurations by hash(address|donConfigId) - set to true if the signer is part of the config
  mapping(bytes32 => bool) private s_signerByAddressAndDonConfigId;

  /// array of DON configs
  DonConfig[] private s_donConfigs;

  /// @notice The address of the verifierProxy
  address public s_feeManager;

  /// @notice The address of the access controller
  address public s_accessController;

  /// @notice The address of the verifierProxy
  IDestinationVerifierProxy public immutable i_verifierProxy;

  /// @notice This error is thrown whenever trying to set a config
  /// with a fault tolerance of 0
  error FaultToleranceMustBePositive();

  /// @notice This error is thrown whenever a report is signed
  /// with more than the max number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param maxSigners The maximum number of signers that can sign a report
  error ExcessSigners(uint256 numSigners, uint256 maxSigners);

  /// @notice This error is thrown whenever a report is signed or expected to be signed with less than the minimum number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param minSigners The minimum number of signers that need to sign a report
  error InsufficientSigners(uint256 numSigners, uint256 minSigners);

  /// @notice This error is thrown whenever a report is submitted with no signatures
  error NoSigners();

  /// @notice This error is thrown whenever a DonConfig already exists
  /// @param donConfigId The ID of the DonConfig that already exists
  error DonConfigAlreadyExists(bytes24 donConfigId);

  /// @notice This error is thrown whenever the R and S signer components
  /// have different lengths
  /// @param rsLength The number of r signature components
  /// @param ssLength The number of s signature components
  error MismatchedSignatures(uint256 rsLength, uint256 ssLength);

  /// @notice This error is thrown whenever setting a config with duplicate signatures
  error NonUniqueSignatures();

  /* @notice This error is thrown whenever a report fails to verify. This error be thrown for multiple reasons and it's purposely like
   * this to prevent information being leaked about the verification process which could be used to enable free verifications maliciously
   */
  error BadVerification();

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown when the fee manager at an address does
  /// not conform to the fee manager interface
  error FeeManagerInvalid();

  /// @notice This error is thrown whenever an address tries
  /// to execute a verification that it is not authorized to do so
  error AccessForbidden();

  /// @notice This error is thrown whenever a config does not exist
  error DonConfigDoesNotExist();

  /// @notice This error is thrown when the activation time is either in the future or less than the current configs
  error BadActivationTime();

  /// @notice This event is emitted when a new report is verified.
  /// It is used to keep a historical record of verified reports.
  event ReportVerified(bytes32 indexed feedId, address requester);

  /// @notice This event is emitted whenever a configuration is activated or deactivated
  event ConfigActivated(bytes24 donConfigId, bool isActive);

  /// @notice This event is emitted whenever a configuration is removed
  event ConfigRemoved(bytes24 donConfigId);

  /// @notice event is emitted whenever a new DON Config is set
  event ConfigSet(
    bytes24 indexed donConfigId,
    address[] signers,
    uint8 f,
    Common.AddressAndWeight[] recipientAddressesAndWeights,
    uint16 donConfigIndex
  );

  /// @notice This event is emitted when a new fee manager is set
  /// @param oldFeeManager The old fee manager address
  /// @param newFeeManager The new fee manager address
  event FeeManagerSet(address oldFeeManager, address newFeeManager);

  /// @notice This event is emitted when a new access controller is set
  /// @param oldAccessController The old access controller address
  /// @param newAccessController The new access controller address
  event AccessControllerSet(address oldAccessController, address newAccessController);

  struct DonConfig {
    // The ID of the DonConfig
    bytes24 donConfigId;
    // Fault tolerance of the DON
    uint8 f;
    // Whether the config is active
    bool isActive;
    // The time the config was set
    uint32 activationTime;
  }

  constructor(address verifierProxy) ConfirmedOwner(msg.sender) {
    if (verifierProxy == address(0)) {
      revert ZeroAddress();
    }

    i_verifierProxy = IDestinationVerifierProxy(verifierProxy);
  }

  /// @inheritdoc IDestinationVerifierProxyVerifier
  function verify(
    bytes calldata signedReport,
    bytes calldata parameterPayload,
    address sender
  ) external payable override onlyProxy checkAccess(sender) returns (bytes memory) {
    (bytes memory verifierResponse, bytes32 donConfigId) = _verify(signedReport, sender);

    address fm = s_feeManager;
    if (fm != address(0)) {
      //process the fee and catch the error
      try
        IDestinationVerifierFeeManager(fm).processFee{value: msg.value}(
          donConfigId,
          signedReport,
          parameterPayload,
          sender
        )
      {
        //do nothing
      } catch {
        // we purposefully obfuscate the error here to prevent information leaking leading to free verifications
        revert BadVerification();
      }
    }

    return verifierResponse;
  }

  /// @inheritdoc IDestinationVerifierProxyVerifier
  function verifyBulk(
    bytes[] calldata signedReports,
    bytes calldata parameterPayload,
    address sender
  ) external payable override onlyProxy checkAccess(sender) returns (bytes[] memory) {
    bytes[] memory verifierResponses = new bytes[](signedReports.length);
    bytes32[] memory donConfigs = new bytes32[](signedReports.length);

    for (uint256 i; i < signedReports.length; ++i) {
      (bytes memory report, bytes32 config) = _verify(signedReports[i], sender);
      verifierResponses[i] = report;
      donConfigs[i] = config;
    }

    address fm = s_feeManager;
    if (fm != address(0)) {
      //process the fee and catch the error
      try
        IDestinationVerifierFeeManager(fm).processFeeBulk{value: msg.value}(
          donConfigs,
          signedReports,
          parameterPayload,
          sender
        )
      {
        //do nothing
      } catch {
        // we purposefully obfuscate the error here to prevent information leaking leading to free verifications
        revert BadVerification();
      }
    }

    return verifierResponses;
  }

  function _verify(bytes calldata signedReport, address sender) internal returns (bytes memory, bytes32) {
    (
      bytes32[3] memory reportContext,
      bytes memory reportData,
      bytes32[] memory rs,
      bytes32[] memory ss,
      bytes32 rawVs
    ) = abi.decode(signedReport, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));

    // Signature lengths must match
    if (rs.length != ss.length) revert MismatchedSignatures(rs.length, ss.length);

    //Must always be at least 1 signer
    if (rs.length == 0) revert NoSigners();

    // The payload is hashed and signed by the oracles - we need to recover the addresses
    bytes32 signedPayload = keccak256(abi.encodePacked(keccak256(reportData), reportContext));
    address[] memory signers = new address[](rs.length);
    for (uint256 i; i < rs.length; ++i) {
      signers[i] = ecrecover(signedPayload, uint8(rawVs[i]) + 27, rs[i], ss[i]);
    }

    // Duplicate signatures are not allowed
    if (Common._hasDuplicateAddresses(signers)) {
      revert BadVerification();
    }

    //We need to know the timestamp the report was generated to lookup the active activeDonConfig
    uint256 reportTimestamp = _decodeReportTimestamp(reportData);

    // Find the latest config for this report
    DonConfig memory activeDonConfig = _findActiveConfig(reportTimestamp);

    // Check a config has been set
    if (activeDonConfig.donConfigId == bytes24(0)) {
      revert BadVerification();
    }

    //check the config is active
    if (!activeDonConfig.isActive) {
      revert BadVerification();
    }

    //check we have enough signatures
    if (signers.length <= activeDonConfig.f) {
      revert BadVerification();
    }

    //check each signer is registered against the active DON
    bytes32 signerDonConfigKey;
    for (uint256 i; i < signers.length; ++i) {
      signerDonConfigKey = keccak256(abi.encodePacked(signers[i], activeDonConfig.donConfigId));
      if (!s_signerByAddressAndDonConfigId[signerDonConfigKey]) {
        revert BadVerification();
      }
    }

    emit ReportVerified(bytes32(reportData), sender);

    return (reportData, activeDonConfig.donConfigId);
  }

  /// @inheritdoc IDestinationVerifier
  function setConfigWithActivationTime(
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights,
    uint32 activationTime
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(signers, f, recipientAddressesAndWeights, activationTime);
  }

  /// @inheritdoc IDestinationVerifier
  function setConfig(
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(signers, f, recipientAddressesAndWeights, uint32(block.timestamp));
  }

  function _setConfig(
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights,
    uint32 activationTime
  ) internal {
    // Duplicate addresses would break protocol rules
    if (Common._hasDuplicateAddresses(signers)) {
      revert NonUniqueSignatures();
    }

    //activation time cannot be in the future
    if (activationTime > block.timestamp) {
      revert BadActivationTime();
    }

    // Sort signers to ensure donConfigId is deterministic
    Common._quickSort(signers, 0, int256(signers.length - 1));

    //DonConfig is made up of hash(signers|f)
    bytes24 donConfigId = bytes24(keccak256(abi.encodePacked(signers, f)));

    // Register the signers for this DON
    for (uint256 i; i < signers.length; ++i) {
      if (signers[i] == address(0)) revert ZeroAddress();
      /** This index is registered so we can efficiently lookup whether a NOP is part of a config without having to
                loop through the entire config each verification. It's effectively a DonConfig <-> Signer
                composite key which keys track of all historic configs for a signer */
      s_signerByAddressAndDonConfigId[keccak256(abi.encodePacked(signers[i], donConfigId))] = true;
    }

    // Check the activation time is greater than the latest config
    uint256 donConfigLength = s_donConfigs.length;
    if (donConfigLength > 0 && s_donConfigs[donConfigLength - 1].activationTime >= activationTime) {
      revert BadActivationTime();
    }

    // Check the config we're setting isn't already set as the current active config as this will increase search costs unnecessarily when verifying historic reports
    if (donConfigLength > 0 && s_donConfigs[donConfigLength - 1].donConfigId == donConfigId) {
      revert DonConfigAlreadyExists(donConfigId);
    }

    // We may want to register these later or skip this step in the unlikely scenario they've previously been registered in the RewardsManager
    if (recipientAddressesAndWeights.length != 0) {
      if (s_feeManager == address(0)) {
        revert FeeManagerInvalid();
      }
      IDestinationVerifierFeeManager(s_feeManager).setFeeRecipients(donConfigId, recipientAddressesAndWeights);
    }

    // push the DonConfig
    s_donConfigs.push(DonConfig(donConfigId, f, true, activationTime));

    emit ConfigSet(donConfigId, signers, f, recipientAddressesAndWeights, uint16(donConfigLength));
  }

  /// @inheritdoc IDestinationVerifier
  function setFeeManager(address feeManager) external override onlyOwner {
    if (!IERC165(feeManager).supportsInterface(type(IDestinationVerifierFeeManager).interfaceId))
      revert FeeManagerInvalid();

    address oldFeeManager = s_feeManager;
    s_feeManager = feeManager;

    emit FeeManagerSet(oldFeeManager, feeManager);
  }

  /// @inheritdoc IDestinationVerifier
  function setAccessController(address accessController) external override onlyOwner {
    address oldAccessController = s_accessController;
    s_accessController = accessController;
    emit AccessControllerSet(oldAccessController, accessController);
  }

  /// @inheritdoc IDestinationVerifier
  function setConfigActive(uint256 donConfigIndex, bool isActive) external onlyOwner {
    // Config must exist
    if (donConfigIndex >= s_donConfigs.length) {
      revert DonConfigDoesNotExist();
    }

    // Update the config
    DonConfig storage config = s_donConfigs[donConfigIndex];
    config.isActive = isActive;

    emit ConfigActivated(config.donConfigId, isActive);
  }

  /// @inheritdoc IDestinationVerifier
  function removeLatestConfig() external onlyOwner {
    if (s_donConfigs.length == 0) {
      revert DonConfigDoesNotExist();
    }

    DonConfig memory config = s_donConfigs[s_donConfigs.length - 1];

    s_donConfigs.pop();

    emit ConfigRemoved(config.donConfigId);
  }

  function _decodeReportTimestamp(bytes memory reportPayload) internal pure returns (uint256) {
    (, , uint256 timestamp) = abi.decode(reportPayload, (bytes32, uint32, uint32));

    return timestamp;
  }

  function _findActiveConfig(uint256 timestamp) internal view returns (DonConfig memory) {
    DonConfig memory activeDonConfig;

    // 99% of the time the signer config will be the last index, however for historic reports generated by a previous configuration we'll need to cycle back
    uint256 i = s_donConfigs.length;
    while (i > 0) {
      --i;
      if (s_donConfigs[i].activationTime <= timestamp) {
        activeDonConfig = s_donConfigs[i];
        break;
      }
    }
    return activeDonConfig;
  }

  modifier checkConfigValid(uint256 numSigners, uint256 f) {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
    if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
    _;
  }

  modifier onlyProxy() {
    if (address(i_verifierProxy) != msg.sender) {
      revert AccessForbidden();
    }
    _;
  }

  modifier checkAccess(address sender) {
    address ac = s_accessController;
    if (address(ac) != address(0) && !IAccessController(ac).hasAccess(sender, msg.data)) revert AccessForbidden();
    _;
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return
      interfaceId == type(IDestinationVerifier).interfaceId ||
      interfaceId == type(IDestinationVerifierProxyVerifier).interfaceId;
  }

  /// @inheritdoc ITypeAndVersion
  function typeAndVersion() external pure override returns (string memory) {
    return "DestinationVerifier 0.4.0";
  }
}
