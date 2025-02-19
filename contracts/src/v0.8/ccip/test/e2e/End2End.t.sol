// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMNRemote} from "../../interfaces/IRMNRemote.sol";

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {NonceManager} from "../../NonceManager.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {MerkleHelper} from "../helpers/MerkleHelper.sol";
import {OnRampHelper} from "../helpers/OnRampHelper.sol";
import {OffRampSetup} from "../offRamp/OffRamp/OffRampSetup.t.sol";
import {OnRampSetup} from "../onRamp/OnRamp/OnRampSetup.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/// @notice This E2E test implements the following scenario:
/// 1. Send multiple messages from multiple source chains to a single destination chain (2 messages from source chain
/// 1 and 1 from source chain 2).
/// 2. Commit multiple merkle roots (1 for each source chain).
/// 3. Batch execute all the committed messages.
contract E2E is OnRampSetup, OffRampSetup {
  using Internal for Internal.Any2EVMRampMessage;

  uint256 internal constant TOKEN_AMOUNT_1 = 9;
  uint256 internal constant TOKEN_AMOUNT_2 = 7;

  Router internal s_sourceRouter2;
  OnRampHelper internal s_onRamp2;
  TokenAdminRegistry internal s_tokenAdminRegistry2;
  NonceManager internal s_nonceManager2;

  bytes32 internal s_metadataHash2;

  mapping(address destPool => address sourcePool) internal s_sourcePoolByDestPool;

  function setUp() public virtual override(OnRampSetup, OffRampSetup) {
    OnRampSetup.setUp();
    OffRampSetup.setUp();

    // Deploy new source router for the new source chain
    s_sourceRouter2 = new Router(s_sourceRouter.getWrappedNative(), address(s_mockRMNRemote));

    // Deploy new TokenAdminRegistry for the new source chain
    s_tokenAdminRegistry2 = new TokenAdminRegistry();

    // Deploy new token pools and set them on the new TokenAdminRegistry
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      address token = s_sourceTokens[i];
      address pool = address(
        new LockReleaseTokenPool(
          IERC20(token),
          DEFAULT_TOKEN_DECIMALS,
          new address[](0),
          address(s_mockRMNRemote),
          true,
          address(s_sourceRouter2)
        )
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

    s_nonceManager2 = new NonceManager(new address[](0));

    (
      // Deploy the new source chain onRamp
      // Outsource to shared helper function with OnRampSetup
      s_onRamp2,
      s_metadataHash2
    ) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR + 1, s_sourceRouter2, address(s_nonceManager2), address(s_tokenAdminRegistry2)
    );

    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(s_onRamp2);
    s_nonceManager2.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );

    // Enable destination chain on new source chain router
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(s_onRamp2)});
    s_sourceRouter2.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));

    // Deploy offRamp
    _deployOffRamp(s_mockRMNRemote, s_inboundNonceManager);

    // Enable source chains on offRamp
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      // Must match OnRamp address
      onRamp: abi.encode(address(s_onRamp)),
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
      isEnabled: true,
      onRamp: abi.encode(address(s_onRamp2)),
      isRMNVerificationDisabled: false
    });

    _setupMultipleOffRampsFromConfigs(sourceChainConfigs);
  }

  function test_E2E_3MessagesMMultiOffRampSuccess_gas() public {
    vm.pauseGasMetering();

    Internal.Any2EVMRampMessage[] memory messages1 = new Internal.Any2EVMRampMessage[](2);
    Internal.Any2EVMRampMessage[] memory messages2 = new Internal.Any2EVMRampMessage[](1);

    // Scoped to sending to reduce stack pressure
    {
      IERC20 token0 = IERC20(s_sourceTokens[0]);
      IERC20 token1 = IERC20(s_sourceTokens[1]);

      uint256 balance0Pre = token0.balanceOf(OWNER);
      uint256 balance1Pre = token1.balanceOf(OWNER);

      // Send messages
      messages1[0] = _sendRequest(1, SOURCE_CHAIN_SELECTOR, 1, s_metadataHash, s_sourceRouter, s_tokenAdminRegistry);
      messages1[1] = _sendRequest(2, SOURCE_CHAIN_SELECTOR, 2, s_metadataHash, s_sourceRouter, s_tokenAdminRegistry);
      messages2[0] =
        _sendRequest(1, SOURCE_CHAIN_SELECTOR + 1, 1, s_metadataHash2, s_sourceRouter2, s_tokenAdminRegistry2);

      uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, _generateTokenMessage());
      // Asserts that the tokens have been sent and the fee has been paid.
      assertEq(
        balance0Pre - (messages1.length + messages2.length) * (TOKEN_AMOUNT_1 + expectedFee), token0.balanceOf(OWNER)
      );
      assertEq(balance1Pre - (messages1.length + messages2.length) * TOKEN_AMOUNT_2, token1.balanceOf(OWNER));
    }

    // Commit

    bytes32[] memory merkleRoots = new bytes32[](2);

    // Scoped to commit to reduce stack pressure
    {
      bytes32[] memory hashedMessages1 = new bytes32[](2);
      hashedMessages1[0] = _hashMessage(messages1[0], abi.encode(address(s_onRamp)));
      hashedMessages1[1] = _hashMessage(messages1[1], abi.encode(address(s_onRamp)));
      bytes32[] memory hashedMessages2 = new bytes32[](1);
      hashedMessages2[0] = _hashMessage(messages2[0], abi.encode(address(s_onRamp2)));

      merkleRoots[0] = MerkleHelper.getMerkleRoot(hashedMessages1);
      merkleRoots[1] = MerkleHelper.getMerkleRoot(hashedMessages2);

      // TODO make these real sigs
      IRMNRemote.Signature[] memory rmnSignatures = new IRMNRemote.Signature[](0);

      Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](2);
      roots[0] = Internal.MerkleRoot({
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRampAddress: abi.encode(address(s_onRamp)),
        minSeqNr: messages1[0].header.sequenceNumber,
        maxSeqNr: messages1[1].header.sequenceNumber,
        merkleRoot: merkleRoots[0]
      });
      roots[1] = Internal.MerkleRoot({
        sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
        onRampAddress: abi.encode(address(s_onRamp2)),
        minSeqNr: messages2[0].header.sequenceNumber,
        maxSeqNr: messages2[0].header.sequenceNumber,
        merkleRoot: merkleRoots[1]
      });

      OffRamp.CommitReport memory report = OffRamp.CommitReport({
        priceUpdates: _getEmptyPriceUpdates(),
        blessedMerkleRoots: roots,
        unblessedMerkleRoots: new Internal.MerkleRoot[](0),
        rmnSignatures: rmnSignatures
      });

      vm.resumeGasMetering();
      _commit(report, ++s_latestSequenceNumber);
      vm.pauseGasMetering();
    }

    // Scoped to RMN and verify to reduce stack pressure
    {
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
    }

    // Execute

    Internal.ExecutionReport[] memory reports = new Internal.ExecutionReport[](2);
    reports[0] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR, messages1);
    reports[1] = _generateReportFromMessages(SOURCE_CHAIN_SELECTOR + 1, messages2);

    vm.resumeGasMetering();
    vm.recordLogs();
    _execute(reports);

    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR,
      messages1[0].header.sequenceNumber,
      messages1[0].header.messageId,
      _hashMessage(messages1[0], abi.encode(address(s_onRamp))),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR,
      messages1[1].header.sequenceNumber,
      messages1[1].header.messageId,
      _hashMessage(messages1[1], abi.encode(address(s_onRamp))),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );

    _assertExecutionStateChangedEventLogs(
      SOURCE_CHAIN_SELECTOR + 1,
      messages2[0].header.sequenceNumber,
      messages2[0].header.messageId,
      _hashMessage(messages2[0], abi.encode(address(s_onRamp2))),
      Internal.MessageExecutionState.SUCCESS,
      ""
    );
  }

  function _sendRequest(
    uint64 expectedSeqNum,
    uint64 sourceChainSelector,
    uint64 nonce,
    bytes32 metadataHash,
    Router router,
    TokenAdminRegistry tokenAdminRegistry
  ) public returns (Internal.Any2EVMRampMessage memory) {
    Client.EVM2AnyMessage memory message = _generateTokenMessage();
    IERC20(s_sourceTokens[0]).approve(address(router), TOKEN_AMOUNT_1 + router.getFee(DEST_CHAIN_SELECTOR, message));
    IERC20(s_sourceTokens[1]).approve(address(router), TOKEN_AMOUNT_2);

    uint256 feeAmount = router.getFee(DEST_CHAIN_SELECTOR, message);

    message.receiver = abi.encode(address(s_receiver));
    Internal.EVM2AnyRampMessage memory msgEvent = _evmMessageToEvent(
      message, sourceChainSelector, expectedSeqNum, nonce, feeAmount, feeAmount, OWNER, metadataHash, tokenAdminRegistry
    );

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, expectedSeqNum, msgEvent);

    vm.resumeGasMetering();
    router.ccipSend(DEST_CHAIN_SELECTOR, message);
    vm.pauseGasMetering();

    Internal.Any2EVMTokenTransfer[] memory any2EVMTokenTransfer =
      new Internal.Any2EVMTokenTransfer[](message.tokenAmounts.length);

    for (uint256 i = 0; i < msgEvent.tokenAmounts.length; ++i) {
      any2EVMTokenTransfer[i] = Internal.Any2EVMTokenTransfer({
        sourcePoolAddress: abi.encode(msgEvent.tokenAmounts[i].sourcePoolAddress),
        destTokenAddress: abi.decode(msgEvent.tokenAmounts[i].destTokenAddress, (address)),
        extraData: msgEvent.tokenAmounts[i].extraData,
        amount: msgEvent.tokenAmounts[i].amount,
        destGasAmount: abi.decode(msgEvent.tokenAmounts[i].destExecData, (uint32))
      });
    }

    return Internal.Any2EVMRampMessage({
      header: Internal.RampMessageHeader({
        messageId: msgEvent.header.messageId,
        sourceChainSelector: sourceChainSelector,
        destChainSelector: DEST_CHAIN_SELECTOR,
        sequenceNumber: msgEvent.header.sequenceNumber,
        nonce: msgEvent.header.nonce
      }),
      sender: abi.encode(msgEvent.sender),
      data: msgEvent.data,
      receiver: abi.decode(msgEvent.receiver, (address)),
      gasLimit: s_feeQuoter.parseEVMExtraArgsFromBytes(msgEvent.extraArgs, DEST_CHAIN_SELECTOR).gasLimit,
      tokenAmounts: any2EVMTokenTransfer
    });
  }

  function _generateTokenMessage() public view returns (Client.EVM2AnyMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = TOKEN_AMOUNT_1;
    tokenAmounts[1].amount = TOKEN_AMOUNT_2;
    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }
}
