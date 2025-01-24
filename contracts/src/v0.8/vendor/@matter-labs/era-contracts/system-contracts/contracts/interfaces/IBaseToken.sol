// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

interface IBaseToken {
    function balanceOf(uint256) external view returns (uint256);

    function transferFromTo(address _from, address _to, uint256 _amount) external;

    function totalSupply() external view returns (uint256);

    function name() external pure returns (string memory);

    function symbol() external pure returns (string memory);

    function decimals() external pure returns (uint8);

    function mint(address _account, uint256 _amount) external;

    function withdraw(address _l1Receiver) external payable;

    function withdrawWithMessage(address _l1Receiver, bytes calldata _additionalData) external payable;

    event Mint(address indexed account, uint256 amount);

    event Transfer(address indexed from, address indexed to, uint256 value);

    event Withdrawal(address indexed _l2Sender, address indexed _l1Receiver, uint256 _amount);

    event WithdrawalWithMessage(
        address indexed _l2Sender,
        address indexed _l1Receiver,
        uint256 _amount,
        bytes _additionalData
    );
}
