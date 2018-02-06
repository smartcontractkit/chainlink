pragma solidity ^0.4.18;

import "../../contracts/Chainlinked.sol";
import "../../contracts/Oracle.sol";

contract SimpleConsumer is Chainlinked {
  uint256 private nonce;
  bytes32 public currentPrice;

  function SimpleConsumer(address _oracle) public {
    oracle = Oracle(_oracle);
  }

  function requestEthereumPrice() public {
    var fid = bytes4(keccak256("fulfill(uint256,bytes32)"));
    var data = '{"url":"https://etherprice.com/api","path":["recent","usd"]}';
    nonce = oracle.requestData("someJobId", this, fid, data);
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
