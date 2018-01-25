pragma solidity ^0.4.18;

import "../../contracts/ChainLinked.sol";
import "../../contracts/Oracle.sol";

contract Consumer is ChainLinked {
  uint256 private nonce;
  bytes32 public currentPrice;

  function Consumer(address _oracle) public {
    oracle = Oracle(_oracle);
  }

  function requestEthereumPrice() public {
    ChainLink.Run memory run = newRun("1234", this, "fulfill(uint256,bytes32)");
    run.add("url", "https://etherprice.com/api");
    run.add("path", "recent,usd");
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
