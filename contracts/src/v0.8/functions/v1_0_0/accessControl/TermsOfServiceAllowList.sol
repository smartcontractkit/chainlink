// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITermsOfServiceAllowList} from "./interfaces/ITermsOfServiceAllowList.sol";
import {IAccessController} from "../../../shared/interfaces/IAccessController.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";

import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";

import {Address} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/Address.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @notice A contract to handle access control of subscription management dependent on signing a Terms of Service
contract TermsOfServiceAllowList is ITermsOfServiceAllowList, IAccessController, ITypeAndVersion, ConfirmedOwner {
  using Address for address;
  using EnumerableSet for EnumerableSet.AddressSet;

  /// @inheritdoc ITypeAndVersion
  string public constant override typeAndVersion = "Functions Terms of Service Allow List v1.0.0";

  EnumerableSet.AddressSet private s_allowedSenders;
  mapping(address => bool) private s_blockedSenders;

  event AddedAccess(address user);
  event BlockedAccess(address user);
  event UnblockedAccess(address user);

  error InvalidSignature();
  error InvalidUsage();
  error RecipientIsBlocked();

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================
  struct Config {
    bool enabled; // ═════════════╗ When enabled, access will be checked against s_allowedSenders. When disabled, all access will be allowed.
    address signerPublicKey; // ══╝ The key pair that needs to sign the acceptance data
  }

  Config private s_config;

  event ConfigUpdated(Config config);

  // ================================================================
  // |                       Initialization                         |
  // ================================================================

  constructor(Config memory config) ConfirmedOwner(msg.sender) {
    updateConfig(config);
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  /// @notice Gets the contracts's configuration
  /// @return config
  function getConfig() external view returns (Config memory) {
    return s_config;
  }

  /// @notice Sets the contracts's configuration
  /// @param config - See the contents of the TermsOfServiceAllowList.Config struct for more information
  function updateConfig(Config memory config) public onlyOwner {
    s_config = config;
    emit ConfigUpdated(config);
  }

  // ================================================================
  // |                      Allow methods                           |
  // ================================================================

  /// @inheritdoc ITermsOfServiceAllowList
  function getMessage(address acceptor, address recipient) public pure override returns (bytes32) {
    return keccak256(abi.encodePacked(acceptor, recipient));
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function acceptTermsOfService(address acceptor, address recipient, bytes32 r, bytes32 s, uint8 v) external override {
    if (s_blockedSenders[recipient]) {
      revert RecipientIsBlocked();
    }

    // Validate that the signature is correct and the correct data has been signed
    bytes32 prefixedMessage = keccak256(
      abi.encodePacked("\x19Ethereum Signed Message:\n32", getMessage(acceptor, recipient))
    );
    if (ecrecover(prefixedMessage, v, r, s) != s_config.signerPublicKey) {
      revert InvalidSignature();
    }

    // If contract, validate that msg.sender == recipient
    // This is to prevent EoAs from claiming contracts that they are not in control of
    // If EoA, validate that msg.sender == acceptor == recipient
    // This is to prevent EoAs from accepting for other EoAs
    if (msg.sender != recipient || (msg.sender != acceptor && !msg.sender.isContract())) {
      revert InvalidUsage();
    }

    // Add recipient to the allow list
    s_allowedSenders.add(recipient);
    emit AddedAccess(recipient);
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function getAllAllowedSenders() external view override returns (address[] memory) {
    return s_allowedSenders.values();
  }

  /// @inheritdoc IAccessController
  function hasAccess(address user, bytes calldata /* data */) external view override returns (bool) {
    if (!s_config.enabled) {
      return true;
    }
    return s_allowedSenders.contains(user);
  }

  // ================================================================
  // |                         Block methods                        |
  // ================================================================

  /// @inheritdoc ITermsOfServiceAllowList
  function isBlockedSender(address sender) external view override returns (bool) {
    if (!s_config.enabled) {
      return false;
    }
    return s_blockedSenders[sender];
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function blockSender(address sender) external override onlyOwner {
    s_allowedSenders.remove(sender);
    s_blockedSenders[sender] = true;
    emit BlockedAccess(sender);
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function unblockSender(address sender) external override onlyOwner {
    s_blockedSenders[sender] = false;
    emit UnblockedAccess(sender);
  }
}
