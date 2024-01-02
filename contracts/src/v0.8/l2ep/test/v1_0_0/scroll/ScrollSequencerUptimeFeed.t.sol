// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {MockScrollL1CrossDomainMessenger} from "../../mocks/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/MockScrollL2CrossDomainMessenger.sol";
import {ScrollSequencerUptimeFeed} from "../../../dev/scroll/ScrollSequencerUptimeFeed.sol";
import {FeedConsumer} from "../../../../../v0.8/tests/FeedConsumer.sol";
import {L2EPTest} from "../L2EPTest.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/scroll/ScrollSequencerUptimeFeed.t.sol
//
contract ScrollSequencerUptimeFeedTest is L2EPTest {
  /// Constants
  uint256 internal constant GAS_USED_DEVIATION = 100;

  /// Helper variables
  address internal s_mockL1OwnerAddr = vm.addr(0x1);
  address internal s_strangerAddr = vm.addr(0x2);

  /// L2EP contracts
  MockScrollL1CrossDomainMessenger internal s_mockScrollL1CrossDomainMessenger;
  MockScrollL2CrossDomainMessenger internal s_mockScrollL2CrossDomainMessenger;
  ScrollSequencerUptimeFeed internal s_scrollSequencerUptimeFeed;

  /// Events
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);
  event RoundUpdated(int256 status, uint64 updatedAt);

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_mockScrollL1CrossDomainMessenger = new MockScrollL1CrossDomainMessenger();
    s_mockScrollL2CrossDomainMessenger = new MockScrollL2CrossDomainMessenger();
    s_scrollSequencerUptimeFeed = new ScrollSequencerUptimeFeed(
      s_mockL1OwnerAddr,
      address(s_mockScrollL2CrossDomainMessenger),
      false
    );

    // Sets mock sender in mock L2 messenger contract
    s_mockScrollL2CrossDomainMessenger.setSender(s_mockL1OwnerAddr);
  }
}

contract Constructor is ScrollSequencerUptimeFeedTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialState() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_mockL1OwnerAddr, s_mockL1OwnerAddr);

    // Checks L1 sender
    address actualL1Addr = s_scrollSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_mockL1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_scrollSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract UpdateStatus is ScrollSequencerUptimeFeedTest {
  /// @notice it should revert if called by an address that is not the L2 Cross Domain Messenger
  function test_RevertIfNotL2CrossDomainMessengerAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Tries to update the status from an unauthorized account
    vm.expectRevert(abi.encodeWithSelector(ScrollSequencerUptimeFeed.InvalidSender.selector));
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(1));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert if called by an address that is not the L2 Cross Domain Messenger and is not the L1 sender
  function test_RevertIfNotL2CrossDomainMessengerAddrAndNotL1SenderAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sets mock sender in mock L2 messenger contract
    s_mockScrollL2CrossDomainMessenger.setSender(s_strangerAddr);

    // Tries to update the status from an unauthorized account
    vm.expectRevert(abi.encodeWithSelector(ScrollSequencerUptimeFeed.InvalidSender.selector));
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(1));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should update status when status has not changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenNoChange() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Fetches the latest timestamp
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Stores the current round data before updating it
    (
      uint80 roundIdBeforeUpdate,
      int256 answerBeforeUpdate,
      uint256 startedAtBeforeUpdate,
      ,
      uint80 answeredInRoundBeforeUpdate
    ) = s_scrollSequencerUptimeFeed.latestRoundData();

    // Submit another status update with the same status
    vm.expectEmit(false, false, false, true);
    emit RoundUpdated(1, uint64(block.timestamp));
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp + 200));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Stores the current round data after updating it
    (
      uint80 roundIdAfterUpdate,
      int256 answerAfterUpdate,
      uint256 startedAtAfterUpdate,
      uint256 updatedAtAfterUpdate,
      uint80 answeredInRoundAfterUpdate
    ) = s_scrollSequencerUptimeFeed.latestRoundData();

    // Verifies the latest round data has been properly updated
    assertEq(roundIdAfterUpdate, roundIdBeforeUpdate);
    assertEq(answerAfterUpdate, answerBeforeUpdate);
    assertEq(startedAtAfterUpdate, startedAtBeforeUpdate);
    assertEq(answeredInRoundAfterUpdate, answeredInRoundBeforeUpdate);
    assertEq(updatedAtAfterUpdate, block.timestamp);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should update status when status has changed and incoming timestamp is newer than the latest
  function test_UpdateStatusWhenStatusChangeAndTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Submits a status update
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp();
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, newer timestamp should update
    timestamp = timestamp + 200;
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(0, 3, timestamp);
    s_scrollSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should update status when status has changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenStatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Fetches the latest timestamp
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, same timestamp should update
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(0, 3, timestamp);
    s_scrollSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should ignore out-of-order updates
  function test_IgnoreOutOfOrderUpdates() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Submits a status update
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 10000;
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_scrollSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Update with different status, but stale timestamp, should be ignored
    timestamp = timestamp - 1000;
    vm.expectEmit(false, false, false, false);
    emit UpdateIgnored(true, 0, true, 0); // arguments are dummy values
    // TODO: how can we check that an AnswerUpdated event was NOT emitted
    s_scrollSequencerUptimeFeed.updateStatus(false, uint64(timestamp));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AggregatorV3Interface is ScrollSequencerUptimeFeedTest {
  /// @notice it should return valid answer from getRoundData and latestRoundData
  function test_AggregatorV3Interface() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables
    uint80 roundId;
    int256 answer;
    uint256 startedAt;
    uint256 updatedAt;
    uint80 answeredInRound;

    // Checks initial state
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_scrollSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Submits status update with different status and newer timestamp, should update
    uint256 timestamp = startedAt + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_scrollSequencerUptimeFeed.getRoundData(2);
    assertEq(roundId, 2);
    assertEq(answer, 1);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, timestamp);
    assertLe(updatedAt, startedAt);

    // Saves round 2 data
    uint80 roundId2 = roundId;
    int256 answer2 = answer;
    uint256 startedAt2 = startedAt;
    uint256 updatedAt2 = updatedAt;
    uint80 answeredInRound2 = answeredInRound;

    // Checks that last round is still returning the correct data
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_scrollSequencerUptimeFeed.getRoundData(1);
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Assert latestRoundData corresponds to latest round id
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_scrollSequencerUptimeFeed.latestRoundData();
    assertEq(roundId2, roundId);
    assertEq(answer2, answer);
    assertEq(startedAt2, startedAt);
    assertEq(updatedAt2, updatedAt);
    assertEq(answeredInRound2, answeredInRound);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert from #getRoundData when round does not yet exist (future roundId)
  function test_RevertGetRoundDataWhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_mockL1OwnerAddr, s_mockL1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(abi.encodeWithSelector(ScrollSequencerUptimeFeed.NoDataPresent.selector));
    s_scrollSequencerUptimeFeed.getRoundData(2);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert from #getAnswer when round does not yet exist (future roundId)
  function test_RevertGetAnswerWhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_mockL1OwnerAddr, s_mockL1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(abi.encodeWithSelector(ScrollSequencerUptimeFeed.NoDataPresent.selector));
    s_scrollSequencerUptimeFeed.getAnswer(2);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert from #getTimestamp when round does not yet exist (future roundId)
  function test_RevertGetTimestampWhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_mockL1OwnerAddr, s_mockL1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(abi.encodeWithSelector(ScrollSequencerUptimeFeed.NoDataPresent.selector));
    s_scrollSequencerUptimeFeed.getTimestamp(2);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract ProtectReadsOnAggregatorV2V3InterfaceFunctions is ScrollSequencerUptimeFeedTest {
  /// @notice it should disallow reads on AggregatorV2V3Interface functions when consuming contract is not whitelisted
  function test_AggregatorV2V3InterfaceDisallowReadsIfConsumingContractIsNotWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_scrollSequencerUptimeFeed));

    // Sanity - consumer is not whitelisted
    assertEq(s_scrollSequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_scrollSequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), false);

    // Asserts reads are not possible from consuming contract
    vm.expectRevert("No access");
    feedConsumer.latestAnswer();
    vm.expectRevert("No access");
    feedConsumer.latestRoundData();
  }

  /// @notice it should allow reads on AggregatorV2V3Interface functions when consuming contract is whitelisted
  function test_AggregatorV2V3InterfaceAllowReadsIfConsumingContractIsWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_scrollSequencerUptimeFeed));

    // Whitelist consumer
    s_scrollSequencerUptimeFeed.addAccess(address(feedConsumer));

    // Sanity - consumer is whitelisted
    assertEq(s_scrollSequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_scrollSequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), true);

    // Asserts reads are possible from consuming contract
    (uint80 roundId, int256 answer, , , ) = feedConsumer.latestRoundData();
    assertEq(feedConsumer.latestAnswer(), 0);
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract GasCosts is ScrollSequencerUptimeFeedTest {
  /// @notice it should consume a known amount of gas for updates @skip-coverage
  function test_GasCosts() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Assert initial conditions
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp();
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 0);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed;
    uint256 gasStart;
    uint256 gasFinal;

    // measures gas used for no update
    expectedGasUsed = 10197; // TODO: used to be 38594
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.updateStatus(false, uint64(timestamp + 1000));
    gasFinal = gasleft();
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 0);
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // measures gas used for update
    expectedGasUsed = 31644; // TODO: Used to be 58458
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp + 1000));
    gasFinal = gasleft();
    assertEq(s_scrollSequencerUptimeFeed.latestAnswer(), 1);
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AggregatorInterfaceGasCosts is ScrollSequencerUptimeFeedTest {
  /// @notice it should consume a known amount of gas for getRoundData(uint80) @skip-coverage
  function test_GasUsageForGetRoundData() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 4504; // TODO: used to be 30952
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.getRoundData(1);
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestRoundData() @skip-coverage
  function test_GasUsageForLatestRoundData() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 2154; // TODO: used to be 28523
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.latestRoundData();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestAnswer() @skip-coverage
  function test_GasUsageForLatestAnswer() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 1566; // TODO: used to be 28229
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.latestAnswer();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestTimestamp() @skip-coverage
  function test_GasUsageForLatestTimestamp() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 1459; // TODO: used to be 28129
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.latestTimestamp();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestRound() @skip-coverage
  function test_GasUsageForLatestRound() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 1470; // TODO: used to be 28145
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.latestRound();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for getAnswer() @skip-coverage
  function test_GasUsageForGetAnswer() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 3929; // TODO: used to be 30682
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.getAnswer(1);
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for getTimestamp() @skip-coverage
  function test_GasUsageForGetTimestamp() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 3817; // TODO: used to be 30570
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_scrollSequencerUptimeFeed.latestTimestamp() + 1000;
    s_scrollSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_scrollSequencerUptimeFeed.getTimestamp(1);
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}
