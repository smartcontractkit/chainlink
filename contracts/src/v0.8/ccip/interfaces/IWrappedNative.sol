// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/token/ERC20/IERC20.sol";

interface IWrappedNative is IERC20 {
  function deposit() external payable;
}
