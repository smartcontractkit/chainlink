// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

contract BurnToAddressMintTokenPool_releaseOrMint is BurnToAddressMintTokenPoolSetup {
  function test_releaseOrMint() public {
    uint256 amount = 1e24;
    address receiver = makeAddr("RECEIVER_ADDRESS");

    vm.startPrank(s_burnMintOffRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), receiver, amount);

    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: receiver,
        amount: amount,
        localToken: address(s_burnMintERC20),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    assertEq(s_burnMintERC20.balanceOf(receiver), amount);
    assertEq(s_pool.getOutstandingTokens(), amount);
  }
}
