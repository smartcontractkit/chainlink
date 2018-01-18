pragma solidity ^0.4.18;

import "../../contracts/ChainLink.sol";

contract Consumer {
  ChainLink private chainLink;
  uint256 private nonce;

  function Consumer(address _chainLink) public {
    chainLink = ChainLink(_chainLink);
  }

  function requestEthereumPrice() public {
    bytes4 fid = bytes4(keccak256("callback(uint256,bytes32)"));
    nonce = chainLink.requestData(this, fid, "https://etherprice.com/api", "recent,usd");
  }

  function callback(uint256 _nonce, bytes32 _data) public {
  }
}
