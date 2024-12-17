// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {INodeInfoProvider} from "../../../../keystone/interfaces/INodeInfoProvider.sol";

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {CCIPHomeHelper} from "../../helpers/CCIPHomeHelper.sol";
import {Test} from "forge-std/Test.sol";

contract CCIPHomeTestSetup is Test {
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
  ) internal returns (CCIPHome.OCR3Config memory) {
    CCIPHome.OCR3Node[] memory nodes = new CCIPHome.OCR3Node[](4);
    bytes32[] memory p2pIds = new bytes32[](4);
    for (uint256 i = 0; i < nodes.length; i++) {
      p2pIds[i] = keccak256(abi.encode("p2pId", i));
      nodes[i] = CCIPHome.OCR3Node({
        p2pId: keccak256(abi.encode("p2pId", i)),
        signerKey: abi.encode("signerKey"),
        transmitterKey: abi.encode("transmitterKey")
      });
    }

    // This is a work-around for not calling mockCall / expectCall with each scenario using _getBaseConfig
    INodeInfoProvider.NodeInfo[] memory nodeInfos = new INodeInfoProvider.NodeInfo[](4);
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, p2pIds),
      abi.encode(nodeInfos)
    );

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
