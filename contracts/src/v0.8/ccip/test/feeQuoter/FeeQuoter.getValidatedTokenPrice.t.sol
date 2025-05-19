// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MockV3Aggregator} from "../../../shared/mocks/MockV3Aggregator.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getValidatedTokenPrice is FeeQuoterSetup {
  function test_GetValidatedTokenPrice() public view {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_feeQuoter.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  function test_GetValidatedTokenPriceFromFeed() public {
    uint256 originalTimestampValue = block.timestamp;

    // Right below staleness threshold
    vm.warp(originalTimestampValue + TWELVE_HOURS);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedOverStalenessPeriod() public {
    uint256 originalTimestampValue = block.timestamp;

    // Right above staleness threshold
    vm.warp(originalTimestampValue + TWELVE_HOURS + 1);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedMaxInt224Value() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, int256(uint256(type(uint224).max)));

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is: uint224.MAX_VALUE * (10 ** (36 - 18 - 18))
    assertEq(tokenPriceAnswer, uint224(type(uint224).max));
  }

  function test_GetValidatedTokenPriceFromFeedErc20Below18Decimals() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 6);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 8, 1e8);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 6);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e6) -> expected 1e30
    assertEq(tokenPriceAnswer, uint224(1e30));
  }

  function test_GetValidatedTokenPriceFromFeedErc20Above18Decimals() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 24);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 8, 1e8);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 24);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e24) -> expected 1e12
    assertEq(tokenPriceAnswer, uint224(1e12));
  }

  function test_GetValidatedTokenPriceFromFeedFeedAt18Decimals() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, 1e18);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedFeedAt0Decimals() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 0);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 0, 1e31);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 0);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e31 (0 decimal token) - unit is (1e18 * 1e18 / 1e0) -> expected 1e36
    assertEq(tokenPriceAnswer, uint224(1e67));
  }

  function test_GetValidatedTokenPriceFromFeedFlippedDecimals() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 20);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 20, 1e18);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 20);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_feeQuoter.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e20) -> expected 1e14
    assertEq(tokenPriceAnswer, uint224(1e14));
  }

  function test_StaleFeeToken() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_feeQuoter.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  // Reverts

  function test_RevertWhen_OverflowFeedPrice() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, int256(uint256(type(uint224).max) + 1));

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    vm.expectRevert(FeeQuoter.DataFeedValueOutOfUint224Range.selector);
    s_feeQuoter.getValidatedTokenPrice(tokenAddress);
  }

  function test_RevertWhen_UnderflowFeedPrice() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, -1);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = _getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    vm.expectRevert(FeeQuoter.DataFeedValueOutOfUint224Range.selector);
    s_feeQuoter.getValidatedTokenPrice(tokenAddress);
  }

  function test_RevertWhen_TokenNotSupported() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_feeQuoter.getValidatedTokenPrice(DUMMY_CONTRACT_ADDRESS);
  }

  function test_RevertWhen_TokenNotSupportedFeed() public {
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
