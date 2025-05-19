// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Client} from "../../../libraries/Client.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_ccipReceive is OffRampSetup {
  function test_RevertWhen_Always() public {
    Client.Any2EVMMessage memory message;

    vm.expectRevert();

    s_offRamp.ccipReceive(message);
  }
}
