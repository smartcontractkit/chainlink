// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ILiquidityContainer} from "../../../liquiditymanager/interfaces/ILiquidityContainer.sol";
import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";
import {IPoolV1} from "../../interfaces/IPool.sol";
import {ITokenMessenger} from "../../pools/USDC/ITokenMessenger.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../Router.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";

import {TokenPool} from "../../pools/TokenPool.sol";
import {HybridLockReleaseUSDCTokenPool} from "../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCBridgeMigrator} from "../../pools/USDC/USDCBridgeMigrator.sol";
import {USDCTokenPool} from "../../pools/USDC/USDCTokenPool.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {USDCTokenPoolHelper} from "../helpers/USDCTokenPoolHelper.sol";
import {MockE2EUSDCTransmitter} from "../mocks/MockE2EUSDCTransmitter.sol";
import {MockUSDCTokenMessenger} from "../mocks/MockUSDCTokenMessenger.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

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
    setUpRamps();

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
  }

  function setUpRamps() internal {
    s_router = new Router(address(s_token), address(s_mockRMN));

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: s_routerAllowedOnRamp});
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    address[] memory offRamps = new address[](1);
    offRamps[0] = s_routerAllowedOffRamp;
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: offRamps[0]});

    s_router.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }

  function _generateUSDCMessage(USDCMessage memory usdcMessage) internal pure returns (bytes memory) {
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

contract HybridUSDCTokenPoolTests is USDCTokenPoolSetup {
  function test_LockOrBurn_onLockReleaseMechanism_Success() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR),
      "Lock/Release mech not configured for outgoing message to DEST_CHAIN_SELECTOR"
    );

    uint256 amount = 1e6;

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    vm.expectEmit();
    emit TokenPool.Locked(s_routerAllowedOnRamp, amount);

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), amount, "Incorrect token amount in the tokenPool");
  }

  function test_MintOrRelease_OnLockReleaseMechanism_Success() public {
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
        sourcePoolData: "",
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

  function test_LockOrBurn_PrimaryMechanism_Success() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));
    uint256 amount = 1;

    vm.startPrank(OWNER);

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    USDCTokenPool.Domain memory expectedDomain = s_usdcTokenPool.getDomain(DEST_CHAIN_SELECTOR);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);

    vm.expectEmit();
    emit ITokenMessenger.DepositForBurn(
      s_mockUSDC.s_nonce(),
      address(s_token),
      amount,
      address(s_usdcTokenPool),
      receiver,
      expectedDomain.domainIdentifier,
      s_mockUSDC.DESTINATION_TOKEN_MESSENGER(),
      expectedDomain.allowedCaller
    );

    vm.expectEmit();
    emit TokenPool.Burned(s_routerAllowedOnRamp, amount);

    Pool.LockOrBurnOutV1 memory poolReturnDataV1 = s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    uint64 nonce = abi.decode(poolReturnDataV1.destPoolData, (uint64));
    assertEq(s_mockUSDC.s_nonce() - 1, nonce);
  }

  // https://etherscan.io/tx/0xac9f501fe0b76df1f07a22e1db30929fd12524bc7068d74012dff948632f0883
  function test_MintOrRelease_incomingMessageWithPrimaryMechanism() public {
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

  function test_LockOrBurn_LocKReleaseMechanism_then_switchToPrimary_Success() public {
    // Test Enabling the LR mechanism and sending an outgoing message
    test_LockOrBurn_PrimaryMechanism_Success();

    // Disable the LR mechanism so that primary CCTP is used and then attempt to send a message
    uint64[] memory destChainRemoves = new uint64[](1);
    destChainRemoves[0] = DEST_CHAIN_SELECTOR;

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LockReleaseDisabled(DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.updateChainSelectorMechanisms(destChainRemoves, new uint64[](0));

    // Send an outgoing message
    test_LockOrBurn_PrimaryMechanism_Success();
  }

  function test_MintOrRelease_OnLockReleaseMechanism_then_switchToPrimary_Success() public {
    test_MintOrRelease_OnLockReleaseMechanism_Success();

    // Disable the LR mechanism so that primary CCTP is used and then attempt to send a message
    uint64[] memory destChainRemoves = new uint64[](1);
    destChainRemoves[0] = SOURCE_CHAIN_SELECTOR;

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LockReleaseDisabled(SOURCE_CHAIN_SELECTOR);

    s_usdcTokenPool.updateChainSelectorMechanisms(destChainRemoves, new uint64[](0));

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LiquidityProviderSet(OWNER, OWNER, SOURCE_CHAIN_SELECTOR);

    s_usdcTokenPool.setLiquidityProvider(SOURCE_CHAIN_SELECTOR, OWNER);

    // Test incoming on the primary mechanism after disable LR, simulating Circle's new support for CCTP on
    // DEST_CHAIN_SELECTOR
    test_MintOrRelease_incomingMessageWithPrimaryMechanism();
  }

  function test_withdrawLiquidity_Success() public {
    uint256 liquidityAmount = 1e12;

    vm.startPrank(OWNER);
    s_token.approve(address(s_usdcTokenPool), type(uint256).max);

    s_usdcTokenPool.setLiquidityProvider(SOURCE_CHAIN_SELECTOR, OWNER);
    assertEq(s_usdcTokenPool.getLiquidityProvider(SOURCE_CHAIN_SELECTOR), OWNER, "Owner is not Rebalancer");

    s_usdcTokenPool.provideLiquidity(SOURCE_CHAIN_SELECTOR, liquidityAmount);

    assertEq(s_usdcTokenPool.getLockedTokensForChain(SOURCE_CHAIN_SELECTOR), liquidityAmount);

    assertEq(
      s_token.balanceOf(address(s_usdcTokenPool)),
      liquidityAmount,
      "Available tokens doesn't match provided liquidity amount"
    );

    uint256 withdrawalAmount = liquidityAmount / 2; // Withdraw half of the liquidity

    vm.expectEmit();
    emit ILiquidityContainer.LiquidityRemoved(OWNER, withdrawalAmount);

    s_usdcTokenPool.withdrawLiquidity(SOURCE_CHAIN_SELECTOR, withdrawalAmount);

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(SOURCE_CHAIN_SELECTOR),
      liquidityAmount - withdrawalAmount,
      "Remaining liquidity incorrect"
    );
    assertEq(
      s_token.balanceOf(address(s_usdcTokenPool)),
      liquidityAmount - withdrawalAmount,
      "Available tokens doesn't match provided liquidity amount"
    );
  }

  function test_LockOrBurn_WhileMigrationPause_Revert() public {
    // Create a fake migration proposal
    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    assertEq(s_usdcTokenPool.getCurrentProposedCCTPChainMigration(), DEST_CHAIN_SELECTOR);

    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR),
      "Lock Release mech not configured for outgoing message to DEST_CHAIN_SELECTOR"
    );

    uint256 amount = 1e6;

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    // Expect the lockOrBurn to fail because a pending CCTP-Migration has paused outgoing messages on CCIP
    vm.expectRevert(
      abi.encodeWithSelector(HybridLockReleaseUSDCTokenPool.LanePausedForCCTPMigration.selector, DEST_CHAIN_SELECTOR)
    );

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }
}

contract HybridUSDCTokenPoolMigrationTests is HybridUSDCTokenPoolTests {
  function test_lockOrBurn_then_BurnInCCTPMigration_Success() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));
    address CIRCLE = makeAddr("CIRCLE CCTP Migrator");

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR),
      "Lock/Release mech not configured for outgoing message to DEST_CHAIN_SELECTOR"
    );

    uint256 amount = 1e6;

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    vm.expectEmit();
    emit TokenPool.Locked(s_routerAllowedOnRamp, amount);

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    // Ensure that the tokens are properly locked
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), amount, "Incorrect token amount in the tokenPool");

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      amount,
      "Internal locked token accounting is incorrect"
    );

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit USDCBridgeMigrator.CircleMigratorAddressSet(CIRCLE);

    s_usdcTokenPool.setCircleMigratorAddress(CIRCLE);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationProposed(DEST_CHAIN_SELECTOR);

    // Propose the migration to CCTP
    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    assertEq(
      s_usdcTokenPool.getCurrentProposedCCTPChainMigration(),
      DEST_CHAIN_SELECTOR,
      "Current proposed chain migration does not match expected for DEST_CHAIN_SELECTOR"
    );

    // Impersonate the set circle address and execute the proposal
    vm.startPrank(CIRCLE);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationExecuted(DEST_CHAIN_SELECTOR, amount);

    // Ensure the call to the burn function is properly
    vm.expectCall(address(s_token), abi.encodeWithSelector(bytes4(keccak256("burn(uint256)")), amount));

    s_usdcTokenPool.burnLockedUSDC();

    // Assert that the tokens were actually burned
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 0, "Tokens were not burned out of the tokenPool");

    // Ensure the proposal slot was cleared and there's no tokens locked for the destination chain anymore
    assertEq(s_usdcTokenPool.getCurrentProposedCCTPChainMigration(), 0, "Proposal Slot should be empty");
    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      0,
      "No tokens should be locked for DEST_CHAIN_SELECTOR after CCTP-approved burn"
    );

    assertFalse(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR), "Lock/Release mech should be disabled after a burn"
    );

    test_LockOrBurn_PrimaryMechanism_Success();
  }

  function test_cancelExistingCCTPMigrationProposal() public {
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationProposed(DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    assertEq(
      s_usdcTokenPool.getCurrentProposedCCTPChainMigration(),
      DEST_CHAIN_SELECTOR,
      "migration proposal should exist, but doesn't"
    );

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationCancelled(DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.cancelExistingCCTPMigrationProposal();

    assertEq(
      s_usdcTokenPool.getCurrentProposedCCTPChainMigration(),
      0,
      "migration proposal exists, but shouldn't after being cancelled"
    );

    vm.expectRevert(USDCBridgeMigrator.NoExistingMigrationProposal.selector);
    s_usdcTokenPool.cancelExistingCCTPMigrationProposal();
  }

  function test_burnLockedUSDC_invalidPermissions_Revert() public {
    address CIRCLE = makeAddr("CIRCLE");

    vm.startPrank(OWNER);

    // Set the circle migrator address for later, but don't start pranking as it yet
    s_usdcTokenPool.setCircleMigratorAddress(CIRCLE);

    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.onlyCircle.selector));

    // Should fail because only Circle can call this function
    s_usdcTokenPool.burnLockedUSDC();

    vm.startPrank(CIRCLE);

    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.ExistingMigrationProposal.selector));
    s_usdcTokenPool.burnLockedUSDC();
  }

  function test_transferLiquidity_Success() public {
    // Set as the OWNER so we can provide liquidity
    vm.startPrank(OWNER);
    s_usdcTokenPoolTransferLiquidity.setLiquidityProvider(DEST_CHAIN_SELECTOR, OWNER);

    s_token.approve(address(s_usdcTokenPoolTransferLiquidity), type(uint256).max);

    uint256 liquidityAmount = 1e9;

    // Provide 1000 USDC as liquidity
    s_usdcTokenPoolTransferLiquidity.provideLiquidity(DEST_CHAIN_SELECTOR, liquidityAmount);

    // Set the new token pool as the rebalancer
    s_usdcTokenPoolTransferLiquidity.setLiquidityProvider(DEST_CHAIN_SELECTOR, address(s_usdcTokenPool));

    vm.expectEmit();
    emit ILiquidityContainer.LiquidityRemoved(address(s_usdcTokenPool), liquidityAmount);

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LiquidityTransferred(
      address(s_usdcTokenPoolTransferLiquidity), DEST_CHAIN_SELECTOR, liquidityAmount
    );

    s_usdcTokenPool.transferLiquidity(address(s_usdcTokenPoolTransferLiquidity), DEST_CHAIN_SELECTOR, liquidityAmount);

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      liquidityAmount,
      "Tokens locked for dest chain doesn't match expected amount in storage"
    );

    assertEq(
      s_usdcTokenPoolTransferLiquidity.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      0,
      "Tokens locked for dest chain in old token pool doesn't match expected amount in storage"
    );

    assertEq(
      s_token.balanceOf(address(s_usdcTokenPool)),
      liquidityAmount,
      "Liquidity amount of tokens should be new in new pool, but aren't"
    );
  }

  function test_cannotModifyLiquidityWithoutPermissions_Revert() public {
    address randomAddr = makeAddr("RANDOM");

    vm.startPrank(randomAddr);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, randomAddr));

    // Revert because there's insufficient permissions for the DEST_CHAIN_SELECTOR to provide liquidity
    s_usdcTokenPool.provideLiquidity(DEST_CHAIN_SELECTOR, 1e6);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.Unauthorized.selector, randomAddr));

    // Revert because there's insufficient permissions for the DEST_CHAIN_SELECTOR to withdraw liquidity
    s_usdcTokenPool.withdrawLiquidity(DEST_CHAIN_SELECTOR, 1e6);
  }

  function test_cannotCancelANonExistentMigrationProposal() public {
    vm.expectRevert(USDCBridgeMigrator.NoExistingMigrationProposal.selector);

    // Proposal to migrate doesn't exist, and so the chain selector is zero, and therefore should revert
    s_usdcTokenPool.cancelExistingCCTPMigrationProposal();
  }
}
