// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IFeeQuoter} from "../../interfaces/IFeeQuoter.sol";

import {KeystoneFeedsPermissionHandler} from "../../../keystone/KeystoneFeedsPermissionHandler.sol";
import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {FeeQuoterHelper} from "../helpers/FeeQuoterHelper.sol";
import {FeeQuoterFeeSetup, FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract FeeQuoter_constructor is FeeQuoterSetup {
  function test_Setup_Success() public virtual {
    address[] memory priceUpdaters = new address[](2);
    priceUpdaters[0] = STRANGER;
    priceUpdaters[1] = OWNER;
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_sourceTokens[1];
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](2);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);
    tokenPriceFeedUpdates[1] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[1], s_dataFeedByToken[s_sourceTokens[1]], 6);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();

    FeeQuoter.StaticConfig memory staticConfig = FeeQuoter.StaticConfig({
      linkToken: s_sourceTokens[0],
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      stalenessThreshold: uint32(TWELVE_HOURS)
    });
    s_feeQuoter = new FeeQuoterHelper(
      staticConfig,
      priceUpdaters,
      feeTokens,
      tokenPriceFeedUpdates,
      s_feeQuoterTokenTransferFeeConfigArgs,
      s_feeQuoterPremiumMultiplierWeiPerEthArgs,
      destChainConfigArgs
    );

    _assertFeeQuoterStaticConfigsEqual(s_feeQuoter.getStaticConfig(), staticConfig);
    assertEq(feeTokens, s_feeQuoter.getFeeTokens());
    assertEq(priceUpdaters, s_feeQuoter.getAllAuthorizedCallers());
    assertEq(s_feeQuoter.typeAndVersion(), "FeeQuoter 1.6.0-dev");

    _assertTokenPriceFeedConfigEquality(
      tokenPriceFeedUpdates[0].feedConfig, s_feeQuoter.getTokenPriceFeedConfig(s_sourceTokens[0])
    );

    _assertTokenPriceFeedConfigEquality(
      tokenPriceFeedUpdates[1].feedConfig, s_feeQuoter.getTokenPriceFeedConfig(s_sourceTokens[1])
    );

    assertEq(
      s_feeQuoterPremiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(s_feeQuoterPremiumMultiplierWeiPerEthArgs[0].token)
    );

    assertEq(
      s_feeQuoterPremiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(s_feeQuoterPremiumMultiplierWeiPerEthArgs[1].token)
    );

    FeeQuoter.TokenTransferFeeConfigArgs memory tokenTransferFeeConfigArg = s_feeQuoterTokenTransferFeeConfigArgs[0];
    for (uint256 i = 0; i < tokenTransferFeeConfigArg.tokenTransferFeeConfigs.length; ++i) {
      FeeQuoter.TokenTransferFeeConfigSingleTokenArgs memory tokenFeeArgs =
        s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[i];

      _assertTokenTransferFeeConfigEqual(
        tokenFeeArgs.tokenTransferFeeConfig,
        s_feeQuoter.getTokenTransferFeeConfig(tokenTransferFeeConfigArg.destChainSelector, tokenFeeArgs.token)
      );
    }

    for (uint256 i = 0; i < destChainConfigArgs.length; ++i) {
      FeeQuoter.DestChainConfig memory expectedConfig = destChainConfigArgs[i].destChainConfig;
      uint64 destChainSelector = destChainConfigArgs[i].destChainSelector;

      _assertFeeQuoterDestChainConfigsEqual(expectedConfig, s_feeQuoter.getDestChainConfig(destChainSelector));
    }
  }

  function test_InvalidStalenessThreshold_Revert() public {
    FeeQuoter.StaticConfig memory staticConfig = FeeQuoter.StaticConfig({
      linkToken: s_sourceTokens[0],
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      stalenessThreshold: 0
    });

    vm.expectRevert(FeeQuoter.InvalidStaticConfig.selector);

    s_feeQuoter = new FeeQuoterHelper(
      staticConfig,
      new address[](0),
      new address[](0),
      new FeeQuoter.TokenPriceFeedUpdate[](0),
      s_feeQuoterTokenTransferFeeConfigArgs,
      s_feeQuoterPremiumMultiplierWeiPerEthArgs,
      new FeeQuoter.DestChainConfigArgs[](0)
    );
  }

  function test_InvalidLinkTokenEqZeroAddress_Revert() public {
    FeeQuoter.StaticConfig memory staticConfig = FeeQuoter.StaticConfig({
      linkToken: address(0),
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      stalenessThreshold: uint32(TWELVE_HOURS)
    });

    vm.expectRevert(FeeQuoter.InvalidStaticConfig.selector);

    s_feeQuoter = new FeeQuoterHelper(
      staticConfig,
      new address[](0),
      new address[](0),
      new FeeQuoter.TokenPriceFeedUpdate[](0),
      s_feeQuoterTokenTransferFeeConfigArgs,
      s_feeQuoterPremiumMultiplierWeiPerEthArgs,
      new FeeQuoter.DestChainConfigArgs[](0)
    );
  }

  function test_InvalidMaxFeeJuelsPerMsg_Revert() public {
    FeeQuoter.StaticConfig memory staticConfig = FeeQuoter.StaticConfig({
      linkToken: s_sourceTokens[0],
      maxFeeJuelsPerMsg: 0,
      stalenessThreshold: uint32(TWELVE_HOURS)
    });

    vm.expectRevert(FeeQuoter.InvalidStaticConfig.selector);

    s_feeQuoter = new FeeQuoterHelper(
      staticConfig,
      new address[](0),
      new address[](0),
      new FeeQuoter.TokenPriceFeedUpdate[](0),
      s_feeQuoterTokenTransferFeeConfigArgs,
      s_feeQuoterPremiumMultiplierWeiPerEthArgs,
      new FeeQuoter.DestChainConfigArgs[](0)
    );
  }
}

contract FeeQuoter_getTokenPrices is FeeQuoterSetup {
  function test_GetTokenPrices_Success() public view {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    address[] memory tokens = new address[](3);
    tokens[0] = s_sourceTokens[0];
    tokens[1] = s_sourceTokens[1];
    tokens[2] = s_weth;

    Internal.TimestampedPackedUint224[] memory tokenPrices = s_feeQuoter.getTokenPrices(tokens);

    assertEq(tokenPrices.length, 3);
    assertEq(tokenPrices[0].value, priceUpdates.tokenPriceUpdates[0].usdPerToken);
    assertEq(tokenPrices[1].value, priceUpdates.tokenPriceUpdates[1].usdPerToken);
    assertEq(tokenPrices[2].value, priceUpdates.tokenPriceUpdates[2].usdPerToken);
  }
}

contract FeeQuoter_getTokenPrice is FeeQuoterSetup {
  function test_GetTokenPriceFromFeed_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Above staleness threshold
    vm.warp(originalTimestampValue + s_feeQuoter.getStaticConfig().stalenessThreshold + 1);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    Internal.TimestampedPackedUint224 memory tokenPriceAnswer = s_feeQuoter.getTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer.value, uint224(1e18));
    assertEq(tokenPriceAnswer.timestamp, uint32(block.timestamp));
  }
}

contract FeeQuoter_getValidatedTokenPrice is FeeQuoterSetup {
  function test_GetValidatedTokenPrice_Success() public view {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_feeQuoter.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  function test_GetValidatedTokenPriceFromFeed_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Right below staleness threshold
    vm.warp(originalTimestampValue + TWELVE_HOURS);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedOverStalenessPeriod_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Right above staleness threshold
    vm.warp(originalTimestampValue + TWELVE_HOURS + 1);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedMaxInt224Value_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, int256(uint256(type(uint224).max)));

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is: uint224.MAX_VALUE * (10 ** (36 - 18 - 18))
    assertEq(tokenPriceAnswer, uint224(type(uint224).max));
  }

  function test_GetValidatedTokenPriceFromFeedErc20Below18Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 6);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 8, 1e8);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 6);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e6) -> expected 1e30
    assertEq(tokenPriceAnswer, uint224(1e30));
  }

  function test_GetValidatedTokenPriceFromFeedErc20Above18Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 24);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 8, 1e8);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 24);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e24) -> expected 1e12
    assertEq(tokenPriceAnswer, uint224(1e12));
  }

  function test_GetValidatedTokenPriceFromFeedFeedAt18Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, 1e18);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedFeedAt0Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 0);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 0, 1e31);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 0);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e31 (0 decimal token) - unit is (1e18 * 1e18 / 1e0) -> expected 1e36
    assertEq(tokenPriceAnswer, uint224(1e67));
  }

  function test_GetValidatedTokenPriceFromFeedFlippedDecimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 20);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 20, 1e18);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 20);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e20) -> expected 1e14
    assertEq(tokenPriceAnswer, uint224(1e14));
  }

  function test_StaleFeeToken_Success() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_feeQuoter.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  // Reverts

  function test_OverflowFeedPrice_Revert() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, int256(uint256(type(uint224).max) + 1));

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    vm.expectRevert(FeeQuoter.DataFeedValueOutOfUint224Range.selector);
    s_feeQuoter.getValidatedTokenPrice(tokenAddress);
  }

  function test_UnderflowFeedPrice_Revert() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, -1);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    vm.expectRevert(FeeQuoter.DataFeedValueOutOfUint224Range.selector);
    s_feeQuoter.getValidatedTokenPrice(tokenAddress);
  }

  function test_TokenNotSupported_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_feeQuoter.getValidatedTokenPrice(DUMMY_CONTRACT_ADDRESS);
  }

  function test_TokenNotSupportedFeed_Revert() public {
    address sourceToken = _initialiseSingleTokenPriceFeed();
    MockV3Aggregator(s_dataFeedByToken[sourceToken]).updateAnswer(0);
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    priceUpdates.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: sourceToken, usdPerToken: 0});

    s_feeQuoter.updatePrices(priceUpdates);

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, sourceToken));
    s_feeQuoter.getValidatedTokenPrice(sourceToken);
  }
}

contract FeeQuoter_applyFeeTokensUpdates is FeeQuoterSetup {
  function test_ApplyFeeTokensUpdates_Success() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];

    vm.expectEmit();
    emit FeeQuoter.FeeTokenAdded(feeTokens[0]);

    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_feeQuoter.getFeeTokens().length, 3);
    assertEq(s_feeQuoter.getFeeTokens()[2], feeTokens[0]);

    // add same feeToken is no-op
    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_feeQuoter.getFeeTokens().length, 3);
    assertEq(s_feeQuoter.getFeeTokens()[2], feeTokens[0]);

    vm.expectEmit();
    emit FeeQuoter.FeeTokenRemoved(feeTokens[0]);

    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_feeQuoter.getFeeTokens().length, 2);

    // removing already removed feeToken is no-op
    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_feeQuoter.getFeeTokens().length, 2);
  }

  function test_OnlyCallableByOwner_Revert() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = STRANGER;
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));
  }
}

contract FeeQuoter_updatePrices is FeeQuoterSetup {
  function test_OnlyTokenPrice_Success() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    update.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(
      update.tokenPriceUpdates[0].sourceToken, update.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );

    s_feeQuoter.updatePrices(update);

    assertEq(s_feeQuoter.getTokenPrice(s_sourceTokens[0]).value, update.tokenPriceUpdates[0].usdPerToken);
  }

  function test_OnlyGasPrice_Success() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](1)
    });
    update.gasPriceUpdates[0] =
      Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});

    vm.expectEmit();
    emit FeeQuoter.UsdPerUnitGasUpdated(
      update.gasPriceUpdates[0].destChainSelector, update.gasPriceUpdates[0].usdPerUnitGas, block.timestamp
    );

    s_feeQuoter.updatePrices(update);

    assertEq(
      s_feeQuoter.getDestinationChainGasPrice(DEST_CHAIN_SELECTOR).value, update.gasPriceUpdates[0].usdPerUnitGas
    );
  }

  function test_UpdateMultiplePrices_Success() public {
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](3);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[1], usdPerToken: 1800e18});
    tokenPriceUpdates[2] = Internal.TokenPriceUpdate({sourceToken: address(12345), usdPerToken: 1e18});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](3);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2e6});
    gasPriceUpdates[1] = Internal.GasPriceUpdate({destChainSelector: SOURCE_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});
    gasPriceUpdates[2] = Internal.GasPriceUpdate({destChainSelector: 12345, usdPerUnitGas: 1e18});

    Internal.PriceUpdates memory update =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: gasPriceUpdates});

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit FeeQuoter.UsdPerTokenUpdated(
        update.tokenPriceUpdates[i].sourceToken, update.tokenPriceUpdates[i].usdPerToken, block.timestamp
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit FeeQuoter.UsdPerUnitGasUpdated(
        update.gasPriceUpdates[i].destChainSelector, update.gasPriceUpdates[i].usdPerUnitGas, block.timestamp
      );
    }

    s_feeQuoter.updatePrices(update);

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      assertEq(
        s_feeQuoter.getTokenPrice(update.tokenPriceUpdates[i].sourceToken).value, tokenPriceUpdates[i].usdPerToken
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      assertEq(
        s_feeQuoter.getDestinationChainGasPrice(update.gasPriceUpdates[i].destChainSelector).value,
        gasPriceUpdates[i].usdPerUnitGas
      );
    }
  }

  function test_UpdatableByAuthorizedCaller_Success() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    priceUpdates.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    // Revert when caller is not authorized
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_feeQuoter.updatePrices(priceUpdates);

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;
    vm.startPrank(OWNER);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );

    // Stranger is now an authorized caller to update prices
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(
      priceUpdates.tokenPriceUpdates[0].sourceToken, priceUpdates.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );
    s_feeQuoter.updatePrices(priceUpdates);

    assertEq(s_feeQuoter.getTokenPrice(s_sourceTokens[0]).value, priceUpdates.tokenPriceUpdates[0].usdPerToken);

    vm.startPrank(OWNER);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: new address[](0), removedCallers: priceUpdaters})
    );

    // Revert when authorized caller is removed
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_feeQuoter.updatePrices(priceUpdates);
  }

  // Reverts

  function test_OnlyCallableByUpdater_Revert() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_feeQuoter.updatePrices(priceUpdates);
  }
}

contract FeeQuoter_convertTokenAmount is FeeQuoterSetup {
  function test_ConvertTokenAmount_Success() public view {
    Internal.PriceUpdates memory initialPriceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    uint256 amount = 3e16;
    uint256 conversionRate = (uint256(initialPriceUpdates.tokenPriceUpdates[2].usdPerToken) * 1e18)
      / uint256(initialPriceUpdates.tokenPriceUpdates[0].usdPerToken);
    uint256 expected = (amount * conversionRate) / 1e18;
    assertEq(s_feeQuoter.convertTokenAmount(s_weth, amount, s_sourceTokens[0]), expected);
  }

  function test_Fuzz_ConvertTokenAmount_Success(
    uint256 feeTokenAmount,
    uint224 usdPerFeeToken,
    uint160 usdPerLinkToken,
    uint224 usdPerUnitGas
  ) public {
    vm.assume(usdPerFeeToken > 0);
    vm.assume(usdPerLinkToken > 0);
    // We bound the max fees to be at most uint96.max link.
    feeTokenAmount = bound(feeTokenAmount, 0, (uint256(type(uint96).max) * usdPerLinkToken) / usdPerFeeToken);

    address feeToken = address(1);
    address linkToken = address(2);
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = feeToken;
    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](2);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: feeToken, usdPerToken: usdPerFeeToken});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: linkToken, usdPerToken: usdPerLinkToken});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: usdPerUnitGas});

    Internal.PriceUpdates memory priceUpdates =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: gasPriceUpdates});

    s_feeQuoter.updatePrices(priceUpdates);

    uint256 linkFee = s_feeQuoter.convertTokenAmount(feeToken, feeTokenAmount, linkToken);
    assertEq(linkFee, (feeTokenAmount * usdPerFeeToken) / usdPerLinkToken);
  }

  // Reverts

  function test_LinkTokenNotSupported_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_feeQuoter.convertTokenAmount(DUMMY_CONTRACT_ADDRESS, 3e16, s_sourceTokens[0]);

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_feeQuoter.convertTokenAmount(s_sourceTokens[0], 3e16, DUMMY_CONTRACT_ADDRESS);
  }
}

contract FeeQuoter_getTokenAndGasPrices is FeeQuoterSetup {
  function test_GetFeeTokenAndGasPrices_Success() public view {
    (uint224 feeTokenPrice, uint224 gasPrice) = s_feeQuoter.getTokenAndGasPrices(s_sourceFeeToken, DEST_CHAIN_SELECTOR);

    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    assertEq(feeTokenPrice, s_sourceTokenPrices[0]);
    assertEq(gasPrice, priceUpdates.gasPriceUpdates[0].usdPerUnitGas);
  }

  function test_ZeroGasPrice_Success() public {
    uint64 zeroGasDestChainSelector = 345678;
    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: zeroGasDestChainSelector, usdPerUnitGas: 0});

    Internal.PriceUpdates memory priceUpdates =
      Internal.PriceUpdates({tokenPriceUpdates: new Internal.TokenPriceUpdate[](0), gasPriceUpdates: gasPriceUpdates});
    s_feeQuoter.updatePrices(priceUpdates);

    (, uint224 gasPrice) = s_feeQuoter.getTokenAndGasPrices(s_sourceFeeToken, zeroGasDestChainSelector);

    assertEq(gasPrice, priceUpdates.gasPriceUpdates[0].usdPerUnitGas);
  }

  function test_UnsupportedChain_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.ChainNotSupported.selector, DEST_CHAIN_SELECTOR + 1));
    s_feeQuoter.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR + 1);
  }

  function test_StaleGasPrice_Revert() public {
    uint256 diff = TWELVE_HOURS + 1;
    vm.warp(block.timestamp + diff);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.StaleGasPrice.selector, DEST_CHAIN_SELECTOR, TWELVE_HOURS, diff));
    s_feeQuoter.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR);
  }
}

contract FeeQuoter_updateTokenPriceFeeds is FeeQuoterSetup {
  function test_ZeroFeeds_Success() public {
    Vm.Log[] memory logEntries = vm.getRecordedLogs();

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](0);
    vm.recordLogs();
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    // Verify no log emissions
    assertEq(logEntries.length, 0);
  }

  function test_SingleFeedUpdate_Success() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    _assertTokenPriceFeedConfigUnconfigured(s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken));

    vm.expectEmit();
    emit FeeQuoter.PriceFeedPerTokenUpdated(tokenPriceFeedUpdates[0].sourceToken, tokenPriceFeedUpdates[0].feedConfig);

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
  }

  function test_MultipleFeedUpdate_Success() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](2);

    for (uint256 i = 0; i < 2; ++i) {
      tokenPriceFeedUpdates[i] =
        _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[i], s_dataFeedByToken[s_sourceTokens[i]], 18);

      _assertTokenPriceFeedConfigUnconfigured(s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[i].sourceToken));

      vm.expectEmit();
      emit FeeQuoter.PriceFeedPerTokenUpdated(tokenPriceFeedUpdates[i].sourceToken, tokenPriceFeedUpdates[i].feedConfig);
    }

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[1].sourceToken), tokenPriceFeedUpdates[1].feedConfig
    );
  }

  function test_FeedUnset_Success() public {
    Internal.TimestampedPackedUint224 memory priceQueryInitial = s_feeQuoter.getTokenPrice(s_sourceTokens[0]);
    assertFalse(priceQueryInitial.value == 0);
    assertFalse(priceQueryInitial.timestamp == 0);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );

    tokenPriceFeedUpdates[0].feedConfig.dataFeedAddress = address(0);
    vm.expectEmit();
    emit FeeQuoter.PriceFeedPerTokenUpdated(tokenPriceFeedUpdates[0].sourceToken, tokenPriceFeedUpdates[0].feedConfig);

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );

    // Price data should remain after a feed has been set->unset
    Internal.TimestampedPackedUint224 memory priceQueryPostUnsetFeed = s_feeQuoter.getTokenPrice(s_sourceTokens[0]);
    assertEq(priceQueryPostUnsetFeed.value, priceQueryInitial.value);
    assertEq(priceQueryPostUnsetFeed.timestamp, priceQueryInitial.timestamp);
  }

  function test_FeedNotUpdated() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
  }

  // Reverts

  function test_FeedUpdatedByNonOwner_Revert() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);
  }
}

contract FeeQuoter_applyDestChainConfigUpdates is FeeQuoterSetup {
  function test_Fuzz_applyDestChainConfigUpdates_Success(
    FeeQuoter.DestChainConfigArgs memory destChainConfigArgs
  ) public {
    vm.assume(destChainConfigArgs.destChainSelector != 0);
    vm.assume(destChainConfigArgs.destChainConfig.maxPerMsgGasLimit != 0);
    destChainConfigArgs.destChainConfig.defaultTxGasLimit = uint32(
      bound(
        destChainConfigArgs.destChainConfig.defaultTxGasLimit, 1, destChainConfigArgs.destChainConfig.maxPerMsgGasLimit
      )
    );
    destChainConfigArgs.destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;

    bool isNewChain = destChainConfigArgs.destChainSelector != DEST_CHAIN_SELECTOR;

    FeeQuoter.DestChainConfigArgs[] memory newDestChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    newDestChainConfigArgs[0] = destChainConfigArgs;

    if (isNewChain) {
      vm.expectEmit();
      emit FeeQuoter.DestChainAdded(destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig);
    } else {
      vm.expectEmit();
      emit FeeQuoter.DestChainConfigUpdated(destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig);
    }

    s_feeQuoter.applyDestChainConfigUpdates(newDestChainConfigArgs);

    _assertFeeQuoterDestChainConfigsEqual(
      destChainConfigArgs.destChainConfig, s_feeQuoter.getDestChainConfig(destChainConfigArgs.destChainSelector)
    );
  }

  function test_applyDestChainConfigUpdates_Success() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](2);
    destChainConfigArgs[0] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[0].destChainConfig.isEnabled = false;
    destChainConfigArgs[1] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    vm.expectEmit();
    emit FeeQuoter.DestChainConfigUpdated(DEST_CHAIN_SELECTOR, destChainConfigArgs[0].destChainConfig);
    vm.expectEmit();
    emit FeeQuoter.DestChainAdded(DEST_CHAIN_SELECTOR + 1, destChainConfigArgs[1].destChainConfig);

    vm.recordLogs();
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    FeeQuoter.DestChainConfig memory gotDestChainConfig0 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
    FeeQuoter.DestChainConfig memory gotDestChainConfig1 = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);

    assertEq(vm.getRecordedLogs().length, 2);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[0].destChainConfig, gotDestChainConfig0);
    _assertFeeQuoterDestChainConfigsEqual(destChainConfigArgs[1].destChainConfig, gotDestChainConfig1);
  }

  function test_applyDestChainConfigUpdatesZeroIntput_Success() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](0);

    vm.recordLogs();
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitEqZero_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.defaultTxGasLimit = 0;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitGtMaxPerMessageGasLimit_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    // Allow setting to the max value
    destChainConfigArg.destChainConfig.defaultTxGasLimit = destChainConfigArg.destChainConfig.maxPerMsgGasLimit;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    // Revert when exceeding max value
    destChainConfigArg.destChainConfig.defaultTxGasLimit = destChainConfigArg.destChainConfig.maxPerMsgGasLimit + 1;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestChainConfigDestChainSelectorEqZero_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainSelector = 0;
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidChainFamilySelector_Revert() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    FeeQuoter.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.chainFamilySelector = bytes4(uint32(1));

    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
  }
}

contract FeeQuoter_getDataAvailabilityCost is FeeQuoterSetup {
  function test_EmptyMessageCalculatesDataAvailabilityCost_Success() public {
    uint256 dataAvailabilityCostUSD =
      s_feeQuoter.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);

    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

    uint256 dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);

    // Test that the cost is destnation chain specific
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR + 1;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityOverheadGas =
      destChainConfig.destDataAvailabilityOverheadGas * 2;
    destChainConfigArgs[0].destChainConfig.destGasPerDataAvailabilityByte =
      destChainConfig.destGasPerDataAvailabilityByte * 2;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps =
      destChainConfig.destDataAvailabilityMultiplierBps * 2;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);
    uint256 dataAvailabilityCostUSD2 =
      s_feeQuoter.getDataAvailabilityCost(DEST_CHAIN_SELECTOR + 1, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);
    dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * Internal.MESSAGE_FIXED_BYTES;
    expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD2);
    assertFalse(dataAvailabilityCostUSD == dataAvailabilityCostUSD2);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCost_Success() public view {
    uint256 dataAvailabilityCostUSD =
      s_feeQuoter.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

    uint256 dataAvailabilityLengthBytes =
      Internal.MESSAGE_FIXED_BYTES + 100 + (5 * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + 50;
    uint256 dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCostUnsupportedDestChainSelector_Success() public view {
    uint256 dataAvailabilityCostUSD = s_feeQuoter.getDataAvailabilityCost(0, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    assertEq(dataAvailabilityCostUSD, 0);
  }

  function test_Fuzz_ZeroDataAvailabilityGasPriceAlwaysCalculatesZeroDataAvailabilityCost_Success(
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public view {
    uint256 dataAvailabilityCostUSD = s_feeQuoter.getDataAvailabilityCost(
      DEST_CHAIN_SELECTOR, 0, messageDataLength, numberOfTokens, tokenTransferBytesOverhead
    );

    assertEq(0, dataAvailabilityCostUSD);
  }

  function test_Fuzz_CalculateDataAvailabilityCost_Success(
    uint64 destChainSelector,
    uint32 destDataAvailabilityOverheadGas,
    uint16 destGasPerDataAvailabilityByte,
    uint16 destDataAvailabilityMultiplierBps,
    uint112 dataAvailabilityGasPrice,
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public {
    vm.assume(destChainSelector != 0);
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(destChainSelector);
    destChainConfigArgs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: destChainSelector, destChainConfig: destChainConfig});
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityOverheadGas = destDataAvailabilityOverheadGas;
    destChainConfigArgs[0].destChainConfig.destGasPerDataAvailabilityByte = destGasPerDataAvailabilityByte;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = destDataAvailabilityMultiplierBps;
    destChainConfigArgs[0].destChainConfig.defaultTxGasLimit = GAS_LIMIT;
    destChainConfigArgs[0].destChainConfig.maxPerMsgGasLimit = GAS_LIMIT;
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 dataAvailabilityCostUSD = s_feeQuoter.getDataAvailabilityCost(
      destChainConfigArgs[0].destChainSelector,
      dataAvailabilityGasPrice,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
    );

    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES + messageDataLength
      + (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + tokenTransferBytesOverhead;

    uint256 dataAvailabilityGas =
      destDataAvailabilityOverheadGas + destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      dataAvailabilityGasPrice * dataAvailabilityGas * destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }
}

contract FeeQuoter_applyPremiumMultiplierWeiPerEthUpdates is FeeQuoterSetup {
  function test_Fuzz_applyPremiumMultiplierWeiPerEthUpdates_Success(
    FeeQuoter.PremiumMultiplierWeiPerEthArgs memory premiumMultiplierWeiPerEthArg
  ) public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = premiumMultiplierWeiPerEthArg;

    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      premiumMultiplierWeiPerEthArg.token, premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth
    );

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(premiumMultiplierWeiPerEthArg.token)
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesSingleToken_Success() public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = s_feeQuoterPremiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);

    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      s_feeQuoterPremiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesMultipleTokens_Success() public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](2);
    premiumMultiplierWeiPerEthArgs[0] = s_feeQuoterPremiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);
    premiumMultiplierWeiPerEthArgs[1].token = vm.addr(2);

    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );
    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(2), premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth
    );

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
    assertEq(
      premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(vm.addr(2))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesZeroInput() public {
    vm.recordLogs();
    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs;
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);
  }
}

contract FeeQuoter_applyTokenTransferFeeConfigUpdates is FeeQuoterSetup {
  function test_Fuzz_ApplyTokenTransferFeeConfig_Success(
    FeeQuoter.TokenTransferFeeConfig[2] memory tokenTransferFeeConfigs
  ) public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(2, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    for (uint256 i = 0; i < tokenTransferFeeConfigArgs.length; ++i) {
      for (uint256 j = 0; j < tokenTransferFeeConfigs.length; ++j) {
        tokenTransferFeeConfigs[j].destBytesOverhead = uint32(
          bound(tokenTransferFeeConfigs[j].destBytesOverhead, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, type(uint32).max)
        );
        address feeToken = s_sourceTokens[j];
        tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs[j].token = feeToken;
        tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs[j].tokenTransferFeeConfig = tokenTransferFeeConfigs[j];

        vm.expectEmit();
        emit FeeQuoter.TokenTransferFeeConfigUpdated(
          tokenTransferFeeConfigArgs[i].destChainSelector, feeToken, tokenTransferFeeConfigs[j]
        );
      }
    }

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    for (uint256 i = 0; i < tokenTransferFeeConfigs.length; ++i) {
      _assertTokenTransferFeeConfigEqual(
        tokenTransferFeeConfigs[i],
        s_feeQuoter.getTokenTransferFeeConfig(
          tokenTransferFeeConfigArgs[0].destChainSelector,
          tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[i].token
        )
      );
    }
  }

  function test_ApplyTokenTransferFeeConfig_Success() public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: 312,
      isEnabled: true
    });
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token = address(11);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 12,
      maxFeeUSDCents: 13,
      deciBps: 14,
      destGasOverhead: 15,
      destBytesOverhead: 394,
      isEnabled: true
    });

    vm.expectEmit();
    emit FeeQuoter.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig
    );
    vm.expectEmit();
    emit FeeQuoter.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig
    );

    FeeQuoter.TokenTransferFeeConfigRemoveArgs[] memory tokensToRemove =
      new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0);
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs, tokensToRemove);

    FeeQuoter.TokenTransferFeeConfig memory config0 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    FeeQuoter.TokenTransferFeeConfig memory config1 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig, config0
    );
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );

    // Remove only the first token and validate only the first token is removed
    tokensToRemove = new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](1);
    tokensToRemove[0] = FeeQuoter.TokenTransferFeeConfigRemoveArgs({
      destChainSelector: tokenTransferFeeConfigArgs[0].destChainSelector,
      token: tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    });

    vm.expectEmit();
    emit FeeQuoter.TokenTransferFeeConfigDeleted(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(new FeeQuoter.TokenTransferFeeConfigArgs[](0), tokensToRemove);

    config0 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    config1 = s_feeQuoter.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    FeeQuoter.TokenTransferFeeConfig memory emptyConfig;

    _assertTokenTransferFeeConfigEqual(emptyConfig, config0);
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );
  }

  function test_ApplyTokenTransferFeeZeroInput() public {
    vm.recordLogs();
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      new FeeQuoter.TokenTransferFeeConfigArgs[](0), new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    vm.startPrank(STRANGER);
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs;

    vm.expectRevert("Only callable by owner");

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }

  function test_InvalidDestBytesOverhead_Revert() public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES - 1),
      isEnabled: true
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        FeeQuoter.InvalidDestBytesOverhead.selector,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.destBytesOverhead
      )
    );

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }
}

contract FeeQuoter_getTokenTransferCost is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_NoTokenTransferChargesZeroFee_Success() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(0, feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  function test_getTokenTransferCost_selfServeUsesDefaults_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_selfServeTokenDefaultPricing, 1000);

    // Get config to assert it isn't set
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    assertFalse(transferFeeConfig.isEnabled);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // Assert that the default values are used
    assertEq(uint256(DEFAULT_TOKEN_FEE_USD_CENTS) * 1e16, feeUSDWei);
    assertEq(DEFAULT_TOKEN_DEST_GAS_OVERHEAD, destGasOverhead);
    assertEq(DEFAULT_TOKEN_BYTES_OVERHEAD, destBytesOverhead);
  }

  function test_SmallTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1000);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(_configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroAmountTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 0);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(_configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_LargeTokenTransferChargesMaxFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(_configUSDCentToWei(transferFeeConfig.maxFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_FeeTokenBpsFee_Success() public view {
    uint256 tokenAmount = 10000e18;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = _calcUSDValueFromTokenAmount(s_feeTokenPrice, tokenAmount);
    uint256 bpsUSDWei = _applyBpsRatio(
      usdWei, s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.deciBps
    );

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_CustomTokenBpsFee_Success() public view {
    uint256 tokenAmount = 200000e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: tokenAmount});

    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = _calcUSDValueFromTokenAmount(s_customTokenPrice, tokenAmount);
    uint256 bpsUSDWei = _applyBpsRatio(
      usdWei, s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig.deciBps
    );

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroFeeConfigChargesMinFee_Success() public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = s_sourceFeeToken;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      isEnabled: true
    });
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // if token charges 0 bps, it should cost minFee to transfer
    assertEq(
      _configUSDCentToWei(
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.minFeeUSDCents
      ),
      feeUSDWei
    );
    assertEq(0, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_Fuzz_TokenTransferFeeDuplicateTokens_Success(uint256 transfers, uint256 amount) public view {
    // It shouldn't be possible to pay materially lower fees by splitting up the transfers.
    // Note it is possible to pay higher fees since the minimum fees are added.
    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
    transfers = bound(transfers, 1, destChainConfig.maxNumberOfTokensPerMsg);
    // Cap amount to avoid overflow
    amount = bound(amount, 0, 1e36);
    Client.EVMTokenAmount[] memory multiple = new Client.EVMTokenAmount[](transfers);
    for (uint256 i = 0; i < transfers; ++i) {
      multiple[i] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: amount});
    }
    Client.EVMTokenAmount[] memory single = new Client.EVMTokenAmount[](1);
    single[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: amount * transfers});

    address feeToken = s_sourceRouter.getWrappedNative();

    (uint256 feeSingleUSDWei, uint32 gasOverheadSingle, uint32 bytesOverheadSingle) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, single);
    (uint256 feeMultipleUSDWei, uint32 gasOverheadMultiple, uint32 bytesOverheadMultiple) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, multiple);

    // Note that there can be a rounding error once per split.
    assertGe(feeMultipleUSDWei, (feeSingleUSDWei - destChainConfig.maxNumberOfTokensPerMsg));
    assertEq(gasOverheadMultiple, gasOverheadSingle * transfers);
    assertEq(bytesOverheadMultiple, bytesOverheadSingle * transfers);
  }

  function test_MixedTokenTransferFee_Success() public view {
    address[3] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative(), CUSTOM_TOKEN];
    uint224[3] memory tokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice, s_customTokenPrice];
    FeeQuoter.TokenTransferFeeConfig[3] memory tokenTransferFeeConfigs = [
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[0]),
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[1]),
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[2])
    ];

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](3),
      feeToken: s_sourceRouter.getWrappedNative(),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    uint256 expectedTotalGas = 0;
    uint256 expectedTotalBytes = 0;

    // Start with small token transfers, total bps fee is lower than min token transfer fee
    for (uint256 i = 0; i < testTokens.length; ++i) {
      message.tokenAmounts[i] = Client.EVMTokenAmount({token: testTokens[i], amount: 1e14});
      FeeQuoter.TokenTransferFeeConfig memory tokenTransferFeeConfig =
        s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[i]);

      expectedTotalGas += tokenTransferFeeConfig.destGasOverhead == 0
        ? DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        : tokenTransferFeeConfig.destGasOverhead;
      expectedTotalBytes += tokenTransferFeeConfig.destBytesOverhead == 0
        ? DEFAULT_TOKEN_BYTES_OVERHEAD
        : tokenTransferFeeConfig.destBytesOverhead;
    }
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);

    uint256 expectedFeeUSDWei = 0;
    for (uint256 i = 0; i < testTokens.length; ++i) {
      expectedFeeUSDWei += _configUSDCentToWei(
        tokenTransferFeeConfigs[i].minFeeUSDCents == 0
          ? DEFAULT_TOKEN_FEE_USD_CENTS
          : tokenTransferFeeConfigs[i].minFeeUSDCents
      );
    }

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 1");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 1");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 1");

    // Set 1st token transfer to a meaningful amount so its bps fee is now between min and max fee
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: testTokens[0], amount: 10000e18});

    uint256 token0USDWei = _applyBpsRatio(
      _calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount), tokenTransferFeeConfigs[0].deciBps
    );
    uint256 token1USDWei = _configUSDCentToWei(DEFAULT_TOKEN_FEE_USD_CENTS);

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = token0USDWei + token1USDWei + _configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 2");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 2");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 2");

    // Set 2nd token transfer to a large amount that is higher than maxFeeUSD
    message.tokenAmounts[2] = Client.EVMTokenAmount({token: testTokens[2], amount: 1e36});

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = token0USDWei + token1USDWei + _configUSDCentToWei(tokenTransferFeeConfigs[2].maxFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 3");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 3");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 3");
  }
}

contract FeeQuoter_getValidatedFee is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (_configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_ZeroDataAvailabilityMultiplier_Success() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
    destChainConfigArgs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: destChainConfig});
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = 0;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);

    uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

    uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
    uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
    uint256 messageFeeUSD = (_configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);

    uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD) / s_feeTokenPrice;
    assertEq(totalPriceInFeeToken, feeAmount);
  }

  function test_HighGasMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = MAX_GAS_LIMIT;
    uint256 customDataSize = MAX_DATA_SIZE;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: new bytes(customDataSize),
        tokenAmounts: new Client.EVMTokenAmount[](0),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });

      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + customDataSize * DEST_GAS_PER_PAYLOAD_BYTE;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (_configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_SingleTokenMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 tokenAmount = 10000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
      message.feeToken = testTokens[i];
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
      uint32 destBytesOverhead =
        s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destBytesOverhead;
      uint32 tokenBytesOverhead =
        destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD
        + s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken));
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_MessageWithDataAndTokenTransfer_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = 1_000_000;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: "",
        tokenAmounts: new Client.EVMTokenAmount[](2),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });
      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: 10000e18}); // feeTokenAmount
      message.tokenAmounts[1] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: 200000e18}); // customTokenAmount
      message.data = "random bits and bytes that should be factored into the cost of the message";

      uint32 tokenGasOverhead = 0;
      uint32 tokenBytesOverhead = 0;
      for (uint256 j = 0; j < message.tokenAmounts.length; ++j) {
        tokenGasOverhead +=
          s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destGasOverhead;
        uint32 destBytesOverhead =
          s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destBytesOverhead;
        tokenBytesOverhead += destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;
      }

      uint256 gasUsed =
        customGasLimit + DEST_GAS_OVERHEAD + message.data.length * DEST_GAS_PER_PAYLOAD_BYTE + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message));
    }
  }

  function test_Fuzz_EnforceOutOfOrder(bool enforce, bool allowOutOfOrderExecution) public {
    // Update config to enforce allowOutOfOrderExecution = defaultVal.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = enforce;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG,
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: allowOutOfOrderExecution})
    );

    // If enforcement is on, only true should be allowed.
    if (enforce && !allowOutOfOrderExecution) {
      vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    }
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Reverts

  function test_DestinationChainNotEnabled_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.DestinationChainNotEnabled.selector, DEST_CHAIN_SELECTOR + 1));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR + 1, _generateEmptyMessage());
  }

  function test_EnforceOutOfOrder_Revert() public {
    // Update config to enforce allowOutOfOrderExecution = true.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = true;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Empty extraArgs to should revert since it enforceOutOfOrder is true.
    message.extraArgs = "";

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(FeeQuoter.UnsupportedNumberOfTokens.selector);
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.MessageGasLimitTooHigh.selector));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_NotAFeeToken_Revert() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.FeeTokenNotSupported.selector, notAFeeToken));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_InvalidEVMAddress_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(type(uint208).max);

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract FeeQuoter_processMessageArgs is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function setUp() public virtual override {
    super.setUp();
  }

  function test_processMessageArgs_WithLinkTokenAmount_Success() public view {
    (
      uint256 msgFeeJuels,
      /* bool isOutOfOrderExecution */
      ,
      /* bytes memory convertedExtraArgs */
      ,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      // LINK
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(msgFeeJuels, MAX_MSG_FEES_JUELS);
  }

  function test_processMessageArgs_WithConvertedTokenAmount_Success() public view {
    address feeToken = s_sourceTokens[1];
    uint256 feeTokenAmount = 10_000 gwei;
    uint256 expectedConvertedAmount = s_feeQuoter.convertTokenAmount(feeToken, feeTokenAmount, s_sourceTokens[0]);

    (
      uint256 msgFeeJuels,
      /* bool isOutOfOrderExecution */
      ,
      /* bytes memory convertedExtraArgs */
      ,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      feeToken,
      feeTokenAmount,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(msgFeeJuels, expectedConvertedAmount);
  }

  function test_processMessageArgs_WithEmptyEVMExtraArgs_Success() public view {
    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(isOutOfOrderExecution, false);
    assertEq(convertedExtraArgs, Client._argsToBytes(s_feeQuoter.parseEVMExtraArgsFromBytes("", DEST_CHAIN_SELECTOR)));
  }

  function test_processMessageArgs_WithEVMExtraArgsV1_Success() public view {
    bytes memory extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 1000}));

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      extraArgs,
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(isOutOfOrderExecution, false);
    assertEq(
      convertedExtraArgs, Client._argsToBytes(s_feeQuoter.parseEVMExtraArgsFromBytes(extraArgs, DEST_CHAIN_SELECTOR))
    );
  }

  function test_processMessageArgs_WitEVMExtraArgsV2_Success() public view {
    bytes memory extraArgs = Client._argsToBytes(Client.EVMExtraArgsV2({gasLimit: 0, allowOutOfOrderExecution: true}));

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      extraArgs,
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(isOutOfOrderExecution, true);
    assertEq(
      convertedExtraArgs, Client._argsToBytes(s_feeQuoter.parseEVMExtraArgsFromBytes(extraArgs, DEST_CHAIN_SELECTOR))
    );
  }

  // Reverts

  function test_processMessageArgs_MessageFeeTooHigh_Revert() public {
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.MessageFeeTooHigh.selector, MAX_MSG_FEES_JUELS + 1, MAX_MSG_FEES_JUELS)
    );

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS + 1,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_processMessageArgs_InvalidExtraArgs_Revert() public {
    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      "abcde",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_processMessageArgs_MalformedEVMExtraArgs_Revert() public {
    // abi.decode error
    vm.expectRevert();

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      abi.encodeWithSelector(Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV1({gasLimit: 100})),
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_processMessageArgs_WithCorrectPoolReturnData_Success() public view {
    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](2);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];
    sourceTokenAmounts[1].amount = 1e18;
    sourceTokenAmounts[1].token = CUSTOM_TOKEN_2;

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](2);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    tokenAmounts[1] = _getSourceTokenData(sourceTokenAmounts[1], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    bytes[] memory expectedDestExecData = new bytes[](2);
    expectedDestExecData[0] = abi.encode(
      s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.destGasOverhead
    );
    expectedDestExecData[1] = abi.encode(DEFAULT_TOKEN_DEST_GAS_OVERHEAD); //expected return data should be abi.encoded  default as isEnabled is false

    // No revert - successful
    ( /* msgFeeJuels */ , /* isOutOfOrderExecution */, /* convertedExtraArgs */, bytes[] memory destExecData) =
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    for (uint256 i = 0; i < destExecData.length; ++i) {
      assertEq(destExecData[i], expectedDestExecData[i]);
    }
  }

  function test_processMessageArgs_TokenAmountArraysMismatching_Revert() public {
    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](2);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);

    // Revert due to index out of bounds access
    vm.expectRevert();

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS,
      "",
      new Internal.EVM2AnyTokenTransfer[](1),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_processMessageArgs_SourceTokenDataTooLarge_Revert() public {
    address sourceETH = s_sourceTokens[1];

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1000;
    sourceTokenAmounts[0].token = sourceETH;

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);

    // No data set, should succeed
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set max data length, should succeed
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES);
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set data to max length +1, should revert
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 1);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set token config to allow larger data
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = sourceETH;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      isEnabled: true
    });
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set the token data larger than the configured token data, should revert
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 32 + 1);

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );
  }

  function test_processMessageArgs_InvalidEVMAddressDestToken_Revert() public {
    bytes memory nonEvmAddress = abi.encode(type(uint208).max);

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    tokenAmounts[0].destTokenAddress = nonEvmAddress;

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, nonEvmAddress));
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );
  }
}

contract FeeQuoter_validateDestFamilyAddress is FeeQuoterSetup {
  function test_ValidEVMAddress_Success() public view {
    bytes memory encodedAddress = abi.encode(address(10000));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, encodedAddress);
  }

  function test_ValidNonEVMAddress_Success() public view {
    s_feeQuoter.validateDestFamilyAddress(bytes4(uint32(1)), abi.encode(type(uint208).max));
  }

  // Reverts

  function test_InvalidEVMAddress_Revert() public {
    bytes memory invalidAddress = abi.encode(type(uint208).max);
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
  }

  function test_InvalidEVMAddressEncodePacked_Revert() public {
    bytes memory invalidAddress = abi.encodePacked(address(234));
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
  }

  function test_InvalidEVMAddressPrecompiles_Revert() public {
    for (uint160 i = 0; i < Internal.PRECOMPILE_SPACE; ++i) {
      bytes memory invalidAddress = abi.encode(address(i));
      vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
      s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
    }

    s_feeQuoter.validateDestFamilyAddress(
      Internal.CHAIN_FAMILY_SELECTOR_EVM, abi.encode(address(uint160(Internal.PRECOMPILE_SPACE)))
    );
  }
}

contract FeeQuoter_parseEVMExtraArgsFromBytes is FeeQuoterSetup {
  FeeQuoter.DestChainConfig private s_destChainConfig;

  function setUp() public virtual override {
    super.setUp();
    s_destChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
  }

  function test_EVMExtraArgsV1_Success() public view {
    Client.EVMExtraArgsV1 memory inputArgs = Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    Client.EVMExtraArgsV2 memory expectedOutputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});

    vm.assertEq(
      abi.encode(s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig)),
      abi.encode(expectedOutputArgs)
    );
  }

  function test_EVMExtraArgsV2_Success() public view {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

    vm.assertEq(
      abi.encode(s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig)), abi.encode(inputArgs)
    );
  }

  function test_EVMExtraArgsDefault_Success() public view {
    Client.EVMExtraArgsV2 memory expectedOutputArgs =
      Client.EVMExtraArgsV2({gasLimit: s_destChainConfig.defaultTxGasLimit, allowOutOfOrderExecution: false});

    vm.assertEq(
      abi.encode(s_feeQuoter.parseEVMExtraArgsFromBytes("", s_destChainConfig)), abi.encode(expectedOutputArgs)
    );
  }

  // Reverts

  function test_EVMExtraArgsInvalidExtraArgsTag_Revert() public {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    // Invalidate selector
    inputExtraArgs[0] = bytes1(uint8(0));

    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);
    s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_EVMExtraArgsEnforceOutOfOrder_Revert() public {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    s_destChainConfig.enforceOutOfOrder = true;

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_EVMExtraArgsGasLimitTooHigh_Revert() public {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: s_destChainConfig.maxPerMsgGasLimit + 1, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

    vm.expectRevert(FeeQuoter.MessageGasLimitTooHigh.selector);
    s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }
}

contract FeeQuoter_KeystoneSetup is FeeQuoterSetup {
  address internal constant FORWARDER_1 = address(0x1);
  address internal constant WORKFLOW_OWNER_1 = address(0x3);
  bytes10 internal constant WORKFLOW_NAME_1 = "workflow1";
  bytes2 internal constant REPORT_NAME_1 = "01";
  address internal onReportTestToken1;
  address internal onReportTestToken2;

  function setUp() public virtual override {
    super.setUp();
    onReportTestToken1 = s_sourceTokens[0];
    onReportTestToken2 = _deploySourceToken("onReportTestToken2", 0, 20);

    KeystoneFeedsPermissionHandler.Permission[] memory permissions = new KeystoneFeedsPermissionHandler.Permission[](1);
    permissions[0] = KeystoneFeedsPermissionHandler.Permission({
      forwarder: FORWARDER_1,
      workflowOwner: WORKFLOW_OWNER_1,
      workflowName: WORKFLOW_NAME_1,
      reportName: REPORT_NAME_1,
      isAllowed: true
    });
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeeds = new FeeQuoter.TokenPriceFeedUpdate[](2);
    tokenPriceFeeds[0] = FeeQuoter.TokenPriceFeedUpdate({
      sourceToken: onReportTestToken1,
      feedConfig: FeeQuoter.TokenPriceFeedConfig({dataFeedAddress: address(0x0), tokenDecimals: 18})
    });
    tokenPriceFeeds[1] = FeeQuoter.TokenPriceFeedUpdate({
      sourceToken: onReportTestToken2,
      feedConfig: FeeQuoter.TokenPriceFeedConfig({dataFeedAddress: address(0x0), tokenDecimals: 20})
    });
    s_feeQuoter.setReportPermissions(permissions);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeeds);
  }
}

contract FeeQuoter_onReport is FeeQuoter_KeystoneSetup {
  function test_onReport_Success() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);

    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](2);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp)});
    report[1] =
      FeeQuoter.ReceivedCCIPFeedReport({token: onReportTestToken2, price: 4e18, timestamp: uint32(block.timestamp)});

    uint224 expectedStoredToken1Price = s_feeQuoter.calculateRebasedValue(18, 18, report[0].price);
    uint224 expectedStoredToken2Price = s_feeQuoter.calculateRebasedValue(18, 20, report[1].price);
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(onReportTestToken1, expectedStoredToken1Price, block.timestamp);
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(onReportTestToken2, expectedStoredToken2Price, block.timestamp);

    changePrank(FORWARDER_1);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));

    vm.assertEq(s_feeQuoter.getTokenPrice(report[0].token).value, expectedStoredToken1Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(report[0].token).timestamp, report[0].timestamp);

    vm.assertEq(s_feeQuoter.getTokenPrice(report[1].token).value, expectedStoredToken2Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(report[1].token).timestamp, report[1].timestamp);
  }

  function test_onReport_InvalidForwarder_Reverts() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_sourceTokens[0], price: 4e18, timestamp: uint32(block.timestamp)});

    vm.expectRevert(
      abi.encodeWithSelector(
        KeystoneFeedsPermissionHandler.ReportForwarderUnauthorized.selector,
        STRANGER,
        WORKFLOW_OWNER_1,
        WORKFLOW_NAME_1,
        REPORT_NAME_1
      )
    );
    changePrank(STRANGER);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }

  function test_onReport_UnsupportedToken_Reverts() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_sourceTokens[1], price: 4e18, timestamp: uint32(block.timestamp)});

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, s_sourceTokens[1]));
    changePrank(FORWARDER_1);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }

  function test_OnReport_StaleUpdate_Revert() public {
    //Creating a correct report
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);

    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp)});

    uint224 expectedStoredTokenPrice = s_feeQuoter.calculateRebasedValue(18, 18, report[0].price);

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(onReportTestToken1, expectedStoredTokenPrice, block.timestamp);

    changePrank(FORWARDER_1);
    //setting the correct price and time with the correct report
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));

    //create a stale report
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp - 1)});
    //expecting a revert
    vm.expectRevert(
      abi.encodeWithSelector(
        FeeQuoter.StaleKeystoneUpdate.selector, onReportTestToken1, block.timestamp - 1, block.timestamp
      )
    );
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }
}
