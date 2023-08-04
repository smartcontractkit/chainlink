// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITermsOfServiceAllowList} from "./interfaces/ITermsOfServiceAllowList.sol";
import {IAccessController} from "../../../../shared/interfaces/IAccessController.sol";
import {ITypeAndVersion} from "../../../../shared/interfaces/ITypeAndVersion.sol";

import {Routable} from "../Routable.sol";
import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";

import {Address} from "../../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/Address.sol";
import {EnumerableSet} from "../../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/structs/EnumerableSet.sol";

/**
 * @notice A contract to handle access control of subscription management dependent on signing a Terms of Service
 */
contract TermsOfServiceAllowList is ITermsOfServiceAllowList, IAccessController, Routable, ConfirmedOwner {
  using Address for address;
  using EnumerableSet for EnumerableSet.AddressSet;

  // @inheritdoc ITypeAndVersion
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
  event ConfigUpdated(Config config);

  Config private s_config;

  // ================================================================
  // |                       Initialization                         |
  // ================================================================

  constructor(address router, Config memory config) Routable(router) ConfirmedOwner(msg.sender) {
    updateConfig(config);
  }

  // ================================================================
  // |                        Configuration                         |
  // ================================================================

  // @inheritdoc ITermsOfServiceAllowList
  function getConfig() external view override returns (Config memory) {
    return s_config;
  }

  // @inheritdoc ITermsOfServiceAllowList
  function updateConfig(Config memory config) public override onlyOwner {
    s_config = config;
    emit ConfigUpdated(config);
  }

  // ================================================================
  // |                  Terms of Service methods                    |
  // ================================================================

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function getMessage(address acceptor, address recipient) public pure override returns (bytes32) {
    return keccak256(abi.encodePacked(acceptor, recipient));
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
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

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function getAllAllowedSenders() external view override returns (address[] memory) {
    return s_allowedSenders.values();
  }

  /**
   * @inheritdoc IAccessController
   */
  function hasAccess(address user, bytes calldata /* data */) external view override returns (bool) {
    if (!s_config.enabled) {
      return true;
    }
    return s_allowedSenders.contains(user);
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function isBlockedSender(address sender) external view override returns (bool) {
    if (!s_config.enabled) {
      return false;
    }
    return s_blockedSenders[sender];
  }

  // ================================================================
  // |                     Owner methods                          |
  // ================================================================

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function blockSender(address sender) external override onlyRouterOwner {
    s_allowedSenders.remove(sender);
    s_blockedSenders[sender] = true;
    emit BlockedAccess(sender);
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function unblockSender(address sender) external override onlyRouterOwner {
    s_blockedSenders[sender] = false;
    emit UnblockedAccess(sender);
  }
}
