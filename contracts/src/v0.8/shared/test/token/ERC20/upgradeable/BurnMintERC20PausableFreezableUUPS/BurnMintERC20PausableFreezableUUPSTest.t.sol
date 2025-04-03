// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ERC1967Proxy} from "../../../../../../vendor/openzeppelin-solidity/v5.0.2/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {BurnMintERC20PausableFreezableUUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20PausableFreezableUUPS.sol";
import {BurnMintERC20UUPS} from "../../../../../token/ERC20/upgradeable/BurnMintERC20UUPS.sol";
import {ERC20UpgradableBaseTest_freeze} from "../ERC20UpgradableBaseTest.freeze.t.sol";
import {ERC20UpgradableBaseTest_unfreeze} from "../ERC20UpgradableBaseTest.unfreeze.t.sol";

contract BurnMintERC20PausableFreezableUUPSTest is ERC20UpgradableBaseTest_freeze, ERC20UpgradableBaseTest_unfreeze {
  BurnMintERC20PausableFreezableUUPS internal s_burnMintERC20PausableFreezableUUPS;

  function setUp() public virtual override {
    address implementation = address(new BurnMintERC20PausableFreezableUUPS());

    address proxy = address(
      new ERC1967Proxy(
        implementation,
        abi.encodeCall(
          BurnMintERC20UUPS.initialize,
          (NAME, SYMBOL, DECIMALS, MAX_SUPPLY, PRE_MINT, DEFAULT_ADMIN, DEFAULT_UPGRADER)
        )
      )
    );

    s_burnMintERC20PausableFreezableUUPS = BurnMintERC20PausableFreezableUUPS(proxy);

    changePrank(DEFAULT_ADMIN);
    s_burnMintERC20PausableFreezableUUPS.grantRole(s_burnMintERC20PausableFreezableUUPS.PAUSER_ROLE(), DEFAULT_PAUSER);
    s_burnMintERC20PausableFreezableUUPS.grantRole(
      s_burnMintERC20PausableFreezableUUPS.FREEZER_ROLE(),
      DEFAULT_FREEZER
    );
  }

  // ================================================================
  // │                           Freeze                             │
  // ================================================================

  function test_Freeze() public {
    should_Freeze(address(s_burnMintERC20PausableFreezableUUPS));
  }

  function test_Freeze_EvenWhenImplementationIsPaused() public {
    should_Freeze_EvenWhenImplementationIsPaused(address(s_burnMintERC20PausableFreezableUUPS));
  }

  function test_Freeze_RevertWhen_CallerDoesNotHaveFreezerRole() public {
    should_Freeze_RevertWhen_CallerDoesNotHaveFreezerRole(
      address(s_burnMintERC20PausableFreezableUUPS),
      s_burnMintERC20PausableFreezableUUPS.FREEZER_ROLE()
    );
  }

  function test_Freeze_RevertWhen_RecipientIsAddressZero() public {
    should_Freeze_RevertWhen_RecipientIsAddressZero(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__InvalidRecipient.selector
    );
  }

  function test_Freeze_RevertWhen_RecipientIsImplementationItself() public {
    should_Freeze_RevertWhen_RecipientIsImplementationItself(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__InvalidRecipient.selector
    );
  }

  function test_Freeze_RevertWhen_AccountIsAlreadyFrozen() public {
    should_Freeze_RevertWhen_AccountIsAlreadyFrozen(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountFrozen.selector
    );
  }

  function test_Mint_RevertWhen_AccountIsFrozen() public {
    should_Mint_RevertWhen_AccountIsFrozen(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountFrozen.selector
    );
  }

  function test_Transfer_RevertWhen_SenderOrRecipientAreFrozen() public {
    should_Transfer_RevertWhen_SenderOrRecipientAreFrozen(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountFrozen.selector
    );
  }

  function test_Approve_RevertWhen_OwnerOrSpenderAreFrozen() public {
    should_Approve_RevertWhen_OwnerOrSpenderAreFrozen(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountFrozen.selector
    );
  }

  // ================================================================
  // │                          Unfreeze                            │
  // ================================================================

  function test_Unfreeze() public {
    should_Unfreeze(address(s_burnMintERC20PausableFreezableUUPS));
  }

  function test_Unfreeze_EvenWhenImplementationIsPaused() public {
    should_Unfreeze_EvenWhenImplementationIsPaused(address(s_burnMintERC20PausableFreezableUUPS));
  }

  function test_Unfreeze_RevertWhen_CallerDoesNotHaveFreezerRole() public {
    should_Unfreeze_RevertWhen_CallerDoesNotHaveFreezerRole(
      address(s_burnMintERC20PausableFreezableUUPS),
      s_burnMintERC20PausableFreezableUUPS.FREEZER_ROLE()
    );
  }

  function test_Unfreeze_RevertWhen_AccountIsNotFrozen() public {
    should_Unfreeze_RevertWhen_AccountIsNotFrozen(
      address(s_burnMintERC20PausableFreezableUUPS),
      BurnMintERC20PausableFreezableUUPS.BurnMintERC20PausableFreezableUUPS__AccountNotFrozen.selector
    );
  }
}
