// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MaybeRevertMessageReceiver} from "../../helpers/receivers/MaybeRevertMessageReceiver.sol";

import {Test} from "forge-std/Test.sol";

contract ERC165CheckerRevertingSetup is Test {
  address internal s_receiver;

  bytes4 internal constant EXAMPLE_INTERFACE_ID = 0xdeadbeef;
  bytes4 internal constant NOT_ENOUGH_GAS_SIG = 0x161c3bf7;

  constructor() {
    s_receiver = address(new MaybeRevertMessageReceiver(false));
  }
}
