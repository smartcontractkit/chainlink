// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import "../BaseTest.t.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {Router} from "../../Router.sol";
import {USDCTokenPool} from "../../pools/USDC/USDCTokenPool.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {MockUSDC} from "../mocks/MockUSDC.sol";

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.0/utils/introspection/IERC165.sol";

contract USDCTokenPoolSetup is BaseTest {
  IBurnMintERC20 internal s_token;
  MockUSDC internal s_mockUSDC;

  address internal s_routerAllowedOnRamp = address(3456);
  address internal s_routerAllowedOffRamp = address(234);
  Router internal s_router;

  USDCTokenPool internal s_usdcTokenPool;
  USDCTokenPool internal s_usdcTokenPoolWithAllowList;
  address[] internal s_allowedList;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_token = new BurnMintERC677("LINK", "LNK", 18, 0);
    deal(address(s_token), OWNER, type(uint256).max);
    setUpRamps();

    s_mockUSDC = new MockUSDC(42);

    USDCTokenPool.USDCConfig memory config = USDCTokenPool.USDCConfig({
      version: s_mockUSDC.messageBodyVersion(),
      tokenMessenger: address(s_mockUSDC),
      messageTransmitter: address(s_mockUSDC)
    });

    s_usdcTokenPool = new USDCTokenPool(config, s_token, new address[](0), address(s_mockARM));

    s_allowedList.push(USER_1);
    s_usdcTokenPoolWithAllowList = new USDCTokenPool(config, s_token, s_allowedList, address(s_mockARM));

    TokenPool.RampUpdate[] memory onRamps = new TokenPool.RampUpdate[](1);
    onRamps[0] = TokenPool.RampUpdate({
      ramp: s_routerAllowedOnRamp,
      allowed: true,
      rateLimiterConfig: rateLimiterConfig()
    });

    TokenPool.RampUpdate[] memory offRamps = new TokenPool.RampUpdate[](1);
    offRamps[0] = TokenPool.RampUpdate({
      ramp: s_routerAllowedOffRamp,
      allowed: true,
      rateLimiterConfig: rateLimiterConfig()
    });

    s_usdcTokenPool.applyRampUpdates(onRamps, offRamps);
    s_usdcTokenPoolWithAllowList.applyRampUpdates(onRamps, offRamps);

    USDCTokenPool.DomainUpdate[] memory domains = new USDCTokenPool.DomainUpdate[](1);
    domains[0] = USDCTokenPool.DomainUpdate({
      destChainSelector: DEST_CHAIN_ID,
      domainIdentifier: 9999,
      allowedCaller: keccak256("allowedCaller"),
      enabled: true
    });

    s_usdcTokenPool.setDomains(domains);
    s_usdcTokenPoolWithAllowList.setDomains(domains);
  }

  function setUpRamps() internal {
    s_router = new Router(address(s_token), address(s_mockARM));

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_ID, onRamp: s_routerAllowedOnRamp});
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    address[] memory offRamps = new address[](1);
    offRamps[0] = s_routerAllowedOffRamp;
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_ID, offRamp: offRamps[0]});

    s_router.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }
}

contract USDCTokenPool_lockOrBurn is USDCTokenPoolSetup {
  error SenderNotAllowed(address sender);
  event DepositForBurn(
    uint64 indexed nonce,
    address indexed burnToken,
    uint256 amount,
    address indexed depositor,
    bytes32 mintRecipient,
    uint32 destinationDomain,
    bytes32 destinationTokenMessenger,
    bytes32 destinationCaller
  );
  event Burned(address indexed sender, uint256 amount);
  event TokensConsumed(uint256 tokens);

  function testFuzz_LockOrBurnSuccess(bytes32 destinationReceiver, uint256 amount) public {
    vm.assume(amount < rateLimiterConfig().capacity);
    vm.assume(amount > 0);
    changePrank(s_routerAllowedOnRamp);
    s_token.approve(address(s_usdcTokenPool), amount);

    USDCTokenPool.Domain memory expectedDomain = s_usdcTokenPool.getDomain(DEST_CHAIN_ID);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit DepositForBurn(
      s_mockUSDC.s_nonce(),
      address(s_token),
      amount,
      address(s_usdcTokenPool),
      destinationReceiver,
      expectedDomain.domainIdentifier,
      s_mockUSDC.i_destinationTokenMessenger(),
      expectedDomain.allowedCaller
    );
    vm.expectEmit();
    emit Burned(s_routerAllowedOnRamp, amount);

    bytes memory encodedNonce = s_usdcTokenPool.lockOrBurn(
      OWNER,
      abi.encodePacked(destinationReceiver),
      amount,
      DEST_CHAIN_ID,
      bytes("")
    );
    uint64 nonce = abi.decode(encodedNonce, (uint64));
    assertEq(s_mockUSDC.s_nonce() - 1, nonce);
  }

  function testFuzz_LockOrBurnWithAllowListSuccess(bytes32 destinationReceiver, uint256 amount) public {
    vm.assume(amount < rateLimiterConfig().capacity);
    vm.assume(amount > 0);
    changePrank(s_routerAllowedOnRamp);
    s_token.approve(address(s_usdcTokenPoolWithAllowList), amount);

    USDCTokenPool.Domain memory expectedDomain = s_usdcTokenPoolWithAllowList.getDomain(DEST_CHAIN_ID);

    vm.expectEmit();
    emit TokensConsumed(amount);
    vm.expectEmit();
    emit DepositForBurn(
      s_mockUSDC.s_nonce(),
      address(s_token),
      amount,
      address(s_usdcTokenPoolWithAllowList),
      destinationReceiver,
      expectedDomain.domainIdentifier,
      s_mockUSDC.i_destinationTokenMessenger(),
      expectedDomain.allowedCaller
    );
    vm.expectEmit();
    emit Burned(s_routerAllowedOnRamp, amount);

    bytes memory encodedNonce = s_usdcTokenPoolWithAllowList.lockOrBurn(
      s_allowedList[0],
      abi.encodePacked(destinationReceiver),
      amount,
      DEST_CHAIN_ID,
      bytes("")
    );
    uint64 nonce = abi.decode(encodedNonce, (uint64));
    assertEq(s_mockUSDC.s_nonce() - 1, nonce);
  }

  // Reverts
  function testUnknownDomainReverts() public {
    uint256 amount = 1000;
    changePrank(s_routerAllowedOnRamp);
    deal(address(s_token), s_routerAllowedOnRamp, amount);
    s_token.approve(address(s_usdcTokenPool), amount);

    uint64 wrongDomain = DEST_CHAIN_ID + 1;

    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.UnknownDomain.selector, wrongDomain));

    s_usdcTokenPool.lockOrBurn(OWNER, abi.encodePacked(address(0)), amount, wrongDomain, bytes(""));
  }

  function testPermissionsErrorReverts() public {
    vm.expectRevert(TokenPool.PermissionsError.selector);

    s_usdcTokenPool.lockOrBurn(OWNER, abi.encodePacked(address(0)), 0, DEST_CHAIN_ID, bytes(""));
  }

  function testLockOrBurnWithAllowListReverts() public {
    changePrank(s_routerAllowedOnRamp);

    vm.expectRevert(abi.encodeWithSelector(SenderNotAllowed.selector, STRANGER));

    s_usdcTokenPoolWithAllowList.lockOrBurn(STRANGER, abi.encodePacked(address(0)), 1000, DEST_CHAIN_ID, bytes(""));
  }
}

contract USDCTokenPool_releaseOrMint is USDCTokenPoolSetup {
  event Minted(address indexed sender, address indexed recipient, uint256 amount);

  function testFuzz_ReleaseOrMintSuccess(address receiver, uint256 amount) public {
    amount = bound(amount, 0, rateLimiterConfig().capacity);
    changePrank(s_routerAllowedOffRamp);

    bytes memory message = bytes("message bytes");
    bytes memory attestation = bytes("attestation bytes");

    bytes memory extraData = abi.encode(
      USDCTokenPool.MessageAndAttestation({message: message, attestation: attestation})
    );

    vm.expectEmit();
    emit Minted(s_routerAllowedOffRamp, receiver, amount);

    vm.expectCall(address(s_mockUSDC), abi.encodeWithSelector(MockUSDC.receiveMessage.selector, message, attestation));

    s_usdcTokenPool.releaseOrMint(abi.encode(OWNER), receiver, amount, SOURCE_CHAIN_ID, extraData);
  }

  // Reverts
  function testUnlockingUSDCFailedReverts() public {
    changePrank(s_routerAllowedOffRamp);
    s_mockUSDC.setShouldSucceed(false);

    bytes memory extraData = abi.encode(
      USDCTokenPool.MessageAndAttestation({message: bytes(""), attestation: bytes("")})
    );

    vm.expectRevert(USDCTokenPool.UnlockingUSDCFailed.selector);

    s_usdcTokenPool.releaseOrMint(abi.encode(OWNER), OWNER, 1, SOURCE_CHAIN_ID, extraData);
  }

  function testTokenMaxCapacityExceededReverts() public {
    uint256 capacity = rateLimiterConfig().capacity;
    uint256 amount = 10 * capacity;
    address receiver = address(1);
    changePrank(s_routerAllowedOffRamp);

    bytes memory extraData = abi.encode(
      USDCTokenPool.MessageAndAttestation({message: bytes(""), attestation: bytes("")})
    );

    vm.expectRevert(
      abi.encodeWithSelector(RateLimiter.TokenMaxCapacityExceeded.selector, capacity, amount, address(s_token))
    );

    s_usdcTokenPool.releaseOrMint(abi.encode(OWNER), receiver, amount, SOURCE_CHAIN_ID, extraData);
  }
}

contract USDCTokenPool_supportsInterface is USDCTokenPoolSetup {
  function testSupportsInterfaceSuccess() public {
    assertTrue(s_usdcTokenPool.supportsInterface(s_usdcTokenPool.getUSDCInterfaceId()));
    assertTrue(s_usdcTokenPool.supportsInterface(type(IPool).interfaceId));
    assertTrue(s_usdcTokenPool.supportsInterface(type(IERC165).interfaceId));
  }
}

contract USDCTokenPool_setDomains is USDCTokenPoolSetup {
  event DomainsSet(USDCTokenPool.DomainUpdate[]);

  mapping(uint64 destChainSelector => USDCTokenPool.Domain domain) private s_chainToDomain;

  // Setting lower fuzz run as 256 runs was causing differing gas results in snapshot.
  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function testFuzz_SetDomainsSuccess(
    bytes32[10] calldata allowedCallers,
    uint32[10] calldata domainIdentifiers,
    uint64[10] calldata destChainSelectors
  ) public {
    uint256 numberOfDomains = allowedCallers.length;
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](numberOfDomains);
    for (uint256 i = 0; i < numberOfDomains; ++i) {
      domainUpdates[i] = USDCTokenPool.DomainUpdate({
        allowedCaller: allowedCallers[i],
        domainIdentifier: domainIdentifiers[i],
        destChainSelector: destChainSelectors[i],
        enabled: true
      });

      s_chainToDomain[destChainSelectors[i]] = USDCTokenPool.Domain({
        domainIdentifier: domainIdentifiers[i],
        allowedCaller: allowedCallers[i],
        enabled: true
      });
    }

    vm.expectEmit();
    emit DomainsSet(domainUpdates);

    s_usdcTokenPool.setDomains(domainUpdates);

    for (uint256 i = 0; i < numberOfDomains; ++i) {
      USDCTokenPool.Domain memory expected = s_chainToDomain[destChainSelectors[i]];
      USDCTokenPool.Domain memory got = s_usdcTokenPool.getDomain(destChainSelectors[i]);
      assertEq(got.allowedCaller, expected.allowedCaller);
      assertEq(got.domainIdentifier, expected.domainIdentifier);
    }
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    USDCTokenPool.DomainUpdate[] memory domainUpdates = new USDCTokenPool.DomainUpdate[](0);

    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");

    s_usdcTokenPool.setDomains(domainUpdates);
  }
}

contract USDCTokenPool_setConfig is USDCTokenPoolSetup {
  event ConfigSet(USDCTokenPool.USDCConfig);

  function testSetConfigSuccess() public {
    USDCTokenPool.USDCConfig memory newConfig = USDCTokenPool.USDCConfig({
      version: 12332,
      tokenMessenger: address(100),
      messageTransmitter: address(123456789)
    });

    USDCTokenPool.USDCConfig memory oldConfig = s_usdcTokenPool.getConfig();

    vm.expectEmit();
    emit ConfigSet(newConfig);
    s_usdcTokenPool.setConfig(newConfig);

    USDCTokenPool.USDCConfig memory gotConfig = s_usdcTokenPool.getConfig();
    assertEq(gotConfig.tokenMessenger, newConfig.tokenMessenger);
    assertEq(gotConfig.messageTransmitter, newConfig.messageTransmitter);
    assertEq(gotConfig.version, newConfig.version);

    assertEq(0, s_usdcTokenPool.getToken().allowance(address(s_usdcTokenPool), oldConfig.tokenMessenger));
    assertEq(
      type(uint256).max,
      s_usdcTokenPool.getToken().allowance(address(s_usdcTokenPool), gotConfig.tokenMessenger)
    );
  }

  // Reverts

  function testInvalidConfigReverts() public {
    USDCTokenPool.USDCConfig memory newConfig = USDCTokenPool.USDCConfig({
      version: 12332,
      tokenMessenger: address(0),
      messageTransmitter: address(123456789)
    });

    vm.expectRevert(USDCTokenPool.InvalidConfig.selector);
    s_usdcTokenPool.setConfig(newConfig);

    newConfig.tokenMessenger = address(235);
    newConfig.messageTransmitter = address(0);

    vm.expectRevert(USDCTokenPool.InvalidConfig.selector);
    s_usdcTokenPool.setConfig(newConfig);
  }

  function testOnlyOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");

    s_usdcTokenPool.setConfig(
      USDCTokenPool.USDCConfig({version: 1, tokenMessenger: address(100), messageTransmitter: address(1)})
    );
  }
}
