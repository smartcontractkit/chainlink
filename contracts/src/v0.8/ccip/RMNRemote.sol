// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/access/Ownable2Step.sol";

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

bytes32 constant RMN_V1_6_ANY2EVM_REPORT = keccak256("RMN_V1_6_ANY2EVM_REPORT");

/// @notice This contract supports verification of RMN reports for any Any2EVM OffRamp.
contract RMNRemote is Ownable2Step, ITypeAndVersion {
  /// @dev temp placeholder to exclude this contract from coverage
  function test() public {}

  string public constant override typeAndVersion = "RMNRemote 1.6.0-dev";

  uint64 internal immutable i_chainSelector;

  constructor(uint64 chainSelector) {
    i_chainSelector = chainSelector;
  }

  struct Signer {
    address onchainPublicKey; // for signing reports
    uint64 nodeIndex; // maps to nodes in home chain config, should be strictly increasing
  }

  struct Config {
    bytes32 rmnHomeContractConfigDigest;
    Signer[] signers;
    uint64 minSigners;
  }

  struct VersionedConfig {
    uint32 version;
    Config config;
  }

  Config s_config;
  uint32 s_configCount;

  mapping(address signer => bool exists) s_signers; // for more gas efficient verify

  function setConfig(Config calldata newConfig) external onlyOwner {
    // sanity checks
    {
      // signers are in ascending order of nodeIndex
      for (uint256 i = 1; i < newConfig.signers.length; ++i) {
        if (!(newConfig.signers[i - 1].nodeIndex < newConfig.signers[i].nodeIndex)) {
          revert InvalidSignerOrder();
        }
      }

      // minSigners is tenable
      if (!(newConfig.minSigners <= newConfig.signers.length)) {
        revert MinSignersTooHigh();
      }
    }

    // clear the old signers
    {
      Config storage oldConfig = s_config;
      while (oldConfig.signers.length > 0) {
        delete s_signers[oldConfig.signers[oldConfig.signers.length - 1].onchainPublicKey];
        oldConfig.signers.pop();
      }
    }

    // set the new signers
    {
      for (uint256 i = 0; i < newConfig.signers.length; ++i) {
        if (s_signers[newConfig.signers[i].onchainPublicKey]) {
          revert DuplicateOnchainPublicKey();
        }
        s_signers[newConfig.signers[i].onchainPublicKey] = true;
      }
    }

    s_config = newConfig;
    uint32 newConfigCount = ++s_configCount;
    emit ConfigSet(VersionedConfig({version: newConfigCount, config: newConfig}));
  }

  function getVersionedConfig() external view returns (VersionedConfig memory) {
    return VersionedConfig({version: s_configCount, config: s_config});
  }

  /// @notice The part of the LaneUpdate for a fixed destination chain and OffRamp, to avoid verbosity in Report
  struct DestLaneUpdate {
    uint64 sourceChainSelector;
    bytes onrampAddress; // generic, to support arbitrary sources; for EVM2EVM, use abi.encodePacked
    uint64 minMsgNr;
    uint64 maxMsgNr;
    bytes32 root;
  }

  struct Report {
    uint256 destChainId; // to guard against chain selector misconfiguration
    uint64 destChainSelector;
    address rmnRemoteContractAddress;
    address offrampAddress;
    bytes32 rmnHomeContractConfigDigest;
    DestLaneUpdate[] destLaneUpdates;
  }

  struct Signature {
    bytes32 r;
    bytes32 s;
  }

  /// @notice Verifies signatures of RMN nodes, on dest lane updates as provided in the CommitReport
  /// @param destLaneUpdates must be well formed, and is a representation of the CommitReport received from the oracles
  /// @param signatures must be sorted in ascending order by signer address
  /// @dev Will revert if verification fails. Needs to be called by the OffRamp for which the signatures are produced,
  /// otherwise verification will fail.
  function verify(DestLaneUpdate[] memory destLaneUpdates, Signature[] memory signatures) external view {
    if (s_configCount == 0) {
      revert ConfigNotSet();
    }

    bytes32 signedHash = keccak256(
      abi.encode(
        RMN_V1_6_ANY2EVM_REPORT,
        Report({
          destChainId: block.chainid,
          destChainSelector: i_chainSelector,
          rmnRemoteContractAddress: address(this),
          offrampAddress: msg.sender,
          rmnHomeContractConfigDigest: s_config.rmnHomeContractConfigDigest,
          destLaneUpdates: destLaneUpdates
        })
      )
    );

    uint256 numSigners = 0;
    address prevAddress = address(0);
    for (uint256 i = 0; i < signatures.length; ++i) {
      Signature memory sig = signatures[i];
      address signerAddress = ecrecover(signedHash, 27, sig.r, sig.s);
      if (signerAddress == address(0)) revert InvalidSignature();
      if (!(prevAddress < signerAddress)) revert OutOfOrderSignatures();
      if (!s_signers[signerAddress]) revert UnexpectedSigner();
      prevAddress = signerAddress;
      ++numSigners;
    }
    if (numSigners < s_config.minSigners) revert ThresholdNotMet();
  }

  ///
  /// Events
  ///

  event ConfigSet(VersionedConfig versionedConfig);

  ///
  /// Errors
  ///

  error InvalidSignature();
  error OutOfOrderSignatures();
  error UnexpectedSigner();
  error ThresholdNotMet();
  error ConfigNotSet();
  error InvalidSignerOrder();
  error MinSignersTooHigh();
  error DuplicateOnchainPublicKey();
}
