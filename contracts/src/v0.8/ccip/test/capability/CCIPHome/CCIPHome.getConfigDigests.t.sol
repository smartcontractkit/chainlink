// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_getConfigDigests is CCIPHomeTestSetup {
  function test_getConfigDigests() public {
    (bytes32 activeDigest, bytes32 candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeDigest, ZERO_DIGEST);
    assertEq(candidateDigest, ZERO_DIGEST);

    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);
    bytes32 firstDigest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeDigest, ZERO_DIGEST);
    assertEq(candidateDigest, firstDigest);

    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, firstDigest, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeDigest, firstDigest);
    assertEq(candidateDigest, ZERO_DIGEST);

    bytes32 secondDigest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    (activeDigest, candidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeDigest, firstDigest);
    assertEq(candidateDigest, secondDigest);

    assertEq(activeDigest, s_ccipHome.getActiveDigest(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE));
    assertEq(candidateDigest, s_ccipHome.getCandidateDigest(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE));
  }
}
