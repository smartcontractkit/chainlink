// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IConfigurator} from "./interfaces/IConfigurator.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

// Subsequent versions of the onchain config must be backwards compatible and only append fields
uint256 constant MIN_SUPPORTED_ONCHAIN_CONFIG_VERSION = 1;

/**
 * @title Configurator
 * @author samsondav
 * @notice This contract is intended to be deployed on the source chain and acts as a OCR3 configurator for LLO/Mercury
 **/

contract Configurator is IConfigurator, ConfirmedOwner, ITypeAndVersion, IERC165 {
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

  /// @notice This error is thrown whenever the onchainConfig length is invalid
  /// (must be at least 64 bytes)
  /// @param onchainConfigLength The (wrong) length of the onchainConfig
  error InvalidOnchainLength(uint256 onchainConfigLength);

  /// @notice This error is thrown if the onchainConfig version is too old.
  /// @param version The version of the onchainConfig
  error UnsupportedOnchainConfigVersion(uint256 version);

  /// @notice This event is emitted when a production config is set with a non-zero predecessor config digest in the on-chain config.
  /// @param predecessorConfigDigest The predecessor config digest
  error NonZeroPredecessorConfigDigest(bytes32 predecessorConfigDigest);

  /// @notice This event is emitted when a staging config is set with a predecessor config digest that does not match the current production config digest.
  /// @param predecessorConfigDigest The predecessor config digest
  error InvalidPredecessorConfigDigest(bytes32 predecessorConfigDigest);

  /// @notice This event is emitted during promoteStagingConfig if the isGreenProduction flag does not match the contract state
  /// @param configId The configId
  /// @param isGreenProductionContractState The current (correct) isGreenProduction state according to the contract
  error IsGreenProductionMustMatchContractState(bytes32 configId, bool isGreenProductionContractState);

  /// @notice This event is emitted during promoteStagingConfig if the configId has never been set
  /// @param configId The configId that has never been set
  error ConfigUnset(bytes32 configId);

  /// @notice This event is emitted during promoteStagingConfig if the configId has never been set as a staging config
  /// @param configId The configId that has never been set as a staging config
  /// @param isGreenProduction The isGreenProduction flag
  error ConfigUnsetStaging(bytes32 configId, bool isGreenProduction);

  /// @notice This event is emitted during promoteStagingConfig if the configId has never been set as a production config
  /// @param configId The configId that has never been set as a production config
  /// @param isGreenProduction The isGreenProduction flag
  error ConfigUnsetProduction(bytes32 configId, bool isGreenProduction);

  struct ConfigurationState {
    // The number of times a configuration (either staging or production) has
    // been set for this configId
    uint64 configCount;
    // The block number of the block the last time
    // the configuration was updated.
    uint32 latestConfigBlockNumber;
    // isGreenProduction is a bit flip that indicates whether blue is production
    // exactly one of blue/green must be production at all times.
    // 0 -> blue is production
    // 1 -> green is production
    //
    // So, to clarify, if isGreenProduction is false (initial state) then:
    // [0](blue) is production and [1](green) is staging/retired
    //
    // and if isGreenProduction is true then:
    // [0](blue) is staging/retired and [1](green) is production
    // State is swapped every time a staging config is promoted to production.
    bool isGreenProduction;
    // The digest of the current configurations (0 is always blue, 1 is always green)
    bytes32[2] configDigest;
  }

  constructor() ConfirmedOwner(msg.sender) {}

  /// @notice Configuration states keyed on config ID
  /// @dev The first element is the blue configuration state
  /// and the second element is the green configuration state
  mapping(bytes32 => ConfigurationState) internal s_configurationStates;

  /// @inheritdoc IConfigurator
  function setProductionConfig(
    bytes32 configId,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    if (onchainConfig.length < 64) revert InvalidOnchainLength(onchainConfig.length);

    // Ensure that predecessorConfigDigest is unset and version is correct
    uint256 version;
    bytes32 predecessorConfigDigest;
    assembly {
      version := mload(add(onchainConfig, 32))
      predecessorConfigDigest := mload(add(onchainConfig, 64))
    }
    if (version < MIN_SUPPORTED_ONCHAIN_CONFIG_VERSION) revert UnsupportedOnchainConfigVersion(version);
    if (predecessorConfigDigest != 0) revert NonZeroPredecessorConfigDigest(predecessorConfigDigest);

    _setConfig(
      configId,
      block.chainid,
      address(this),
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig,
      true
    );
  }

  /// @inheritdoc IConfigurator
  function setStagingConfig(
    bytes32 configId,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    if (onchainConfig.length < 64) revert InvalidOnchainLength(onchainConfig.length);

    // Ensure that predecessorConfigDigest is set and corresponds to an
    // existing production instance
    uint256 version;
    bytes32 predecessorConfigDigest;
    assembly {
      version := mload(add(onchainConfig, 32))
      predecessorConfigDigest := mload(add(onchainConfig, 64))
    }
    if (version < MIN_SUPPORTED_ONCHAIN_CONFIG_VERSION) revert UnsupportedOnchainConfigVersion(version);

    ConfigurationState memory configurationState = s_configurationStates[configId];
    if (
      predecessorConfigDigest == bytes32(0) ||
      predecessorConfigDigest !=
      s_configurationStates[configId].configDigest[configurationState.isGreenProduction ? 1 : 0]
    ) revert InvalidPredecessorConfigDigest(predecessorConfigDigest);

    _setConfig(
      configId,
      block.chainid,
      address(this),
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig,
      false
    );
  }

  /// @inheritdoc IConfigurator
  // This will trigger the following:
  // - Offchain ShouldRetireCache will start returning true for the old (production)
  //   protocol instance
  // - Once the old production instance retires it will generate a handover
  //   retirement report
  // - The staging instance will become the new production instance once
  //   any honest oracle that is on both instances forward the retirement
  //   report from the old instance to the new instance via the
  //   PredecessorRetirementReportCache
  //
  // Note: the promotion flow only works if the previous production instance
  // is working correctly & generating reports. If that's not the case, the
  // owner is expected to "setProductionConfig" directly instead. This will
  // cause "gaps" to be created, but that seems unavoidable in such a scenario.
  function promoteStagingConfig(bytes32 configId, bool isGreenProduction) external onlyOwner {
    ConfigurationState storage configurationState = s_configurationStates[configId];
    if (isGreenProduction != configurationState.isGreenProduction)
      revert IsGreenProductionMustMatchContractState(configId, !isGreenProduction);
    if (configurationState.configCount == 0) revert ConfigUnset(configId);
    if (configurationState.configDigest[isGreenProduction ? 0 : 1] == bytes32(0))
      revert ConfigUnsetStaging(configId, isGreenProduction);
    bytes32 retiredConfigDigest = configurationState.configDigest[isGreenProduction ? 1 : 0];
    if (retiredConfigDigest == bytes32(0)) revert ConfigUnsetProduction(configId, isGreenProduction);

    configurationState.isGreenProduction = !isGreenProduction; // flip blue<->green
    emit PromoteStagingConfig(configId, retiredConfigDigest, !isGreenProduction);
  }

  /// @notice Sets config based on the given arguments
  /// @param configId config ID to set config for
  /// @param sourceChainId Chain ID of source config
  /// @param sourceAddress Address of source config Verifier
  /// @param signers addresses with which oracles sign the reports
  /// @param offchainTransmitters CSA key for the ith Oracle
  /// @param f number of faulty oracles the system can tolerate
  /// @param onchainConfig serialized configuration used by the contract (and possibly oracles)
  /// @param offchainConfigVersion version number for offchainEncoding schema
  /// @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
  function _setConfig(
    bytes32 configId,
    uint256 sourceChainId,
    address sourceAddress,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    bool isProduction
  ) internal {
    ConfigurationState storage configurationState = s_configurationStates[configId];

    uint64 newConfigCount = ++configurationState.configCount;

    bytes32 configDigest = _configDigestFromConfigData(
      configId,
      sourceChainId,
      sourceAddress,
      newConfigCount,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    if (isProduction) {
      emit ProductionConfigSet(
        configId,
        configurationState.latestConfigBlockNumber,
        configDigest,
        newConfigCount,
        signers,
        offchainTransmitters,
        f,
        onchainConfig,
        offchainConfigVersion,
        offchainConfig,
        configurationState.isGreenProduction
      );
      s_configurationStates[configId].configDigest[configurationState.isGreenProduction ? 1 : 0] = configDigest;
    } else {
      emit StagingConfigSet(
        configId,
        configurationState.latestConfigBlockNumber,
        configDigest,
        newConfigCount,
        signers,
        offchainTransmitters,
        f,
        onchainConfig,
        offchainConfigVersion,
        offchainConfig,
        configurationState.isGreenProduction
      );
      s_configurationStates[configId].configDigest[configurationState.isGreenProduction ? 0 : 1] = configDigest;
    }

    configurationState.latestConfigBlockNumber = uint32(block.number);
  }

  /// @notice Generates the config digest from config data
  /// @param configId config ID to set config for
  /// @param sourceChainId Chain ID of configurator contract
  /// @param sourceAddress Address of configurator contract
  /// @param configCount ordinal number of this config setting among all config settings over the life of this contract
  /// @param signers ith element is address ith oracle uses to sign a report
  /// @param offchainTransmitters ith element is address ith oracle used to transmit reports (in this case used for flexible additional field, such as CSA pub keys)
  /// @param f maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly
  /// @param onchainConfig serialized configuration used by the contract (and possibly oracles)
  /// @param offchainConfigVersion version of the serialization format used for "offchainConfig" parameter
  /// @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
  /// @dev This function is a modified version of the method from OCR2Abstract
  function _configDigestFromConfigData(
    bytes32 configId,
    uint256 sourceChainId,
    address sourceAddress,
    uint64 configCount,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          configId,
          sourceChainId,
          sourceAddress,
          configCount,
          signers,
          offchainTransmitters,
          f,
          onchainConfig,
          offchainConfigVersion,
          offchainConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    // 0x0009 corresponds to ConfigDigestPrefixLLO in libocr
    uint256 prefix = 0x0009 << (256 - 16); // 0x000900..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool isVerifier) {
    return interfaceId == type(IConfigurator).interfaceId;
  }

  /// @inheritdoc ITypeAndVersion
  function typeAndVersion() external pure override returns (string memory) {
    return "Configurator 0.5.0";
  }

  modifier checkConfigValid(uint256 numSigners, uint256 f) {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
    if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
    _;
  }
}
