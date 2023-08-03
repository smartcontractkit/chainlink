// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITermsOfServiceAllowList} from "./interfaces/ITermsOfServiceAllowList.sol";
import {Routable} from "../Routable.sol";
import {IAccessController} from "../../../../shared/interfaces/IAccessController.sol";
import {ITypeAndVersion} from "../../../../shared/interfaces/ITypeAndVersion.sol";

import {Address} from "../../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/Address.sol";
import {EnumerableSet} from "../../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/structs/EnumerableSet.sol";

/**
 * @notice A contract to handle access control of subscription management dependent on signing a Terms of Service
 */
contract TermsOfServiceAllowList is Routable, ITermsOfServiceAllowList, IAccessController {
  using Address for address;
  using EnumerableSet for EnumerableSet.AddressSet;

  EnumerableSet.AddressSet private s_allowedSenders;
  mapping(address => bool) private s_blockedSenders;

  error InvalidSignature();
  error InvalidUsage();
  error RecipientIsBlocked();

  // ================================================================
  // |                     Configuration state                      |
  // ================================================================

  struct Config {
    bool enabled;
    address signerPublicKey;
  }

  Config private s_config;

  event ConfigSet(bool enabled);

  // ================================================================
  // |                       Initialization                         |
  // ================================================================

  constructor(address router, bytes memory config) Routable(router, config) {}

  // ================================================================
  // |                    Configuration methods                     |
  // ================================================================

  /**
   * @notice Sets the configuration
   * @param config bytes of config data to set the following:
   *  - enabled: boolean representing if the allow list is active, when disabled all usage will be allowed
   *  - signerPublicKey: public key of the signer of the proof
   */
  function _updateConfig(bytes memory config) internal override {
    (bool enabled, address signerPublicKey) = abi.decode(config, (bool, address));
    s_config = Config({enabled: enabled, signerPublicKey: signerPublicKey});
    emit ConfigSet(enabled);
  }

  /**
   * @inheritdoc ITypeAndVersion
   */
  string public constant override typeAndVersion = "Functions Terms of Service Allow List v1.0.0";

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
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function unblockSender(address sender) external override onlyRouterOwner {
    s_blockedSenders[sender] = false;
  }
}
