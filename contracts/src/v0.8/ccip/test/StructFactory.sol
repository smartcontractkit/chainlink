// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../interfaces/pools/IPool.sol";

import {ARM} from "../ARM.sol";
import {EVM2EVMOffRamp} from "../offRamp/EVM2EVMOffRamp.sol";
import {EVM2EVMOnRamp} from "../onRamp/EVM2EVMOnRamp.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {Internal} from "../libraries/Internal.sol";

contract StructFactory {
  // Addresses
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant STRANGER = address(999999);
  address internal constant DUMMY_CONTRACT_ADDRESS = 0x1111111111111111111111111111111111111112;
  address internal constant ON_RAMP_ADDRESS = 0x11118e64e1FB0c487f25dD6D3601FF6aF8d32E4e;
  address internal constant ZERO_ADDRESS = address(0);
  address internal constant BLESS_VOTER_1 = address(1);
  address internal constant CURSE_VOTER_1 = address(10);
  address internal constant CURSE_UNVOTER_1 = address(110);
  address internal constant BLESS_VOTER_2 = address(2);
  address internal constant CURSE_VOTER_2 = address(12);
  address internal constant CURSE_UNVOTER_2 = address(112);
  address internal constant BLESS_VOTER_3 = address(3);
  address internal constant CURSE_VOTER_3 = address(13);
  address internal constant CURSE_UNVOTER_3 = address(113);
  address internal constant BLESS_VOTER_4 = address(4);
  address internal constant CURSE_VOTER_4 = address(14);
  address internal constant CURSE_UNVOTER_4 = address(114);

  address internal constant USER_1 = address(1);
  address internal constant USER_2 = address(2);
  address internal constant USER_3 = address(3);
  address internal constant USER_4 = address(4);

  // Arm
  function armConstructorArgs() internal pure returns (ARM.Config memory) {
    ARM.Voter[] memory voters = new ARM.Voter[](4);
    voters[0] = ARM.Voter({
      blessVoteAddr: BLESS_VOTER_1,
      curseVoteAddr: CURSE_VOTER_1,
      curseUnvoteAddr: CURSE_UNVOTER_1,
      blessWeight: WEIGHT_1,
      curseWeight: WEIGHT_1
    });
    voters[1] = ARM.Voter({
      blessVoteAddr: BLESS_VOTER_2,
      curseVoteAddr: CURSE_VOTER_2,
      curseUnvoteAddr: CURSE_UNVOTER_2,
      blessWeight: WEIGHT_10,
      curseWeight: WEIGHT_10
    });
    voters[2] = ARM.Voter({
      blessVoteAddr: BLESS_VOTER_3,
      curseVoteAddr: CURSE_VOTER_3,
      curseUnvoteAddr: CURSE_UNVOTER_3,
      blessWeight: WEIGHT_20,
      curseWeight: WEIGHT_20
    });
    voters[3] = ARM.Voter({
      blessVoteAddr: BLESS_VOTER_4,
      curseVoteAddr: CURSE_VOTER_4,
      curseUnvoteAddr: CURSE_UNVOTER_4,
      blessWeight: WEIGHT_40,
      curseWeight: WEIGHT_40
    });
    return
      ARM.Config({
        voters: voters,
        blessWeightThreshold: WEIGHT_10 + WEIGHT_20 + WEIGHT_40,
        curseWeightThreshold: WEIGHT_1 + WEIGHT_10 + WEIGHT_20 + WEIGHT_40
      });
  }

  uint8 internal constant ZERO = 0;
  uint8 internal constant WEIGHT_1 = 1;
  uint8 internal constant WEIGHT_10 = 10;
  uint8 internal constant WEIGHT_20 = 20;
  uint8 internal constant WEIGHT_40 = 40;

  // Message info
  uint64 internal constant SOURCE_CHAIN_SELECTOR = 1;
  uint64 internal constant DEST_CHAIN_SELECTOR = 2;
  uint64 internal constant GAS_LIMIT = 200_000;

  // Timing
  uint256 internal constant BLOCK_TIME = 1234567890;
  uint32 internal constant TWELVE_HOURS = 60 * 60 * 12;

  // Onramp
  uint96 internal constant MAX_NOP_FEES_JUELS = 1e27;
  uint32 internal constant DEST_GAS_OVERHEAD = 350_000;
  uint16 internal constant DEST_GAS_PER_PAYLOAD_BYTE = 16;

  // Use 16 gas per data availability byte in our tests.
  // This is an overstimation in OP stack, it ignores 4 gas per 0 byte rule.
  // Arbitrum on the other hand, does always use 16 gas per data availability byte.
  // This value may be substantially decreased after EIP 4844.
  uint16 internal constant DEST_GAS_PER_DATA_AVAILABILITY_BYTE = 16;

  // Total L1 data availability overhead estimate is 33_596 gas.
  // This value includes complete CommitStore and OffRamp call data.
  uint32 internal constant DEST_DATA_AVAILABILITY_OVERHEAD_GAS =
    188 + // Fixed data availability overhead in OP stack.
      (32 * 31 + 4) *
      DEST_GAS_PER_DATA_AVAILABILITY_BYTE + // CommitStore single-root transmission takes up about 31 slots, plus selector.
      (32 * 34 + 4) *
      DEST_GAS_PER_DATA_AVAILABILITY_BYTE; // OffRamp transmission excluding EVM2EVMMessage takes up about 34 slots, plus selector.

  // Multiples of bps, or 0.0001, use 6840 to be same as OP mainnet compression factor of 0.684.
  uint16 internal constant DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS = 6840;

  // OffRamp
  uint256 internal constant POOL_BALANCE = 1e25;
  uint32 internal constant EXECUTION_DELAY_SECONDS = 0;
  uint32 internal constant MAX_DATA_SIZE = 30_000;
  uint16 internal constant MAX_TOKENS_LENGTH = 5;
  uint32 internal constant MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS = 200_000;
  uint16 internal constant GAS_FOR_CALL_EXACT_CHECK = 5000;
  uint32 internal constant PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS = 500;
  uint32 internal constant MAX_GAS_LIMIT = 4_000_000;

  function generateManualGasLimit(uint256 callDataLength) internal view returns (uint256) {
    return ((gasleft() - 2 * (16 * callDataLength + GAS_FOR_CALL_EXACT_CHECK)) * 62) / 64;
  }

  function generateDynamicOffRampConfig(
    address router,
    address priceRegistry
  ) internal pure returns (EVM2EVMOffRamp.DynamicConfig memory) {
    return
      EVM2EVMOffRamp.DynamicConfig({
        permissionLessExecutionThresholdSeconds: PERMISSION_LESS_EXECUTION_THRESHOLD_SECONDS,
        router: router,
        priceRegistry: priceRegistry,
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        maxDataBytes: MAX_DATA_SIZE,
        maxPoolReleaseOrMintGas: MAX_TOKEN_POOL_RELEASE_OR_MINT_GAS
      });
  }

  function generateDynamicOnRampConfig(
    address router,
    address priceRegistry
  ) internal pure returns (EVM2EVMOnRamp.DynamicConfig memory) {
    return
      EVM2EVMOnRamp.DynamicConfig({
        router: router,
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        destGasOverhead: DEST_GAS_OVERHEAD,
        destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
        destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
        destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
        destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
        priceRegistry: priceRegistry,
        maxDataBytes: MAX_DATA_SIZE,
        maxPerMsgGasLimit: MAX_GAS_LIMIT
      });
  }

  function getTokensAndPools(
    address[] memory sourceTokens,
    IPool[] memory pools
  ) internal pure returns (Internal.PoolUpdate[] memory) {
    Internal.PoolUpdate[] memory tokensAndPools = new Internal.PoolUpdate[](sourceTokens.length);
    for (uint256 i = 0; i < sourceTokens.length; ++i) {
      tokensAndPools[i] = Internal.PoolUpdate({token: sourceTokens[i], pool: address(pools[i])});
    }
    return tokensAndPools;
  }

  function getNopsAndWeights() internal pure returns (EVM2EVMOnRamp.NopAndWeight[] memory) {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](3);
    nopsAndWeights[0] = EVM2EVMOnRamp.NopAndWeight({nop: USER_1, weight: 19284});
    nopsAndWeights[1] = EVM2EVMOnRamp.NopAndWeight({nop: USER_2, weight: 52935});
    nopsAndWeights[2] = EVM2EVMOnRamp.NopAndWeight({nop: USER_3, weight: 8});
    return nopsAndWeights;
  }

  // Rate limiter
  address internal constant ADMIN = 0x11118e64e1FB0c487f25dD6D3601FF6aF8d32E4e;

  function getOutboundRateLimiterConfig() internal pure returns (RateLimiter.Config memory) {
    return RateLimiter.Config({isEnabled: true, capacity: 100e28, rate: 1e15});
  }

  function getInboundRateLimiterConfig() internal pure returns (RateLimiter.Config memory) {
    return RateLimiter.Config({isEnabled: true, capacity: 222e30, rate: 1e18});
  }

  function getSingleTokenPriceUpdateStruct(
    address token,
    uint224 price
  ) internal pure returns (Internal.PriceUpdates memory) {
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](1);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: token, usdPerToken: price});

    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    return priceUpdates;
  }

  function getSingleGasPriceUpdateStruct(
    uint64 chainSelector,
    uint224 usdPerUnitGas
  ) internal pure returns (Internal.PriceUpdates memory) {
    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: chainSelector, usdPerUnitGas: usdPerUnitGas});

    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: gasPriceUpdates
    });

    return priceUpdates;
  }

  function getSingleTokenAndGasPriceUpdateStruct(
    address token,
    uint224 price,
    uint64 chainSelector,
    uint224 usdPerUnitGas
  ) internal pure returns (Internal.PriceUpdates memory) {
    Internal.PriceUpdates memory update = getSingleTokenPriceUpdateStruct(token, price);
    update.gasPriceUpdates = getSingleGasPriceUpdateStruct(chainSelector, usdPerUnitGas).gasPriceUpdates;
    return update;
  }

  function getPriceUpdatesStruct(
    address[] memory tokens,
    uint224[] memory prices
  ) internal pure returns (Internal.PriceUpdates memory) {
    uint256 length = tokens.length;

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](length);
    for (uint256 i = 0; i < length; ++i) {
      tokenPriceUpdates[i] = Internal.TokenPriceUpdate({sourceToken: tokens[i], usdPerToken: prices[i]});
    }
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    return priceUpdates;
  }

  // OffRamp
  function getEmptyPriceUpdates() internal pure returns (Internal.PriceUpdates memory priceUpdates) {
    return
      Internal.PriceUpdates({
        tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
        gasPriceUpdates: new Internal.GasPriceUpdate[](0)
      });
  }
}
