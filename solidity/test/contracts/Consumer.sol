pragma solidity ^0.4.18;

import "../../contracts/Oracle.sol";
import "../../contracts/ChainLinked.sol";

contract Consumer is ChainLinked {
  Oracle private oracle;
  uint256 private nonce;
  bytes32 public currentPrice;

  function Consumer(address _oracle) public {
    oracle = Oracle(_oracle);
  }

  function requestEthereumPrice() public {
    bytes4 fid = bytes4(keccak256("fulfill(uint256,bytes32)"));
    Json.Params memory ps;
    ps.add("url", "https://etherprice.com/api");
    ps.add("path", "recent,usd");
    nonce = oracle.requestData(this, fid, ps.close());
  }

  function fulfill(uint256 _nonce, bytes32 _data)
    public
    checkSender
    checkNonce(_nonce)
  {
    currentPrice = _data;
  }

  modifier checkSender() {
    require(msg.sender == address(oracle));
    _;
  }

  modifier checkNonce(uint256 _nonce) {
    require(nonce == _nonce);
    _;
  }
}
