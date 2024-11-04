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

contract HybridLockReleaseUSDCTokenPool_releaseOrMint is USDCTokenPoolSetup {
  function test_OnLockReleaseMechanism_Success() public {
    address recipient = address(1234);

    // Designate the SOURCE_CHAIN as not using native-USDC, and so the L/R mechanism must be used instead
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = SOURCE_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(SOURCE_CHAIN_SELECTOR),
      "Lock/Release mech not configured for incoming message from SOURCE_CHAIN_SELECTOR"
    );

    vm.startPrank(OWNER);
    s_usdcTokenPool.setLiquidityProvider(SOURCE_CHAIN_SELECTOR, OWNER);

    // Add 1e12 liquidity so that there's enough to release
    vm.startPrank(s_usdcTokenPool.getLiquidityProvider(SOURCE_CHAIN_SELECTOR));

    s_token.approve(address(s_usdcTokenPool), type(uint256).max);

    uint256 liquidityAmount = 1e12;
    s_usdcTokenPool.provideLiquidity(SOURCE_CHAIN_SELECTOR, liquidityAmount);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    uint256 amount = 1e6;

    vm.startPrank(s_routerAllowedOffRamp);

    vm.expectEmit();
    emit TokenPool.Released(s_routerAllowedOffRamp, recipient, amount);

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

    // Simulate the off-ramp forwarding tokens to the recipient on destination chain
    // s_token.transfer(recipient, amount);

    assertEq(
      s_token.balanceOf(address(s_usdcTokenPool)),
      liquidityAmount - amount,
      "Incorrect remaining liquidity in TokenPool"
    );
    assertEq(s_token.balanceOf(recipient), amount, "Tokens not transferred to recipient");
  }

  // https://etherscan.io/tx/0xac9f501fe0b76df1f07a22e1db30929fd12524bc7068d74012dff948632f0883
  function test_incomingMessageWithPrimaryMechanism() public {
    bytes memory encodedUsdcMessage =
      hex"000000000000000300000000000000000000127a00000000000000000000000019330d10d9cc8751218eaf51e8885d058642e08a000000000000000000000000bd3fa81b58ba92a82136038b25adec7066af3155000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000af88d065e77c8cc2239327c5edb3a432268e58310000000000000000000000004af08f56978be7dce2d1be3c65c005b41e79401c000000000000000000000000000000000000000000000000000000002057ff7a0000000000000000000000003a23f943181408eac424116af7b7790c94cb97a50000000000000000000000000000000000000000000000000000000000000000000000000000008274119237535fd659626b090f87e365ff89ebc7096bb32e8b0e85f155626b73ae7c4bb2485c184b7cc3cf7909045487890b104efb62ae74a73e32901bdcec91df1bb9ee08ccb014fcbcfe77b74d1263fd4e0b0e8de05d6c9a5913554364abfd5ea768b222f50c715908183905d74044bb2b97527c7e70ae7983c443a603557cac3b1c000000000000000000000000000000000000000000000000000000000000";
    bytes memory attestation = bytes("attestation bytes");

    uint32 nonce = 4730;
    uint32 sourceDomain = 3;
    uint256 amount = 100;

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: nonce, sourceDomain: sourceDomain})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    // The mocked receiver does not release the token to the pool, so we manually do it here
    deal(address(s_token), address(s_usdcTokenPool), amount);

    bytes memory offchainTokenData =
      abi.encode(USDCTokenPool.MessageAndAttestation({message: encodedUsdcMessage, attestation: attestation}));

    vm.expectCall(
      address(s_mockUSDCTransmitter),
      abi.encodeWithSelector(MockE2EUSDCTransmitter.receiveMessage.selector, encodedUsdcMessage, attestation)
    );

    vm.startPrank(s_routerAllowedOffRamp);
    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }

  function test_WhileMigrationPause_Revert() public {
    address recipient = address(1234);

    // Designate the SOURCE_CHAIN as not using native-USDC, and so the L/R mechanism must be used instead
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = SOURCE_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(SOURCE_CHAIN_SELECTOR),
      "Lock/Release mech not configured for incoming message from SOURCE_CHAIN_SELECTOR"
    );

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationProposed(SOURCE_CHAIN_SELECTOR);

    // Propose the migration to CCTP
    s_usdcTokenPool.proposeCCTPMigration(SOURCE_CHAIN_SELECTOR);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory sourcePoolDataLockRelease = abi.encode(LOCK_RELEASE_FLAG);

    uint256 amount = 1e6;

    vm.startPrank(s_routerAllowedOffRamp);

    // Expect revert because the lane is paused and no incoming messages should be allowed
    vm.expectRevert(
      abi.encodeWithSelector(HybridLockReleaseUSDCTokenPool.LanePausedForCCTPMigration.selector, SOURCE_CHAIN_SELECTOR)
    );

    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourcePoolDataLockRelease,
        offchainTokenData: ""
      })
    );
  }
}
