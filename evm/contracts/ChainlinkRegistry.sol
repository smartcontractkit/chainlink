pragma solidity 0.4.24;

contract ChainlinkRegistry {
  address public getChainlinkTokenAddress;

  constructor(address _link) public {
    getChainlinkTokenAddress = _link;
  }
}