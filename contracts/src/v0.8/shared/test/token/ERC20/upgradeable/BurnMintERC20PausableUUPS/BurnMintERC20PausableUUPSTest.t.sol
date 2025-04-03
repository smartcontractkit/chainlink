// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC1967Proxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {BurnMintERC20PausableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableUUPS.sol";
import {BurnMintERC20UUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {ERC20UpgradableBaseTest_pausing} from "../ERC20UpgradableBaseTest.pausing.t.sol";

contract BurnMintERC20PausableUUPSTest is ERC20UpgradableBaseTest_pausing {
  BurnMintERC20PausableUUPS internal s_burnMintERC20PausableUUPS;

  function setUp() public virtual override {
    address implementation = address(new BurnMintERC20PausableUUPS());

    address proxy = address(
      new ERC1967Proxy(
        implementation,
        abi.encodeCall(
          BurnMintERC20UUPS.initialize,
          (NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN, DEFAULT_UPGRADER)
        )
      )
    );

    s_burnMintERC20PausableUUPS = BurnMintERC20PausableUUPS(proxy);

    changePrank(DEFAULT_ADMIN);
    s_burnMintERC20PausableUUPS.grantRole(s_burnMintERC20PausableUUPS.PAUSER_ROLE(), DEFAULT_PAUSER);
    s_burnMintERC20PausableUUPS.grantMintAndBurnRoles(i_mockPool);
  }

  // ================================================================
  // │                          Pausing                             │
  // ================================================================

  function test_Pause() public {
    should_Pause(address(s_burnMintERC20PausableUUPS));
  }

  function test_Unpause() public {
    should_Unpause(address(s_burnMintERC20PausableUUPS));
  }

  function test_Pause_RevertWhen_CallerDoesNotHavePauserRole() public {
    should_Pause_RevertWhen_CallerDoesNotHavePauserRole(
      address(s_burnMintERC20PausableUUPS),
      s_burnMintERC20PausableUUPS.PAUSER_ROLE()
    );
  }

  function test_Unpause_RevertWhen_CallerDoesNotHaveDefaultAdminRole() public {
    should_Unpause_RevertWhen_CallerDoesNotHaveDefaultAdminRole(
      address(s_burnMintERC20PausableUUPS),
      s_burnMintERC20PausableUUPS.DEFAULT_ADMIN_ROLE()
    );
  }

  // ================================================================
  // │                      Burning & minting                       │
  // ================================================================

  function test_Mint_RevertWhen_ImplementationIsPaused() public {
    should_Mint_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableUUPS));
  }

  function test_Burn_RevertWhen_ImplementationIsPaused() public {
    should_Burn_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableUUPS));
  }

  function test_BurnFrom_RevertWhen_ImplementationIsPaused() public {
    should_BurnFrom_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableUUPS));
  }

  function test_BurnFrom_alias_RevertWhen_ImplementationIsPaused() public {
    should_BurnFrom_alias_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableUUPS));
  }

  // ================================================================
  // │                            ERC20                             │
  // ================================================================

  function test_Transfer_RevertWhen_ImplementationIsPaused() public {
    should_Transfer_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableUUPS));
  }

  function test_Approve_RevertWhen_ImplementationIsPaused() public {
    should_Approve_RevertWhen_ImplementationIsPaused(address(s_burnMintERC20PausableUUPS));
  }
}
