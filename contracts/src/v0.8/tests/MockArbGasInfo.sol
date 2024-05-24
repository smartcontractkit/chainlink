pragma solidity 0.8.6;

contract MockArbGasInfo {
  function getCurrentTxL1GasFees() external view returns (uint256) {
    return 1000000;
  }

  function getPricesInWei() external view returns (uint256, uint256, uint256, uint256, uint256, uint256) {
    return (0, 1000, 0, 0, 0, 0);
  }
}
