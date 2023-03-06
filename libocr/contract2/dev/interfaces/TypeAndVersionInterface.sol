// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface TypeAndVersionInterface{
  function typeAndVersion()
    external
    pure
    returns (string memory);
}