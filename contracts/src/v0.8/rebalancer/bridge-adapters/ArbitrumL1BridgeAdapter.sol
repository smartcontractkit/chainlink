// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBridgeAdapter} from "../interfaces/IBridge.sol";

import {IL1GatewayRouter} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/ethereum/gateway/IL1GatewayRouter.sol";
import {IGatewayRouter} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/libraries/gateway/IGatewayRouter.sol";
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
contract ArbitrumL1BridgeAdapter is IBridgeAdapter {
  using SafeERC20 for IERC20;

  IL1GatewayRouter internal immutable i_l1GatewayRouter;
  IOutbox internal immutable i_l1Outbox;

  // TODO not static?
  uint256 public constant MAX_GAS = 100_000;
  uint256 public constant GAS_PRICE_BID = 300_000_000;
  uint256 public constant MAX_SUBMISSION_COST = 8e14;

  // Nonce to use for L2 deposits to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  error NoGatewayForToken(address token);

  constructor(IL1GatewayRouter l1GatewayRouter, IOutbox l1Outbox) {
    if (address(l1GatewayRouter) == address(0) || address(l1Outbox) == address(0)) {
      revert BridgeAddressCannotBeZero();
    }
    i_l1GatewayRouter = l1GatewayRouter;
    i_l1Outbox = l1Outbox;
  }

  /// @inheritdoc IBridgeAdapter
  function sendERC20(
    address localToken,
    address /* remoteToken */,
    address recipient,
    uint256 amount
  ) external payable override returns (bytes memory) {
    // receive the token transfer from the msg.sender
    IERC20(localToken).safeTransferFrom(msg.sender, address(this), amount);

    // Note: the gateway router could return 0x0 for the gateway address
    // if that token is not yet registered
    address gateway = IGatewayRouter(address(i_l1GatewayRouter)).getGateway(localToken);
    if (gateway == address(0)) {
      revert NoGatewayForToken(localToken);
    }

    // approve the gateway to transfer the token amount sent to the adapter
    IERC20(localToken).safeApprove(gateway, amount);

    uint256 wantedNativeFeeCoin = getBridgeFeeInNative();
    if (msg.value < wantedNativeFeeCoin) {
      revert InsufficientEthValue(wantedNativeFeeCoin, msg.value);
    }

    // TODO: return data bombs?
    // The router will route the call to the gateway that we approved
    // above. The gateway will then transfer the tokens to the L2.
    return
      i_l1GatewayRouter.outboundTransferCustomRefund{value: msg.value}(
        localToken,
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

  /// @notice Finalize an L2 -> L1 transfer.
  /// @param remoteSender sender if original message (i.e., caller of ArbSys.sendTxToL1)
  /// @param localReceiver destination address for L1 contract call
  function finalizeWithdrawERC20(
    address remoteSender,
    address localReceiver,
    bytes calldata arbitrumFinalizationPayload
  ) external {
    ArbitrumFinalizationPayload memory payload = abi.decode(arbitrumFinalizationPayload, (ArbitrumFinalizationPayload));
    i_l1Outbox.executeTransaction(
      payload.proof,
      payload.index,
      remoteSender,
      localReceiver,
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
