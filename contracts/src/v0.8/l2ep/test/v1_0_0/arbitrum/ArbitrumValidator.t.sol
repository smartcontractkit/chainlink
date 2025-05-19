// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";

import {SimpleWriteAccessController} from "../../../../shared/access/SimpleWriteAccessController.sol";
import {ArbitrumSequencerUptimeFeed} from "../../../arbitrum/ArbitrumSequencerUptimeFeed.sol";
import {ArbitrumValidator} from "../../../arbitrum/ArbitrumValidator.sol";
import {BaseValidator} from "../../../base/BaseValidator.sol";
import {MockArbitrumInbox} from "../../mocks/MockArbitrumInbox.sol";
import {MockAggregatorV2V3} from "../../mocks/MockAggregatorV2V3.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ArbitrumValidatorTest is L2EPTest {
  /// Helper constants
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = 0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b;
  uint256 internal constant GAS_PRICE_BID = 1000000;
  uint256 internal constant BASE_FEE = 14000000000;
  uint256 internal constant MAX_GAS = 1000000;

  /// L2EP contracts
  AccessControllerInterface internal s_accessController;
  MockArbitrumInbox internal s_mockArbitrumInbox;
  ArbitrumValidator internal s_arbitrumValidator;
  MockAggregatorV2V3 internal s_l1GasFeed;

  /// Events
  event RetryableTicketNoRefundAliasRewriteCreated(
    address destAddr,
    uint256 arbTxCallValue,
    uint256 maxSubmissionCost,
    address submissionRefundAddress,
    address valueRefundAddress,
    uint256 maxGas,
    uint256 gasPriceBid,
    bytes data
  );

  /// Setup
  function setUp() public {
    s_accessController = new SimpleWriteAccessController();
    s_mockArbitrumInbox = new MockArbitrumInbox();
    s_l1GasFeed = new MockAggregatorV2V3();
    s_arbitrumValidator = new ArbitrumValidator(
      address(s_mockArbitrumInbox),
      L2_SEQ_STATUS_RECORDER_ADDRESS,
      address(s_accessController),
      MAX_GAS,
      GAS_PRICE_BID,
      BASE_FEE,
      address(s_l1GasFeed),
      ArbitrumValidator.PaymentStrategy.L1
    );
  }
}

contract ArbitrumValidator_Validate is ArbitrumValidatorTest {
  /// @notice it post sequencer offline
  function test_PostSequencerOffline() public {
    // Gives access to the s_eoaValidator
    s_arbitrumValidator.addAccess(s_eoaValidator);

    // Gets the ArbitrumValidator L2 address
    address arbitrumValidatorL2Addr = toArbitrumL2AliasAddress(address(s_arbitrumValidator));

    // Sets block.timestamp to a later date, funds the ArbitrumValidator contract, and sets msg.sender and tx.origin
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.warp(futureTimestampInSeconds);
    vm.deal(address(s_arbitrumValidator), 1 ether);
    vm.startPrank(s_eoaValidator);

    // Sets up the expected event data
    vm.expectEmit();
    emit RetryableTicketNoRefundAliasRewriteCreated(
      L2_SEQ_STATUS_RECORDER_ADDRESS, // destAddr
      0, // arbTxCallValue
      25312000000000, // maxSubmissionCost
      arbitrumValidatorL2Addr, // submissionRefundAddress
      arbitrumValidatorL2Addr, // valueRefundAddress
      MAX_GAS, // maxGas
      GAS_PRICE_BID, // gasPriceBid
      abi.encodeWithSelector(ArbitrumSequencerUptimeFeed.updateStatus.selector, true, futureTimestampInSeconds) // data
    );

    uint256 previousRoundId = 0;
    int256 previousAnswer = 0;
    uint256 currentRoundId = 1;
    int256 currentAnswer = 1;

    vm.expectEmit(address(s_arbitrumValidator));
    emit BaseValidator.ValidatedStatus(previousRoundId, previousAnswer, currentRoundId, currentAnswer);
    // Runs the function (which produces the event to test)
    s_arbitrumValidator.validate(previousRoundId, previousAnswer, currentRoundId, currentAnswer);
  }
}
