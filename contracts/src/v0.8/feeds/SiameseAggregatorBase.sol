// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

abstract contract SiameseAggregatorBase {
  struct Report {
    int192 juelsPerFeeCoin;
    uint32 observationsTimestamp;
    bytes observers; // ith element is the index of the ith observer
    int192[] observations; // ith element is the ith observation
  }

  struct Transmission {
    int192 answer;
    uint32 observationsTimestamp;
    uint32 recordedTimestamp; // renamed from transmissionTimestamp
    bool locked; // NB: New
  }

  mapping(uint32 /* aggregator round ID */ => Transmission) internal s_transmissions;

  address internal s_siameseAggregator;

  function recordSiameseReport(Report memory report) public virtual;

  function _duplicateReport(Report memory report, Transmission memory transmission) internal pure returns (bool) {
    // Reports don't have Round IDs so compare based on observation timestamp and answer.
    int192 reportAnswer = report.observations[report.observations.length / 2];

    return
      report.observationsTimestamp == transmission.observationsTimestamp &&
      transmission.answer == reportAnswer &&
      transmission.locked;
  }
}
