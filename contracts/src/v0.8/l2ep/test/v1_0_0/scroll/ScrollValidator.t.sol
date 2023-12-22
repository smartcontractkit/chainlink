// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {MockScrollL1CrossDomainMessenger} from "../../mocks/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/MockScrollL2CrossDomainMessenger.sol";
import {ScrollSequencerUptimeFeed} from "../../../dev/scroll/ScrollSequencerUptimeFeed.sol";
import {ScrollValidator} from "../../../dev/scroll/ScrollValidator.sol";
import {L2EPTest} from "../L2EPTest.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/scroll/ScrollValidator.t.sol
//
contract ScrollValidatorTest is L2EPTest {
  /// Sets a fake L2 target and the initial gas limit
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = 0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b;
  uint32 internal constant INIT_GAS_LIMIT = 1900000;

  /// L2EP contracts
  MockScrollL1CrossDomainMessenger internal s_mockScrollL1CrossDomainMessenger;
  MockScrollL2CrossDomainMessenger internal s_mockScrollL2CrossDomainMessenger;
  ScrollSequencerUptimeFeed internal s_scrollSequencerUptimeFeed;
  ScrollValidator internal s_scrollValidator;

  /// https://github.com/scroll-tech/scroll/blob/03089eaeee1193ff44c532c7038611ae123e7ef3/contracts/src/libraries/IScrollMessenger.sol#L22
  event SentMessage(
    address indexed sender,
    address indexed target,
    uint256 value,
    uint256 messageNonce,
    uint256 gasLimit,
    bytes message
  );

  /// Setup
  function setUp() public {
    s_mockScrollL1CrossDomainMessenger = new MockScrollL1CrossDomainMessenger();
    s_mockScrollL2CrossDomainMessenger = new MockScrollL2CrossDomainMessenger();

    s_scrollSequencerUptimeFeed = new ScrollSequencerUptimeFeed(
      address(s_mockScrollL1CrossDomainMessenger),
      address(s_mockScrollL2CrossDomainMessenger),
      true
    );

    s_scrollValidator = new ScrollValidator(
      address(s_mockScrollL1CrossDomainMessenger),
      address(s_scrollSequencerUptimeFeed),
      INIT_GAS_LIMIT
    );
  }
}

contract ScrollValidator_SetAndGetGasLimit is ScrollValidatorTest {
  function test_SetAndGetGasLimit() public {
    uint32 newGasLimit = 2000000;

    assertEq(s_scrollValidator.getGasLimit(), INIT_GAS_LIMIT);
    s_scrollValidator.setGasLimit(newGasLimit);
    assertEq(s_scrollValidator.getGasLimit(), newGasLimit);
  }
}

contract ScrollValidator_CheckValidateAccessControl is ScrollValidatorTest {
  function test_CheckValidateAccessControl() public {
    address strangerAddress = vm.addr(0x2);

    vm.startPrank(strangerAddress, strangerAddress);
    vm.expectRevert("No access");
    s_scrollValidator.validate(0, 0, 1, 1);
    vm.stopPrank();
  }
}

contract ScrollValidator_SequencerOnline is ScrollValidatorTest {
  function test_CheckValidateAccessControl() public {
    // Gives access to the eoaValidator
    address eoaValidator = vm.addr(0x1);
    s_scrollValidator.addAccess(eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.startPrank(eoaValidator, eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      address(s_scrollValidator), // sender
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      0, // value
      0, // nonce
      INIT_GAS_LIMIT, // gas limit
      abi.encodeWithSelector(ScrollSequencerUptimeFeed.updateStatus.selector, false, futureTimestampInSeconds) // message
    );

    // Runs the function (which produces the event to test)
    s_scrollValidator.validate(0, 0, 0, 0);
    vm.stopPrank();
  }
}

contract ScrollValidator_SequencerOffline is ScrollValidatorTest {
  function test_CheckValidateAccessControl() public {
    // Gives access to the eoaValidator
    address eoaValidator = vm.addr(0x1);
    s_scrollValidator.addAccess(eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 10000;
    vm.startPrank(eoaValidator, eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      address(s_scrollValidator), // sender
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      0, // value
      0, // nonce
      INIT_GAS_LIMIT, // gas limit
      abi.encodeWithSelector(ScrollSequencerUptimeFeed.updateStatus.selector, true, futureTimestampInSeconds) // message
    );

    // Runs the function (which produces the event to test)
    s_scrollValidator.validate(0, 0, 1, 1);
    vm.stopPrank();
  }
}
