// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getTokenPrice is FeeQuoterSetup {
  function test_GetTokenPriceFromFeed_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Above staleness threshold
    vm.warp(originalTimestampValue + s_feeQuoter.getStaticConfig().tokenPriceStalenessThreshold + 1);

    address sourceToken = _initialiseSingleTokenPriceFeed();

    vm.expectCall(s_dataFeedByToken[sourceToken], abi.encodeWithSelector(MockV3Aggregator.latestRoundData.selector));

    Internal.TimestampedPackedUint224 memory tokenPriceAnswer = s_feeQuoter.getTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer.value, uint224(1e18));
    assertEq(tokenPriceAnswer.timestamp, uint32(originalTimestampValue));
  }

  function test_GetTokenPrice_LocalMoreRecent_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    update.tokenPriceUpdates[0] =
      Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: uint32(originalTimestampValue + 5)});

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(
      update.tokenPriceUpdates[0].sourceToken, update.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );

    s_feeQuoter.updatePrices(update);

    vm.warp(originalTimestampValue + s_feeQuoter.getStaticConfig().tokenPriceStalenessThreshold + 10);

    Internal.TimestampedPackedUint224 memory tokenPriceAnswer = s_feeQuoter.getTokenPrice(s_sourceTokens[0]);

    //Assert that the returned price is the local price, not the oracle price
    assertEq(tokenPriceAnswer.value, update.tokenPriceUpdates[0].usdPerToken);
    assertEq(tokenPriceAnswer.timestamp, uint32(originalTimestampValue));
  }
}
