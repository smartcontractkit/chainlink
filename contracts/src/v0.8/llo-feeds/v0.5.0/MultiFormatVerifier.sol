// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IVerifier} from "./interfaces/IVerifier.sol";
import {IVerifierProxy} from "./interfaces/IVerifierProxy.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";

uint32 constant REPORT_FORMAT_EVM_PREMIUM_LEGACY = 1;
uint32 constant REPORT_FORMAT_EVM_ABI_ENCODE_UNPACKED = 3;
uint32 constant REPORT_FORMAT_EVM_STREAMLINED = 6;

uint32[] constant SUPPORTED_FORMATS = [
    REPORT_FORMAT_EVM_PREMIUM_LEGACY,
    REPORT_FORMAT_EVM_ABI_ENCODE_UNPACKED,
    REPORT_FORMAT_EVM_STREAMLINED
];

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/*
 **/
contract Verifier is IVerifier, ConfirmedOwner, ITypeAndVersion {
    // The first byte of the mask can be 0, because we only ever have 31 oracles
    uint256 internal constant ORACLE_MASK = 0x0001010101010101010101010101010101010101010101010101010101010101;

    enum Role {
        // Default role for an oracle address.  This means that the oracle address
        // is not a signer
        Unset,
        // Role given to an oracle address that is allowed to sign a report
        Signer
    }

    struct Signer {
        // Index of oracle in a configuration
        uint8 index;
        // The oracle's role
        Role role;
    }

    struct VerifierState {
        // The block number of the block the last time the configuration was updated.
        uint32 latestConfigBlockNumber;
        // Whether the config is deactivated
        bool isActive;
        // Fault tolerance
        uint8 f;
        // Number of signers
        uint8 oracleCount;
        // Map of signer addresses to oracles
        mapping(address => Signer) oracles;
    }

    /// @notice This event is emitted when a new report is verified.
    /// It is used to keep a historical record of verified reports.
    event ReportVerified(bytes32 indexed feedId, address requester);

    /// @notice This event is emitted whenever a new DON configuration is set.
    event ConfigSet(bytes32 indexed configDigest, address[] signers, uint8 f);

    /// @notice This event is
    event ConfigUpdated(bytes32 indexed configDigest, address[] prevSigners, address[] newSigners);

    /// @notice This event is emitted whenever a configuration is deactivated
    event ConfigDeactivated(bytes32 indexed configDigest);

    /// @notice This event is emitted whenever a configuration is activated
    event ConfigActivated(bytes32 indexed configDigest);

    /// @notice This error is thrown whenever an address tries
    /// to exeecute a transaction that it is not authorized to do so
    error AccessForbidden();

    /// @notice This error is thrown whenever a zero address is passed
    error ZeroAddress();

    /// @notice This error is thrown whenever the config digest
    /// is empty
    error DigestEmpty();

    /// @notice This error is thrown whenever the config digest
    /// passed in has not been set in this verifier
    /// @param configDigest The config digest that has not been set
    error DigestNotSet(bytes32 configDigest);

    /// @notice This error is thrown whenever the config digest
    /// has been deactivated
    /// @param configDigest The config digest that is inactive
    error DigestInactive(bytes32 configDigest);

    /// @notice This error is thrown whenever trying to set a config
    /// with a fault tolerance of 0
    error FaultToleranceMustBePositive();

    /// @notice This error is thrown whenever a report is signed
    /// with more than the max number of signers
    /// @param numSigners The number of signers who have signed the report
    /// @param maxSigners The maximum number of signers that can sign a report
    error ExcessSigners(uint256 numSigners, uint256 maxSigners);

    /// @notice This error is thrown whenever a report is signed
    /// with less than the minimum number of signers
    /// @param numSigners The number of signers who have signed the report
    /// @param minSigners The minimum number of signers that need to sign a report
    error InsufficientSigners(uint256 numSigners, uint256 minSigners);

    /// @notice This error is thrown whenever a report is signed
    /// with an incorrect number of signers
    /// @param numSigners The number of signers who have signed the report
    /// @param expectedNumSigners The expected number of signers that need to sign
    /// a report
    error IncorrectSignatureCount(uint256 numSigners, uint256 expectedNumSigners);

    /// @notice This error is thrown whenever the R and S signer components
    /// have different lengths
    /// @param rsLength The number of r signature components
    /// @param ssLength The number of s signature components
    error MismatchedSignatures(uint256 rsLength, uint256 ssLength);

    /// @notice This error is thrown whenever setting a config with duplicate signatures
    error NonUniqueSignatures();

    /// @notice This error is thrown whenever a report fails to verify due to bad or duplicate signatures
    error BadVerification();

    /// @notice This error is thrown whenever a config digest is already set when setting the configuration
    error ConfigDigestAlreadySet();

    /// @notice The address of the verifier proxy
    address private immutable i_verifierProxyAddr;

    /// @notice Verifier states keyed on config digest
    mapping(bytes32 => VerifierState) internal s_verifierStates;

    /// @param verifierProxyAddr The address of the VerifierProxy contract
    constructor(address verifierProxyAddr) ConfirmedOwner(msg.sender) {
        if (verifierProxyAddr == address(0)) revert ZeroAddress();
        i_verifierProxyAddr = verifierProxyAddr;
    }

    modifier checkConfigValid(uint256 numSigners, uint256 f) {
        if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
        if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
        _;
    }

    /// @inheritdoc IERC165
    function supportsInterface(bytes4 interfaceId) external pure override returns (bool isVerifier) {
        return interfaceId == this.verify.selector;
    }

    /// @inheritdoc ITypeAndVersion
    function typeAndVersion() external pure override returns (string memory) {
        return "MultiFormatVerifier 1.0.0";
    }

    /// @inheritdoc IVerifier
    function verify(
        bytes calldata payload,
        address sender
    ) external override returns (bytes memory verifierResponse) {
        if (msg.sender != i_verifierProxyAddr) revert AccessForbidden();
        require(payload.length >= 36, "Input must be at least 36 bytes");

        bytes32 configDigest;
        uint32 reportFormat;

        assembly {
        // Load the first 32 bytes (config digest) from payload calldata.
        // Note: payload.offset gives the start position of payload data.
            configDigest := calldataload(payload.offset)
        // Load the next 32 bytes starting at payload.offset + 32,
        // then shift right by 224 bits to isolate the first 4 bytes
        // (which are bytes 32-35 of the payload).
            reportFormat := shr(224, calldataload(add(payload.offset, 32)))

        // Extract the final 4 bytes of configDigest.
            let last4 := and(configDigest, 0xFFFFFFFF)
        // XOR the final 4 bytes with reportFormat.
            let newLast4 := xor(last4, reportFormat)
        // Clear the final 4 bytes of configDigest and insert the modified value.
            configDigest := or(and(configDigest, not(0xFFFFFFFF)), newLast4)
        }

        // Now, depending on the resulting reportFormat, route to the proper verification.
        if (reportFormat == REPORT_FORMAT_EVM_STREAMLINED) {
            // need to XOR with report format for non-legacy formats
            return _verifyReportStreamlined(configDigest, payload[36:] sender);
        } else {
            revert("Unsupported report format");
        }
    }

    function _verifyReportStreamlined(
        bytes32 configDigest,
        bytes calldata signedReport,
        address sender
    ) private returns (bytes memory verifierResponse) {
        (
            bytes32[3] memory reportContext,
            bytes memory reportData,
            bytes32[] memory rs,
            bytes32[] memory ss,
            bytes32 rawVs
        ) = abi.decode(signedReport, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));

        // reportContext consists of:
        // reportContext[0]: ConfigDigest
        // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
        // reportContext[2]: ExtraHash
        bytes32 configDigest = reportContext[0];

        VerifierState storage verifierState = s_verifierStates[configDigest];

        _validateReport(configDigest, rs, ss, verifierState);

        bytes32 hashedReport = keccak256(reportData);

        _verifySignatures(hashedReport, reportContext, rs, ss, rawVs, verifierState);
        emit ReportVerified(bytes32(reportData), sender);

        return reportData;
    }

    /// @notice Parses the payload and returns the extracted fields.
    /// @param payload The packed payload:
    ///        [uint16 reportLen][bytes report][uint8 numSigs][bytes sigs]
    /// @return reportLen The length of the report.
    /// @return report The report data.
    /// @return numSigs The number of signatures.
    /// @return sigs The signatures data.
    function _parsePayload(bytes calldata payload)
    internal
    pure
    returns (
        uint16 reportLen,
        bytes memory report,
        uint8 numSigs,
        bytes memory sigs
    )
    {
        // Minimal length must account for the 2-byte length, 1 byte for numSigs, and at least zero bytes for report/sigs.
        require(payload.length >= 3, "Payload too short");

        assembly {
        // p is the calldata offset of the payload data.
            let p := payload.offset

        // --- Extract the uint16 report length ---
        // Load 32 bytes from calldata starting at p.
            let word := calldataload(p)
        // The report length is stored in the first 2 bytes (most-significant bytes).
        // Shift right by 240 (256-16) to obtain the uint16 value.
            reportLen := shr(240, word)

        // --- Extract the report bytes ---
        // The report data starts immediately after the 2-byte report length.
            let reportDataOffset := add(p, 2)
        // Allocate memory for the report. The standard memory layout for dynamic bytes is:
        // [length (32 bytes)][data ...].
            let memPtr := mload(0x40)
        // Store the report length in memory.
            mstore(memPtr, reportLen)
        // Copy reportLen bytes from calldata (starting at reportDataOffset) to memory (starting at memPtr + 32).
            calldatacopy(add(memPtr, 32), reportDataOffset, reportLen)
        // Set the output variable.
            report := memPtr

        // --- Extract the uint8 number of signatures ---
        // The numSigs byte is located right after the report.
        // Its offset is: p + 2 (report length field) + reportLen.
            let numSigsOffset := add(reportDataOffset, reportLen)
        // Use the `byte` opcode to extract the first byte from the 32-byte word loaded from this offset.
            numSigs := byte(0, calldataload(numSigsOffset))

        // --- Extract the signatures bytes ---
        // The signatures start immediately after the numSigs byte.
            let sigsOffset := add(numSigsOffset, 1)
        // Calculate the signatures length:
        // Total payload length minus the bytes consumed: 2 (for uint16) + reportLen + 1 (for uint8).
            let sigsLen := sub(payload.length, add(reportLen, 3))
        // Allocate memory for the signatures bytes.
            let sigsMemPtr := mload(0x40)
        // Store the length.
            mstore(sigsMemPtr, sigsLen)
        // Copy the signatures data from calldata.
            calldatacopy(add(sigsMemPtr, 32), sigsOffset, sigsLen)
        // Set the output variable.
            sigs := sigsMemPtr

        // Update the free memory pointer.
            mstore(0x40, add(sigsMemPtr, add(32, mul(div(add(sigsLen, 31), 32), 32))))
        }
    }
}


/// @notice Validates parameters of the report
    /// @param configDigest Config digest from the report
    /// @param rs R components from the report
    /// @param ss S components from the report
    /// @param config Config for the given digest
    function _validateReport(
        bytes32 configDigest,
        bytes32[] memory rs,
        bytes32[] memory ss,
        VerifierState storage config
    ) private view {
        uint8 expectedNumSignatures = config.f + 1;

        if (!config.isActive) revert DigestInactive(configDigest);
        if (rs.length != expectedNumSignatures) revert IncorrectSignatureCount(rs.length, expectedNumSignatures);
        if (rs.length != ss.length) revert MismatchedSignatures(rs.length, ss.length);
    }

    /// @notice Verifies that a report has been signed by the correct
    /// signers and that enough signers have signed the reports.
    /// @param hashedReport The keccak256 hash of the raw report's bytes
    /// @param reportContext The context the report was signed in
    /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
    /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
    /// @param rawVs ith element is the the V component of the ith signature
    /// @param config The config digest the report was signed for
    function _verifySignatures(
        bytes32 hashedReport,
        bytes32[3] memory reportContext,
        bytes32[] memory rs,
        bytes32[] memory ss,
        bytes32 rawVs,
        VerifierState storage config
    ) private view {
        bytes32 h = keccak256(abi.encodePacked(hashedReport, reportContext));
        // i-th byte counts number of sigs made by i-th signer
        uint256 signedCount;

        Signer memory o;
        address signerAddress;
        uint256 numSigners = rs.length;
        for (uint256 i; i < numSigners; ++i) {
            signerAddress = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
            o = config.oracles[signerAddress];
            if (o.role != Role.Signer) revert BadVerification();
            unchecked {
                signedCount += 1 << (8 * o.index);
            }
        }

        if (signedCount & ORACLE_MASK != signedCount) revert BadVerification();
    }

    /// @inheritdoc IVerifier
    function updateConfig(
        bytes32 configDigest,
        address[] calldata prevSigners,
        address[] calldata newSigners,
        uint8 f
    ) external override checkConfigValid(newSigners.length, f) onlyOwner {
        VerifierState storage config = s_verifierStates[configDigest];

        if (config.f == 0) revert DigestNotSet(configDigest);

        // We must be removing the number of signers that were originally set
        if (config.oracleCount != prevSigners.length) {
            revert NonUniqueSignatures();
        }

        for (uint256 i; i < prevSigners.length; ++i) {
            // Check the signers being removed are not zero address or duplicates
            if (config.oracles[prevSigners[i]].role == Role.Unset) {
                revert NonUniqueSignatures();
            }

            delete config.oracles[prevSigners[i]];
        }

        // Once signers have been cleared we can set the new signers
        _setConfig(configDigest, newSigners, f, new Common.AddressAndWeight[](0), true);

        emit ConfigUpdated(configDigest, prevSigners, newSigners);
    }

    /// @inheritdoc IVerifier
    function setConfig(
        bytes32 configDigest,
        address[] calldata signers,
        uint8 f,
        Common.AddressAndWeight[] memory recipientAddressesAndWeights
    ) external override checkConfigValid(signers.length, f) onlyOwner {
        _setConfig(configDigest, signers, f, recipientAddressesAndWeights, false);
    }

    function _setConfig(
        bytes32 configDigest,
        address[] calldata signers,
        uint8 f,
        Common.AddressAndWeight[] memory recipientAddressesAndWeights,
        bool _updateConfig
    ) internal {
        VerifierState storage verifierState = s_verifierStates[configDigest];

        if (verifierState.f > 0 && !_updateConfig) {
            revert ConfigDigestAlreadySet();
        }

        verifierState.latestConfigBlockNumber = uint32(block.number);
        verifierState.f = f;
        verifierState.isActive = true;
        verifierState.oracleCount = uint8(signers.length);

        for (uint8 i; i < signers.length; ++i) {
            address signerAddr = signers[i];
            if (signerAddr == address(0)) revert ZeroAddress();

            // All signer roles are unset by default for a new config digest.
            // Here the contract checks to see if a signer's address has already
            // been set to ensure that the group of signer addresses that will
            // sign reports with the config digest are unique.
            bool isSignerAlreadySet = verifierState.oracles[signerAddr].role != Role.Unset;
            if (isSignerAlreadySet) revert NonUniqueSignatures();
            verifierState.oracles[signerAddr] = Signer({role: Role.Signer, index: i});
        }

        if (!_updateConfig) {
            IVerifierProxy(i_verifierProxyAddr).setVerifier(bytes32(0), configDigest, recipientAddressesAndWeights);

            emit ConfigSet(configDigest, signers, f);
        }
    }

    /// @inheritdoc IVerifier
    function activateConfig(bytes32 configDigest) external onlyOwner {
        VerifierState storage verifierState = s_verifierStates[configDigest];

        if (configDigest == bytes32("")) revert DigestEmpty();
        if (verifierState.f == 0) revert DigestNotSet(configDigest);
        verifierState.isActive = true;
        emit ConfigActivated(configDigest);
    }

    /// @inheritdoc IVerifier
    function deactivateConfig(bytes32 configDigest) external onlyOwner {
        VerifierState storage verifierState = s_verifierStates[configDigest];

        if (configDigest == bytes32("")) revert DigestEmpty();
        if (verifierState.f == 0) revert DigestNotSet(configDigest);
        verifierState.isActive = false;
        emit ConfigDeactivated(configDigest);
    }

    /// @inheritdoc IVerifier
    function latestConfigDetails(bytes32 configDigest) external view override returns (uint32 blockNumber) {
        VerifierState storage verifierState = s_verifierStates[configDigest];
        return (verifierState.latestConfigBlockNumber);
    }
}
