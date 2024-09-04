// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IBridgehub, L2TransactionRequestDirect, L2TransactionRequestTwoBridgesOuter} from "@zksync/contracts/l1-contracts/contracts/bridgehub/IBridgehub.sol";
import {IL1SharedBridge} from "@zksync/contracts/l1-contracts/contracts/bridge/interfaces/IL1SharedBridge.sol";
import {L2Message, L2Log, TxStatus} from "@zksync/contracts/l1-contracts/contracts/common/Messaging.sol";

contract MockBridgehub is IBridgehub {
  address public pendingAdmin;
  address public admin;
  address public sharedBridgeAddr;

  mapping(address stateTransitionManager => bool stateTransitionManagerIsRegistered)
    public registeredStateTransitionManagers;
  mapping(uint256 chainId => address stateTransitionManagerAddress) public stateTransitionManagers;
  mapping(address baseToken => bool tokenIsRegistered) public registeredTokens;
  mapping(uint256 chainId => address baseToken) public baseTokens;
  mapping(uint256 chainId => address hyperChain) public hyperchains;

  /// Generic error for unauthorized actions
  error NotAuthorized(string msg);

  /// Fake event that will get emitted when `requestL2TransactionDirect` is called
  event SentMessage(address indexed sender, bytes message);

  /// Admin functions
  function setPendingAdmin(address _newPendingAdmin) external override {
    emit NewPendingAdmin(pendingAdmin, _newPendingAdmin);
    pendingAdmin = _newPendingAdmin;
  }

  function acceptAdmin() external override {
    if (msg.sender != pendingAdmin) {
      revert NotAuthorized("Only pending admin can accept");
    }

    emit NewAdmin(admin, pendingAdmin);
    admin = pendingAdmin;
    pendingAdmin = address(0);
  }

  /// Getters
  function stateTransitionManagerIsRegistered(address _stateTransitionManager) external view override returns (bool) {
    return registeredStateTransitionManagers[_stateTransitionManager];
  }

  function stateTransitionManager(uint256 _chainId) external view override returns (address) {
    return stateTransitionManagers[_chainId];
  }

  function tokenIsRegistered(address _baseToken) external view override returns (bool) {
    return registeredTokens[_baseToken];
  }

  function baseToken(uint256 _chainId) external view override returns (address) {
    return baseTokens[_chainId];
  }

  function sharedBridge() external view override returns (IL1SharedBridge) {
    return IL1SharedBridge(sharedBridgeAddr);
  }

  function getHyperchain(uint256 _chainId) external view override returns (address) {
    return hyperchains[_chainId];
  }

  /// Mailbox forwarder
  function proveL2MessageInclusion(
    uint256,
    uint256,
    uint256,
    L2Message calldata,
    bytes32[] calldata
  ) external pure override returns (bool) {
    return true;
  }

  function proveL2LogInclusion(
    uint256,
    uint256,
    uint256,
    L2Log memory,
    bytes32[] calldata
  ) external pure override returns (bool) {
    return true;
  }

  function proveL1ToL2TransactionStatus(
    uint256,
    bytes32,
    uint256,
    uint256,
    uint16,
    bytes32[] calldata,
    TxStatus
  ) external pure override returns (bool) {
    return true;
  }

  function requestL2TransactionDirect(
    L2TransactionRequestDirect calldata txRequest
  ) external payable override returns (bytes32) {
    emit SentMessage(msg.sender, txRequest.l2Calldata);
    return keccak256(abi.encodePacked("L2TransactionDirect"));
  }

  function requestL2TransactionTwoBridges(
    L2TransactionRequestTwoBridgesOuter calldata
  ) external payable override returns (bytes32) {
    return keccak256(abi.encodePacked("L2TransactionTwoBridges"));
  }

  function l2TransactionBaseCost(uint256, uint256, uint256, uint256) external pure override returns (uint256) {
    return 0;
  }

  /// Registry
  function createNewChain(
    uint256 _chainId,
    address _stateTransitionManager,
    address _baseToken,
    uint256,
    address,
    bytes calldata
  ) external override returns (uint256 chainId) {
    hyperchains[_chainId] = _stateTransitionManager;
    baseTokens[_chainId] = _baseToken;
    emit NewChain(_chainId, _stateTransitionManager, address(this));
    return _chainId;
  }

  function addStateTransitionManager(address _stateTransitionManager) external override {
    registeredStateTransitionManagers[_stateTransitionManager] = true;
  }

  function removeStateTransitionManager(address _stateTransitionManager) external override {
    registeredStateTransitionManagers[_stateTransitionManager] = false;
  }

  function addToken(address _token) external override {
    registeredTokens[_token] = true;
  }

  function setSharedBridge(address _sharedBridgeAddr) external override {
    sharedBridgeAddr = _sharedBridgeAddr;
  }
}
