// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RMNHome} from "../../../rmn/RMNHome.sol";

import {RMNHomeTestSetup} from "./RMNHomeTestSetup.t.sol";

contract RMNHome_revokeCandidate is RMNHomeTestSetup {
  // Sets two configs
  function setUp() public {
    Config memory config = _getBaseConfig();
    bytes32 digest = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
    s_rmnHome.promoteCandidateAndRevokeActive(digest, ZERO_DIGEST);

    config.dynamicConfig.sourceChains[1].fObserve--;
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_revokeCandidate() public {
    (bytes32 priorActiveDigest, bytes32 priorCandidateDigest) = s_rmnHome.getConfigDigests();

    vm.expectEmit();
    emit RMNHome.CandidateConfigRevoked(priorCandidateDigest);

    s_rmnHome.revokeCandidate(priorCandidateDigest);

    (RMNHome.VersionedConfig memory storedVersionedConfig, bool ok) = s_rmnHome.getConfig(priorCandidateDigest);
    assertFalse(ok);
    // Ensure no old data is returned, even though it's still in storage
    assertEq(storedVersionedConfig.version, 0);
    assertEq(storedVersionedConfig.staticConfig.nodes.length, 0);
    assertEq(storedVersionedConfig.dynamicConfig.sourceChains.length, 0);

    // Asser the active digest is unaffected but the candidate digest is set to zero
    (bytes32 activeDigest, bytes32 candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, priorActiveDigest);
    assertEq(candidateDigest, ZERO_DIGEST);
    assertTrue(candidateDigest != priorCandidateDigest);
  }

  function test_RevertWhen_revokeCandidate_ConfigDigestMismatch() public {
    (, bytes32 priorCandidateDigest) = s_rmnHome.getConfigDigests();

    bytes32 wrongDigest = keccak256("wrong_digest");
    vm.expectRevert(abi.encodeWithSelector(RMNHome.ConfigDigestMismatch.selector, priorCandidateDigest, wrongDigest));
    s_rmnHome.revokeCandidate(wrongDigest);
  }

  function test_RevertWhen_revokeCandidate_RevokingZeroDigestNotAllowed() public {
    vm.expectRevert(RMNHome.RevokingZeroDigestNotAllowed.selector);
    s_rmnHome.revokeCandidate(ZERO_DIGEST);
  }

  function test_RevertWhen_revokeCandidate_OnlyOwner() public {
    vm.startPrank(address(0));

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_rmnHome.revokeCandidate(keccak256("configDigest"));
  }
}
