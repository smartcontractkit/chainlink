// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";

import {CCIPCapabilityConfiguration} from "../../capability/CCIPCapabilityConfiguration.sol";
import {ICapabilityRegistry} from "../../capability/interfaces/ICapabilityRegistry.sol";
import {CCIPCapabilityConfigurationHelper} from "../helpers/CCIPCapabilityConfigurationHelper.sol";

contract CCIPCapabilityConfigurationSetup is Test {
  address public constant OWNER = 0x82ae2B4F57CA5C1CBF8f744ADbD3697aD1a35AFe;
  address public constant CAPABILITY_REGISTRY = 0x272aF4BF7FBFc4944Ed59F914Cd864DfD912D55e;

  CCIPCapabilityConfigurationHelper public s_ccipCC;

  function setUp() public {
    changePrank(OWNER);
    s_ccipCC = new CCIPCapabilityConfigurationHelper(CAPABILITY_REGISTRY);
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
      arr[i] = abi.encodePacked(keccak256(abi.encode(i, 1, seed)));
    }
    return arr;
  }

  function _subset(bytes32[] memory arr, uint256 start, uint256 end) internal pure returns (bytes32[] memory) {
    bytes32[] memory subset = new bytes32[](end - start);
    for (uint256 i = start; i < end; i++) {
      subset[i - start] = arr[i];
    }
    return subset;
  }

  function _addChainConfig(uint256 numNodes)
    internal
    returns (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters)
  {
    p2pIds = _makeBytes32Array(numNodes, 0);
    signers = _makeBytesArray(numNodes, 10);
    transmitters = _makeBytesArray(numNodes, 20);
    for (uint256 i = 0; i < numNodes; i++) {
      vm.mockCall(
        CAPABILITY_REGISTRY,
        abi.encodeWithSelector(ICapabilityRegistry.getNode.selector, p2pIds[i]),
        abi.encode(
          ICapabilityRegistry.NodeInfo({
            nodeOperatorId: 1,
            signer: bytes32(signers[i]),
            p2pId: p2pIds[i],
            hashedCapabilityIds: new bytes32[](0)
          }),
          uint32(1)
        )
      );
    }
    // Add chain selector for chain 1.
    CCIPCapabilityConfiguration.ChainConfigInfo[] memory adds = new CCIPCapabilityConfiguration.ChainConfigInfo[](1);
    adds[0] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: p2pIds, fChain: 1, config: bytes("config1")})
    });

    vm.expectEmit();
    emit CCIPCapabilityConfiguration.ChainConfigSet(1, adds[0].chainConfig);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);

    return (p2pIds, signers, transmitters);
  }

  function test_getCapabilityConfiguration_Success() public {
    bytes memory capConfig = s_ccipCC.getCapabilityConfiguration(42 /* doesn't matter, not used */ );
    assertEq(capConfig.length, 0, "capability config length must be 0");
  }
}

contract CCIPCapabilityConfiguration_chainConfig is CCIPCapabilityConfigurationSetup {
  // Successes.

  function test_applyChainConfigUpdates_addChainConfigs_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPCapabilityConfiguration.ChainConfigInfo[] memory adds = new CCIPCapabilityConfiguration.ChainConfigInfo[](2);
    adds[0] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 2,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });

    vm.mockCall(
      CAPABILITY_REGISTRY,
      abi.encodeWithSelector(ICapabilityRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilityRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          hashedCapabilityIds: new bytes32[](0)
        }),
        uint32(1)
      )
    );

    vm.expectEmit();
    emit CCIPCapabilityConfiguration.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPCapabilityConfiguration.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);

    CCIPCapabilityConfiguration.ChainConfigInfo[] memory configs = s_ccipCC.getAllChainConfigs();
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");
  }

  function test_applyChainConfigUpdates_removeChainConfigs_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPCapabilityConfiguration.ChainConfigInfo[] memory adds = new CCIPCapabilityConfiguration.ChainConfigInfo[](2);
    adds[0] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 2,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });

    vm.mockCall(
      CAPABILITY_REGISTRY,
      abi.encodeWithSelector(ICapabilityRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilityRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          hashedCapabilityIds: new bytes32[](0)
        }),
        uint32(1)
      )
    );

    vm.expectEmit();
    emit CCIPCapabilityConfiguration.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPCapabilityConfiguration.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);

    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectEmit();
    emit CCIPCapabilityConfiguration.ChainConfigRemoved(1);
    s_ccipCC.applyChainConfigUpdates(removes, new CCIPCapabilityConfiguration.ChainConfigInfo[](0));
  }

  // Reverts.

  function test_applyChainConfigUpdates_selectorNotFound_Reverts() public {
    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.ChainSelectorNotFound.selector, 1));
    s_ccipCC.applyChainConfigUpdates(removes, new CCIPCapabilityConfiguration.ChainConfigInfo[](0));
  }

  function test_applyChainConfigUpdates_nodeNotInRegistry_Reverts() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPCapabilityConfiguration.ChainConfigInfo[] memory adds = new CCIPCapabilityConfiguration.ChainConfigInfo[](1);
    adds[0] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 1, config: abi.encode(1, 2, 3)})
    });

    vm.mockCall(
      CAPABILITY_REGISTRY,
      abi.encodeWithSelector(ICapabilityRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilityRegistry.NodeInfo({
          nodeOperatorId: 0,
          signer: bytes32(0),
          p2pId: bytes32(uint256(0)),
          hashedCapabilityIds: new bytes32[](0)
        }),
        uint32(1)
      )
    );

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.NodeNotInRegistry.selector, chainReaders[0]));
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);
  }

  function test__applyChainConfigUpdates_FChainNotPositive_Reverts() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPCapabilityConfiguration.ChainConfigInfo[] memory adds = new CCIPCapabilityConfiguration.ChainConfigInfo[](2);
    adds[0] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPCapabilityConfiguration.ChainConfigInfo({
      chainSelector: 2,
      chainConfig: CCIPCapabilityConfiguration.ChainConfig({readers: chainReaders, fChain: 0, config: bytes("config2")}) // bad fChain
    });

    vm.mockCall(
      CAPABILITY_REGISTRY,
      abi.encodeWithSelector(ICapabilityRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilityRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          hashedCapabilityIds: new bytes32[](0)
        }),
        uint32(1)
      )
    );

    vm.expectRevert(CCIPCapabilityConfiguration.FChainMustBePositive.selector);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);
  }
}

contract CCIPCapabilityConfiguration_validateConfig is CCIPCapabilityConfigurationSetup {
  // Successes.

  function test__validateConfig_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });
    s_ccipCC.validateConfig(config);
  }

  // Reverts.

  function test__validateConfig_ChainSelectorNotSet_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 0, // invalid
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.ChainSelectorNotSet.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_OfframpAddressCannotBeZero_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: bytes(""), // invalid
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.OfframpAddressCannotBeZero.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_ChainSelectorNotFound_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 2, // not set
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.ChainSelectorNotFound.selector, 2));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_TooManySigners_Reverts() public {
    // 32 > 31 (max num oracles)
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(32);

    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.TooManySigners.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_TooManyTransmitters_Reverts() public {
    // 32 > 31 (max num oracles)
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(32);

    // truncate signers but keep transmitters > 31
    assembly {
      mstore(signers, 30)
    }

    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.TooManyTransmitters.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_NotEnoughTransmitters_Reverts() public {
    // 32 > 31 (max num oracles)
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(31);

    // truncate transmitters to < 3 * fChain + 1
    // since fChain is 1 in this case, we need to truncate to 3 transmitters.
    assembly {
      mstore(transmitters, 3)
    }

    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.NotEnoughTransmitters.selector, 3, 4));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_FMustBePositive_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 0,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.FMustBePositive.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_FTooHigh_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 2,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.FTooHigh.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_P2PIdsLengthNotMatching_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // truncate the p2pIds length
    assembly {
      mstore(p2pIds, 3)
    }

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        CCIPCapabilityConfiguration.P2PIdsLengthNotMatching.selector, uint256(3), uint256(4), uint256(4)
      )
    );
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_TooManyBootstrapP2PIds_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _makeBytes32Array(5, 0), // too many bootstrap p2pIds, 5 > 4
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(CCIPCapabilityConfiguration.TooManyBootstrapP2PIds.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_NodeNotInRegistry_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    bytes32 nonExistentP2PId = keccak256("notInRegistry");
    p2pIds[0] = nonExistentP2PId;

    vm.mockCall(
      CAPABILITY_REGISTRY,
      abi.encodeWithSelector(ICapabilityRegistry.getNode.selector, nonExistentP2PId),
      abi.encode(
        ICapabilityRegistry.NodeInfo({
          nodeOperatorId: 0,
          signer: bytes32(0),
          p2pId: bytes32(uint256(0)),
          hashedCapabilityIds: new bytes32[](0)
        }),
        uint32(1)
      )
    );

    // Config is for 4 nodes, so f == 1.
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.NodeNotInRegistry.selector, nonExistentP2PId));
    s_ccipCC.validateConfig(config);
  }
}

contract CCIPCapabilityConfiguration_ConfigStateMachine is CCIPCapabilityConfigurationSetup {
  // Successful cases.

  function test__stateFromConfigLength_Success() public {
    uint256 configLen = 0;
    CCIPCapabilityConfiguration.ConfigState state = s_ccipCC.stateFromConfigLength(configLen);
    assertEq(uint256(state), uint256(CCIPCapabilityConfiguration.ConfigState.Init));

    configLen = 1;
    state = s_ccipCC.stateFromConfigLength(configLen);
    assertEq(uint256(state), uint256(CCIPCapabilityConfiguration.ConfigState.Running));

    configLen = 2;
    state = s_ccipCC.stateFromConfigLength(configLen);
    assertEq(uint256(state), uint256(CCIPCapabilityConfiguration.ConfigState.Staging));
  }

  function test__validateConfigStateTransition_Success() public {
    s_ccipCC.validateConfigStateTransition(
      CCIPCapabilityConfiguration.ConfigState.Init, CCIPCapabilityConfiguration.ConfigState.Running
    );

    s_ccipCC.validateConfigStateTransition(
      CCIPCapabilityConfiguration.ConfigState.Running, CCIPCapabilityConfiguration.ConfigState.Staging
    );

    s_ccipCC.validateConfigStateTransition(
      CCIPCapabilityConfiguration.ConfigState.Staging, CCIPCapabilityConfiguration.ConfigState.Running
    );
  }

  function test__computeConfigDigest_Success() public {
    // config digest must change upon:
    // - ocr config change (e.g plugin type, chain selector, etc.)
    // - don id change
    // - config count change
    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(2, 10);
    bytes[] memory transmitters = _makeBytesArray(2, 20);
    CCIPCapabilityConfiguration.OCR3Config memory config = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });
    uint32 donId = 1;
    uint32 configCount = 1;

    bytes32 configDigest1 = s_ccipCC.computeConfigDigest(donId, configCount, config);

    donId = 2;
    bytes32 configDigest2 = s_ccipCC.computeConfigDigest(donId, configCount, config);

    donId = 1;
    configCount = 2;
    bytes32 configDigest3 = s_ccipCC.computeConfigDigest(donId, configCount, config);

    configCount = 1;
    config.pluginType = CCIPCapabilityConfiguration.PluginType.Execution;
    bytes32 configDigest4 = s_ccipCC.computeConfigDigest(donId, configCount, config);

    assertNotEq(configDigest1, configDigest2, "config digests 1 and 2 must not match");
    assertNotEq(configDigest1, configDigest3, "config digests 1 and 3 must not match");
    assertNotEq(configDigest1, configDigest4, "config digests 1 and 4 must not match");

    assertNotEq(configDigest2, configDigest3, "config digests 2 and 3 must not match");
    assertNotEq(configDigest2, configDigest4, "config digests 2 and 4 must not match");
  }

  function test_Fuzz__groupByPluginType_Success(uint256 numCommitCfgs, uint256 numExecCfgs) public {
    vm.assume(numCommitCfgs >= 0 && numCommitCfgs < 3);
    vm.assume(numExecCfgs >= 0 && numExecCfgs < 3);

    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(4, 10);
    bytes[] memory transmitters = _makeBytesArray(4, 20);
    CCIPCapabilityConfiguration.OCR3Config[] memory cfgs =
      new CCIPCapabilityConfiguration.OCR3Config[](numCommitCfgs + numExecCfgs);
    for (uint256 i = 0; i < numCommitCfgs; i++) {
      cfgs[i] = CCIPCapabilityConfiguration.OCR3Config({
        pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
        offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
        chainSelector: 1,
        bootstrapP2PIds: _subset(p2pIds, 0, 1),
        p2pIds: p2pIds,
        signers: signers,
        transmitters: transmitters,
        F: 1,
        offchainConfigVersion: 30,
        offchainConfig: abi.encode("commit", i)
      });
    }
    for (uint256 i = 0; i < numExecCfgs; i++) {
      cfgs[numCommitCfgs + i] = CCIPCapabilityConfiguration.OCR3Config({
        pluginType: CCIPCapabilityConfiguration.PluginType.Execution,
        offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
        chainSelector: 1,
        bootstrapP2PIds: _subset(p2pIds, 0, 1),
        p2pIds: p2pIds,
        signers: signers,
        transmitters: transmitters,
        F: 1,
        offchainConfigVersion: 30,
        offchainConfig: abi.encode("exec", numCommitCfgs + i)
      });
    }
    (
      CCIPCapabilityConfiguration.OCR3Config[] memory commitCfgs,
      CCIPCapabilityConfiguration.OCR3Config[] memory execCfgs
    ) = s_ccipCC.groupByPluginType(cfgs);

    assertEq(commitCfgs.length, numCommitCfgs, "commitCfgs length must match");
    assertEq(execCfgs.length, numExecCfgs, "execCfgs length must match");
    for (uint256 i = 0; i < commitCfgs.length; i++) {
      assertEq(
        uint8(commitCfgs[i].pluginType),
        uint8(CCIPCapabilityConfiguration.PluginType.Commit),
        "plugin type must be commit"
      );
      assertEq(commitCfgs[i].offchainConfig, abi.encode("commit", i), "offchain config must match");
    }
    for (uint256 i = 0; i < execCfgs.length; i++) {
      assertEq(
        uint8(execCfgs[i].pluginType),
        uint8(CCIPCapabilityConfiguration.PluginType.Execution),
        "plugin type must be execution"
      );
      assertEq(execCfgs[i].offchainConfig, abi.encode("exec", numCommitCfgs + i), "offchain config must match");
    }
  }

  function test__computeNewConfigWithMeta_InitToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](0);
    CCIPCapabilityConfiguration.OCR3Config[] memory newConfig = new CCIPCapabilityConfiguration.OCR3Config[](1);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.ConfigState currentState = CCIPCapabilityConfiguration.ConfigState.Init;
    CCIPCapabilityConfiguration.ConfigState newState = CCIPCapabilityConfiguration.ConfigState.Running;
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfigWithMeta =
      s_ccipCC.computeNewConfigWithMeta(donId, currentConfig, newConfig, currentState, newState);
    assertEq(newConfigWithMeta.length, 1, "new config with meta length must be 1");
    assertEq(newConfigWithMeta[0].configCount, uint64(1), "config count must be 1");
    assertEq(uint8(newConfigWithMeta[0].config.pluginType), uint8(newConfig[0].pluginType), "plugin type must match");
    assertEq(newConfigWithMeta[0].config.offchainConfig, newConfig[0].offchainConfig, "offchain config must match");
    assertEq(
      newConfigWithMeta[0].configDigest,
      s_ccipCC.computeConfigDigest(donId, 1, newConfig[0]),
      "config digest must match"
    );

    // This ensures that the test case is using correct inputs.
    s_ccipCC.validateConfigTransition(currentConfig, newConfigWithMeta);
  }

  function test__computeNewConfigWithMeta_RunningToStaging_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    CCIPCapabilityConfiguration.OCR3Config[] memory newConfig = new CCIPCapabilityConfiguration.OCR3Config[](2);
    // existing blue config first.
    newConfig[0] = blueConfig;
    // green config next.
    newConfig[1] = greenConfig;

    CCIPCapabilityConfiguration.ConfigState currentState = CCIPCapabilityConfiguration.ConfigState.Running;
    CCIPCapabilityConfiguration.ConfigState newState = CCIPCapabilityConfiguration.ConfigState.Staging;

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfigWithMeta =
      s_ccipCC.computeNewConfigWithMeta(donId, currentConfig, newConfig, currentState, newState);
    assertEq(newConfigWithMeta.length, 2, "new config with meta length must be 2");

    assertEq(newConfigWithMeta[0].configCount, uint64(1), "config count of blue must be 1");
    assertEq(
      uint8(newConfigWithMeta[0].config.pluginType), uint8(blueConfig.pluginType), "plugin type of blue must match"
    );
    assertEq(
      newConfigWithMeta[0].config.offchainConfig, blueConfig.offchainConfig, "offchain config of blue must match"
    );
    assertEq(
      newConfigWithMeta[0].configDigest,
      s_ccipCC.computeConfigDigest(donId, 1, blueConfig),
      "config digest of blue must match"
    );

    assertEq(newConfigWithMeta[1].configCount, uint64(2), "config count of green must be 2");
    assertEq(
      uint8(newConfigWithMeta[1].config.pluginType), uint8(greenConfig.pluginType), "plugin type of green must match"
    );
    assertEq(
      newConfigWithMeta[1].config.offchainConfig, greenConfig.offchainConfig, "offchain config of green must match"
    );
    assertEq(
      newConfigWithMeta[1].configDigest,
      s_ccipCC.computeConfigDigest(donId, 2, greenConfig),
      "config digest of green must match"
    );

    // This ensures that the test case is using correct inputs.
    s_ccipCC.validateConfigTransition(currentConfig, newConfigWithMeta);
  }

  function test__computeNewConfigWithMeta_StagingToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](2);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    currentConfig[1] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory newConfig = new CCIPCapabilityConfiguration.OCR3Config[](1);
    newConfig[0] = greenConfig;

    CCIPCapabilityConfiguration.ConfigState currentState = CCIPCapabilityConfiguration.ConfigState.Staging;
    CCIPCapabilityConfiguration.ConfigState newState = CCIPCapabilityConfiguration.ConfigState.Running;

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfigWithMeta =
      s_ccipCC.computeNewConfigWithMeta(donId, currentConfig, newConfig, currentState, newState);

    assertEq(newConfigWithMeta.length, 1, "new config with meta length must be 1");
    assertEq(newConfigWithMeta[0].configCount, uint64(2), "config count must be 2");
    assertEq(uint8(newConfigWithMeta[0].config.pluginType), uint8(greenConfig.pluginType), "plugin type must match");
    assertEq(newConfigWithMeta[0].config.offchainConfig, greenConfig.offchainConfig, "offchain config must match");
    assertEq(
      newConfigWithMeta[0].configDigest, s_ccipCC.computeConfigDigest(donId, 2, greenConfig), "config digest must match"
    );

    // This ensures that the test case is using correct inputs.
    s_ccipCC.validateConfigTransition(currentConfig, newConfigWithMeta);
  }

  function test__validateConfigTransition_InitToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](0);

    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_RunningToStaging_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](2);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    newConfig[1] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_StagingToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](2);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    currentConfig[1] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  // Reverts.

  function test_Fuzz__stateFromConfigLength_Reverts(uint256 configLen) public {
    vm.assume(configLen > 2);
    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.InvalidConfigLength.selector, configLen));
    s_ccipCC.stateFromConfigLength(configLen);
  }

  function test__groupByPluginType_threeCommitConfigs_Reverts() public {
    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(4, 10);
    bytes[] memory transmitters = _makeBytesArray(4, 20);
    CCIPCapabilityConfiguration.OCR3Config[] memory cfgs = new CCIPCapabilityConfiguration.OCR3Config[](3);
    for (uint256 i = 0; i < 3; i++) {
      cfgs[i] = CCIPCapabilityConfiguration.OCR3Config({
        pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
        offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
        chainSelector: 1,
        bootstrapP2PIds: _subset(p2pIds, 0, 1),
        p2pIds: p2pIds,
        signers: signers,
        transmitters: transmitters,
        F: 1,
        offchainConfigVersion: 30,
        offchainConfig: abi.encode("commit", i)
      });
    }
    vm.expectRevert();
    s_ccipCC.groupByPluginType(cfgs);
  }

  function test__groupByPluginType_threeExecutionConfigs_Reverts() public {
    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(4, 10);
    bytes[] memory transmitters = _makeBytesArray(4, 20);
    CCIPCapabilityConfiguration.OCR3Config[] memory cfgs = new CCIPCapabilityConfiguration.OCR3Config[](3);
    for (uint256 i = 0; i < 3; i++) {
      cfgs[i] = CCIPCapabilityConfiguration.OCR3Config({
        pluginType: CCIPCapabilityConfiguration.PluginType.Execution,
        offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
        chainSelector: 1,
        bootstrapP2PIds: _subset(p2pIds, 0, 1),
        p2pIds: p2pIds,
        signers: signers,
        transmitters: transmitters,
        F: 1,
        offchainConfigVersion: 30,
        offchainConfig: abi.encode("exec", i)
      });
    }
    vm.expectRevert();
    s_ccipCC.groupByPluginType(cfgs);
  }

  function test__groupByPluginType_TooManyOCR3Configs_Reverts() public {
    CCIPCapabilityConfiguration.OCR3Config[] memory cfgs = new CCIPCapabilityConfiguration.OCR3Config[](5);
    vm.expectRevert(CCIPCapabilityConfiguration.TooManyOCR3Configs.selector);
    s_ccipCC.groupByPluginType(cfgs);
  }

  function test__validateConfigTransition_InitToRunning_WrongConfigCount_Reverts() public {
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 0,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](0);

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.WrongConfigCount.selector, 0, 1));
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_RunningToStaging_WrongConfigDigestBlueGreen_Reverts() public {
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](2);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 3, blueConfig) // wrong config digest (due to diff config count)
    });
    newConfig[1] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        CCIPCapabilityConfiguration.WrongConfigDigestBlueGreen.selector,
        s_ccipCC.computeConfigDigest(donId, 3, blueConfig),
        s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
      )
    );
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_RunningToStaging_WrongConfigCount_Reverts() public {
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](2);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    newConfig[1] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 3, // wrong config count
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 3, greenConfig)
    });

    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.WrongConfigCount.selector, 3, 2));
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_StagingToRunning_WrongConfigDigest_Reverts() public {
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(_makeBytes32Array(4, 0), 0, 1),
      p2pIds: _makeBytes32Array(4, 0),
      signers: _makeBytesArray(4, 10),
      transmitters: _makeBytesArray(4, 20),
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](2);
    currentConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    currentConfig[1] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPCapabilityConfiguration.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 3, greenConfig) // wrong config digest
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        CCIPCapabilityConfiguration.WrongConfigDigest.selector,
        s_ccipCC.computeConfigDigest(donId, 3, greenConfig),
        s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
      )
    );
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_NonExistentConfigTransition_Reverts() public {
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory currentConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](3);
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory newConfig =
      new CCIPCapabilityConfiguration.OCR3ConfigWithMeta[](1);
    vm.expectRevert(CCIPCapabilityConfiguration.NonExistentConfigTransition.selector);
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }
}

contract CCIPCapabilityConfiguration__updatePluginConfig is CCIPCapabilityConfigurationSetup {
  // Successes.

  function test__updatePluginConfig_InitToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory configs = new CCIPCapabilityConfiguration.OCR3Config[](1);
    configs[0] = blueConfig;

    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, configs);

    // should see the updated config in the contract state.
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedConfig =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit);
    assertEq(storedConfig.length, 1, "don config length must be 1");
    assertEq(storedConfig[0].configCount, uint64(1), "config count must be 1");
    assertEq(uint256(storedConfig[0].config.pluginType), uint256(blueConfig.pluginType), "plugin type must match");
  }

  function test__updatePluginConfig_RunningToStaging_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // add blue config.
    uint32 donId = 1;
    CCIPCapabilityConfiguration.PluginType pluginType = CCIPCapabilityConfiguration.PluginType.Commit;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory startConfigs = new CCIPCapabilityConfiguration.OCR3Config[](1);
    startConfigs[0] = blueConfig;

    // add blue AND green config to indicate an update.
    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, startConfigs);
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory blueAndGreen = new CCIPCapabilityConfiguration.OCR3Config[](2);
    blueAndGreen[0] = blueConfig;
    blueAndGreen[1] = greenConfig;

    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, blueAndGreen);

    // should see the updated config in the contract state.
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedConfig =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit);
    assertEq(storedConfig.length, 2, "don config length must be 2");
    // 0 index is blue config, 1 index is green config.
    assertEq(storedConfig[1].configCount, uint64(2), "config count must be 2");
    assertEq(
      uint256(storedConfig[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must match"
    );
    assertEq(
      uint256(storedConfig[1].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must match"
    );
    assertEq(storedConfig[0].config.offchainConfig, bytes("commit"), "blue offchain config must match");
    assertEq(storedConfig[1].config.offchainConfig, bytes("commit-new"), "green offchain config must match");
  }

  function test__updatePluginConfig_StagingToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // add blue config.
    uint32 donId = 1;
    CCIPCapabilityConfiguration.PluginType pluginType = CCIPCapabilityConfiguration.PluginType.Commit;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory startConfigs = new CCIPCapabilityConfiguration.OCR3Config[](1);
    startConfigs[0] = blueConfig;

    // add blue AND green config to indicate an update.
    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, startConfigs);
    CCIPCapabilityConfiguration.OCR3Config memory greenConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit-new")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory blueAndGreen = new CCIPCapabilityConfiguration.OCR3Config[](2);
    blueAndGreen[0] = blueConfig;
    blueAndGreen[1] = greenConfig;

    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, blueAndGreen);

    // should see the updated config in the contract state.
    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedConfig =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit);
    assertEq(storedConfig.length, 2, "don config length must be 2");
    // 0 index is blue config, 1 index is green config.
    assertEq(storedConfig[1].configCount, uint64(2), "config count must be 2");
    assertEq(
      uint256(storedConfig[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must match"
    );
    assertEq(
      uint256(storedConfig[1].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must match"
    );
    assertEq(storedConfig[0].config.offchainConfig, bytes("commit"), "blue offchain config must match");
    assertEq(storedConfig[1].config.offchainConfig, bytes("commit-new"), "green offchain config must match");

    // promote green to blue.
    CCIPCapabilityConfiguration.OCR3Config[] memory promote = new CCIPCapabilityConfiguration.OCR3Config[](1);
    promote[0] = greenConfig;

    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, promote);

    // should see the updated config in the contract state.
    storedConfig = s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit);
    assertEq(storedConfig.length, 1, "don config length must be 1");
    assertEq(storedConfig[0].configCount, uint64(2), "config count must be 2");
    assertEq(
      uint256(storedConfig[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must match"
    );
    assertEq(storedConfig[0].config.offchainConfig, bytes("commit-new"), "green offchain config must match");
  }

  // Reverts.
  function test__updatePluginConfig_InvalidConfigLength_Reverts() public {
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config[] memory newConfig = new CCIPCapabilityConfiguration.OCR3Config[](3);
    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.InvalidConfigLength.selector, uint256(3)));
    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, newConfig);
  }

  function test__updatePluginConfig_InvalidConfigStateTransition_Reverts() public {
    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config[] memory newConfig = new CCIPCapabilityConfiguration.OCR3Config[](2);
    // 0 -> 2 is an invalid state transition.
    vm.expectRevert(abi.encodeWithSelector(CCIPCapabilityConfiguration.InvalidConfigStateTransition.selector, 0, 2));
    s_ccipCC.updatePluginConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit, newConfig);
  }
}

contract CCIPCapabilityConfiguration_beforeCapabilityConfigSet is CCIPCapabilityConfigurationSetup {
  // Successes.
  function test_beforeCapabilityConfigSet_ZeroLengthConfig_Success() public {
    changePrank(CAPABILITY_REGISTRY);

    CCIPCapabilityConfiguration.OCR3Config[] memory configs = new CCIPCapabilityConfiguration.OCR3Config[](0);
    bytes memory encodedConfigs = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encodedConfigs, 1, 1);
  }

  function test_beforeCapabilityConfigSet_CommitConfigOnly_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    changePrank(CAPABILITY_REGISTRY);

    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory configs = new CCIPCapabilityConfiguration.OCR3Config[](1);
    configs[0] = blueConfig;

    bytes memory encoded = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encoded, 1, donId);

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedConfigs =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit);
    assertEq(storedConfigs.length, 1, "config length must be 1");
    assertEq(storedConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedConfigs[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must be commit"
    );
  }

  function test_beforeCapabilityConfigSet_ExecConfigOnly_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    changePrank(CAPABILITY_REGISTRY);

    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Execution,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("exec")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory configs = new CCIPCapabilityConfiguration.OCR3Config[](1);
    configs[0] = blueConfig;

    bytes memory encoded = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encoded, 1, donId);

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedConfigs =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Execution);
    assertEq(storedConfigs.length, 1, "config length must be 1");
    assertEq(storedConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedConfigs[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Execution),
      "plugin type must be execution"
    );
  }

  function test_beforeCapabilityConfigSet_CommitAndExecConfig_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    changePrank(CAPABILITY_REGISTRY);

    uint32 donId = 1;
    CCIPCapabilityConfiguration.OCR3Config memory blueCommitConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("commit")
    });
    CCIPCapabilityConfiguration.OCR3Config memory blueExecConfig = CCIPCapabilityConfiguration.OCR3Config({
      pluginType: CCIPCapabilityConfiguration.PluginType.Execution,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 1),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("exec")
    });
    CCIPCapabilityConfiguration.OCR3Config[] memory configs = new CCIPCapabilityConfiguration.OCR3Config[](2);
    configs[0] = blueExecConfig;
    configs[1] = blueCommitConfig;

    bytes memory encoded = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encoded, 1, donId);

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedExecConfigs =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Execution);
    assertEq(storedExecConfigs.length, 1, "config length must be 1");
    assertEq(storedExecConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedExecConfigs[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Execution),
      "plugin type must be execution"
    );

    CCIPCapabilityConfiguration.OCR3ConfigWithMeta[] memory storedCommitConfigs =
      s_ccipCC.getOCRConfig(donId, CCIPCapabilityConfiguration.PluginType.Commit);
    assertEq(storedCommitConfigs.length, 1, "config length must be 1");
    assertEq(storedCommitConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedCommitConfigs[0].config.pluginType),
      uint256(CCIPCapabilityConfiguration.PluginType.Commit),
      "plugin type must be commit"
    );
  }

  // Reverts.

  function test_beforeCapabilityConfigSet_OnlyCapabilityRegistryCanCall_Reverts() public {
    bytes32[] memory nodes = new bytes32[](0);
    bytes memory config = bytes("");
    uint64 configCount = 1;
    uint32 donId = 1;
    vm.expectRevert(CCIPCapabilityConfiguration.OnlyCapabilityRegistryCanCall.selector);
    s_ccipCC.beforeCapabilityConfigSet(nodes, config, configCount, donId);
  }
}
