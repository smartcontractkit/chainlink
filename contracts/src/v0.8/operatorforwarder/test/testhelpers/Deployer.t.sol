// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

import {MockReceiver} from "./MockReceiver.sol";
import {AuthorizedForwarder} from "../../AuthorizedForwarder.sol";
import {Operator} from "../../Operator.sol";
import {OperatorFactory} from "../../OperatorFactory.sol";
import {LinkToken} from "../../../shared/token/ERC677/LinkToken.sol";

abstract contract Deployer is Test {
  OperatorFactory internal s_factory;
  LinkToken internal s_link;
  MockReceiver internal s_mockReceiver;

  address internal constant ALICE = address(0x101);
  address internal constant BOB = address(0x102);
  address internal constant SENDER_1 = address(0x103);
  address internal constant SENDER_2 = address(0x104);
  address internal constant SENDER_3 = address(0x105);

  function _setUp() internal {
    _deploy();
  }

  function _deploy() internal {
    s_link = new LinkToken();
    s_factory = new OperatorFactory(address(s_link));

    s_mockReceiver = new MockReceiver();
  }
}
