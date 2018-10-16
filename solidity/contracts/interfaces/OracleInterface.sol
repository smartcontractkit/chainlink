pragma solidity ^0.4.24;

interface OracleInterface {
  function cancel(bytes32 _externalId) external;
}