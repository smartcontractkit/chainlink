// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_beforeCapabilityConfigSet is CCIPHomeTestSetup {
  function setUp() public virtual override {
    super.setUp();
    vm.stopPrank();
    vm.startPrank(address(CAPABILITIES_REGISTRY));
  }

  function test_beforeCapabilityConfigSet() public {
    // first set a config
    bytes memory callData = abi.encodeCall(
      CCIPHome.setCandidate,
      (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST)
    );

    vm.expectCall(address(s_ccipHome), callData);

    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);

    // Then revoke the config
    bytes32 candidateDigest = s_ccipHome.getCandidateDigest(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertNotEq(candidateDigest, ZERO_DIGEST);

    callData = abi.encodeCall(CCIPHome.revokeCandidate, (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, candidateDigest));

    vm.expectCall(address(s_ccipHome), callData);

    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);

    // Then set a new config
    callData = abi.encodeCall(
      CCIPHome.setCandidate,
      (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST)
    );

    vm.expectCall(address(s_ccipHome), callData);

    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);

    // Then promote the new config

    bytes32 newCandidateDigest = s_ccipHome.getCandidateDigest(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertNotEq(newCandidateDigest, ZERO_DIGEST);

    callData = abi.encodeCall(
      CCIPHome.promoteCandidateAndRevokeActive, (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, newCandidateDigest, ZERO_DIGEST)
    );

    vm.expectCall(address(s_ccipHome), callData);

    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);

    bytes32 activeDigest = s_ccipHome.getActiveDigest(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);
    assertEq(activeDigest, newCandidateDigest);
  }

  function test_RevertWhen_beforeCapabilityConfigSet_OnlyCapabilitiesRegistryCanCall() public {
    bytes memory callData = abi.encodeCall(
      CCIPHome.setCandidate,
      (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST)
    );

    vm.stopPrank();

    vm.expectRevert(CCIPHome.OnlyCapabilitiesRegistryCanCall.selector);

    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);
  }

  function test_RevertWhen_beforeCapabilityConfigSet_InvalidSelector() public {
    bytes memory callData = abi.encodeCall(CCIPHome.getConfigDigests, (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE));

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.InvalidSelector.selector, CCIPHome.getConfigDigests.selector));
    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);
  }

  function test_RevertWhen_beforeCapabilityConfigSet_DONIdMismatch() public {
    uint32 wrongDonId = DEFAULT_DON_ID + 1;

    bytes memory callData = abi.encodeCall(
      CCIPHome.setCandidate,
      (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST)
    );

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.DONIdMismatch.selector, DEFAULT_DON_ID, wrongDonId));
    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, wrongDonId);
  }

  function test_RevertWhen_beforeCapabilityConfigSet_InnerCallReverts() public {
    bytes memory callData = abi.encodeCall(CCIPHome.revokeCandidate, (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, ZERO_DIGEST));

    vm.expectRevert(CCIPHome.RevokingZeroDigestNotAllowed.selector);
    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);
  }
}
