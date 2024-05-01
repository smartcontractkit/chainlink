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

  address private s_owner = makeAddr("owner");
  address internal s_alice = makeAddr("alice");
  address internal s_bob = makeAddr("bob");

  address public s_sender1 = makeAddr("sender1");
  address public s_sender2 = makeAddr("sender2");
  address public s_sender3 = makeAddr("sender3");

  function _setUp() internal {
    _deploy();
  }

  function _deploy() internal {
    s_link = new LinkToken();
    s_factory = new OperatorFactory(address(s_link));

    s_mockReceiver = new MockReceiver();
  }
}
