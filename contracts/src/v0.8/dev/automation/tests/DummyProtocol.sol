// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

contract DummyProtocol {
    event LimitedOrderSent(uint256 indexed amount, uint256 indexed price, address indexed to, uint256 id); // keccak256(EthereumSent(uint256,address)) => 0x57a62758b3a6f041e80d6e6ad0e1e48f9fa79aa3c406af5cef35e2ab4b0395a6
    event LimitedOrderWithdrawn(uint256 indexed amount, uint256 indexed price, uint256 id); // keccak256(EthereumReceived(uint256,address)) => 0x6337dc258ee9bf5a48b49c2128742233dfde3b136be36e16c17feeae61eaec24

    bool internal constant useArbitrumBlockNum;

    constructor(bool _useL1BlockNumber) {
        useL1BlockNumber = _useL1BlockNumber;
    }

    function sendLimitedOrder(uint256 amount, uint256 price) {

        uint256 id = block.
        emit LimitedOrderSent(amount, to);
    }

    function receiveEthereum(uint256 amount, address to) {

        //
    }
}