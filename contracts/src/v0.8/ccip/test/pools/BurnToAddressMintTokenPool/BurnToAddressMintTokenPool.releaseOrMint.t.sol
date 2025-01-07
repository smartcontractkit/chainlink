// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

contract BurnToAddressMintTokenPool_releaseOrMint is BurnToAddressMintTokenPoolSetup {
  function test_releaseOrMint() public {
    uint256 amount = 1e19;
    address receiver = makeAddr("receiver_address");

    deal(address(s_burnMintERC20), address(s_pool), amount);

    vm.startPrank(s_burnMintOnRamp);

    // Lock some tokens
    s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_burnMintERC20)
      })
    );

    uint256 lockedAmountBefore = s_pool.getLockedTokens();

    vm.startPrank(s_burnMintOffRamp);

    uint256 releaseAmount = amount / 2;

    vm.expectEmit();
    emit IERC20.Transfer(address(0), receiver, releaseAmount);

    // relwaese
    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: receiver,
        amount: releaseAmount,
        localToken: address(s_burnMintERC20),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    assertEq(s_burnMintERC20.balanceOf(receiver), releaseAmount);
    assertEq(s_pool.getLockedTokens(), lockedAmountBefore - releaseAmount);
  }

  function test_ReleaseOrMint_RevertWhen_LockedTokensUnderflows() public {
    uint256 burnAmount = s_initialTokenAmount + 1;

    deal(address(s_burnMintERC20), address(s_pool), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    // Call should revert due to underflow error due to trying to burn more tokens than are locked via CCIP.
    vm.expectRevert(abi.encodeWithSignature("Panic(uint256)", 0x11));

    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: address(0xdead),
        amount: burnAmount,
        localToken: address(s_burnMintERC20),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }
}
