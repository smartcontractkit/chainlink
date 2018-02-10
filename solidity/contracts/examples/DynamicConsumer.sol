pragma solidity ^0.4.18;

import "../../contracts/Chainlinked.sol";
import "../../contracts/Oracle.sol";

contract DynamicConsumer is Chainlinked {
  uint256 private nonce;
  bytes32 public currentPrice;

  function DynamicConsumer(address _oracle) public {
    oracle = Oracle(_oracle);
  }

  function requestEthereumPrice(string _currency) public {
    Chainlink.Run memory run = newRun("someJobId", this, "fulfill(uint256,bytes32)");
    run.add("url", "https://etherprice.com/api");
    string[] memory path = new string[](2);
    path[0] = "recent";
    path[1] = _currency;
    run.add("path", path);
    nonce = chainlinkRequest(run);
  }

  function fulfill(uint256 _nonce, bytes32 _data)
    public
    onlyOracle
    checkNonce(_nonce)
  {
    currentPrice = _data;
  }

  modifier checkNonce(uint256 _nonce) {
    require(nonce == _nonce);
    _;
  }

}
