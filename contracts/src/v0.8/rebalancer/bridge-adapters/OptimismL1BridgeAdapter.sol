// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBridgeAdapter} from "../interfaces/IBridge.sol";
import {IWrappedNative} from "../../ccip/interfaces/IWrappedNative.sol";

import {IL1StandardBridge} from "@eth-optimism/contracts/L1/messaging/IL1StandardBridge.sol";
import {IL1CrossDomainMessenger} from "@eth-optimism/contracts/L1/messaging/IL1CrossDomainMessenger.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract OptimismL1BridgeAdapter is IBridgeAdapter {
  using SafeERC20 for IERC20;

  IL1StandardBridge internal immutable i_L1Bridge;
  IL1CrossDomainMessenger internal immutable i_L1CrossDomainMessenger;
  IWrappedNative internal immutable i_wrappedNative;

  // Nonce to use for L2 deposits to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  constructor(
    IL1StandardBridge l1Bridge,
    IWrappedNative wrappedNative,
    IL1CrossDomainMessenger l1CrossDomainMessenger
  ) {
    if (address(l1Bridge) == address(0) || address(wrappedNative) == address(0)) {
      revert BridgeAddressCannotBeZero();
    }
    i_L1Bridge = l1Bridge;
    i_L1CrossDomainMessenger = l1CrossDomainMessenger;
    i_wrappedNative = wrappedNative;
  }

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

    // If the token is the wrapped native, we unwrap it and deposit native
    if (localToken == address(i_wrappedNative)) {
      i_wrappedNative.withdraw(amount);
      _depositNativeToL2(recipient, amount);
      return "";
    }

    // Token is normal ERC20
    IERC20(localToken).approve(address(i_L1Bridge), amount);
    i_L1Bridge.depositERC20To(localToken, remoteToken, recipient, amount, 0, abi.encode(s_nonce++));
    return "";
  }

  /// @notice Bridging to Optimism is paid for with gas
  /// @dev Since the gas amount charged is dynamic, the gas burn can change from block to block.
  /// You should always add a buffer of at least 20% to the gas limit for your L1 to L2 transaction
  /// to avoid running out of gas.
  function getBridgeFeeInNative() public pure returns (uint256) {
    return 0;
  }

  function depositNativeToL2(address recipient, uint256 amount) public payable {
    if (msg.value != amount) {
      revert MsgValueDoesNotMatchAmount(msg.value, amount);
    }

    _depositNativeToL2(recipient, amount);
  }

  function _depositNativeToL2(address recipient, uint256 amount) internal {
    i_L1Bridge.depositETHTo{value: amount}(recipient, 0, "");
  }

  struct OptimismFinalizationPayload {
    address l1Token;
    address l2Token;
    uint256 amount;
  }

  function finalizeWithdrawERC20(address remoteSender, address localReceiver, bytes calldata data) external override {
    OptimismFinalizationPayload memory payload = abi.decode(data, (OptimismFinalizationPayload));
    i_L1Bridge.finalizeERC20Withdrawal(
      payload.l1Token,
      payload.l2Token,
      remoteSender,
      localReceiver,
      payload.amount,
      data
    );
  }

  function finalizeWithdrawNativeFromL2(address from, address to, uint256 amount, bytes calldata data) external {
    i_L1Bridge.finalizeETHWithdrawal(from, to, amount, data);
  }

  function relayMessageFromL2ToL1(
    address target,
    address sender,
    bytes memory message,
    uint256 messageNonce,
    IL1CrossDomainMessenger.L2MessageInclusionProof memory proof
  ) external {
    i_L1CrossDomainMessenger.relayMessage(target, sender, message, messageNonce, proof);
    // TODO
  }

  /// @notice returns the address of the
  function getWrappedNative() external view returns (address) {
    return address(i_wrappedNative);
  }
}
