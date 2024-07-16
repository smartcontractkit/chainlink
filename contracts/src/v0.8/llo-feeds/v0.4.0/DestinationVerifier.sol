// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IDestinationVerifier} from "./interfaces/IDestinationVerifier.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";
import {IAccessController} from "../../shared/interfaces/IAccessController.sol";
import {IDestinationVerifierProxy} from "./interfaces/IDestinationVerifierProxy.sol";
import {IDestinationFeeManager} from "./interfaces/IDestinationFeeManager.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/**
 * @title DestinationVerifier
 * @author Michael Fletcher
 * @notice This contract will be used to verify reports based on the oracle signatures
 */
contract DestinationVerifier is IDestinationVerifier, ConfirmedOwner, TypeAndVersionInterface, IERC165 {

    /// @notice The list of DON configurations by hash(address|DONConfigID) - set to true if the signer is part of the config
    mapping(bytes32 => bool) private s_SignerByAddressAndDONConfigId;

    /// array of DON configs
    DONConfig[] private s_DONConfigs;

    /// @notice The address of the verifierProxy
    IDestinationFeeManager private s_feeManager;

    /// @notice The address of the access controller
    IAccessController private s_accessController;

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

    /// @notice This error is thrown whenever a DONConfig already exists
    /// @param DONConfigID The ID of the DONConfig that already exists
    error DONConfigAlreadyExists(bytes24 DONConfigID);

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
    error DONConfigDoesNotExist();

    /// @notice this error is thrown when the verifierProxy is incorrect when initialising
    error VerifierProxyInvalid();

    /// @notice This event is emitted when a new report is verified.
    /// It is used to keep a historical record of verified reports.
    event ReportVerified(bytes32 indexed feedId, address requester);

    /// @notice This event is emitted whenever a configuration is activated or deactivated
    event ConfigActivated(bytes24 DONConfigID, bool isActive);

    /// @notice event is emitted whenever a new DON Config is set
    event ConfigSet(bytes24 indexed DONConfigID, address[] signers, uint8 f, Common.AddressAndWeight[] recipientAddressesAndWeights);

    /// @notice This event is emitted when a new fee manager is set
    /// @param oldFeeManager The old fee manager address
    /// @param newFeeManager The new fee manager address
    event FeeManagerSet(address oldFeeManager, address newFeeManager);

    /// @notice This event is emitted when a new access controller is set
    /// @param oldAccessController The old access controller address
    /// @param newAccessController The new access controller address
    event AccessControllerSet(address oldAccessController, address newAccessController);

    struct DONConfig {
        // The ID of the DONConfig
        bytes24 DONConfigID;
        // Fault tolerance of the DON
        uint8 f;
        // Whether the config is active
        bool isActive;
        // The time the config was set
        uint32 activationTime;
    }

    constructor(address verifierProxy) ConfirmedOwner(msg.sender) {
        if(verifierProxy == address(0)) {
            revert ZeroAddress();
        }

        //check it supports the interface
        // ask how to mock this
    //    if(!IERC165(verifierProxy).supportsInterface(type(IDestinationVerifierProxy).interfaceId)) {
      //      revert VerifierProxyInvalid();
       // }

        i_verifierProxy = IDestinationVerifierProxy(verifierProxy);
    }

    /// @inheritdoc IDestinationVerifier
    function verify(
        bytes calldata signedReport,
        bytes calldata parameterPayload,
        address sender
    ) external override checkValidProxy checkAccess(sender) payable returns (bytes memory) {
        (bytes memory verifierResponse, bytes32 DONConfigId) = _verify(signedReport, sender);

        if(address(s_feeManager) != address(0)){
            //process the fee and catch the error
            try s_feeManager.processFee(DONConfigId, signedReport, parameterPayload, sender) {
                //do nothing
            } catch {
                // we purposefully obfuscate the error here to prevent information leaking leading to free verifications
                revert BadVerification();
            }
        }

        return verifierResponse;
    }

    /// @inheritdoc IDestinationVerifier
    function verifyBulk(
        bytes[] calldata signedReports,
        bytes calldata parameterPayload,
        address sender
    ) external override checkValidProxy checkAccess(sender) payable returns (bytes[] memory) {
        bytes[] memory verifierResponses = new bytes[](signedReports.length);
        bytes32[] memory DONConfigs = new bytes32[](signedReports.length);

        for(uint i; i < signedReports.length; ++i)  {
            (bytes memory report, bytes32 config) = _verify(signedReports[i], sender);
            verifierResponses[i] = report;
            DONConfigs[i] = config;
        }

        if(address(s_feeManager) != address(0)){
            //process the fee and catch the error
            try s_feeManager.processFeeBulk(DONConfigs, signedReports, parameterPayload, sender) {
                //do nothing
            } catch {
                // we purposefully obfuscate the error here to prevent information leaking leading to free verifications
                revert BadVerification();
            }
        }

        return verifierResponses;
    }

    function _verify(
        bytes calldata signedReport,
        address sender
    ) internal returns (bytes memory, bytes32)  {
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
        if(rs.length == 0) revert NoSigners();

        // The payload is hashed and signed by the oracles - we need to recover the addresses
        bytes32 signedPayload = keccak256(abi.encodePacked(keccak256(reportData), reportContext));
        address[] memory signers = new address[](rs.length);
        for(uint i; i < rs.length; ++i) {
            signers[i] = ecrecover(signedPayload, uint8(rawVs[i]) + 27, rs[i], ss[i]);
        }

        // Duplicate signatures are not allowed
        if(Common._hasDuplicateAddresses(signers)) {
            revert BadVerification();
        }

        //We need to know the timestamp the report was generated to lookup the active activeDONConfig
        uint256 reportTimestamp = decodeReportTimestamp(reportData);

        // Find the latest config for this report
        DONConfig memory activeDONConfig = findActiveConfig(reportTimestamp);

        // Check a config has been set
        if(activeDONConfig.DONConfigID == bytes24(0)) {
            revert BadVerification();
        }

        //check the config is active
        if(!activeDONConfig.isActive) {
            revert BadVerification();
        }

        //check we have enough signatures
        if(signers.length <= activeDONConfig.f) {
            revert BadVerification();
        }

        //check each signer is registered against the active DON
        bytes32 signerDONConfigKey;
        for(uint i; i < signers.length; ++i) {
            signerDONConfigKey = keccak256(abi.encodePacked(signers[i], activeDONConfig.DONConfigID));
            if(!s_SignerByAddressAndDONConfigId[signerDONConfigKey]) {
                revert BadVerification();
            }
        }

        emit ReportVerified(bytes32(reportData), sender);

        return (reportData, activeDONConfig.DONConfigID);
    }

    /// @inheritdoc IDestinationVerifier
    function setConfig(
        address[] memory signers,
        uint8 f,
        Common.AddressAndWeight[] memory recipientAddressesAndWeights
    ) external override checkConfigValid(signers.length, f) onlyOwner {
        // Duplicate addresses would break protocol rules
        if(Common._hasDuplicateAddresses(signers)) {
            revert NonUniqueSignatures();
        }

        // Sort signers to ensure DONConfigID is deterministic
        Common._quickSort(signers, 0, int256(signers.length - 1));

        //DONConfig is made up of hash(signers|f)
        bytes24 DONConfigID = bytes24(keccak256(abi.encodePacked(signers, f)));

        // Register the signers for this DON
        for(uint i; i < signers.length; ++i) {
            if(signers[i] == address(0))
                revert ZeroAddress();
            /** This index is registered so we can efficiently lookup whether a NOP is part of a config without having to
                loop through the entire config each verification. It's effectively a DONConfig <-> Signer
                composite key which keys track of all historic configs for a signer */
            s_SignerByAddressAndDONConfigId[keccak256(abi.encodePacked(signers[i], DONConfigID))] = true;
        }

        // Check the config we're setting isn't already set as the current active config as this will increase search costs unnecessarily when verifying historic reports
        if(s_DONConfigs.length > 0 && s_DONConfigs[s_DONConfigs.length - 1].DONConfigID == DONConfigID) {
            revert DONConfigAlreadyExists(DONConfigID);
        }

        // We may want to register these later or skip this step in the unlikely scenario they've previously been registered in the RewardsManager
        if(recipientAddressesAndWeights.length != 0) {
          s_feeManager.setFeeRecipients(DONConfigID, recipientAddressesAndWeights);
        }

        // push the DONConfig
        s_DONConfigs.push(DONConfig(DONConfigID, f, true, uint32(block.timestamp)));

        emit ConfigSet(DONConfigID, signers, f, recipientAddressesAndWeights);
    }

    /// @inheritdoc IDestinationVerifier
    function setFeeManager(address feeManager) external override onlyOwner {
        if (!IERC165(feeManager).supportsInterface(IDestinationFeeManager.processFee.selector) ||
            !IERC165(feeManager).supportsInterface(IDestinationFeeManager.processFeeBulk.selector) ||
            !IERC165(feeManager).supportsInterface(IDestinationFeeManager.setFeeRecipients.selector)) revert FeeManagerInvalid();

        address oldFeeManager = address(s_feeManager);
        s_feeManager = IDestinationFeeManager(feeManager);

        emit FeeManagerSet(oldFeeManager, feeManager);
    }

    /// @inheritdoc IDestinationVerifier
    function setAccessController(address accessController) external override onlyOwner {
        address oldAccessController = address(s_accessController);
        s_accessController = IAccessController(accessController);
        emit AccessControllerSet(oldAccessController, accessController);
    }

    /// @inheritdoc IDestinationVerifier
    function setConfigActive(uint256 DONConfigIndex, bool isActive) external onlyOwner {
        // Config must exist
        if(DONConfigIndex >= s_DONConfigs.length) {
            revert DONConfigDoesNotExist();
        }

        // Update the config
        DONConfig storage config = s_DONConfigs[DONConfigIndex];
        config.isActive = isActive;

        emit ConfigActivated(config.DONConfigID, isActive);
    }

    function decodeReportTimestamp(bytes memory reportPayload) internal pure returns (uint256) {
        (,,uint256 timestamp) = abi.decode(reportPayload, (bytes32, uint32, uint32));

        return timestamp;
    }

    function findActiveConfig(uint256 timestamp) internal view returns (DONConfig memory){
        DONConfig memory activeDONConfig;

        // 99% of the time the signer config will be the last index, however for historic reports generated by a previous configuration we'll need to cycle back
        for(uint256 i = s_DONConfigs.length - 1; i >= 0; --i) {
            activeDONConfig = s_DONConfigs[i];
            if(activeDONConfig.activationTime <= timestamp) {
                break;
            }
        }
        return activeDONConfig;
    }

    modifier checkConfigValid(uint256 numSigners, uint256 f){
        if(f == 0) revert FaultToleranceMustBePositive();
        if(numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
        if(numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
        _;
    }

    modifier checkValidProxy() {
        if(address(i_verifierProxy) != msg.sender) {
            revert AccessForbidden();
        }
        _;
    }

    modifier checkAccess(address sender) {
        IAccessController ac = s_accessController;
        if (address(ac) != address(0) && !ac.hasAccess(sender, msg.data)) revert AccessForbidden();
        _;
    }

    /// @inheritdoc IERC165
    function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
       return interfaceId == this.verify.selector ||
              interfaceId == this.verifyBulk.selector ||
              interfaceId == this.getAccessController.selector ||
              interfaceId == this.getFeeManager.selector ||
              interfaceId == this.setConfig.selector ||
              interfaceId == this.setFeeManager.selector ||
              interfaceId == this.setAccessController.selector ||
              interfaceId == this.setConfigActive.selector;
    }

    /// @inheritdoc TypeAndVersionInterface
    function typeAndVersion() external pure override returns (string memory) {
        return "DestinationVerifier 1.0.0";
    }

    /// @inheritdoc IDestinationVerifier
    function getAccessController() external view override returns (address) {
        return address(s_accessController);
    }

    /// @inheritdoc IDestinationVerifier
    function getFeeManager() external view override returns (address) {
        return address(s_feeManager);
    }

    function getDONConfigs() public view returns (DONConfig[] memory) {
        return s_DONConfigs;
    }

     function getSignerConfigByAddressAndDONConfigId(bytes32 _id) public view returns (bool) {
        return s_SignerByAddressAndDONConfigId[_id];
    }
}

