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
  address internal constant DUMMY_CONTRACT_ADDRESS = 0x1111111111111111111111111111111111111112;
  address internal constant ON_RAMP_ADDRESS = 0x11118e64e1FB0c487f25dD6D3601FF6aF8d32E4e;
  address internal constant ZERO_ADDRESS = address(0);
  address internal constant FEE_AGGREGATOR = 0xa33CDB32eAEce34F6affEfF4899cef45744EDea3;

  address internal constant USER_1 = address(1);
  address internal constant USER_2 = address(2);
  address internal constant USER_3 = address(3);
  address internal constant USER_4 = address(4);

  // Message info
  uint64 internal constant SOURCE_CHAIN_SELECTOR = 1;
  uint64 internal constant DEST_CHAIN_SELECTOR = 2;
  uint32 internal constant GAS_LIMIT = 200_000;

  // Timing
  uint256 internal constant BLOCK_TIME = 1234567890;
  uint32 internal constant TWELVE_HOURS = 60 * 60 * 12;

  // Onramp
  uint96 internal constant MAX_NOP_FEES_JUELS = 1e27;
  uint96 internal constant MAX_MSG_FEES_JUELS = 1e18;
  uint32 internal constant DEST_GAS_OVERHEAD = 300_000;
  uint16 internal constant DEST_GAS_PER_PAYLOAD_BYTE = 16;

  uint16 internal constant DEFAULT_TOKEN_FEE_USD_CENTS = 50;
  uint32 internal constant DEFAULT_TOKEN_DEST_GAS_OVERHEAD = 90_000;
  uint32 internal constant DEFAULT_TOKEN_BYTES_OVERHEAD = 32;

  bool private s_baseTestInitialized;

  // Use 16 gas per data availability byte in our tests.
  // This is an overestimation in OP stack, it ignores 4 gas per 0 byte rule.
  // Arbitrum on the other hand, does always use 16 gas per data availability byte.
  // This value may be substantially decreased after EIP 4844.
  uint16 internal constant DEST_GAS_PER_DATA_AVAILABILITY_BYTE = 16;

  // Total L1 data availability overhead estimate is 33_596 gas.
  // This value includes complete CommitStore and OffRamp call data.
  uint32 internal constant DEST_DATA_AVAILABILITY_OVERHEAD_GAS = 188 // Fixed data availability overhead in OP stack.
    + (32 * 31 + 4) * DEST_GAS_PER_DATA_AVAILABILITY_BYTE // CommitStore single-root transmission takes up about 31 slots, plus selector.
    + (32 * 34 + 4) * DEST_GAS_PER_DATA_AVAILABILITY_BYTE; // OffRamp transmission excluding EVM2EVMMessage takes up about 34 slots, plus selector.

  // Multiples of bps, or 0.0001, use 6840 to be same as OP mainnet compression factor of 0.684.
  uint16 internal constant DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS = 6840;

  // OffRamp
  uint32 internal constant MAX_DATA_SIZE = 30_000;
  uint16 internal constant MAX_TOKENS_LENGTH = 5;
  uint16 internal constant GAS_FOR_CALL_EXACT_CHECK = 5000;
  uint32 internal constant PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS = 500;
  uint32 internal constant MAX_GAS_LIMIT = 4_000_000;

  // Rate limiter
  address internal constant ADMIN = 0x11118e64e1FB0c487f25dD6D3601FF6aF8d32E4e;

  MockRMN internal s_mockRMN;
  IRMNRemote internal s_mockRMNRemote;

  // nonce for pseudo-random number generation, not to be exposed to test suites
  uint256 private randNonce;

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

  function _setMockRMNGlobalCurse(bool isCursed) internal {
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
    return keccak256(abi.encodePacked(++randNonce));
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
