// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRouter} from "../../../interfaces/IRouter.sol";

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {NonceManager} from "../../../NonceManager.sol";
import {Router} from "../../../Router.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {TokenAdminRegistry} from "../../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {FeeQuoterFeeSetup} from "../../feeQuoter/FeeQuoterSetup.t.sol";
import {OnRampHelper} from "../../helpers/OnRampHelper.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OnRampSetup is FeeQuoterFeeSetup {
  address internal constant FEE_AGGREGATOR = 0xa33CDB32eAEce34F6affEfF4899cef45744EDea3;

  bytes32 internal s_metadataHash;

  OnRampHelper internal s_onRamp;
  NonceManager internal s_outboundNonceManager;

  function setUp() public virtual override {
    super.setUp();

    s_outboundNonceManager = new NonceManager(new address[](0));
    (s_onRamp, s_metadataHash) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry)
    );

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(s_onRamp)});

    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: makeAddr("offRamp0")});
    offRampUpdates[1] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: makeAddr("offRamp1")});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);

    // Pre approve the first token so the gas estimates of the tests
    // only cover actual gas usage from the ramps
    IERC20(s_sourceTokens[0]).approve(address(s_sourceRouter), 2 ** 128);
    IERC20(s_sourceTokens[1]).approve(address(s_sourceRouter), 2 ** 128);
  }

  /// @dev a helper function to compose EVM2AnyRampMessage messages
  /// @dev it is assumed that LINK is the payment token because feeTokenAmount == feeValueJuels
  function _messageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    address originalSender
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    return _messageToEvent(
      message,
      seqNum,
      nonce,
      feeTokenAmount, // fee paid
      feeTokenAmount, // conversion to jules is the same
      originalSender
    );
  }

  function _messageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    uint256 feeValueJuels,
    address originalSender
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    bytes4 chainFamilySelector = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR).chainFamilySelector;
    if (chainFamilySelector == Internal.CHAIN_FAMILY_SELECTOR_EVM) {
      return _evmMessageToEvent(
        message,
        SOURCE_CHAIN_SELECTOR,
        seqNum,
        nonce,
        feeTokenAmount,
        feeValueJuels,
        originalSender,
        s_metadataHash,
        s_tokenAdminRegistry
      );
    } else {
      return _svmMessageToEvent(
        message,
        SOURCE_CHAIN_SELECTOR,
        seqNum,
        nonce,
        feeTokenAmount,
        feeValueJuels,
        originalSender,
        s_metadataHash,
        s_tokenAdminRegistry
      );
    }
  }

  function _evmMessageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 sourceChainSelector,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    uint256 feeValueJuels,
    address originalSender,
    bytes32 metadataHash,
    TokenAdminRegistry tokenAdminRegistry
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    Client.GenericExtraArgsV2 memory extraArgs =
      s_feeQuoter.parseEVMExtraArgsFromBytes(message.extraArgs, DEST_CHAIN_SELECTOR);

    Internal.EVM2AnyRampMessage memory messageEvent = Internal.EVM2AnyRampMessage({
      header: Internal.RampMessageHeader({
        messageId: "",
        sourceChainSelector: sourceChainSelector,
        destChainSelector: DEST_CHAIN_SELECTOR,
        sequenceNumber: seqNum,
        nonce: extraArgs.allowOutOfOrderExecution ? 0 : nonce
      }),
      sender: originalSender,
      data: message.data,
      receiver: message.receiver,
      extraArgs: Client._argsToBytes(extraArgs),
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      feeValueJuels: feeValueJuels,
      tokenAmounts: new Internal.EVM2AnyTokenTransfer[](message.tokenAmounts.length)
    });

    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      messageEvent.tokenAmounts[i] =
        _getSourceTokenData(message.tokenAmounts[i], tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    }

    messageEvent.header.messageId = Internal._hash(messageEvent, metadataHash);
    return messageEvent;
  }

  function _svmMessageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 sourceChainSelector,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    uint256 feeValueJuels,
    address originalSender,
    bytes32 metadataHash,
    TokenAdminRegistry tokenAdminRegistry
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    Client.SVMExtraArgsV1 memory extraArgs =
      s_feeQuoter.parseSVMExtraArgsFromBytes(message.extraArgs, s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR));

    Internal.EVM2AnyRampMessage memory messageEvent = Internal.EVM2AnyRampMessage({
      header: Internal.RampMessageHeader({
        messageId: "",
        sourceChainSelector: sourceChainSelector,
        destChainSelector: DEST_CHAIN_SELECTOR,
        sequenceNumber: seqNum,
        nonce: extraArgs.allowOutOfOrderExecution ? 0 : nonce
      }),
      sender: originalSender,
      data: message.data,
      receiver: message.receiver,
      extraArgs: Client._svmArgsToBytes(extraArgs),
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      feeValueJuels: feeValueJuels,
      tokenAmounts: new Internal.EVM2AnyTokenTransfer[](message.tokenAmounts.length)
    });

    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      messageEvent.tokenAmounts[i] =
        _getSourceTokenData(message.tokenAmounts[i], tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    }

    messageEvent.header.messageId = Internal._hash(messageEvent, metadataHash);
    return messageEvent;
  }

  function _generateDynamicOnRampConfig(
    address feeQuoter
  ) internal pure returns (OnRamp.DynamicConfig memory) {
    return OnRamp.DynamicConfig({
      feeQuoter: feeQuoter,
      reentrancyGuardEntered: false,
      messageInterceptor: address(0),
      feeAggregator: FEE_AGGREGATOR,
      allowlistAdmin: address(0)
    });
  }

  function _generateDestChainConfigArgs(
    IRouter router
  ) internal pure returns (OnRamp.DestChainConfigArgs[] memory) {
    OnRamp.DestChainConfigArgs[] memory destChainConfigs = new OnRamp.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      OnRamp.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, router: router, allowlistEnabled: false});
    return destChainConfigs;
  }

  function _deployOnRamp(
    uint64 sourceChainSelector,
    IRouter router,
    address nonceManager,
    address tokenAdminRegistry
  ) internal returns (OnRampHelper, bytes32 metadataHash) {
    OnRampHelper onRamp = new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: sourceChainSelector,
        rmnRemote: s_mockRMNRemote,
        nonceManager: nonceManager,
        tokenAdminRegistry: tokenAdminRegistry
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(router)
    );

    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(onRamp);

    NonceManager(nonceManager).applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );

    return (
      onRamp,
      keccak256(abi.encode(Internal.EVM_2_ANY_MESSAGE_HASH, sourceChainSelector, DEST_CHAIN_SELECTOR, address(onRamp)))
    );
  }
}
