// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RMNHome} from "../../../rmn/RMNHome.sol";

import {RMNHomeTestSetup} from "./RMNHomeTestSetup.t.sol";

contract RMNHome_setDynamicConfig is RMNHomeTestSetup {
  function setUp() public {
    Config memory config = _getBaseConfig();
    s_rmnHome.setCandidate(config.staticConfig, config.dynamicConfig, ZERO_DIGEST);
  }

  function test_setDynamicConfig() public {
    (bytes32 priorActiveDigest,) = s_rmnHome.getConfigDigests();

    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[1].fObserve--;

    (, bytes32 candidateConfigDigest) = s_rmnHome.getConfigDigests();

    vm.expectEmit();
    emit RMNHome.DynamicConfigSet(candidateConfigDigest, config.dynamicConfig);

    s_rmnHome.setDynamicConfig(config.dynamicConfig, candidateConfigDigest);

    (RMNHome.VersionedConfig memory storedVersionedConfig, bool ok) = s_rmnHome.getConfig(candidateConfigDigest);
    assertTrue(ok);
    assertEq(
      storedVersionedConfig.dynamicConfig.sourceChains[0].fObserve, config.dynamicConfig.sourceChains[0].fObserve
    );

    // Asser the digests don't change when updating the dynamic config
    (bytes32 activeDigest, bytes32 candidateDigest) = s_rmnHome.getConfigDigests();
    assertEq(activeDigest, priorActiveDigest);
    assertEq(candidateDigest, candidateConfigDigest);
  }

  // Asserts the validation function is being called
  function test_RevertWhen_setDynamicConfig_MinObserversTooHigh() public {
    Config memory config = _getBaseConfig();
    config.dynamicConfig.sourceChains[0].fObserve++;

    vm.expectRevert(abi.encodeWithSelector(RMNHome.DigestNotFound.selector, ZERO_DIGEST));
    s_rmnHome.setDynamicConfig(config.dynamicConfig, ZERO_DIGEST);
  }

  function test_RevertWhen_setDynamicConfig_DigestNotFound() public {
    // Zero always reverts
    vm.expectRevert(abi.encodeWithSelector(RMNHome.DigestNotFound.selector, ZERO_DIGEST));
    s_rmnHome.setDynamicConfig(_getBaseConfig().dynamicConfig, ZERO_DIGEST);

    // Non-existent digest reverts
    bytes32 nonExistentDigest = keccak256("nonExistentDigest");
    vm.expectRevert(abi.encodeWithSelector(RMNHome.DigestNotFound.selector, nonExistentDigest));
    s_rmnHome.setDynamicConfig(_getBaseConfig().dynamicConfig, nonExistentDigest);
  }

  function test_RevertWhen_setDynamicConfig_OnlyOwner() public {
    Config memory config = _getBaseConfig();

    vm.startPrank(address(0));

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_rmnHome.setDynamicConfig(config.dynamicConfig, keccak256("configDigest"));
  }
}
