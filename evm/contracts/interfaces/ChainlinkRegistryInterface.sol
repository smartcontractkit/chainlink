pragma solidity 0.4.24;

interface ChainlinkRegistryInterface {
  function getChainlinkTokenAddress() external view returns (address);
}