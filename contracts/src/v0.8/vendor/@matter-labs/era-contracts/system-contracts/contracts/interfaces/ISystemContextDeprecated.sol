// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice The interface with deprecated functions of the SystemContext contract. It is aimed for backward compatibility.
 */
interface ISystemContextDeprecated {
    function currentBlockInfo() external view returns (uint256);

    function getBlockNumberAndTimestamp() external view returns (uint256 blockNumber, uint256 blockTimestamp);

    function blockHash(uint256 _blockNumber) external view returns (bytes32 hash);
}
