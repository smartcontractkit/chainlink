pragma solidity ^0.8.0;

import {ERC20Mock} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";

// mock ERC20 with 6 decimals
contract ERC20Mock6Decimals is ERC20Mock {
  constructor(
    string memory name,
    string memory symbol,
    address initialAccount,
    uint256 initialBalance
  ) payable ERC20Mock(name, symbol, initialAccount, initialBalance) {}

  function decimals() public view virtual override returns (uint8) {
    return 6;
  }
}
