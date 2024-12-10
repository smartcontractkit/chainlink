// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BaseValidator} from "../../../base/BaseValidator.sol";
import {MockBaseValidator} from "../../mocks/MockBaseValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract BaseValidator_Setup is L2EPTest {
  address internal immutable L2_SEQ_STATUS_RECORDER_ADDRESS = makeAddr("L2_SEQ_STATUS_RECORDER_ADDRESS");
  address internal immutable DUMMY_L1_XDOMAIN_MSNGR_ADDR = makeAddr("DUMMY_L1_XDOMAIN_MSNGR_ADDR");
  address internal immutable DUMMY_L2_UPTIME_FEED_ADDR = makeAddr("DUMMY_L2_UPTIME_FEED_ADDR");
  uint32 internal constant INIT_GAS_LIMIT = 1900000;

  BaseValidator internal s_baseValidator;

  /// Fake event that will get emitted when `requestL2TransactionDirect` is called
  /// Definition is taken from MockZKSyncL1Bridge
  event SentMessage(address indexed sender, bytes message);

  /// Setup
  function setUp() public {
    s_baseValidator = new MockBaseValidator(
      DUMMY_L1_XDOMAIN_MSNGR_ADDR,
      L2_SEQ_STATUS_RECORDER_ADDRESS,
      INIT_GAS_LIMIT
    );
  }
}

contract BaseValidator_Constructor is BaseValidator_Setup {
  /// @notice Reverts when L1 bridge address is zero
  function test_Constructor_RevertWhen_L1BridgeAddressIsZero() public {
    vm.expectRevert(BaseValidator.L1CrossDomainMessengerAddressZero.selector);
    new MockBaseValidator(address(0), DUMMY_L2_UPTIME_FEED_ADDR, INIT_GAS_LIMIT);
  }

  /// @notice Reverts when L2 Uptime feed address is zero
  function test_Constructor_RevertWhen_L2UptimeFeedAddressIsZero() public {
    vm.expectRevert(BaseValidator.L2UptimeFeedAddrZero.selector);
    new MockBaseValidator(DUMMY_L1_XDOMAIN_MSNGR_ADDR, address(0), INIT_GAS_LIMIT);
  }
}

contract BaseValidator_GetAndSetGasLimit is BaseValidator_Setup {
  /// @notice Verifies the correct retrieval and update of the gas limit
  function test_GetAndSetGasLimit_CorrectlyHandlesGasLimit() public {
    assertEq(s_baseValidator.getGasLimit(), INIT_GAS_LIMIT);

    uint32 newGasLimit = INIT_GAS_LIMIT + 1;

    vm.expectEmit();
    emit BaseValidator.GasLimitUpdated(newGasLimit);
    s_baseValidator.setGasLimit(newGasLimit);

    assertEq(s_baseValidator.getGasLimit(), newGasLimit);
  }
}
