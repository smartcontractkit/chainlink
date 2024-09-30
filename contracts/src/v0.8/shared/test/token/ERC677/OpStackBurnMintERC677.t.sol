// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IBurnMintERC20} from "../../../token/ERC20/IBurnMintERC20.sol";
import {IOptimismMintableERC20Minimal, IOptimismMintableERC20} from "../../../token/ERC20/IOptimismMintableERC20.sol";
import {IERC677} from "../../../token/ERC677/IERC677.sol";

import {BurnMintERC677} from "../../../token/ERC677/BurnMintERC677.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {OpStackBurnMintERC677} from "../../../token/ERC677/OpStackBurnMintERC677.sol";

import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

contract OpStackBurnMintERC677Setup is BaseTest {
  address internal s_l1Token = address(897352983527);
  address internal s_l2Bridge = address(1928235235);

  OpStackBurnMintERC677 internal s_opStackBurnMintERC677;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_opStackBurnMintERC677 = new OpStackBurnMintERC677("Chainlink Token", "LINK", 18, 1e27, s_l1Token, s_l2Bridge);
  }
}

contract OpStackBurnMintERC677_constructor is OpStackBurnMintERC677Setup {
  function testConstructorSuccess() public {
    string memory name = "Chainlink token l2";
    string memory symbol = "LINK L2";
    uint8 decimals = 18;
    uint256 maxSupply = 1e33;
    s_opStackBurnMintERC677 = new OpStackBurnMintERC677(name, symbol, decimals, maxSupply, s_l1Token, s_l2Bridge);

    assertEq(name, s_opStackBurnMintERC677.name());
    assertEq(symbol, s_opStackBurnMintERC677.symbol());
    assertEq(decimals, s_opStackBurnMintERC677.decimals());
    assertEq(maxSupply, s_opStackBurnMintERC677.maxSupply());
    assertEq(s_l1Token, s_opStackBurnMintERC677.remoteToken());
    assertEq(s_l2Bridge, s_opStackBurnMintERC677.bridge());
  }
}

contract OpStackBurnMintERC677_supportsInterface is OpStackBurnMintERC677Setup {
  function testConstructorSuccess() public view {
    assertTrue(s_opStackBurnMintERC677.supportsInterface(type(IOptimismMintableERC20Minimal).interfaceId));
    assertTrue(s_opStackBurnMintERC677.supportsInterface(type(IERC677).interfaceId));
    assertTrue(s_opStackBurnMintERC677.supportsInterface(type(IBurnMintERC20).interfaceId));
    assertTrue(s_opStackBurnMintERC677.supportsInterface(type(IERC165).interfaceId));
  }
}

contract OpStackBurnMintERC677_interfaceCompatibility is OpStackBurnMintERC677Setup {
  event Transfer(address indexed from, address indexed to, uint256 value);

  IOptimismMintableERC20 internal s_opStackToken;

  function setUp() public virtual override {
    OpStackBurnMintERC677Setup.setUp();
    s_opStackToken = IOptimismMintableERC20(address(s_opStackBurnMintERC677));
  }

  function testStaticFunctionsCompatibility() public {
    assertEq(s_l1Token, s_opStackToken.remoteToken());
    assertEq(s_l2Bridge, s_opStackToken.bridge());
  }

  function testMintCompatibility() public {
    // Ensure roles work
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC677.SenderNotMinter.selector, OWNER));
    s_opStackToken.mint(OWNER, 1);

    // Use the actual contract to grant mint
    s_opStackBurnMintERC677.grantMintRole(OWNER);

    // Ensure zero address check works
    vm.expectRevert("ERC20: mint to the zero address");
    s_opStackToken.mint(address(0x0), 0);

    address mintToAddress = address(0x1);
    uint256 mintAmount = 1;

    vm.expectEmit();
    emit Transfer(address(0), mintToAddress, mintAmount);

    s_opStackToken.mint(mintToAddress, mintAmount);
  }

  function testBurnCompatibility() public {
    // Ensure roles work
    vm.expectRevert(abi.encodeWithSelector(BurnMintERC677.SenderNotBurner.selector, OWNER));
    s_opStackToken.burn(address(0x0), 1);

    // Use the actual contract to grant burn
    s_opStackBurnMintERC677.grantBurnRole(OWNER);

    // Ensure zero address check works
    vm.expectRevert("ERC20: approve from the zero address");
    s_opStackToken.burn(address(0x0), 0);

    address burnFromAddress = address(0x1);
    uint256 burnAmount = 1;

    // Ensure `burn(address, amount)` works like burnFrom and requires allowance
    vm.expectRevert("ERC20: insufficient allowance");
    s_opStackToken.burn(burnFromAddress, burnAmount);

    changePrank(burnFromAddress);
    deal(address(s_opStackToken), burnFromAddress, burnAmount);
    s_opStackBurnMintERC677.approve(OWNER, burnAmount);
    changePrank(OWNER);

    vm.expectEmit();
    emit Transfer(burnFromAddress, address(0x0), burnAmount);

    s_opStackToken.burn(burnFromAddress, burnAmount);
  }
}
