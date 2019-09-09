pragma solidity 0.5.0;

interface LinkExInterface {
  function currentRate() external view returns (uint256);
  function update(uint256 rate) external;
}