pragma solidity ^0.4.24;

interface OracleInterface {
  function fulfillData(uint256 _internalId, bytes32 _data) external returns (bool);
  function cancel(bytes32 _externalId) external;
}