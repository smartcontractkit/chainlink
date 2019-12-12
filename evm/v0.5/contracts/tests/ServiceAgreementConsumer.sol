pragma solidity 0.5.0;

import "../ChainlinkClient.sol";

contract ServiceAgreementConsumer is ChainlinkClient {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK;

  bytes32 internal sAId;
  bytes32 public currentPrice;

  constructor(address _link, address _coordinator, bytes32 _sAId) public {
    setChainlinkToken(_link);
    setChainlinkOracle(_coordinator);
    sAId = _sAId;
  }

  function requestEthereumPrice(string memory _currency) public {
    Chainlink.Request memory req = buildChainlinkRequest(sAId, address(this), this.fulfill.selector);
    req.add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    req.add("path", _currency);
    sendChainlinkRequest(req, ORACLE_PAYMENT);
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    recordChainlinkFulfillment(_requestId)
  {
    currentPrice = _price;
  }
}
