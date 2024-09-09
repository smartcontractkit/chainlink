// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockOptimismL1CrossDomainMessenger} from "../../../../tests/MockOptimismL1CrossDomainMessenger.sol";
import {MockOptimismL2CrossDomainMessenger} from "../../../../tests/MockOptimismL2CrossDomainMessenger.sol";
import {BaseSequencerUptimeFeed} from "../../../dev/shared/BaseSequencerUptimeFeed.sol";
import {OptimismSequencerUptimeFeed} from "../../../dev/optimism/OptimismSequencerUptimeFeed.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract OptimismSequencerUptimeFeedTestWrapper is OptimismSequencerUptimeFeed {
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) OptimismSequencerUptimeFeed(l1SenderAddress, l2CrossDomainMessengerAddr, initialStatus) {}

  /// @notice it exposes the internal _validateSender function for testing
  function validateSenderTestWrapper(address l1Sender) external view {
    super._validateSender(l1Sender);
  }
}

contract OptimismSequencerUptimeFeedTest is L2EPTest {
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);

  /// L2EP contracts
  MockOptimismL1CrossDomainMessenger internal s_mockOptimismL1CrossDomainMessenger;
  MockOptimismL2CrossDomainMessenger internal s_mockOptimismL2CrossDomainMessenger;
  OptimismSequencerUptimeFeedTestWrapper internal s_optimismSequencerUptimeFeed;

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_mockOptimismL1CrossDomainMessenger = new MockOptimismL1CrossDomainMessenger();
    s_mockOptimismL2CrossDomainMessenger = new MockOptimismL2CrossDomainMessenger();
    s_optimismSequencerUptimeFeed = new OptimismSequencerUptimeFeedTestWrapper(
      s_l1OwnerAddr,
      address(s_mockOptimismL2CrossDomainMessenger),
      false
    );

    // Sets mock sender in mock L2 messenger contract
    s_mockOptimismL2CrossDomainMessenger.setSender(s_l1OwnerAddr);
  }
}

contract OptimismSequencerUptimeFeed_Constructor is OptimismSequencerUptimeFeedTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialState() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    new OptimismSequencerUptimeFeedTestWrapper(s_l1OwnerAddr, address(s_mockOptimismL2CrossDomainMessenger), false);

    // Checks L1 sender
    address actualL1Addr = s_optimismSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_optimismSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract OptimismSequencerUptimeFeed_ValidateSender is OptimismSequencerUptimeFeedTest {
  /// @notice it should revert if called by an address that is not the L2 Cross Domain Messenger
  function test_RevertIfSenderIsNotL2CrossDomainMessengerAddr() public {
    address l2MessengerAddr = address(s_mockOptimismL2CrossDomainMessenger);
    // Sets msg.sender to a different address
    vm.startPrank(s_strangerAddr, l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_optimismSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }

  /// @notice it should revert if the L1 sender address is not the L1 Cross Domain Messenger Sender
  function test_RevertIfL1CrossDomainMessengerAddrIsNotL1SenderAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    address l2MessengerAddr = address(s_mockOptimismL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_optimismSequencerUptimeFeed.validateSenderTestWrapper(s_strangerAddr);
  }

  /// @notice it should update status when status has changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenStatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockOptimismL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    s_optimismSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }
}
