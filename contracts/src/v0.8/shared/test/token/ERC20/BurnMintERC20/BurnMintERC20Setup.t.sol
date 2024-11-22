// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BurnMintERC20} from "../../../../token/ERC20/BurnMintERC20.sol";
import {BaseTest} from "../../../BaseTest.t.sol";

contract BurnMintERC20Setup is BaseTest {
  BurnMintERC20 internal s_burnMintERC20;

  address internal s_mockPool = makeAddr("s_mockPool");
  uint256 internal s_amount = 1e18;

  address internal s_alice;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_alice = makeAddr("alice");

    s_burnMintERC20 = new BurnMintERC20("Chainlink Token", "LINK", 18, 1e27, 0);

    // Set s_mockPool to be a burner and minter
    s_burnMintERC20.grantMintAndBurnRoles(s_mockPool);
    deal(address(s_burnMintERC20), OWNER, s_amount);
  }
}
