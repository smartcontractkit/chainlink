pragma solidity ^0.4.18;

import "../../contracts/ChainLink.sol";

contract Consumer {
  ChainLink private chainLink;
  uint256 private nonce;
  bytes32 public currentPrice;

  function Consumer(address _chainLink) public {
    chainLink = ChainLink(_chainLink);
  }

  function requestEthereumPrice() public {
    bytes4 fid = bytes4(keccak256("fulfill(uint256,bytes32)"));
    nonce = chainLink.requestData(this, fid, "https://etherprice.com/api", "recent,usd");
  }

  function fulfill(uint256 _nonce, bytes32 _data)
    public
    checkSender
    checkNonce(_nonce)
  {
    currentPrice = _data;
  }

  modifier checkSender() {
    require(msg.sender == address(chainLink));
    _;
  }

  modifier checkNonce(uint256 _nonce) {
    require(nonce == _nonce);
    _;
  }
}
