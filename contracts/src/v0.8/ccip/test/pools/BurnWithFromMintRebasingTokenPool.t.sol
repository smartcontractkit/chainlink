// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {Pool} from "../../libraries/Pool.sol";
import {BurnWithFromMintRebasingTokenPool} from "../../pools/BurnWithFromMintRebasingTokenPool.sol";
import {BurnMintSetup} from "./BurnMintSetup.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {ERC20RebasingHelper} from "../helpers/ERC20RebasingHelper.sol";

contract BurnWithFromMintRebasingTokenPoolSetup is BurnMintSetup {
  BurnWithFromMintRebasingTokenPool internal s_pool;
  ERC20RebasingHelper internal s_rebasingToken;

  function setUp() public virtual override {
    BurnMintSetup.setUp();

    s_rebasingToken = new ERC20RebasingHelper();

    s_pool = new BurnWithFromMintRebasingTokenPool(
      IBurnMintERC20(address(s_rebasingToken)), new address[](0), address(s_mockRMN), address(s_sourceRouter)
    );

    _applyChainUpdates(address(s_pool));

    deal(address(s_rebasingToken), OWNER, 1e18);

    vm.startPrank(s_burnMintOffRamp);
  }
}

contract BurnWithFromMintTokenPool_releaseOrMint is BurnWithFromMintRebasingTokenPoolSetup {
  function test_Setup_Success() public view {
    assertEq(address(s_rebasingToken), address(s_pool.getToken()));
    assertEq(address(s_mockRMN), s_pool.getRmnProxy());
    assertEq(false, s_pool.getAllowListEnabled());
    assertEq(type(uint256).max, s_rebasingToken.allowance(address(s_pool), address(s_pool)));
    assertEq("BurnWithFromMintRebasingTokenPool 1.5.0", s_pool.typeAndVersion());
  }

  function test_releaseOrMint_Success() public {
    uint256 amount = 1000;
    uint256 balancePre = s_rebasingToken.balanceOf(address(OWNER));

    Pool.ReleaseOrMintOutV1 memory releaseOrMintOut = s_pool.releaseOrMint(_getReleaseOrMintIn(amount));

    assertEq(amount, releaseOrMintOut.destinationAmount);
    assertEq(balancePre + amount, s_rebasingToken.balanceOf(address(OWNER)));
  }

  function testFuzz_releaseOrMint_rebasing_success(uint16 multiplierPercentage) public {
    uint256 amount = 1000;
    uint256 expectedAmount = amount * multiplierPercentage / 100;
    s_rebasingToken.setMultiplierPercentage(multiplierPercentage);

    uint256 balancePre = s_rebasingToken.balanceOf(address(OWNER));

    Pool.ReleaseOrMintOutV1 memory releaseOrMintOut = s_pool.releaseOrMint(_getReleaseOrMintIn(amount));

    assertEq(expectedAmount, releaseOrMintOut.destinationAmount);
    assertEq(balancePre + expectedAmount, s_rebasingToken.balanceOf(address(OWNER)));
  }

  function test_releaseOrMint_NegativeMintAmount_reverts() public {
    uint256 amount = 1000;
    s_rebasingToken.setMintShouldBurn(true);

    vm.expectRevert(abi.encodeWithSelector(BurnWithFromMintRebasingTokenPool.NegativeMintAmount.selector, amount));

    s_pool.releaseOrMint(_getReleaseOrMintIn(amount));
  }

  function _getReleaseOrMintIn(uint256 amount) internal view returns (Pool.ReleaseOrMintInV1 memory) {
    return Pool.ReleaseOrMintInV1({
      originalSender: bytes(""),
      receiver: OWNER,
      amount: amount,
      localToken: address(s_rebasingToken),
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
      sourcePoolData: "",
      offchainTokenData: ""
    });
  }
}
