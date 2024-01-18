// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IL1BridgeAdapter} from "../interfaces/IBridge.sol";

import {IL1GatewayRouter} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/ethereum/gateway/IL1GatewayRouter.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

interface IOutbox {
  /**
   * @notice Executes a messages in an Outbox entry.
   * @dev Reverts if dispute period hasn't expired, since the outbox entry
   *      is only created once the rollup confirms the respective assertion.
   * @dev it is not possible to execute any L2-to-L1 transaction which contains data
   *      to a contract address without any code (as enforced by the Bridge contract).
   * @param proof Merkle proof of message inclusion in send root
   * @param index Merkle path to message
   * @param l2Sender sender if original message (i.e., caller of ArbSys.sendTxToL1)
   * @param to destination address for L1 contract call
   * @param l2Block l2 block number at which sendTxToL1 call was made
   * @param l1Block l1 block number at which sendTxToL1 call was made
   * @param l2Timestamp l2 Timestamp at which sendTxToL1 call was made
   * @param value wei in L1 message
   * @param data abi-encoded L1 message data
   */
  function executeTransaction(
    bytes32[] calldata proof,
    uint256 index,
    address l2Sender,
    address to,
    uint256 l2Block,
    uint256 l1Block,
    uint256 l2Timestamp,
    uint256 value,
    bytes calldata data
  ) external;
}

/// @notice Arbitrum L1 Bridge adapter
/// @dev Auto unwraps and re-wraps wrapped eth in the bridge.
contract ArbitrumL1BridgeAdapter is IL1BridgeAdapter {
  using SafeERC20 for IERC20;

  IL1GatewayRouter internal immutable i_l1GatewayRouter;
  address internal immutable i_l1ERC20Gateway;
  IOutbox internal immutable i_l1Outbox;

  // TODO not static?
  uint256 public constant MAX_GAS = 100_000;
  uint256 public constant GAS_PRICE_BID = 300_000_000;
  uint256 public constant MAX_SUBMISSION_COST = 8e14;

  // Nonce to use for L2 deposits to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  constructor(IL1GatewayRouter l1GatewayRouter, IOutbox l1Outbox, address l1ERC20Gateway) {
    if (
      address(l1GatewayRouter) == address(0) || address(l1Outbox) == address(0) || address(l1ERC20Gateway) == address(0)
    ) {
      revert BridgeAddressCannotBeZero();
    }
    i_l1GatewayRouter = l1GatewayRouter;
    i_l1Outbox = l1Outbox;
    i_l1ERC20Gateway = l1ERC20Gateway;
  }

  function sendERC20(address l1Token, address, address recipient, uint256 amount) external payable {
    IERC20(l1Token).safeTransferFrom(msg.sender, address(this), amount);

    IERC20(l1Token).approve(i_l1ERC20Gateway, amount);

    uint256 wantedNativeFeeCoin = getBridgeFeeInNative();
    if (msg.value < wantedNativeFeeCoin) {
      revert InsufficientEthValue(wantedNativeFeeCoin, msg.value);
    }

    i_l1GatewayRouter.outboundTransferCustomRefund{value: msg.value}(
      l1Token,
      recipient,
      recipient,
      amount,
      MAX_GAS,
      GAS_PRICE_BID,
      abi.encode(MAX_SUBMISSION_COST, bytes(""))
    );
  }

  function getBridgeFeeInNative() public pure returns (uint256) {
    return MAX_SUBMISSION_COST + MAX_GAS * GAS_PRICE_BID;
  }

  /// @param proof Merkle proof of message inclusion in send root
  /// @param index Merkle path to message
  /// @param l2Block l2 block number at which sendTxToL1 call was made
  /// @param l1Block l1 block number at which sendTxToL1 call was made
  /// @param l2Timestamp l2 Timestamp at which sendTxToL1 call was made
  /// @param value wei in L1 message
  /// @param data abi-encoded L1 message data
  struct ArbitrumFinalizationPayload {
    bytes32[] proof;
    uint256 index;
    uint256 l2Block;
    uint256 l1Block;
    uint256 l2Timestamp;
    uint256 value;
    bytes data;
  }

  /// @param l2Sender sender if original message (i.e., caller of ArbSys.sendTxToL1)
  /// @param l1Receiver destination address for L1 contract call
  function finalizeWithdrawERC20FromL2(
    address l2Sender,
    address l1Receiver,
    bytes calldata arbitrumFinalizationPayload
  ) external {
    ArbitrumFinalizationPayload memory payload = abi.decode(arbitrumFinalizationPayload, (ArbitrumFinalizationPayload));
    i_l1Outbox.executeTransaction(
      payload.proof,
      payload.index,
      l2Sender,
      l1Receiver,
      payload.l2Block,
      payload.l1Block,
      payload.l2Timestamp,
      payload.value,
      payload.data
    );
  }

  function getL2Token(address l1Token) external view returns (address) {
    return i_l1GatewayRouter.calculateL2TokenAddress(l1Token);
  }
}
