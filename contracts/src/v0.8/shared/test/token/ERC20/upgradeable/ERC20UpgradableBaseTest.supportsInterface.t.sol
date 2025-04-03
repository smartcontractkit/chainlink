// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IGetCCIPAdmin} from "../../../../../ccip/interfaces/IGetCCIPAdmin.sol";
import {IBurnMintERC20Upgradeable} from "../../../../../shared/token/ERC20/upgradeable/IBurnMintERC20Upgradeable.sol";
import {IAccessControl} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/access/IAccessControl.sol";
import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC20.sol";
import {IERC165} from "../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";

import {ERC20UpgradableBaseTest} from "./ERC20UpgradableBaseTest.t.sol";

contract ERC20UpgradableBaseTest_supportsInterface is ERC20UpgradableBaseTest {
  function should_SupportsInterface(IERC165 implementation) public view {
    assertTrue(implementation.supportsInterface(type(IERC20).interfaceId));
    assertTrue(implementation.supportsInterface(type(IBurnMintERC20Upgradeable).interfaceId));
    assertTrue(implementation.supportsInterface(type(IERC165).interfaceId));
    assertTrue(implementation.supportsInterface(type(IAccessControl).interfaceId));
    assertTrue(implementation.supportsInterface(type(IGetCCIPAdmin).interfaceId));
  }
}
