// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IBurnMintERC20} from "../../../../shared/token/ERC20/IBurnMintERC20.sol";
import {IOwner} from "../../../interfaces/IOwner.sol";

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";

import {Router} from "../../../Router.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {BurnFromMintTokenPool} from "../../../pools/BurnFromMintTokenPool.sol";
import {BurnMintTokenPool} from "../../../pools/BurnMintTokenPool.sol";
import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {RegistryModuleOwnerCustom} from "../../../tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {TokenPoolFactory} from "../../../tokenAdminRegistry/TokenPoolFactory/TokenPoolFactory.sol";
import {TokenPoolFactorySetup} from "./TokenPoolFactorySetup.t.sol";

import {IERC20Metadata} from
  "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {Create2} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";

contract TokenPoolFactory_createTokenPool is TokenPoolFactorySetup {
  using Create2 for bytes32;

  uint8 private constant LOCAL_TOKEN_DECIMALS = 18;
  uint8 private constant REMOTE_TOKEN_DECIMALS = 6;

  address internal s_burnMintOffRamp = makeAddr("burn_mint_offRamp");

  function setUp() public override {
    TokenPoolFactorySetup.setUp();

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: DEST_CHAIN_SELECTOR, offRamp: s_burnMintOffRamp});
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRampUpdates);
  }

  function test_createTokenPool_WithNoExistingTokenOnRemoteChain() public {
    vm.startPrank(OWNER);

    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    address predictedTokenAddress =
      Create2.computeAddress(dynamicSalt, keccak256(s_tokenInitCode), address(s_tokenPoolFactory));

    // Create the constructor params for the predicted pool
    bytes memory poolCreationParams =
      abi.encode(predictedTokenAddress, LOCAL_TOKEN_DECIMALS, new address[](0), s_rmnProxy, s_sourceRouter);

    // Predict the address of the pool before we make the tx by using the init code and the params
    bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, poolCreationParams);

    address predictedPoolAddress =
      dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(s_tokenPoolFactory));

    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), LOCAL_TOKEN_DECIMALS, s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

    assertNotEq(address(0), tokenAddress, "Token Address should not be 0");
    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    assertEq(predictedTokenAddress, tokenAddress, "Token Address should have been predicted");
    assertEq(predictedPoolAddress, poolAddress, "Pool Address should have been predicted");

    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    Ownable2Step(tokenAddress).acceptOwnership();
    Ownable2Step(poolAddress).acceptOwnership();

    assertEq(poolAddress, s_tokenAdminRegistry.getPool(tokenAddress), "Token Pool should be set");
    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be owned by the owner");
    assertEq(IOwner(poolAddress).owner(), OWNER, "Token should be owned by the owner");
  }

  function test_createTokenPool_WithNoExistingRemoteContracts_predict() public {
    vm.startPrank(OWNER);
    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    // We have to create a new factory, registry module, and token admin registry to simulate the other chain
    TokenAdminRegistry newTokenAdminRegistry = new TokenAdminRegistry();
    RegistryModuleOwnerCustom newRegistryModule = new RegistryModuleOwnerCustom(address(newTokenAdminRegistry));

    // We want to deploy a new factory and Owner Module.
    TokenPoolFactory newTokenPoolFactory =
      new TokenPoolFactory(newTokenAdminRegistry, newRegistryModule, s_rmnProxy, address(s_destRouter));

    newTokenAdminRegistry.addRegistryModule(address(newRegistryModule));

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig = TokenPoolFactory.RemoteChainConfig(
      address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy), LOCAL_TOKEN_DECIMALS
    );

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
      LOCAL_TOKEN_DECIMALS, // 18 decimal token
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
        abi.encode(predictedTokenAddress, LOCAL_TOKEN_DECIMALS, new address[](0), s_rmnProxy, address(s_destRouter));

      // Take the init code and concat the destination params to it, the initCode shouldn't change
      bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, predictedPoolCreationParams);

      // Predict the address of the pool on the DESTINATION chain
      address predictedPoolAddress =
        dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(newTokenPoolFactory));

      // Assert that the address set for the remote pool is the same as the predicted address
      assertEq(
        abi.encode(predictedPoolAddress),
        TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
        "Pool Address should have been predicted"
      );
    }

    // On the new token pool factory, representing a destination chain,
    // deploy a new token and a new pool
    (address newTokenAddress, address newPoolAddress) = newTokenPoolFactory.deployTokenAndTokenPool(
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), LOCAL_TOKEN_DECIMALS, s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

    assertEq(
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
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

    Ownable2Step(tokenAddress).acceptOwnership();
    Ownable2Step(poolAddress).acceptOwnership();

    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be controlled by the OWNER");
    assertEq(IOwner(poolAddress).owner(), OWNER, "Pool should be controlled by the OWNER");
  }

  function test_createTokenPool_ExistingRemoteToken_AndPredictPool() public {
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

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig = TokenPoolFactory.RemoteChainConfig(
      address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy), LOCAL_TOKEN_DECIMALS
    );

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
    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      remoteTokenPools, LOCAL_TOKEN_DECIMALS, s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

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
      abi.encode(address(newRemoteToken), LOCAL_TOKEN_DECIMALS, new address[](0), s_rmnProxy, address(s_destRouter));

    // Take the init code and concat the destination params to it, the initCode shouldn't change
    bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, predictedPoolCreationParams);

    // Predict the address of the pool on the DESTINATION chain
    address predictedPoolAddress =
      dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(newTokenPoolFactory));

    // Assert that the address set for the remote pool is the same as the predicted address
    assertEq(
      abi.encode(predictedPoolAddress),
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      "Pool Address should have been predicted"
    );

    // On the new token pool factory, representing a destination chain,
    // deploy a new token and a new pool
    address newPoolAddress = newTokenPoolFactory.deployTokenPoolWithExistingToken(
      address(newRemoteToken),
      LOCAL_TOKEN_DECIMALS,
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
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      abi.encode(newPoolAddress),
      "New Pool Address should have been deployed correctly"
    );
  }

  function test_createTokenPool_WithRemoteTokenAndRemotePool() public {
    vm.startPrank(OWNER);

    bytes memory RANDOM_TOKEN_ADDRESS = abi.encode(makeAddr("RANDOM_TOKEN"));
    bytes memory RANDOM_POOL_ADDRESS = abi.encode(makeAddr("RANDOM_POOL"));

    // Create an array of remote pools with some fake addresses
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      RANDOM_POOL_ADDRESS, // remotePoolAddress
      type(BurnMintTokenPool).creationCode, // remotePoolInitCode
      TokenPoolFactory.RemoteChainConfig(address(0), address(0), address(0), 0), // remoteChainConfig
      TokenPoolFactory.PoolType.BURN_MINT, // poolType
      RANDOM_TOKEN_ADDRESS, // remoteTokenAddress
      "", // remoteTokenInitCode
      RateLimiter.Config(false, 0, 0) // rateLimiterConfig
    );

    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      remoteTokenPools, LOCAL_TOKEN_DECIMALS, s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

    assertNotEq(address(0), tokenAddress, "Token Address should not be 0");
    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    Ownable2Step(tokenAddress).acceptOwnership();
    Ownable2Step(poolAddress).acceptOwnership();

    assertEq(
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      RANDOM_TOKEN_ADDRESS,
      "Remote Token Address should have been set"
    );

    assertEq(
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      RANDOM_POOL_ADDRESS,
      "Remote Pool Address should have been set"
    );

    assertEq(poolAddress, s_tokenAdminRegistry.getPool(tokenAddress), "Token Pool should be set");

    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be owned by the owner");

    assertEq(IOwner(poolAddress).owner(), OWNER, "Token should be owned by the owner");
  }

  function test_createTokenPoolLockRelease_ExistingToken_predict() public {
    vm.startPrank(OWNER);

    // We have to create a new factory, registry module, and token admin registry to simulate the other chain
    TokenAdminRegistry newTokenAdminRegistry = new TokenAdminRegistry();
    RegistryModuleOwnerCustom newRegistryModule = new RegistryModuleOwnerCustom(address(newTokenAdminRegistry));

    // We want to deploy a new factory and Owner Module.
    TokenPoolFactory newTokenPoolFactory =
      new TokenPoolFactory(newTokenAdminRegistry, newRegistryModule, s_rmnProxy, address(s_destRouter));

    newTokenAdminRegistry.addRegistryModule(address(newRegistryModule));

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig = TokenPoolFactory.RemoteChainConfig(
      address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy), LOCAL_TOKEN_DECIMALS
    );

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
      LOCAL_TOKEN_DECIMALS,
      remoteTokenPools,
      type(LockReleaseTokenPool).creationCode,
      FAKE_SALT,
      TokenPoolFactory.PoolType.LOCK_RELEASE
    );

    // Check that the pool was correctly deployed on the local chain first

    // Accept the ownership which was transferred
    Ownable2Step(poolAddress).acceptOwnership();

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
      LOCAL_TOKEN_DECIMALS,
      new TokenPoolFactory.RemoteTokenPoolInfo[](0), // No existing remote pools
      type(LockReleaseTokenPool).creationCode, // Pool Init Code
      FAKE_SALT, // Salt
      TokenPoolFactory.PoolType.LOCK_RELEASE
    );

    assertEq(
      LockReleaseTokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
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

  function test_createTokenPool_BurnFromMintTokenPool() public {
    vm.startPrank(OWNER);

    bytes memory RANDOM_TOKEN_ADDRESS = abi.encode(makeAddr("RANDOM_TOKEN"));
    bytes memory RANDOM_POOL_ADDRESS = abi.encode(makeAddr("RANDOM_POOL"));

    // Create an array of remote pools with some fake addresses
    TokenPoolFactory.RemoteTokenPoolInfo[] memory remoteTokenPools = new TokenPoolFactory.RemoteTokenPoolInfo[](1);

    remoteTokenPools[0] = TokenPoolFactory.RemoteTokenPoolInfo(
      DEST_CHAIN_SELECTOR, // remoteChainSelector
      RANDOM_POOL_ADDRESS, // remotePoolAddress
      type(BurnFromMintTokenPool).creationCode, // remotePoolInitCode
      TokenPoolFactory.RemoteChainConfig(address(0), address(0), address(0), 0), // remoteChainConfig
      TokenPoolFactory.PoolType.BURN_MINT, // poolType
      RANDOM_TOKEN_ADDRESS, // remoteTokenAddress
      "", // remoteTokenInitCode
      RateLimiter.Config(false, 0, 0) // rateLimiterConfig
    );

    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      remoteTokenPools, LOCAL_TOKEN_DECIMALS, s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

    assertNotEq(address(0), tokenAddress, "Token Address should not be 0");
    assertNotEq(address(0), poolAddress, "Pool Address should not be 0");

    s_tokenAdminRegistry.acceptAdminRole(tokenAddress);
    Ownable2Step(tokenAddress).acceptOwnership();
    Ownable2Step(poolAddress).acceptOwnership();

    assertEq(
      TokenPool(poolAddress).getRemoteToken(DEST_CHAIN_SELECTOR),
      RANDOM_TOKEN_ADDRESS,
      "Remote Token Address should have been set"
    );

    assertEq(
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      RANDOM_POOL_ADDRESS,
      "Remote Pool Address should have been set"
    );

    assertEq(poolAddress, s_tokenAdminRegistry.getPool(tokenAddress), "Token Pool should be set");

    assertEq(IOwner(tokenAddress).owner(), OWNER, "Token should be owned by the owner");

    assertEq(IOwner(poolAddress).owner(), OWNER, "Token should be owned by the owner");
  }

  function test_createTokenPool_RemoteTokenHasDifferentDecimals() public {
    vm.startPrank(OWNER);
    bytes32 dynamicSalt = keccak256(abi.encodePacked(FAKE_SALT, OWNER));

    // Deploy the "remote" token which has a different decimal value than the local token
    FactoryBurnMintERC20 newRemoteToken =
      new FactoryBurnMintERC20("TestToken", "TT", 6, type(uint256).max, PREMINT_AMOUNT, OWNER);

    // We have to create a new factory, registry module, and token admin registry to simulate the other chain
    TokenAdminRegistry newTokenAdminRegistry = new TokenAdminRegistry();
    RegistryModuleOwnerCustom newRegistryModule = new RegistryModuleOwnerCustom(address(newTokenAdminRegistry));

    // We want to deploy a new factory and Owner Module.
    TokenPoolFactory newTokenPoolFactory =
      new TokenPoolFactory(newTokenAdminRegistry, newRegistryModule, s_rmnProxy, address(s_destRouter));

    newTokenAdminRegistry.addRegistryModule(address(newRegistryModule));

    TokenPoolFactory.RemoteChainConfig memory remoteChainConfig = TokenPoolFactory.RemoteChainConfig(
      address(newTokenPoolFactory), address(s_destRouter), address(s_rmnProxy), REMOTE_TOKEN_DECIMALS
    );

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
    (address tokenAddress, address poolAddress) = s_tokenPoolFactory.deployTokenAndTokenPool(
      remoteTokenPools, LOCAL_TOKEN_DECIMALS, s_tokenInitCode, s_poolInitCode, FAKE_SALT
    );

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
      abi.encode(address(newRemoteToken), REMOTE_TOKEN_DECIMALS, new address[](0), s_rmnProxy, address(s_destRouter));

    // Take the init code and concat the destination params to it, the initCode shouldn't change
    bytes memory predictedPoolInitCode = abi.encodePacked(s_poolInitCode, predictedPoolCreationParams);

    // Predict the address of the pool on the DESTINATION chain
    address predictedPoolAddress =
      dynamicSalt.computeAddress(keccak256(predictedPoolInitCode), address(newTokenPoolFactory));

    // Assert that the address set for the remote pool is the same as the predicted address
    assertEq(
      abi.encode(predictedPoolAddress),
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      "Pool Address should have been predicted"
    );

    // On the new token pool factory, representing a destination chain,
    // deploy a new token and a new pool
    address newPoolAddress = newTokenPoolFactory.deployTokenPoolWithExistingToken(
      address(newRemoteToken),
      REMOTE_TOKEN_DECIMALS,
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
      TokenPool(poolAddress).getRemotePools(DEST_CHAIN_SELECTOR)[0],
      abi.encode(newPoolAddress),
      "New Pool Address should have been deployed correctly"
    );

    assertEq(TokenPool(poolAddress).getTokenDecimals(), LOCAL_TOKEN_DECIMALS, "Local token pool should use 18 decimals");

    // Assert the local token has 18 decimals
    assertEq(IERC20Metadata(tokenAddress).decimals(), LOCAL_TOKEN_DECIMALS, "Token Decimals should be 18");

    // Check configs on the remote pool and remote token decimals
    assertEq(TokenPool(newPoolAddress).getTokenDecimals(), REMOTE_TOKEN_DECIMALS, "Token Decimals should be 6");
    assertEq(address(TokenPool(newPoolAddress).getToken()), address(newRemoteToken), "Token Address should be set");
    assertEq(IERC20Metadata(newRemoteToken).decimals(), REMOTE_TOKEN_DECIMALS, "Token Decimals should be 6");
  }
}
