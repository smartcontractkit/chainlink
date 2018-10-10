pragma solidity ^0.4.24;

interface ILinkToken {
  function balanceOf(address _owner) external returns (uint256 balance);
  function transfer(address _to, uint _value) external returns (bool success);
  function transferAndCall(address _to, uint _value, bytes _data) external returns (bool success);
}