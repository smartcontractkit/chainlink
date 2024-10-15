// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {OCR2Abstract} from "./ocr/OCR2Abstract.sol";

/// @notice OCR2Base provides config management compatible with OCR3
contract OCR3Capability is OwnerIsCreator, OCR2Abstract {
  error InvalidConfig(string message);
  error ReportingUnsupported();

  string public constant override typeAndVersion = "Keystone 1.0.0";

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest, to prevent replay attacks.
  uint32 internal s_configCount;
  // makes it easier for offchain systems to extract config from logs.
  uint32 internal s_latestConfigBlockNumber;

  // Storing these fields used on the hot path in a ConfigInfo variable reduces the
  // retrieval of all of them to a single SLOAD. If any further fields are
  // added, make sure that storage of the struct still takes at most 32 bytes.
  struct ConfigInfo {
    bytes32 latestConfigDigest;
    uint8 f; // TODO: could be optimized by squeezing into one slot
    uint8 n;
  }
  ConfigInfo internal s_configInfo;

  // Reverts transaction if config args are invalid
  modifier checkConfigValid(uint256 numSigners, uint256 numTransmitters, uint256 f) {
    if (numSigners > MAX_NUM_ORACLES) revert InvalidConfig("too many signers");
    if (f == 0) revert InvalidConfig("f must be positive");
    if (numSigners != numTransmitters) revert InvalidConfig("oracle addresses out of registration");
    if (numSigners <= 3 * f) revert InvalidConfig("faulty-oracle f too high");
    _;
  }

  /// @notice sets offchain reporting protocol configuration incl. participating oracles
  /// @param _signers addresses with which oracles sign the reports
  /// @param _transmitters addresses oracles use to transmit the reports
  /// @param _f number of faulty oracles the system can tolerate
  /// @param _onchainConfig encoded on-chain contract configuration
  /// @param _offchainConfigVersion version number for offchainEncoding schema
  /// @param _offchainConfig encoded off-chain oracle configuration
  /// @dev signer = [ 1 byte type | 2 byte len | n byte value ]...
  function setConfig(
    bytes[] calldata _signers,
    address[] calldata _transmitters,
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _offchainConfigVersion,
    bytes memory _offchainConfig
  ) external override checkConfigValid(_signers.length, _transmitters.length, _f) onlyOwner {
    // Bounded by MAX_NUM_ORACLES in OCR2Abstract.sol
    for (uint256 i = 0; i < _signers.length; ++i) {
      if (_transmitters[i] == address(0)) revert InvalidConfig("transmitter must not be empty");
      // add new signers
      bytes calldata publicKeys = _signers[i];
      uint256 offset = 0;
      uint256 publicKeysLength = uint16(publicKeys.length);
      // scan through public keys to validate encoded format
      while (offset < publicKeysLength) {
        if (offset + 3 > publicKeysLength) revert InvalidConfig("invalid signer pubKey encoding");
        uint256 keyLen = uint256(uint8(publicKeys[offset + 1])) + (uint256(uint8(publicKeys[offset + 2])) << 8);
        if (offset + 3 + keyLen > publicKeysLength) revert InvalidConfig("invalid signer pubKey encoding");
        offset += 3 + keyLen;
      }
    }
    s_configInfo.f = _f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    {
      s_configInfo.latestConfigDigest = _configDigestFromConfigData(
        block.chainid,
        address(this),
        s_configCount,
        _signers,
        _transmitters,
        _f,
        _onchainConfig,
        _offchainConfigVersion,
        _offchainConfig
      );
    }
    s_configInfo.n = uint8(_signers.length);

    emit ConfigSet(
      previousConfigBlockNumber,
      s_configInfo.latestConfigDigest,
      s_configCount,
      _signers,
      _transmitters,
      _f,
      _onchainConfig,
      _offchainConfigVersion,
      _offchainConfig
    );
  }

  function _configDigestFromConfigData(
    uint256 _chainId,
    address _contractAddress,
    uint64 _configCount,
    bytes[] calldata _signers,
    address[] calldata _transmitters,
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _encodedConfigVersion,
    bytes memory _encodedConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          _chainId,
          _contractAddress,
          _configCount,
          _signers,
          _transmitters,
          _f,
          _onchainConfig,
          _encodedConfigVersion,
          _encodedConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x000e << (256 - 16); // 0x000e00..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  /// @notice information about current offchain reporting protocol configuration
  /// @return configCount ordinal number of current config, out of all configs applied to this contract so far
  /// @return blockNumber block at which this config was set
  /// @return configDigest domain-separation tag for current config (see __configDigestFromConfigData)
  function latestConfigDetails()
    external
    view
    override
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest)
  {
    return (s_configCount, s_latestConfigBlockNumber, s_configInfo.latestConfigDigest);
  }

  function transmit(
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes32[3] calldata /* reportContext */,
    bytes calldata /* report */,
    bytes32[] calldata /* rs */,
    bytes32[] calldata /* ss */,
    bytes32 /* rawVs */ // signatures
  ) external pure override {
    revert ReportingUnsupported();
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDigestAndEpoch()
    external
    view
    virtual
    override
    returns (bool scanLogs, bytes32 configDigest, uint32 epoch)
  {
    return (true, bytes32(0), uint32(0));
  }
}
