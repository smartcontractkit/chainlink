// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_constructor is BurnMintERC20Setup {
  function test_Constructor() public {
    string memory name = "Chainlink token v2";
    string memory symbol = "LINK2";
    uint8 decimals = 19;
    uint256 maxSupply = 1e33;

    s_burnMintERC20 = new FactoryBurnMintERC20(name, symbol, decimals, maxSupply, 1e18, s_alice);

    assertEq(name, s_burnMintERC20.name());
    assertEq(symbol, s_burnMintERC20.symbol());
    assertEq(decimals, s_burnMintERC20.decimals());
    assertEq(maxSupply, s_burnMintERC20.maxSupply());

    assertTrue(s_burnMintERC20.isMinter(s_alice));
    assertTrue(s_burnMintERC20.isBurner(s_alice));
    assertEq(s_burnMintERC20.balanceOf(s_alice), 1e18);
    assertEq(s_burnMintERC20.totalSupply(), 1e18);
  }
}
