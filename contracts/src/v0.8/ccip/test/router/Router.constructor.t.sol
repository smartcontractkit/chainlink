// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

import {Router} from "../../Router.sol";
import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {IRouter} from "../../interfaces/IRouter.sol";
import {IRouterClient} from "../../interfaces/IRouterClient.sol";
import {IWrappedNative} from "../../interfaces/IWrappedNative.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";

import {OnRamp} from "../../onRamp/OnRamp.sol";
import {OnRampSetup} from "../onRamp/onRamp/OnRampSetup.t.sol";

contract Router_constructor is OnRampSetup {
  function test_Constructor_Success() public view {
    assertEq("Router 1.2.0", s_sourceRouter.typeAndVersion());
    assertEq(OWNER, s_sourceRouter.owner());
  }
}
