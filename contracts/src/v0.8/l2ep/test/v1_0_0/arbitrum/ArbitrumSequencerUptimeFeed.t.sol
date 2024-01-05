// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {SimpleWriteAccessController} from "../../../../shared/access/SimpleWriteAccessController.sol";
import {ArbitrumSequencerUptimeFeed} from "../../../dev/arbitrum/ArbitrumSequencerUptimeFeed.sol";
import {FeedConsumer} from "../../../../../v0.8/tests/FeedConsumer.sol";
import {MockAggregatorV2V3} from "../../mocks/MockAggregatorV2V3.sol";
import {Flags} from "../../../dev/Flags.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/arbitrum/ArbitrumSequencerUptimeFeed.t.sol
//
contract ArbitrumSequencerUptimeFeedTest is L2EPTest {
  /// Constants
  uint256 internal constant GAS_USED_DEVIATION = 100;

  /// Helper variables
  address internal s_mockL1OwnerAddr = vm.addr(0x1);
  address internal s_strangerAddr = vm.addr(0x2);
  address internal s_deployerAddr = vm.addr(0x3);
  address internal s_l2MessengerAddr =
    address(uint160(s_mockL1OwnerAddr) + uint160(0x1111000000000000000000000000000000001111));

  /// L2EP contracts
  ArbitrumSequencerUptimeFeed internal s_arbitrumSequencerUptimeFeed;
  SimpleWriteAccessController internal s_accessController;
  MockAggregatorV2V3 internal s_l1GasFeed;
  Flags internal s_flags;

  /// Events
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);
  event RoundUpdated(int256 status, uint64 updatedAt);
  event Initialized();

  /// Setup
  function setUp() public {
    vm.startPrank(s_deployerAddr, s_deployerAddr);

    s_accessController = new SimpleWriteAccessController();
    s_flags = new Flags(address(s_accessController), address(s_accessController));
    s_arbitrumSequencerUptimeFeed = new ArbitrumSequencerUptimeFeed(address(s_flags), s_mockL1OwnerAddr);

    s_accessController.addAccess(address(s_arbitrumSequencerUptimeFeed));
    s_accessController.addAccess(address(s_flags));
    s_accessController.addAccess(s_deployerAddr);
    s_flags.addAccess(address(s_arbitrumSequencerUptimeFeed));

    vm.expectEmit(false, false, false, true);
    emit Initialized();
    s_arbitrumSequencerUptimeFeed.initialize();

    vm.stopPrank();
  }
}

contract Constants is ArbitrumSequencerUptimeFeedTest {
  /// @notice it should have the correct value for FLAG_L2_SEQ_OFFLINE'
  function test_InitialState() public {
    assertEq(s_arbitrumSequencerUptimeFeed.FLAG_L2_SEQ_OFFLINE(), 0xa438451D6458044c3c8CD2f6f31c91ac882A6d91);
  }
}

contract UpdateStatus is ArbitrumSequencerUptimeFeedTest {
  /// @notice it should revert if called by an address that is not the L2 Cross Domain Messenger
  function test_RevertIfNotL2CrossDomainMessengerAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Tries to update the status from an unauthorized account
    vm.expectRevert(abi.encodeWithSelector(ArbitrumSequencerUptimeFeed.InvalidSender.selector));
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(1));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should update status when status has changed and incoming timestamp is newer than the latest
  function test_UpdateStatusWhenStatusChangeAndTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Submits a status update
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp();
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_arbitrumSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, newer timestamp should update
    timestamp = timestamp + 200;
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(0, 3, timestamp);
    s_arbitrumSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_arbitrumSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should update status when status has changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenStatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Fetches the latest timestamp
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_arbitrumSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, same timestamp should update
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(0, 3, timestamp);
    s_arbitrumSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_arbitrumSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should ignore out-of-order updates
  function test_IgnoreOutOfOrderUpdates() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Submits a status update
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 10000;
    vm.expectEmit(false, false, false, true);
    emit AnswerUpdated(1, 2, timestamp);
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_arbitrumSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Update with different status, but stale timestamp, should be ignored
    timestamp = timestamp - 1000;
    vm.expectEmit(false, false, false, false);
    emit UpdateIgnored(true, 0, true, 0); // arguments are dummy values
    // TODO: how can we check that an AnswerUpdated event was NOT emitted
    s_arbitrumSequencerUptimeFeed.updateStatus(false, uint64(timestamp));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AggregatorV3Interface is ArbitrumSequencerUptimeFeedTest {
  /// @notice it should return valid answer from getRoundData and latestRoundData
  function test_AggregatorV3Interface() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables
    uint80 roundId;
    int256 answer;
    uint256 startedAt;
    uint256 updatedAt;
    uint80 answeredInRound;

    // Checks initial state
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_arbitrumSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Submits status update with different status and newer timestamp, should update
    uint256 timestamp = startedAt + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_arbitrumSequencerUptimeFeed.getRoundData(2);
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
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_arbitrumSequencerUptimeFeed.getRoundData(1);
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Assert latestRoundData corresponds to latest round id
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_arbitrumSequencerUptimeFeed.latestRoundData();
    assertEq(roundId2, roundId);
    assertEq(answer2, answer);
    assertEq(startedAt2, startedAt);
    assertEq(updatedAt2, updatedAt);
    assertEq(answeredInRound2, answeredInRound);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert from #getRoundData when round does not yet exist (future roundId)
  function test_Return0WhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_mockL1OwnerAddr, s_mockL1OwnerAddr);

    // Gets data from a round that has not happened yet
    (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    ) = s_arbitrumSequencerUptimeFeed.getRoundData(2);

    // Validates round data
    assertEq(roundId, 2);
    assertEq(answer, 0);
    assertEq(startedAt, 0);
    assertEq(updatedAt, 0);
    assertEq(answeredInRound, 2);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract ProtectReadsOnAggregatorV2V3InterfaceFunctions is ArbitrumSequencerUptimeFeedTest {
  /// @notice it should disallow reads on AggregatorV2V3Interface functions when consuming contract is not whitelisted
  function test_AggregatorV2V3InterfaceDisallowReadsIfConsumingContractIsNotWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_arbitrumSequencerUptimeFeed));

    // Sanity - consumer is not whitelisted
    assertEq(s_arbitrumSequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_arbitrumSequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), false);

    // Asserts reads are not possible from consuming contract
    vm.expectRevert("No access");
    feedConsumer.latestAnswer();
    vm.expectRevert("No access");
    feedConsumer.latestRoundData();
  }

  /// @notice it should allow reads on AggregatorV2V3Interface functions when consuming contract is whitelisted
  function test_AggregatorV2V3InterfaceAllowReadsIfConsumingContractIsWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_arbitrumSequencerUptimeFeed));

    // Whitelist consumer
    vm.startPrank(s_deployerAddr, s_deployerAddr);
    s_arbitrumSequencerUptimeFeed.addAccess(address(feedConsumer));
    vm.stopPrank();

    // Sanity - consumer is whitelisted
    assertEq(s_arbitrumSequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_arbitrumSequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), true);

    // Asserts reads are possible from consuming contract
    (uint80 roundId, int256 answer, , , ) = feedConsumer.latestRoundData();
    assertEq(feedConsumer.latestAnswer(), 0);
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract GasCosts is ArbitrumSequencerUptimeFeedTest {
  /// @notice it should consume a known amount of gas for updates @skip-coverage
  function test_GasCosts() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Assert initial conditions
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp();
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 0);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed;
    uint256 gasStart;
    uint256 gasFinal;

    // measures gas used for no update
    expectedGasUsed = 5507; // TODO: used to be 28300
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.updateStatus(false, uint64(timestamp + 1000));
    gasFinal = gasleft();
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 0);
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // measures gas used for update
    expectedGasUsed = 68198; // TODO: Used to be 93015
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp + 1000));
    gasFinal = gasleft();
    assertEq(s_arbitrumSequencerUptimeFeed.latestAnswer(), 1);
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AggregatorInterfaceGasCosts is ArbitrumSequencerUptimeFeedTest {
  /// @notice it should consume a known amount of gas for getRoundData(uint80) @skip-coverage
  function test_GasUsageForGetRoundData() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 4658; // TODO: used to be 31157
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.getRoundData(1);
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestRoundData() @skip-coverage
  function test_GasUsageForLatestRoundData() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 2154; // TODO: used to be 28523
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.latestRoundData();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestAnswer() @skip-coverage
  function test_GasUsageForLatestAnswer() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 1722; // TODO: used to be 28329
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.latestAnswer();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestTimestamp() @skip-coverage
  function test_GasUsageForLatestTimestamp() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 1652; // TODO: used to be 28229
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.latestTimestamp();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for latestRound() @skip-coverage
  function test_GasUsageForLatestRound() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 1632; // TODO: used to be 28245
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.latestRound();
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for getAnswer() @skip-coverage
  function test_GasUsageForGetAnswer() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 4059; // TODO: used to be 30799
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.getAnswer(1);
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should consume a known amount of gas for getTimestamp() @skip-coverage
  function test_GasUsageForGetTimestamp() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l2MessengerAddr, s_l2MessengerAddr);

    // Defines helper variables for measuring gas usage
    uint256 expectedGasUsed = 4024; // TODO: used to be 30753
    uint256 gasStart;
    uint256 gasFinal;

    // Initializes a round
    uint256 timestamp = s_arbitrumSequencerUptimeFeed.latestTimestamp() + 1000;
    s_arbitrumSequencerUptimeFeed.updateStatus(true, uint64(timestamp));

    // Measures gas usage
    gasStart = gasleft();
    s_arbitrumSequencerUptimeFeed.getTimestamp(1);
    gasFinal = gasleft();

    // Checks that gas usage is within expected range
    assertGasUsageIsCloseTo(expectedGasUsed, gasStart, gasFinal, GAS_USED_DEVIATION);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}
