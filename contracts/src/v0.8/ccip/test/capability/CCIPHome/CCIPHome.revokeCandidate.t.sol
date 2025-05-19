// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_revokeCandidate is CCIPHomeTestSetup {
  // Sets two configs
  function setUp() public virtual override {
    super.setUp();
    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);
    bytes32 digest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);
    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, digest, ZERO_DIGEST);

    config.offrampAddress = abi.encode("new_offrampAddress");
    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);
  }

  function test_revokeCandidate() public {
    (bytes32 priorActiveDigest, bytes32 priorCandidateDigest) =
      s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);

    vm.expectEmit();
    emit CCIPHome.CandidateConfigRevoked(priorCandidateDigest);

    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, priorCandidateDigest);

    (CCIPHome.VersionedConfig memory storedVersionedConfig, bool ok) =
      s_ccipHome.getConfig(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, priorCandidateDigest);
    assertFalse(ok);
    // Ensure no old data is returned, even though it's still in storage
    assertEq(storedVersionedConfig.version, 0);
    assertEq(storedVersionedConfig.config.chainSelector, 0);
    assertEq(storedVersionedConfig.config.FRoleDON, 0);

    // Asser the active digest is unaffected but the candidate digest is set to zero
    (bytes32 activeDigest, bytes32 candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeDigest, priorActiveDigest);
    assertEq(candidateDigest, ZERO_DIGEST);
    assertTrue(candidateDigest != priorCandidateDigest);
  }

  function test_RevertWhen_revokeCandidate_ConfigDigestMismatch() public {
    (, bytes32 priorCandidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);

    bytes32 wrongDigest = keccak256("wrong_digest");
    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ConfigDigestMismatch.selector, priorCandidateDigest, wrongDigest));
    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, wrongDigest);
  }

  function test_RevertWhen_revokeCandidate_RevokingZeroDigestNotAllowed() public {
    vm.expectRevert(CCIPHome.RevokingZeroDigestNotAllowed.selector);
    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, ZERO_DIGEST);
  }

  function test_RevertWhen_revokeCandidate_CanOnlySelfCall() public {
    vm.startPrank(address(0));

    vm.expectRevert(CCIPHome.CanOnlySelfCall.selector);
    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, keccak256("configDigest"));
  }
}
