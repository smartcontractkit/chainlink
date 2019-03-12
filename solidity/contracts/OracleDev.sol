pragma solidity 0.4.24;

import "./Oracle.sol";
import "./interfaces/LinkExInterface.sol";

contract OracleDev is Oracle {
  LinkExInterface internal usdFeed;
  LinkExInterface internal ethFeed;

  constructor(
    address _link,
    address _usdPriceFeed,
    address _ethPriceFeed
  ) public Oracle(_link) {
    usdFeed = LinkExInterface(_usdPriceFeed);
    ethFeed = LinkExInterface(_ethPriceFeed);
  }

  function getEthPriceFeed() public view returns (address) {
    return address(ethFeed);
  }

  function getUsdPriceFeed() public view returns (address) {
    return address(usdFeed);
  }
}