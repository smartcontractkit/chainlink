// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";
import {IPoolV1} from "../../interfaces/IPool.sol";
import {ITokenMessenger} from "../../pools/USDC/ITokenMessenger.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {Router} from "../../Router.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {USDCTokenPool} from "../../pools/USDC/USDCTokenPool.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {USDCTokenPoolHelper} from "../helpers/USDCTokenPoolHelper.sol";
import {MockE2EUSDCTransmitter} from "../mocks/MockE2EUSDCTransmitter.sol";
import {MockUSDCTokenMessenger} from "../mocks/MockUSDCTokenMessenger.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

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

  USDCTokenPoolHelper internal s_usdcTokenPool;
  USDCTokenPoolHelper internal s_usdcTokenPoolWithAllowList;
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
      new USDCTokenPoolHelper(s_mockUSDC, s_token, new address[](0), address(s_mockRMN), address(s_router));
    usdcToken.grantMintAndBurnRoles(address(s_mockUSDC));

    s_allowedList.push(USER_1);
    s_usdcTokenPoolWithAllowList =
      new USDCTokenPoolHelper(s_mockUSDC, s_token, s_allowedList, address(s_mockRMN), address(s_router));

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
    s_usdcTokenPoolWithAllowList.applyChainUpdates(chainUpdates);

    USDCTokenPool.DomainUpdate[] memory domains = new USDCTokenPool.DomainUpdate[](1);
    domains[0] = USDCTokenPool.DomainUpdate({
      destChainSelector: DEST_CHAIN_SELECTOR,
      domainIdentifier: 9999,
      allowedCaller: keccak256("allowedCaller"),
      enabled: true
    });

    s_usdcTokenPool.setDomains(domains);
    s_usdcTokenPoolWithAllowList.setDomains(domains);
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

contract USDCTokenPool_lockOrBurn is USDCTokenPoolSetup {
  // Base test case, included for PR gas comparisons as fuzz tests are excluded from forge snapshot due to being flaky.
  function test_LockOrBurn_Success() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));
    uint256 amount = 1;
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

  function test_Fuzz_LockOrBurn_Success(bytes32 destinationReceiver, uint256 amount) public {
    vm.assume(destinationReceiver != bytes32(0));
    amount = bound(amount, 1, _getOutboundRateLimiterConfig().capacity);
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
      destinationReceiver,
      expectedDomain.domainIdentifier,
      s_mockUSDC.DESTINATION_TOKEN_MESSENGER(),
      expectedDomain.allowedCaller
    );

    vm.expectEmit();
    emit TokenPool.Burned(s_routerAllowedOnRamp, amount);

    Pool.LockOrBurnOutV1 memory poolReturnDataV1 = s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(destinationReceiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    uint64 nonce = abi.decode(poolReturnDataV1.destPoolData, (uint64));
    assertEq(s_mockUSDC.s_nonce() - 1, nonce);
    assertEq(poolReturnDataV1.destTokenAddress, abi.encode(DEST_CHAIN_USDC_TOKEN));
  }

  function test_Fuzz_LockOrBurnWithAllowList_Success(bytes32 destinationReceiver, uint256 amount) public {
    vm.assume(destinationReceiver != bytes32(0));
    amount = bound(amount, 1, _getOutboundRateLimiterConfig().capacity);
    s_token.transfer(address(s_usdcTokenPoolWithAllowList), amount);
    vm.startPrank(s_routerAllowedOnRamp);

    USDCTokenPool.Domain memory expectedDomain = s_usdcTokenPoolWithAllowList.getDomain(DEST_CHAIN_SELECTOR);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);
    vm.expectEmit();
    emit ITokenMessenger.DepositForBurn(
      s_mockUSDC.s_nonce(),
      address(s_token),
      amount,
      address(s_usdcTokenPoolWithAllowList),
      destinationReceiver,
      expectedDomain.domainIdentifier,
      s_mockUSDC.DESTINATION_TOKEN_MESSENGER(),
      expectedDomain.allowedCaller
    );
    vm.expectEmit();
    emit TokenPool.Burned(s_routerAllowedOnRamp, amount);

    Pool.LockOrBurnOutV1 memory poolReturnDataV1 = s_usdcTokenPoolWithAllowList.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: s_allowedList[0],
        receiver: abi.encodePacked(destinationReceiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
    uint64 nonce = abi.decode(poolReturnDataV1.destPoolData, (uint64));
    assertEq(s_mockUSDC.s_nonce() - 1, nonce);
    assertEq(poolReturnDataV1.destTokenAddress, abi.encode(DEST_CHAIN_USDC_TOKEN));
  }

  // Reverts
  function test_UnknownDomain_Revert() public {
    uint64 wrongDomain = DEST_CHAIN_SELECTOR + 1;
    // We need to setup the wrong chainSelector so it reaches the domain check
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: wrongDomain, onRamp: s_routerAllowedOnRamp});
    s_router.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: wrongDomain,
      remotePoolAddress: abi.encode(address(1)),
      remoteTokenAddress: abi.encode(address(2)),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_usdcTokenPool.applyChainUpdates(chainUpdates);

    uint256 amount = 1000;
    vm.startPrank(s_routerAllowedOnRamp);
    deal(address(s_token), s_routerAllowedOnRamp, amount);
    s_token.approve(address(s_usdcTokenPool), amount);

    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.UnknownDomain.selector, wrongDomain));

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(address(0)),
        amount: amount,
        remoteChainSelector: wrongDomain,
        localToken: address(s_token)
      })
    );
  }

  function test_CallerIsNotARampOnRouter_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(TokenPool.CallerIsNotARampOnRouter.selector, OWNER));

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(address(0)),
        amount: 0,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }

  function test_LockOrBurnWithAllowList_Revert() public {
    vm.startPrank(s_routerAllowedOnRamp);

    vm.expectRevert(abi.encodeWithSelector(TokenPool.SenderNotAllowed.selector, STRANGER));

    s_usdcTokenPoolWithAllowList.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: abi.encodePacked(address(0)),
        amount: 1000,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }
}

contract USDCTokenPool_releaseOrMint is USDCTokenPoolSetup {
  // From https://github.com/circlefin/evm-cctp-contracts/blob/377c9bd813fb86a42d900ae4003599d82aef635a/src/messages/BurnMessage.sol#L57
  function _formatMessage(
    uint32 _version,
    bytes32 _burnToken,
    bytes32 _mintRecipient,
    uint256 _amount,
    bytes32 _messageSender
  ) internal pure returns (bytes memory) {
    return abi.encodePacked(_version, _burnToken, _mintRecipient, _amount, _messageSender);
  }

  function test_Fuzz_ReleaseOrMint_Success(address recipient, uint256 amount) public {
    vm.assume(recipient != address(0) && recipient != address(s_token));
    amount = bound(amount, 0, _getInboundRateLimiterConfig().capacity);

    USDCMessage memory usdcMessage = USDCMessage({
      version: 0,
      sourceDomain: SOURCE_DOMAIN_IDENTIFIER,
      destinationDomain: DEST_DOMAIN_IDENTIFIER,
      nonce: 0x060606060606,
      sender: SOURCE_CHAIN_TOKEN_SENDER,
      recipient: bytes32(uint256(uint160(recipient))),
      destinationCaller: bytes32(uint256(uint160(address(s_usdcTokenPool)))),
      messageBody: _formatMessage(
        0,
        bytes32(uint256(uint160(address(s_token)))),
        bytes32(uint256(uint160(recipient))),
        amount,
        bytes32(uint256(uint160(OWNER)))
      )
    });

    bytes memory message = _generateUSDCMessage(usdcMessage);
    bytes memory attestation = bytes("attestation bytes");

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(
        USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})
      ),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory offchainTokenData =
      abi.encode(USDCTokenPool.MessageAndAttestation({message: message, attestation: attestation}));

    // The mocked receiver does not release the token to the pool, so we manually do it here
    deal(address(s_token), address(s_usdcTokenPool), amount);

    vm.expectEmit();
    emit TokenPool.Minted(s_routerAllowedOffRamp, recipient, amount);

    vm.expectCall(
      address(s_mockUSDCTransmitter),
      abi.encodeWithSelector(MockE2EUSDCTransmitter.receiveMessage.selector, message, attestation)
    );

    vm.startPrank(s_routerAllowedOffRamp);
    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }

  // https://etherscan.io/tx/0xac9f501fe0b76df1f07a22e1db30929fd12524bc7068d74012dff948632f0883
  function test_ReleaseOrMintRealTx_Success() public {
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

  // Reverts
  function test_UnlockingUSDCFailed_Revert() public {
    vm.startPrank(s_routerAllowedOffRamp);
    s_mockUSDCTransmitter.setShouldSucceed(false);

    uint256 amount = 13255235235;

    USDCMessage memory usdcMessage = USDCMessage({
      version: 0,
      sourceDomain: SOURCE_DOMAIN_IDENTIFIER,
      destinationDomain: DEST_DOMAIN_IDENTIFIER,
      nonce: 0x060606060606,
      sender: SOURCE_CHAIN_TOKEN_SENDER,
      recipient: bytes32(uint256(uint160(address(s_mockUSDC)))),
      destinationCaller: bytes32(uint256(uint160(address(s_usdcTokenPool)))),
      messageBody: _formatMessage(
        0,
        bytes32(uint256(uint160(address(s_token)))),
        bytes32(uint256(uint160(OWNER))),
        amount,
        bytes32(uint256(uint160(OWNER)))
      )
    });

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(
        USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})
      ),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory offchainTokenData = abi.encode(
      USDCTokenPool.MessageAndAttestation({message: _generateUSDCMessage(usdcMessage), attestation: bytes("")})
    );

    vm.expectRevert(USDCTokenPool.UnlockingUSDCFailed.selector);

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

  function test_TokenMaxCapacityExceeded_Revert() public {
    uint256 capacity = _getInboundRateLimiterConfig().capacity;
    uint256 amount = 10 * capacity;
    address recipient = address(1);
    vm.startPrank(s_routerAllowedOffRamp);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    bytes memory offchainTokenData =
      abi.encode(USDCTokenPool.MessageAndAttestation({message: bytes(""), attestation: bytes("")}));

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.TokenMaxCapacityExceeded.selector, capacity, amount, address(s_token))
    );

    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: sourceTokenData.extraData,
        offchainTokenData: offchainTokenData
      })
    );
  }
}

contract USDCTokenPool_supportsInterface is USDCTokenPoolSetup {
  function test_SupportsInterface_Success() public view {
    assertTrue(s_usdcTokenPool.supportsInterface(type(IPoolV1).interfaceId));
    assertTrue(s_usdcTokenPool.supportsInterface(type(IERC165).interfaceId));
  }
}

contract USDCTokenPool_setDomains is USDCTokenPoolSetup {
  mapping(uint64 destChainSelector => USDCTokenPool.Domain domain) private s_chainToDomain;

  // Setting lower fuzz run as 256 runs was causing differing gas results in snapshot.
  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function test_Fuzz_SetDomains_Success(
    bytes32[5] calldata allowedCallers,
    uint32[5] calldata domainIdentifiers,
    uint64[5] calldata destChainSelectors
  ) public {
    uint256 numberOfDomains = allowedCallers.length;
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](numberOfDomains);
    for (uint256 i = 0; i < numberOfDomains; ++i) {
      vm.assume(allowedCallers[i] != bytes32(0) && domainIdentifiers[i] != 0 && destChainSelectors[i] != 0);

      domainUpdates[i] = USDCTokenPool.DomainUpdate({
        allowedCaller: allowedCallers[i],
        domainIdentifier: domainIdentifiers[i],
        destChainSelector: destChainSelectors[i],
        enabled: true
      });

      s_chainToDomain[destChainSelectors[i]] =
        USDCTokenPool.Domain({domainIdentifier: domainIdentifiers[i], allowedCaller: allowedCallers[i], enabled: true});
    }

    vm.expectEmit();
    emit USDCTokenPool.DomainsSet(domainUpdates);

    s_usdcTokenPool.setDomains(domainUpdates);

    for (uint256 i = 0; i < numberOfDomains; ++i) {
      USDCTokenPool.Domain memory expected = s_chainToDomain[destChainSelectors[i]];
      USDCTokenPool.Domain memory got = s_usdcTokenPool.getDomain(destChainSelectors[i]);
      assertEq(got.allowedCaller, expected.allowedCaller);
      assertEq(got.domainIdentifier, expected.domainIdentifier);
    }
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](0);

    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");

    s_usdcTokenPool.setDomains(domainUpdates);
  }

  function test_InvalidDomain_Revert() public {
    bytes32 validCaller = bytes32(uint256(25));
    // Ensure valid domain works
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](1);
    domainUpdates[0] = USDCTokenPool.DomainUpdate({
      allowedCaller: validCaller,
      domainIdentifier: 0, // ensures 0 is valid, as this is eth mainnet
      destChainSelector: 45690,
      enabled: true
    });

    s_usdcTokenPool.setDomains(domainUpdates);

    // Make update invalid on allowedCaller
    domainUpdates[0].allowedCaller = bytes32(0);
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidDomain.selector, domainUpdates[0]));

    s_usdcTokenPool.setDomains(domainUpdates);

    // Make valid again
    domainUpdates[0].allowedCaller = validCaller;

    // Make invalid on destChainSelector
    domainUpdates[0].destChainSelector = 0;
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidDomain.selector, domainUpdates[0]));

    s_usdcTokenPool.setDomains(domainUpdates);
  }
}

contract USDCTokenPool__validateMessage is USDCTokenPoolSetup {
  function test_Fuzz_ValidateMessage_Success(uint32 sourceDomain, uint64 nonce) public {
    vm.pauseGasMetering();
    USDCMessage memory usdcMessage = USDCMessage({
      version: 0,
      sourceDomain: sourceDomain,
      destinationDomain: DEST_DOMAIN_IDENTIFIER,
      nonce: nonce,
      sender: SOURCE_CHAIN_TOKEN_SENDER,
      recipient: bytes32(uint256(299999)),
      destinationCaller: bytes32(uint256(uint160(address(s_usdcTokenPool)))),
      messageBody: bytes("")
    });

    bytes memory encodedUsdcMessage = _generateUSDCMessage(usdcMessage);

    vm.resumeGasMetering();
    s_usdcTokenPool.validateMessage(
      encodedUsdcMessage, USDCTokenPool.SourceTokenDataPayload({nonce: nonce, sourceDomain: sourceDomain})
    );
  }

  // Reverts

  function test_ValidateInvalidMessage_Revert() public {
    USDCMessage memory usdcMessage = USDCMessage({
      version: 0,
      sourceDomain: 1553252,
      destinationDomain: DEST_DOMAIN_IDENTIFIER,
      nonce: 387289284924,
      sender: SOURCE_CHAIN_TOKEN_SENDER,
      recipient: bytes32(uint256(92398429395823)),
      destinationCaller: bytes32(uint256(uint160(address(s_usdcTokenPool)))),
      messageBody: bytes("")
    });

    USDCTokenPool.SourceTokenDataPayload memory sourceTokenData =
      USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: usdcMessage.sourceDomain});

    bytes memory encodedUsdcMessage = _generateUSDCMessage(usdcMessage);

    s_usdcTokenPool.validateMessage(encodedUsdcMessage, sourceTokenData);

    uint32 expectedSourceDomain = usdcMessage.sourceDomain + 1;

    vm.expectRevert(
      abi.encodeWithSelector(USDCTokenPool.InvalidSourceDomain.selector, expectedSourceDomain, usdcMessage.sourceDomain)
    );
    s_usdcTokenPool.validateMessage(
      encodedUsdcMessage,
      USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: expectedSourceDomain})
    );

    uint64 expectedNonce = usdcMessage.nonce + 1;

    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidNonce.selector, expectedNonce, usdcMessage.nonce));
    s_usdcTokenPool.validateMessage(
      encodedUsdcMessage,
      USDCTokenPool.SourceTokenDataPayload({nonce: expectedNonce, sourceDomain: usdcMessage.sourceDomain})
    );

    usdcMessage.destinationDomain = DEST_DOMAIN_IDENTIFIER + 1;
    vm.expectRevert(
      abi.encodeWithSelector(
        USDCTokenPool.InvalidDestinationDomain.selector, DEST_DOMAIN_IDENTIFIER, usdcMessage.destinationDomain
      )
    );

    s_usdcTokenPool.validateMessage(
      _generateUSDCMessage(usdcMessage),
      USDCTokenPool.SourceTokenDataPayload({nonce: usdcMessage.nonce, sourceDomain: usdcMessage.sourceDomain})
    );
    usdcMessage.destinationDomain = DEST_DOMAIN_IDENTIFIER;

    uint32 wrongVersion = usdcMessage.version + 1;

    usdcMessage.version = wrongVersion;
    encodedUsdcMessage = _generateUSDCMessage(usdcMessage);

    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidMessageVersion.selector, wrongVersion));
    s_usdcTokenPool.validateMessage(encodedUsdcMessage, sourceTokenData);
  }
}
