// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

interface IWrappedNative is IERC20 {
  function deposit() external payable;

  function withdraw(uint256 wad) external;
}
