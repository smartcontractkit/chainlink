// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IRouterClient} from "../interfaces/IRouterClient.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Client} from "../libraries/Client.sol";
import {CCIPReceiver} from "./CCIPReceiver.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

// @notice Example of a client which supports EVM/non-EVM chains
// @dev If chain specific logic is required for different chain families (e.g. particular
// decoding the bytes sender for authorization checks), it may be required to point to a helper
// authorization contract unless all chain families are known up front.
// @dev If contract does not implement IAny2EVMMessageReceiver and IERC165,
// and tokens are sent to it, ccipReceive will not be called but tokens will be transferred.
// @dev If the client is upgradeable you have significantly more flexibility and
// can avoid storage based options like the below contract uses. However it's
// worth carefully considering how the trust assumptions of your client dapp will
// change if you introduce upgradeability. An immutable dapp building on top of CCIP
// like the example below will inherit the trust properties of CCIP (i.e. the oracle network).
// @dev The receiver's are encoded offchain and passed as direct arguments to permit supporting
// new chain family receivers (e.g. a Solana encoded receiver address) without upgrading.
contract CCIPClientExample is CCIPReceiver, OwnerIsCreator {
  error InvalidChain(uint64 chainSelector);

  event MessageSent(bytes32 messageId);
  event MessageReceived(bytes32 messageId);

  // Current feeToken
  IERC20 public s_feeToken;
  // Below is a simplistic example (same params for all messages) of using storage to allow for new options without
  // upgrading the dapp. Note that extra args are chain family specific (e.g. gasLimit is EVM specific etc.).
  // and will always be backwards compatible i.e. upgrades are opt-in.
  // Offchain we can compute the V1 extraArgs:
  //    Client.EVMExtraArgsV1 memory extraArgs = Client.EVMExtraArgsV1({gasLimit: 300_000});
  //    bytes memory encodedV1ExtraArgs = Client._argsToBytes(extraArgs);
  // Then later compute V2 extraArgs, for example if a refund feature was added:
  //    Client.EVMExtraArgsV2 memory extraArgs = Client.EVMExtraArgsV2({gasLimit: 300_000, destRefundAddress: 0x1234});
  //    bytes memory encodedV2ExtraArgs = Client._argsToBytes(extraArgs);
  // and update storage with the new args.
  // If different options are required for different messages, for example different gas limits,
  // one can simply key based on (chainSelector, messageType) instead of only chainSelector.
  mapping(uint64 destChainSelector => bytes extraArgsBytes) public s_chains;

  constructor(IRouterClient router, IERC20 feeToken) CCIPReceiver(address(router)) {
    s_feeToken = feeToken;
    s_feeToken.approve(address(router), type(uint256).max);
  }

  function enableChain(uint64 chainSelector, bytes memory extraArgs) external onlyOwner {
    s_chains[chainSelector] = extraArgs;
  }

  function disableChain(
    uint64 chainSelector
  ) external onlyOwner {
    delete s_chains[chainSelector];
  }

  function ccipReceive(
    Client.Any2EVMMessage calldata message
  ) external virtual override onlyRouter validChain(message.sourceChainSelector) {
    // Extremely important to ensure only router calls this.
    // Tokens in message if any will be transferred to this contract
    // TODO: Validate sender/origin chain and process message and/or tokens.
    _ccipReceive(message);
  }

  function _ccipReceive(
    Client.Any2EVMMessage memory message
  ) internal override {
    emit MessageReceived(message.messageId);
  }

  /// @notice sends data to receiver on dest chain. Assumes address(this) has sufficient native asset.
  function sendDataPayNative(
    uint64 destChainSelector,
    bytes memory receiver,
    bytes memory data
  ) external validChain(destChainSelector) {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](0);
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: receiver,
      data: data,
      tokenAmounts: tokenAmounts,
      extraArgs: s_chains[destChainSelector],
      feeToken: address(0) // We leave the feeToken empty indicating we'll pay raw native.
    });
    bytes32 messageId = IRouterClient(i_ccipRouter).ccipSend{
      value: IRouterClient(i_ccipRouter).getFee(destChainSelector, message)
    }(destChainSelector, message);
    emit MessageSent(messageId);
  }

  /// @notice sends data to receiver on dest chain. Assumes address(this) has sufficient feeToken.
  function sendDataPayFeeToken(
    uint64 destChainSelector,
    bytes memory receiver,
    bytes memory data
  ) external validChain(destChainSelector) {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](0);
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: receiver,
      data: data,
      tokenAmounts: tokenAmounts,
      extraArgs: s_chains[destChainSelector],
      feeToken: address(s_feeToken)
    });
    // Optional uint256 fee = i_ccipRouter.getFee(destChainSelector, message);
    // Can decide if fee is acceptable.
    // address(this) must have sufficient feeToken or the send will revert.
    bytes32 messageId = IRouterClient(i_ccipRouter).ccipSend(destChainSelector, message);
    emit MessageSent(messageId);
  }

  /// @notice sends data to receiver on dest chain. Assumes address(this) has sufficient native token.
  function sendDataAndTokens(
    uint64 destChainSelector,
    bytes memory receiver,
    bytes memory data,
    Client.EVMTokenAmount[] memory tokenAmounts
  ) external validChain(destChainSelector) {
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      IERC20(tokenAmounts[i].token).transferFrom(msg.sender, address(this), tokenAmounts[i].amount);
      IERC20(tokenAmounts[i].token).approve(i_ccipRouter, tokenAmounts[i].amount);
    }
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: receiver,
      data: data,
      tokenAmounts: tokenAmounts,
      extraArgs: s_chains[destChainSelector],
      feeToken: address(s_feeToken)
    });
    // Optional uint256 fee = i_ccipRouter.getFee(destChainSelector, message);
    // Can decide if fee is acceptable.
    // address(this) must have sufficient feeToken or the send will revert.
    bytes32 messageId = IRouterClient(i_ccipRouter).ccipSend(destChainSelector, message);
    emit MessageSent(messageId);
  }

  // @notice user sends tokens to a receiver
  // Approvals can be optimized with a whitelist of tokens and inf approvals if desired.
  function sendTokens(
    uint64 destChainSelector,
    bytes memory receiver,
    Client.EVMTokenAmount[] memory tokenAmounts
  ) external validChain(destChainSelector) {
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      IERC20(tokenAmounts[i].token).transferFrom(msg.sender, address(this), tokenAmounts[i].amount);
      IERC20(tokenAmounts[i].token).approve(i_ccipRouter, tokenAmounts[i].amount);
    }
    bytes memory data;
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: receiver,
      data: data,
      tokenAmounts: tokenAmounts,
      extraArgs: s_chains[destChainSelector],
      feeToken: address(s_feeToken)
    });
    // Optional uint256 fee = i_ccipRouter.getFee(destChainSelector, message);
    // Can decide if fee is acceptable.
    // address(this) must have sufficient feeToken or the send will revert.
    bytes32 messageId = IRouterClient(i_ccipRouter).ccipSend(destChainSelector, message);
    emit MessageSent(messageId);
  }

  modifier validChain(
    uint64 chainSelector
  ) {
    if (s_chains[chainSelector].length == 0) revert InvalidChain(chainSelector);
    _;
  }
}
