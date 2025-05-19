// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AddressAliasHelper} from "../../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import {ZKSyncSequencerUptimeFeed} from "../../../zksync/ZKSyncSequencerUptimeFeed.sol";
import {BaseSequencerUptimeFeed} from "../../../base/BaseSequencerUptimeFeed.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ZKSyncSequencerUptimeFeed_TestWrapper is ZKSyncSequencerUptimeFeed {
  constructor(address l1SenderAddress, bool initialStatus) ZKSyncSequencerUptimeFeed(l1SenderAddress, initialStatus) {}

  /// @notice Exposes the internal _validateSender function for testing
  function validateSenderTestWrapper(address l1Sender) external view {
    super._validateSender(l1Sender);
  }
}

contract ZKSyncSequencerUptimeFeed_Setup is L2EPTest {
  /// Helper Variables
  address internal l1SenderAddress = address(5);
  address internal s_aliasedL1SenderAddress = AddressAliasHelper.applyL1ToL2Alias(l1SenderAddress);

  /// L2EP contracts
  ZKSyncSequencerUptimeFeed_TestWrapper internal s_zksyncSequencerUptimeFeed;

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_zksyncSequencerUptimeFeed = new ZKSyncSequencerUptimeFeed_TestWrapper(l1SenderAddress, false);
  }
}

contract ZKSyncSequencerUptimeFeed_ValidateSender is ZKSyncSequencerUptimeFeed_Setup {
  /// @notice Reverts when pass is not valid
  function test_ValidateSender_RevertWhen_PassIsNotValid() public {
    // Sets msg.sender and tx.origin to an authorized address
    vm.startPrank(s_aliasedL1SenderAddress, s_aliasedL1SenderAddress);

    // Tries to update the status from an unauthorized account
    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_zksyncSequencerUptimeFeed.validateSenderTestWrapper(address(6));
  }

  /// @notice Passes when sender is valid
  function test_ValidateSender_SuccessWhen_SenderIsValid() public {
    // Sets msg.sender and tx.origin to an authorized address
    vm.startPrank(s_aliasedL1SenderAddress, s_aliasedL1SenderAddress);

    // Tries to update the status from an authorized account
    s_zksyncSequencerUptimeFeed.validateSenderTestWrapper(l1SenderAddress);
  }
}
