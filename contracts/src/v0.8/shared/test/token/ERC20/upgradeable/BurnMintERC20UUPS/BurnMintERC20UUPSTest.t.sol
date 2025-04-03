// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IERC1822Proxiable} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/draft-IERC1822.sol";
import {ERC1967Proxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {BurnMintERC20UUPS, Initializable, UUPSUpgradeable} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";

import {ERC20UpgradableBaseTest_approve} from "../ERC20UpgradableBaseTest.approve.t.sol";
import {ERC20UpgradableBaseTest_burn} from "../ERC20UpgradableBaseTest.burn.t.sol";
import {ERC20UpgradableBaseTest_burnFrom} from "../ERC20UpgradableBaseTest.burnFrom.t.sol";
import {ERC20UpgradableBaseTest_burnFrom_alias} from "../ERC20UpgradableBaseTest.burnFrom_alias.t.sol";
import {ERC20UpgradableBaseTest_initialize} from "../ERC20UpgradableBaseTest.initialize.t.sol";
import {ERC20UpgradableBaseTest_mint} from "../ERC20UpgradableBaseTest.mint.t.sol";
import {ERC20UpgradableBaseTest_roles} from "../ERC20UpgradableBaseTest.roles.t.sol";
import {ERC20UpgradableBaseTest_supportsInterface} from "../ERC20UpgradableBaseTest.supportsInterface.t.sol";
import {ERC20UpgradableBaseTest} from "../ERC20UpgradableBaseTest.t.sol";

contract BurnMintERC20UUPSTest is
  ERC20UpgradableBaseTest_initialize,
  ERC20UpgradableBaseTest_supportsInterface,
  ERC20UpgradableBaseTest_approve,
  ERC20UpgradableBaseTest_mint,
  ERC20UpgradableBaseTest_burn,
  ERC20UpgradableBaseTest_burnFrom,
  ERC20UpgradableBaseTest_burnFrom_alias,
  ERC20UpgradableBaseTest_roles
{
  BurnMintERC20UUPS internal s_burnMintERC20UUPS;

  // ================================================================
  // │                           Set Up                             │
  // ================================================================

  function deployBurnMintERC20UUPS(
    string memory name,
    string memory symbol,
    uint8 decimals,
    uint256 maxSupply,
    uint256 preMint,
    address defaultAdmin,
    address defaultUpgrader
  ) public returns (BurnMintERC20UUPS) {
    address implementation = address(new BurnMintERC20UUPS());

    address proxy = address(
      new ERC1967Proxy(
        implementation,
        abi.encodeCall(
          BurnMintERC20UUPS.initialize,
          (name, symbol, decimals, maxSupply, preMint, defaultAdmin, defaultUpgrader)
        )
      )
    );

    return BurnMintERC20UUPS(proxy);
  }

  function setUp() public virtual override {
    ERC20UpgradableBaseTest.setUp();

    s_burnMintERC20UUPS = deployBurnMintERC20UUPS(
      NAME,
      SYMBOL,
      DECIMALS,
      MAX_SUPPLY,
      PRE_MINT,
      DEFAULT_ADMIN,
      DEFAULT_UPGRADER
    );

    // Grant mint and burn role to the mock pool address
    s_burnMintERC20UUPS.grantMintAndBurnRoles(i_mockPool);
  }

  // ================================================================
  // │                            UUPS                              │
  // ================================================================

  function test_Initialize() public view {
    should_Initialize(address(s_burnMintERC20UUPS), s_burnMintERC20UUPS.DEFAULT_ADMIN_ROLE());
    assertTrue(s_burnMintERC20UUPS.hasRole(s_burnMintERC20UUPS.UPGRADER_ROLE(), DEFAULT_UPGRADER));
  }

  function test_Initialize_WithPreMint() public {
    uint256 newPreMint = 1e18;
    BurnMintERC20UUPS newBurnMintERC20UUPS = deployBurnMintERC20UUPS(
      NAME,
      SYMBOL,
      DECIMALS,
      MAX_SUPPLY,
      newPreMint,
      DEFAULT_ADMIN,
      DEFAULT_UPGRADER
    );

    should_Initialize_WithPreMint(address(newBurnMintERC20UUPS), newPreMint);
  }

  function test_Initialize_RevertWhen_PreMintExceedsMaxSupply() public {
    uint256 newPreMint = MAX_SUPPLY + 1;
    address implementation = address(new BurnMintERC20UUPS());

    vm.expectRevert(
      abi.encodeWithSelector(BurnMintERC20UUPS.BurnMintERC20UUPS__MaxSupplyExceeded.selector, newPreMint)
    );

    new ERC1967Proxy(
      implementation,
      abi.encodeCall(
        BurnMintERC20UUPS.initialize,
        (NAME, SYMBOL, DECIMALS, MAX_SUPPLY, newPreMint, DEFAULT_ADMIN, DEFAULT_UPGRADER)
      )
    );
  }

  function test_Initialize_RevertWhen_AlreadyInitialized() public {
    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));
    s_burnMintERC20UUPS.initialize(NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN, DEFAULT_UPGRADER);
  }

  /// @dev Adding _disableInitializers() function to implementation's constructor ensures that no one can call initialize directly on the implementation.
  /// @dev The initialize should be only callable through Proxy.
  /// @dev This test tests that case.
  function test_Initialize_RevertWhen_CallIsNotThroughProxy() public {
    BurnMintERC20UUPS newBurnMintERC20UUPS = new BurnMintERC20UUPS();

    vm.expectRevert(abi.encodeWithSelector(Initializable.InvalidInitialization.selector));
    newBurnMintERC20UUPS.initialize(NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN, DEFAULT_UPGRADER);
  }

  function test_ProxiableUUID() public {
    BurnMintERC20UUPS implementation = new BurnMintERC20UUPS();
    bytes32 proxiableUUID = IERC1822Proxiable(address(implementation)).proxiableUUID();
    bytes32 expectedProxiableUUID = bytes32(uint256(keccak256("eip1967.proxy.implementation")) - 1);
    assertEq(proxiableUUID, expectedProxiableUUID);
  }

  function test_ProxiableUUID_RevertWhen_CalledThroughProxy() public {
    vm.expectRevert(abi.encodeWithSelector(UUPSUpgradeable.UUPSUnauthorizedCallContext.selector));
    IERC1822Proxiable(s_burnMintERC20UUPS).proxiableUUID();
  }

  // ================================================================
  // │                           ERC165                             │
  // ================================================================

  function test_SupportsInterface() public view {
    should_SupportsInterface(s_burnMintERC20UUPS);
    assertTrue(s_burnMintERC20UUPS.supportsInterface(type(IERC1822Proxiable).interfaceId));
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  function test_Approve() public {
    should_Approve(address(s_burnMintERC20UUPS));
  }

  function test_Approve_RevertWhen_RecipientIsImplementationItself() public {
    should_Approve_RevertWhen_RecipientIsImplementationItself(
      address(s_burnMintERC20UUPS),
      BurnMintERC20UUPS.BurnMintERC20UUPS__InvalidRecipient.selector
    );
  }

  // ================================================================
  // │                      Burning & minting                       │
  // ================================================================

  function test_Mint() public {
    should_Mint(address(s_burnMintERC20UUPS));
  }

  function test_Mint_RevertWhen_CallerDoesNotHaveMinterRole() public {
    should_Mint_RevertWhen_CallerDoesNotHaveMinterRole(address(s_burnMintERC20UUPS), s_burnMintERC20UUPS.MINTER_ROLE());
  }

  function test_Mint_RevertWhen_AmountExceedsMaxSupply() public {
    should_Mint_RevertWhen_AmountExceedsMaxSupply(
      address(s_burnMintERC20UUPS),
      BurnMintERC20UUPS.BurnMintERC20UUPS__MaxSupplyExceeded.selector
    );
  }

  function test_Mint_RevertWhen_RecipientIsImplementationItself() public {
    should_Mint_RevertWhen_RecipientIsImplementationItself(
      address(s_burnMintERC20UUPS),
      BurnMintERC20UUPS.BurnMintERC20UUPS__InvalidRecipient.selector
    );
  }

  function test_Burn() public {
    should_Burn(address(s_burnMintERC20UUPS));
  }

  function test_Burn_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    should_Burn_RevertWhen_CallerDoesNotHaveBurnerRole(address(s_burnMintERC20UUPS), s_burnMintERC20UUPS.BURNER_ROLE());
  }

  function test_BurnFrom() public {
    should_BurnFrom(address(s_burnMintERC20UUPS));
  }

  function test_BurnFrom_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    should_BurnFrom_RevertWhen_CallerDoesNotHaveBurnerRole(
      address(s_burnMintERC20UUPS),
      s_burnMintERC20UUPS.BURNER_ROLE()
    );
  }

  function test_BurnFrom_alias() public {
    should_BurnFrom_alias(address(s_burnMintERC20UUPS));
  }

  function test_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole() public {
    should_BurnFrom_alias_RevertWhen_CallerDoesNotHaveBurnerRole(
      address(s_burnMintERC20UUPS),
      s_burnMintERC20UUPS.BURNER_ROLE()
    );
  }

  // ================================================================
  // │                            Roles                             │
  // ================================================================

  function test_GrantMintAndBurnRoles() public {
    should_GrantMintAndBurnRoles(
      address(s_burnMintERC20UUPS),
      s_burnMintERC20UUPS.MINTER_ROLE(),
      s_burnMintERC20UUPS.BURNER_ROLE()
    );
  }

  function test_GetCCIPAdmin() public view {
    should_GetCCIPAdmin(address(s_burnMintERC20UUPS));
  }

  function test_SetCCIPAdmin() public {
    should_SetCCIPAdmin(address(s_burnMintERC20UUPS));
  }
}
