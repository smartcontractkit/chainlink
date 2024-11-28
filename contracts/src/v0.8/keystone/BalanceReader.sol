// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

/// @notice BalanceReader is used to read native currency balances from one or more accounts
/// using a contract method instead of an RPC "eth_getBalance" call.
contract BalanceReader {
  function getNativeBalances(address[] memory addresses) public view returns (uint256[] memory) {
    uint256[] memory balances = new uint256[](addresses.length);
    for (uint256 i = 0; i < addresses.length; ++i) {
      balances[i] = addresses[i].balance;
    }
    return balances;
  }
}
