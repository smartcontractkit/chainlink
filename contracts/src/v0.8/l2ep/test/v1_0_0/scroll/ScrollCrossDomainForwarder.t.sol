// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {MockScrollL1CrossDomainMessenger} from "../../mocks/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/MockScrollL2CrossDomainMessenger.sol";
import {ScrollCrossDomainForwarder} from "../../../dev/scroll/ScrollCrossDomainForwarder.sol";
import {ScrollValidator} from "../../../dev/scroll/ScrollValidator.sol";
import {L2EPTest} from "../L2EPTest.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/scroll/ScrollCrossDomainForwarder.t.sol
//
contract ScrollCrossDomainForwarderTest is L2EPTest {
  /// L2EP contracts
  MockScrollL1CrossDomainMessenger internal s_mockScrollL1CrossDomainMessenger;
  MockScrollL2CrossDomainMessenger internal s_mockScrollL2CrossDomainMessenger;
  ScrollCrossDomainForwarder internal s_scrollCrossDomainForwarder;

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_mockScrollL1CrossDomainMessenger = new MockScrollL1CrossDomainMessenger();
    s_mockScrollL2CrossDomainMessenger = new MockScrollL2CrossDomainMessenger();

    // TODO:
    // s_scrollCrossDomainForwarder = new ScrollCrossDomainForwarder();
  }
}

/// @notice it should have been deployed with the correct initial state
contract ScrollCrossDomainForwarder_CheckInitialState is ScrollCrossDomainForwarderTest {
  function test_CheckInitialState() public {}
}
