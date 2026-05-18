// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";
import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";

contract BurnMintERC20_constructor is BurnMintERC20Setup {
  function test_Constructor() public {
    vm.startPrank(s_alice);

    string memory name = "Chainlink token v2";
    string memory symbol = "LINK2";
    uint8 decimals = 19;
    uint256 maxSupply = 1e33;

    s_burnMintERC20 = new BurnMintERC20(name, symbol, decimals, maxSupply, 1e18);

    assertEq(name, s_burnMintERC20.name());
    assertEq(symbol, s_burnMintERC20.symbol());
    assertEq(decimals, s_burnMintERC20.decimals());
    assertEq(maxSupply, s_burnMintERC20.maxSupply());

    assertTrue(s_burnMintERC20.hasRole(s_burnMintERC20.DEFAULT_ADMIN_ROLE(), s_alice));
    assertEq(s_burnMintERC20.balanceOf(s_alice), 1e18);
    assertEq(s_burnMintERC20.totalSupply(), 1e18);
  }
}
