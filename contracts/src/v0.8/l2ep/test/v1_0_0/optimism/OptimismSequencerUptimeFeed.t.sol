// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockOptimismL1CrossDomainMessenger} from "../../mocks/MockOptimismL1CrossDomainMessenger.sol";
import {MockOptimismL2CrossDomainMessenger} from "../../mocks/MockOptimismL2CrossDomainMessenger.sol";
import {OptimismSequencerUptimeFeed} from "../../../optimism/OptimismSequencerUptimeFeed.sol";
import {BaseSequencerUptimeFeed} from "../../../base/BaseSequencerUptimeFeed.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract OptimismSequencerUptimeFeed_TestWrapper is OptimismSequencerUptimeFeed {
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) OptimismSequencerUptimeFeed(l1SenderAddress, l2CrossDomainMessengerAddr, initialStatus) {}

  /// @notice Exposes the internal `_validateSender` function for testing
  function validateSenderTestWrapper(address l1Sender) external view {
    super._validateSender(l1Sender);
  }
}

contract OptimismSequencerUptimeFeed_Setup is L2EPTest {
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);

  /// L2EP contracts
  MockOptimismL1CrossDomainMessenger internal s_mockOptimismL1CrossDomainMessenger;
  MockOptimismL2CrossDomainMessenger internal s_mockOptimismL2CrossDomainMessenger;
  OptimismSequencerUptimeFeed_TestWrapper internal s_optimismSequencerUptimeFeed;

  /// Setup
  function setUp() public {
    // Deploy contracts
    s_mockOptimismL1CrossDomainMessenger = new MockOptimismL1CrossDomainMessenger();
    s_mockOptimismL2CrossDomainMessenger = new MockOptimismL2CrossDomainMessenger();
    s_optimismSequencerUptimeFeed = new OptimismSequencerUptimeFeed_TestWrapper(
      s_l1OwnerAddr,
      address(s_mockOptimismL2CrossDomainMessenger),
      false
    );

    // Sets mock sender in mock L2 messenger contract
    s_mockOptimismL2CrossDomainMessenger.setSender(s_l1OwnerAddr);
  }
}

contract OptimismSequencerUptimeFeed_Constructor is OptimismSequencerUptimeFeed_Setup {
  /// @notice Tests the initial state of the contract
  function test_Constructor_InitialState() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    new OptimismSequencerUptimeFeed_TestWrapper(s_l1OwnerAddr, address(s_mockOptimismL2CrossDomainMessenger), false);

    // Checks L1 sender
    address actualL1Addr = s_optimismSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_optimismSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract OptimismSequencerUptimeFeed_ValidateSender is OptimismSequencerUptimeFeed_Setup {
  /// @notice Reverts if called by an address that is not the L2 Cross Domain Messenger
  function test_ValidateSender_RevertWhen_SenderIsNotL2CrossDomainMessengerAddr() public {
    address l2MessengerAddr = address(s_mockOptimismL2CrossDomainMessenger);
    // Sets msg.sender to a different address
    vm.startPrank(s_strangerAddr, l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_optimismSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }

  /// @notice Reverts if the L1 sender address is not the L1 Cross Domain Messenger Sender
  function test_ValidateSender_RevertWhen_L1CrossDomainMessengerAddrIsNotL1SenderAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    address l2MessengerAddr = address(s_mockOptimismL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_optimismSequencerUptimeFeed.validateSenderTestWrapper(s_strangerAddr);
  }

  /// @notice Updates status when status has changed and incoming timestamp is the same as the latest
  function test_ValidateSender_UpdateStatusWhen_StatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockOptimismL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    s_optimismSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }
}
