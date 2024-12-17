// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract FeeQuoter_updateTokenPriceFeeds is FeeQuoterSetup {
  function test_ZeroFeeds() public {
    Vm.Log[] memory logEntries = vm.getRecordedLogs();

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](0);
    vm.recordLogs();
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    // Verify no log emissions
    assertEq(logEntries.length, 0);
  }

  function test_SingleFeedUpdate() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    _assertTokenPriceFeedConfigNotConfigured(s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken));

    vm.expectEmit();
    emit FeeQuoter.PriceFeedPerTokenUpdated(tokenPriceFeedUpdates[0].sourceToken, tokenPriceFeedUpdates[0].feedConfig);

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
  }

  function test_MultipleFeedUpdate() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](2);

    for (uint256 i = 0; i < 2; ++i) {
      tokenPriceFeedUpdates[i] =
        _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[i], s_dataFeedByToken[s_sourceTokens[i]], 18);

      _assertTokenPriceFeedConfigNotConfigured(
        s_feeQuoter.getTokenPriceFeedConfig(tokenPriceFeedUpdates[i].sourceToken)
      );

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

  function test_FeedUnset() public {
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

  function test_RevertWhen_FeedUpdatedByNonOwner() public {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);
  }

  function _assertTokenPriceFeedConfigNotConfigured(
    FeeQuoter.TokenPriceFeedConfig memory config
  ) internal pure virtual {
    _assertTokenPriceFeedConfigEquality(
      config, FeeQuoter.TokenPriceFeedConfig({dataFeedAddress: address(0), tokenDecimals: 0, isEnabled: false})
    );
  }
}
