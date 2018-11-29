pragma solidity 0.4.24;

import "../interfaces/OracleInterface.sol";

contract EmptyOracle is OracleInterface {

  function cancel(bytes32) external {}
  function fulfillData(uint256, bytes32) external returns (bool) {}
  function getAuthorizationStatus(address) external view returns (bool) { return false; }
  function onTokenTransfer(address, uint256, bytes) external pure {}
  function requestData(address, uint256, uint256, bytes32, address, bytes4, bytes32, bytes) external {}
  function setFulfillmentPermission(address, bool) external {}
  function withdraw(address, uint256) external {}

}
