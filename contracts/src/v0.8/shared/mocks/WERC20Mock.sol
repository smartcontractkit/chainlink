// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";

contract WERC20Mock is ERC20 {
  constructor() ERC20("WERC20Mock", "WERC") {}

  event Deposit(address indexed dst, uint256 wad);
  event Withdrawal(address indexed src, uint256 wad);

  receive() external payable {
    deposit();
  }

  function deposit() public payable {
    _mint(msg.sender, msg.value);
    emit Deposit(msg.sender, msg.value);
  }

  function withdraw(uint256 wad) public {
    // solhint-disable-next-line gas-custom-errors, reason-string
    require(balanceOf(msg.sender) >= wad);
    _burn(msg.sender, wad);
    payable(msg.sender).transfer(wad);
    emit Withdrawal(msg.sender, wad);
  }

  function mint(address account, uint256 amount) external {
    _mint(account, amount);
  }

  function burn(address account, uint256 amount) external {
    _burn(account, amount);
  }
}
