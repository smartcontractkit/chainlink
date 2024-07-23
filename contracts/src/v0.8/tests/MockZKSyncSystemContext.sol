pragma solidity 0.8.6;

contract MockZKSyncSystemContext {
  function gasPrice() external view returns (uint256) {
    return 250000000; // 0.25 gwei
  }

  function gasPerPubdataByte() external view returns (uint256) {
    return 500;
  }

  function getCurrentPubdataSpent() external view returns (uint256 currentPubdataSpent) {
    return 1000;
  }
}
