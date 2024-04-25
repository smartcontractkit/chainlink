// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";

import "../../AuthorizedForwarder.sol";
import "../../Operator.sol";
import "../../OperatorFactory.sol";

import "./MockReceiver.sol";

import "../../../mocks/MockLinkToken.sol";

abstract contract Deployer is Test {

    OperatorFactory factory;
    Operator operator;
    AuthorizedForwarder forwarder;

    MockLinkToken public link;
    MockReceiver mockReceiver;

    address owner = address(uint160(uint256(keccak256("owner"))));
    address alice = address(uint160(uint256(keccak256("alice"))));
    address bob = address(uint160(uint256(keccak256("bob"))));

    address sender1 = address(uint160(uint256(keccak256("sender1"))));
    address sender2 = address(uint160(uint256(keccak256("sender2"))));
    address sender3 = address(uint160(uint256(keccak256("sender3"))));

    function _setUp() internal {
        _deploy();
    }

    function _deploy() internal {
        vm.startPrank(owner);

        link = new MockLinkToken();
        factory = new OperatorFactory(address(link));

        mockReceiver = new MockReceiver();
        
        vm.stopPrank();

    }
}