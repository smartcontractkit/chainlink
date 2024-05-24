// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IForwarder} from "./interfaces/IForwarder.sol";
import {IReceiver} from "./interfaces/IReceiver.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";

/// @notice This is an entry point for `write_${chain}` Target capability. It
/// allows nodes to determine if reports have been processed (successfully or
/// not) in a decentralized and product-agnostic way by recording processed
/// reports.
contract KeystoneForwarder is IForwarder, ConfirmedOwner, TypeAndVersionInterface {
  error ReentrantCall();

  /// @notice This error is returned when the report is shorter than
  /// REPORT_METADATA_LENGTH, which is the minimum length of a report.
  error InvalidReport();

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

  /// @notice This error is thrown whenever a report specifies a DON ID that
  /// does not have a configuration.
  /// @param donId The DON ID that was provided in the report
  error InvalidDonId(uint32 donId);

  /// @notice This error is thrown whenever a signer address is not in the
  /// configuration.
  /// @param signer The signer address that was not in the configuration
  error InvalidSigner(address signer);

  /// @notice This error is thrown whenever a signature is invalid.
  /// @param signature The signature that was invalid
  error InvalidSignature(bytes signature);

  /// @notice This error is thrown whenever a report has already been processed.
  /// @param reportId The ID of the report that was already processed
  error ReportAlreadyProcessed(bytes32 reportId);

  bool internal s_reentrancyGuard; // guard against reentrancy

  /// @notice Contains the signing address of each oracle
  struct OracleSet {
    uint8 f; // Number of faulty nodes allowed
    address[] signers;
    mapping(address => uint256) _positions; // 1-indexed to detect unset values
  }

  /// @notice Contains the configuration for each DON ID
  mapping(uint32 donId => OracleSet) internal s_configs;

  struct DeliveryStatus {
    address transmitter;
    bool success;
  }

  mapping(bytes32 reportId => DeliveryStatus status) internal s_reports;

  /// @notice Emitted when a report is processed
  /// @param receiver The address of the receiver contract
  /// @param workflowOwner The address of the workflow owner
  /// @param workflowExecutionId The ID of the workflow execution
  /// @param result The result of the attempted delivery. True if successful.
  event ReportProcessed(
    address indexed receiver,
    address indexed workflowOwner,
    bytes32 indexed workflowExecutionId,
    bool result
  );

  constructor() ConfirmedOwner(msg.sender) {}

  uint256 internal constant MAX_ORACLES = 31;
  // 32 bytes for workflowId, 4 bytes for donId, 32 bytes for
  // workflowExecutionId, 20 bytes for workflowOwner
  uint256 internal constant REPORT_METADATA_LENGTH = 88;
  uint256 internal constant SIGNATURE_LENGTH = 65;

  function setConfig(uint32 donId, uint8 f, address[] calldata signers) external nonReentrant {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (signers.length > MAX_ORACLES) revert ExcessSigners(signers.length, MAX_ORACLES);
    if (signers.length <= 3 * f) revert InsufficientSigners(signers.length, 3 * f + 1);

    // TODO: how does setConfig handle expiration? e.g. if the signer set changes

    // remove any old signer addresses
    for (uint256 i; i < s_configs[donId].signers.length; ++i) {
      address signer = s_configs[donId].signers[i];
      delete s_configs[donId]._positions[signer];
    }

    // add new signer addresses
    s_configs[donId].signers = signers;
    for (uint256 i; i < signers.length; ++i) {
      // assign indices, detect duplicates
      address signer = signers[i];
      if (s_configs[donId]._positions[signer] != 0) revert DuplicateSigner(signer);
      s_configs[donId]._positions[signer] = uint8(i) + 1;
      s_configs[donId].signers.push(signer);
    }
    s_configs[donId].f = f;
  }

  // send a report to receiver
  function report(
    address receiverAddress,
    bytes calldata rawReport,
    bytes[] calldata signatures
  ) external nonReentrant {
    if (rawReport.length < REPORT_METADATA_LENGTH) {
      revert InvalidReport();
    }

    (bytes32 workflowId, uint32 donId, bytes32 workflowExecutionId, address workflowOwner) = _getMetadata(rawReport);

    // f can never be 0, so this means the config doesn't actually exist
    if (s_configs[donId].f == 0) revert InvalidDonId(donId);

    bytes32 reportId = _reportId(receiverAddress, workflowExecutionId);
    if (s_reports[reportId].transmitter != address(0)) revert ReportAlreadyProcessed(reportId);

    if (s_configs[donId].f + 1 != signatures.length)
      revert InvalidSignatureCount(s_configs[donId].f + 1, signatures.length);

    // validate signatures
    {
      bytes32 hash = keccak256(rawReport);

      address[MAX_ORACLES] memory signed;
      uint8 index;
      for (uint256 i; i < signatures.length; ++i) {
        // TODO: is libocr-style multiple bytes32 arrays more optimal, gas-wise?
        (bytes32 r, bytes32 s, uint8 v) = _splitSignature(signatures[i]);
        address signer = ecrecover(hash, v, r, s);

        // validate signer is trusted and signature is unique
        index = uint8(s_configs[donId]._positions[signer]);
        if (index == 0) revert InvalidSigner(signer); // index is 1-indexed so we can detect unset signers
        index -= 1;
        if (signed[index] != address(0)) revert DuplicateSigner(signer);
        signed[index] = signer;
      }
    }

    IReceiver receiver = IReceiver(receiverAddress);
    bool success;
    try receiver.onReport(workflowId, workflowOwner, rawReport[REPORT_METADATA_LENGTH:]) {
      success = true;
    } catch {
      // Do nothing, success is already false
    }

    s_reports[reportId] = DeliveryStatus(msg.sender, success);

    emit ReportProcessed(receiverAddress, workflowOwner, workflowExecutionId, success);
  }

  function _reportId(address receiver, bytes32 workflowExecutionId) internal pure returns (bytes32) {
    // TODO: gas savings: could we just use a bytes key and avoid another keccak256 call
    return keccak256(bytes.concat(bytes20(uint160(receiver)), workflowExecutionId));
  }

  // get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(address receiver, bytes32 workflowExecutionId) external view returns (address) {
    bytes32 reportId = _reportId(receiver, workflowExecutionId);
    return s_reports[reportId].transmitter;
  }

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

  function _getMetadata(
    bytes memory rawReport
  ) internal pure returns (bytes32 workflowId, uint32 donId, bytes32 workflowExecutionId, address workflowOwner) {
    assembly {
      // skip first 32 bytes, contains length of the report
      // first 32 bytes is the workflowId
      workflowId := mload(add(rawReport, 32))
      // next 4 bytes is donId. We shift right by 28 bytes to get the actual value
      donId := shr(mul(28, 8), mload(add(rawReport, 64)))
      // next 32 bytes is the workflowExecutionId
      workflowExecutionId := mload(add(rawReport, 68))
      // next 20 bytes is the workflowOwner. We shift right by 12 bytes to get
      // the actual value
      workflowOwner := shr(mul(12, 8), mload(add(rawReport, 100)))
    }
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "KeystoneForwarder 1.0.0";
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
