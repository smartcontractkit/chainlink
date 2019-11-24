pragma solidity ^0.5.0;

interface WithdrawalInterface {
  /**
   * @notice transfer LINK to another address
   * @param recipient is the address to send the LINK to
   * @param amount is the amount of LINK to send
   */
  function withdraw(address recipient, uint256 amount) external;

  /**
   * @notice query the available amount of LINK to withdraw
   */
  function withdrawable() external view returns (uint256);
}
