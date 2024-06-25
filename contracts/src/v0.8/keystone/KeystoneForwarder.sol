// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IReceiver} from "./interfaces/IReceiver.sol";
import {IRouter} from "./interfaces/IRouter.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

/// @notice This is an entry point for `write_${chain}` Target capability. It
/// allows nodes to determine if reports have been processed (successfully or
/// not) in a decentralized and product-agnostic way by recording processed
/// reports.
contract KeystoneForwarder is OwnerIsCreator, ITypeAndVersion, IRouter {
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

  /// @notice This error is thrown whenever a report specifies a configuration that
  /// does not exist.
  /// @param configId (uint64(donId) << 32) | configVersion
  error InvalidConfig(uint64 configId);

  /// @notice This error is thrown whenever a signer address is not in the
  /// configuration.
  /// @param signer The signer address that was not in the configuration
  error InvalidSigner(address signer);

  /// @notice This error is thrown whenever a signature is invalid.
  /// @param signature The signature that was invalid
  error InvalidSignature(bytes signature);

  /// @notice Contains the signing address of each oracle
  struct OracleSet {
    uint8 f; // Number of faulty nodes allowed
    address[] signers;
    mapping(address signer => uint256 position) _positions; // 1-indexed to detect unset values
  }

  /// @notice Contains the configuration for each DON ID
  // @param configId (uint64(donId) << 32) | configVersion
  mapping(uint64 configId => OracleSet) internal s_configs;

  event ConfigSet(uint32 indexed donId, uint32 indexed configVersion, uint8 f, address[] signers);

  /// @notice Emitted when a report is processed
  /// @param result The result of the attempted delivery. True if successful.
  event ReportProcessed(
    address indexed receiver,
    bytes32 indexed workflowExecutionId,
    bytes2 indexed reportId,
    bool result
  );

  string public constant override typeAndVersion = "Forwarder and Router 1.0.0";

  constructor() OwnerIsCreator() {
    s_forwarders[address(this)] = true;
  }

  uint256 internal constant MAX_ORACLES = 31;
  uint256 internal constant METADATA_LENGTH = 109;
  uint256 internal constant FORWARDER_METADATA_LENGTH = 45;
  uint256 internal constant SIGNATURE_LENGTH = 65;

  // ================================================================
  // │                          Router                              │
  // ================================================================

  mapping(address forwarder => bool) internal s_forwarders;
  mapping(bytes32 transmissionId => TransmissionInfo) internal s_transmissions;

  function addForwarder(address forwarder) external onlyOwner {
    s_forwarders[forwarder] = true;
    emit ForwarderAdded(forwarder);
  }

  function removeForwarder(address forwarder) external onlyOwner {
    s_forwarders[forwarder] = false;
    emit ForwarderRemoved(forwarder);
  }

  function route(
    bytes32 transmissionId,
    address transmitter,
    address receiver,
    bytes calldata metadata,
    bytes calldata validatedReport
  ) public returns (bool) {
    if (!s_forwarders[msg.sender]) {
      revert UnauthorizedForwarder();
    }

    if (s_transmissions[transmissionId].transmitter != address(0)) revert AlreadyAttempted(transmissionId);
    s_transmissions[transmissionId].transmitter = transmitter;

    if (receiver.code.length == 0) return false;

    try IReceiver(receiver).onReport(metadata, validatedReport) {
      s_transmissions[transmissionId].success = true;
      return true;
    } catch {
      return false;
    }
  }

  function getTransmissionId(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) public pure returns (bytes32) {
    // This is slightly cheaper compared to
    // keccak256(abi.encode(receiver, workflowExecutionId, reportId));
    return keccak256(bytes.concat(bytes20(uint160(receiver)), workflowExecutionId, reportId));
  }

  /// @notice Get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (address) {
    return s_transmissions[getTransmissionId(receiver, workflowExecutionId, reportId)].transmitter;
  }

  /// @notice Get delivery status of a given report
  function getTransmissionState(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (IRouter.TransmissionState) {
    bytes32 transmissionId = getTransmissionId(receiver, workflowExecutionId, reportId);

    if (s_transmissions[transmissionId].transmitter == address(0)) return IRouter.TransmissionState.NOT_ATTEMPTED;
    return
      s_transmissions[transmissionId].success ? IRouter.TransmissionState.SUCCEEDED : IRouter.TransmissionState.FAILED;
  }

  function isForwarder(address forwarder) external view returns (bool) {
    return s_forwarders[forwarder];
  }

  // ================================================================
  // │                          Forwarder                           │
  // ================================================================

  function setConfig(uint32 donId, uint32 configVersion, uint8 f, address[] calldata signers) external onlyOwner {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (signers.length > MAX_ORACLES) revert ExcessSigners(signers.length, MAX_ORACLES);
    if (signers.length <= 3 * f) revert InsufficientSigners(signers.length, 3 * f + 1);

    uint64 configId = (uint64(donId) << 32) | configVersion;

    // remove any old signer addresses
    for (uint256 i = 0; i < s_configs[configId].signers.length; ++i) {
      delete s_configs[configId]._positions[s_configs[configId].signers[i]];
    }

    // add new signer addresses
    for (uint256 i = 0; i < signers.length; ++i) {
      // assign indices, detect duplicates
      address signer = signers[i];
      if (s_configs[configId]._positions[signer] != 0) revert DuplicateSigner(signer);
      s_configs[configId]._positions[signer] = i + 1;
    }
    s_configs[configId].signers = signers;
    s_configs[configId].f = f;

    emit ConfigSet(donId, configVersion, f, signers);
  }

  function clearConfig(uint32 donId, uint32 configVersion) external onlyOwner {
    // We are not removing old signer positions, because it is sufficient to
    // clear the f value for `report` function. If we decide to restore
    // the configId in the future, the setConfig function clears the positions.
    s_configs[(uint64(donId) << 32) | configVersion].f = 0;

    emit ConfigSet(donId, configVersion, 0, new address[](0));
  }

  // send a report to receiver
  function report(
    address receiver,
    bytes calldata rawReport,
    bytes calldata reportContext,
    bytes[] calldata signatures
  ) external {
    if (rawReport.length < METADATA_LENGTH) {
      revert InvalidReport();
    }

    bytes32 workflowExecutionId;
    bytes2 reportId;
    {
      uint64 configId;
      (workflowExecutionId, configId, reportId) = _getMetadata(rawReport);
      OracleSet storage config = s_configs[configId];

      uint8 f = config.f;
      // f can never be 0, so this means the config doesn't actually exist
      if (f == 0) revert InvalidConfig(configId);
      if (f + 1 != signatures.length) revert InvalidSignatureCount(f + 1, signatures.length);

      // validate signatures
      bytes32 completeHash = keccak256(abi.encodePacked(keccak256(rawReport), reportContext));
      address[MAX_ORACLES + 1] memory signed;
      for (uint256 i = 0; i < signatures.length; ++i) {
        bytes calldata signature = signatures[i];
        if (signature.length != SIGNATURE_LENGTH) revert InvalidSignature(signature);
        address signer = ecrecover(
          completeHash,
          uint8(signature[64]) + 27,
          bytes32(signature[0:32]),
          bytes32(signature[32:64])
        );

        // validate signer is trusted and signature is unique
        uint256 index = config._positions[signer];
        if (index == 0) revert InvalidSigner(signer); // index is 1-indexed so we can detect unset signers
        if (signed[index] != address(0)) revert DuplicateSigner(signer);
        signed[index] = signer;
      }
    }

    bool success = this.route(
      getTransmissionId(receiver, workflowExecutionId, reportId),
      msg.sender,
      receiver,
      rawReport[FORWARDER_METADATA_LENGTH:METADATA_LENGTH],
      rawReport[METADATA_LENGTH:]
    );

    emit ReportProcessed(receiver, workflowExecutionId, reportId, success);
  }

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getMetadata(
    bytes memory rawReport
  ) internal pure returns (bytes32 workflowExecutionId, uint64 configId, bytes2 reportId) {
    // (first 32 bytes of memory contain length of the report)
    // version                  // offset  32, size  1
    // workflow_execution_id    // offset  33, size 32
    // timestamp                // offset  65, size  4
    // don_id                   // offset  69, size  4
    // don_config_version,	    // offset  73, size  4
    // workflow_cid             // offset  77, size 32
    // workflow_name            // offset 109, size 10
    // workflow_owner           // offset 119, size 20
    // report_id              // offset 139, size  2
    assembly {
      workflowExecutionId := mload(add(rawReport, 33))
      // shift right by 24 bytes to get the combined don_id and don_config_version
      configId := shr(mul(24, 8), mload(add(rawReport, 69)))
      reportId := mload(add(rawReport, 139))
    }
  }
}
