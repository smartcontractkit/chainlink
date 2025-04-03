// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BaseTest} from "../../../BaseTest.t.sol";

import {IGetCCIPAdmin} from "../../../../../ccip/interfaces/IGetCCIPAdmin.sol";
import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {IERC20Metadata} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/extensions/IERC20Metadata.sol";

interface IERC20UpgradeableBase is IBurnMintERC20Upgradeable, IERC20Metadata, IGetCCIPAdmin {
  event CCIPAdminTransferred(address indexed previousAdmin, address indexed newAdmin);
  event AccountFrozen(address indexed account);
  event AccountUnfrozen(address indexed account);

  function maxSupply() external view returns (uint256);
  function grantMintAndBurnRoles(address burnAndMinter) external;
  function setCCIPAdmin(address newAdmin) external;

  function pause() external;
  function unpause() external;

  function freeze(address account) external;
  function unfreeze(address account) external;
  function isFrozen(address account) external view returns (bool);
}

contract ERC20UpgradableBaseTest is BaseTest {
  string internal constant NAME = "CCIP-BnM Upgradeable";
  string internal constant SYMBOL = "CCIP-BnM";
  uint8 internal constant DECIMALS = 18;
  uint256 internal constant MAX_SUPPLY = 1e27;
  uint256 internal constant PRE_MINT = 0;
  address internal constant DEFAULT_ADMIN = OWNER;
  address internal constant DEFAULT_UPGRADER = OWNER;
  address internal constant DEFAULT_PAUSER = OWNER;
  address internal constant INITIAL_OWNER_ADDRESS_FOR_PROXY_ADMIN = OWNER;
  address internal constant DEFAULT_FREEZER = OWNER;

  uint256 internal constant AMOUNT = 1e18;
  address internal immutable i_mockPool = makeAddr("i_mockPool");

  function setUp() public virtual override {
    BaseTest.setUp();
  }
}
