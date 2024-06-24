// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {NonceManager} from "../NonceManager.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {EVM2EVMMultiOnRamp} from "../onRamp/EVM2EVMMultiOnRamp.sol";
import {EVM2EVMOnRamp} from "../onRamp/EVM2EVMOnRamp.sol";
import {EVM2EVMMultiOnRampHelper} from "./helpers/EVM2EVMMultiOnRampHelper.sol";
import {EVM2EVMOnRampHelper} from "./helpers/EVM2EVMOnRampHelper.sol";
import {EVM2EVMMultiOnRampSetup} from "./onRamp/EVM2EVMMultiOnRampSetup.t.sol";

contract NonceManagerTest_getIncrementedOutboundNonce is EVM2EVMMultiOnRampSetup {
  function test_getIncrementedOutboundNonce_Success() public {
    vm.startPrank(address(s_onRamp));
    address sender = address(this);

    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 0);

    uint64 outboundNonce = s_nonceManager.getIncrementedOutboundNonce(DEST_CHAIN_SELECTOR, sender);
    assertEq(outboundNonce, 1);
  }
}

contract NonceManager_applyPreviousRampsUpdates is EVM2EVMMultiOnRampSetup {
  function test_SingleRampUpdate() public {
    address prevOnRamp = vm.addr(1);
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] = NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp));

    vm.expectEmit();
    emit NonceManager.PreviousOnRampUpdated(DEST_CHAIN_SELECTOR, prevOnRamp);

    s_nonceManager.applyPreviousRampsUpdates(previousRamps);

    _assertPreviousRampsEqual(s_nonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR), previousRamps[0].prevRamps);
  }

  function test_MultipleRampsUpdates() public {
    address prevOnRamp1 = vm.addr(1);
    address prevOnRamp2 = vm.addr(2);
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](2);
    previousRamps[0] = NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(prevOnRamp1));
    previousRamps[1] = NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR + 1, NonceManager.PreviousRamps(prevOnRamp2));

    vm.expectEmit();
    emit NonceManager.PreviousOnRampUpdated(DEST_CHAIN_SELECTOR, prevOnRamp1);
    vm.expectEmit();
    emit NonceManager.PreviousOnRampUpdated(DEST_CHAIN_SELECTOR + 1, prevOnRamp2);

    s_nonceManager.applyPreviousRampsUpdates(previousRamps);

    _assertPreviousRampsEqual(s_nonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR), previousRamps[0].prevRamps);
    _assertPreviousRampsEqual(s_nonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR + 1), previousRamps[1].prevRamps);
  }

  function test_ZeroInput() public {
    vm.recordLogs();
    s_nonceManager.applyPreviousRampsUpdates(new NonceManager.PreviousRampsArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_PreviousRampAlreadySetOnRamp_Revert() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(vm.addr(1))));

    s_nonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(vm.addr(2))));

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_nonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function _assertPreviousRampsEqual(
    NonceManager.PreviousRamps memory a,
    NonceManager.PreviousRamps memory b
  ) internal pure {
    assertEq(a.prevOnRamp, b.prevOnRamp);
  }
}

contract NonceManager_onRampUpgrade is EVM2EVMMultiOnRampSetup {
  uint256 internal constant FEE_AMOUNT = 1234567890;
  EVM2EVMOnRampHelper internal s_prevOnRamp;

  function setUp() public virtual override {
    EVM2EVMMultiOnRampSetup.setUp();

    EVM2EVMOnRamp.FeeTokenConfigArgs[] memory feeTokenConfigArgs = new EVM2EVMOnRamp.FeeTokenConfigArgs[](1);
    feeTokenConfigArgs[0] = EVM2EVMOnRamp.FeeTokenConfigArgs({
      token: s_sourceFeeToken,
      networkFeeUSDCents: 1_00, // 1 USD
      gasMultiplierWeiPerEth: 1e18, // 1x
      premiumMultiplierWeiPerEth: 5e17, // 0.5x
      enabled: true
    });

    EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfig =
      new EVM2EVMOnRamp.TokenTransferFeeConfigArgs[](1);

    tokenTransferFeeConfig[0] = EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
      token: s_sourceFeeToken,
      minFeeUSDCents: 1_00, // 1 USD
      maxFeeUSDCents: 1000_00, // 1,000 USD
      deciBps: 2_5, // 2.5 bps, or 0.025%
      destGasOverhead: 40_000,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      aggregateRateLimitEnabled: true
    });

    s_prevOnRamp = new EVM2EVMOnRampHelper(
      EVM2EVMOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        destChainSelector: DEST_CHAIN_SELECTOR,
        defaultTxGasLimit: GAS_LIMIT,
        maxNopFeesJuels: MAX_NOP_FEES_JUELS,
        prevOnRamp: address(0),
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      EVM2EVMOnRamp.DynamicConfig({
        router: address(s_sourceRouter),
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        destGasOverhead: DEST_GAS_OVERHEAD,
        destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
        destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
        destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
        destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
        priceRegistry: address(s_priceRegistry),
        maxDataBytes: MAX_DATA_SIZE,
        maxPerMsgGasLimit: MAX_GAS_LIMIT,
        defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
        defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        defaultTokenDestBytesOverhead: DEFAULT_TOKEN_BYTES_OVERHEAD,
        enforceOutOfOrder: false
      }),
      RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e15}),
      feeTokenConfigArgs,
      tokenTransferFeeConfig,
      new EVM2EVMOnRamp.NopAndWeight[](0)
    );

    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] =
      NonceManager.PreviousRampsArgs(DEST_CHAIN_SELECTOR, NonceManager.PreviousRamps(address(s_prevOnRamp)));
    s_nonceManager.applyPreviousRampsUpdates(previousRamps);

    EVM2EVMMultiOnRamp.DestChainConfigArgs[] memory destChainConfigArgs = _generateDestChainConfigArgs();
    destChainConfigArgs[0].prevOnRamp = address(s_prevOnRamp);

    (s_onRamp, s_metadataHash) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, address(s_sourceRouter), address(s_nonceManager), address(s_tokenAdminRegistry)
    );

    vm.startPrank(address(s_sourceRouter));
  }

  function test_Upgrade_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, FEE_AMOUNT, OWNER));
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);
  }

  function test_UpgradeSenderNoncesReadsPreviousRamp_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 startNonce = s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);

    for (uint64 i = 1; i < 4; ++i) {
      s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, 0, OWNER);

      assertEq(startNonce + i, s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));
    }
  }

  function test_UpgradeNonceStartsAtV1Nonce_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint64 startNonce = s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER);

    // send 1 message from previous onramp
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 1, s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // new onramp nonce should start from 2, while sequence number start from 1
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, startNonce + 2, FEE_AMOUNT, OWNER)
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 2, s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));

    // after another send, nonce should be 3, and sequence number be 2
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(message, 2, startNonce + 3, FEE_AMOUNT, OWNER)
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    assertEq(startNonce + 3, s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, OWNER));
  }

  function test_UpgradeNonceNewSenderStartsAtZero_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    // send 1 message from previous onramp from OWNER
    s_prevOnRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, OWNER);

    address newSender = address(1234567);
    // new onramp nonce should start from 1 for new sender
    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(
      DEST_CHAIN_SELECTOR, _messageToEvent(message, 1, 1, FEE_AMOUNT, newSender)
    );
    s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, FEE_AMOUNT, newSender);
  }
}
