// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

// Imports to any non-library are not allowed due to the significant cascading
// compile time increase they cause when imported into this base test.
import {IRMNRemote} from "../interfaces/IRMNRemote.sol";

import {Router} from "../Router.sol";
import {Internal} from "../libraries/Internal.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";

import {WETH9} from "../../vendor/canonical-weth/WETH9.sol";
import {Test} from "forge-std/Test.sol";

contract BaseTest is Test {
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant STRANGER = address(999999);

  // Timing
  uint256 internal constant BLOCK_TIME = 1234567890;
  uint32 internal constant TWELVE_HOURS = 60 * 60 * 12;

  // Message info
  uint64 internal constant SOURCE_CHAIN_SELECTOR = 1;
  uint64 internal constant DEST_CHAIN_SELECTOR = 2;
  uint32 internal constant GAS_LIMIT = 200_000;

  uint32 internal constant DEFAULT_TOKEN_DEST_GAS_OVERHEAD = 90_000;
  uint8 internal constant DEFAULT_TOKEN_DECIMALS = 18;
  uint16 internal constant GAS_FOR_CALL_EXACT_CHECK = 5_000;

  bool private s_baseTestInitialized;

  IRMNRemote internal s_mockRMNRemote;
  Router internal s_sourceRouter;
  Router internal s_destRouter;

  function setUp() public virtual {
    // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;

    // Set the sender to OWNER permanently
    vm.startPrank(OWNER);
    deal(OWNER, 1e20);
    vm.label(OWNER, "Owner");
    vm.label(STRANGER, "Stranger");

    // Set the block time to a constant known value
    vm.warp(BLOCK_TIME);

    // setup RMNRemote
    s_mockRMNRemote = IRMNRemote(makeAddr("MOCK RMN REMOTE"));
    vm.etch(address(s_mockRMNRemote), bytes("fake bytecode"));
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed()"), abi.encode(false));
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed(bytes16)"), abi.encode(false)); // no curses by defaule

    s_sourceRouter = new Router(address(new WETH9()), address(s_mockRMNRemote));
    vm.label(address(s_sourceRouter), "sourceRouter");
    // Deploy a destination router
    s_destRouter = new Router(address(new WETH9()), address(s_mockRMNRemote));
    vm.label(address(s_destRouter), "destRouter");
  }

  function _setMockRMNChainCurse(uint64 chainSelector, bool isCursed) internal {
    vm.mockCall(
      address(s_mockRMNRemote),
      abi.encodeWithSignature("isCursed(bytes16)", bytes16(uint128(chainSelector))),
      abi.encode(isCursed)
    );
  }

  function _getOutboundRateLimiterConfig() internal pure returns (RateLimiter.Config memory) {
    return RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e15});
  }

  function _getInboundRateLimiterConfig() internal pure returns (RateLimiter.Config memory) {
    return RateLimiter.Config({isEnabled: true, capacity: 222e30, rate: 1e18});
  }

  function _getSingleTokenPriceUpdateStruct(
    address token,
    uint224 price
  ) internal pure returns (Internal.PriceUpdates memory) {
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](1);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: token, usdPerToken: price});

    Internal.PriceUpdates memory priceUpdates =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: new Internal.GasPriceUpdate[](0)});

    return priceUpdates;
  }

  function _generateSourceTokenData() internal pure returns (Internal.SourceTokenData memory) {
    return Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(address(12312412312)),
      destTokenAddress: abi.encode(address(9809808909)),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });
  }
}
