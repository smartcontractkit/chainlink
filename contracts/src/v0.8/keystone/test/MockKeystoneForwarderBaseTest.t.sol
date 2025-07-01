// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {Receiver} from "./mocks/Receiver.sol";
import {MockKeystoneForwarder} from "../MockKeystoneForwarder.sol";

contract MockBaseTest is Test {
  address internal ADMIN = address(1);
  address internal constant TRANSMITTER = address(50);
  uint32 internal DON_ID = 0x01020304;
  uint32 internal CONFIG_VERSION = 1;

  MockKeystoneForwarder internal s_mockForwarder;
  Receiver internal s_receiver;

  function setUp() public virtual {
    vm.startPrank(ADMIN);
    s_mockForwarder = new MockKeystoneForwarder();
    s_receiver = new Receiver();
  }
}