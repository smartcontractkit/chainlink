// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {MockOptimismL1CrossDomainMessenger} from "../../../../tests/MockOptimismL1CrossDomainMessenger.sol";
import {MockOptimismL2CrossDomainMessenger} from "../../../../tests/MockOptimismL2CrossDomainMessenger.sol";
import {OptimismSequencerUptimeFeed} from "../../../dev/optimism/OptimismSequencerUptimeFeed.sol";
import {OptimismValidator} from "../../../dev/optimism/OptimismValidator.sol";
import {SequencerUptimeFeed} from "../../../dev/SequencerUptimeFeed.sol";
import {GasLimitValidator} from "../../../dev/GasLimitValidator.sol";
import {GasLimitValidatorTest} from "../GasLimitValidator.t.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract OptimismValidatorTest is GasLimitValidatorTest {
  event SentMessage(address indexed target, address sender, bytes message, uint256 messageNonce, uint256 gasLimit);

  function newValidator() internal override returns (GasLimitValidator validator) {
    MockOptimismL1CrossDomainMessenger mockOptimismL1CrossDomainMessenger = new MockOptimismL1CrossDomainMessenger();
    MockOptimismL2CrossDomainMessenger mockOptimismL2CrossDomainMessenger = new MockOptimismL2CrossDomainMessenger();
    OptimismSequencerUptimeFeed optimismSequencerUptimeFeed = new OptimismSequencerUptimeFeed(
      address(mockOptimismL1CrossDomainMessenger),
      address(mockOptimismL2CrossDomainMessenger),
      true
    );
    return
      new OptimismValidator(
        address(mockOptimismL1CrossDomainMessenger),
        address(optimismSequencerUptimeFeed),
        INIT_GAS_LIMIT
      );
  }

  function emitExpectedSentMessageEvent(
    address validatorAddress,
    bool status,
    uint256 futureTimestampInSeconds
  ) internal override {
    emit SentMessage(
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      validatorAddress, // sender
      abi.encodeWithSelector(SequencerUptimeFeed.updateStatus.selector, status, futureTimestampInSeconds), // message
      0, // nonce
      INIT_GAS_LIMIT // gas limit
    );
  }

}

