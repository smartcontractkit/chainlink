// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsClient} from "../../dev/v1_X/FunctionsClient.sol";
import {FunctionsRouter} from "../../dev/v1_X/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/v1_X/FunctionsSubscriptions.sol";
import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";

import {FunctionsFulfillmentSetup} from "./Setup.t.sol";

import {ArbGasInfo} from "../../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {GasPriceOracle} from "../../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";

/// @notice #_getL1FeeUpperLimit Arbitrum
/// @dev Arbitrum gas formula = L2 Gas Price * (Gas used on L2 + Extra Buffer for L1 cost)
/// @dev where Extra Buffer for L1 cost = (L1 Estimated Cost / L2 Gas Price)
contract ChainSpecificUtil__getL1FeeUpperLimit_Arbitrum is FunctionsFulfillmentSetup {
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  uint256 private constant WEI_PER_L1_CALLDATA_BYTE = 2_243_708_528;

  uint256 private constant L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI = 6_408_031_555_968;
  uint256 private constant L1_FEE_WEI = 3_697_631_654_144;

  uint96 l1FeeEstimateJuels = uint96((1e18 * L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI) / uint256(LINK_ETH_RATE));
  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    uint256 unused = 0;
    vm.mockCall(
      ARBGAS_ADDR,
      abi.encodeWithSelector(ArbGasInfo.getPricesInWei.selector),
      abi.encode(unused, WEI_PER_L1_CALLDATA_BYTE, unused, unused, unused, unused)
    );
  }

  function test__getL1FeeUpperLimit_SuccessWhenArbitrumMainnet() public {
    // Set the chainID
    vm.chainId(42161);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Arbitrum mainnet"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Arbitrum mainnet");
  }

  function test__getL1FeeUpperLimit_SuccessWhenArbitrumGoerli() public {
    // Set the chainID
    vm.chainId(421613);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Arbitrum Goerli"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Arbitrum Goerli");
  }

  function test__getL1FeeUpperLimit_SuccessWhenArbitrumSepolia() public {
    // Set the chainID
    vm.chainId(421614);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Arbitrum Sepolia"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Arbitrum Sepolia");
  }
}

/// @notice #_getL1FeeUpperLimit Optimism
/// @dev Optimism gas formula: https://docs.optimism.io/stack/transactions/fees#ecotone
/// @dev Note that the docs fail to mention the calculation also requires a division by 10^6
/// @dev See here: https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/exec-engine.md#ecotone-l1-cost-fee-changes-eip-4844-da
/// @dev Also, we conservatively assume all non-zero bytes: tx_compressed_size = tx_data_size_bytes
contract ChainSpecificUtil__getL1FeeUpperLimit_Optimism is FunctionsFulfillmentSetup {
  address private constant GAS_PRICE_ORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  GasPriceOracle private constant GAS_PRICE_ORACLE = GasPriceOracle(GAS_PRICE_ORACLE_ADDR);

  uint256 private constant L1_FEE_WEI = 1_648_118_885_618;
  uint256 private constant L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI = (L1_FEE_WEI * (10_000 + 5_000)) / 10_000;

  uint96 l1FeeEstimateJuels = uint96((1e18 * L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI) / uint256(LINK_ETH_RATE));
  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(
      GAS_PRICE_ORACLE_ADDR,
      abi.encodeWithSelector(GAS_PRICE_ORACLE.getL1FeeUpperBound.selector),
      abi.encode(L1_FEE_WEI)
    );
  }

  function test__getL1FeeUpperLimit_SuccessWhenOptimismMainnet() public {
    // Set the chainID
    vm.chainId(10);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Optimism mainnet"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Optimism mainnet");
  }

  function test__getL1FeeUpperLimit_SuccessWhenOptimismGoerli() public {
    // Set the chainID
    vm.chainId(420);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Optimism Goerli"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Optimism Goerli");
  }

  function test__getL1FeeUpperLimit_SuccessWhenOptimismSepolia() public {
    // Set the chainID
    vm.chainId(11155420);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Optimism Sepolia"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Optimism Sepolia");
  }
}

/// @notice #_getL1FeeUpperLimit Base
/// @dev Base gas formula uses Optimism formula since it is build on the OP chain stack (See comments above for Optimism tests)
contract ChainSpecificUtil__getL1FeeUpperLimit_Base is FunctionsFulfillmentSetup {
  address private constant GAS_PRICE_ORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  GasPriceOracle private constant GAS_PRICE_ORACLE = GasPriceOracle(GAS_PRICE_ORACLE_ADDR);

  uint256 private constant L1_FEE_WEI = 1_648_118_885_618;
  uint256 private constant L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI = (L1_FEE_WEI * (10_000 + 5_000)) / 10_000;

  uint96 l1FeeEstimateJuels = uint96((1e18 * L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI) / uint256(LINK_ETH_RATE));
  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(
      GAS_PRICE_ORACLE_ADDR,
      abi.encodeWithSelector(GAS_PRICE_ORACLE.getL1FeeUpperBound.selector),
      abi.encode(L1_FEE_WEI)
    );
  }

  function test__getL1FeeUpperLimit_SuccessWhenBaseMainnet() public {
    // Set the chainID
    vm.chainId(8453);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Base mainnet"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Base mainnet");
  }

  function test__getL1FeeUpperLimit_SuccessWhenBaseGoerli() public {
    // Set the chainID
    vm.chainId(84531);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Base Goerli"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Base Goerli");
  }

  function test__getL1FeeUpperLimit_SuccessWhenBaseSepolia() public {
    // Set the chainID
    vm.chainId(84532);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(
      s_requests[1].commitment.estimatedTotalCostJuels,
      expectedEstimatedTotalCostJuels,
      "Estimated cost mismatch for Base Sepolia"
    );

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Base Sepolia");
  }
}
