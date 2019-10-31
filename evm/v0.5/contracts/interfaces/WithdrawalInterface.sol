pragma solidity ^0.5.0;

interface WithdrawalInterface {
  function withdraw(address recipient, uint256 amount) external;
  function withdrawable() external view returns (uint256);
}
