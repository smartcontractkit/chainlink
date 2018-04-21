pragma solidity ^0.4.21;

import "../Chainlinked.sol";

contract Consumer is Chainlinked {
  bytes32 internal externalId;
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
    externalId = chainlinkRequest(run);
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    onlyOracle
    checkRequestId(_externalId)
  {
    currentPrice = _data;
  }

  modifier checkRequestId(bytes32 _externalId) {
    require(externalId == _externalId);
    _;
  }

}
