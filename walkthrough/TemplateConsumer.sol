pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract TemplateConsumer is Chainlinked {
  bytes32 internal specId;
  bytes32 public currentPrice;

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    specId = _specId;
  }

  function requestEthereumPrice(string _currency) public {
    chainlinkLib.Run memory run = newRun(specId, this, "reportPrice(bytes32,uint256)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    run.add("path", _currency); // "USD"
    chainlinkRequest(run, LINK(1));
  }

  function reportPrice(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    // ...
  }

}
