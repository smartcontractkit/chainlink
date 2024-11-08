// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../../keystone/interfaces/ICapabilityConfiguration.sol";
import {INodeInfoProvider} from "../../../keystone/interfaces/INodeInfoProvider.sol";

import {CCIPHome} from "../../capability/CCIPHome.sol";
import {Internal} from "../../libraries/Internal.sol";
import {CCIPHomeHelper} from "../helpers/CCIPHomeHelper.sol";
import {Test} from "forge-std/Test.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";
import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_getAllConfigs is CCIPHomeTestSetup {
  function test_getAllConfigs_success() public {
    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);
    bytes32 firstDigest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    (CCIPHome.VersionedConfig memory activeConfig, CCIPHome.VersionedConfig memory candidateConfig) =
      s_ccipHome.getAllConfigs(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeConfig.configDigest, ZERO_DIGEST);
    assertEq(candidateConfig.configDigest, firstDigest);

    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, firstDigest, ZERO_DIGEST);

    (activeConfig, candidateConfig) = s_ccipHome.getAllConfigs(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeConfig.configDigest, firstDigest);
    assertEq(candidateConfig.configDigest, ZERO_DIGEST);

    bytes32 secondDigest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    (activeConfig, candidateConfig) = s_ccipHome.getAllConfigs(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeConfig.configDigest, firstDigest);
    assertEq(candidateConfig.configDigest, secondDigest);

    (activeConfig, candidateConfig) = s_ccipHome.getAllConfigs(DEFAULT_DON_ID + 1, DEFAULT_PLUGIN_TYPE);
    assertEq(activeConfig.configDigest, ZERO_DIGEST);
    assertEq(candidateConfig.configDigest, ZERO_DIGEST);

    (activeConfig, candidateConfig) = s_ccipHome.getAllConfigs(DEFAULT_DON_ID, Internal.OCRPluginType.Execution);
    assertEq(activeConfig.configDigest, ZERO_DIGEST);
    assertEq(candidateConfig.configDigest, ZERO_DIGEST);
  }
}
