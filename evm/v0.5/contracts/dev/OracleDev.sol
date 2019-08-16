pragma solidity 0.5.0;

import "../Oracle.sol";
import "./LinkExInterface.sol";

contract OracleDev is Oracle {

  LinkExInterface internal priceFeed;

  mapping(bytes32 => LinkExInterface) public priceFeeds;

  constructor(address _link) public Oracle(_link) {} // solhint-disable-line no-empty-blocks

  function currentRate(bytes32 _currency) public view returns (uint256) {
    return priceFeeds[_currency].currentRate();
  }

  function setPriceFeed(address _priceFeed, bytes32 _currency) external onlyOwner {
    priceFeeds[_currency] = LinkExInterface(_priceFeed);
  }
}
