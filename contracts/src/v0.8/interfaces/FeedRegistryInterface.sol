// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma abicoder v2;

import "./AggregatorV2V3Interface.sol";

interface FeedRegistryInterface {
  struct Phase {
    uint16 phaseId;
    uint80 startingAggregatorRoundId;
    uint80 endingAggregatorRoundId;
  }

  event FeedProposed(
    address indexed asset,
    address indexed denomination,
    address indexed proposedAggregator,
    address currentAggregator,
    address sender
  );
  event FeedConfirmed(
    address indexed asset,
    address indexed denomination,
    address indexed latestAggregator,
    address previousAggregator,
    uint16 nextPhaseId,
    address sender
  );

  // V3 AggregatorV3Interface

  function decimals(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint8
    );

  function description(
    address asset,
    address denomination
  )
    external
    view
    returns (
      string memory
    );

  function version(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint256
    );

  function latestRoundData(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );

  function getRoundData(
    address asset,
    address denomination,
    uint80 _roundId
  )
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );

  // V2 AggregatorInterface

  function latestAnswer(
    address asset,
    address denomination
  )
    external
    view
    returns (
      int256 answer
    );

  function latestTimestamp(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint256 timestamp
    );

  function latestRound(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint256 roundId
    );

  function getAnswer(
    address asset,
    address denomination,
    uint256 roundId
  )
    external
    view
    returns (
      int256 answer
    );

  function getTimestamp(
    address asset,
    address denomination,
    uint256 roundId
  )
    external
    view
    returns (
      uint256 timestamp
    );

  // Registry getters

  function getFeed(
    address asset,
    address denomination
  )
    external
    view
    returns (
      AggregatorV2V3Interface aggregator
    );

  function getPhaseFeed(
    address asset,
    address denomination,
    uint16 phaseId
  )
    external
    view
    returns (
      AggregatorV2V3Interface aggregator
    );

  function isFeedEnabled(
    address aggregator
  )
    external
    view
    returns (
      bool
    );

  function getPhase(
    address asset,
    address denomination,
    uint16 phaseId
  )
    external
    view
    returns (
      Phase memory phase
    );

  // Round helpers

  function getRoundFeed(
    address asset,
    address denomination,
    uint80 roundId
  )
    external
    view
    returns (
      AggregatorV2V3Interface aggregator
    );

  function getPhaseRange(
    address asset,
    address denomination,
    uint16 phaseId
  )
    external
    view
    returns (
      uint80 startingRoundId,
      uint80 endingRoundId
    );

  function getPreviousRoundId(
    address asset,
    address denomination,
    uint80 roundId
  ) external
    view
    returns (
      uint80 previousRoundId
    );

  function getNextRoundId(
    address asset,
    address denomination,
    uint80 roundId
  ) external
    view
    returns (
      uint80 nextRoundId
    );

  // Feed management

  function proposeFeed(
    address asset,
    address denomination,
    address aggregator
  ) external;

  function confirmFeed(
    address asset,
    address denomination,
    address aggregator
  ) external;

  // Proposed aggregator

  function getProposedFeed(
    address asset,
    address denomination
  )
    external
    view
    returns (
      AggregatorV2V3Interface proposedAggregator
    );

  function proposedGetRoundData(
    address asset,
    address denomination,
    uint80 roundId
  )
    external
    view
    returns (
      uint80 id,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );

  function proposedLatestRoundData(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint80 id,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );

  // Phases
  function getCurrentPhaseId(
    address asset,
    address denomination
  )
    external
    view
    returns (
      uint16 currentPhaseId
    );
}
