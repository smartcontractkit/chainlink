// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMN} from "../../../interfaces/IRMN.sol";
import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";

import {Internal} from "../../../libraries/Internal.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract RMNRemoteSetup is BaseTest {
  RMNRemote public s_rmnRemote;
  address public OFF_RAMP_ADDRESS;

  RMNRemote.Signer[] public s_signers;
  Vm.Wallet[] public s_signerWallets;

  Internal.MerkleRoot[] internal s_merkleRoots;
  IRMNRemote.Signature[] internal s_signatures;

  bytes16 internal constant CURSE_SUBJ_1 = bytes16(keccak256("subject 1"));
  bytes16 internal constant CURSE_SUBJ_2 = bytes16(keccak256("subject 2"));
  bytes16[] internal s_curseSubjects;

  IRMN internal s_legacyRMN = IRMN(makeAddr("legacyRMN"));

  function setUp() public virtual override {
    super.setUp();
    s_rmnRemote = new RMNRemote(1, s_legacyRMN);
    OFF_RAMP_ADDRESS = makeAddr("OFF RAMP");
    s_curseSubjects = [CURSE_SUBJ_1, CURSE_SUBJ_2];

    _setupSigners(10);
  }

  /// @notice sets up a list of signers with strictly increasing onchain public keys
  /// @dev signers do not have to be in order when configured, but they do when generating signatures
  /// rather than sort signers every time, we do it once here and store the sorted list
  function _setupSigners(
    uint256 numSigners
  ) internal {
    // remove any existing config
    while (s_signerWallets.length > 0) {
      s_signerWallets.pop();
    }
    while (s_signers.length > 0) {
      s_signers.pop();
    }

    for (uint256 i = 0; i < numSigners; ++i) {
      s_signerWallets.push(vm.createWallet(vm.randomUint()));
    }

    _sort(s_signerWallets);

    for (uint256 i = 0; i < numSigners; ++i) {
      s_signers.push(RMNRemote.Signer({onchainPublicKey: s_signerWallets[i].addr, nodeIndex: uint64(i)}));
    }
  }

  /// @notice generates n merkleRoots and matching valid signatures and populates them into
  /// the shared storage vars
  function _generatePayloadAndSigs(uint256 numUpdates, uint256 numSigs) internal {
    vm.assertTrue(numUpdates > 0, "need at least 1 dest lane update");
    vm.assertTrue(numSigs <= s_signerWallets.length, "cannot generate more sigs than signers");

    // remove any existing merkleRoots and sigs
    while (s_merkleRoots.length > 0) {
      s_merkleRoots.pop();
    }
    while (s_signatures.length > 0) {
      s_signatures.pop();
    }

    for (uint256 i = 0; i < numUpdates; i++) {
      s_merkleRoots.push(_generateRandomDestLaneUpdate());
    }

    for (uint256 i = 0; i < numSigs; i++) {
      s_signatures.push(_signDestLaneUpdate(s_merkleRoots, s_signerWallets[i]));
    }
  }

  /// @notice generates a random dest lane update
  function _generateRandomDestLaneUpdate() private returns (Internal.MerkleRoot memory) {
    uint64 minSeqNum = uint32(vm.randomUint());
    uint64 maxSeqNum = minSeqNum + 100;
    return Internal.MerkleRoot({
      sourceChainSelector: uint64(vm.randomUint()),
      onRampAddress: abi.encode(vm.randomAddress()),
      minSeqNr: minSeqNum,
      maxSeqNr: maxSeqNum,
      merkleRoot: _randomBytes32()
    });
  }

  /// @notice signs the provided payload with the provided wallet
  /// @return sig the signature
  function _signDestLaneUpdate(
    Internal.MerkleRoot[] memory merkleRoots,
    Vm.Wallet memory wallet
  ) private returns (IRMNRemote.Signature memory) {
    (, RMNRemote.Config memory config) = s_rmnRemote.getVersionedConfig();
    bytes32 digest = keccak256(
      abi.encode(
        s_rmnRemote.getReportDigestHeader(),
        RMNRemote.Report({
          destChainId: block.chainid,
          destChainSelector: s_rmnRemote.getLocalChainSelector(),
          rmnRemoteContractAddress: address(s_rmnRemote),
          offrampAddress: OFF_RAMP_ADDRESS,
          rmnHomeContractConfigDigest: config.rmnHomeContractConfigDigest,
          merkleRoots: merkleRoots
        })
      )
    );
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(wallet, digest);
    // RMNRemote only supports sigs with v=27, so adjust if necessary
    // Any valid ECDSA sig (r, s, v) can be "flipped" into (r, s*, v*) without knowing the private key (where v=27 or 28 for secp256k1)
    // https://github.com/kadenzipfel/smart-contract-vulnerabilities/blob/master/vulnerabilities/signature-malleability.md
    if (v == 28) {
      uint256 N = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141;
      s = bytes32(N - uint256(s));
    }
    return (IRMNRemote.Signature({r: r, s: s}));
  }

  /// @notice bubble sort on a storage array of wallets
  function _sort(
    Vm.Wallet[] storage wallets
  ) private {
    bool swapped;
    for (uint256 i = 1; i < wallets.length; i++) {
      swapped = false;
      for (uint256 j = 0; j < wallets.length - i; j++) {
        Vm.Wallet memory next = wallets[j + 1];
        Vm.Wallet memory actual = wallets[j];
        if (next.addr < actual.addr) {
          wallets[j] = next;
          wallets[j + 1] = actual;
          swapped = true;
        }
      }
      if (!swapped) {
        return;
      }
    }
  }

  /// @dev returns a pseudo-random bytes32
  function _randomBytes32() internal returns (bytes32) {
    return bytes32(vm.randomUint());
  }
}
