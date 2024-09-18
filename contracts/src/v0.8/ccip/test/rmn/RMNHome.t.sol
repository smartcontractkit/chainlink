// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {RMNHome} from "../../rmn/RMNHome.sol";
import {BaseTest} from "../BaseTest.t.sol";

contract RMNHomeTest is BaseTest {
  RMNHome public s_rmnHome;

  function setUp() public virtual override {
    super.setUp();
    s_rmnHome = new RMNHome();
  }
}
