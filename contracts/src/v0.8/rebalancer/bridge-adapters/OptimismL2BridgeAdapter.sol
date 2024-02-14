// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBridgeAdapter} from "../interfaces/IBridge.sol";
import {IWrappedNative} from "../../ccip/interfaces/IWrappedNative.sol";

import {L2StandardBridge} from "@eth-optimism/contracts/L2/messaging/L2StandardBridge.sol";
import {Lib_PredeployAddresses} from "@eth-optimism/contracts/libraries/constants/Lib_PredeployAddresses.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract OptimismL2BridgeAdapter is IBridgeAdapter {
  using SafeERC20 for IERC20;

  L2StandardBridge internal immutable i_L2Bridge = L2StandardBridge(Lib_PredeployAddresses.L2_STANDARD_BRIDGE);
  IWrappedNative internal immutable i_wrappedNative;

  // Nonce to use for L1 withdrawals to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  constructor(IWrappedNative wrappedNative) {
    // Wrapped native can be address zero, this means that auto-wrapping is disabled.
    i_wrappedNative = wrappedNative;
  }

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

    // If the token is the wrapped native, we unwrap it and deposit native
    if (localToken == address(i_wrappedNative)) {
      i_wrappedNative.withdraw(amount);
      _depositNativeToL1(recipient, amount);
      return "";
    }

    // Token is normal ERC20
    IERC20(localToken).approve(address(i_L2Bridge), amount);
    i_L2Bridge.withdrawTo(localToken, recipient, amount, 0, abi.encode(s_nonce++));
    return "";
  }

  /// @notice No-op since L1 -> L2 transfers do not need finalization.
  function finalizeWithdrawERC20(
    address /* remoteSender */,
    address /* localReceiver */,
    bytes calldata /* bridgeSpecificPayload */
  ) external override {}

  /// @notice There are no fees to bridge back to L1
  function getBridgeFeeInNative() external pure returns (uint256) {
    return 0;
  }

  function depositNativeToL1(address recipient) public payable {
    _depositNativeToL1(recipient, msg.value);
  }

  function _depositNativeToL1(address recipient, uint256 amount) internal {
    i_L2Bridge.withdrawTo(Lib_PredeployAddresses.OVM_ETH, recipient, amount, 0, abi.encode(s_nonce++));
  }

  /// @notice returns the address of the
  function getWrappedNative() external view returns (address) {
    return address(i_wrappedNative);
  }
}
