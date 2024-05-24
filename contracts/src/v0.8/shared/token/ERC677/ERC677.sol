// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {IERC677} from "./IERC677.sol";
import {IERC677Receiver} from "../../interfaces/IERC677Receiver.sol";

import {ERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";

contract ERC677 is IERC677, ERC20 {
  constructor(string memory name, string memory symbol) ERC20(name, symbol) {}

  /// @inheritdoc IERC677
  function transferAndCall(address to, uint256 amount, bytes memory data) public returns (bool success) {
    super.transfer(to, amount);
    emit Transfer(msg.sender, to, amount, data);
    if (to.code.length > 0) {
      IERC677Receiver(to).onTokenTransfer(msg.sender, amount, data);
    }
    return true;
  }
}
