// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {AggregatorV2V3Interface} from "../shared/interfaces/AggregatorV2V3Interface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {OCR2Abstract} from "../shared/ocr2/OCR2Abstract.sol";
import {LinkTokenInterface} from "../shared/interfaces/LinkTokenInterface.sol";
import {AccessControllerInterface} from "../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorValidatorInterface} from "../shared/interfaces/AggregatorValidatorInterface.sol";
import {PrimaryAggregator} from "./PrimaryAggregator.sol";

contract SecondaryAggregator is PrimaryAggregator {
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_
  ) PrimaryAggregator(
    link,
    minAnswer_,
    maxAnswer_,
    billingAccessController,
    requesterAccessController,
    decimals_,
    description_
  ) {}
}

