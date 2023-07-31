// SPDX-License-Identifier: MIT

pragma solidity ^0.8.19;

import {IConfigurable} from "./interfaces/IConfigurable.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IOwnableRouter} from "./interfaces/IOwnableRouter.sol";

abstract contract HasRouter is ITypeAndVersion, IConfigurable {
  bytes32 internal s_configHash;

  IOwnableRouter private s_router;

  error RouterMustBeSet();
  error OnlyCallableByRouter();
  error OnlyCallableByRouterOwner();

  /**
   * @dev Initializes the contract.
   */
  constructor(address router, bytes memory config) {
    if (router == address(0)) {
      revert RouterMustBeSet();
    }
    s_router = IOwnableRouter(router);
    updateConfig(config);
  }

  function _getRouter() internal view returns (IOwnableRouter router) {
    return s_router;
  }

  /**
   * @inheritdoc IConfigurable
   */
  function getConfigHash() external view override returns (bytes32 config) {
    return s_configHash;
  }

  /**
   * @dev Must be implemented by inheriting contract
   * Use to set configuration state
   */
  function _updateConfig(bytes memory config) internal virtual;

  /**
   * @inheritdoc IConfigurable
   * @dev Only callable by the Router
   */
  function updateConfig(bytes memory config) public override onlyRouter {
    _updateConfig(config);
    s_configHash = keccak256(config);
  }

  /**
   * @notice Reverts if called by anyone other than the router.
   */
  modifier onlyRouter() {
    if (msg.sender != address(s_router)) {
      revert OnlyCallableByRouter();
    }
    _;
  }

  /**
   * @notice Reverts if called by anyone other than the router owner.
   */
  modifier onlyRouterOwner() {
    if (msg.sender != IOwnableRouter(address(s_router)).owner()) {
      revert OnlyCallableByRouterOwner();
    }
    _;
  }
}
