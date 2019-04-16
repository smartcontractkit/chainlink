pragma solidity 0.4.24;

import "./Oracle.sol";
import "./interfaces/LinkExInterface.sol";

contract OracleDev is Oracle {

  LinkExInterface internal priceFeed;

  mapping(bytes32 => LinkExInterface) public priceFeeds;

  constructor(address _link) public Oracle(_link) {} // solium-disable-line no-empty-blocks

  function currentRate(bytes32 _currency) public view returns (uint256) {
    return priceFeeds[_currency].currentRate();
  }

  function setPriceFeed(address _priceFeed, bytes32 _currency) external onlyOwner {
    priceFeeds[_currency] = LinkExInterface(_priceFeed);
  }
}