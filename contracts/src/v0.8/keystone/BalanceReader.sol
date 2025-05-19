// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";

/// @notice BalanceReader is used to read native currency balances from one or more accounts
/// using a contract method instead of an RPC "eth_getBalance" call.
contract BalanceReader is ITypeAndVersion {
  string public constant override typeAndVersion = "BalanceReader 1.0.0";

  function getNativeBalances(address[] memory addresses) public view returns (uint256[] memory) {
    uint256[] memory balances = new uint256[](addresses.length);
    for (uint256 i = 0; i < addresses.length; ++i) {
      balances[i] = addresses[i].balance;
    }
    return balances;
  }
}
