// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC1967Proxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {BurnMintERC20PausableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20UUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {ERC20UpgradableBaseTest_pausing} from "../ERC20UpgradableBaseTest.pausing.t.sol";

interface IUpgradeableProxy {
  function upgradeToAndCall(address, bytes memory) external payable;
}

/// @dev Mock contract with limited functionality used for testing only.
/// @dev It adds a new freeze functionality to the BurnMintERC20PausableUUPS contract.
/// @dev We want to test that the new implementation works as expected and keeps track of the balances from previous version.
contract MockBurnMintERC20PausableUUPSV2 is BurnMintERC20PausableUUPS {
  error MockBurnMintERC20PausableUUPSV2__AccountFrozen(address account);

  bytes32 public constant FREEZER_ROLE = keccak256("FREEZER_ROLE");

  /// @custom:storage-location erc7201:chainlink.storage.MockBurnMintERC20PausableUUPSV2
  struct MockBurnMintERC20PausableUUPSV2Storage {
    mapping(address => bool) s_isFrozen;
  }

  // keccak256(abi.encode(uint256(keccak256("chainlink.storage.MockBurnMintERC20PausableUUPSV2")) - 1)) & ~bytes32(uint256(0xff));
  bytes32 private constant V2_STORAGE_LOCATION = 0x98bca5456fc57bb77324f8627b5055944605eb027b3a0652fea6ac1ede88a400;

  function _getV2Storage() private pure returns (MockBurnMintERC20PausableUUPSV2Storage storage $) {
    assembly {
      $.slot := V2_STORAGE_LOCATION
    }
  }

  function initializeFreezerRole(address defaultFreezer) public onlyRole(UPGRADER_ROLE) {
    _grantRole(FREEZER_ROLE, defaultFreezer);
  }

  function freeze(address account) public onlyRole(FREEZER_ROLE) {
    MockBurnMintERC20PausableUUPSV2Storage storage $ = _getV2Storage();
    $.s_isFrozen[account] = true;
  }

  function _update(address from, address to, uint256 value) internal virtual override {
    MockBurnMintERC20PausableUUPSV2Storage storage $ = _getV2Storage();
    if ($.s_isFrozen[from]) revert MockBurnMintERC20PausableUUPSV2__AccountFrozen(from);
    if ($.s_isFrozen[to]) revert MockBurnMintERC20PausableUUPSV2__AccountFrozen(to);

    super._update(from, to, value);
  }
}

contract BurnMintERC20PausableUUPS_upgrade is ERC20UpgradableBaseTest_pausing {
  BurnMintERC20PausableUUPS internal s_burnMintERC20PausableUUPS;
  address internal s_uupsProxy;

  function setUp() public virtual override {
    address implementation = address(new BurnMintERC20PausableUUPS());

    s_uupsProxy = address(
      new ERC1967Proxy(
        implementation,
        abi.encodeCall(
          BurnMintERC20UUPS.initialize,
          (NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN, DEFAULT_UPGRADER)
        )
      )
    );

    s_burnMintERC20PausableUUPS = BurnMintERC20PausableUUPS(s_uupsProxy);

    changePrank(DEFAULT_ADMIN);
    s_burnMintERC20PausableUUPS.grantRole(s_burnMintERC20PausableUUPS.PAUSER_ROLE(), DEFAULT_PAUSER);
    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(i_mockPool);
  }

  function test_Upgrade() public {
    address defaultFreezer = makeAddr("defaultFreezer");

    changePrank(DEFAULT_ADMIN);
    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(DEFAULT_ADMIN);
    s_burnMintERC20PausableUUPS.mint(STRANGER, AMOUNT);

    assertEq(s_burnMintERC20PausableUUPS.balanceOf(STRANGER), AMOUNT);
    assertEq(s_burnMintERC20PausableUUPS.totalSupply(), AMOUNT);

    // Upgrade to the new version
    changePrank(DEFAULT_UPGRADER);
    MockBurnMintERC20PausableUUPSV2 newImplementation = new MockBurnMintERC20PausableUUPSV2();

    IUpgradeableProxy(s_uupsProxy).upgradeToAndCall(
      address(newImplementation),
      abi.encodeCall(MockBurnMintERC20PausableUUPSV2.initializeFreezerRole, (defaultFreezer))
    );

    newImplementation = MockBurnMintERC20PausableUUPSV2(s_uupsProxy);

    // Validate that the new implementation keeps track of the balances from the previous version
    assertEq(newImplementation.balanceOf(STRANGER), AMOUNT);
    assertEq(newImplementation.totalSupply(), AMOUNT);

    // Validate that the new functionality works as expected
    changePrank(defaultFreezer);
    newImplementation.freeze(STRANGER);

    changePrank(STRANGER);
    vm.expectRevert(
      abi.encodeWithSelector(
        MockBurnMintERC20PausableUUPSV2.MockBurnMintERC20PausableUUPSV2__AccountFrozen.selector,
        STRANGER
      )
    );
    newImplementation.transfer(OWNER, AMOUNT);
  }
}
