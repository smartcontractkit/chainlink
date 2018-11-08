pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract ServiceAgreementConsumer is Chainlinked {
  // NB: `Chainlinked` covers both the single-Oracle and Service-Agreement
  // frameworks. Some of the terminology there is in terms of a single Oracle.
  bytes32 internal sAId;
  bytes32 public currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId,
    bytes32 indexed data
  );

  constructor(address _link, address _coordinator, bytes32 _sAId) public {
    setLinkToken(_link);
    setOracle(_coordinator);
    sAId = _sAId;
  }

  function requestEthereumPrice(string _currency) public {
    ChainlinkLib.Run memory run = newRun(sAId, this, "fulfill(bytes32,bytes32)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    run.addStringArray("path", path);
    serviceRequest(run, LINK(1));
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    // TODO: Example which deals with multiple responses. (Here, or in Coordinator?)
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  function cancelRequest(bytes32 _requestId) public {
    cancelChainlinkRequest(_requestId);
  }

  function withdrawLink() public {
    LinkTokenInterface link = LinkTokenInterface(chainlinkToken());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }
}
