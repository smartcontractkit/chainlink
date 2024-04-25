// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../../interfaces/IPool.sol";

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
    return bytes("");
  }

  function releaseOrMint(
    bytes memory,
    address,
    uint256,
    uint64,
    IPool.SourceTokenData memory,
    bytes memory
  ) external view override returns (address) {
    return address(i_token);
  }

  function getToken() public view override returns (IERC20 token) {
    return IERC20(i_token);
  }

  function getRemotePool(uint64) public pure override returns (bytes memory) {
    return abi.encode(address(1));
  }
}
