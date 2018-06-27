pragma solidity ^0.4.24;

import "../Chainlinked.sol";
import "../lib/Ownable.sol";

contract SpecAndRunRequester is Chainlinked, Ownable {
  uint256 public currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId,
    uint256 indexed price
  );

  constructor(address _link, address _oracle) Ownable() public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function requestEthereumPrice(string _currency) public {
    string[] memory tasks = new string[](5);
    tasks[0] = "httpget";
    tasks[1] = "jsonparse";
    tasks[2] = "multiply";
    tasks[3] = "ethuint256";
    tasks[4] = "ethtx";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfill(bytes32,uint256)");
    spec.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    spec.addStringArray("path", path);
    spec.addInt("times", 100);
    chainlinkRequest(spec, LINK(1));
  }

  function fulfill(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  function withdrawLink() onlyOwner public {
    require(link.transfer(owner, link.balanceOf(address(this))), "Unable to transfer");
  }

}

