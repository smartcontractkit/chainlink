// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IEVM2AnyOnRamp} from "../../interfaces/IEVM2AnyOnRamp.sol";

import {NonceManager} from "../../NonceManager.sol";
import {Client} from "../../libraries/Client.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";
import {OnRampHelper} from "../helpers/OnRampHelper.sol";
import {OnRampSetup} from "../onRamp/OnRamp/OnRampSetup.t.sol";

contract NonceManager_getOutboundNonce is OnRampSetup {
  uint256 internal constant FEE_AMOUNT = 1234567890;
  OnRampHelper internal s_prevOnRamp;

  function setUp() public virtual override {
    super.setUp();

    (s_prevOnRamp,) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry)
    );

    // Since the previous onRamp is not a 1.5 ramp it doesn't have the getSenderNonce function. We mock it to return 0
    vm.mockCall(address(s_prevOnRamp), abi.encodeWithSelector(IEVM2AnyOnRamp.getSenderNonce.selector), abi.encode(0));

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(address(s_prevOnRamp), address(0)),
      overrideExistingRamps: false
    });
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    (s_onRamp, s_metadataHash) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry)
    );

    vm.startPrank(address(s_sourceRouter));
  }

  function test_getOutboundNonce_Upgrade() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, FEE_AMOUNT, OWNER));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);
  }

  function test_getOutboundNonce_UpgradeSenderNoncesReadsPreviousRamp() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 startNonce = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);
    uint64 prevRampNextOutboundNonce = IEVM2AnyOnRamp(address(s_prevOnRamp)).getSenderNonce(OWNER);

    assertEq(startNonce, prevRampNextOutboundNonce);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      assertEq(startNonce + i, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));
    }
  }

  function test_getOutboundNonce_UpgradeNonceStartsAtV1Nonce() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 startNonce = s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);

    // send 1 message from previous onRamp
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 1, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // new onRamp nonce should start from 2, while sequence number start from 1
    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, startNonce + 2, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 2, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // after another send, nonce should be 3, and sequence number be 2
    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 2, _messageToEvent(message, 2, startNonce + 3, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 3, s_outboundNonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));
  }

  function test_getOutboundNonce_UpgradeNonceNewSenderStartsAtZero() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    // send 1 message from previous onRamp from OWNER
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    address newSender = address(1234567);
    // new onRamp nonce should start from 1 for new sender
    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, _messageToEvent(message, 1, 1, FEE_AMOUNT, newSender));

    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, newSender);
  }
}
