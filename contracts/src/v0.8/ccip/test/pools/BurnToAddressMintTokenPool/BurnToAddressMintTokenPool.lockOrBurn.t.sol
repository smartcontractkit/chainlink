// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

contract BurnToAddressMintTokenPool_lockOrBurn is BurnToAddressMintTokenPoolSetup {
  uint256 public constant AMOUNT = 1e24;

  function test_LockOrBurn() public {
    deal(address(s_burnMintERC20), address(s_pool), AMOUNT);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), AMOUNT);

    vm.startPrank(s_burnMintOnRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_pool), BURN_ADDRESS, AMOUNT);

    vm.expectCall(address(s_burnMintERC20), abi.encodeWithSelector(IERC20.transfer.selector, BURN_ADDRESS, AMOUNT));

    s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: bytes(""),
        amount: AMOUNT,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_burnMintERC20)
      })
    );

    assertEq(s_burnMintERC20.balanceOf(s_pool.getBurnAddress()), AMOUNT);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), 0);
  }
}
