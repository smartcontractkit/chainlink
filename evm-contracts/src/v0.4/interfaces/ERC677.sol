pragma solidity ^0.4.8;

import "./linkERC20.sol";

contract ERC677 is linkERC20 {
  function transferAndCall(address to, uint value, bytes data) returns (bool success);

  event Transfer(address indexed from, address indexed to, uint value, bytes data);
}
