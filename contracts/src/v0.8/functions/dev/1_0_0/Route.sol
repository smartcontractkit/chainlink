// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import {IConfigurable} from "./interfaces/IConfigurable.sol";
import {Versioned} from "./Versioned.sol";
import {IRouterBase} from "./interfaces/IRouterBase.sol";
import {IOwnable} from "../../../shared/interfaces/IOwnable.sol";

abstract contract Route is Versioned, IConfigurable {
  bytes32 internal s_config_hash;

  IRouterBase internal s_router;

  error RouterMustBeSet();
  error OnlyCallableByRouter();
  error OnlyCallableByRouterOwner();

  /**
   * @dev Initializes the contract.
   */
  constructor(
    string memory id,
    uint16 version,
    address router,
    bytes memory config
  ) Versioned(id, version) {
    if (router == address(0)) {
      revert RouterMustBeSet();
    }
    s_router = IRouterBase(router);
    _setConfig(config);
    s_config_hash = keccak256(config);
  }

  /**
   * @inheritdoc IConfigurable
   */
  function getConfigHash() external view override returns (bytes32 config) {
    return s_config_hash;
  }

  function _setConfig(bytes memory config) internal virtual;

  /**
   * @inheritdoc IConfigurable
   */
  function setConfig(bytes memory config) external onlyRouter {
    _setConfig(config);
    s_config_hash = keccak256(config);
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
    if (msg.sender != IOwnable(address(s_router)).owner()) {
      revert OnlyCallableByRouterOwner();
    }
    _;
  }
}
