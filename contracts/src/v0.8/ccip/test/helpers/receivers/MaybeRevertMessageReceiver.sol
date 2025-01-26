// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";
import {Client} from "../../../libraries/Client.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";
import {IERC20} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/IERC20.sol";

contract MaybeRevertMessageReceiver is IAny2EVMMessageReceiver, IERC165 {
  error ReceiveRevert();
  error CustomError(bytes err);
  error Unauthorized();
  error InsufficientBalance(uint256 available, uint256 required);
  error TransferFailed();

  event ValueReceived(uint256 amount);
  event FundsWithdrawn(address indexed owner, uint256 amount);
  event TokensWithdrawn(address indexed token, address indexed owner, uint256 amount);
  event MessageReceived(
    bytes32 messageId,
    uint64 sourceChainSelector,
    bytes sender,
    bytes data,
    Client.EVMTokenAmount[] destTokenAmounts
  );

  address private immutable s_manager;
  bool public s_toRevert;
  bytes private s_err;

  constructor(bool toRevert) {
    s_manager = msg.sender;
    s_toRevert = toRevert;
  }

  modifier onlyManager() {
    if (msg.sender != s_manager) {
      revert Unauthorized();
    }
    _;
  }

  function setRevert(bool toRevert) external onlyManager {
    s_toRevert = toRevert;
  }

  function setErr(bytes memory err) external onlyManager {
    s_err = err;
  }

  /// @notice IERC165 supports an interfaceId
  /// @param interfaceId The interfaceId to check
  /// @return true if the interfaceId is supported
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IAny2EVMMessageReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }

  function ccipReceive(Client.Any2EVMMessage calldata message) external override {
    if (s_toRevert) {
      revert CustomError(s_err);
    }

    emit MessageReceived(
      message.messageId,
      message.sourceChainSelector,
      message.sender,
      message.data,
      message.destTokenAmounts
    );
  }

  receive() external payable {
    if (s_toRevert) {
      revert ReceiveRevert();
    }

    emit ValueReceived(msg.value);
  }

  /// @notice Allows the manager (deployer) to withdraw all Ether from the contract
  function withdrawFunds() external onlyManager {
    uint256 balance = address(this).balance;
    if (balance == 0) {
      revert InsufficientBalance(0, 1);
    }

    (bool success, ) = s_manager.call{value: balance}("");
    if (!success) {
      revert TransferFailed();
    }

    emit FundsWithdrawn(s_manager, balance);
  }

  /// @notice Allows the manager to withdraw ERC-20 tokens from the contract
  /// @param token The address of the ERC-20 token contract
  /// @param amount The amount of tokens to withdraw
  function withdrawTokens(address token, uint256 amount) external onlyManager {
    IERC20 erc20 = IERC20(token);
    uint256 balance = erc20.balanceOf(address(this));
    if (balance < amount) {
      revert InsufficientBalance(balance, amount);
    }

    bool success = erc20.transfer(s_manager, amount);
    if (!success) {
      revert TransferFailed();
    }

    emit TokensWithdrawn(token, s_manager, amount);
  }

  /// @notice Fetches the balance of an ERC-20 token held by the contract
  /// @param token The address of the ERC-20 token contract
  /// @return The balance of the specified ERC-20 token
  function balanceOfToken(address token) external view returns (uint256) {
    IERC20 erc20 = IERC20(token);
    return erc20.balanceOf(address(this));
  }
}
