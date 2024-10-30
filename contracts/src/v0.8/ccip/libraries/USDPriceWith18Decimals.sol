// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

library USDPriceWith18Decimals {
  /// @notice Takes a price in USD, with 18 decimals per 1e18 token amount, and amount of the smallest token
  /// denomination, calculates the value in USD with 18 decimals.
  /// @param tokenPrice The USD price of the token.
  /// @param tokenAmount Amount of the smallest token denomination.
  /// @return USD value with 18 decimals.
  /// @dev this function assumes that no more than 1e59 US dollar worth of token is passed in. If more is sent, this
  /// function will overflow and revert. Since there isn't even close to 1e59 dollars, this is ok for all legit tokens.
  function _calcUSDValueFromTokenAmount(uint224 tokenPrice, uint256 tokenAmount) internal pure returns (uint256) {
    /// LINK Example:
    /// tokenPrice:         8e18 -> $8/LINK, as 1e18 token amount is 1 LINK, worth 8 USD, or 8e18 with 18 decimals
    /// tokenAmount:        2e18 -> 2 LINK
    /// result:             8e18 * 2e18 / 1e18 -> 16e18 with 18 decimals = $16

    /// USDC Example:
    /// tokenPrice:         1e30 -> $1/USDC, as 1e18 token amount is 1e12 USDC, worth 1e12 USD, or 1e30 with 18 decimals
    /// tokenAmount:        5e6  -> 5 USDC
    /// result:             1e30 * 5e6 / 1e18 -> 5e18 with 18 decimals = $5
    return (tokenPrice * tokenAmount) / 1e18;
  }

  /// @notice Takes a price in USD, with 18 decimals per 1e18 token amount, and USD value with 18 decimals, calculates
  /// amount of the smallest token denomination.
  /// @param tokenPrice The USD price of the token.
  /// @param usdValue USD value with 18 decimals.
  /// @return Amount of the smallest token denomination.
  function _calcTokenAmountFromUSDValue(uint224 tokenPrice, uint256 usdValue) internal pure returns (uint256) {
    /// LINK Example:
    /// tokenPrice:          8e18 -> $8/LINK, as 1e18 token amount is 1 LINK, worth 8 USD, or 8e18 with 18 decimals
    /// usdValue:           16e18 -> $16
    /// result:             16e18 * 1e18 / 8e18 -> 2e18 = 2 LINK

    /// USDC Example:
    /// tokenPrice:         1e30 -> $1/USDC, as 1e18 token amount is 1e12 USDC, worth 1e12 USD, or 1e30 with 18 decimals
    /// usdValue:           5e18 -> $5
    /// result:             5e18 * 1e18 / 1e30 -> 5e6 = 5 USDC
    return (usdValue * 1e18) / tokenPrice;
  }
}
