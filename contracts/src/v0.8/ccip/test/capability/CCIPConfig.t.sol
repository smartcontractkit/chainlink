// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {SortedSetValidationUtil} from "../../../shared/util/SortedSetValidationUtil.sol";
import {CCIPConfig} from "../../capability/CCIPConfig.sol";
import {ICapabilitiesRegistry} from "../../capability/interfaces/ICapabilitiesRegistry.sol";
import {CCIPConfigTypes} from "../../capability/libraries/CCIPConfigTypes.sol";
import {Internal} from "../../libraries/Internal.sol";
import {CCIPConfigHelper} from "../helpers/CCIPConfigHelper.sol";
import {Test} from "forge-std/Test.sol";

contract CCIPConfigSetup is Test {
  address public constant OWNER = 0x82ae2B4F57CA5C1CBF8f744ADbD3697aD1a35AFe;
  address public constant CAPABILITIES_REGISTRY = 0x272aF4BF7FBFc4944Ed59F914Cd864DfD912D55e;

  CCIPConfigHelper public s_ccipCC;

  function setUp() public {
    changePrank(OWNER);
    s_ccipCC = new CCIPConfigHelper(CAPABILITIES_REGISTRY);
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

  //TODO: Use OZ's Arrays.sort when we upgrade to OZ v5
  function _sort(bytes32[] memory arr, int256 left, int256 right) private pure {
    int256 i = left;
    int256 j = right;
    if (i == j) return;
    bytes32 pivot = arr[uint256(left + (right - left) / 2)];
    while (i <= j) {
      while (arr[uint256(i)] < pivot) i++;
      while (pivot < arr[uint256(j)]) j--;
      if (i <= j) {
        (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
        i++;
        j--;
      }
    }
    if (left < j) _sort(arr, left, j);
    if (i < right) _sort(arr, i, right);
  }

  function _addChainConfig(uint256 numNodes)
    internal
    returns (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters)
  {
    p2pIds = _makeBytes32Array(numNodes, 0);
    _sort(p2pIds, 0, int256(numNodes - 1));
    signers = _makeBytesArray(numNodes, 10);
    transmitters = _makeBytesArray(numNodes, 20);
    for (uint256 i = 0; i < numNodes; i++) {
      vm.mockCall(
        CAPABILITIES_REGISTRY,
        abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, p2pIds[i]),
        abi.encode(
          ICapabilitiesRegistry.NodeInfo({
            nodeOperatorId: 1,
            signer: bytes32(signers[i]),
            p2pId: p2pIds[i],
            hashedCapabilityIds: new bytes32[](0),
            configCount: uint32(1),
            workflowDONId: uint32(1),
            capabilitiesDONIds: new uint256[](0)
          })
        )
      );
    }
    // Add chain selector for chain 1.
    CCIPConfigTypes.ChainConfigInfo[] memory adds = new CCIPConfigTypes.ChainConfigInfo[](1);
    adds[0] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: p2pIds, fChain: 1, config: bytes("config1")})
    });

    vm.expectEmit();
    emit CCIPConfig.ChainConfigSet(1, adds[0].chainConfig);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);

    return (p2pIds, signers, transmitters);
  }

  function test_getCapabilityConfiguration_Success() public view {
    bytes memory capConfig = s_ccipCC.getCapabilityConfiguration(42 /* doesn't matter, not used */ );
    assertEq(capConfig.length, 0, "capability config length must be 0");
  }
}

contract CCIPConfig_constructor is Test {
  // Successes.

  function test_constructor_Success() public {
    address capabilitiesRegistry = makeAddr("capabilitiesRegistry");
    CCIPConfigHelper ccipCC = new CCIPConfigHelper(capabilitiesRegistry);
    assertEq(address(ccipCC.getCapabilityRegistry()), capabilitiesRegistry);
    assertEq(ccipCC.typeAndVersion(), "CCIPConfig 1.6.0-dev");
  }

  // Reverts.

  function test_constructor_ZeroAddressNotAllowed_Revert() public {
    vm.expectRevert(CCIPConfig.ZeroAddressNotAllowed.selector);
    new CCIPConfigHelper(address(0));
  }
}

contract CCIPConfig_chainConfig is CCIPConfigSetup {
  // Successes.

  function test_applyChainConfigUpdates_addChainConfigs_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPConfigTypes.ChainConfigInfo[] memory adds = new CCIPConfigTypes.ChainConfigInfo[](2);
    adds[0] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 2,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectEmit();
    emit CCIPConfig.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPConfig.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);

    CCIPConfigTypes.ChainConfigInfo[] memory configs = s_ccipCC.getAllChainConfigs();
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");
  }

  function test_applyChainConfigUpdates_removeChainConfigs_Success() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPConfigTypes.ChainConfigInfo[] memory adds = new CCIPConfigTypes.ChainConfigInfo[](2);
    adds[0] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 2,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectEmit();
    emit CCIPConfig.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPConfig.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);

    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectEmit();
    emit CCIPConfig.ChainConfigRemoved(1);
    s_ccipCC.applyChainConfigUpdates(removes, new CCIPConfigTypes.ChainConfigInfo[](0));
  }

  // Reverts.

  function test_applyChainConfigUpdates_selectorNotFound_Reverts() public {
    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.ChainSelectorNotFound.selector, 1));
    s_ccipCC.applyChainConfigUpdates(removes, new CCIPConfigTypes.ChainConfigInfo[](0));
  }

  function test_applyChainConfigUpdates_nodeNotInRegistry_Reverts() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPConfigTypes.ChainConfigInfo[] memory adds = new CCIPConfigTypes.ChainConfigInfo[](1);
    adds[0] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 1, config: abi.encode(1, 2, 3)})
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 0,
          signer: bytes32(0),
          p2pId: bytes32(uint256(0)),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.NodeNotInRegistry.selector, chainReaders[0]));
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);
  }

  function test__applyChainConfigUpdates_FChainNotPositive_Reverts() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPConfigTypes.ChainConfigInfo[] memory adds = new CCIPConfigTypes.ChainConfigInfo[](2);
    adds[0] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 1,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPConfigTypes.ChainConfigInfo({
      chainSelector: 2,
      chainConfig: CCIPConfigTypes.ChainConfig({readers: chainReaders, fChain: 0, config: bytes("config2")}) // bad fChain
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, chainReaders[0]),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 1,
          signer: bytes32(uint256(1)),
          p2pId: chainReaders[0],
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    vm.expectRevert(CCIPConfig.FChainMustBePositive.selector);
    s_ccipCC.applyChainConfigUpdates(new uint64[](0), adds);
  }
}

contract CCIPConfig_validateConfig is CCIPConfigSetup {
  // Successes.

  function test__validateConfig_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.ChainSelectorNotSet.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_OfframpAddressCannotBeZero_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.OfframpAddressCannotBeZero.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_ChainSelectorNotFound_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.ChainSelectorNotFound.selector, 2));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_TooManySigners_Reverts() public {
    // 32 > 31 (max num oracles)
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(32);

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.TooManySigners.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_TooManyTransmitters_Reverts() public {
    // 32 > 31 (max num oracles)
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(32);

    // truncate signers but keep transmitters > 31
    assembly {
      mstore(signers, 30)
    }

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.TooManyTransmitters.selector);
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

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.NotEnoughTransmitters.selector, 3, 4));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_FMustBePositive_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.FMustBePositive.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_FTooHigh_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.FTooHigh.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_P2PIdsLengthNotMatching_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // truncate the p2pIds length
    assembly {
      mstore(p2pIds, 3)
    }

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
      abi.encodeWithSelector(CCIPConfig.P2PIdsLengthNotMatching.selector, uint256(3), uint256(4), uint256(4))
    );
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_TooManyBootstrapP2PIds_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(CCIPConfig.TooManyBootstrapP2PIds.selector);
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_NodeNotInRegistry_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    bytes32 nonExistentP2PId = keccak256("notInRegistry");
    p2pIds[0] = nonExistentP2PId;

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(ICapabilitiesRegistry.getNode.selector, nonExistentP2PId),
      abi.encode(
        ICapabilitiesRegistry.NodeInfo({
          nodeOperatorId: 0,
          signer: bytes32(0),
          p2pId: bytes32(uint256(0)),
          hashedCapabilityIds: new bytes32[](0),
          configCount: uint32(1),
          workflowDONId: uint32(1),
          capabilitiesDONIds: new uint256[](0)
        })
      )
    );

    // Config is for 4 nodes, so f == 1.
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.NodeNotInRegistry.selector, nonExistentP2PId));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_P2PIdsNotSorted_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // Config is for 4 nodes, so f == 1.

    //swapping two adjacent p2pIds to make it unsorted
    (p2pIds[2], p2pIds[3]) = (p2pIds[3], p2pIds[2]);

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, p2pIds));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_BootstrapP2PIdsNotSorted_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // Config is for 4 nodes, so f == 1.

    bytes32[] memory bootstrapP2PIds = _subset(p2pIds, 0, 2);

    //swapping bootstrapP2PIds to make it unsorted
    (bootstrapP2PIds[0], bootstrapP2PIds[1]) = (bootstrapP2PIds[1], bootstrapP2PIds[0]);

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: bootstrapP2PIds,
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, bootstrapP2PIds));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_P2PIdsHasDuplicates_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // Config is for 4 nodes, so f == 1.

    //forcing duplicate p2pIds
    p2pIds[1] = p2pIds[2];

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: _subset(p2pIds, 0, 2),
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, p2pIds));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_BootstrapP2PIdsHasDuplicates_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // Config is for 4 nodes, so f == 1.

    bytes32[] memory bootstrapP2PIds = _subset(p2pIds, 0, 2);
    //forcing duplicate bootstrapP2PIds
    bootstrapP2PIds[1] = bootstrapP2PIds[0];

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: bootstrapP2PIds,
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASortedSet.selector, bootstrapP2PIds));
    s_ccipCC.validateConfig(config);
  }

  function test__validateConfig_BootstrapP2PIdsNotASubsetOfP2PIds_Reverts() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // Config is for 4 nodes, so f == 1.

    //forcing invalid bootstrapP2PIds where the bootstrapP2PIds is sorted, but one of the element is not in the p2pIdsSet
    bytes32[] memory bootstrapP2PIds = _subset(p2pIds, 0, 2);
    p2pIds[1] = bytes32(uint256(p2pIds[0]) + 100);

    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
      offrampAddress: abi.encodePacked(keccak256(abi.encode("offramp"))),
      chainSelector: 1,
      bootstrapP2PIds: bootstrapP2PIds,
      p2pIds: p2pIds,
      signers: signers,
      transmitters: transmitters,
      F: 1,
      offchainConfigVersion: 30,
      offchainConfig: bytes("offchainConfig")
    });

    vm.expectRevert(abi.encodeWithSelector(SortedSetValidationUtil.NotASubset.selector, bootstrapP2PIds, p2pIds));
    s_ccipCC.validateConfig(config);
  }
}

contract CCIPConfig_ConfigStateMachine is CCIPConfigSetup {
  // Successful cases.

  function test__stateFromConfigLength_Success() public view {
    uint256 configLen = 0;
    CCIPConfigTypes.ConfigState state = s_ccipCC.stateFromConfigLength(configLen);
    assertEq(uint256(state), uint256(CCIPConfigTypes.ConfigState.Init));

    configLen = 1;
    state = s_ccipCC.stateFromConfigLength(configLen);
    assertEq(uint256(state), uint256(CCIPConfigTypes.ConfigState.Running));

    configLen = 2;
    state = s_ccipCC.stateFromConfigLength(configLen);
    assertEq(uint256(state), uint256(CCIPConfigTypes.ConfigState.Staging));
  }

  function test__validateConfigStateTransition_Success() public view {
    s_ccipCC.validateConfigStateTransition(CCIPConfigTypes.ConfigState.Init, CCIPConfigTypes.ConfigState.Running);

    s_ccipCC.validateConfigStateTransition(CCIPConfigTypes.ConfigState.Running, CCIPConfigTypes.ConfigState.Staging);

    s_ccipCC.validateConfigStateTransition(CCIPConfigTypes.ConfigState.Staging, CCIPConfigTypes.ConfigState.Running);
  }

  function test__computeConfigDigest_Success() public view {
    // config digest must change upon:
    // - ocr config change (e.g plugin type, chain selector, etc.)
    // - don id change
    // - config count change
    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(2, 10);
    bytes[] memory transmitters = _makeBytesArray(2, 20);
    CCIPConfigTypes.OCR3Config memory config = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    config.pluginType = Internal.OCRPluginType.Execution;
    bytes32 configDigest4 = s_ccipCC.computeConfigDigest(donId, configCount, config);

    assertNotEq(configDigest1, configDigest2, "config digests 1 and 2 must not match");
    assertNotEq(configDigest1, configDigest3, "config digests 1 and 3 must not match");
    assertNotEq(configDigest1, configDigest4, "config digests 1 and 4 must not match");

    assertNotEq(configDigest2, configDigest3, "config digests 2 and 3 must not match");
    assertNotEq(configDigest2, configDigest4, "config digests 2 and 4 must not match");
  }

  function test_Fuzz__groupByPluginType_Success(uint256 numCommitCfgs, uint256 numExecCfgs) public view {
    numCommitCfgs = bound(numCommitCfgs, 0, 2);
    numExecCfgs = bound(numExecCfgs, 0, 2);

    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(4, 10);
    bytes[] memory transmitters = _makeBytesArray(4, 20);
    CCIPConfigTypes.OCR3Config[] memory cfgs = new CCIPConfigTypes.OCR3Config[](numCommitCfgs + numExecCfgs);
    for (uint256 i = 0; i < numCommitCfgs; i++) {
      cfgs[i] = CCIPConfigTypes.OCR3Config({
        pluginType: Internal.OCRPluginType.Commit,
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
      cfgs[numCommitCfgs + i] = CCIPConfigTypes.OCR3Config({
        pluginType: Internal.OCRPluginType.Execution,
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
    (CCIPConfigTypes.OCR3Config[] memory commitCfgs, CCIPConfigTypes.OCR3Config[] memory execCfgs) =
      s_ccipCC.groupByPluginType(cfgs);

    assertEq(commitCfgs.length, numCommitCfgs, "commitCfgs length must match");
    assertEq(execCfgs.length, numExecCfgs, "execCfgs length must match");
    for (uint256 i = 0; i < commitCfgs.length; i++) {
      assertEq(uint8(commitCfgs[i].pluginType), uint8(Internal.OCRPluginType.Commit), "plugin type must be commit");
      assertEq(commitCfgs[i].offchainConfig, abi.encode("commit", i), "offchain config must match");
    }
    for (uint256 i = 0; i < execCfgs.length; i++) {
      assertEq(uint8(execCfgs[i].pluginType), uint8(Internal.OCRPluginType.Execution), "plugin type must be execution");
      assertEq(execCfgs[i].offchainConfig, abi.encode("exec", numCommitCfgs + i), "offchain config must match");
    }
  }

  function test__computeNewConfigWithMeta_InitToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](0);
    CCIPConfigTypes.OCR3Config[] memory newConfig = new CCIPConfigTypes.OCR3Config[](1);
    newConfig[0] = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.ConfigState currentState = CCIPConfigTypes.ConfigState.Init;
    CCIPConfigTypes.ConfigState newState = CCIPConfigTypes.ConfigState.Running;
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta =
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
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    CCIPConfigTypes.OCR3Config[] memory newConfig = new CCIPConfigTypes.OCR3Config[](2);
    // existing blue config first.
    newConfig[0] = blueConfig;
    // green config next.
    newConfig[1] = greenConfig;

    CCIPConfigTypes.ConfigState currentState = CCIPConfigTypes.ConfigState.Running;
    CCIPConfigTypes.ConfigState newState = CCIPConfigTypes.ConfigState.Staging;

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta =
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
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](2);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    currentConfig[1] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });
    CCIPConfigTypes.OCR3Config[] memory newConfig = new CCIPConfigTypes.OCR3Config[](1);
    newConfig[0] = greenConfig;

    CCIPConfigTypes.ConfigState currentState = CCIPConfigTypes.ConfigState.Staging;
    CCIPConfigTypes.ConfigState newState = CCIPConfigTypes.ConfigState.Running;

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta =
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
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](0);

    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_RunningToStaging_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](2);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    newConfig[1] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_StagingToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](2);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    currentConfig[1] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  // Reverts.

  function test_Fuzz__stateFromConfigLength_Reverts(uint256 configLen) public {
    vm.assume(configLen > 2);
    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.InvalidConfigLength.selector, configLen));
    s_ccipCC.stateFromConfigLength(configLen);
  }

  function test__groupByPluginType_threeCommitConfigs_Reverts() public {
    bytes32[] memory p2pIds = _makeBytes32Array(4, 0);
    bytes[] memory signers = _makeBytesArray(4, 10);
    bytes[] memory transmitters = _makeBytesArray(4, 20);
    CCIPConfigTypes.OCR3Config[] memory cfgs = new CCIPConfigTypes.OCR3Config[](3);
    for (uint256 i = 0; i < 3; i++) {
      cfgs[i] = CCIPConfigTypes.OCR3Config({
        pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory cfgs = new CCIPConfigTypes.OCR3Config[](3);
    for (uint256 i = 0; i < 3; i++) {
      cfgs[i] = CCIPConfigTypes.OCR3Config({
        pluginType: Internal.OCRPluginType.Execution,
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
    CCIPConfigTypes.OCR3Config[] memory cfgs = new CCIPConfigTypes.OCR3Config[](5);
    vm.expectRevert(CCIPConfig.TooManyOCR3Configs.selector);
    s_ccipCC.groupByPluginType(cfgs);
  }

  function test__validateConfigTransition_InitToRunning_WrongConfigCount_Reverts() public {
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 0,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](0);

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.WrongConfigCount.selector, 0, 1));
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_RunningToStaging_WrongConfigDigestBlueGreen_Reverts() public {
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](2);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 3, blueConfig) // wrong config digest (due to diff config count)
    });
    newConfig[1] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        CCIPConfig.WrongConfigDigestBlueGreen.selector,
        s_ccipCC.computeConfigDigest(donId, 3, blueConfig),
        s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
      )
    );
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_RunningToStaging_WrongConfigCount_Reverts() public {
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](2);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    newConfig[1] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 3, // wrong config count
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 3, greenConfig)
    });

    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.WrongConfigCount.selector, 3, 2));
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_StagingToRunning_WrongConfigDigest_Reverts() public {
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](2);
    currentConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 1,
      config: blueConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 1, blueConfig)
    });
    currentConfig[1] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
    });

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    newConfig[0] = CCIPConfigTypes.OCR3ConfigWithMeta({
      configCount: 2,
      config: greenConfig,
      configDigest: s_ccipCC.computeConfigDigest(donId, 3, greenConfig) // wrong config digest
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        CCIPConfig.WrongConfigDigest.selector,
        s_ccipCC.computeConfigDigest(donId, 3, greenConfig),
        s_ccipCC.computeConfigDigest(donId, 2, greenConfig)
      )
    );
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }

  function test__validateConfigTransition_NonExistentConfigTransition_Reverts() public {
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](3);
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfig = new CCIPConfigTypes.OCR3ConfigWithMeta[](1);
    vm.expectRevert(CCIPConfig.NonExistentConfigTransition.selector);
    s_ccipCC.validateConfigTransition(currentConfig, newConfig);
  }
}

contract CCIPConfig_updatePluginConfig is CCIPConfigSetup {
  // Successes.

  function test__updatePluginConfig_InitToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory configs = new CCIPConfigTypes.OCR3Config[](1);
    configs[0] = blueConfig;

    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, configs);

    // should see the updated config in the contract state.
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedConfig =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Commit);
    assertEq(storedConfig.length, 1, "don config length must be 1");
    assertEq(storedConfig[0].configCount, uint64(1), "config count must be 1");
    assertEq(uint256(storedConfig[0].config.pluginType), uint256(blueConfig.pluginType), "plugin type must match");
  }

  function test__updatePluginConfig_RunningToStaging_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // add blue config.
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory startConfigs = new CCIPConfigTypes.OCR3Config[](1);
    startConfigs[0] = blueConfig;

    // add blue AND green config to indicate an update.
    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, startConfigs);
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory blueAndGreen = new CCIPConfigTypes.OCR3Config[](2);
    blueAndGreen[0] = blueConfig;
    blueAndGreen[1] = greenConfig;

    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, blueAndGreen);

    // should see the updated config in the contract state.
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedConfig =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Commit);
    assertEq(storedConfig.length, 2, "don config length must be 2");
    // 0 index is blue config, 1 index is green config.
    assertEq(storedConfig[1].configCount, uint64(2), "config count must be 2");
    assertEq(
      uint256(storedConfig[0].config.pluginType), uint256(Internal.OCRPluginType.Commit), "plugin type must match"
    );
    assertEq(
      uint256(storedConfig[1].config.pluginType), uint256(Internal.OCRPluginType.Commit), "plugin type must match"
    );
    assertEq(storedConfig[0].config.offchainConfig, bytes("commit"), "blue offchain config must match");
    assertEq(storedConfig[1].config.offchainConfig, bytes("commit-new"), "green offchain config must match");
  }

  function test__updatePluginConfig_StagingToRunning_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    // add blue config.
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory startConfigs = new CCIPConfigTypes.OCR3Config[](1);
    startConfigs[0] = blueConfig;

    // add blue AND green config to indicate an update.
    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, startConfigs);
    CCIPConfigTypes.OCR3Config memory greenConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory blueAndGreen = new CCIPConfigTypes.OCR3Config[](2);
    blueAndGreen[0] = blueConfig;
    blueAndGreen[1] = greenConfig;

    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, blueAndGreen);

    // should see the updated config in the contract state.
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedConfig =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Commit);
    assertEq(storedConfig.length, 2, "don config length must be 2");
    // 0 index is blue config, 1 index is green config.
    assertEq(storedConfig[1].configCount, uint64(2), "config count must be 2");
    assertEq(
      uint256(storedConfig[0].config.pluginType), uint256(Internal.OCRPluginType.Commit), "plugin type must match"
    );
    assertEq(
      uint256(storedConfig[1].config.pluginType), uint256(Internal.OCRPluginType.Commit), "plugin type must match"
    );
    assertEq(storedConfig[0].config.offchainConfig, bytes("commit"), "blue offchain config must match");
    assertEq(storedConfig[1].config.offchainConfig, bytes("commit-new"), "green offchain config must match");

    // promote green to blue.
    CCIPConfigTypes.OCR3Config[] memory promote = new CCIPConfigTypes.OCR3Config[](1);
    promote[0] = greenConfig;

    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, promote);

    // should see the updated config in the contract state.
    storedConfig = s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Commit);
    assertEq(storedConfig.length, 1, "don config length must be 1");
    assertEq(storedConfig[0].configCount, uint64(2), "config count must be 2");
    assertEq(
      uint256(storedConfig[0].config.pluginType), uint256(Internal.OCRPluginType.Commit), "plugin type must match"
    );
    assertEq(storedConfig[0].config.offchainConfig, bytes("commit-new"), "green offchain config must match");
  }

  // Reverts.
  function test__updatePluginConfig_InvalidConfigLength_Reverts() public {
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config[] memory newConfig = new CCIPConfigTypes.OCR3Config[](3);
    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.InvalidConfigLength.selector, uint256(3)));
    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, newConfig);
  }

  function test__updatePluginConfig_InvalidConfigStateTransition_Reverts() public {
    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config[] memory newConfig = new CCIPConfigTypes.OCR3Config[](2);
    // 0 -> 2 is an invalid state transition.
    vm.expectRevert(abi.encodeWithSelector(CCIPConfig.InvalidConfigStateTransition.selector, 0, 2));
    s_ccipCC.updatePluginConfig(donId, Internal.OCRPluginType.Commit, newConfig);
  }
}

contract CCIPConfig_beforeCapabilityConfigSet is CCIPConfigSetup {
  // Successes.
  function test_beforeCapabilityConfigSet_ZeroLengthConfig_Success() public {
    changePrank(CAPABILITIES_REGISTRY);

    CCIPConfigTypes.OCR3Config[] memory configs = new CCIPConfigTypes.OCR3Config[](0);
    bytes memory encodedConfigs = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encodedConfigs, 1, 1);
  }

  function test_beforeCapabilityConfigSet_CommitConfigOnly_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    changePrank(CAPABILITIES_REGISTRY);

    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config[] memory configs = new CCIPConfigTypes.OCR3Config[](1);
    configs[0] = blueConfig;

    bytes memory encoded = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encoded, 1, donId);

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedConfigs =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Commit);
    assertEq(storedConfigs.length, 1, "config length must be 1");
    assertEq(storedConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedConfigs[0].config.pluginType), uint256(Internal.OCRPluginType.Commit), "plugin type must be commit"
    );
  }

  function test_beforeCapabilityConfigSet_ExecConfigOnly_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    changePrank(CAPABILITIES_REGISTRY);

    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Execution,
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
    CCIPConfigTypes.OCR3Config[] memory configs = new CCIPConfigTypes.OCR3Config[](1);
    configs[0] = blueConfig;

    bytes memory encoded = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encoded, 1, donId);

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedConfigs =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Execution);
    assertEq(storedConfigs.length, 1, "config length must be 1");
    assertEq(storedConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedConfigs[0].config.pluginType),
      uint256(Internal.OCRPluginType.Execution),
      "plugin type must be execution"
    );
  }

  function test_beforeCapabilityConfigSet_CommitAndExecConfig_Success() public {
    (bytes32[] memory p2pIds, bytes[] memory signers, bytes[] memory transmitters) = _addChainConfig(4);
    changePrank(CAPABILITIES_REGISTRY);

    uint32 donId = 1;
    CCIPConfigTypes.OCR3Config memory blueCommitConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Commit,
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
    CCIPConfigTypes.OCR3Config memory blueExecConfig = CCIPConfigTypes.OCR3Config({
      pluginType: Internal.OCRPluginType.Execution,
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
    CCIPConfigTypes.OCR3Config[] memory configs = new CCIPConfigTypes.OCR3Config[](2);
    configs[0] = blueExecConfig;
    configs[1] = blueCommitConfig;

    bytes memory encoded = abi.encode(configs);
    s_ccipCC.beforeCapabilityConfigSet(new bytes32[](0), encoded, 1, donId);

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedExecConfigs =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Execution);
    assertEq(storedExecConfigs.length, 1, "config length must be 1");
    assertEq(storedExecConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedExecConfigs[0].config.pluginType),
      uint256(Internal.OCRPluginType.Execution),
      "plugin type must be execution"
    );

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory storedCommitConfigs =
      s_ccipCC.getOCRConfig(donId, Internal.OCRPluginType.Commit);
    assertEq(storedCommitConfigs.length, 1, "config length must be 1");
    assertEq(storedCommitConfigs[0].configCount, uint64(1), "config count must be 1");
    assertEq(
      uint256(storedCommitConfigs[0].config.pluginType),
      uint256(Internal.OCRPluginType.Commit),
      "plugin type must be commit"
    );
  }

  // Reverts.

  function test_beforeCapabilityConfigSet_OnlyCapabilitiesRegistryCanCall_Reverts() public {
    bytes32[] memory nodes = new bytes32[](0);
    bytes memory config = bytes("");
    uint64 configCount = 1;
    uint32 donId = 1;
    vm.expectRevert(CCIPConfig.OnlyCapabilitiesRegistryCanCall.selector);
    s_ccipCC.beforeCapabilityConfigSet(nodes, config, configCount, donId);
  }
}
