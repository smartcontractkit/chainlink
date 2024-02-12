// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {FacadeClient} from "./FacadeClient.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract ReentrantMaliciousTokenPool is TokenPool {
  address private i_facade;

  bool private s_attacked;

  constructor(
    address facade,
    IERC20 token,
    address armProxy,
    address router
  ) TokenPool(token, new address[](0), armProxy, router) {
    i_facade = facade;
  }

  /// @dev Calls into Facade to reenter Router exactly 1 time
  function lockOrBurn(
    address,
    bytes calldata,
    uint256 amount,
    uint64,
    bytes calldata
  ) external override returns (bytes memory) {
    if (s_attacked) {
      return "";
    }

    s_attacked = true;

    FacadeClient(i_facade).send(amount);
    emit Burned(msg.sender, amount);
    return "";
  }

  function releaseOrMint(bytes memory, address receiver, uint256 amount, uint64, bytes memory) external override {}
}
