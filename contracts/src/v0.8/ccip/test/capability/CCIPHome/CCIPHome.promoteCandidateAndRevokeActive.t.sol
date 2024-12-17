// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_promoteCandidateAndRevokeActive is CCIPHomeTestSetup {
  function test_promoteCandidateAndRevokeActive_multiplePlugins() public {
    promoteCandidateAndRevokeActive(Internal.OCRPluginType.Commit);
    promoteCandidateAndRevokeActive(Internal.OCRPluginType.Execution);

    // check that the two plugins have only active configs and no candidates.
    (bytes32 activeDigest, bytes32 candidateDigest) =
      s_ccipHome.getConfigDigests(DEFAULT_DON_ID, Internal.OCRPluginType.Commit);
    assertTrue(activeDigest != ZERO_DIGEST);
    assertEq(candidateDigest, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, Internal.OCRPluginType.Execution);
    assertTrue(activeDigest != ZERO_DIGEST);
    assertEq(candidateDigest, ZERO_DIGEST);
  }

  function promoteCandidateAndRevokeActive(
    Internal.OCRPluginType pluginType
  ) public {
    CCIPHome.OCR3Config memory config = _getBaseConfig(pluginType);
    bytes32 firstConfigToPromote = s_ccipHome.setCandidate(DEFAULT_DON_ID, pluginType, config, ZERO_DIGEST);

    vm.expectEmit();
    emit CCIPHome.ConfigPromoted(firstConfigToPromote);

    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, pluginType, firstConfigToPromote, ZERO_DIGEST);

    // Assert the active digest is updated and the candidate digest is set to zero
    (bytes32 activeDigest, bytes32 candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, pluginType);
    assertEq(activeDigest, firstConfigToPromote);
    assertEq(candidateDigest, ZERO_DIGEST);

    // Set a new candidate to promote over a non-zero active config.
    config.offchainConfig = abi.encode("new_offchainConfig_config");
    bytes32 secondConfigToPromote = s_ccipHome.setCandidate(DEFAULT_DON_ID, pluginType, config, ZERO_DIGEST);

    vm.expectEmit();
    emit CCIPHome.ActiveConfigRevoked(firstConfigToPromote);

    vm.expectEmit();
    emit CCIPHome.ConfigPromoted(secondConfigToPromote);

    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, pluginType, secondConfigToPromote, firstConfigToPromote);

    (CCIPHome.VersionedConfig memory activeConfig, CCIPHome.VersionedConfig memory candidateConfig) =
      s_ccipHome.getAllConfigs(DEFAULT_DON_ID, pluginType);
    assertEq(activeConfig.configDigest, secondConfigToPromote);
    assertEq(candidateConfig.configDigest, ZERO_DIGEST);
    assertEq(keccak256(abi.encode(activeConfig.config)), keccak256(abi.encode(config)));
  }

  function test_RevertWhen_promoteCandidateAndRevokeActive_NoOpStateTransitionNotAllowed() public {
    vm.expectRevert(CCIPHome.NoOpStateTransitionNotAllowed.selector);
    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, ZERO_DIGEST, ZERO_DIGEST);
  }

  function test_RevertWhen_promoteCandidateAndRevokeActive_ConfigDigestMismatch() public {
    (bytes32 priorActiveDigest, bytes32 priorCandidateDigest) =
      s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    bytes32 wrongActiveDigest = keccak256("wrongActiveDigest");
    bytes32 wrongCandidateDigest = keccak256("wrongCandidateDigest");

    vm.expectRevert(
      abi.encodeWithSelector(CCIPHome.ConfigDigestMismatch.selector, priorActiveDigest, wrongCandidateDigest)
    );
    s_ccipHome.promoteCandidateAndRevokeActive(
      DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, wrongCandidateDigest, wrongActiveDigest
    );

    vm.expectRevert(
      abi.encodeWithSelector(CCIPHome.ConfigDigestMismatch.selector, priorActiveDigest, wrongActiveDigest)
    );

    s_ccipHome.promoteCandidateAndRevokeActive(
      DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, priorCandidateDigest, wrongActiveDigest
    );
  }

  function test_RevertWhen_promoteCandidateAndRevokeActive_CanOnlySelfCall() public {
    vm.stopPrank();

    vm.expectRevert(CCIPHome.CanOnlySelfCall.selector);
    s_ccipHome.promoteCandidateAndRevokeActive(
      DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, keccak256("toPromote"), keccak256("ToRevoke")
    );
  }
}
