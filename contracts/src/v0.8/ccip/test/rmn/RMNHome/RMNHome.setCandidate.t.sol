// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RMNHome} from "../../../rmn/RMNHome.sol";

import {RMNHomeTestSetup} from "./RMNHomeTestSetup.t.sol";

contract RMNHome_setCandidate is RMNHomeTestSetup {
  function test_setCandidate() public {
    Config memory config = _getBaseConfig();
    RMNHome.VersionedConfig memory versionedConfig = RMNHome.VersionedConfig({
      version: 1,
      staticConfig: config.staticConfig,
      dynamicConfig: config.dynamicConfig,
      configDigest: ZERO_DIGEST
    });

    versionedConfig.configDigest = _getConfigDigest(abi.encode(versionedConfig.staticConfig), versionedConfig.version);

    vm.expectEmit();
    emit RMNHome.ConfigSet(
      versionedConfig.configDigest, versionedConfig.version, versionedConfig.staticConfig, versionedConfig.dynamicConfig
    );

    s_rmnHome.setCandidate(versionedConfig.staticConfig, versionedConfig.dynamicConfig, ZERO_DIGEST);

    (RMNHome.VersionedConfig memory storedVersionedConfig, bool ok) = s_rmnHome.getConfig(versionedConfig.configDigest);
    assertTrue(ok);
    assertEq(storedVersionedConfig.version, versionedConfig.version);
    RMNHome.StaticConfig memory storedStaticConfig = storedVersionedConfig.staticConfig;
    RMNHome.DynamicConfig memory storedDynamicConfig = storedVersionedConfig.dynamicConfig;

    assertEq(storedStaticConfig.nodes.length, versionedConfig.staticConfig.nodes.length);
    for (uint256 i = 0; i < storedStaticConfig.nodes.length; i++) {
      RMNHome.Node memory storedNode = storedStaticConfig.nodes[i];
      assertEq(storedNode.peerId, versionedConfig.staticConfig.nodes[i].peerId);
      assertEq(storedNode.offchainPublicKey, versionedConfig.staticConfig.nodes[i].offchainPublicKey);
    }

    assertEq(storedDynamicConfig.sourceChains.length, versionedConfig.dynamicConfig.sourceChains.length);
    for (uint256 i = 0; i < storedDynamicConfig.sourceChains.length; i++) {
      RMNHome.SourceChain memory storedSourceChain = storedDynamicConfig.sourceChains[i];
      assertEq(storedSourceChain.chainSelector, versionedConfig.dynamicConfig.sourceChains[i].chainSelector);
      assertEq(storedSourceChain.fObserve, versionedConfig.dynamicConfig.sourceChains[i].fObserve);
      assertEq(storedSourceChain.observerNodesBitmap, versionedConfig.dynamicConfig.sourceChains[i].observerNodesBitmap);
    }
    assertEq(storedDynamicConfig.offchainConfig, versionedConfig.dynamicConfig.offchainConfig);
    assertEq(storedStaticConfig.offchainConfig, versionedConfig.staticConfig.offchainConfig);
  }

  function test_RevertWhen_setCandidate_ConfigDigestMismatch() public {
    Config memory config = _getBaseConfig();

    bytes32 digest = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    vm.expectRevert(abi.encodeWithSelector(RMNHome.ConfigDigestMismatch.selector, digest, ZERO_DIGEST));
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    vm.expectEmit();
    emit RMNHome.CandidateConfigRevoked(digest);

    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, digest);
  }

  function test_RevertWhen_setCandidate_OnlyOwner() public {
    Config memory config = _getBaseConfig();

    vm.startPrank(address(0));

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }
}
