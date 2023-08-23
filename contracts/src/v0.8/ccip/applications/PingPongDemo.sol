// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../interfaces/IRouterClient.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Client} from "../libraries/Client.sol";
import {CCIPReceiver} from "./CCIPReceiver.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/contracts/token/ERC20/IERC20.sol";

/// @title PingPongDemo - A simple ping-pong contract for demonstrating cross-chain communication
contract PingPongDemo is CCIPReceiver, OwnerIsCreator {
  event Ping(uint256 pingPongCount);
  event Pong(uint256 pingPongCount);

  // The chain ID of the counterpart ping pong contract
  uint64 private s_counterpartChainSelector;
  // The contract address of the counterpart ping pong contract
  address private s_counterpartAddress;

  // Pause ping-ponging
  bool private s_isPaused;
  IERC20 private s_feeToken;

  constructor(address router, IERC20 feeToken) CCIPReceiver(router) {
    s_isPaused = false;
    s_feeToken = feeToken;
    s_feeToken.approve(address(router), 2 ** 256 - 1);
  }

  function setCounterpart(uint64 counterpartChainSelector, address counterpartAddress) external onlyOwner {
    s_counterpartChainSelector = counterpartChainSelector;
    s_counterpartAddress = counterpartAddress;
  }

  function startPingPong() external onlyOwner {
    s_isPaused = false;
    _respond(1);
  }

  function _respond(uint256 pingPongCount) private {
    if (pingPongCount & 1 == 1) {
      emit Ping(pingPongCount);
    } else {
      emit Pong(pingPongCount);
    }

    bytes memory data = abi.encode(pingPongCount);
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_counterpartAddress),
      data: data,
      tokenAmounts: new Client.EVMTokenAmount[](0),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000})),
      feeToken: address(s_feeToken)
    });
    IRouterClient(getRouter()).ccipSend(s_counterpartChainSelector, message);
  }

  function _ccipReceive(Client.Any2EVMMessage memory message) internal override {
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

  function setCounterpartChainSelector(uint64 chainSelector) external onlyOwner {
    s_counterpartChainSelector = chainSelector;
  }

  function getCounterpartAddress() external view returns (address) {
    return s_counterpartAddress;
  }

  function setCounterpartAddress(address addr) external onlyOwner {
    s_counterpartAddress = addr;
  }

  function isPaused() external view returns (bool) {
    return s_isPaused;
  }

  function setPaused(bool pause) external onlyOwner {
    s_isPaused = pause;
  }
}
