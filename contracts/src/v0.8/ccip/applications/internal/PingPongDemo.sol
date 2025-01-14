// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../libraries/Client.sol";
import {CCIPClient} from "../external/CCIPClient.sol";

import {IRouter} from "../../interfaces/IRouter.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/// @title PingPongDemo - A simple ping-pong contract for demonstrating cross-chain communication
contract PingPongDemo is CCIPClient {
  event Ping(uint256 pingPongCount);
  event Pong(uint256 pingPongCount);
  event OutOfOrderExecutionChange(bool isOutOfOrder);

  // The chain ID of the counterpart ping pong contract
  uint64 internal s_counterpartChainSelector;

  // The contract address of the counterpart ping pong contract
  address internal s_counterpartAddress;

  // Pause ping-ponging
  bool private s_isPaused;

  bool private s_allowOutOfOrderExecution;

  // CCIPClient will handle the token approval so there's no need to do it here
  constructor(address router, IERC20 feeToken) CCIPClient(router, feeToken, true) {}

  function typeAndVersion() external pure virtual returns (string memory) {
    return "PingPongDemo 1.6.0-dev";
  }

  function startPingPong() external onlyOwner {
    s_isPaused = false;

    // Start the game
    _respond(1);
  }

  function _respond(uint256 pingPongCount) internal virtual {
    if (pingPongCount & 1 == 1) {
      emit Ping(pingPongCount);
    } else {
      emit Pong(pingPongCount);
    }

    bytes memory data = abi.encode(pingPongCount);

    ccipSend({destChainSelector: s_counterpartChainSelector, tokenAmounts: new Client.EVMTokenAmount[](0), data: data});
  }

  /// @notice This function the entrypoint for this contract to process messages.
  /// @param message The message to process.
  /// @dev This example just sends the tokens to the owner of this contracts. More
  /// interesting functions could be implemented.
  /// @dev It has to be external because of the try/catch.
  function processMessage(
    Client.Any2EVMMessage calldata message
  ) external override onlySelf isValidSender(message.sourceChainSelector, message.sender) {
    if (!s_isPaused) {
      _respond(abi.decode(message.data, (uint256)) + 1);
    }
  }

  // ================================================================
  // │                     Admin Functions                          │
  // ================================================================

  function setCounterpart(uint64 counterpartChainSelector, address counterpartAddress) external onlyOwner {
    if (counterpartAddress == address(0) || counterpartChainSelector == 0) revert ZeroAddressNotAllowed();

    s_counterpartChainSelector = counterpartChainSelector;
    s_counterpartAddress = counterpartAddress;

    // Approve the counterpart contract under validSender
    s_chainConfigs[counterpartChainSelector].approvedSender[abi.encode(counterpartAddress)] = true;

    // Approve the counterpart Chain selector under validChain
    s_chainConfigs[counterpartChainSelector].recipient = abi.encode(counterpartAddress);
  }

  function setPaused(bool pause) external onlyOwner {
    s_isPaused = pause;
  }

  // ================================================================
  // │                      State Management                        │
  // ================================================================

  function getCounterpartChainSelector() external view returns (uint64) {
    return s_counterpartChainSelector;
  }

  function getCounterpartAddress() external view returns (address) {
    return s_counterpartAddress;
  }

  function isPaused() external view returns (bool) {
    return s_isPaused;
  }

  function getOutOfOrderExecution() external view virtual returns (bool) {
    return s_allowOutOfOrderExecution;
  }

  function setOutOfOrderExecution(bool outOfOrderExecution) external virtual onlyOwner {
    // An additional storage slot is used for code simplicity. The current storage value can be
    // retrieved by parsing the extraArgsBytes field of the chain configuration, but this is not recommended
    // as it is more expensive and error-prone by requiring additional parsing logic in raw assembly.
    s_allowOutOfOrderExecution = outOfOrderExecution;

    address onRamp = IRouter(s_ccipRouter).getOnRamp(s_counterpartChainSelector);
    // OnRamp.StaticConfig memory staticConfig = OnRamp(onRamp).getStaticConfig();

    address feeQuoter = OnRamp(onRamp).getDynamicConfig().feeQuoter;
    uint64 defaultTxGasLimit = FeeQuoter(feeQuoter).getDestChainConfig(s_counterpartChainSelector).defaultTxGasLimit;

    // Enabling out of order execution also requires setting a manual gas limit, therefore the on-ramp default
    // gas limit is used to ensure consistency, but can be overwritten manually by the contract owner using
    // the applyChainUpdates function.
    s_chainConfigs[s_counterpartChainSelector].extraArgsBytes = Client._argsToBytes(
      Client.EVMExtraArgsV2({gasLimit: defaultTxGasLimit, allowOutOfOrderExecution: outOfOrderExecution})
    );

    emit OutOfOrderExecutionChange(outOfOrderExecution);
  }
}
