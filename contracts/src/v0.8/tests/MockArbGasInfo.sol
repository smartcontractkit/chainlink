pragma solidity 0.8.6;

contract MockArbGasInfo {
  function getCurrentTxL1GasFees() external view returns (uint256) {
    return 1000000;
  }
}
