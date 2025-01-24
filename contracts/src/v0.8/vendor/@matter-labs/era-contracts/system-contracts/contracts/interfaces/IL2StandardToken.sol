// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

interface IL2StandardToken {
    event BridgeMint(address indexed _account, uint256 _amount);

    event BridgeBurn(address indexed _account, uint256 _amount);

    function bridgeMint(address _account, uint256 _amount) external;

    function bridgeBurn(address _account, uint256 _amount) external;

    function l1Address() external view returns (address);

    function l2Bridge() external view returns (address);
}
