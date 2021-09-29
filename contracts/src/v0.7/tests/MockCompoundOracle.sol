// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../interfaces/UniswapAnchoredView.sol";

contract MockCompoundOracle is UniswapAnchoredView {
  struct OracleDetails {
    uint256 price;
    uint256 decimals;
  }

  mapping(string => OracleDetails) s_oracleDetails;

  function price(string memory symbol) external view override returns (uint256) {
    return s_oracleDetails[symbol].price;
  }

  function setPrice(
    string memory symbol,
    uint256 newPrice,
    uint256 newDecimals
  ) public {
    OracleDetails memory details = s_oracleDetails[symbol];
    details.price = newPrice;
    details.decimals = newDecimals;
    s_oracleDetails[symbol] = details;
  }
}
