// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {OCR2Abstract} from "./OCR2Abstract.sol";

/**
 * @notice Onchain verification of reports from the offchain reporting protocol
 * @dev For details on its operation, see the offchain reporting protocol design
 * doc, which refers to this contract as simply the "contract".
 */
abstract contract OCR2Base is ConfirmedOwner, OCR2Abstract {
  error InvalidConfig(string message);

  constructor() ConfirmedOwner(msg.sender) {}

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest, to prevent replay attacks.
  uint32 internal s_configCount;
  uint32 internal s_latestConfigBlockNumber; // makes it easier for offchain systems
  // to extract config from logs.

  // Storing these fields used on the hot path in a ConfigInfo variable reduces the
  // retrieval of all of them to a single SLOAD. If any further fields are
  // added, make sure that storage of the struct still takes at most 32 bytes.
  struct ConfigInfo {
    bytes32 latestConfigDigest;
    uint8 f; // TODO: could be optimized by squeezing into one slot
    uint8 n;
  }
  ConfigInfo internal s_configInfo;

  /// @notice triggers a new run of the offchain reporting protocol
  /// @param previousConfigBlockNumber block in which the previous config was set, to simplify historic analysis
  /// @param configDigest configDigest of this configuration
  /// @param configCount ordinal number of this config setting among all config settings over the life of this contract
  /// @param signers ith element is address ith oracle uses to sign a report
  /// @param transmitters ith element is address ith oracle uses to transmit a report via the transmit method
  /// @param f maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly
  /// @param onchainConfig serialized configuration used by the contract (and possibly oracles)
  /// @param offchainConfigVersion version of the serialization format used for "offchainConfig" parameter
  /// @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    bytes[] signers,
    address[] transmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );

  // s_signers contains the signing address of each oracle
  bytes[] internal s_signers;

  /*
   * Config logic
   */

  // Reverts transaction if config args are invalid
  modifier checkConfigValid(
    uint256 numSigners,
    uint256 numTransmitters,
    uint256 f
  ) {
    if (numSigners > MAX_NUM_ORACLES) revert InvalidConfig("too many signers");
    if (f == 0) revert InvalidConfig("f must be positive");
    if (numSigners <= 3 * f) revert InvalidConfig("faulty-oracle f too high");
    _;
  }

  // solhint-disable-next-line gas-struct-packing
  struct SetConfigArgs {
    bytes[] signers;
    address[] transmitters;
    uint8 f;
    bytes onchainConfig;
    uint64 offchainConfigVersion;
    bytes offchainConfig;
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

  // signer = [ 1 byte type | 2 byte len | n byte value ]...

  function setConfig(
    address[] memory _signers,
    address[] memory _transmitters,
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _offchainConfigVersion,
    bytes memory _offchainConfig
  ) external override checkConfigValid(_signers.length, _transmitters.length, _f) onlyOwner {}

  /**
   * @notice sets offchain reporting protocol configuration incl. participating oracles
   * @param _signers addresses with which oracles sign the reports
   * @param _transmitters addresses oracles use to transmit the reports
   * @param _f number of faulty oracles the system can tolerate
   * @param _onchainConfig encoded on-chain contract configuration
   * @param _offchainConfigVersion version number for offchainEncoding schema
   * @param _offchainConfig encoded off-chain oracle configuration
   */
  function setConfig(
    bytes[] calldata _signers,
    address[] calldata _transmitters, // TODO: remove, use fake static addrs to satisfy offchain
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _offchainConfigVersion,
    bytes memory _offchainConfig
  ) external checkConfigValid(_signers.length, _transmitters.length, _f) onlyOwner {
    SetConfigArgs memory args = SetConfigArgs({ // TODO: is this saving gas? unused fields
      signers: _signers,
      transmitters: _transmitters,
      f: _f,
      onchainConfig: _onchainConfig,
      offchainConfigVersion: _offchainConfigVersion,
      offchainConfig: _offchainConfig
    });

    while (s_signers.length != 0) {
      // remove any old signer addresses
      s_signers.pop();
    }

    // Bounded by MAX_NUM_ORACLES in OCR2Abstract.sol
    for (uint256 i = 0; i < args.signers.length; i++) {
      if (args.transmitters[i] == address(0)) revert InvalidConfig("transmitter must not be empty");
      // add new signers
      bytes calldata publicKeys = _signers[i];
      uint16 offset = 32; // skip the array length
      uint16 len = uint16(publicKeys.length);

      // chainlink/core/capabilities/ccip/ocrimpls/config_tracker.go TODO: libocr better signer validation for key subtypes checkIdentityListsHaveNoDuplicates
      // after PublicConfigFromContractConfig call

      // parse encoded public keys to validate uniqueness
      // TODO: could we just trust libocr to ensure uniqueness?
      while (offset < len) {
        uint8 keyType = uint8(publicKeys[offset]);
        uint16 keyLen = uint16(uint8(publicKeys[offset + 1])) + (uint16(uint8(publicKeys[offset + 2])) << 8);
        bytes calldata publicKey = publicKeys[offset + 3:offset + 3 + keyLen];
        offset += 3 + keyLen;
        // TODO: uniq the signers() list
        // if (s_oracles[args.signers[i]].role != Role.Unset) revert InvalidConfig("repeated signer address");
      }
      // TODO: remove transmitters altogether?
      s_signers.push(args.signers[i]);
    }
    s_configInfo.f = args.f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    {
      s_configInfo.latestConfigDigest = _configDigestFromConfigData(
        block.chainid,
        address(this),
        s_configCount,
        args.signers,
        args.transmitters,
        args.f,
        args.onchainConfig,
        args.offchainConfigVersion,
        args.offchainConfig
      );
    }
    s_configInfo.n = uint8(args.signers.length);

    emit ConfigSet(
      previousConfigBlockNumber,
      s_configInfo.latestConfigDigest,
      s_configCount,
      args.signers,
      args.transmitters,
      args.f,
      args.onchainConfig,
      args.offchainConfigVersion,
      args.offchainConfig
    );
  }

  function _configDigestFromConfigData(
    uint256 _chainId,
    address _contractAddress,
    uint64 _configCount,
    bytes[] memory _signers,
    address[] memory _transmitters,
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
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  /**
   * @notice information about current offchain reporting protocol configuration
   * @return configCount ordinal number of current config, out of all configs applied to this contract so far
   * @return blockNumber block at which this config was set
   * @return configDigest domain-separation tag for current config (see __configDigestFromConfigData)
   */
  function latestConfigDetails()
    external
    view
    override
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest)
  {
    return (s_configCount, s_latestConfigBlockNumber, s_configInfo.latestConfigDigest);
  }
}
