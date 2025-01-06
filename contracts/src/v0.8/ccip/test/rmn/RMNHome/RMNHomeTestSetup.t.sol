// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNHome} from "../../../rmn/RMNHome.sol";
import {Test} from "forge-std/Test.sol";

contract RMNHomeTestSetup is Test {
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
    sourceChains[0] =
      RMNHome.SourceChain({chainSelector: 9000, fObserve: 1, observerNodesBitmap: 1 << 0 | 1 << 1 | 1 << 2});
    // Observers 0, 1 and 2 for source chain 9001
    sourceChains[1] =
      RMNHome.SourceChain({chainSelector: 9001, fObserve: 1, observerNodesBitmap: 1 << 0 | 1 << 1 | 1 << 2});

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
