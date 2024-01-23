// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {MockScrollL1CrossDomainMessenger} from "../../mocks/scroll/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/scroll/MockScrollL2CrossDomainMessenger.sol";
import {ScrollSequencerUptimeFeed} from "../../../dev/scroll/ScrollSequencerUptimeFeed.sol";
import {SequencerUptimeFeed} from "../../../dev/SequencerUptimeFeed.sol";
import {ScrollValidator} from "../../../dev/scroll/ScrollValidator.sol";
import {GasLimitValidator} from "../../../dev/GasLimitValidator.sol";
import {GasLimitValidatorTest} from "../GasLimitValidator.t.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ScrollValidatorTest is GasLimitValidatorTest {
  /// https://github.com/scroll-tech/scroll/blob/03089eaeee1193ff44c532c7038611ae123e7ef3/contracts/src/libraries/IScrollMessenger.sol#L22
  event SentMessage(
    address indexed sender,
    address indexed target,
    uint256 value,
    uint256 messageNonce,
    uint256 gasLimit,
    bytes message
  );

  function newValidator() internal override returns (GasLimitValidator validator) {
    MockScrollL1CrossDomainMessenger mockScrollL1CrossDomainMessenger = new MockScrollL1CrossDomainMessenger();
    MockScrollL2CrossDomainMessenger mockScrollL2CrossDomainMessenger = new MockScrollL2CrossDomainMessenger();
    ScrollSequencerUptimeFeed scrollSequencerUptimeFeed = new ScrollSequencerUptimeFeed(
      address(mockScrollL1CrossDomainMessenger),
      address(mockScrollL2CrossDomainMessenger),
      true
    );
    return
      new ScrollValidator(
        address(mockScrollL1CrossDomainMessenger),
        address(scrollSequencerUptimeFeed),
        INIT_GAS_LIMIT
      );
  }

  function emitExpectedSentMessageEvent(
    address validatorAddress,
    bool status,
    uint256 futureTimestampInSeconds
  ) internal override {
    emit SentMessage(
      validatorAddress, // sender
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      0, // value
      0, // nonce
      INIT_GAS_LIMIT, // gas limit
      abi.encodeWithSelector(SequencerUptimeFeed.updateStatus.selector, status, futureTimestampInSeconds) // message
    );
  }

}

