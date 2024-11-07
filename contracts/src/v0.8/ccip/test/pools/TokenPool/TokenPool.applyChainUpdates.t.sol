// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_applyChainUpdates is TokenPoolSetup {
  function assertState(
    TokenPool.ChainUpdate[] memory chainUpdates
  ) public view {
    uint64[] memory chainSelectors = s_tokenPool.getSupportedChains();
    for (uint256 i = 0; i < chainUpdates.length; i++) {
      assertEq(chainUpdates[i].remoteChainSelector, chainSelectors[i]);
    }

    for (uint256 i = 0; i < chainUpdates.length; ++i) {
      assertTrue(s_tokenPool.isSupportedChain(chainUpdates[i].remoteChainSelector));
      RateLimiter.TokenBucket memory bkt =
        s_tokenPool.getCurrentOutboundRateLimiterState(chainUpdates[i].remoteChainSelector);
      assertEq(bkt.capacity, chainUpdates[i].outboundRateLimiterConfig.capacity);
      assertEq(bkt.rate, chainUpdates[i].outboundRateLimiterConfig.rate);
      assertEq(bkt.isEnabled, chainUpdates[i].outboundRateLimiterConfig.isEnabled);

      bkt = s_tokenPool.getCurrentInboundRateLimiterState(chainUpdates[i].remoteChainSelector);
      assertEq(bkt.capacity, chainUpdates[i].inboundRateLimiterConfig.capacity);
      assertEq(bkt.rate, chainUpdates[i].inboundRateLimiterConfig.rate);
      assertEq(bkt.isEnabled, chainUpdates[i].inboundRateLimiterConfig.isEnabled);
    }
  }

  function test_applyChainUpdates_Success() public {
    RateLimiter.Config memory outboundRateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18});
    RateLimiter.Config memory inboundRateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e29, rate: 1e19});
    RateLimiter.Config memory outboundRateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 100e26, rate: 1e16});
    RateLimiter.Config memory inboundRateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 100e27, rate: 1e17});

    // EVM chain, which uses the 160 bit evm address space
    uint64 evmChainSelector = 1;
    bytes memory evmRemotePool = abi.encode(makeAddr("evm_remote_pool"));
    bytes memory evmRemoteToken = abi.encode(makeAddr("evm_remote_token"));

    // Non EVM chain, which uses the full 256 bits
    uint64 nonEvmChainSelector = type(uint64).max;
    bytes memory nonEvmRemotePool = abi.encode(keccak256("non_evm_remote_pool"));
    bytes memory nonEvmRemoteToken = abi.encode(keccak256("non_evm_remote_token"));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](2);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: evmChainSelector,
      remotePoolAddress: evmRemotePool,
      remoteTokenAddress: evmRemoteToken,
      allowed: true,
      outboundRateLimiterConfig: outboundRateLimit1,
      inboundRateLimiterConfig: inboundRateLimit1
    });
    chainUpdates[1] = TokenPool.ChainUpdate({
      remoteChainSelector: nonEvmChainSelector,
      remotePoolAddress: nonEvmRemotePool,
      remoteTokenAddress: nonEvmRemoteToken,
      allowed: true,
      outboundRateLimiterConfig: outboundRateLimit2,
      inboundRateLimiterConfig: inboundRateLimit2
    });

    // Assert configuration is applied
    vm.expectEmit();
    emit TokenPool.ChainAdded(
      chainUpdates[0].remoteChainSelector,
      chainUpdates[0].remoteTokenAddress,
      chainUpdates[0].outboundRateLimiterConfig,
      chainUpdates[0].inboundRateLimiterConfig
    );
    vm.expectEmit();
    emit TokenPool.ChainAdded(
      chainUpdates[1].remoteChainSelector,
      chainUpdates[1].remoteTokenAddress,
      chainUpdates[1].outboundRateLimiterConfig,
      chainUpdates[1].inboundRateLimiterConfig
    );
    s_tokenPool.applyChainUpdates(chainUpdates);
    // on1: rateLimit1, on2: rateLimit2, off1: rateLimit1, off2: rateLimit3
    assertState(chainUpdates);

    // Removing an non-existent chain should revert
    TokenPool.ChainUpdate[] memory chainRemoves = new TokenPool.ChainUpdate[](1);
    uint64 strangerChainSelector = 120938;
    chainRemoves[0] = TokenPool.ChainUpdate({
      remoteChainSelector: strangerChainSelector,
      remotePoolAddress: evmRemotePool,
      remoteTokenAddress: evmRemoteToken,
      allowed: false,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0})
    });
    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, strangerChainSelector));
    s_tokenPool.applyChainUpdates(chainRemoves);
    // State remains
    assertState(chainUpdates);

    // Can remove a chain
    chainRemoves[0].remoteChainSelector = evmChainSelector;

    vm.expectEmit();
    emit TokenPool.ChainRemoved(chainRemoves[0].remoteChainSelector);

    s_tokenPool.applyChainUpdates(chainRemoves);

    // State updated, only chain 2 remains
    TokenPool.ChainUpdate[] memory singleChainConfigured = new TokenPool.ChainUpdate[](1);
    singleChainConfigured[0] = chainUpdates[1];
    assertState(singleChainConfigured);

    // Cannot reset already configured ramp
    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.ChainAlreadyExists.selector, singleChainConfigured[0].remoteChainSelector)
    );
    s_tokenPool.applyChainUpdates(singleChainConfigured);
  }

  // Reverts

  function test_applyChainUpdates_OnlyCallableByOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_tokenPool.applyChainUpdates(new TokenPool.ChainUpdate[](0));
  }

  function test_applyChainUpdates_ZeroAddressNotAllowed_Revert() public {
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddress: "",
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18})
    });

    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddress: abi.encode(address(2)),
      remoteTokenAddress: "",
      allowed: true,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18})
    });

    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);
    s_tokenPool.applyChainUpdates(chainUpdates);
  }

  function test_applyChainUpdates_DisabledNonZeroRateLimit_Revert() public {
    RateLimiter.Config memory outboundRateLimit = RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18});
    RateLimiter.Config memory inboundRateLimit = RateLimiter.Config({isEnabled: true, capacity: 100e22, rate: 1e12});
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: outboundRateLimit,
      inboundRateLimiterConfig: inboundRateLimit
    });

    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].allowed = false;
    chainUpdates[0].outboundRateLimiterConfig = RateLimiter.Config({isEnabled: false, capacity: 10, rate: 1});
    chainUpdates[0].inboundRateLimiterConfig = RateLimiter.Config({isEnabled: false, capacity: 10, rate: 1});

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.DisabledNonZeroRateLimit.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);
  }

  function test_applyChainUpdates_NonExistentChain_Revert() public {
    RateLimiter.Config memory outboundRateLimit = RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0});
    RateLimiter.Config memory inboundRateLimit = RateLimiter.Config({isEnabled: false, capacity: 0, rate: 0});
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: false,
      outboundRateLimiterConfig: outboundRateLimit,
      inboundRateLimiterConfig: inboundRateLimit
    });

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainUpdates[0].remoteChainSelector));
    s_tokenPool.applyChainUpdates(chainUpdates);
  }

  function test_applyChainUpdates_InvalidRateLimitRate_Revert() public {
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e22, rate: 1e12})
    });

    // Outbound

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.rate = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.capacity = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.capacity = 101;

    s_tokenPool.applyChainUpdates(chainUpdates);

    // Change the chain selector as adding the same one would revert
    chainUpdates[0].remoteChainSelector = 2;

    // Inbound

    chainUpdates[0].inboundRateLimiterConfig.capacity = 0;
    chainUpdates[0].inboundRateLimiterConfig.rate = 0;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.rate = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.capacity = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.capacity = 101;

    s_tokenPool.applyChainUpdates(chainUpdates);
  }
}
