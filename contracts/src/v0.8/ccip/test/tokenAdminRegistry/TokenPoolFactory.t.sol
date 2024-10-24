// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";
import {IOwner} from "../../interfaces/IOwner.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {OwnerIsCreator} from "../../../shared/access/OwnerIsCreator.sol";

import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";

import {RegistryModuleOwnerCustom} from "../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {FactoryBurnMintERC20} from "../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {TokenPoolFactory} from "../../tokenAdminRegistry/TokenPoolFactory/TokenPoolFactory.sol";
import {TokenAdminRegistrySetup} from "./TokenAdminRegistry.t.sol";

import {Create2} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";

contract TokenPoolFactorySetup is TokenAdminRegistrySetup {
  using Create2 for bytes32;

  TokenPoolFactory internal s_tokenPoolFactory;
  RegistryModuleOwnerCustom internal s_registryModuleOwnerCustom;

  bytes internal s_poolInitCode;
  bytes internal s_poolInitArgs;

  bytes32 internal constant FAKE_SALT = keccak256(abi.encode("FAKE_SALT"));

  address internal s_rmnProxy = address(0x1234);

  bytes internal s_tokenCreationParams;
  bytes internal s_tokenInitCode;

  uint256 public constant PREMINT_AMOUNT = 100 ether;

  function setUp() public virtual override {
    TokenAdminRegistrySetup.setUp();

    s_registryModuleOwnerCustom = new RegistryModuleOwnerCustom(address(s_tokenAdminRegistry));
    s_tokenAdminRegistry.addRegistryModule(address(s_registryModuleOwnerCustom));

    s_tokenPoolFactory =
      new TokenPoolFactory(s_tokenAdminRegistry, s_registryModuleOwnerCustom, s_rmnProxy, address(s_sourceRouter));

    // Create Init Code for BurnMintERC20 TestToken with 18 decimals and supply cap of max uint256 value
    s_tokenCreationParams = abi.encode("TestToken", "TT", 18, type(uint256).max, PREMINT_AMOUNT, OWNER);

    s_tokenInitCode = abi.encodePacked(type(FactoryBurnMintERC20).creationCode, s_tokenCreationParams);

    s_poolInitCode = type(BurnMintTokenPool).creationCode;

    // Create Init Args for BurnMintTokenPool with no allowlist minus the token address
    address[] memory allowlist = new address[](1);
    allowlist[0] = OWNER;
    s_poolInitArgs = abi.encode(allowlist, address(0x1234), s_sourceRouter);
  }
}

contract TokenPoolFactoryTests is TokenPoolFactorySetup {
  using Create2 for bytes32;

  function test_TokenPoolFactory_Constructor_Revert() public {
    // Revert cause the tokenAdminRegistry is address(0)
    vm.expectRevert(TokenPoolFactory.InvalidZeroAddress.selector);
    new TokenPoolFactory(ITokenAdminRegistry(address(0)), RegistryModuleOwnerCustom(address(0)), address(0), address(0));

    new TokenPoolFactory(
      ITokenAdminRegistry(address(0xdeadbeef)),
      RegistryModuleOwnerCustom(address(0xdeadbeef)),
      address(0xdeadbeef),
      address(0xdeadbeef)
    );
  }

  function test_createTokenPool_WithNoExistingTokenOnRemoteChain_Success() public {
    vm.startPrank(OWNER);

    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    address predictedTokenAddress =
      Create2.computeAddress(dynamicSalt, keccak256(s_tokenInitCode), address(s_tokenPoolFactory));

    // Create the constructor params for the predicted pool
    bytes memory poolCreationParams = abi.encode(predictedTokenAddress, new address[](0), s_rmnProxy, s_sourceRouter);

    // Predict the address of the pool before we make the tx by using the init code and the params
    bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, poolCreationParams);

    address predictedPoolAddress =
      dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(s_tokenPoolFactory));

    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

    assertNotEq(address(0), tokenAddress, "Token Address should not be 0");
    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    assertEq(predictedTokenAddress, tokenAddress, "Token Address should have been predicted");
    assertEq(predictedPoolAddress, poolAddress, "Pool Address should have been predicted");

    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    OwnerIsCreator(tokenAddress).acceptOwnership();
    OwnerIsCreator(poolAddress).acceptOwnership();

    assertEq(poolAddress, s_tokenAdminRegistry.getPool(tokenAddress), "Token Pool should be set");
    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be owned by the owner");
    assertEq(IOwner(poolAddress).owner(), OWNER, "Token should be owned by the owner");
  }

  function test_createTokenPool_WithNoExistingRemoteContracts_predict_Success() public {
    vm.startPrank(OWNER);
    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    // We have to create a new factory, registry module, and token admin registry to simulate the other chain
    TokenAdminRegistry newTokenAdminRegistry = new TokenAdminRegistry();
    RegistryModuleOwnerCustom newRegistryModule = new RegistryModuleOwnerCustom(address(newTokenAdminRegistry));

    // We want to deploy a new factory and Owner Module.
    TokenPoolFactory newTokenPoolFactory =
      new TokenPoolFactory(newTokenAdminRegistry, newRegistryModule, s_rmnProxy, address(s_destRouter));

    newTokenAdminRegistry.addRegistryModule(address(newRegistryModule));

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig =
      TokenPoolFactory.RemoteChainConfig(address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy));

    // Create an array of remote pools where nothing exists yet, but we want to predict the address for
    // the new pool and token on DEST_CHAIN_SELECTOR
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    // The only field that matters is DEST_CHAIN_SELECTOR because we dont want any existing token pool or token
    // on the remote chain
    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      "", // remotePoolAddress
      type(BurnMintTokenPool).creationCode, // remotePoolInitCode
      remoteChainConfig, // remoteChainConfig
      TokenPoolFactory.PoolType.BURN_MINT, // poolType
      "", // remoteTokenAddress
      s_tokenInitCode, // remoteTokenInitCode
      RateLimiter.Config(false, 0, 0)
    );

    // Predict the address of the token and pool on the DESTINATION chain
    address predictedTokenAddress = dynamicSalt.computeAddress(keccak256(s_tokenInitCode), address(newTokenPoolFactory));

    // Since the remote chain information was provided, we should be able to get the information from the newly
    // deployed token pool using the available getter functions
    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      remoteTokenPools, // No existing remote pools
      s_tokenInitCode, // Token Init Code
      s_poolInitCode, // Pool Init Code
      FAKE_SALT // Salt
    );

    // Ensure that the remote Token was set to the one we predicted
    assertEq(
      abi.encode(predictedTokenAddress),
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      "Token Address should have been predicted"
    );

    {
      // Create the constructor params for the predicted pool
      // The predictedTokenAddress is NOT abi-encoded since the raw evm-address
      // is used in the constructor params
      bytes memory predictedPoolCreationParams =
        abi.encode(predictedTokenAddress, new address[](0), s_rmnProxy, address(s_destRouter));

      // Take the init code and concat the destination params to it, the initCode shouldn't change
      bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, predictedPoolCreationParams);

      // Predict the address of the pool on the DESTINATION chain
      address predictedPoolAddress =
        dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(newTokenPoolFactory));

      // Assert that the address set for the remote pool is the same as the predicted address
      assertEq(
        abi.encode(predictedPoolAddress),
        TokenPool(poolAddress).getRemotePool(DEST_CHAIN_SELECTOR),
        "Pool Address should have been predicted"
      );
    }

    // On the new token pool factory, representing a destination chain,
    // deploy a new token and a new pool
    (address newTokenAddress, address newPoolAddress) = newTokenPoolFactory.deployTokenAndTokenPool(
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

    assertEq(
      TokenPool(poolAddress).getRemotePool(DEST_CHAIN_SELECTOR),
      abi.encode(newPoolAddress),
      "New Pool Address should have been deployed correctly"
    );

    assertEq(
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      abi.encode(newTokenAddress),
      "New Token Address should have been deployed correctly"
    );

    // Check that the token pool has the correct permissions
    vm.startPrank(poolAddress);
    IBurnMintERC20(tokenAddress).mint(poolAddress, 1e18);

    assertEq(1e18, IBurnMintERC20(tokenAddress).balanceOf(poolAddress), "Balance should be 1e18");

    IBurnMintERC20(tokenAddress).burn(1e18);
    assertEq(0, IBurnMintERC20(tokenAddress).balanceOf(poolAddress), "Balance should be 0");

    vm.stopPrank();

    assertEq(s_tokenAdminRegistry.getPool(tokenAddress), poolAddress, "Token Pool should be set");

    // Check the token admin registry for config
    TokenAdminRegistry.TokenConfig memory tokenConfig = s_tokenAdminRegistry.getTokenConfig(tokenAddress);
    assertEq(tokenConfig.administrator, address(s_tokenPoolFactory), "Administrator should be set");
    assertEq(tokenConfig.pendingAdministrator, OWNER, "Pending Administrator should be 0");
    assertEq(tokenConfig.tokenPool, poolAddress, "Pool Address should be set");

    // Accept Ownership of the token, pool, and adminRegistry
    vm.startPrank(OWNER);
    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    assertEq(s_tokenAdminRegistry.getTokenConfig(tokenAddress).administrator, OWNER, "Administrator should be set");
    assertEq(
      s_tokenAdminRegistry.getTokenConfig(tokenAddress).pendingAdministrator, address(0), "Administrator should be set"
    );

    OwnerIsCreator(tokenAddress).acceptOwnership();
    OwnerIsCreator(poolAddress).acceptOwnership();

    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be controlled by the OWNER");
    assertEq(IOwner(poolAddress).owner(), OWNER, "Pool should be controlled by the OWNER");
  }

  function test_createTokenPool_ExistingRemoteToken_AndPredictPool_Success() public {
    vm.startPrank(OWNER);
    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    FactoryBurnMintERC20 newRemoteToken =
      new FactoryBurnMintERC20("TestToken", "TT", 18, type(uint256).max, PREMINT_AMOUNT, OWNER);

    // We have to create a new factory, registry module, and token admin registry to simulate the other chain
    TokenAdminRegistry newTokenAdminRegistry = new TokenAdminRegistry();
    RegistryModuleOwnerCustom newRegistryModule = new RegistryModuleOwnerCustom(address(newTokenAdminRegistry));

    // We want to deploy a new factory and Owner Module.
    TokenPoolFactory newTokenPoolFactory =
      new TokenPoolFactory(newTokenAdminRegistry, newRegistryModule, s_rmnProxy, address(s_destRouter));

    newTokenAdminRegistry.addRegistryModule(address(newRegistryModule));

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig =
      TokenPoolFactory.RemoteChainConfig(address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy));

    // Create an array of remote pools where nothing exists yet, but we want to predict the address for
    // the new pool and token on DEST_CHAIN_SELECTOR
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    // The only field that matters is DEST_CHAIN_SELECTOR because we dont want any existing token pool or token
    // on the remote chain
    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      "", // remotePoolAddress
      type(BurnMintTokenPool).creationCode, // remotePoolInitCode
      remoteChainConfig, // remoteChainConfig
      TokenPoolFactory.PoolType.BURN_MINT, // poolType
      abi.encode(address(newRemoteToken)), // remoteTokenAddress
      s_tokenInitCode, // remoteTokenInitCode
      RateLimiter.Config(false, 0, 0) // rateLimiterConfig
    );

    // Since the remote chain information was provided, we should be able to get the information from the newly
    // deployed token pool using the available getter functions
    (address tokenAddress, address poolAddress) =
      s_tokenPoolFactory.deployTokenAndTokenPool(remoteTokenPools, s_tokenInitCode, s_poolInitCode, FAKE_SALT);

    assertEq(address(TokenPool(poolAddress).getToken()), tokenAddress, "Token Address should have been set locally");

    // Ensure that the remote Token was set to the one we predicted
    assertEq(
      abi.encode(address(newRemoteToken)),
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      "Token Address should have been predicted"
    );

    // Create the constructor params for the predicted pool
    // The predictedTokenAddress is NOT abi-encoded since the raw evm-address
    // is used in the constructor params
    bytes memory predictedPoolCreationParams =
      abi.encode(address(newRemoteToken), new address[](0), s_rmnProxy, address(s_destRouter));

    // Take the init code and concat the destination params to it, the initCode shouldn't change
    bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, predictedPoolCreationParams);

    // Predict the address of the pool on the DESTINATION chain
    address predictedPoolAddress =
      dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(newTokenPoolFactory));

    // Assert that the address set for the remote pool is the same as the predicted address
    assertEq(
      abi.encode(predictedPoolAddress),
      TokenPool(poolAddress).getRemotePool(DEST_CHAIN_SELECTOR),
      "Pool Address should have been predicted"
    );

    // On the new token pool factory, representing a destination chain,
    // deploy a new token and a new pool
    address newPoolAddress = newTokenPoolFactory.deployTokenPoolWithExistingToken(
      address(newRemoteToken),
      new TokenPoolFactory.RemoteTokenPoolInfo[](0),
      s_poolInitCode,
      FAKE_SALT,
      TokenPoolFactory.PoolType.BURN_MINT
    );

    assertEq(
      abi.encode(newRemoteToken),
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      "Remote Token Address should have been set correctly"
    );

    assertEq(
      TokenPool(poolAddress).getRemotePool(DEST_CHAIN_SELECTOR),
      abi.encode(newPoolAddress),
      "New Pool Address should have been deployed correctly"
    );
  }

  function test_createTokenPool_WithRemoteTokenAndRemotePool_Success() public {
    vm.startPrank(OWNER);

    bytes memory RANDOM_TOKEN_ADDRESS = abi.encode(makeAddr("RANDOM_TOKEN"));
    bytes memory RANDOM_POOL_ADDRESS = abi.encode(makeAddr("RANDOM_POOL"));

    // Create an array of remote pools with some fake addresses
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      RANDOM_POOL_ADDRESS, // remotePoolAddress
      type(BurnMintTokenPool).creationCode, // remotePoolInitCode
      TokenPoolFactory.RemoteChainConfig(address(0), address(0), address(0)), // remoteChainConfig
      TokenPoolFactory.PoolType.BURN_MINT, // poolType
      RANDOM_TOKEN_ADDRESS, // remoteTokenAddress
      "", // remoteTokenInitCode
      RateLimiter.Config(false, 0, 0) // rateLimiterConfig
    );

    (address tokenAddress, address poolAddress) =
      s_tokenPoolFactory.deployTokenAndTokenPool(remoteTokenPools, s_tokenInitCode, s_poolInitCode, FAKE_SALT);

    assertNotEq(address(0), tokenAddress, "Token Address should not be 0");
    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    OwnerIsCreator(tokenAddress).acceptOwnership();
    OwnerIsCreator(poolAddress).acceptOwnership();

    assertEq(
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      RANDOM_TOKEN_ADDRESS,
      "Remote Token Address should have been set"
    );

    assertEq(
      TokenPool(poolAddress).getRemotePool(DEST_CHAIN_SELECTOR),
      RANDOM_POOL_ADDRESS,
      "Remote Pool Address should have been set"
    );

    assertEq(poolAddress, s_tokenAdminRegistry.getPool(tokenAddress), "Token Pool should be set");

    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be owned by the owner");

    assertEq(IOwner(poolAddress).owner(), OWNER, "Token should be owned by the owner");
  }

  function test_createTokenPoolLockRelease_ExistingToken_predict_Success() public {
    vm.startPrank(OWNER);

    // We have to create a new factory, registry module, and token admin registry to simulate the other chain
    TokenAdminRegistry newTokenAdminRegistry = new TokenAdminRegistry();
    RegistryModuleOwnerCustom newRegistryModule = new RegistryModuleOwnerCustom(address(newTokenAdminRegistry));

    // We want to deploy a new factory and Owner Module.
    TokenPoolFactory newTokenPoolFactory =
      new TokenPoolFactory(newTokenAdminRegistry, newRegistryModule, s_rmnProxy, address(s_destRouter));

    newTokenAdminRegistry.addRegistryModule(address(newRegistryModule));

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig =
      TokenPoolFactory.RemoteChainConfig(address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy));

    FactoryBurnMintERC20 newLocalToken =
      new FactoryBurnMintERC20("TestToken", "TEST", 18, type(uint256).max, PREMINT_AMOUNT, OWNER);

    FactoryBurnMintERC20 newRemoteToken =
      new FactoryBurnMintERC20("TestToken", "TEST", 18, type(uint256).max, PREMINT_AMOUNT, OWNER);

    // Create an array of remote pools where nothing exists yet, but we want to predict the address for
    // the new pool and token on DEST_CHAIN_SELECTOR
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    // The only field that matters is DEST_CHAIN_SELECTOR because we dont want any existing token pool or token
    // on the remote chain
    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      "", // remotePoolAddress
      type(LockReleaseTokenPool).creationCode, // remotePoolInitCode
      remoteChainConfig, // remoteChainConfig
      TokenPoolFactory.PoolType.LOCK_RELEASE, // poolType
      abi.encode(address(newRemoteToken)), // remoteTokenAddress
      s_tokenInitCode, // remoteTokenInitCode
      RateLimiter.Config(false, 0, 0)
    );

    // Since the remote chain information was provided, we should be able to get the information from the newly
    // deployed token pool using the available getter functions
    address poolAddress = s_tokenPoolFactory.deployTokenPoolWithExistingToken(
      address(newLocalToken),
      remoteTokenPools,
      type(LockReleaseTokenPool).creationCode,
      FAKE_SALT,
      TokenPoolFactory.PoolType.LOCK_RELEASE
    );

    // Check that the pool was correctly deployed on the local chain first

    // Accept the ownership which was transfered
    OwnerIsCreator(poolAddress).acceptOwnership();

    // Ensure that the remote Token was set to the one we predicted
    assertEq(
      address(LockReleaseTokenPool(poolAddress).getToken()),
      address(newLocalToken),
      "Token Address should have been set"
    );

    LockReleaseTokenPool(poolAddress).setRebalancer(OWNER);
    assertEq(OWNER, LockReleaseTokenPool(poolAddress).getRebalancer(), "Rebalancer should be set");

    // Deploy the Lock-Release Token Pool on the destination chain with the existing remote token
    (address newPoolAddress) = newTokenPoolFactory.deployTokenPoolWithExistingToken(
      address(newRemoteToken),
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), // No existing remote pools
      type(LockReleaseTokenPool).creationCode, // Pool Init Code
      FAKE_SALT, // Salt
      TokenPoolFactory.PoolType.LOCK_RELEASE
    );

    assertEq(
      LockReleaseTokenPool(poolAddress).getRemotePool(DEST_CHAIN_SELECTOR),
      abi.encode(newPoolAddress),
      "New Pool Address should have been deployed correctly"
    );

    assertEq(
      LockReleaseTokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      abi.encode(address(newRemoteToken)),
      "New Token Address should have been deployed correctly"
    );

    assertEq(
      address(LockReleaseTokenPool(newPoolAddress).getToken()),
      address(newRemoteToken),
      "New Remote Token should be set correctly"
    );
  }
}
