// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

// Imports to any non-library are not allowed due to the significant cascading
// compile time increase they cause when imported into this base test.

import {IRMNRemote} from "../interfaces/IRMNRemote.sol";
import {Internal} from "../libraries/Internal.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {MockRMN} from "./mocks/MockRMN.sol";
import {Test} from "forge-std/Test.sol";

contract BaseTest is Test {
  // Addresses
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant STRANGER = address(999999);

  address internal constant USER_1 = address(1);

  // Message info
  uint64 internal constant SOURCE_CHAIN_SELECTOR = 1;
  uint64 internal constant DEST_CHAIN_SELECTOR = 2;
  uint32 internal constant GAS_LIMIT = 200_000;

  // Timing
  uint256 internal constant BLOCK_TIME = 1234567890;
  uint32 internal constant TWELVE_HOURS = 60 * 60 * 12;

  // Onramp
  uint96 internal constant MAX_MSG_FEES_JUELS = 1_000e18;
  uint32 internal constant DEST_GAS_OVERHEAD = 300_000;
  uint16 internal constant DEST_GAS_PER_PAYLOAD_BYTE = 16;

  uint16 internal constant DEFAULT_TOKEN_FEE_USD_CENTS = 50;
  uint32 internal constant DEFAULT_TOKEN_DEST_GAS_OVERHEAD = 90_000;
  uint32 internal constant DEFAULT_TOKEN_BYTES_OVERHEAD = 32;

  bool private s_baseTestInitialized;

  // OffRamp
  uint32 internal constant MAX_DATA_SIZE = 30_000;
  uint16 internal constant MAX_TOKENS_LENGTH = 5;
  uint16 internal constant GAS_FOR_CALL_EXACT_CHECK = 5000;
  uint32 internal constant MAX_GAS_LIMIT = 4_000_000;

  MockRMN internal s_mockRMN;
  IRMNRemote internal s_mockRMNRemote;

  // nonce for pseudo-random number generation, not to be exposed to test suites
  uint256 private s_randNonce;

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

    // setup mock RMN & RMNRemote
    s_mockRMN = new MockRMN();
    s_mockRMNRemote = IRMNRemote(makeAddr("MOCK RMN REMOTE"));
    vm.etch(address(s_mockRMNRemote), bytes("fake bytecode"));
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));
    _setMockRMNGlobalCurse(false);
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed(bytes16)"), abi.encode(false)); // no curses by defaule
  }

  function _setMockRMNGlobalCurse(
    bool isCursed
  ) internal {
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed()"), abi.encode(isCursed));
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

  /// @dev returns a pseudo-random bytes32
  function _randomBytes32() internal returns (bytes32) {
    return keccak256(abi.encodePacked(++s_randNonce));
  }

  /// @dev returns a pseudo-random number
  function _randomNum() internal returns (uint256) {
    return uint256(_randomBytes32());
  }

  /// @dev returns a pseudo-random address
  function _randomAddress() internal returns (address) {
    return address(uint160(_randomNum()));
  }
}
