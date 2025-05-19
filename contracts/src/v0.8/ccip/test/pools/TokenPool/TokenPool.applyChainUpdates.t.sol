// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {BurnMintERC677} from "../../../../shared/token/ERC677/BurnMintERC677.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPool_applyChainUpdates is BaseTest {
  IERC20 internal s_token;
  TokenPoolHelper internal s_tokenPool;

  function setUp() public virtual override {
    super.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);

    s_tokenPool = new TokenPoolHelper(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );
  }

  function assertState(
    TokenPool.ChainUpdate[] memory chainUpdates
  ) public view {
    uint64[] memory chainSelectors = s_tokenPool.getSupportedChains();
    for (uint256 i = 0; i < chainUpdates.length; ++i) {
      assertEq(chainUpdates[i].remoteChainSelector, chainSelectors[i], "Chain selector mismatch");
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

  function test_applyChainUpdates() public {
    RateLimiter.Config memory outboundRateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18});
    RateLimiter.Config memory inboundRateLimit1 = RateLimiter.Config({isEnabled: true, capacity: 100e29, rate: 1e19});
    RateLimiter.Config memory outboundRateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 100e26, rate: 1e16});
    RateLimiter.Config memory inboundRateLimit2 = RateLimiter.Config({isEnabled: true, capacity: 100e27, rate: 1e17});

    // EVM chain, which uses the 160 bit evm address space
    uint64 evmChainSelector = 1789142;
    bytes memory evmRemotePool = abi.encode(makeAddr("evm_remote_pool"));
    bytes memory evmRemoteToken = abi.encode(makeAddr("evm_remote_token"));

    // Non EVM chain, which uses the full 256 bits
    uint64 nonEvmChainSelector = type(uint64).max;
    bytes memory nonEvmRemotePool = abi.encode(keccak256("non_evm_remote_pool"));
    bytes memory nonEvmRemoteToken = abi.encode(keccak256("non_evm_remote_token"));

    bytes[] memory evmRemotePools = new bytes[](1);
    evmRemotePools[0] = evmRemotePool;

    bytes[] memory nonEvmRemotePools = new bytes[](1);
    nonEvmRemotePools[0] = nonEvmRemotePool;

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](2);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: evmChainSelector,
      remotePoolAddresses: evmRemotePools,
      remoteTokenAddress: evmRemoteToken,
      outboundRateLimiterConfig: outboundRateLimit1,
      inboundRateLimiterConfig: inboundRateLimit1
    });
    chainUpdates[1] = TokenPool.ChainUpdate({
      remoteChainSelector: nonEvmChainSelector,
      remotePoolAddresses: nonEvmRemotePools,
      remoteTokenAddress: nonEvmRemoteToken,
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
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);
    // on1: rateLimit1, on2: rateLimit2, off1: rateLimit1, off2: rateLimit3
    assertState(chainUpdates);

    // Removing an non-existent chain should revert
    uint64 strangerChainSelector = 120938;

    uint64[] memory chainRemoves = new uint64[](1);
    chainRemoves[0] = strangerChainSelector;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, strangerChainSelector));
    s_tokenPool.applyChainUpdates(chainRemoves, new TokenPool.ChainUpdate[](0));
    // State remains
    assertState(chainUpdates);

    // Can remove a chain
    chainRemoves[0] = evmChainSelector;

    vm.expectEmit();
    emit TokenPool.ChainRemoved(chainRemoves[0]);

    s_tokenPool.applyChainUpdates(chainRemoves, new TokenPool.ChainUpdate[](0));

    // State updated, only chain 2 remains
    TokenPool.ChainUpdate[] memory singleChainConfigured = new TokenPool.ChainUpdate[](1);
    singleChainConfigured[0] = chainUpdates[1];
    assertState(singleChainConfigured);

    // Cannot reset already configured ramp
    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.ChainAlreadyExists.selector, singleChainConfigured[0].remoteChainSelector)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), singleChainConfigured);
  }

  function test_applyChainUpdates_UpdatesRemotePoolHashes() public {
    assertEq(s_tokenPool.getRemotePools(DEST_CHAIN_SELECTOR).length, 0);

    uint64 selector1 = 789;
    uint64 selector2 = 123;
    uint64 selector3 = 456;

    bytes memory pool1 = abi.encode(makeAddr("pool1"));
    bytes memory pool2 = abi.encode(makeAddr("pool2"));
    bytes memory pool3 = abi.encode(makeAddr("pool3"));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](3);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: selector1,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: pool1,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    chainUpdates[1] = TokenPool.ChainUpdate({
      remoteChainSelector: selector2,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: pool2,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    chainUpdates[2] = TokenPool.ChainUpdate({
      remoteChainSelector: selector3,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: pool3,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    // This adds 3 for the first chain, 2 for the second, and 1 for the third for a total of 6.
    for (uint256 i = 0; i < chainUpdates.length; ++i) {
      for (uint256 j = i; j < chainUpdates.length; ++j) {
        s_tokenPool.addRemotePool(chainUpdates[i].remoteChainSelector, abi.encode(i, j));
      }
      assertEq(s_tokenPool.getRemotePools(chainUpdates[i].remoteChainSelector).length, 3 - i);
    }

    // Removing a chain should remove all associated pool hashes
    uint64[] memory chainRemoves = new uint64[](1);
    chainRemoves[0] = selector1;

    s_tokenPool.applyChainUpdates(chainRemoves, new TokenPool.ChainUpdate[](0));

    assertEq(s_tokenPool.getRemotePools(selector1).length, 0);

    chainRemoves[0] = selector2;

    s_tokenPool.applyChainUpdates(chainRemoves, new TokenPool.ChainUpdate[](0));

    assertEq(s_tokenPool.getRemotePools(selector2).length, 0);

    // The above deletions should not have affected the third chain
    assertEq(s_tokenPool.getRemotePools(selector3).length, 1);
  }

  // Reverts

  function test_RevertWhen_applyChainUpdates_OnlyCallableByOwner() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_tokenPool.applyChainUpdates(new uint64[](0), new TokenPool.ChainUpdate[](0));
  }

  function test_RevertWhen_applyChainUpdates_ZeroAddressNotAllowed() public {
    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = "";

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: abi.encode(address(2)),
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18})
    });

    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: 1,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: "",
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e18})
    });

    vm.expectRevert(TokenPool.ZeroAddressNotAllowed.selector);
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);
  }

  function test_RevertWhen_applyChainUpdates_NonExistentChain() public {
    uint64[] memory chainRemoves = new uint64[](1);
    chainRemoves[0] = 1;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.NonExistentChain.selector, chainRemoves[0]));
    s_tokenPool.applyChainUpdates(chainRemoves, new TokenPool.ChainUpdate[](0));
  }

  function test_RevertWhen_applyChainUpdates_InvalidRateLimitRate() public {
    uint64 unusedChainSelector = 2 ** 64 - 1;

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: unusedChainSelector,
      remotePoolAddresses: new bytes[](0),
      remoteTokenAddress: abi.encode(address(2)),
      outboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 0, rate: 0}),
      inboundRateLimiterConfig: RateLimiter.Config({isEnabled: true, capacity: 100e22, rate: 1e12})
    });

    // Outbound

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.rate = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.capacity = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].outboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates[0].outboundRateLimiterConfig.capacity = 101;

    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    // Change the chain selector as adding the same one would revert
    chainUpdates[0].remoteChainSelector = unusedChainSelector - 1;

    // Inbound

    chainUpdates[0].inboundRateLimiterConfig.capacity = 0;
    chainUpdates[0].inboundRateLimiterConfig.rate = 0;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.rate = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.capacity = 100;

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.InvalidRateLimitRate.selector, chainUpdates[0].inboundRateLimiterConfig)
    );
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);

    chainUpdates[0].inboundRateLimiterConfig.capacity = 101;

    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdates);
  }
}
