// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNHomeTestSetup} from "./RMNHomeTestSetup.t.sol";

contract RMNHome_getConfigDigests is RMNHomeTestSetup {
  function test_getConfigDigests() public {
    (bytes32 activeDigest, bytes32 candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, ZERO_DIGEST);
    assertEq(candidateDigest, ZERO_DIGEST);

    Config memory config = _getBaseConfig();
    bytes32 firstDigest = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, ZERO_DIGEST);
    assertEq(candidateDigest, firstDigest);

    s_rmnHome.promoteCandidateAndRevokeActive(firstDigest, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, firstDigest);
    assertEq(candidateDigest, ZERO_DIGEST);

    bytes32 secondDigest = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, firstDigest);
    assertEq(candidateDigest, secondDigest);

    assertEq(activeDigest, s_rmnHome.getActiveDigest());
    assertEq(candidateDigest, s_rmnHome.getCandidateDigest());
  }
}
