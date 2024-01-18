// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IL1BridgeAdapter} from "../interfaces/IBridge.sol";
import {IWrappedNative} from "../../ccip/interfaces/IWrappedNative.sol";

import {IL1StandardBridge} from "@eth-optimism/contracts/L1/messaging/IL1StandardBridge.sol";
import {IL1CrossDomainMessenger} from "@eth-optimism/contracts/L1/messaging/IL1CrossDomainMessenger.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract OptimismL1BridgeAdapter is IL1BridgeAdapter {
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

  function sendERC20(address l1Token, address l2Token, address recipient, uint256 amount) external payable {
    IERC20(l1Token).safeTransferFrom(msg.sender, address(this), amount);

    if (msg.value != 0) {
      revert MsgShouldNotContainValue(msg.value);
    }

    // If the token is the wrapped native, we unwrap it and deposit native
    if (l1Token == address(i_wrappedNative)) {
      i_wrappedNative.withdraw(amount);
      depositNativeToL2(recipient, amount);
      return;
    }

    // Token is normal ERC20
    IERC20(l1Token).approve(address(i_L1Bridge), amount);
    i_L1Bridge.depositERC20To(l1Token, l2Token, recipient, amount, 0, abi.encode(s_nonce++));
  }

  /// @notice Bridging to Optimism is free.
  function getBridgeFeeInNative() public pure returns (uint256) {
    return 0;
  }

  function depositNativeToL2(address recipient, uint256 amount) public payable {
    if (msg.value != amount) {
      revert MsgValueDoesNotMatchAmount(msg.value, amount);
    }

    i_L1Bridge.depositETHTo{value: msg.value}(recipient, 0, abi.encode(s_nonce++));
  }

  struct OptimismFinalizationPayload {
    address l1Token;
    address l2Token;
    uint256 amount;
  }

  function finalizeWithdrawERC20FromL2(address from, address to, bytes calldata data) external {
    OptimismFinalizationPayload memory payload = abi.decode(data, (OptimismFinalizationPayload));
    i_L1Bridge.finalizeERC20Withdrawal(payload.l1Token, payload.l2Token, from, to, payload.amount, data);
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
