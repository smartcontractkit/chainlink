// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBridgeAdapter} from "../interfaces/IBridge.sol";
import {IWrappedNative} from "../../ccip/interfaces/IWrappedNative.sol";

import {Lib_PredeployAddresses} from "@eth-optimism/contracts/libraries/constants/Lib_PredeployAddresses.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @dev copy/pasted from https://github.com/ethereum-optimism/optimism/blob/f707883038d527cbf1e9f8ea513fe33255deadbc/packages/contracts-bedrock/src/L2/L2StandardBridge.sol#L114-L122.
/// We can't import it because of hard pin solidity version in the pragma (0.8.15).
interface IL2StandardBridge {
  /// @custom:legacy
  /// @notice Initiates a withdrawal from L2 to L1 to a target account on L1.
  ///         Note that if ETH is sent to a contract on L1 and the call fails, then that ETH will
  ///         be locked in the L1StandardBridge. ETH may be recoverable if the call can be
  ///         successfully replayed by increasing the amount of gas supplied to the call. If the
  ///         call will fail for any amount of gas, then the ETH will be locked permanently.
  ///         This function only works with OptimismMintableERC20 tokens or ether. Use the
  ///         `bridgeERC20To` function to bridge native L2 tokens to L1.
  /// @param _l2Token     Address of the L2 token to withdraw.
  /// @param _to          Recipient account on L1.
  /// @param _amount      Amount of the L2 token to withdraw.
  /// @param _minGasLimit Minimum gas limit to use for the transaction.
  /// @param _extraData   Extra data attached to the withdrawal.
  function withdrawTo(
    address _l2Token,
    address _to,
    uint256 _amount,
    uint32 _minGasLimit,
    bytes calldata _extraData
  ) external payable;
}

/// @notice OptimismL2BridgeAdapter implements IBridgeAdapter for the Optimism L2<=>L1 bridge.
/// @dev We have to unwrap WETH into ether before withdrawing it to L1. Therefore this bridge adapter bridges
/// WETH to ether. The receiver on L1 must wrap the ether back into WETH.
contract OptimismL2BridgeAdapter is IBridgeAdapter {
  using SafeERC20 for IERC20;

  IL2StandardBridge internal immutable i_L2Bridge = IL2StandardBridge(Lib_PredeployAddresses.L2_STANDARD_BRIDGE);
  IWrappedNative internal immutable i_wrappedNative;

  // Nonce to use for L1 withdrawals to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  constructor(IWrappedNative wrappedNative) {
    // Wrapped native can be address zero, this means that auto-wrapping is disabled.
    i_wrappedNative = wrappedNative;
  }

  /// @notice The WETH withdraw requires this be present otherwise withdraws will fail.
  receive() external payable {}

  /// @inheritdoc IBridgeAdapter
  function sendERC20(
    address localToken,
    address /* remoteToken */,
    address recipient,
    uint256 amount,
    bytes calldata /* bridgeSpecificPayload */
  ) external payable override returns (bytes memory) {
    if (msg.value != 0) {
      revert MsgShouldNotContainValue(msg.value);
    }

    IERC20(localToken).safeTransferFrom(msg.sender, address(this), amount);

    // Extra data for the L2 withdraw.
    // We encode the nonce in the extra data so that we can track the L2 withdraw offchain.
    bytes memory extraData = abi.encode(s_nonce++);

    // If the token is the wrapped native, we unwrap it and withdraw native
    if (localToken == address(i_wrappedNative)) {
      i_wrappedNative.withdraw(amount);
      // XXX: Lib_PredeployAddresses.OVM_ETH is actually 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000.
      // This code path still works because the L2 bridge is hardcoded to handle this specific address.
      // The better approach might be to use the bridgeEthTo function, which is on the StandardBridge
      // abstract contract, inherited by both L1StandardBridge and L2StandardBridge.
      // This is also marked as legacy, so it might mean that this will be deprecated soon.
      i_L2Bridge.withdrawTo{value: amount}(Lib_PredeployAddresses.OVM_ETH, recipient, amount, 0, extraData);
      return extraData;
    }

    // Token is normal ERC20
    IERC20(localToken).approve(address(i_L2Bridge), amount);
    i_L2Bridge.withdrawTo(localToken, recipient, amount, 0, extraData);
    return extraData;
  }

  /// @notice No-op since L1 -> L2 transfers do not need finalization.
  /// @return true always.
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address /* localReceiver */,
    bytes calldata /* bridgeSpecificPayload */
  ) external pure override returns (bool) {
    return true;
  }

  /// @notice There are no fees to bridge back to L1
  function getBridgeFeeInNative() external pure returns (uint256) {
    return 0;
  }

  /// @notice returns the address of the WETH token used by this adapter.
  /// @return the address of the WETH token used by this adapter.
  function getWrappedNative() external view returns (address) {
    return address(i_wrappedNative);
  }

  /// @notice returns the address of the L2 bridge used by this adapter.
  /// @return the address of the L2 bridge used by this adapter.
  function getL2Bridge() external view returns (address) {
    return address(i_L2Bridge);
  }
}
