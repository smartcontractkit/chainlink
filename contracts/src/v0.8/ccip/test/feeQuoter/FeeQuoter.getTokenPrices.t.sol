// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getTokenPrices is FeeQuoterSetup {
  function test_GetTokenPrices() public view {
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
