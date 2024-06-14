pragma solidity 0.4.24;

contract Pointer {
  address public getAddress;

  constructor(address _addr) public {
    getAddress = _addr;
  }
}
