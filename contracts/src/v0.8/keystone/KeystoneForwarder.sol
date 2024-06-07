// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouter} from "./interfaces/IRouter.sol";
import {IForwarder} from "./interfaces/IForwarder.sol";
import {IReceiver} from "./interfaces/IReceiver.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

/// @notice This is an entry point for `write_${chain}` Target capability. It
/// allows nodes to determine if reports have been processed (successfully or
/// not) in a decentralized and product-agnostic way by recording processed
/// reports.
contract KeystoneForwarder is IForwarder, OwnerIsCreator, ITypeAndVersion {
  error ReentrantCall();

  /// @notice This error is returned when the report is shorter than
  /// REPORT_METADATA_LENGTH, which is the minimum length of a report.
  error InvalidReport();

  /// @notice This error is returned when the metadata version is not supported.
  error InvalidVersion(uint8 version);

  /// @notice This error is thrown whenever trying to set a config with a fault
  /// tolerance of 0.
  error FaultToleranceMustBePositive();

  /// @notice This error is thrown whenever configuration provides more signers
  /// than the maximum allowed number.
  /// @param numSigners The number of signers who have signed the report
  /// @param maxSigners The maximum number of signers that can sign a report
  error ExcessSigners(uint256 numSigners, uint256 maxSigners);

  /// @notice This error is thrown whenever a configuration is provided with
  /// less than the minimum number of signers.
  /// @param numSigners The number of signers provided
  /// @param minSigners The minimum number of signers expected
  error InsufficientSigners(uint256 numSigners, uint256 minSigners);

  /// @notice This error is thrown whenever a duplicate signer address is
  /// provided in the configuration.
  /// @param signer The signer address that was duplicated.
  error DuplicateSigner(address signer);

  /// @notice This error is thrown whenever a report has an incorrect number of
  /// signatures.
  /// @param expected The number of signatures expected, F + 1
  /// @param received The number of signatures received
  error InvalidSignatureCount(uint256 expected, uint256 received);

  /// @notice This error is thrown whenever a report specifies a configuration that
  /// does not exist.
  /// @param donId The DON ID that was provided in the report
  /// @param configVersion The config version that was provided in the report
  error InvalidConfig(uint32 donId, uint32 configVersion);

  /// @notice This error is thrown whenever a signer address is not in the
  /// configuration.
  /// @param signer The signer address that was not in the configuration
  error InvalidSigner(address signer);

  /// @notice This error is thrown whenever a signature is invalid.
  /// @param signature The signature that was invalid
  error InvalidSignature(bytes signature);

  /// @notice This error is thrown whenever a message has already been processed.
  /// @param messageId The ID of the message that was already processed
  error AlreadyProcessed(bytes32 messageId);

  bool internal s_reentrancyGuard; // guard against reentrancy

  /// @notice Contains the signing address of each oracle
  struct OracleSet {
    uint8 f; // Number of faulty nodes allowed
    address[] signers;
    mapping(address => uint256) _positions; // 1-indexed to detect unset values
  }

  address internal s_router;

  /// @notice Contains the configuration for each DON ID
  // @param configId keccak256(donId, donConfigVersion)
  mapping(bytes32 configId => OracleSet) internal s_configs;

  event ConfigSet(
    uint32 indexed donId,
    uint32 indexed configVersion,
    uint8 f,
    address[] signers
  );

  /// @notice Emitted when a report is processed
  /// @param receiver The address of the receiver contract
  /// @param workflowExecutionId The ID of the workflow execution
  /// @param result The result of the attempted delivery. True if successful.
  event ReportProcessed(address indexed receiver, bytes32 indexed workflowExecutionId, bool result);

  string public constant override typeAndVersion = "KeystoneForwarder 1.0.0";

  constructor(address router) OwnerIsCreator() {
    s_router = router;
  }

  uint256 internal constant MAX_ORACLES = 31;
  uint256 internal constant METADATA_LENGTH = 109;
  uint256 internal constant FORWARDER_METADATA_LENGTH = 45;
  uint256 internal constant SIGNATURE_LENGTH = 65;

  function setConfig(uint32 donId, uint32 configVersion, uint8 f, address[] calldata signers) external onlyOwner {
    // TODO: configSet event

    if (f == 0) revert FaultToleranceMustBePositive();
    if (signers.length > MAX_ORACLES) revert ExcessSigners(signers.length, MAX_ORACLES);
    if (signers.length <= 3 * f) revert InsufficientSigners(signers.length, 3 * f + 1);

    bytes32 configId = keccak256(abi.encode(donId, configVersion));

    // remove any old signer addresses
    for (uint256 i; i < s_configs[configId].signers.length; ++i) {
      address signer = s_configs[configId].signers[i];
      delete s_configs[configId]._positions[signer];
    }

    // add new signer addresses
    s_configs[configId].signers = signers;
    for (uint256 i; i < signers.length; ++i) {
      // assign indices, detect duplicates
      address signer = signers[i];
      if (s_configs[configId]._positions[signer] != 0) revert DuplicateSigner(signer);
      s_configs[configId]._positions[signer] = uint8(i) + 1;
      s_configs[configId].signers.push(signer);
    }
    s_configs[configId].f = f;

    emit ConfigSet(donId, configVersion, f, signers);
  }

  function clearConfig(uint32 donId, uint32 configVersion) external onlyOwner {
    bytes32 configId = keccak256(abi.encode(donId, configVersion));

    // remove any old signer addresses
    for (uint256 i; i < s_configs[configId].signers.length; ++i) {
      address signer = s_configs[configId].signers[i];
      delete s_configs[configId]._positions[signer];
    }

    s_configs[configId].f = 0;

    emit ConfigSet(donId, configVersion, 0, new address[](0));
  }

  // send a report to receiver
  function report(
    address receiverAddress,
    bytes calldata rawReport,
    bytes calldata reportContext,
    bytes[] calldata signatures
  ) external nonReentrant {
    if (rawReport.length < METADATA_LENGTH) {
      revert InvalidReport();
    }

    bytes32 workflowExecutionId;
    bytes32 combinedId;
    {
      bytes32 configId;
      {
        uint32 donId;
        uint32 configVersion;
        bytes2 reportId;
        (workflowExecutionId, donId, configVersion, reportId) = _getMetadata(rawReport);

        configId = keccak256(abi.encode(donId, configVersion));

        uint8 f = s_configs[configId].f;
        // f can never be 0, so this means the config doesn't actually exist
        if (f == 0) revert InvalidConfig(donId, configVersion);
        if (f + 1 != signatures.length) revert InvalidSignatureCount(f + 1, signatures.length);

        combinedId = _combinedId(receiverAddress, workflowExecutionId, reportId);
      }

      // validate signatures
      {
        bytes32 completeHash = keccak256(abi.encodePacked(keccak256(rawReport), reportContext));

        address[MAX_ORACLES] memory signed;
        uint8 index;
        for (uint256 i; i < signatures.length; ++i) {
          (bytes32 r, bytes32 s, uint8 v) = _splitSignature(signatures[i]);
          address signer = ecrecover(completeHash, v + 27, r, s);

          // validate signer is trusted and signature is unique
          index = uint8(s_configs[configId]._positions[signer]);
          if (index == 0) revert InvalidSigner(signer); // index is 1-indexed so we can detect unset signers
          index -= 1;
          if (signed[index] != address(0)) revert DuplicateSigner(signer);
          signed[index] = signer;
        }
      }
    }

    // TODO: make router address a compile time address to save on read from storage?
    bool success = IRouter(s_router).route(
      combinedId,
      msg.sender,
      receiverAddress,
      rawReport[FORWARDER_METADATA_LENGTH:METADATA_LENGTH],
      rawReport[METADATA_LENGTH:]
    );

    emit ReportProcessed(receiverAddress, workflowExecutionId, success);
  }

  function _combinedId(address receiver, bytes32 workflowExecutionId, bytes2 reportId) internal pure returns (bytes32) {
    return keccak256(bytes.concat(bytes20(uint160(receiver)), workflowExecutionId, reportId));
  }

  function getTransmitter(
    address receiverAddress,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (address) {
    bytes32 combinedId = _combinedId(receiverAddress, workflowExecutionId, reportId);
    return IRouter(s_router).getTransmitter(combinedId);
  }

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _splitSignature(bytes memory sig) internal pure returns (bytes32 r, bytes32 s, uint8 v) {
    if (sig.length != SIGNATURE_LENGTH) revert InvalidSignature(sig);

    assembly {
      /*
      First 32 bytes stores the length of the signature

      add(sig, 32) = pointer of sig + 32
      effectively, skips first 32 bytes of signature

      mload(p) loads next 32 bytes starting at the memory address p into memory
      */

      // first 32 bytes, after the length prefix
      r := mload(add(sig, 32))
      // second 32 bytes
      s := mload(add(sig, 64))
      // final byte (first byte of the next 32 bytes)
      v := byte(0, mload(add(sig, 96)))
    }
  }

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getMetadata(
    bytes memory rawReport
  ) internal pure returns (bytes32 workflowExecutionId, uint32 donId, uint32 donConfigVersion, bytes2 reportId) {
    // (first 32 bytes of memory contain length of the report)
    // version                  // offset  32, size  1
    // workflow_execution_id    // offset  33, size 32
    // timestamp                // offset  65, size  4
    // don_id                   // offset  69, size  4
    // don_config_version,	    // offset  73, size  4
    // workflow_cid             // offset  77, size 32
    // workflow_name            // offset 109, size 10
    // workflow_owner           // offset 119, size 20
    // report_name              // offset 139, size  2
    if (uint8(rawReport[0]) != 1) {
      revert InvalidVersion(uint8(rawReport[0]));
    }
    assembly {
      workflowExecutionId := mload(add(rawReport, 33))
      // shift right by 28 bytes to get the actual value
      donId := shr(mul(28, 8), mload(add(rawReport, 69)))
      // shift right by 28 bytes to get the actual value
      donConfigVersion := shr(mul(28, 8), mload(add(rawReport, 73)))
      reportId := mload(add(rawReport, 139))
    }
  }

  /**
   * @dev replicates Open Zeppelin's ReentrancyGuard but optimized to fit our storage
   */
  modifier nonReentrant() {
    if (s_reentrancyGuard) revert ReentrantCall();
    s_reentrancyGuard = true;
    _;
    s_reentrancyGuard = false;
  }
}
