// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {FactoryBurnMintERC20} from "../../tokenAdminRegistry/TokenPoolFactory/FactoryBurnMintERC20.sol";
import {BaseTest} from "../BaseTest.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

contract BurnMintERC20Setup is BaseTest {
  FactoryBurnMintERC20 internal s_burnMintERC20;

  address internal s_mockPool = makeAddr("s_mockPool");
  uint256 internal s_amount = 1e18;

  address internal s_alice;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_alice = makeAddr("alice");

    s_burnMintERC20 = new FactoryBurnMintERC20("Chainlink Token", "LINK", 18, 1e27, 0, s_alice);

    // Set s_mockPool to be a burner and minter
    s_burnMintERC20.grantMintAndBurnRoles(s_mockPool);
    deal(address(s_burnMintERC20), OWNER, s_amount);
  }
}

contract FactoryBurnMintERC20constructor is BurnMintERC20Setup {
  function test_Constructor_Success() public {
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

contract FactoryBurnMintERC20approve is BurnMintERC20Setup {
  function test_Approve_Success() public {
    uint256 balancePre = s_burnMintERC20.balanceOf(STRANGER);
    uint256 sendingAmount = s_amount / 2;

    s_burnMintERC20.approve(STRANGER, sendingAmount);

    changePrank(STRANGER);

    s_burnMintERC20.transferFrom(OWNER, STRANGER, sendingAmount);

    assertEq(sendingAmount + balancePre, s_burnMintERC20.balanceOf(STRANGER));
  }

  // Reverts

  function test_InvalidAddress_Reverts() public {
    vm.expectRevert();

    s_burnMintERC20.approve(address(s_burnMintERC20), s_amount);
  }
}

contract FactoryBurnMintERC20transfer is BurnMintERC20Setup {
  function test_Transfer_Success() public {
    uint256 balancePre = s_burnMintERC20.balanceOf(STRANGER);
    uint256 sendingAmount = s_amount / 2;

    s_burnMintERC20.transfer(STRANGER, sendingAmount);

    assertEq(sendingAmount + balancePre, s_burnMintERC20.balanceOf(STRANGER));
  }

  // Reverts

  function test_InvalidAddress_Reverts() public {
    vm.expectRevert();

    s_burnMintERC20.transfer(address(s_burnMintERC20), s_amount);
  }
}

contract FactoryBurnMintERC20mint is BurnMintERC20Setup {
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

contract FactoryBurnMintERC20burn is BurnMintERC20Setup {
  function test_BasicBurn_Success() public {
    s_burnMintERC20.grantBurnRole(OWNER);
    deal(address(s_burnMintERC20), OWNER, s_amount);

    vm.expectEmit();
    emit IERC20.Transfer(OWNER, address(0), s_amount);

    s_burnMintERC20.burn(s_amount);

    assertEq(0, s_burnMintERC20.balanceOf(OWNER));
  }

  // Revert

  function test_SenderNotBurner_Reverts() public {
    vm.expectRevert(abi.encodeWithSelector(FactoryBurnMintERC20.SenderNotBurner.selector, OWNER));

    s_burnMintERC20.burnFrom(STRANGER, s_amount);
  }

  function test_ExceedsBalance_Reverts() public {
    changePrank(s_mockPool);

    vm.expectRevert("ERC20: burn amount exceeds balance");

    s_burnMintERC20.burn(s_amount * 2);
  }

  function test_BurnFromZeroAddress_Reverts() public {
    s_burnMintERC20.grantBurnRole(address(0));
    changePrank(address(0));

    vm.expectRevert("ERC20: burn from the zero address");

    s_burnMintERC20.burn(0);
  }
}

contract FactoryBurnMintERC20burnFromAlias is BurnMintERC20Setup {
  function setUp() public virtual override {
    BurnMintERC20Setup.setUp();
  }

  function test_BurnFrom_Success() public {
    s_burnMintERC20.approve(s_mockPool, s_amount);

    changePrank(s_mockPool);

    s_burnMintERC20.burn(OWNER, s_amount);

    assertEq(0, s_burnMintERC20.balanceOf(OWNER));
  }

  // Reverts

  function test_SenderNotBurner_Reverts() public {
    vm.expectRevert(abi.encodeWithSelector(FactoryBurnMintERC20.SenderNotBurner.selector, OWNER));

    s_burnMintERC20.burn(OWNER, s_amount);
  }

  function test_InsufficientAllowance_Reverts() public {
    changePrank(s_mockPool);

    vm.expectRevert("ERC20: insufficient allowance");

    s_burnMintERC20.burn(OWNER, s_amount);
  }

  function test_ExceedsBalance_Reverts() public {
    s_burnMintERC20.approve(s_mockPool, s_amount * 2);

    changePrank(s_mockPool);

    vm.expectRevert("ERC20: burn amount exceeds balance");

    s_burnMintERC20.burn(OWNER, s_amount * 2);
  }
}

contract FactoryBurnMintERC20burnFrom is BurnMintERC20Setup {
  function setUp() public virtual override {
    BurnMintERC20Setup.setUp();
  }

  function test_BurnFrom_Success() public {
    s_burnMintERC20.approve(s_mockPool, s_amount);

    changePrank(s_mockPool);

    s_burnMintERC20.burnFrom(OWNER, s_amount);

    assertEq(0, s_burnMintERC20.balanceOf(OWNER));
  }

  // Reverts

  function test_SenderNotBurner_Reverts() public {
    vm.expectRevert(abi.encodeWithSelector(FactoryBurnMintERC20.SenderNotBurner.selector, OWNER));

    s_burnMintERC20.burnFrom(OWNER, s_amount);
  }

  function test_InsufficientAllowance_Reverts() public {
    changePrank(s_mockPool);

    vm.expectRevert("ERC20: insufficient allowance");

    s_burnMintERC20.burnFrom(OWNER, s_amount);
  }

  function test_ExceedsBalance_Reverts() public {
    s_burnMintERC20.approve(s_mockPool, s_amount * 2);

    changePrank(s_mockPool);

    vm.expectRevert("ERC20: burn amount exceeds balance");

    s_burnMintERC20.burnFrom(OWNER, s_amount * 2);
  }
}

contract FactoryBurnMintERC20grantRole is BurnMintERC20Setup {
  function test_GrantMintAccess_Success() public {
    assertFalse(s_burnMintERC20.isMinter(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.MintAccessGranted(STRANGER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.isMinter(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.MintAccessRevoked(STRANGER);

    s_burnMintERC20.revokeMintRole(STRANGER);

    assertFalse(s_burnMintERC20.isMinter(STRANGER));
  }

  function test_GrantBurnAccess_Success() public {
    assertFalse(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.BurnAccessGranted(STRANGER);

    s_burnMintERC20.grantBurnRole(STRANGER);

    assertTrue(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.BurnAccessRevoked(STRANGER);

    s_burnMintERC20.revokeBurnRole(STRANGER);

    assertFalse(s_burnMintERC20.isBurner(STRANGER));
  }

  function test_GrantMany_Success() public {
    // Since alice was already granted mint and burn roles in the setup, we will revoke them
    // and then grant them again for the purposes of the test
    s_burnMintERC20.revokeMintRole(s_alice);
    s_burnMintERC20.revokeBurnRole(s_alice);

    uint256 numberOfPools = 10;
    address[] memory permissionedAddresses = new address[](numberOfPools + 1);
    permissionedAddresses[0] = s_mockPool;

    for (uint160 i = 0; i < numberOfPools; ++i) {
      permissionedAddresses[i + 1] = address(i);
      s_burnMintERC20.grantMintAndBurnRoles(address(i));
    }

    assertEq(permissionedAddresses, s_burnMintERC20.getBurners());
    assertEq(permissionedAddresses, s_burnMintERC20.getMinters());
  }
}

contract FactoryBurnMintERC20grantMintAndBurnRoles is BurnMintERC20Setup {
  function test_GrantMintAndBurnRoles_Success() public {
    assertFalse(s_burnMintERC20.isMinter(STRANGER));
    assertFalse(s_burnMintERC20.isBurner(STRANGER));

    vm.expectEmit();
    emit FactoryBurnMintERC20.MintAccessGranted(STRANGER);
    vm.expectEmit();
    emit FactoryBurnMintERC20.BurnAccessGranted(STRANGER);

    s_burnMintERC20.grantMintAndBurnRoles(STRANGER);

    assertTrue(s_burnMintERC20.isMinter(STRANGER));
    assertTrue(s_burnMintERC20.isBurner(STRANGER));
  }
}

contract FactoryBurnMintERC20decreaseApproval is BurnMintERC20Setup {
  function test_DecreaseApproval_Success() public {
    s_burnMintERC20.approve(s_mockPool, s_amount);
    uint256 allowance = s_burnMintERC20.allowance(OWNER, s_mockPool);
    assertEq(allowance, s_amount);
    s_burnMintERC20.decreaseApproval(s_mockPool, s_amount);
    assertEq(s_burnMintERC20.allowance(OWNER, s_mockPool), allowance - s_amount);
  }
}

contract FactoryBurnMintERC20increaseApproval is BurnMintERC20Setup {
  function test_IncreaseApproval_Success() public {
    s_burnMintERC20.approve(s_mockPool, s_amount);
    uint256 allowance = s_burnMintERC20.allowance(OWNER, s_mockPool);
    assertEq(allowance, s_amount);
    s_burnMintERC20.increaseApproval(s_mockPool, s_amount);
    assertEq(s_burnMintERC20.allowance(OWNER, s_mockPool), allowance + s_amount);
  }
}

contract FactoryBurnMintERC20supportsInterface is BurnMintERC20Setup {
  function test_SupportsInterface_Success() public view {
    assertTrue(s_burnMintERC20.supportsInterface(type(IERC20).interfaceId));
    assertTrue(s_burnMintERC20.supportsInterface(type(IBurnMintERC20).interfaceId));
    assertTrue(s_burnMintERC20.supportsInterface(type(IERC165).interfaceId));
  }
}
