// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

interface IAccountCodeStorage {
    function storeAccountConstructingCodeHash(address _address, bytes32 _hash) external;

    function storeAccountConstructedCodeHash(address _address, bytes32 _hash) external;

    function markAccountCodeHashAsConstructed(address _address) external;

    function getRawCodeHash(address _address) external view returns (bytes32 codeHash);

    function getCodeHash(uint256 _input) external view returns (bytes32 codeHash);

    function getCodeSize(uint256 _input) external view returns (uint256 codeSize);
}
