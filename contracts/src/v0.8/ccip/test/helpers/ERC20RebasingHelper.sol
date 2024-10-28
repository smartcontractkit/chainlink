// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";

contract ERC20RebasingHelper is ERC20 {
  uint16 public s_multiplierPercentage = 100;
  bool public s_mintShouldBurn = false;

  constructor() ERC20("Rebasing", "REB") {}

  function mint(address to, uint256 amount) external {
    if (!s_mintShouldBurn) {
      _mint(to, amount * s_multiplierPercentage / 100);
      return;
    }
    _burn(to, amount * s_multiplierPercentage / 100);
  }

  function setMultiplierPercentage(
    uint16 multiplierPercentage
  ) external {
    s_multiplierPercentage = multiplierPercentage;
  }

  function setMintShouldBurn(
    bool mintShouldBurn
  ) external {
    s_mintShouldBurn = mintShouldBurn;
  }
}
