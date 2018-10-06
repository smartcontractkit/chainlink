pragma solidity ^0.4.24;

contract BoringConsumer {
  bytes32 internal specId;
  bytes32 public currentPrice;

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    specId = _specId;
  }

  function requestEthereumPrice(string _currency) public {
    chainlinkLib.Run memory run = newRun(specId, this, "fulfill(bytes32,bytes32)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    run.addStringArray("path", _currency); // "USD"
    chainlinkRequest(run, LINK(1));
  }

  function cancelRequest(bytes32 _requestId) public {
    cancelChainlinkRequest(_requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

}
