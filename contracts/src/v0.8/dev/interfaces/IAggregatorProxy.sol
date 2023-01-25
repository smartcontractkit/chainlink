// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

interface IAggregatorProxy {
  function aggregator() external view returns (address);
}
