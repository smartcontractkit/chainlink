// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {BurnMintERC677} from "../../shared/token/ERC677/BurnMintERC677.sol";
import {Client} from "../libraries/Client.sol";
import {BurnMintTokenPool} from "../pools/BurnMintTokenPool.sol";
import {LockReleaseTokenPool} from "../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../pools/TokenPool.sol";
import {TokenAdminRegistry} from "../tokenAdminRegistry/TokenAdminRegistry.sol";
import {MaybeRevertingBurnMintTokenPool} from "./helpers/MaybeRevertingBurnMintTokenPool.sol";
import {RouterSetup} from "./router/RouterSetup.t.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenSetup is RouterSetup {
  address[] internal s_sourceTokens;
  address[] internal s_destTokens;

  address internal s_sourceFeeToken;
  address internal s_destFeeToken;

  TokenAdminRegistry internal s_tokenAdminRegistry;

  mapping(address sourceToken => address sourcePool) internal s_sourcePoolByToken;
  mapping(address sourceToken => address destPool) internal s_destPoolBySourceToken;
  mapping(address destToken => address destPool) internal s_destPoolByToken;
  mapping(address sourceToken => address destToken) internal s_destTokenBySourceToken;

  function _deploySourceToken(string memory tokenName, uint256 dealAmount, uint8 decimals) internal returns (address) {
    BurnMintERC677 token = new BurnMintERC677(tokenName, tokenName, decimals, 0);
    s_sourceTokens.push(address(token));
    deal(address(token), OWNER, dealAmount);
    return address(token);
  }

  function _deployDestToken(string memory tokenName, uint256 dealAmount) internal returns (address) {
    BurnMintERC677 token = new BurnMintERC677(tokenName, tokenName, 18, 0);
    s_destTokens.push(address(token));
    deal(address(token), OWNER, dealAmount);
    return address(token);
  }

  function _deployLockReleasePool(address token, bool isSourcePool) internal {
    address router = address(s_sourceRouter);
    if (!isSourcePool) {
      router = address(s_destRouter);
    }

    LockReleaseTokenPool pool =
      new LockReleaseTokenPool(IERC20(token), new address[](0), address(s_mockRMN), true, router);

    if (isSourcePool) {
      s_sourcePoolByToken[address(token)] = address(pool);
    } else {
      s_destPoolByToken[address(token)] = address(pool);
      s_destPoolBySourceToken[s_sourceTokens[s_destTokens.length - 1]] = address(pool);
    }
  }

  function _deployTokenAndBurnMintPool(address token, bool isSourcePool) internal {
    address router = address(s_sourceRouter);
    if (!isSourcePool) {
      router = address(s_destRouter);
    }

    BurnMintTokenPool pool =
      new MaybeRevertingBurnMintTokenPool(BurnMintERC677(token), new address[](0), address(s_mockRMN), router);
    BurnMintERC677(token).grantMintAndBurnRoles(address(pool));

    if (isSourcePool) {
      s_sourcePoolByToken[address(token)] = address(pool);
    } else {
      s_destPoolByToken[address(token)] = address(pool);
      s_destPoolBySourceToken[s_sourceTokens[s_destTokens.length - 1]] = address(pool);
    }
  }

  function setUp() public virtual override {
    RouterSetup.setUp();

    bool isSetup = s_sourceTokens.length != 0;
    if (isSetup) {
      return;
    }

    // Source tokens & pools
    address sourceLink = _deploySourceToken("sLINK", type(uint256).max, 18);
    _deployLockReleasePool(sourceLink, true);
    s_sourceFeeToken = sourceLink;

    address sourceEth = _deploySourceToken("sETH", 2 ** 128, 18);
    _deployTokenAndBurnMintPool(sourceEth, true);

    // Destination tokens & pools
    address destLink = _deployDestToken("dLINK", type(uint256).max);
    _deployLockReleasePool(destLink, false);
    s_destFeeToken = destLink;

    s_destTokenBySourceToken[sourceLink] = destLink;

    address destEth = _deployDestToken("dETH", 2 ** 128);
    _deployTokenAndBurnMintPool(destEth, false);

    s_destTokenBySourceToken[sourceEth] = destEth;

    // Float the dest link lock release pool with funds
    IERC20(destLink).transfer(s_destPoolByToken[destLink], 1000 ether);

    s_tokenAdminRegistry = new TokenAdminRegistry();

    // Set pools in the registry
    for (uint256 i = 0; i < s_sourceTokens.length; ++i) {
      address token = s_sourceTokens[i];
      address pool = s_sourcePoolByToken[token];

      _setPool(
        s_tokenAdminRegistry, token, pool, DEST_CHAIN_SELECTOR, s_destPoolByToken[s_destTokens[i]], s_destTokens[i]
      );
    }

    for (uint256 i = 0; i < s_destTokens.length; ++i) {
      address token = s_destTokens[i];
      address pool = s_destPoolByToken[token];
      s_tokenAdminRegistry.proposeAdministrator(token, OWNER);
      s_tokenAdminRegistry.acceptAdminRole(token);
      s_tokenAdminRegistry.setPool(token, pool);

      _setPool(
        s_tokenAdminRegistry,
        token,
        pool,
        SOURCE_CHAIN_SELECTOR,
        s_sourcePoolByToken[s_sourceTokens[i]],
        s_sourceTokens[i]
      );
    }
  }

  function _getCastedSourceEVMTokenAmountsWithZeroAmounts()
    internal
    view
    returns (Client.EVMTokenAmount[] memory tokenAmounts)
  {
    tokenAmounts = new Client.EVMTokenAmount[](s_sourceTokens.length);
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].token = s_sourceTokens[i];
    }
    return tokenAmounts;
  }

  function _setPool(
    TokenAdminRegistry tokenAdminRegistry,
    address token,
    address pool,
    uint64 remoteChainSelector,
    address remotePoolAddress,
    address remoteToken
  ) internal {
    if (!tokenAdminRegistry.isAdministrator(token, OWNER)) {
      tokenAdminRegistry.proposeAdministrator(token, OWNER);
      tokenAdminRegistry.acceptAdminRole(token);
    }

    tokenAdminRegistry.setPool(token, pool);

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: remoteChainSelector,
      remotePoolAddress: abi.encode(remotePoolAddress),
      remoteTokenAddress: abi.encode(remoteToken),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });

    TokenPool(pool).applyChainUpdates(chainUpdates);
  }
}
