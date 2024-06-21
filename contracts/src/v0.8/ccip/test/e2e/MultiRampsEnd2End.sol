// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import "../helpers/MerkleHelper.sol";
import "../offRamp/EVM2EVMMultiOffRampSetup.t.sol";
import "../onRamp/EVM2EVMMultiOnRampSetup.t.sol";

/// @notice This E2E test implements the following scenario:
/// 1. Send multiple messages from multiple source chains to a single destination chain (2 messages from source chain 1 and 1 from
/// source chain 2).
/// 2. Commit multiple merkle roots (1 for each source chain).
/// 3. Batch execute all the committed messages.
contract MultiRampsE2E is EVM2EVMMultiOnRampSetup, EVM2EVMMultiOffRampSetup {
  using Internal for Internal.EVM2EVMMessage;

  Router internal s_sourceRouter2;
  EVM2EVMMultiOnRampHelper internal s_onRamp2;
  TokenAdminRegistry internal s_tokenAdminRegistry2;

  bytes32 internal s_metadataHash2;

  mapping(address destPool => address sourcePool) internal s_sourcePoolByDestPool;

  function setUp() public virtual override(EVM2EVMMultiOnRampSetup, EVM2EVMMultiOffRampSetup) {
    EVM2EVMMultiOnRampSetup.setUp();
    EVM2EVMMultiOffRampSetup.setUp();

    // Deploy new source router for the new source chain
    s_sourceRouter2 = new Router(s_sourceRouter.getWrappedNative(), address(s_mockRMN));

    // Deploy new TokenAdminRegistry for the new source chain
    s_tokenAdminRegistry2 = new TokenAdminRegistry();

    // Deploy new token pools and set them on the new TokenAdminRegistry
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      address token = s_sourceTokens[i];
      address pool = address(
        new LockReleaseTokenPool(IERC20(token), new address[](0), address(s_mockRMN), true, address(s_sourceRouter2))
      );

      s_sourcePoolByDestPool[s_destPoolBySourceToken[token]] = pool;

      _setPool(
        s_tokenAdminRegistry2, token, pool, DEST_CHAIN_SELECTOR, s_destPoolByToken[s_destTokens[i]], s_destTokens[i]
      );
    }

    for (uint256 i = 0; i < s_destTokens.length; ++i) {
      address token = s_destTokens[i];
      address pool = s_destPoolByToken[token];

      _setPool(
        s_tokenAdminRegistry2, token, pool, SOURCE_CHAIN_SELECTOR + 1, s_sourcePoolByDestPool[pool], s_sourceTokens[i]
      );
    }

    // Deploy the new source chain onramp
    // Outsource to shared helper function with EVM2EVMMultiOnRampSetup
    (s_onRamp2, s_metadataHash2) =
      _deployOnRamp(SOURCE_CHAIN_SELECTOR + 1, address(s_sourceRouter2), address(s_tokenAdminRegistry2));

    // Enable destination chain on new source chain router
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(s_onRamp2)});
    s_sourceRouter2.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    // Deploy offramp
    _deployOffRamp(s_destRouter, s_mockRMN);

    // Enable source chains on offramp
    EVM2EVMMultiOffRamp.SourceChainConfigArgs[] memory sourceChainConfigs =
      new EVM2EVMMultiOffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      prevOffRamp: address(0),
      onRamp: address(s_onRamp)
    });
    sourceChainConfigs[1] = EVM2EVMMultiOffRamp.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
      isEnabled: true,
      prevOffRamp: address(0),
      onRamp: address(s_onRamp2)
    });

    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);
  }

  function test_E2E_3MessagesSuccess_gas() public {
    vm.pauseGasMetering();
    IERC20 token0 = IERC20(s_sourceTokens[0]);
    IERC20 token1 = IERC20(s_sourceTokens[1]);
    uint256 balance0Pre = token0.balanceOf(OWNER);
    uint256 balance1Pre = token1.balanceOf(OWNER);

    // Send messages
    Internal.EVM2EVMMessage[] memory messages1 = new Internal.EVM2EVMMessage[](2);
    messages1[0] = _sendRequest(1, SOURCE_CHAIN_SELECTOR, 1, s_metadataHash, s_sourceRouter, s_tokenAdminRegistry);
    messages1[1] = _sendRequest(2, SOURCE_CHAIN_SELECTOR, 2, s_metadataHash, s_sourceRouter, s_tokenAdminRegistry);
    Internal.EVM2EVMMessage[] memory messages2 = new Internal.EVM2EVMMessage[](1);
    messages2[0] =
      _sendRequest(1, SOURCE_CHAIN_SELECTOR + 1, 1, s_metadataHash2, s_sourceRouter2, s_tokenAdminRegistry2);

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, _generateTokenMessage());
    // Asserts that the tokens have been sent and the fee has been paid.
    assertEq(
      balance0Pre - (messages1.length + messages2.length) * (i_tokenAmount0 + expectedFee), token0.balanceOf(OWNER)
    );
    assertEq(balance1Pre - (messages1.length + messages2.length) * i_tokenAmount1, token1.balanceOf(OWNER));

    // Commit
    bytes32[] memory hashedMessages1 = new bytes32[](2);
    hashedMessages1[0] = messages1[0]._hash(s_metadataHash);
    messages1[0].messageId = hashedMessages1[0];
    hashedMessages1[1] = messages1[1]._hash(s_metadataHash);
    messages1[1].messageId = hashedMessages1[1];
    bytes32[] memory hashedMessages2 = new bytes32[](1);
    hashedMessages2[0] = messages2[0]._hash(s_metadataHash2);
    messages2[0].messageId = hashedMessages2[0];

    bytes32[] memory merkleRoots = new bytes32[](2);
    merkleRoots[0] = MerkleHelper.getMerkleRoot(hashedMessages1);
    merkleRoots[1] = MerkleHelper.getMerkleRoot(hashedMessages2);

    EVM2EVMMultiOffRamp.MerkleRoot[] memory roots = new EVM2EVMMultiOffRamp.MerkleRoot[](2);
    roots[0] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: EVM2EVMMultiOffRamp.Interval(messages1[0].sequenceNumber, messages1[1].sequenceNumber),
      merkleRoot: merkleRoots[0]
    });
    roots[1] = EVM2EVMMultiOffRamp.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
      interval: EVM2EVMMultiOffRamp.Interval(messages2[0].sequenceNumber, messages2[0].sequenceNumber),
      merkleRoot: merkleRoots[1]
    });

    EVM2EVMMultiOffRamp.CommitReport memory report =
      EVM2EVMMultiOffRamp.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.resumeGasMetering();
    _commit(report, ++s_latestSequenceNumber);
    vm.pauseGasMetering();

    bytes32[] memory proofs = new bytes32[](0);
    bytes32[] memory hashedLeaves = new bytes32[](1);
    hashedLeaves[0] = merkleRoots[0];
    uint256 timestamp = s_offRamp.verify(SOURCE_CHAIN_SELECTOR, hashedLeaves, proofs, 2 ** 2 - 1);
    assertEq(BLOCK_TIME, timestamp);
    hashedLeaves[0] = merkleRoots[1];
    timestamp = s_offRamp.verify(SOURCE_CHAIN_SELECTOR + 1, hashedLeaves, proofs, 2 ** 2 - 1);
    assertEq(BLOCK_TIME, timestamp);

    // We change the block time so when execute would e.g. use the current
    // block time instead of the committed block time the value would be
    // incorrect in the checks below.
    vm.warp(BLOCK_TIME + 2000);

    // Execute
    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR,
      messages1[0].sequenceNumber,
      messages1[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR,
      messages1[1].sequenceNumber,
      messages1[1].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    vm.expectEmit();
    emit EVM2EVMMultiOffRamp.ExecutionStateChanged(
      SOURCE_CHAIN_SELECTOR + 1,
      messages2[0].sequenceNumber,
      messages2[0].messageId,
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    Internal.ExecutionReportSingleChain[] memory reports = new Internal.ExecutionReportSingleChain[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR + 1, messages2);

    vm.resumeGasMetering();
    s_offRamp.batchExecute(reports, new uint256[][](0));
  }

  function _sendRequest(
    uint64 expectedSeqNum,
    uint64 sourceChainSelector,
    uint64 nonce,
    bytes32 metadataHash,
    Router router,
    TokenAdminRegistry tokenAdminRegistry
  ) public returns (Internal.EVM2EVMMessage memory) {
    Client.EVM2AnyMessage memory message = _generateTokenMessage();
    uint256 expectedFee = router.getFee(DEST_CHAIN_SELECTOR, message);

    IERC20(s_sourceTokens[0]).approve(address(router), i_tokenAmount0 + expectedFee);
    IERC20(s_sourceTokens[1]).approve(address(router), i_tokenAmount1);

    message.receiver = abi.encode(address(s_receiver));
    Internal.EVM2EVMMessage memory msgEvent = _messageToEvent(
      message, sourceChainSelector, expectedSeqNum, nonce, expectedFee, OWNER, metadataHash, tokenAdminRegistry
    );

    vm.expectEmit();
    emit EVM2EVMMultiOnRamp.CCIPSendRequested(DEST_CHAIN_SELECTOR, msgEvent);

    vm.resumeGasMetering();
    router.ccipSend(DEST_CHAIN_SELECTOR, message);
    vm.pauseGasMetering();

    return msgEvent;
  }
}
