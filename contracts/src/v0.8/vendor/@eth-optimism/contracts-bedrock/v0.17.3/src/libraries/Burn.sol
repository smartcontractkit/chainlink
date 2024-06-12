// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

/// @title Burn
/// @notice Utilities for burning stuff.
library Burn {
  /// @notice Burns a given amount of ETH.
  /// @param _amount Amount of ETH to burn.
  function eth(uint256 _amount) internal {
    new Burner{value: _amount}();
  }

  /// @notice Burns a given amount of gas.
  /// @param _amount Amount of gas to burn.
  function gas(uint256 _amount) internal view {
    uint256 i = 0;
    uint256 initialGas = gasleft();
    while (initialGas - gasleft() < _amount) {
      ++i;
    }
  }
}

/// @title Burner
/// @notice Burner self-destructs on creation and sends all ETH to itself, removing all ETH given to
///         the contract from the circulating supply. Self-destructing is the only way to remove ETH
///         from the circulating supply.
contract Burner {
  constructor() payable {
    // solhint-disable-next-line avoid-low-level-calls
    selfdestruct(payable(address(this)));
  }
}
