// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNHome} from "../../../rmn/RMNHome.sol";

import {RMNHomeTestSetup} from "./RMNHomeTestSetup.t.sol";

contract RMNHome_validateStaticAndDynamicConfig is RMNHomeTestSetup {
  function test_RevertWhen_validateStaticAndDynamicConfig_OutOfBoundsNodesLength() public {
    Config memory config = _getBaseConfig();
    config.staticConfig.nodes = new RMNHome.Node[](257);

    vm.expectRevert(RMNHome.OutOfBoundsNodesLength.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_RevertWhen_validateStaticAndDynamicConfig_DuplicatePeerId() public {
    Config memory config = _getBaseConfig();
    config.staticConfig.nodes[1].peerId = config.staticConfig.nodes[0].peerId;

    vm.expectRevert(RMNHome.DuplicatePeerId.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_RevertWhen_validateStaticAndDynamicConfig_DuplicateOffchainPublicKey() public {
    Config memory config = _getBaseConfig();
    config.staticConfig.nodes[1].offchainPublicKey = config.staticConfig.nodes[0].offchainPublicKey;

    vm.expectRevert(RMNHome.DuplicateOffchainPublicKey.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_RevertWhen_validateStaticAndDynamicConfig_DuplicateSourceChain() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[1].chainSelector = config.dynamicConfig.sourceChains[0].chainSelector;

    vm.expectRevert(RMNHome.DuplicateSourceChain.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_RevertWhen_validateStaticAndDynamicConfig_OutOfBoundsObserverNodeIndex() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].observerNodesBitmap = 1 << config.staticConfig.nodes.length;

    vm.expectRevert(RMNHome.OutOfBoundsObserverNodeIndex.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_RevertWhen_validateStaticAndDynamicConfig_NotEnoughObservers() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].fObserve++;

    vm.expectRevert(RMNHome.NotEnoughObservers.selector);
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }
}
