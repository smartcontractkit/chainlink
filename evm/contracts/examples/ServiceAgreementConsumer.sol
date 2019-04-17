pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract ServiceAgreementConsumer is Chainlinked {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  bytes32 internal sAId;
  bytes32 public currentPrice;

  constructor(address _link, address _coordinator, bytes32 _sAId) public {
    setLinkToken(_link);
    setOracle(_coordinator);
    sAId = _sAId;
  }

  function requestEthereumPrice(string _currency) public {
    Chainlink.Request memory req = newRequest(sAId, this, this.fulfill.selector);
    req.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    req.add("path", _currency);
    chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    recordChainlinkFulfillment(_requestId)
  {
    currentPrice = _price;
  }
}
