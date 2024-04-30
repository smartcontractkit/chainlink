// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITermsOfServiceAllowList, TermsOfServiceAllowListConfig} from "./interfaces/ITermsOfServiceAllowList.sol";
import {IAccessController} from "../../../../shared/interfaces/IAccessController.sol";
import {ITypeAndVersion} from "../../../../shared/interfaces/ITypeAndVersion.sol";

import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";

import {Address} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/Address.sol";
import {EnumerableSet} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @notice A contract to handle access control of subscription management dependent on signing a Terms of Service
contract TermsOfServiceAllowList is ITermsOfServiceAllowList, IAccessController, ITypeAndVersion, ConfirmedOwner {
  using Address for address;
  using EnumerableSet for EnumerableSet.AddressSet;

  /// @inheritdoc ITypeAndVersion
  string public constant override typeAndVersion = "Functions Terms of Service Allow List v1.1.0";

  EnumerableSet.AddressSet private s_allowedSenders;
  EnumerableSet.AddressSet private s_blockedSenders;

  event AddedAccess(address user);
  event BlockedAccess(address user);
  event UnblockedAccess(address user);

  error InvalidSignature();
  error InvalidUsage();
  error RecipientIsBlocked();
  error InvalidCalldata();

  TermsOfServiceAllowListConfig private s_config;

  event ConfigUpdated(TermsOfServiceAllowListConfig config);

  // ================================================================
  // |                       Initialization                         |
  // ================================================================

  constructor(
    TermsOfServiceAllowListConfig memory config,
    address[] memory initialAllowedSenders,
    address[] memory initialBlockedSenders
  ) ConfirmedOwner(msg.sender) {
    updateConfig(config);

    for (uint256 i = 0; i < initialAllowedSenders.length; ++i) {
      s_allowedSenders.add(initialAllowedSenders[i]);
    }

    for (uint256 j = 0; j < initialBlockedSenders.length; ++j) {
      if (s_allowedSenders.contains(initialBlockedSenders[j])) {
        // Allowed senders cannot also be blocked
        revert InvalidCalldata();
      }
      s_blockedSenders.add(initialBlockedSenders[j]);
    }
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  /// @notice Gets the contracts's configuration
  /// @return config
  function getConfig() external view returns (TermsOfServiceAllowListConfig memory) {
    return s_config;
  }

  /// @notice Sets the contracts's configuration
  /// @param config - See the contents of the TermsOfServiceAllowListConfig struct in ITermsOfServiceAllowList.sol for more information
  function updateConfig(TermsOfServiceAllowListConfig memory config) public onlyOwner {
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
    if (s_blockedSenders.contains(recipient)) {
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
    if (s_allowedSenders.add(recipient)) {
      emit AddedAccess(recipient);
    }
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function getAllAllowedSenders() external view override returns (address[] memory) {
    return s_allowedSenders.values();
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function getAllowedSendersCount() external view override returns (uint64) {
    return uint64(s_allowedSenders.length());
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function getAllowedSendersInRange(
    uint64 allowedSenderIdxStart,
    uint64 allowedSenderIdxEnd
  ) external view override returns (address[] memory allowedSenders) {
    if (allowedSenderIdxStart > allowedSenderIdxEnd || allowedSenderIdxEnd >= s_allowedSenders.length()) {
      revert InvalidCalldata();
    }

    allowedSenders = new address[]((allowedSenderIdxEnd - allowedSenderIdxStart) + 1);
    for (uint256 i = 0; i <= allowedSenderIdxEnd - allowedSenderIdxStart; ++i) {
      allowedSenders[i] = s_allowedSenders.at(uint256(allowedSenderIdxStart + i));
    }

    return allowedSenders;
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
    return s_blockedSenders.contains(sender);
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function blockSender(address sender) external override onlyOwner {
    s_allowedSenders.remove(sender);
    s_blockedSenders.add(sender);
    emit BlockedAccess(sender);
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function unblockSender(address sender) external override onlyOwner {
    s_blockedSenders.remove(sender);
    emit UnblockedAccess(sender);
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function getBlockedSendersCount() external view override returns (uint64) {
    return uint64(s_blockedSenders.length());
  }

  /// @inheritdoc ITermsOfServiceAllowList
  function getBlockedSendersInRange(
    uint64 blockedSenderIdxStart,
    uint64 blockedSenderIdxEnd
  ) external view override returns (address[] memory blockedSenders) {
    if (
      blockedSenderIdxStart > blockedSenderIdxEnd ||
      blockedSenderIdxEnd >= s_blockedSenders.length() ||
      s_blockedSenders.length() == 0
    ) {
      revert InvalidCalldata();
    }

    blockedSenders = new address[]((blockedSenderIdxEnd - blockedSenderIdxStart) + 1);
    for (uint256 i = 0; i <= blockedSenderIdxEnd - blockedSenderIdxStart; ++i) {
      blockedSenders[i] = s_blockedSenders.at(uint256(blockedSenderIdxStart + i));
    }

    return blockedSenders;
  }
}
