// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "./interfaces/IReceiver.sol";
import {IRouter} from "./interfaces/IRouter.sol";
import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

import {ERC165Checker} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/ERC165Checker.sol";

/// @notice This is an entry point for `write_${chain}` Target capability. It allows nodes to
/// determine if reports have been processed (successfully or not) in a decentralized and
/// product-agnostic way by recording processed reports.
contract KeystoneForwarder is OwnerIsCreator, ITypeAndVersion, IRouter {
  /// @notice This error is returned when the report is shorter than REPORT_METADATA_LENGTH,
  /// which is the minimum length of a report.
  error InvalidReport();

  /// @notice This error is thrown whenever trying to set a config with a fault tolerance of 0.
  error FaultToleranceMustBePositive();

  /// @notice This error is thrown whenever configuration provides more signers than the maximum allowed number.
  /// @param numSigners The number of signers who have signed the report
  /// @param maxSigners The maximum number of signers that can sign a report
  error ExcessSigners(uint256 numSigners, uint256 maxSigners);

  /// @notice This error is thrown whenever a configuration is provided with less than the minimum number of signers.
  /// @param numSigners The number of signers provided
  /// @param minSigners The minimum number of signers expected
  error InsufficientSigners(uint256 numSigners, uint256 minSigners);

  /// @notice This error is thrown whenever a duplicate signer address is provided in the configuration.
  /// @param signer The signer address that was duplicated.
  error DuplicateSigner(address signer);

  /// @notice This error is thrown whenever a report has an incorrect number of signatures.
  /// @param expected The number of signatures expected, F + 1
  /// @param received The number of signatures received
  error InvalidSignatureCount(uint256 expected, uint256 received);

  /// @notice This error is thrown whenever a report specifies a configuration that does not exist.
  /// @param configId (uint64(donId) << 32) | configVersion
  error InvalidConfig(uint64 configId);

  /// @notice This error is thrown whenever a signer address is not in the configuration or
  /// when trying to set a zero address as a signer.
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

  struct Transmission {
    address transmitter;
    // This is true if the receiver is not a contract or does not implement the `IReceiver` interface.
    bool invalidReceiver;
    // Whether the transmission attempt was successful. If `false`, the transmission can be retried
    // with an increased gas limit.
    bool success;
    // The amount of gas allocated for the `IReceiver.onReport` call. uint80 allows storing gas for known EVM block
    // gas limits. Ensures that the minimum gas requested by the user is available during the transmission attempt.
    // If the transmission fails (indicated by a `false` success state), it can be retried with an increased gas limit.
    uint80 gasLimit;
  }

  /// @notice Emitted when a report is processed
  /// @param result The result of the attempted delivery. True if successful.
  event ReportProcessed(
    address indexed receiver,
    bytes32 indexed workflowExecutionId,
    bytes2 indexed reportId,
    bool result
  );

  /// @notice Contains the configuration for each DON ID
  /// configId (uint64(donId) << 32) | configVersion
  mapping(uint64 configId => OracleSet oracleSet) internal s_configs;

  event ConfigSet(uint32 indexed donId, uint32 indexed configVersion, uint8 f, address[] signers);

  string public constant override typeAndVersion = "KeystoneForwarder 1.0.0";

  constructor() OwnerIsCreator() {
    s_forwarders[address(this)] = true;
  }

  uint256 internal constant MAX_ORACLES = 31;
  uint256 internal constant METADATA_LENGTH = 109;
  uint256 internal constant FORWARDER_METADATA_LENGTH = 45;
  uint256 internal constant SIGNATURE_LENGTH = 65;

  /// @dev This is the gas required to store `success` after the report is processed.
  /// It is a warm storage write because of the packed struct. In practice it will cost less.
  uint256 internal constant INTERNAL_GAS_REQUIREMENTS_AFTER_REPORT = 5_000;
  /// @dev This is the gas required to store the transmission struct and perform other checks.
  uint256 internal constant INTERNAL_GAS_REQUIREMENTS = 25_000 + INTERNAL_GAS_REQUIREMENTS_AFTER_REPORT;
  /// @dev This is the minimum gas required to route a report. This includes internal gas requirements
  /// as well as the minimum gas that the user contract will receive. 30k * 3 gas is to account for
  /// cases where consumers need close to the 30k limit provided in the supportsInterface check.
  uint256 internal constant MINIMUM_GAS_LIMIT = INTERNAL_GAS_REQUIREMENTS + 30_000 * 3 + 10_000;

  // ================================================================
  // │                          Router                              │
  // ================================================================

  mapping(address forwarder => bool isForwarder) internal s_forwarders;
  mapping(bytes32 transmissionId => Transmission transmission) internal s_transmissions;

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
    if (!s_forwarders[msg.sender]) revert UnauthorizedForwarder();

    uint256 gasLimit = gasleft() - INTERNAL_GAS_REQUIREMENTS;
    if (gasLimit < MINIMUM_GAS_LIMIT) revert InsufficientGasForRouting(transmissionId);

    Transmission memory transmission = s_transmissions[transmissionId];
    if (transmission.success || transmission.invalidReceiver) revert AlreadyAttempted(transmissionId);

    s_transmissions[transmissionId].transmitter = transmitter;
    s_transmissions[transmissionId].gasLimit = uint80(gasLimit);

    // This call can consume up to 90k gas.
    if (!ERC165Checker.supportsInterface(receiver, type(IReceiver).interfaceId)) {
      s_transmissions[transmissionId].invalidReceiver = true;
      return false;
    }

    bool success;
    bytes memory payload = abi.encodeCall(IReceiver.onReport, (metadata, validatedReport));

    uint256 remainingGas = gasleft() - INTERNAL_GAS_REQUIREMENTS_AFTER_REPORT;
    assembly {
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(remainingGas, receiver, 0, add(payload, 0x20), mload(payload), 0x0, 0x0)
    }

    if (success) {
      s_transmissions[transmissionId].success = true;
    }
    return success;
  }

  function getTransmissionId(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) public pure returns (bytes32) {
    // This is slightly cheaper compared to `keccak256(abi.encode(receiver, workflowExecutionId, reportId));`
    return keccak256(bytes.concat(bytes20(uint160(receiver)), workflowExecutionId, reportId));
  }

  function getTransmissionInfo(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (TransmissionInfo memory) {
    bytes32 transmissionId = getTransmissionId(receiver, workflowExecutionId, reportId);

    Transmission memory transmission = s_transmissions[transmissionId];

    TransmissionState state;

    if (transmission.transmitter == address(0)) {
      state = IRouter.TransmissionState.NOT_ATTEMPTED;
    } else if (transmission.invalidReceiver) {
      state = IRouter.TransmissionState.INVALID_RECEIVER;
    } else {
      state = transmission.success ? IRouter.TransmissionState.SUCCEEDED : IRouter.TransmissionState.FAILED;
    }

    return
      TransmissionInfo({
        gasLimit: transmission.gasLimit,
        invalidReceiver: transmission.invalidReceiver,
        state: state,
        success: transmission.success,
        transmissionId: transmissionId,
        transmitter: transmission.transmitter
      });
  }

  /// @notice Get transmitter of a given report or 0x0 if it wasn't transmitted yet
  function getTransmitter(
    address receiver,
    bytes32 workflowExecutionId,
    bytes2 reportId
  ) external view returns (address) {
    return s_transmissions[getTransmissionId(receiver, workflowExecutionId, reportId)].transmitter;
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
      if (signer == address(0)) revert InvalidSigner(signer);
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
    // version                offset  32, size  1
    // workflow_execution_id  offset  33, size 32
    // timestamp              offset  65, size  4
    // don_id                 offset  69, size  4
    // don_config_version,    offset  73, size  4
    // workflow_cid           offset  77, size 32
    // workflow_name          offset 109, size 10
    // workflow_owner         offset 119, size 20
    // report_id              offset 139, size  2
    assembly {
      workflowExecutionId := mload(add(rawReport, 33))
      // shift right by 24 bytes to get the combined don_id and don_config_version
      configId := shr(mul(24, 8), mload(add(rawReport, 69)))
      reportId := mload(add(rawReport, 139))
    }
  }
}
