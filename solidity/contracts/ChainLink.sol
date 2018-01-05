pragma solidity ^0.4.18;

import "./zeppelin/Ownable.sol";
contract ChainLink is Ownable {

  bytes32 public value;
  uint public nonce;
  event Request(
    uint indexed nonce,
    address indexed to,
    bytes4 indexed fid
  );

  function ChainLink() public {
    value = "Hello World!";
  }

  function setValue(bytes32 _value) public {
    value = _value;
  }

  function requestData(address _address, bytes4 _fid) public {
    Request(nonce, _address, _fid);
    nonce += 1;
  }
}
