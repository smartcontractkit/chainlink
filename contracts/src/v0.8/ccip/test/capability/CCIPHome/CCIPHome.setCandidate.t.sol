// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_setCandidate is CCIPHomeTestSetup {
  function test_setCandidate() public {
    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);
    CCIPHome.VersionedConfig memory versionedConfig =
      CCIPHome.VersionedConfig({version: 1, config: config, configDigest: ZERO_DIGEST});

    versionedConfig.configDigest =
      _getConfigDigest(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, abi.encode(versionedConfig.config), versionedConfig.version);

    vm.expectEmit();
    emit CCIPHome.ConfigSet(versionedConfig.configDigest, versionedConfig.version, versionedConfig.config);

    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, versionedConfig.config, ZERO_DIGEST);

    (CCIPHome.VersionedConfig memory storedVersionedConfig, bool ok) =
      s_ccipHome.getConfig(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, versionedConfig.configDigest);
    assertTrue(ok);
    assertEq(storedVersionedConfig.version, versionedConfig.version);
    assertEq(storedVersionedConfig.configDigest, versionedConfig.configDigest);
    assertEq(keccak256(abi.encode(storedVersionedConfig.config)), keccak256(abi.encode(versionedConfig.config)));
  }

  function test_RevertWhen_setCandidate_ConfigDigestMismatch() public {
    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);

    bytes32 digest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ConfigDigestMismatch.selector, digest, ZERO_DIGEST));
    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    vm.expectEmit();
    emit CCIPHome.CandidateConfigRevoked(digest);

    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, digest);
  }

  function test_RevertWhen_setCandidate_CanOnlySelfCall() public {
    vm.stopPrank();

    vm.expectRevert(CCIPHome.CanOnlySelfCall.selector);
    s_ccipHome.setCandidate(
      DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST
    );
  }
}
