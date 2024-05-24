// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IOwnableFunctionsRouter} from "./interfaces/IOwnableFunctionsRouter.sol";

/// @title This abstract should be inherited by contracts that will be used
/// as the destinations to a route (id=>contract) on the Router.
/// It provides a Router getter and modifiers.
abstract contract Routable is ITypeAndVersion {
  IOwnableFunctionsRouter private immutable i_functionsRouter;

  error RouterMustBeSet();
  error OnlyCallableByRouter();
  error OnlyCallableByRouterOwner();

  /// @dev Initializes the contract.
  constructor(address router) {
    if (router == address(0)) {
      revert RouterMustBeSet();
    }
    i_functionsRouter = IOwnableFunctionsRouter(router);
  }

  /// @notice Return the Router
  function _getRouter() internal view returns (IOwnableFunctionsRouter router) {
    return i_functionsRouter;
  }

  /// @notice Reverts if called by anyone other than the router.
  modifier onlyRouter() {
    if (msg.sender != address(i_functionsRouter)) {
      revert OnlyCallableByRouter();
    }
    _;
  }

  /// @notice Reverts if called by anyone other than the router owner.
  modifier onlyRouterOwner() {
    if (msg.sender != i_functionsRouter.owner()) {
      revert OnlyCallableByRouterOwner();
    }
    _;
  }
}
