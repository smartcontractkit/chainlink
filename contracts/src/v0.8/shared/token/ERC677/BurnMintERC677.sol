// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IERC677Receiver} from "../../interfaces/IERC677Receiver.sol";
import {IERC677} from "./IERC677.sol";

import {BurnMintERC20} from "../ERC20/BurnMintERC20.sol";

/// @notice A basic ERC677 compatible token contract with burn and minting roles.
/// @dev The total supply can be limited during deployment.
contract BurnMintERC677 is BurnMintERC20, IERC677 {
  constructor(
    string memory name,
    string memory symbol,
    uint8 decimals,
    uint256 maxSupply
  ) BurnMintERC20(name, symbol, decimals, maxSupply) {}

  /// @inheritdoc IERC677
  function transferAndCall(address to, uint256 amount, bytes memory data) public returns (bool success) {
    transfer(to, amount);
    emit Transfer(msg.sender, to, amount, data);
    if (to.code.length > 0) {
      IERC677Receiver(to).onTokenTransfer(msg.sender, amount, data);
    }
    return true;
  }

  /// @inheritdoc BurnMintERC20
  function supportsInterface(bytes4 interfaceId) public pure virtual override returns (bool) {
    return interfaceId == type(IERC677).interfaceId || super.supportsInterface(interfaceId);
  }
}
