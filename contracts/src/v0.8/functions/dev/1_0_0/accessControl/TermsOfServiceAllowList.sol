// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {EnumerableSet} from "../../../../shared/vendor/openzeppelin-solidity/v.4.8.0/contracts/utils/structs/EnumerableSet.sol";
import {ITermsOfServiceAllowList} from "./interfaces/ITermsOfServiceAllowList.sol";
import {Routable, ITypeAndVersion} from "../Routable.sol";
import {IOwnable} from "../../../../shared/interfaces/IOwnable.sol";

/**
 * @notice A contract to handle access control of subscription management dependent on signing a Terms of Service
 */

contract TermsOfServiceAllowList is Routable, ITermsOfServiceAllowList {
  using EnumerableSet for EnumerableSet.AddressSet;

  EnumerableSet.AddressSet private s_allowedSenders;
  EnumerableSet.AddressSet private s_blockedSenders;

  error InvalidProof();
  error RecipientIsBlocked();
  // ================================================================
  // |                     Configuration state                      |
  // ================================================================
  struct Config {
    bool enabled;
    address proofSignerPublicKey;
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
   *  - proofSignerPublicKey: public key of the signer of the proof
   */
  function _setConfig(bytes memory config) internal override {
    (bool enabled, address proofSignerPublicKey) = abi.decode(config, (bool, address));
    s_config = Config({enabled: enabled, proofSignerPublicKey: proofSignerPublicKey});
    emit ConfigSet(enabled);
  }

  /**
   * @inheritdoc ITypeAndVersion
   */
  function typeAndVersion() public pure override returns (string memory) {
    return "Functions Terms of Service Allow List v1";
  }

  // ================================================================
  // |                  Terms of Service methods                    |
  // ================================================================
  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function getMessageHash(address acceptor, address recipient) public pure override returns (bytes32) {
    return keccak256(abi.encodePacked(acceptor, recipient));
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function getEthSignedMessageHash(bytes32 messageHash) public pure override returns (bytes32) {
    return keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash));
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function acceptTermsOfService(address acceptor, address recipient, bytes calldata proof) external override {
    if (s_blockedSenders.contains(recipient)) {
      revert RecipientIsBlocked();
    }

    // Validate that the proof is correct and has been signed
    bytes32 messageHash = getMessageHash(acceptor, recipient);
    bytes32 ethSignedMessageHash = getEthSignedMessageHash(messageHash);
    address proofSigner = _getSigner(ethSignedMessageHash, proof);
    if (proofSigner != s_config.proofSignerPublicKey) {
      revert InvalidProof();
    }

    // Check if msg.sender is an EoA or contract
    bool callerIsContractAccount = _isContract(msg.sender);
    // If EoA, validate that msg.sender == acceptor == recipient
    // This is to prevent EoAs from accepting for other EoAs
    if (callerIsContractAccount == false && (msg.sender != acceptor || msg.sender != recipient)) {
      revert InvalidProof();
    }

    // If contract, validate that msg.sender == recipient
    // This is to prevent EoAs from claiming contracts that they are not in control of
    if (callerIsContractAccount && msg.sender != recipient) {
      revert InvalidProof();
    }

    // Add recipient to the allow list
    s_allowedSenders.add(recipient);
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function isAllowedSender(address sender) public view override returns (bool) {
    if (s_config.enabled == false) {
      return true;
    }
    return s_allowedSenders.contains(sender);
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function isBlockedSender(address sender) public view override returns (bool) {
    if (s_config.enabled == false) {
      return false;
    }
    return s_blockedSenders.contains(sender);
  }

  // ================================================================
  // |                     Owner methods                          |
  // ================================================================
  /**

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function blockSender(address sender) external override onlyRouterOwner {
    s_allowedSenders.remove(sender);
    s_blockedSenders.add(sender);
  }

  /**
   * @inheritdoc ITermsOfServiceAllowList
   */
  function unblockSender(address sender) external override onlyRouterOwner {
    s_blockedSenders.remove(sender);
  }

  // ================================================================
  // |                     Internal helpers                          |
  // ================================================================

  function _isContract(address _addr) private view returns (bool isContract) {
    uint32 size;
    // solhint-disable-next-line no-inline-assembly
    assembly {
      size := extcodesize(_addr)
    }
    return (size > 0);
  }

  function _getSigner(bytes32 _ethSignedMessageHash, bytes memory signature) private pure returns (address) {
    (bytes32 r, bytes32 s, uint8 v) = _splitSignature(signature);
    return ecrecover(_ethSignedMessageHash, v, r, s);
  }

  function _splitSignature(bytes memory sig) private pure returns (bytes32 r, bytes32 s, uint8 v) {
    if (sig.length != 65) {
      revert InvalidProof();
    }
    // solhint-disable-next-line no-inline-assembly
    assembly {
      /*/ 
        First 32 bytes stores the length of the signature

        add(sig, 32) = pointer of sig + 32
        effectively, skips first 32 bytes of signature

        mload(p) loads next 32 bytes starting at the memory address p into memory
      */
      // first 32 bytes, after the length prefix
      r := mload(add(sig, 32))
      // second 32 bytes
      s := mload(add(sig, 64))
      // final byte (first byte of the next 32 bytes)
      v := byte(0, mload(add(sig, 96)))
    }
  }
}
