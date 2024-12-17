// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {INodeInfoProvider} from "../../../../keystone/interfaces/INodeInfoProvider.sol";

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {CCIPHomeHelper} from "../../helpers/CCIPHomeHelper.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome__validateConfig is CCIPHomeTestSetup {
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
    INodeInfoProvider.NodeInfo[] memory nodeInfos = new INodeInfoProvider.NodeInfo[](numNodes);
    for (uint256 i = 0; i < numNodes; i++) {
      nodes[i] = CCIPHome.OCR3Node({p2pId: p2pIds[i], signerKey: signers[i], transmitterKey: transmitters[i]});
      nodeInfos[i] = INodeInfoProvider.NodeInfo({
        nodeOperatorId: 1,
        signer: bytes32(signers[i]),
        p2pId: p2pIds[i],
        encryptionPublicKey: keccak256("encryptionPublicKey"),
        hashedCapabilityIds: new bytes32[](0),
        configCount: uint32(1),
        workflowDONId: uint32(1),
        capabilitiesDONIds: new uint256[](0)
      });
    }
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, p2pIds),
      abi.encode(nodeInfos)
    );
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

  function test__validateConfig() public {
    s_ccipHome.validateConfig(_getCorrectOCR3Config());
  }

  function test__validateConfigLessTransmittersThanSigners() public {
    // fChain is 1, so there should be at least 4 transmitters.
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config(5, 1);
    config.nodes[1].transmitterKey = bytes("");

    s_ccipHome.validateConfig(config);
  }

  function test__validateConfigSmallerFChain() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config(11, 3);

    // Set fChain to 2
    _addChainConfig(4, 2);

    s_ccipHome.validateConfig(config);
  }

  // Reverts

  function test_RevertWhen__validateConfig_ChainSelectorNotSet() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.chainSelector = 0; // invalid

    vm.expectRevert(CCIPHome.ChainSelectorNotSet.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_OfframpAddressCannotBeZero() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.offrampAddress = ""; // invalid

    vm.expectRevert(CCIPHome.OfframpAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_ABIEncodedAddress_OfframpAddressCannotBeZero() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.offrampAddress = abi.encode(address(0)); // invalid

    vm.expectRevert(CCIPHome.OfframpAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_RMNHomeAddressCannotBeZero() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.rmnHomeAddress = ""; // invalid

    vm.expectRevert(CCIPHome.RMNHomeAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_ABIEncodedAddress_RMNHomeAddressCannotBeZero() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.rmnHomeAddress = abi.encode(address(0)); // invalid

    vm.expectRevert(CCIPHome.RMNHomeAddressCannotBeZero.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_ChainSelectorNotFound() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.chainSelector = 2; // not set

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ChainSelectorNotFound.selector, 2));
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_NotEnoughTransmitters() public {
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

  function test_RevertWhen__validateConfig_NotEnoughTransmittersEmptyAddresses() public {
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

  function test_RevertWhen__validateConfig_TooManySigners() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes = new CCIPHome.OCR3Node[](257);

    vm.expectRevert(CCIPHome.TooManySigners.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_FChainTooHigh() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.FRoleDON = 2; // too low

    // Set fChain to 3
    _addChainConfig(4, 3);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.FChainTooHigh.selector, 3, 2));
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_FMustBePositive() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.FRoleDON = 0; // not positive

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.FChainTooHigh.selector, 1, 0));
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_FTooHigh() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.FRoleDON = 2; // too high

    vm.expectRevert(CCIPHome.FTooHigh.selector);
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_ZeroP2PId() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes[1].p2pId = bytes32(0);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.InvalidNode.selector, config.nodes[1]));
    s_ccipHome.validateConfig(config);
  }

  function test_RevertWhen__validateConfig_ZeroSignerKey() public {
    CCIPHome.OCR3Config memory config = _getCorrectOCR3Config();
    config.nodes[2].signerKey = bytes("");

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.InvalidNode.selector, config.nodes[2]));
    s_ccipHome.validateConfig(config);
  }
}
