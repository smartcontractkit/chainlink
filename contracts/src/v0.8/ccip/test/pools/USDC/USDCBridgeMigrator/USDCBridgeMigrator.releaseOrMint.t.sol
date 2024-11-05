// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../../../shared/token/ERC20/IBurnMintERC20.sol";

import {BurnMintERC677} from "../../../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../../../Router.sol";
import {Internal} from "../../../../libraries/Internal.sol";
import {Pool} from "../../../../libraries/Pool.sol";

import {TokenPool} from "../../../../pools/TokenPool.sol";
import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {LOCK_RELEASE_FLAG} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {BaseTest} from "../../../BaseTest.t.sol";
import {MockE2EUSDCTransmitter} from "../../../mocks/MockE2EUSDCTransmitter.sol";
import {MockUSDCTokenMessenger} from "../../../mocks/MockUSDCTokenMessenger.sol";
import {HybridLockReleaseUSDCTokenPool_releaseOrMint} from
  "../HybridLockReleaseUSDCTokenPool/HybridLockReleaseUSDCTokenPool.releaseOrMint.t.sol";

contract USDCTokenPoolSetup is BaseTest {
  IBurnMintERC20 internal s_token;
  MockUSDCTokenMessenger internal s_mockUSDC;
  MockE2EUSDCTransmitter internal s_mockUSDCTransmitter;
  uint32 internal constant USDC_DEST_TOKEN_GAS = 150_000;

  struct USDCMessage {
    uint32 version;
    uint32 sourceDomain;
    uint32 destinationDomain;
    uint64 nonce;
    bytes32 sender;
    bytes32 recipient;
    bytes32 destinationCaller;
    bytes messageBody;
  }

  uint32 internal constant SOURCE_DOMAIN_IDENTIFIER = 0x02020202;
  uint32 internal constant DEST_DOMAIN_IDENTIFIER = 0;

  bytes32 internal constant SOURCE_CHAIN_TOKEN_SENDER = bytes32(uint256(uint160(0x01111111221)));
  address internal constant SOURCE_CHAIN_USDC_POOL = address(0x23789765456789);
  address internal constant DEST_CHAIN_USDC_POOL = address(0x987384873458734);
  address internal constant DEST_CHAIN_USDC_TOKEN = address(0x23598918358198766);

  address internal s_routerAllowedOnRamp = address(3456);
  address internal s_routerAllowedOffRamp = address(234);
  Router internal s_router;

  HybridLockReleaseUSDCTokenPool internal s_usdcTokenPool;
  HybridLockReleaseUSDCTokenPool internal s_usdcTokenPoolTransferLiquidity;
  address[] internal s_allowedList;

  function setUp() public virtual override {
    BaseTest.setUp();
    BurnMintERC677 usdcToken = new BurnMintERC677("LINK", "LNK", 18, 0);
    s_token = usdcToken;
    deal(address(s_token), OWNER, type(uint256).max);
    _setUpRamps();

    s_mockUSDCTransmitter = new MockE2EUSDCTransmitter(0, DEST_DOMAIN_IDENTIFIER, address(s_token));
    s_mockUSDC = new MockUSDCTokenMessenger(0, address(s_mockUSDCTransmitter));

    usdcToken.grantMintAndBurnRoles(address(s_mockUSDCTransmitter));

    s_usdcTokenPool =
      new HybridLockReleaseUSDCTokenPool(s_mockUSDC, s_token, new address[](0), address(s_mockRMN), address(s_router));

    s_usdcTokenPoolTransferLiquidity =
      new HybridLockReleaseUSDCTokenPool(s_mockUSDC, s_token, new address[](0), address(s_mockRMN), address(s_router));

    usdcToken.grantMintAndBurnRoles(address(s_mockUSDC));
    usdcToken.grantMintAndBurnRoles(address(s_usdcTokenPool));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](2);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: SOURCE_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      remoteTokenAddress: abi.encode(address(s_token)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    chainUpdates[1] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(DEST_CHAIN_USDC_POOL),
      remoteTokenAddress: abi.encode(DEST_CHAIN_USDC_TOKEN),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_usdcTokenPool.applyChainUpdates(chainUpdates);

    USDCTokenPool.DomainUpdate[] memory domains = new USDCTokenPool.DomainUpdate[](1);
    domains[0] = USDCTokenPool.DomainUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      domainIdentifier: 9999,
      allowedCaller: keccak256("allowedCaller"),
      enabled: true
    });

    s_usdcTokenPool.setDomains(domains);

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LiquidityProviderSet(address(0), OWNER, DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.setLiquidityProvider(DEST_CHAIN_SELECTOR, OWNER);
    s_usdcTokenPool.setLiquidityProvider(SOURCE_CHAIN_SELECTOR, OWNER);
  }

  function _setUpRamps() internal {
    s_router = new Router(address(s_token), address(s_mockRMN));

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: s_routerAllowedOnRamp});
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    address[] memory offRamps = new address[](1);
    offRamps[0] = s_routerAllowedOffRamp;
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: offRamps[0]});

    s_router.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }

  function _generateUSDCMessage(
    USDCMessage memory usdcMessage
  ) internal pure returns (bytes memory) {
    return abi.encodePacked(
      usdcMessage.version,
      usdcMessage.sourceDomain,
      usdcMessage.destinationDomain,
      usdcMessage.nonce,
      usdcMessage.sender,
      usdcMessage.recipient,
      usdcMessage.destinationCaller,
      usdcMessage.messageBody
    );
  }
}

contract USDCBridgeMigrator_releaseOrMint is HybridLockReleaseUSDCTokenPool_releaseOrMint {
  function test_unstickManualTxAfterMigration_destChain_Success() public {
    address recipient = address(1234);
    // Test the edge case where a tx is stuck in the manual tx queue and the destination chain is the one that
    // should process is after a migration. I.E the message will have the Lock-Release flag set in the OffChainData,
    // which should tell it to use the lock-release mechanism with the tokens provided.

    // We want the released amount to be 1e6, so to simulate the workflow, we sent those tokens to the contract as
    // liquidity
    uint256 amount = 1e6;
    // Add 1e12 liquidity so that there's enough to release
    vm.startPrank(s_usdcTokenPool.getLiquidityProvider(SOURCE_CHAIN_SELECTOR));

    s_token.approve(address(s_usdcTokenPool), type(uint256).max);
    s_usdcTokenPool.provideLiquidity(SOURCE_CHAIN_SELECTOR, amount);

    // By Default, the source chain will be indicated as use-CCTP so we need to change that. We create a message
    // that will use the Lock-Release flag in the offchain data to indicate that the tokens should be released
    // instead of minted since there's no attestation for us to use.

    vm.startPrank(s_routerAllowedOffRamp);

    vm.expectEmit();
    emit TokenPool.Released(s_routerAllowedOffRamp, recipient, amount);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    Pool.ReleaseOrMintOutV1 memory poolReturnDataV1 = s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    // By this point, the tx should have executed, with the Lock-Release taking over, and being forwaded to the
    // recipient

    assertEq(poolReturnDataV1.destinationAmount, amount, "destinationAmount and actual amount transferred differ");
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 0, "Tokens should be transferred out of the pool");
    assertEq(s_token.balanceOf(recipient), amount, "Tokens should be transferred to the recipient");

    // We also want to check that the system uses CCTP Burn/Mint for all other messages that don't have that flag
    // which after a migration will mean all new messages.

    // The message should fail without an error because it failed to decode a non-existent attestation which would
    // revert without an error
    vm.expectRevert();

    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_unstickManualTxAfterMigration_homeChain_Success() public {
    address CIRCLE = makeAddr("CIRCLE");
    address recipient = address(1234);

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = SOURCE_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    // Test the edge case where a tx is stuck in the manual tx queue and the source chain (mainnet) needs unsticking
    // In this test we want 1e6 worth of tokens to be stuck, so first we provide liquidity to the pool >1e6

    uint256 amount = 1e6;
    // Add 1e12 liquidity so that there's enough to release
    vm.startPrank(s_usdcTokenPool.getLiquidityProvider(SOURCE_CHAIN_SELECTOR));

    s_token.approve(address(s_usdcTokenPool), type(uint256).max);

    // I picked 3x the amount to be stuck so that we can have enough to release with a buffer
    s_usdcTokenPool.provideLiquidity(SOURCE_CHAIN_SELECTOR, amount * 3);

    // At this point in the process, the router will lock new messages, so we want to simulate excluding tokens
    // stuck coming back from the destination, to the home chain. This way they can be released and not minted
    // since there's no corresponding attestation to use for minting.
    vm.startPrank(OWNER);

    s_usdcTokenPool.proposeCCTPMigration(SOURCE_CHAIN_SELECTOR);

    // Exclude the tokens from being burned and check for the event
    vm.expectEmit();
    emit USDCBridgeMigrator.TokensExcludedFromBurn(SOURCE_CHAIN_SELECTOR, amount, (amount * 3) - amount);

    s_usdcTokenPool.excludeTokensFromBurn(SOURCE_CHAIN_SELECTOR, amount);

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(SOURCE_CHAIN_SELECTOR),
      (amount * 3),
      "Tokens locked minus ones excluded from the burn should be 2e6"
    );

    assertEq(
      s_usdcTokenPool.getExcludedTokensByChain(SOURCE_CHAIN_SELECTOR),
      1e6,
      "1e6 tokens should be excluded from the burn"
    );

    s_usdcTokenPool.setCircleMigratorAddress(CIRCLE);

    vm.startPrank(CIRCLE);

    s_usdcTokenPool.burnLockedUSDC();

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(SOURCE_CHAIN_SELECTOR), 0, "All tokens should be burned out of the pool"
    );

    assertEq(
      s_usdcTokenPool.getExcludedTokensByChain(SOURCE_CHAIN_SELECTOR),
      1e6,
      "There should still be 1e6 tokens excluded from the burn"
    );

    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 1e6, "All tokens minus the excluded should be in the pool");

    // Now that the burn is successful, we can release the tokens that were excluded from the burn
    vm.startPrank(s_routerAllowedOffRamp);

    vm.expectEmit();
    emit TokenPool.Released(s_routerAllowedOffRamp, recipient, amount);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    Pool.ReleaseOrMintOutV1 memory poolReturnDataV1 = s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    assertEq(poolReturnDataV1.destinationAmount, amount, "destinationAmount and actual amount transferred differ");
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 0, "Tokens should be transferred out of the pool");
    assertEq(s_token.balanceOf(recipient), amount, "Tokens should be transferred to the recipient");
    assertEq(
      s_usdcTokenPool.getExcludedTokensByChain(SOURCE_CHAIN_SELECTOR),
      0,
      "All tokens should be released from the exclusion list"
    );

    // We also want to check that the system uses CCTP Burn/Mint for all other messages that don't have that flag
    test_incomingMessageWithPrimaryMechanism();
  }
}
