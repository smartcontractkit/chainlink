// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {FactoryBurnMintERC20} from "../../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {BurnMintERC20Setup} from "./BurnMintERC20Setup.t.sol";

contract FactoryBurnMintERC20_mint is BurnMintERC20Setup {
  function test_BasicMint_Success() public {
    uint256 balancePre = s_burnMintERC20.balanceOf(OWNER);

    s_burnMintERC20.grantMintAndBurnRoles(OWNER);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), OWNER, s_amount);

    s_burnMintERC20.mint(OWNER, s_amount);

    assertEq(balancePre + s_amount, s_burnMintERC20.balanceOf(OWNER));
  }

  // Revert

  function test_SenderNotMinter_Reverts() public {
    vm.expectRevert(abi.encodeWithSelector(FactoryBurnMintERC20.SenderNotMinter.selector, OWNER));
    s_burnMintERC20.mint(STRANGER, 1e18);
  }

  function test_MaxSupplyExceeded_Reverts() public {
    changePrank(s_mockPool);

    // Mint max supply
    s_burnMintERC20.mint(OWNER, s_burnMintERC20.maxSupply());

    vm.expectRevert(
      abi.encodeWithSelector(FactoryBurnMintERC20.MaxSupplyExceeded.selector, s_burnMintERC20.maxSupply() + 1)
    );

    // Attempt to mint 1 more than max supply
    s_burnMintERC20.mint(OWNER, 1);
  }
}
