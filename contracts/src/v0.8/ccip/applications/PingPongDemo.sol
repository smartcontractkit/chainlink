// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IRouterClient} from "../interfaces/IRouterClient.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Client} from "../libraries/Client.sol";
import {CCIPReceiver} from "./CCIPReceiver.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/// @title PingPongDemo - A simple ping-pong contract for demonstrating cross-chain communication
contract PingPongDemo is CCIPReceiver, OwnerIsCreator, ITypeAndVersion {
  event Ping(uint256 pingPongCount);
  event Pong(uint256 pingPongCount);
  event OutOfOrderExecutionChange(bool isOutOfOrder);

  // Default gas limit used for EVMExtraArgsV2 construction
  uint64 private constant DEFAULT_GAS_LIMIT = 200_000;

  // The chain ID of the counterpart ping pong contract
  uint64 internal s_counterpartChainSelector;
  // The contract address of the counterpart ping pong contract
  address internal s_counterpartAddress;
  // Pause ping-ponging
  bool private s_isPaused;
  // The fee token used to pay for CCIP transactions
  IERC20 internal s_feeToken;
  // Allowing out of order execution
  bool private s_outOfOrderExecution;

  constructor(address router, IERC20 feeToken) CCIPReceiver(router) {
    s_isPaused = false;
    s_feeToken = feeToken;
    s_feeToken.approve(address(router), type(uint256).max);
  }

  function typeAndVersion() external pure virtual returns (string memory) {
    return "PingPongDemo 1.5.0";
  }

  function setCounterpart(uint64 counterpartChainSelector, address counterpartAddress) external onlyOwner {
    s_counterpartChainSelector = counterpartChainSelector;
    s_counterpartAddress = counterpartAddress;
  }

  function startPingPong() external onlyOwner {
    s_isPaused = false;
    _respond(1);
  }

  function _respond(
    uint256 pingPongCount
  ) internal virtual {
    if (pingPongCount & 1 == 1) {
      emit Ping(pingPongCount);
    } else {
      emit Pong(pingPongCount);
    }
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_counterpartAddress),
      data: abi.encode(pingPongCount),
      tokenAmounts: new Client.EVMTokenAmount[](0),
      extraArgs: Client._argsToBytes(
        Client.EVMExtraArgsV2({gasLimit: uint256(DEFAULT_GAS_LIMIT), allowOutOfOrderExecution: s_outOfOrderExecution})
      ),
      feeToken: address(s_feeToken)
    });
    IRouterClient(getRouter()).ccipSend(s_counterpartChainSelector, message);
  }

  function _ccipReceive(
    Client.Any2EVMMessage memory message
  ) internal override {
    uint256 pingPongCount = abi.decode(message.data, (uint256));
    if (!s_isPaused) {
      _respond(pingPongCount + 1);
    }
  }

  /////////////////////////////////////////////////////////////////////
  // Plumbing
  /////////////////////////////////////////////////////////////////////

  function getCounterpartChainSelector() external view returns (uint64) {
    return s_counterpartChainSelector;
  }

  function setCounterpartChainSelector(
    uint64 chainSelector
  ) external onlyOwner {
    s_counterpartChainSelector = chainSelector;
  }

  function getCounterpartAddress() external view returns (address) {
    return s_counterpartAddress;
  }

  function getFeeToken() external view returns (IERC20) {
    return s_feeToken;
  }

  function setCounterpartAddress(
    address addr
  ) external onlyOwner {
    s_counterpartAddress = addr;
  }

  function isPaused() external view returns (bool) {
    return s_isPaused;
  }

  function setPaused(
    bool pause
  ) external onlyOwner {
    s_isPaused = pause;
  }

  function getOutOfOrderExecution() external view returns (bool) {
    return s_outOfOrderExecution;
  }

  function setOutOfOrderExecution(
    bool outOfOrderExecution
  ) external onlyOwner {
    s_outOfOrderExecution = outOfOrderExecution;
    emit OutOfOrderExecutionChange(outOfOrderExecution);
  }
}
