// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RMNHome} from "../../../rmn/RMNHome.sol";

import {RMNHomeTestSetup} from "./RMNHomeTestSetup.t.sol";

contract RMNHome_promoteCandidateAndRevokeActive is RMNHomeTestSetup {
  function test_promoteCandidateAndRevokeActive() public {
    Config memory config = _getBaseConfig();
    bytes32 firstConfigToPromote = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    vm.expectEmit();
    emit RMNHome.ConfigPromoted(firstConfigToPromote);

    s_rmnHome.promoteCandidateAndRevokeActive(firstConfigToPromote, ZERO_DIGEST);

    // Assert the active digest is updated and the candidate digest is set to zero
    (bytes32 activeDigest, bytes32 candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, firstConfigToPromote);
    assertEq(candidateDigest, ZERO_DIGEST);

    // Set a new candidate to promote over a non-zero active config.
    config.staticConfig.offchainConfig = abi.encode("new_static_config");
    config.dynamicConfig.offchainConfig = abi.encode("new_dynamic_config");
    bytes32 secondConfigToPromote = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    vm.expectEmit();
    emit RMNHome.ActiveConfigRevoked(firstConfigToPromote);

    vm.expectEmit();
    emit RMNHome.ConfigPromoted(secondConfigToPromote);

    s_rmnHome.promoteCandidateAndRevokeActive(secondConfigToPromote, firstConfigToPromote);

    (RMNHome.VersionedConfig memory activeConfig, RMNHome.VersionedConfig memory candidateConfig) =
      s_rmnHome.getAllConfigs();
    assertEq(activeConfig.configDigest, secondConfigToPromote);
    assertEq(activeConfig.staticConfig.offchainConfig, config.staticConfig.offchainConfig);
    assertEq(activeConfig.dynamicConfig.offchainConfig, config.dynamicConfig.offchainConfig);

    assertEq(candidateConfig.configDigest, ZERO_DIGEST);
  }

  function test_RevertWhen_promoteCandidateAndRevokeActive_NoOpStateTransitionNotAllowed() public {
    vm.expectRevert(RMNHome.NoOpStateTransitionNotAllowed.selector);
    s_rmnHome.promoteCandidateAndRevokeActive(ZERO_DIGEST, ZERO_DIGEST);
  }

  function test_RevertWhen_promoteCandidateAndRevokeActive_ConfigDigestMismatch() public {
    (bytes32 priorActiveDigest, bytes32 priorCandidateDigest) = s_rmnHome.getConfigDigests();
    bytes32 wrongActiveDigest = keccak256("wrongActiveDigest");
    bytes32 wrongCandidateDigest = keccak256("wrongCandidateDigest");

    vm.expectRevert(
      abi.encodeWithSelector(RMNHome.ConfigDigestMismatch.selector, priorActiveDigest, wrongCandidateDigest)
    );
    s_rmnHome.promoteCandidateAndRevokeActive(wrongCandidateDigest, wrongActiveDigest);

    vm.expectRevert(abi.encodeWithSelector(RMNHome.ConfigDigestMismatch.selector, priorActiveDigest, wrongActiveDigest));

    s_rmnHome.promoteCandidateAndRevokeActive(priorCandidateDigest, wrongActiveDigest);
  }

  function test_RevertWhen_promoteCandidateAndRevokeActive_OnlyOwner() public {
    vm.startPrank(address(0));

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_rmnHome.promoteCandidateAndRevokeActive(keccak256("toPromote"), keccak256("ToRevoke"));
  }
}
