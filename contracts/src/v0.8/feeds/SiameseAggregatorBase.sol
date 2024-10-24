// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {AggregatorV2V3Interface} from "../shared/interfaces/AggregatorV2V3Interface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {OCR2Abstract} from "../shared/ocr2/OCR2Abstract.sol";
import {LinkTokenInterface} from "../shared/interfaces/LinkTokenInterface.sol";
import {AccessControllerInterface} from "../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../shared/interfaces/AggregatorValidatorInterface.sol";

abstract contract SiameseAggregatorBase {
  struct Report {
    uint32 observationsTimestamp;
    bytes observers; // ith element is the index of the ith observer
    int192[] observations; // ith element is the ith observation
    int192 juelsPerFeeCoin;
  }

  struct Transmission {
    int192 answer;
    uint32 observationsTimestamp;
    uint32 recordedTimestamp; // renamed from transmissionTimestamp
    bool locked;  // NB: New
  }

  mapping(uint32 /* aggregator round ID */ => Transmission) internal s_transmissions;

  address private s_siameseAggregator;

  function recordSiameseReport(Report memory report) public virtual;

  function duplicateReport(Report memory report, Transmission memory transmission) internal returns (bool) {
    // Reports don't have Round IDs so compare based on observation timestamp and answer.
    int192 reportAnswer = report.observations[report.observations.length/2];

    return report.observationsTimestamp == transmission.observationsTimestamp &&
      transmission.answer == reportAnswer &&
      transmission.locked;
  }
}

