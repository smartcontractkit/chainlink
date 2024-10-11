// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMNRemote} from "../../interfaces/IRMNRemote.sol";
import {Internal} from "../../libraries/Internal.sol";
import {RMNRemote} from "../../rmn/RMNRemote.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {Vm} from "forge-std/Vm.sol";

import "forge-std/console.sol";

contract RMNRemoteSetup is BaseTest {
  RMNRemote public s_rmnRemote;
  address public OFF_RAMP_ADDRESS;

  RMNRemote.Signer[] public s_signers;
  Vm.Wallet[] public s_signerWallets;

  Internal.MerkleRoot[] s_merkleRoots;
  IRMNRemote.Signature[] s_signatures;
  uint256 internal s_v;

  bytes16 internal constant curseSubj1 = bytes16(keccak256("subject 1"));
  bytes16 internal constant curseSubj2 = bytes16(keccak256("subject 2"));
  bytes16[] internal s_curseSubjects;

  function setUp() public virtual override {
    super.setUp();
    s_rmnRemote = new RMNRemote(1);
    OFF_RAMP_ADDRESS = makeAddr("OFF RAMP");
    s_curseSubjects = [curseSubj1, curseSubj2];

    _setupSigners(10);
  }

  /// @notice sets up a list of signers with strictly increasing onchain public keys
  /// @dev signers do not have to be in order when configured, but they do when generating signatures
  /// rather than sort signers every time, we do it once here and store the sorted list
  function _setupSigners(uint256 numSigners) internal {
    // remove any existing config
    while (s_signerWallets.length > 0) {
      s_signerWallets.pop();
    }
    while (s_signers.length > 0) {
      s_signers.pop();
    }

    for (uint256 i = 0; i < numSigners; i++) {
      s_signerWallets.push(vm.createWallet(_randomNum()));
    }

    _sort(s_signerWallets);

    for (uint256 i = 0; i < numSigners; i++) {
      s_signers.push(RMNRemote.Signer({onchainPublicKey: s_signerWallets[i].addr, nodeIndex: uint64(i)}));
    }
  }

  /// @notice generates n merkleRoots and matching valid signatures and populates them into
  /// the shared storage vars
  function _generatePayloadAndSigs(uint256 numUpdates, uint256 numSigs) internal {
    require(numUpdates > 0, "need at least 1 dest lane update");
    require(numSigs <= s_signerWallets.length, "cannot generate more sigs than signers");

    // remove any existing merkleRoots and sigs
    while (s_merkleRoots.length > 0) {
      s_merkleRoots.pop();
    }
    while (s_signatures.length > 0) {
      s_signatures.pop();
    }
    s_v = 0;

    for (uint256 i = 0; i < numUpdates; i++) {
      s_merkleRoots.push(_generateRandomDestLaneUpdate());
    }

    for (uint256 i = 0; i < numSigs; i++) {
      (uint8 v, IRMNRemote.Signature memory sig) = _signDestLaneUpdate(s_merkleRoots, s_signerWallets[i]);
      s_signatures.push(sig);
      if (v == 28) {
        s_v += 1 << i;
      }
    }
  }

  /// @notice generates a random dest lane update
  function _generateRandomDestLaneUpdate() private returns (Internal.MerkleRoot memory) {
    uint64 minSeqNum = uint32(_randomNum());
    uint64 maxSeqNum = minSeqNum + 100;
    return Internal.MerkleRoot({
      sourceChainSelector: uint64(_randomNum()),
      onRampAddress: abi.encode(_randomAddress()),
      minSeqNr: minSeqNum,
      maxSeqNr: maxSeqNum,
      merkleRoot: _randomBytes32()
    });
  }

  /// @notice signs the provided payload with the provided wallet
  /// @return sigV v, either 27 of 28
  /// @return sig the signature
  function _signDestLaneUpdate(
    Internal.MerkleRoot[] memory merkleRoots,
    Vm.Wallet memory wallet
  ) private returns (uint8 sigV, IRMNRemote.Signature memory) {
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
    return (v, IRMNRemote.Signature({r: r, s: s}));
  }

  /// @notice bubble sort on a storage array of wallets
  function _sort(Vm.Wallet[] storage wallets) private {
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
}
