pragma solidity ^0.4.8;


contract ERC677Receiver {
  function onTokenTransfer(address _sender, uint _value, bytes _data);
}
