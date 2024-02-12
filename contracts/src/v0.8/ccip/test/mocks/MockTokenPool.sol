// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../../interfaces/pools/IPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract MockTokenPool is IPool {
  address public immutable i_token;

  constructor(address token) {
    i_token = token;
  }

  function lockOrBurn(
    address,
    bytes calldata,
    uint256,
    uint64,
    bytes calldata
  ) external pure override returns (bytes memory) {
    return "";
  }

  function releaseOrMint(bytes memory, address, uint256, uint64, bytes memory) external override {}

  function getToken() public view override returns (IERC20 token) {
    return IERC20(i_token);
  }
}
