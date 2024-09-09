// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BaseValidator} from "../../../dev/base/BaseValidator.sol";
import {MockBaseValidator} from "../../mocks/MockBaseValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract BaseValidatorTest is L2EPTest {
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = address(0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b);
  address internal constant DUMMY_L1_XDOMAIN_MSNGR_ADDR = address(0xa04Fc18f012B1a5A8231c7Ee4b916Dd6dbd271b6);
  address internal constant DUMMY_L2_UPTIME_FEED_ADDR = address(0xFe31891940A2e5f04B76eD8bD1038E44127d1512);
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

contract BaseValidator_Constructor is BaseValidatorTest {
  /// @notice it correctly validates that the L1 bridge address is not zero
  function test_ConstructingRevertedWithZeroL1BridgeAddress() public {
    vm.expectRevert(BaseValidator.L1CrossDomainMessengerAddressZero.selector);
    new MockBaseValidator(address(0), DUMMY_L2_UPTIME_FEED_ADDR, INIT_GAS_LIMIT);
  }

  /// @notice it correctly validates that the L2 Uptime feed address is not zero
  function test_ConstructingRevertedWithZeroL2UpdateFeedAddress() public {
    vm.expectRevert(BaseValidator.L2UptimeFeedAddrZero.selector);
    new MockBaseValidator(DUMMY_L1_XDOMAIN_MSNGR_ADDR, address(0), INIT_GAS_LIMIT);
  }
}

contract BaseValidator_GetAndSetGasLimit is BaseValidatorTest {
  function test_CorrectlyGetsGasLimit() public {
    assertEq(s_baseValidator.getGasLimit(), INIT_GAS_LIMIT);

    uint32 newGasLimit = INIT_GAS_LIMIT + 1;

    vm.expectEmit();
    emit BaseValidator.GasLimitUpdated(newGasLimit);
    s_baseValidator.setGasLimit(newGasLimit);

    assertEq(s_baseValidator.getGasLimit(), newGasLimit);
  }
}
