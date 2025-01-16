// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {BurnMintERC20} from "../../../../shared/token/ERC20/BurnMintERC20.sol";
import {Router} from "../../../Router.sol";
import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {BaseTest} from "../../BaseTest.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract LockReleaseTokenPoolSetup is BaseTest {
  IERC20 internal s_token;
  LockReleaseTokenPool internal s_lockReleaseTokenPool;
  LockReleaseTokenPool internal s_lockReleaseTokenPoolWithAllowList;
  address[] internal s_allowedList;

  address internal s_allowedOnRamp = address(123);
  address internal s_allowedOffRamp = address(234);

  address internal s_destPoolAddress = address(2736782345);
  address internal s_sourcePoolAddress = address(53852352095);

  function setUp() public virtual override {
    super.setUp();
    s_token = new BurnMintERC20("LINK", "LNK", 18, 0, 0);
    deal(address(s_token), OWNER, type(uint256).max);
    s_lockReleaseTokenPool = new LockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), true, address(s_sourceRouter)
    );

    s_allowedList.push(vm.randomAddress());
    s_allowedList.push(OWNER);
    s_lockReleaseTokenPoolWithAllowList = new LockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, s_allowedList, address(s_mockRMNRemote), true, address(s_sourceRouter)
    );

    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = abi.encode(s_destPoolAddress);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: abi.encode(address(2)),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    s_lockReleaseTokenPool.applyChainUpdates(new uint64[](0), chainUpdate);
    s_lockReleaseTokenPoolWithAllowList.applyChainUpdates(new uint64[](0), chainUpdate);
    s_lockReleaseTokenPool.setRebalancer(OWNER);

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: s_allowedOnRamp});
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: s_allowedOffRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
  }
}
