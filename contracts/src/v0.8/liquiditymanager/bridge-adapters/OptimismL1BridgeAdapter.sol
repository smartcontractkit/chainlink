// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBridgeAdapter} from "../interfaces/IBridge.sol";
import {IWrappedNative} from "../../ccip/interfaces/IWrappedNative.sol";
import {Types} from "../interfaces/optimism/Types.sol";
import {IOptimismPortal} from "../interfaces/optimism/IOptimismPortal.sol";

import {IL1StandardBridge} from "@eth-optimism/contracts/L1/messaging/IL1StandardBridge.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice OptimismL1BridgeAdapter implements IBridgeAdapter for the Optimism L1<=>L2 bridge.
/// @dev L1 -> L2 deposits are done via the depositERC20To and depositETHTo functions on the L1StandardBridge.
/// The amount of gas provided for the transaction must be buffered - the Optimism SDK recommends a 20% buffer.
/// The Optimism Bridge implements 2-step withdrawals from L2 to L1. Once a withdrawal transaction is included
/// in the L2 chain, it must be proven on L1 before it can be finalized. There is a buffer between the transaction
/// being posted on L2 before it can be proven, and similarly, there is a buffer in the time it takes to prove
/// the transaction before it can be finalized.
/// See https://blog.oplabs.co/two-step-withdrawals/ for more details on this mechanism.
/// @dev We have to unwrap WETH into ether before depositing it to L2. Therefore this bridge adapter bridges
/// WETH to ether. The receiver on L2 must wrap the ether back into WETH.
contract OptimismL1BridgeAdapter is IBridgeAdapter {
  using SafeERC20 for IERC20;

  /// @notice used when the action in the payload is invalid.
  error InvalidFinalizationAction();

  /// @notice Payload for proving a withdrawal from L2 on L1 via finalizeWithdrawERC20.
  /// @param withdrawalTransaction The withdrawal transaction, see its docstring for more details.
  /// @param l2OutputIndex The index of the output in the L2 block, or the dispute game index post fault proof upgrade.
  /// @param outputRootProof The inclusion proof of the L2ToL1MessagePasser contract's storage root.
  /// @param withdrawalProof The Merkle proof of the withdrawal key presence in the L2ToL1MessagePasser contract's state trie.
  struct OptimismProveWithdrawalPayload {
    Types.WithdrawalTransaction withdrawalTransaction;
    uint256 l2OutputIndex;
    Types.OutputRootProof outputRootProof;
    bytes[] withdrawalProof;
  }

  /// @notice Payload for finalizing a withdrawal from L2 on L1.
  /// Note that the withdrawal must be proven first before it can be finalized.
  /// @param withdrawalTransaction The withdrawal transaction, see its docstring for more details.
  struct OptimismFinalizationPayload {
    Types.WithdrawalTransaction withdrawalTransaction;
  }

  /// @notice The action to take when finalizing a withdrawal.
  /// Optimism implements two-step withdrawals, so we need to specify the action to take
  /// each time the finalizeWithdrawERC20 function is called.
  enum FinalizationAction {
    ProveWithdrawal,
    FinalizeWithdrawal
  }

  /// @notice Payload for interacting with the finalizeWithdrawERC20 function.
  /// Since Optimism has 2-step withdrawals, we cannot finalize and get the funds on L1 in the same transaction.
  /// @param action The action to take; either ProveWithdrawal or FinalizeWithdrawal.
  /// @param data The payload for the action. If ProveWithdrawal, it must be an abi-encoded OptimismProveWithdrawalPayload.
  ///        If FinalizeWithdrawal, it must be an abi-encoded OptimismFinalizationPayload.
  struct FinalizeWithdrawERC20Payload {
    FinalizationAction action;
    bytes data;
  }

  /// @dev Reference to the L1StandardBridge contract. Deposits to L2 go through this contract.
  IL1StandardBridge internal immutable i_L1Bridge;

  /// @dev Reference to the WrappedNative contract. Optimism bridges ether directly rather than WETH,
  /// so we need to unwrap WETH into ether before depositing it to L2.
  IWrappedNative internal immutable i_wrappedNative;

  /// @dev Reference to the OptimismPortal contract, which is used to prove and finalize withdrawals.
  IOptimismPortal internal immutable i_optimismPortal;

  /// @dev Nonce to use for L2 deposits to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  constructor(IL1StandardBridge l1Bridge, IWrappedNative wrappedNative, IOptimismPortal optimismPortal) {
    if (
      address(l1Bridge) == address(0) || address(wrappedNative) == address(0) || address(optimismPortal) == address(0)
    ) {
      revert BridgeAddressCannotBeZero();
    }
    i_L1Bridge = l1Bridge;
    i_wrappedNative = wrappedNative;
    i_optimismPortal = optimismPortal;
  }

  /// @notice The WETH withdraw requires this be present otherwise withdraws will fail.
  receive() external payable {}

  /// @inheritdoc IBridgeAdapter
  function sendERC20(
    address localToken,
    address remoteToken,
    address recipient,
    uint256 amount,
    bytes calldata /* bridgeSpecificPayload */
  ) external payable override returns (bytes memory) {
    IERC20(localToken).safeTransferFrom(msg.sender, address(this), amount);

    if (msg.value != 0) {
      revert MsgShouldNotContainValue(msg.value);
    }

    // Extra data for the L2 deposit.
    // We encode the nonce in the extra data so that we can track the L2 deposit offchain.
    bytes memory extraData = abi.encode(s_nonce++);

    // If the token is the wrapped native, we unwrap it and deposit native
    if (localToken == address(i_wrappedNative)) {
      i_wrappedNative.withdraw(amount);
      i_L1Bridge.depositETHTo{value: amount}(recipient, 0, extraData);
      return extraData;
    }

    // Token is a normal ERC20.
    IERC20(localToken).safeApprove(address(i_L1Bridge), amount);
    i_L1Bridge.depositERC20To(localToken, remoteToken, recipient, amount, 0, extraData);

    return extraData;
  }

  /// @notice Bridging to Optimism is paid for with gas
  /// @dev Since the gas amount charged is dynamic, the gas burn can change from block to block.
  /// You should always add a buffer of at least 20% to the gas limit for your L1 to L2 transaction
  /// to avoid running out of gas.
  function getBridgeFeeInNative() public pure returns (uint256) {
    return 0;
  }

  /// @notice Prove or finalize an ERC20 withdrawal from L2.
  /// The action to take is specified in the payload. See the docstring of FinalizeWithdrawERC20Payload for more details.
  /// @param data The payload for the action. This is an abi.encode'd FinalizeWithdrawERC20Payload with the appropriate data.
  /// @return true iff finalization is successful, and false for proving a withdrawal. If either of these fail,
  /// the call to this function will revert.
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address /* localReceiver */,
    bytes calldata data
  ) external override returns (bool) {
    // decode the data into FinalizeWithdrawERC20Payload first and extract the action.
    FinalizeWithdrawERC20Payload memory payload = abi.decode(data, (FinalizeWithdrawERC20Payload));
    if (payload.action == FinalizationAction.ProveWithdrawal) {
      // The action being ProveWithdrawal indicates that this is a withdrawal proof payload.
      // Decode the data into OptimismProveWithdrawalPayload and call the proveWithdrawal function.
      OptimismProveWithdrawalPayload memory provePayload = abi.decode(payload.data, (OptimismProveWithdrawalPayload));
      _proveWithdrawal(provePayload);
      return false;
    } else if (payload.action == FinalizationAction.FinalizeWithdrawal) {
      // decode the data into OptimismFinalizationPayload and call the finalizeWithdrawal function.
      OptimismFinalizationPayload memory finalizePayload = abi.decode(payload.data, (OptimismFinalizationPayload));
      // NOTE: finalizing ether withdrawals will currently send ether to the receiver address as indicated by the
      // withdrawal tx. However, this is problematic because we need to re-wrap it into WETH.
      // However, we can't do that from within this adapter because it doesn't actually have the ether.
      // So its up to the caller to rectify this by re-wrapping the ether.
      _finalizeWithdrawal(finalizePayload);
      return true;
    } else {
      revert InvalidFinalizationAction();
    }
  }

  function _proveWithdrawal(OptimismProveWithdrawalPayload memory payload) internal {
    // will revert if the proof is invalid or the output index is not yet included on L1.
    i_optimismPortal.proveWithdrawalTransaction(
      payload.withdrawalTransaction,
      payload.l2OutputIndex,
      payload.outputRootProof,
      payload.withdrawalProof
    );
  }

  function _finalizeWithdrawal(OptimismFinalizationPayload memory payload) internal {
    i_optimismPortal.finalizeWithdrawalTransaction(payload.withdrawalTransaction);
  }

  /// @notice returns the address of the WETH token used by this adapter.
  /// @return the address of the WETH token used by this adapter.
  function getWrappedNative() external view returns (address) {
    return address(i_wrappedNative);
  }

  /// @notice returns the address of the Optimism portal contract.
  /// @return the address of the Optimism portal contract.
  function getOptimismPortal() external view returns (address) {
    return address(i_optimismPortal);
  }

  /// @notice returns the address of the Optimism L1StandardBridge bridge contract.
  /// @return the address of the Optimism L1StandardBridge bridge contract.
  function getL1Bridge() external view returns (address) {
    return address(i_L1Bridge);
  }
}
