// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IGetCCIPAdmin} from "../../../../ccip/interfaces/IGetCCIPAdmin.sol";
import {IBurnMintERC20Upgradeable} from "../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";

import {AccessControlDefaultAdminRulesUpgradeable} from "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/access/extensions/AccessControlDefaultAdminRulesUpgradeable.sol";
import {Initializable} from "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/proxy/utils/Initializable.sol";
import {ERC20BurnableUpgradeable} from "../../../../vendor/openzeppelin-solidity-upgradeable/v5.0.2/contracts/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";
import {IAccessControl} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

contract BurnMintERC20Transparent is
  Initializable,
  IBurnMintERC20Upgradeable,
  IGetCCIPAdmin,
  IERC165,
  ERC20BurnableUpgradeable,
  AccessControlDefaultAdminRulesUpgradeable
{
  error BurnMintERC20Transparent__MaxSupplyExceeded(uint256 supplyAfterMint);
  error BurnMintERC20Transparent__InvalidRecipient(address recipient);

  event CCIPAdminTransferred(address indexed previousAdmin, address indexed newAdmin);

  bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
  bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");

  // ================================================================
  // │                         Storage                              │
  // ================================================================

  /// @custom:storage-location erc7201:chainlink.storage.BurnMintERC20Transparent
  struct BurnMintERC20TransparentStorage {
    /// @dev the CCIPAdmin can be used to register with the CCIP token admin registry, but has no other special powers, and can only be transferred by the owner.
    address ccipAdmin;
    /// @dev The number of decimals for the token
    uint8 decimals;
    /// @dev The maximum supply of the token, 0 if unlimited
    uint256 maxSupply;
  }

  // keccak256(abi.encode(uint256(keccak256("chainlink.storage.BurnMintERC20Transparent")) - 1)) & ~bytes32(uint256(0xff));
  bytes32 private constant BURN_MINT_ERC20_TRANSPARENT_STORAGE_LOCATION =
    0xc5ce4c6194754ec56151469c4af5ff17dd2a95dab96bf61ba95b3ff079048900;

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getBurnMintERC20TransparentStorage() private pure returns (BurnMintERC20TransparentStorage storage $) {
    assembly {
      $.slot := BURN_MINT_ERC20_TRANSPARENT_STORAGE_LOCATION
    }
  }

  // ================================================================
  // │                         Transparent                          │
  // ================================================================

  /// @custom:oz-upgrades-unsafe-allow constructor
  constructor() {
    _disableInitializers();
  }

  /// @dev the underscores in parameter names are used to suppress compiler warnings about shadowing ERC20 functions
  function initialize(
    string memory name,
    string memory symbol,
    uint8 decimals_,
    uint256 maxSupply_,
    uint256 preMint,
    address defaultAdmin
  ) public initializer {
    __ERC20_init(name, symbol);
    __ERC20Burnable_init();
    __AccessControl_init();

    BurnMintERC20TransparentStorage storage $ = _getBurnMintERC20TransparentStorage();

    $.decimals = decimals_;
    $.maxSupply = maxSupply_;

    $.ccipAdmin = defaultAdmin;

    if (preMint != 0) {
      if (preMint > maxSupply_) {
        revert BurnMintERC20Transparent__MaxSupplyExceeded(preMint);
      }
      _mint(defaultAdmin, preMint);
    }

    _grantRole(DEFAULT_ADMIN_ROLE, defaultAdmin);
  }

  // ================================================================
  // │                           ERC165                             │
  // ================================================================

  /// @inheritdoc IERC165
  function supportsInterface(
    bytes4 interfaceId
  ) public pure virtual override(AccessControlDefaultAdminRulesUpgradeable, IERC165) returns (bool) {
    return
      interfaceId == type(IERC20).interfaceId ||
      interfaceId == type(IBurnMintERC20Upgradeable).interfaceId ||
      interfaceId == type(IERC165).interfaceId ||
      interfaceId == type(IAccessControl).interfaceId ||
      interfaceId == type(IGetCCIPAdmin).interfaceId;
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  /// @dev Returns the number of decimals used in its user representation.
  function decimals() public view virtual override returns (uint8) {
    BurnMintERC20TransparentStorage storage $ = _getBurnMintERC20TransparentStorage();
    return $.decimals;
  }

  /// @dev Returns the max supply of the token, 0 if unlimited.
  function maxSupply() public view virtual returns (uint256) {
    BurnMintERC20TransparentStorage storage $ = _getBurnMintERC20TransparentStorage();
    return $.maxSupply;
  }

  /// @dev Disallows minting and transferring to address(this).
  function _update(address from, address to, uint256 value) internal virtual override {
    if (to == address(this)) revert BurnMintERC20Transparent__InvalidRecipient(to);

    super._update(from, to, value);
  }

  /// @dev Disallows approving for address(this)
  function _approve(address owner, address spender, uint256 value, bool emitEvent) internal virtual override {
    if (spender == address(this)) revert BurnMintERC20Transparent__InvalidRecipient(spender);

    super._approve(owner, spender, value, emitEvent);
  }

  // ================================================================
  // │                      Burning & minting                       │
  // ================================================================

  /// @inheritdoc ERC20BurnableUpgradeable
  /// @dev Uses OZ ERC20Upgradeable _burn to disallow burning from address(0).
  /// @dev Decreases the total supply.
  function burn(
    uint256 amount
  ) public override(IBurnMintERC20Upgradeable, ERC20BurnableUpgradeable) onlyRole(BURNER_ROLE) {
    super.burn(amount);
  }

  /// @inheritdoc IBurnMintERC20Upgradeable
  /// @dev Alias for BurnFrom for compatibility with the older naming convention.
  /// @dev Uses burnFrom for all validation & logic.
  function burn(address account, uint256 amount) public virtual override {
    burnFrom(account, amount);
  }

  /// @inheritdoc ERC20BurnableUpgradeable
  /// @dev Uses OZ ERC20Upgradeable _burn to disallow burning from address(0).
  /// @dev Decreases the total supply.
  function burnFrom(
    address account,
    uint256 amount
  ) public override(IBurnMintERC20Upgradeable, ERC20BurnableUpgradeable) onlyRole(BURNER_ROLE) {
    super.burnFrom(account, amount);
  }

  /// @inheritdoc IBurnMintERC20Upgradeable
  /// @dev Uses OZ ERC20Upgradeable _mint to disallow minting to address(0).
  /// @dev Disallows minting to address(this) via _beforeTokenTransfer hook.
  /// @dev Increases the total supply.
  function mint(address account, uint256 amount) external override onlyRole(MINTER_ROLE) {
    BurnMintERC20TransparentStorage storage $ = _getBurnMintERC20TransparentStorage();
    uint256 _maxSupply = $.maxSupply;
    uint256 _totalSupply = totalSupply();

    if (_maxSupply != 0 && _totalSupply + amount > _maxSupply) {
      revert BurnMintERC20Transparent__MaxSupplyExceeded(_totalSupply + amount);
    }

    _mint(account, amount);
  }

  // ================================================================
  // │                            Roles                             │
  // ================================================================

  /// @notice grants both mint and burn roles to `burnAndMinter`.
  /// @dev calls public functions so this function does not require
  /// access controls. This is handled in the inner functions.
  function grantMintAndBurnRoles(address burnAndMinter) external {
    grantRole(MINTER_ROLE, burnAndMinter);
    grantRole(BURNER_ROLE, burnAndMinter);
  }

  /// @notice Returns the current CCIPAdmin
  function getCCIPAdmin() external view returns (address) {
    BurnMintERC20TransparentStorage storage $ = _getBurnMintERC20TransparentStorage();
    return $.ccipAdmin;
  }

  /// @notice Transfers the CCIPAdmin role to a new address
  /// @dev only the owner can call this function, NOT the current ccipAdmin, and 1-step ownership transfer is used.
  /// @param newAdmin The address to transfer the CCIPAdmin role to. Setting to address(0) is a valid way to revoke
  /// the role
  function setCCIPAdmin(address newAdmin) external onlyRole(DEFAULT_ADMIN_ROLE) {
    BurnMintERC20TransparentStorage storage $ = _getBurnMintERC20TransparentStorage();
    address currentAdmin = $.ccipAdmin;

    $.ccipAdmin = newAdmin;

    emit CCIPAdminTransferred(currentAdmin, newAdmin);
  }
}
