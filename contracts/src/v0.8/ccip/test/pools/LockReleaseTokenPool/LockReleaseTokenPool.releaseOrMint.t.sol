// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";

import {TokenPool} from "../../../pools/TokenPool.sol";
import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

contract LockReleaseTokenPool_releaseOrMint is LockReleaseTokenPoolSetup {
  function setUp() public virtual override {
    LockReleaseTokenPoolSetup.setUp();

    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = abi.encode(s_sourcePoolAddress);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: abi.encode(address(2)),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_lockReleaseTokenPool.applyChainUpdates(new uint64[](0), chainUpdate);
    s_lockReleaseTokenPoolWithAllowList.applyChainUpdates(new uint64[](0), chainUpdate);
  }

  function test_ReleaseOrMint() public {
    vm.startPrank(s_allowedOffRamp);

    uint256 amount = 100;
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit TokenPool.Released(s_allowedOffRamp, OWNER, amount);

    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_sourcePoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function testFuzz_ReleaseOrMint_Success(address recipient, uint256 amount) public {
    // Since the owner already has tokens this would break the checks
    vm.assume(recipient != OWNER);
    vm.assume(recipient != address(0));
    vm.assume(recipient != address(s_token));

    // Makes sure the pool always has enough funds
    deal(address(s_token), address(s_lockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOffRamp);

    uint256 capacity = _getInboundRateLimiterConfig().capacity;
    // Determine if we hit the rate limit or the txs should succeed.
    if (amount > capacity) {
      vm.expectRevert(
        abi.encodeWithSelector(RateLimiter.TokenMaxCapacityExceeded.selector, capacity, amount, address(s_token))
      );
    } else {
      // Only rate limit if the amount is >0
      if (amount > 0) {
        vm.expectEmit();
        emit RateLimiter.TokensConsumed(amount);
      }

      vm.expectEmit();
      emit TokenPool.Released(s_allowedOffRamp, recipient, amount);
    }

    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_sourcePoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_RevertWhen_ChainNotAllowed() public {
    uint64[] memory chainsToRemove = new uint64[](1);
    chainsToRemove[0] = SOURCE_CHAIN_SELECTOR;

    s_lockReleaseTokenPool.applyChainUpdates(chainsToRemove, new TokenPool.ChainUpdate[](0));

    vm.startPrank(s_allowedOffRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, SOURCE_CHAIN_SELECTOR));
    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: 1e5,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_sourcePoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_RevertWhen_PoolMintNotHealthy() public {
    // Should not mint tokens if cursed.
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed(bytes16)"), abi.encode(true));
    uint256 before = s_token.balanceOf(OWNER);
    vm.startPrank(s_allowedOffRamp);
    vm.expectRevert(TokenPool.CursedByRMN.selector);
    s_lockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: 1e5,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: _generateSourceTokenData().sourcePoolAddress,
        sourcePoolData: _generateSourceTokenData().extraData,
        offchainTokenData: ""
      })
    );

    assertEq(s_token.balanceOf(OWNER), before);
  }
}
