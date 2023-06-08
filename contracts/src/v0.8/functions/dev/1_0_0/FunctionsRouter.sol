// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {RouterBase, ITypeAndVersion} from "./RouterBase.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {IVersioned} from "./interfaces/IVersioned.sol";
import {IFunctionsCoordinator} from "./interfaces/IFunctionsCoordinator.sol";
import {AuthorizedOriginReceiver} from "./accessControl/AuthorizedOriginReceiver.sol";
import {IFunctionsSubscriptions, FunctionsSubscriptions, IFunctionsBilling} from "./FunctionsSubscriptions.sol";

contract FunctionsRouter is RouterBase, IFunctionsRouter, AuthorizedOriginReceiver, FunctionsSubscriptions {
  // ================================================================
  // |                    Configuration state                       |
  // ================================================================
  struct Config {
    // Flat fee (in Juels of LINK) that will be paid to the Router owner for operation of the network
    uint96 adminFee;
  }
  Config private s_config;
  event ConfigSet(uint96 adminFee);

  error OnlyCallableByRoute();

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(
    uint16 timelockBlocks,
    uint16 maximumTimelockBlocks,
    bool useAllowList,
    address linkToken,
    bytes32[] memory initialJobIds,
    address[] memory initialAddresses,
    bytes memory config
  )
    RouterBase(msg.sender, timelockBlocks, maximumTimelockBlocks, initialJobIds, initialAddresses, config)
    AuthorizedOriginReceiver(useAllowList)
    FunctionsSubscriptions(linkToken)
  {}

  /**
   * @inheritdoc ITypeAndVersion
   */
  function typeAndVersion() public pure override returns (string memory) {
    return "Functions Router v1";
  }

  // ================================================================
  // |                 Configuration methods                        |
  // ================================================================
  /**
   * @notice Sets the configuration for FunctionsRouter specific state
   * @param config bytes of config data to set the following:
   *  - adminFee: fee that will be paid to the Router owner for operating the network
   */
  function _setConfig(bytes memory config) internal override {
    uint96 adminFee = abi.decode(config, (uint96));
    s_config = Config({adminFee: adminFee});
    emit ConfigSet(adminFee);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function getAdminFee() external view override returns (uint96) {
    return s_config.adminFee;
  }

  // ================================================================
  // |                      Request methods                         |
  // ================================================================

  function _sendRequest(
    bytes32 jobId,
    bool isProposed,
    uint64 subscriptionId,
    bytes memory data,
    uint32 gasLimit
  ) internal returns (bytes32) {
    _isValidSubscription(subscriptionId);
    _isValidConsumer(msg.sender, subscriptionId);

    address route = this.getRoute(jobId, isProposed);
    IFunctionsCoordinator coordinator = IFunctionsCoordinator(route);

    (, , address owner, , ) = this.getSubscription(subscriptionId);

    (bytes32 requestId, uint96 estimatedCost, uint256 requestTimeoutSeconds) = coordinator.sendRequest(
      subscriptionId,
      data,
      gasLimit,
      msg.sender,
      owner
    );

    _blockBalance(msg.sender, subscriptionId, estimatedCost, requestId, route, requestTimeoutSeconds);

    return requestId;
  }

  function _smoke(bytes32 jobId, bytes calldata data) internal override onlyAuthorizedUsers returns (bytes32) {
    (uint64 subscriptionId, bytes memory reqData, uint32 gasLimit) = abi.decode(data, (uint64, bytes, uint32));
    return _sendRequest(jobId, true, subscriptionId, reqData, gasLimit);
  }

  /**
   * @inheritdoc IFunctionsRouter
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    bytes32 jobId
  ) external override onlyAuthorizedUsers returns (bytes32) {
    return _sendRequest(jobId, false, subscriptionId, data, gasLimit);
  }

  // ================================================================
  // |                    Modifier Overrides                        |
  // ================================================================

  function _canSetAuthorizedSenders() internal view override onlyOwner returns (bool) {
    return msg.sender == owner();
  }

  modifier onlyAuthorizedUsers() override {
    _validateIsAuthorizedSender();
    _;
  }

  modifier nonReentrant() override {
    // TODO
    // address route = this.getRoute("FunctionsCoordinator", true);
    // IFunctionsBilling coordinator = IFunctionsBilling(route);
    // if (coordinator.isReentrancyLocked()) {
    //   revert Reentrant();
    // }
    _;
  }

  modifier onlyRouterOwner() override {
    _validateOwnership();
    _;
  }
}
