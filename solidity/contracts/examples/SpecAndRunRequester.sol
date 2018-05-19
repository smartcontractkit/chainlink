pragma solidity ^0.4.23;

import "../Chainlinked.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract SpecAndRunRequester is Chainlinked, Ownable {
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
    tasks[0] = "httpget";
    tasks[1] = "jsonparse";
    tasks[2] = "ethint256";
    tasks[3] = "ethtx";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfill(bytes32,bytes32)");
    spec.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    spec.addStringArray("path", path);
    chainlinkRequest(spec, LINK(1));
  }

  function cancelRequest(bytes32 _requestId) public onlyOwner {
    cancelChainlinkRequest(_requestId);
  }

  function fulfill(bytes32 _requestId, bytes32 _price)
    public
    checkChainlinkRequest(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  modifier checkRequestId(bytes32 _requestId) {
    require(_requestId == _requestId);
    _;
  }

}

