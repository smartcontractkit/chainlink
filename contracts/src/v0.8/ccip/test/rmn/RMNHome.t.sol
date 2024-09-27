// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {RMNHome} from "../../rmn/RMNHome.sol";
import {Test} from "forge-std/Test.sol";
import {Vm} from "forge-std/Vm.sol";

contract RMNHomeTest is Test {
  struct Config {
    RMNHome.StaticConfig staticConfig;
    RMNHome.DynamicConfig dynamicConfig;
  }

  bytes32 internal constant ZERO_DIGEST = bytes32(uint256(0));
  RMNHome public s_rmnHome = new RMNHome();

  function _getBaseConfig() internal pure returns (Config memory) {
    RMNHome.Node[] memory nodes = new RMNHome.Node[](3);
    nodes[0] = RMNHome.Node({peerId: keccak256("peerId_0"), offchainPublicKey: keccak256("offchainPublicKey_0")});
    nodes[1] = RMNHome.Node({peerId: keccak256("peerId_1"), offchainPublicKey: keccak256("offchainPublicKey_1")});
    nodes[2] = RMNHome.Node({peerId: keccak256("peerId_2"), offchainPublicKey: keccak256("offchainPublicKey_2")});

    RMNHome.SourceChain[] memory sourceChains = new RMNHome.SourceChain[](2);
    // Observer 0 for source chain 9000
    sourceChains[0] = RMNHome.SourceChain({chainSelector: 9000, minObservers: 1, observerNodesBitmap: 1 << 0});
    // Observers 1 and 2 for source chain 9001
    sourceChains[1] = RMNHome.SourceChain({chainSelector: 9001, minObservers: 2, observerNodesBitmap: 1 << 1 | 1 << 2});

    return Config({
      staticConfig: RMNHome.StaticConfig({nodes: nodes, offchainConfig: abi.encode("static_config")}),
      dynamicConfig: RMNHome.DynamicConfig({sourceChains: sourceChains, offchainConfig: abi.encode("dynamic_config")})
    });
  }

  uint256 private constant PREFIX_MASK = type(uint256).max << (256 - 16); // 0xFFFF00..00
  uint256 private constant PREFIX = 0x000b << (256 - 16); // 0x000b00..00

  function _getConfigDigest(bytes memory staticConfig, uint32 version) internal view returns (bytes32) {
    return bytes32(
      (PREFIX & PREFIX_MASK)
        | (
          uint256(
            keccak256(bytes.concat(abi.encode(bytes32("EVM"), block.chainid, address(s_rmnHome), version), staticConfig))
          ) & ~PREFIX_MASK
        )
    );
  }
}

contract RMNHome_getConfigDigests is RMNHomeTest {
  function test_getConfigDigests_success() public {
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

contract RMNHome_setCandidate is RMNHomeTest {
  function test_setCandidate_success() public {
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
      assertEq(storedSourceChain.minObservers, versionedConfig.dynamicConfig.sourceChains[i].minObservers);
      assertEq(storedSourceChain.observerNodesBitmap, versionedConfig.dynamicConfig.sourceChains[i].observerNodesBitmap);
    }
    assertEq(storedDynamicConfig.offchainConfig, versionedConfig.dynamicConfig.offchainConfig);
    assertEq(storedStaticConfig.offchainConfig, versionedConfig.staticConfig.offchainConfig);
  }

  function test_setCandidate_ConfigDigestMismatch_reverts() public {
    Config memory config = _getBaseConfig();

    bytes32 digest = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    vm.expectRevert(abi.encodeWithSelector(RMNHome.ConfigDigestMismatch.selector, digest, ZERO_DIGEST));
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);

    vm.expectEmit();
    emit RMNHome.CandidateConfigRevoked(digest);

    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, digest);
  }

  function test_setCandidate_OnlyOwner_reverts() public {
    Config memory config = _getBaseConfig();

    vm.startPrank(address(0));

    vm.expectRevert("Only callable by owner");
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }
}

contract RMNHome_revokeCandidate is RMNHomeTest {
  // Sets two configs
  function setUp() public {
    Config memory config = _getBaseConfig();
    bytes32 digest = s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
    s_rmnHome.promoteCandidateAndRevokeActive(digest, ZERO_DIGEST);

    config.dynamicConfig.sourceChains[0].minObservers--;
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_revokeCandidate_success() public {
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

  function test_revokeCandidate_ConfigDigestMismatch_reverts() public {
    (, bytes32 priorCandidateDigest) = s_rmnHome.getConfigDigests();

    bytes32 wrongDigest = keccak256("wrong_digest");
    vm.expectRevert(abi.encodeWithSelector(RMNHome.ConfigDigestMismatch.selector, priorCandidateDigest, wrongDigest));
    s_rmnHome.revokeCandidate(wrongDigest);
  }

  function test_revokeCandidate_RevokingZeroDigestNotAllowed_reverts() public {
    vm.expectRevert(RMNHome.RevokingZeroDigestNotAllowed.selector);
    s_rmnHome.revokeCandidate(ZERO_DIGEST);
  }

  function test_revokeCandidate_OnlyOwner_reverts() public {
    vm.startPrank(address(0));

    vm.expectRevert("Only callable by owner");
    s_rmnHome.revokeCandidate(keccak256("configDigest"));
  }
}

contract RMNHome_promoteCandidateAndRevokeActive is RMNHomeTest {
  function test_promoteCandidateAndRevokeActive_success() public {
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

  function test_promoteCandidateAndRevokeActive_NoOpStateTransitionNotAllowed_reverts() public {
    vm.expectRevert(RMNHome.NoOpStateTransitionNotAllowed.selector);
    s_rmnHome.promoteCandidateAndRevokeActive(ZERO_DIGEST, ZERO_DIGEST);
  }

  function test_promoteCandidateAndRevokeActive_ConfigDigestMismatch_reverts() public {
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

  function test_promoteCandidateAndRevokeActive_OnlyOwner_reverts() public {
    vm.startPrank(address(0));

    vm.expectRevert("Only callable by owner");
    s_rmnHome.promoteCandidateAndRevokeActive(keccak256("toPromote"), keccak256("ToRevoke"));
  }
}

contract RMNHome__validateStaticAndDynamicConfig is RMNHomeTest {
  function test_validateStaticAndDynamicConfig_OutOfBoundsNodesLength_reverts() public {
    Config memory config = _getBaseConfig();
    config.staticConfig.nodes = new RMNHome.Node[](257);

    vm.expectRevert(RMNHome.OutOfBoundsNodesLength.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_validateStaticAndDynamicConfig_DuplicatePeerId_reverts() public {
    Config memory config = _getBaseConfig();
    config.staticConfig.nodes[1].peerId = config.staticConfig.nodes[0].peerId;

    vm.expectRevert(RMNHome.DuplicatePeerId.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_validateStaticAndDynamicConfig_DuplicateOffchainPublicKey_reverts() public {
    Config memory config = _getBaseConfig();
    config.staticConfig.nodes[1].offchainPublicKey = config.staticConfig.nodes[0].offchainPublicKey;

    vm.expectRevert(RMNHome.DuplicateOffchainPublicKey.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_validateStaticAndDynamicConfig_DuplicateSourceChain_reverts() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[1].chainSelector = config.dynamicConfig.sourceChains[0].chainSelector;

    vm.expectRevert(RMNHome.DuplicateSourceChain.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_validateStaticAndDynamicConfig_OutOfBoundsObserverNodeIndex_reverts() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].observerNodesBitmap = 1 << config.staticConfig.nodes.length;

    vm.expectRevert(RMNHome.OutOfBoundsObserverNodeIndex.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_validateStaticAndDynamicConfig_MinObserversTooHigh_reverts() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].minObservers++;

    vm.expectRevert(RMNHome.MinObserversTooHigh.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }
}

contract RMNHome_setDynamicConfig is RMNHomeTest {
  function setUp() public {
    Config memory config = _getBaseConfig();
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_setDynamicConfig_success() public {
    (bytes32 priorActiveDigest,) = s_rmnHome.getConfigDigests();

    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].minObservers--;

    (, bytes32 candidateConfigDigest) = s_rmnHome.getConfigDigests();

    vm.expectEmit();
    emit RMNHome.DynamicConfigSet(candidateConfigDigest, config.dynamicConfig);

    s_rmnHome.setDynamicConfig(config.dynamicConfig, candidateConfigDigest);

    (RMNHome.VersionedConfig memory storedVersionedConfig, bool ok) = s_rmnHome.getConfig(candidateConfigDigest);
    assertTrue(ok);
    assertEq(
      storedVersionedConfig.dynamicConfig.sourceChains[0].minObservers,
      config.dynamicConfig.sourceChains[0].minObservers
    );

    // Asser the digests don't change when updating the dynamic config
    (bytes32 activeDigest, bytes32 candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, priorActiveDigest);
    assertEq(candidateDigest, candidateConfigDigest);
  }

  // Asserts the validation function is being called
  function test_setDynamicConfig_MinObserversTooHigh_reverts() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].minObservers++;

    vm.expectRevert(abi.encodeWithSelector(RMNHome.DigestNotFound.selector, ZERO_DIGEST));
    s_rmnHome.setDynamicConfig(config.dynamicConfig, ZERO_DIGEST);
  }

  function test_setDynamicConfig_DigestNotFound_reverts() public {
    // Zero always reverts
    vm.expectRevert(abi.encodeWithSelector(RMNHome.DigestNotFound.selector, ZERO_DIGEST));
    s_rmnHome.setDynamicConfig(_getBaseConfig().dynamicConfig, ZERO_DIGEST);

    // Non-existent digest reverts
    bytes32 nonExistentDigest = keccak256("nonExistentDigest");
    vm.expectRevert(abi.encodeWithSelector(RMNHome.DigestNotFound.selector, nonExistentDigest));
    s_rmnHome.setDynamicConfig(_getBaseConfig().dynamicConfig, nonExistentDigest);
  }

  function test_setDynamicConfig_OnlyOwner_reverts() public {
    Config memory config = _getBaseConfig();

    vm.startPrank(address(0));

    vm.expectRevert("Only callable by owner");
    s_rmnHome.setDynamicConfig(config.dynamicConfig, keccak256("configDigest"));
  }
}
