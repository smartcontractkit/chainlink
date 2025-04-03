// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {TransparentUpgradeableProxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {BurnMintERC20Transparent, Initializable} from "../../../../../token/ERC20/upgradeable/BurnMintERC20Transparent.sol";

import {ERC20UpgradableBaseTest_approve} from "../ERC20UpgradableBaseTest.approve.t.sol";
import {ERC20UpgradableBaseTest_burn} from "../ERC20UpgradableBaseTest.burn.t.sol";
import {ERC20UpgradableBaseTest_burnFrom} from "../ERC20UpgradableBaseTest.burnFrom.t.sol";
import {ERC20UpgradableBaseTest_burnFrom_alias} from "../ERC20UpgradableBaseTest.burnFrom_alias.t.sol";
import {ERC20UpgradableBaseTest_initialize} from "../ERC20UpgradableBaseTest.initialize.t.sol";
import {ERC20UpgradableBaseTest_mint} from "../ERC20UpgradableBaseTest.mint.t.sol";
import {ERC20UpgradableBaseTest_roles} from "../ERC20UpgradableBaseTest.roles.t.sol";
import {ERC20UpgradableBaseTest_supportsInterface} from "../ERC20UpgradableBaseTest.supportsInterface.t.sol";
import {ERC20UpgradableBaseTest} from "../ERC20UpgradableBaseTest.t.sol";

contract BurnMintERC20TransparentTest is
  ERC20UpgradableBaseTest_initialize,
  ERC20UpgradableBaseTest_supportsInterface,
  ERC20UpgradableBaseTest_approve,
  ERC20UpgradableBaseTest_mint,
  ERC20UpgradableBaseTest_burn,
  ERC20UpgradableBaseTest_burnFrom,
  ERC20UpgradableBaseTest_burnFrom_alias,
  ERC20UpgradableBaseTest_roles
{
  BurnMintERC20Transparent internal s_burnMintERC20Transparent;

  // ================================================================
  // │                           Set Up                             │
  // ================================================================

  function deployBurnMintERC20Transparent(
    string memory name,
    string memory symbol,
    uint8 decimals,
    uint256 maxSupply,
    uint256 preMint,
    address defaultAdmin
  ) public returns (BurnMintERC20Transparent) {
    address implementation = address(new BurnMintERC20Transparent());

    address proxy = address(
      new TransparentUpgradeableProxy(
        implementation,
        INITIAL_OWNER_ADDRESS_FOR_PROXY_ADMIN,
        abi.encodeCall(BurnMintERC20Transparent.initialize, (name, symbol, decimals, maxSupply, preMint, defaultAdmin))
      )
    );

    return BurnMintERC20Transparent(proxy);
  }

  function setUp() public virtual override {
    ERC20UpgradableBaseTest.setUp();

    s_burnMintERC20Transparent = deployBurnMintERC20Transparent(
      NAME,
      SYMBOL,
      DECIMALS,
      MAX_SUPPLY,
      PRE_MINT,
      DEFAULT_ADMIN
    );

    s_burnMintERC20Transparent.grantMintAndBurnRoles(i_mockPool);
  }

  // ================================================================
  // │                         Transparent                          │
  // ================================================================

  function test_Initialize() public view {
    should_Initialize(address(s_burnMintERC20Transparent), s_burnMintERC20Transparent.DEFAULT_ADMIN_ROLE());
  }

  function test_Initialize_WithPreMint() public {
    uint256 newPreMint = 1e18;
    BurnMintERC20Transparent newBurnMintERC20Transparent = deployBurnMintERC20Transparent(
      NAME,
      SYMBOL,
      DECIMALS,
      MAX_SUPPLY,
      newPreMint,
      DEFAULT_ADMIN
    );

    should_Initialize_WithPreMint(address(newBurnMintERC20Transparent), newPreMint);
  }

  function test_Initialize_RevertWhen_PreMintExceedsMaxSupply() public {
    uint256 newPreMint = MAX_SUPPLY + 1;
    address implementation = address(new BurnMintERC20Transparent());

    vm.expectRevert(
      abi.encodeWithSelector(BurnMintERC20Transparent.BurnMintERC20Transparent__MaxSupplyExceeded.selector, newPreMint)
    );

    new TransparentUpgradeableProxy(
      implementation,
      INITIAL_OWNER_ADDRESS_FOR_PROXY_ADMIN,
      abi.encodeCall(
        BurnMintERC20Transparent.initialize,
        (NAME, SYMBOL, DECIMALS, MAX_SUPPLY, newPreMint, DEFAULT_ADMIN)
      )
    );
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));
    s_burnMintERC20Transparent.initialize(NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN);
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20Transparent newBurnMintERC20Transparent = new BurnMintERC20Transparent();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));
    newBurnMintERC20Transparent.initialize(NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN);
  }

  // ================================================================
  // │                           ERC165                             │
  // ================================================================

  function test_SupportsInterface() public view {
    should_SupportsInterface(s_burnMintERC20Transparent);
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  function test_Approve() public {
    should_Approve(address(s_burnMintERC20Transparent));
  }

  function test_Approve_RevertWhen_RecipientIsImplementationItself() public {
    should_Approve_RevertWhen_RecipientIsImplementationItself(
      address(s_burnMintERC20Transparent),
      BurnMintERC20Transparent.BurnMintERC20Transparent__InvalidRecipient.selector
    );
  }

  // ================================================================
  // │                      Burning & minting                       │
  // ================================================================

  function test_Mint() public {
    should_Mint(address(s_burnMintERC20Transparent));
  }

  function test_Mint_RevertWhen_CallerDoesNotHaveMinterRole() public {
    should_Mint_RevertWhen_CallerDoesNotHaveMinterRole(
      address(s_burnMintERC20Transparent),
      s_burnMintERC20Transparent.MINTER_ROLE()
    );
  }

  function test_Mint_RevertWhen_AmountExceedsMaxSupply() public {
    should_Mint_RevertWhen_AmountExceedsMaxSupply(
      address(s_burnMintERC20Transparent),
      BurnMintERC20Transparent.BurnMintERC20Transparent__MaxSupplyExceeded.selector
    );
  }

  function test_Mint_RevertWhen_RecipientIsImplementationItself() public {
    should_Mint_RevertWhen_RecipientIsImplementationItself(
      address(s_burnMintERC20Transparent),
      BurnMintERC20Transparent.BurnMintERC20Transparent__InvalidRecipient.selector
    );
  }

  function test_Burn() public {
    should_Burn(address(s_burnMintERC20Transparent));
  }

  function test_Burn_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    should_Burn_RevertWhen_CallerDoesNotHaveBurnerRole(
      address(s_burnMintERC20Transparent),
      s_burnMintERC20Transparent.BURNER_ROLE()
    );
  }

  function test_BurnFrom() public {
    should_BurnFrom(address(s_burnMintERC20Transparent));
  }

  function test_BurnFrom_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    should_BurnFrom_RevertWhen_CallerDoesNotHaveBurnerRole(
      address(s_burnMintERC20Transparent),
      s_burnMintERC20Transparent.BURNER_ROLE()
    );
  }

  function test_BurnFrom_alias() public {
    should_BurnFrom_alias(address(s_burnMintERC20Transparent));
  }

  function test_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    should_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole(
      address(s_burnMintERC20Transparent),
      s_burnMintERC20Transparent.BURNER_ROLE()
    );
  }

  // ================================================================
  // │                            Roles                             │
  // ================================================================

  function test_GrantMintAndBurnRoles() public {
    should_GrantMintAndBurnRoles(
      address(s_burnMintERC20Transparent),
      s_burnMintERC20Transparent.MINTER_ROLE(),
      s_burnMintERC20Transparent.BURNER_ROLE()
    );
  }

  function test_GetCCIPAdmin() public view {
    should_GetCCIPAdmin(address(s_burnMintERC20Transparent));
  }

  function test_SetCCIPAdmin() public {
    should_SetCCIPAdmin(address(s_burnMintERC20Transparent));
  }
}
