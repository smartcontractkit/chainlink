// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../../keystone/interfaces/ICapabilityConfiguration.sol";
import {ICapabilitiesRegistry} from "../../interfaces/ICapabilitiesRegistry.sol";

import {CCIPHome} from "../../capability/CCIPHome.sol";
import {Internal} from "../../libraries/Internal.sol";
import {CCIPHomeHelper} from "../helpers/CCIPHomeHelper.sol";
import {Test} from "forge-std/Test.sol";
import {Vm} from "forge-std/Vm.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";

contract CCIPHomeTest is Test {
  //  address internal constant OWNER = address(0x0000000123123123123);
  bytes32 internal constant ZERO_DIGEST = bytes32(uint256(0));
  address internal constant CAPABILITIES_REGISTRY = address(0x0000000123123123123);
  Internal.OCRPluginType internal constant DEFAULT_PLUGIN_TYPE = Internal.OCRPluginType.Commit;
  uint32 internal constant DEFAULT_DON_ID = 78978987;

  CCIPHomeHelper public s_ccipHome;

  uint256 private constant PREFIX_MASK = type(uint256).max << (256 - 16); // 0xFFFF00..00
  uint256 private constant PREFIX = 0x000a << (256 - 16); // 0x000b00..00

  uint64 private constant DEFAULT_CHAIN_SELECTOR = 9381579735;

  function setUp() public virtual {
    s_ccipHome = new CCIPHomeHelper(CAPABILITIES_REGISTRY);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), _getBaseChainConfigs());

    ICapabilitiesRegistry.NodeInfo memory nodeInfo = ICapabilitiesRegistry.NodeInfo({
      p2pId: keccak256("p2pId"),
      signer: keccak256("signer"),
      nodeOperatorId: 1,
      configCount: 1,
      workflowDONId: 1,
      encryptionPublicKey: keccak256("encryptionPublicKey"),
      hashedCapabilityIds: new bytes32[](0),
      capabilitiesDONIds: new uint256[](0)
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY, abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector), abi.encode(nodeInfo)
    );

    vm.startPrank(address(s_ccipHome));
  }

  function _getBaseChainConfigs() internal pure returns (CCIPHome.ChainConfigArgs[] memory) {
    CCIPHome.ChainConfigArgs[] memory configs = new CCIPHome.ChainConfigArgs[](1);
    CCIPHome.ChainConfig memory chainConfig =
      CCIPHome.ChainConfig({readers: new bytes32[](0), fChain: 1, config: abi.encode("chainConfig")});
    configs[0] = CCIPHome.ChainConfigArgs({chainSelector: DEFAULT_CHAIN_SELECTOR, chainConfig: chainConfig});

    return configs;
  }

  function _getConfigDigest(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    bytes memory config,
    uint32 version
  ) internal view returns (bytes32) {
    return bytes32(
      (PREFIX & PREFIX_MASK)
        | (
          uint256(
            keccak256(
              bytes.concat(
                abi.encode(bytes32("EVM"), block.chainid, address(s_ccipHome), donId, pluginType, version), config
              )
            )
          ) & ~PREFIX_MASK
        )
    );
  }

  function _getBaseConfig(
    Internal.OCRPluginType pluginType
  ) internal pure returns (CCIPHome.OCR3Config memory) {
    CCIPHome.OCR3Node[] memory nodes = new CCIPHome.OCR3Node[](4);
    for (uint256 i = 0; i < nodes.length; i++) {
      nodes[i] = CCIPHome.OCR3Node({
        p2pId: keccak256(abi.encode("p2pId", i)),
        signerKey: abi.encode("signerKey"),
        transmitterKey: abi.encode("transmitterKey")
      });
    }

    return CCIPHome.OCR3Config({
      pluginType: pluginType,
      chainSelector: DEFAULT_CHAIN_SELECTOR,
      FRoleDON: 1,
      offchainConfigVersion: 98765,
      offrampAddress: abi.encode("offrampAddress"),
      rmnHomeAddress: abi.encode("rmnHomeAddress"),
      nodes: nodes,
      offchainConfig: abi.encode("offchainConfig")
    });
  }
}

contract CCIPHome_constructor is CCIPHomeTest {
  function test_constructor_success() public {
    CCIPHome ccipHome = new CCIPHome(CAPABILITIES_REGISTRY);

    assertEq(address(ccipHome.getCapabilityRegistry()), CAPABILITIES_REGISTRY);
  }

  function test_supportsInterface_success() public view {
    assertTrue(s_ccipHome.supportsInterface(type(IERC165).interfaceId));
    assertTrue(s_ccipHome.supportsInterface(type(ICapabilityConfiguration).interfaceId));
  }

  function test_getCapabilityConfiguration_success() public view {
    bytes memory config = s_ccipHome.getCapabilityConfiguration(DEFAULT_DON_ID);
    assertEq(config.length, 0);
  }

  function test_constructor_CapabilitiesRegistryAddressZero_reverts() public {
    vm.expectRevert(CCIPHome.ZeroAddressNotAllowed.selector);
    new CCIPHome(address(0));
  }
}

contract CCIPHome_beforeCapabilityConfigSet is CCIPHomeTest {
  function setUp() public virtual override {
    super.setUp();
    vm.stopPrank();
    vm.startPrank(address(CAPABILITIES_REGISTRY));
  }

  function test_beforeCapabilityConfigSet_success() public {
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

  function test_beforeCapabilityConfigSet_OnlyCapabilitiesRegistryCanCall_reverts() public {
    bytes memory callData = abi.encodeCall(
      CCIPHome.setCandidate,
      (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST)
    );

    vm.stopPrank();

    vm.expectRevert(CCIPHome.OnlyCapabilitiesRegistryCanCall.selector);

    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);
  }

  function test_beforeCapabilityConfigSet_InvalidSelector_reverts() public {
    bytes memory callData = abi.encodeCall(CCIPHome.getConfigDigests, (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE));

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.InvalidSelector.selector, CCIPHome.getConfigDigests.selector));
    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);
  }

  function test_beforeCapabilityConfigSet_DONIdMismatch_reverts() public {
    uint32 wrongDonId = DEFAULT_DON_ID + 1;

    bytes memory callData = abi.encodeCall(
      CCIPHome.setCandidate,
      (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST)
    );

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.DONIdMismatch.selector, DEFAULT_DON_ID, wrongDonId));
    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, wrongDonId);
  }

  function test_beforeCapabilityConfigSet_InnerCallReverts_reverts() public {
    bytes memory callData = abi.encodeCall(CCIPHome.revokeCandidate, (DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, ZERO_DIGEST));

    vm.expectRevert(CCIPHome.RevokingZeroDigestNotAllowed.selector);
    s_ccipHome.beforeCapabilityConfigSet(new bytes32[](0), callData, 0, DEFAULT_DON_ID);
  }
}

contract CCIPHome_getConfigDigests is CCIPHomeTest {
  function test_getConfigDigests_success() public {
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

contract CCIPHome_getAllConfigs is CCIPHomeTest {
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

contract CCIPHome_setCandidate is CCIPHomeTest {
  function test_setCandidate_success() public {
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

  function test_setCandidate_ConfigDigestMismatch_reverts() public {
    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);

    bytes32 digest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ConfigDigestMismatch.selector, digest, ZERO_DIGEST));
    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);

    vm.expectEmit();
    emit CCIPHome.CandidateConfigRevoked(digest);

    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, digest);
  }

  function test_setCandidate_CanOnlySelfCall_reverts() public {
    vm.stopPrank();

    vm.expectRevert(CCIPHome.CanOnlySelfCall.selector);
    s_ccipHome.setCandidate(
      DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, _getBaseConfig(Internal.OCRPluginType.Commit), ZERO_DIGEST
    );
  }
}

contract CCIPHome_revokeCandidate is CCIPHomeTest {
  // Sets two configs
  function setUp() public virtual override {
    super.setUp();
    CCIPHome.OCR3Config memory config = _getBaseConfig(Internal.OCRPluginType.Commit);
    bytes32 digest = s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);
    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, digest, ZERO_DIGEST);

    config.offrampAddress = abi.encode("new_offrampAddress");
    s_ccipHome.setCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, config, ZERO_DIGEST);
  }

  function test_revokeCandidate_success() public {
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

  function test_revokeCandidate_ConfigDigestMismatch_reverts() public {
    (, bytes32 priorCandidateDigest) = s_ccipHome.getConfigDigests(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE);

    bytes32 wrongDigest = keccak256("wrong_digest");
    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ConfigDigestMismatch.selector, priorCandidateDigest, wrongDigest));
    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, wrongDigest);
  }

  function test_revokeCandidate_RevokingZeroDigestNotAllowed_reverts() public {
    vm.expectRevert(CCIPHome.RevokingZeroDigestNotAllowed.selector);
    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, ZERO_DIGEST);
  }

  function test_revokeCandidate_CanOnlySelfCall_reverts() public {
    vm.startPrank(address(0));

    vm.expectRevert(CCIPHome.CanOnlySelfCall.selector);
    s_ccipHome.revokeCandidate(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, keccak256("configDigest"));
  }
}

contract CCIPHome_promoteCandidateAndRevokeActive is CCIPHomeTest {
  function test_promoteCandidateAndRevokeActive_multiplePlugins_success() public {
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

  function test_promoteCandidateAndRevokeActive_NoOpStateTransitionNotAllowed_reverts() public {
    vm.expectRevert(CCIPHome.NoOpStateTransitionNotAllowed.selector);
    s_ccipHome.promoteCandidateAndRevokeActive(DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, ZERO_DIGEST, ZERO_DIGEST);
  }

  function test_promoteCandidateAndRevokeActive_ConfigDigestMismatch_reverts() public {
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

  function test_promoteCandidateAndRevokeActive_CanOnlySelfCall_reverts() public {
    vm.stopPrank();

    vm.expectRevert(CCIPHome.CanOnlySelfCall.selector);
    s_ccipHome.promoteCandidateAndRevokeActive(
      DEFAULT_DON_ID, DEFAULT_PLUGIN_TYPE, keccak256("toPromote"), keccak256("ToRevoke")
    );
  }
}

contract CCIPHome__validateConfig is CCIPHomeTest {
  function setUp() public virtual override {
    s_ccipHome = new CCIPHomeHelper(CAPABILITIES_REGISTRY);
  }

  function _addChainConfig(
    uint256 numNodes
  ) internal returns (CCIPHome.OCR3Node[] memory nodes) {
    return _addChainConfig(numNodes, 1);
  }

  function _makeBytes32Array(uint256 length, uint256 seed) internal pure returns (bytes32[] memory arr) {
    arr = new bytes32[](length);
    for (uint256 i = 0; i < length; i++) {
      arr[i] = keccak256(abi.encode(i, 1, seed));
    }
    return arr;
  }

  function _makeBytesArray(uint256 length, uint256 seed) internal pure returns (bytes[] memory arr) {
    arr = new bytes[](length);
    for (uint256 i = 0; i < length; i++) {
      arr[i] = abi.encode(keccak256(abi.encode(i, 1, seed)));
    }
    return arr;
  }

  function _addChainConfig(uint256 numNodes, uint8 fChain) internal returns (CCIPHome.OCR3Node[] memory nodes) {
    bytes32[] memory p2pIds = _makeBytes32Array(numNodes, 0);
    bytes[] memory signers = _makeBytesArray(numNodes, 10);
    bytes[] memory transmitters = _makeBytesArray(numNodes, 20);

    nodes = new CCIPHome.OCR3Node[](numNodes);

    for (uint256 i = 0; i < numNodes; i++) {
      nodes[i] = CCIPHome.OCR3Node({p2pId: p2pIds[i], signerKey: signers[i], transmitterKey: transmitters[i]});

      vm.mockCall(
        CAPABILITIES_REGISTRY,
        abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, p2pIds[i]),
        abi.encode(
          ICapabilitiesRegistry.NodeInfo({
            nodeOperatorId: 1,
            signer: bytes32(signers[i]),
            p2pId: p2pIds[i],
            encryptionPublicKey: keccak256("encryptionPublicKey"),
            hashedCapabilityIds: new bytes32[](0),
            configCount: uint32(1),
            workflowDONId: uint32(1),
            capabilitiesDONIds: new uint256[](0)
          })
        )
      );
    }
    // Add chain selector for chain 1.
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](1);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: p2pIds, fChain: fChain, config: bytes("config1")})
    });

    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(1, adds[0].chainConfig);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    return nodes;
  }

  function _getCorrectOCR3Config(uint8 numNodes, uint8 FRoleDON) internal returns (CCIPHome.OCR3Config memory) {
    CCIPHome.OCR3Node[] memory nodes = _addChainConfig(numNodes);

    return CCIPHome.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
      offrampAddress: abi.encode(keccak256(abi.encode("offramp"))),
      rmnHomeAddress: abi.encode(keccak256(abi.encode("rmnHome"))),
      chainSelector: 1,
      nodes: nodes,
      FRoleDON: FRoleDON,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });
  }

  function _getCorrectOCR3Config() internal returns (CCIPHome.OCR3Config memory) {
    return _getCorrectOCR3Config(4, 1);
  }

  // Successes.

  function test__validateConfig_Success() public {
    s_ccipHome.validateConfig(_getCorrectOCR3Config());
  }

  function test__validateConfigLessTransmittersThanSigners_Success() public {
    // fChain is 1, so there should be at least 4 transmitters.
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config(5, 1);
    config.nodes[1].transmitterKey = bytes("");

    s_ccipHome.validateConfig(config);
  }

  function test__validateConfigSmallerFChain_Success() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config(11, 3);

    // Set fChain to 2
    _addChainConfig(4, 2);

    s_ccipHome.validateConfig(config);
  }

  // Reverts

  function test__validateConfig_ChainSelectorNotSet_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.chainSelector = 0; // invalid

    vm.expectRevert(CCIPHome.ChainSelectorNotSet.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_OfframpAddressCannotBeZero_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.offrampAddress = ""; // invalid

    vm.expectRevert(CCIPHome.OfframpAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_ABIEncodedAddress_OfframpAddressCannotBeZero_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.offrampAddress = abi.encode(address(0)); // invalid

    vm.expectRevert(CCIPHome.OfframpAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_RMNHomeAddressCannotBeZero_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.rmnHomeAddress = ""; // invalid

    vm.expectRevert(CCIPHome.RMNHomeAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_ABIEncodedAddress_RMNHomeAddressCannotBeZero_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.rmnHomeAddress = abi.encode(address(0)); // invalid

    vm.expectRevert(CCIPHome.RMNHomeAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_ChainSelectorNotFound_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.chainSelector = 2; // not set

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ChainSelectorNotFound.selector, 2));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_NotEnoughTransmitters_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    uint256 numberOfTransmitters = 3;

    // 32 > 31 (max num oracles)
    CCIPHome.OCR3Node[] memory nodes = _addChainConfig(31);

    // truncate transmitters to < 3 * fChain + 1
    // since fChain is 1 in this case, we need to truncate to 3 transmitters.
    for (uint256 i = numberOfTransmitters; i < nodes.length; ++i) {
      nodes[i].transmitterKey = bytes("");
    }

    config.nodes = nodes;
    vm.expectRevert(abi.encodeWithSelector(CCIPHome.NotEnoughTransmitters.selector, numberOfTransmitters, 4));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_NotEnoughTransmittersEmptyAddresses_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes[0].transmitterKey = bytes("");

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.NotEnoughTransmitters.selector, 3, 4));
    s_ccipHome.validateConfig(config);

    // Zero out remaining transmitters to verify error changes
    for (uint256 i = 1; i < config.nodes.length; ++i) {
      config.nodes[i].transmitterKey = bytes("");
    }

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.NotEnoughTransmitters.selector, 0, 4));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_TooManySigners_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes = new CCIPHome.OCR3Node[](257);

    vm.expectRevert(CCIPHome.TooManySigners.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_FChainTooHigh_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.FRoleDON = 2; // too low

    // Set fChain to 3
    _addChainConfig(4, 3);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.FChainTooHigh.selector, 3, 2));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_FMustBePositive_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.FRoleDON = 0; // not positive

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.FChainTooHigh.selector, 1, 0));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_FTooHigh_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.FRoleDON = 2; // too high

    vm.expectRevert(CCIPHome.FTooHigh.selector);
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_ZeroP2PId_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes[1].p2pId = bytes32(0);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.InvalidNode.selector, config.nodes[1]));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_ZeroSignerKey_Reverts() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes[2].signerKey = bytes("");

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.InvalidNode.selector, config.nodes[2]));
    s_ccipHome.validateConfig(config);
  }

  function test__validateConfig_NodeNotInRegistry_Reverts() public {
    CCIPHome.OCR3Node[] memory nodes = _addChainConfig(4);
    bytes32 nonExistentP2PId = keccak256("notInRegistry");
    nodes[0].p2pId = nonExistentP2PId;

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, nonExistentP2PId),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 0,
          signer: bytes32(0),
          p2pId: bytes32(uint256(0)),
          encryptionPublicKey: keccak256("encryptionPublicKey"),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes = nodes;

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.NodeNotInRegistry.selector, nonExistentP2PId));
    s_ccipHome.validateConfig(config);
  }
}

contract CCIPHome_applyChainConfigUpdates is CCIPHomeTest {
  function setUp() public virtual override {
    s_ccipHome = new CCIPHomeHelper(CAPABILITIES_REGISTRY);
  }

  function test_applyChainConfigUpdates_addChainConfigs_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          encryptionPublicKey: keccak256("encryptionPublicKey"),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );
    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    CCIPHome.ChainConfigArgs[] memory configs = s_ccipHome.getAllChainConfigs(0, 2);
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");
    assertEq(s_ccipHome.getNumChainConfigurations(), 2, "total chain configs must be 2");
  }

  function test_getPaginatedCCIPHomes_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          encryptionPublicKey: keccak256("encryptionPublicKey"),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    CCIPHome.ChainConfigArgs[] memory configs = s_ccipHome.getAllChainConfigs(0, 2);
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");

    configs = s_ccipHome.getAllChainConfigs(0, 1);
    assertEq(configs.length, 1, "chain configs length must be 1");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");

    configs = s_ccipHome.getAllChainConfigs(0, 10);
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");

    configs = s_ccipHome.getAllChainConfigs(1, 1);
    assertEq(configs.length, 1, "chain configs length must be 1");

    configs = s_ccipHome.getAllChainConfigs(1, 2);
    assertEq(configs.length, 0, "chain configs length must be 0");
  }

  function test_applyChainConfigUpdates_removeChainConfigs_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          encryptionPublicKey: keccak256("encryptionPublicKey"),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    assertEq(s_ccipHome.getNumChainConfigurations(), 2, "total chain configs must be 2");

    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectEmit();
    emit CCIPHome.ChainConfigRemoved(1);
    s_ccipHome.applyChainConfigUpdates(removes, new CCIPHome.ChainConfigArgs[](0));

    assertEq(s_ccipHome.getNumChainConfigurations(), 1, "total chain configs must be 1");
  }

  // Reverts.

  function test_applyChainConfigUpdates_selectorNotFound_Reverts() public {
    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ChainSelectorNotFound.selector, 1));
    s_ccipHome.applyChainConfigUpdates(removes, new CCIPHome.ChainConfigArgs[](0));
  }

  function test_applyChainConfigUpdates_nodeNotInRegistry_Reverts() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](1);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: abi.encode(1, 2, 3)})
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 0,
          signer: bytes32(0),
          p2pId: bytes32(uint256(0)),
          encryptionPublicKey: keccak256("encryptionPublicKey"),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.NodeNotInRegistry.selector, chainReaders[0]));
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);
  }

  function test__applyChainConfigUpdates_FChainNotPositive_Reverts() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 0, config: bytes("config2")}) // bad fChain
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          encryptionPublicKey: keccak256("encryptionPublicKey"),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectRevert(CCIPHome.FChainMustBePositive.selector);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);
  }
}
