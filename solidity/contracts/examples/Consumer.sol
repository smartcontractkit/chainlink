pragma solidity ^0.4.23;

import "../Chainlinked.sol";

contract Consumer is Chainlinked {
  bytes32 internal requestId;
  bytes32 public currentPrice;

  function Consumer(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function requestEthereumPrice(string _currency)
    public
  {
    ChainlinkLib.Run memory run = newRun("someJobId", this, "fulfill(bytes32,bytes32)");
    run.add("url", "https://etherprice.com/api");
    string[] memory path = new string[](2);
    path[0] = "recent";
    path[1] = _currency;
    run.addStringArray("path", path);
    requestId = chainlinkRequest(run, 1 szabo);
  }

  function fulfill(bytes32 _requestId, bytes32 _data)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    currentPrice = _data;
  }

  modifier checkRequestId(bytes32 _requestId) {
    require(requestId == _requestId);
    _;
  }

}
