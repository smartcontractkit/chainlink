pragma solidity 0.4.24;

import "./Oracle.sol";
import "./interfaces/LinkExInterface.sol";

contract OracleDev is Oracle {

  mapping(bytes32 => address) public priceFeeds;

  constructor(address _link) public Oracle(_link) {}

  function currentRate(bytes32 _currency) public view returns (uint256) {
    LinkExInterface priceFeed = LinkExInterface(priceFeeds[_currency]);
    return priceFeed.currentRate();
  }

  function setPriceFeed(address _priceFeed, bytes32 _currency) external onlyOwner {
    priceFeeds[_currency] = _priceFeed;
  }
}