// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

interface IBridgeAdapter {
  error BridgeAddressCannotBeZero();
  error MsgValueDoesNotMatchAmount(uint256 msgValue, uint256 amount);
  error InsufficientEthValue(uint256 wanted, uint256 got);
  error MsgShouldNotContainValue(uint256 value);

  function sendERC20(address l1Token, address l2Token, address recipient, uint256 amount) external payable;

  function getBridgeFeeInNative() external view returns (uint256);
}

interface IL1BridgeAdapter is IBridgeAdapter {
  function finalizeWithdrawERC20FromL2(
    address l2Sender,
    address l1Receiver,
    bytes calldata bridgeSpecificPayload
  ) external;
}
