// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IConfigurator} from "../../interfaces/IConfigurator.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/**
 * @title Configurator
 * @author samsondav
 * @notice This contract is intended to be deployed on the source chain and acts as a OCR3 configurator for LLO/Mercury
 **/

contract Configurator is IConfigurator, ConfirmedOwner, TypeAndVersionInterface, IERC165 {
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

  struct ConfigurationState {
    // The number of times a new configuration
    // has been set
    uint64 configCount;
    // The block number of the block the last time
    /// the configuration was updated.
    uint32 latestConfigBlockNumber;
  }

  constructor() ConfirmedOwner(msg.sender) {}

  /// @notice Configuration states keyed on DON ID
  mapping(bytes32 => ConfigurationState) internal s_configurationStates;

  /// @inheritdoc IConfigurator
  function setConfig(
    bytes32 donId,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(
      donId,
      block.chainid,
      address(this),
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  /// @notice Sets config based on the given arguments
  /// @param donId DON ID to set config for
  /// @param sourceChainId Chain ID of source config
  /// @param sourceAddress Address of source config Verifier
  /// @param signers addresses with which oracles sign the reports
  /// @param offchainTransmitters CSA key for the ith Oracle
  /// @param f number of faulty oracles the system can tolerate
  /// @param onchainConfig serialized configuration used by the contract (and possibly oracles)
  /// @param offchainConfigVersion version number for offchainEncoding schema
  /// @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
  function _setConfig(
    bytes32 donId,
    uint256 sourceChainId,
    address sourceAddress,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal {
    ConfigurationState storage configurationState = s_configurationStates[donId];

    uint64 newConfigCount = ++configurationState.configCount;

    bytes32 configDigest = _configDigestFromConfigData(
      donId,
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

    emit ConfigSet(
      donId,
      configurationState.latestConfigBlockNumber,
      configDigest,
      newConfigCount,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    configurationState.latestConfigBlockNumber = uint32(block.number);
  }

  /// @notice Generates the config digest from config data
  /// @param donId DON ID to set config for
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
    bytes32 donId,
    uint256 sourceChainId,
    address sourceAddress,
    uint64 configCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          donId,
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
    // 0x0006 corresponds to ConfigDigestPrefixLLO in libocr
    uint256 prefix = 0x0009 << (256 - 16); // 0x000900..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool isVerifier) {
    return interfaceId == type(IConfigurator).interfaceId;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "Configurator 0.4.0";
  }

  modifier checkConfigValid(uint256 numSigners, uint256 f) {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
    if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
    _;
  }
}
