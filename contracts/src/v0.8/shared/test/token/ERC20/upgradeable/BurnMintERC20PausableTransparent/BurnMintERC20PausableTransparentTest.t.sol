// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {TransparentUpgradeableProxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {BurnMintERC20PausableTransparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableTransparent.sol";
import {BurnMintERC20Transparent} from "../../../../../token/ERC20/upgradeable/BurnMintERC20Transparent.sol";
import {ERC20UpgradableBaseTest_pausing} from "../ERC20UpgradableBaseTest.pausing.t.sol";

contract BurnMintERC20PausableTransparentTest is ERC20UpgradableBaseTest_pausing {
  BurnMintERC20PausableTransparent internal s_burnMintERC20PausableTransparent;

  function setUp() public virtual override {
    address implementation = address(new BurnMintERC20PausableTransparent());

    address proxy = address(
      new TransparentUpgradeableProxy(
        implementation,
        INITIAL_OWNER_ADDRESS_FOR_PROXY_ADMIN,
        abi.encodeCall(
          BurnMintERC20Transparent.initialize,
          (NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN)
        )
      )
    );

    s_burnMintERC20PausableTransparent = BurnMintERC20PausableTransparent(proxy);

    changePrank(DEFAULT_ADMIN);
    s_burnMintERC20PausableTransparent.grantRole(s_burnMintERC20PausableTransparent.PAUSER_ROLE(), DEFAULT_PAUSER);
    s_burnMintERC20PausableTransparent.grantMintAndBurnRoles(i_mockPool);
  }

  // ================================================================
  // │                          Pausing                             │
  // ================================================================

  function test_Pause() public {
    should_Pause(address(s_burnMintERC20PausableTransparent));
  }

  function test_Unpause() public {
    should_Unpause(address(s_burnMintERC20PausableTransparent));
  }

  function test_Pause_RevertWhen_CallerDoesNotHavePauserRole() public {
    should_Pause_RevertWhen_CallerDoesNotHavePauserRole(
      address(s_burnMintERC20PausableTransparent),
      s_burnMintERC20PausableTransparent.PAUSER_ROLE()
    );
  }

  function test_Unpause_RevertWhen_CallerDoesNotHaveDefaultAdminRole() public {
    should_Unpause_RevertWhen_CallerDoesNotHaveDefaultAdminRole(
      address(s_burnMintERC20PausableTransparent),
      s_burnMintERC20PausableTransparent.DEFAULT_ADMIN_ROLE()
    );
  }

  // ================================================================
  // │                      Burning & minting                       │
  // ================================================================

  function test_Mint_RevertWhen_ImplementationIsPaused() public {
    should_Mint_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableTransparent));
  }

  function test_Burn_RevertWhen_ImplementationIsPaused() public {
    should_Burn_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableTransparent));
  }

  function test_BurnFrom_RevertWhen_ImplementationIsPaused() public {
    should_BurnFrom_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableTransparent));
  }

  function test_BurnFrom_alias_RevertWhen_ImplementationIsPaused() public {
    should_BurnFrom_alias_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableTransparent));
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  function test_Transfer_RevertWhen_ImplementationIsPaused() public {
    should_Transfer_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableTransparent));
  }

  function test_Approve_RevertWhen_ImplementationIsPaused() public {
    should_Approve_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableTransparent));
  }
}
