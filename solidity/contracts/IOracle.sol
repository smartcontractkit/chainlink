pragma solidity ^0.4.24;

interface IOracle {
  function cancel(bytes32 _externalId) external;
}