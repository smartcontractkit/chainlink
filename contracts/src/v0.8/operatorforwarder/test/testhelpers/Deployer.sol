// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

import {AuthorizedForwarder} from "../../AuthorizedForwarder.sol";
import {Operator} from "../../Operator.sol";
import {OperatorFactory} from "../../OperatorFactory.sol";

import {MockReceiver} from "./MockReceiver.sol";

import {MockLinkToken} from "../../../mocks/MockLinkToken.sol";

abstract contract Deployer is Test {
  OperatorFactory public factory;
  Operator public operator;
  AuthorizedForwarder public forwarder;

  MockLinkToken public link;
  MockReceiver public mockReceiver;

  address public owner = address(uint160(uint256(keccak256("owner"))));
  address public alice = address(uint160(uint256(keccak256("alice"))));
  address public bob = address(uint160(uint256(keccak256("bob"))));

  address public sender1 = address(uint160(uint256(keccak256("sender1"))));
  address public sender2 = address(uint160(uint256(keccak256("sender2"))));
  address public sender3 = address(uint160(uint256(keccak256("sender3"))));

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
