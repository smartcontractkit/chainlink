// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {BurnMintERC20} from "../../../../shared/token/ERC20/BurnMintERC20.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenPoolSetup is BaseTest {
  IERC20 internal s_token;
  TokenPoolHelper internal s_tokenPool;

  address internal s_initialRemotePool = makeAddr("initialRemotePool");
  address internal s_initialRemoteToken = makeAddr("initialRemoteToken");

  function setUp() public virtual override {
    super.setUp();
    s_token = new BurnMintERC20("LINK", "LNK", 18, 0, 0);
    deal(address(s_token), OWNER, type(uint256).max);

    s_tokenPool = new TokenPoolHelper(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );

    bytes[] memory remotePoolAddresses = new bytes[](1);
    remotePoolAddresses[0] = abi.encode(s_initialRemotePool);

    TokenPool.ChainUpdate[] memory chainUpdate = new TokenPool.ChainUpdate[](1);
    chainUpdate[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddresses: remotePoolAddresses,
      remoteTokenAddress: abi.encode(s_initialRemoteToken),
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_tokenPool.applyChainUpdates(new uint64[](0), chainUpdate);
  }
}
