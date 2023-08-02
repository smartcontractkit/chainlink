// SPDX-License-Identifier: MIT

pragma solidity ^0.8.19;

import {IConfigurable} from "./interfaces/IConfigurable.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IOwnableFunctionsRouter} from "./interfaces/IOwnableFunctionsRouter.sol";

abstract contract HasRouter is ITypeAndVersion, IConfigurable {
  IOwnableFunctionsRouter private immutable s_router;

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
    s_router = IOwnableFunctionsRouter(router);
    _updateConfig(config);
  }

  function _getRouter() internal view returns (IOwnableFunctionsRouter router) {
    return s_router;
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
    if (msg.sender != s_router.owner()) {
      revert OnlyCallableByRouterOwner();
    }
    _;
  }
}
