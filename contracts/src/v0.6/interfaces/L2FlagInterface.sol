// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

interface L2FlagInterface {
    function raiseFlag() external;

    function lowerFlag() external;

    function isRaised() external view returns (bool);

    function setRaisingAccessController(address) external;
}
