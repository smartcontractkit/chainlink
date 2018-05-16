pragma solidity ^0.4.23;

import "../Chainlinked.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract SpecAndRunRequester is Chainlinked, Ownable {
  bytes32 internal requestId;
  bytes32 public currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId,
    bytes32 indexed price
  );

  constructor(address _link, address _oracle) Ownable() public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function requestEthereumPrice(string _currency) public {
    string[] memory tasks = new string[](4);
    tasks[0] = "httpGet";
    tasks[1] = "jsonParse";
    tasks[2] = "ethint256";
    tasks[3] = "ethtx";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfill(bytes32,bytes32)");
    spec.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    spec.addStringArray("path", path);
    requestId = chainlinkRequest(spec, LINK(1));
  }

  function cancelRequest() public onlyOwner {
    oracle.cancel(requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  modifier checkRequestId(bytes32 _requestId) {
    require(requestId == _requestId);
    _;
  }

}

