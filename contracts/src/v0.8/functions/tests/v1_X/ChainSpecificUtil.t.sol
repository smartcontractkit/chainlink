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
import {L1Block} from "../../../vendor/@eth-optimism/contracts-bedrock/v0.17.1/src/L2/L1Block.sol";

import "forge-std/console.sol";

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
      abi.encodeWithSelector(
        ArbGasInfo.getPricesInWei.selector),
        abi.encode(unused, WEI_PER_L1_CALLDATA_BYTE, unused, unused, unused, unused)
    );
  }

  function test__getL1FeeUpperLimit_SuccessWhenArbitrumMainnet() public {
    // Set the chainID
    vm.chainId(42161);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    uint96 baseCostEstimate = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit);

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeEstimateJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Arbitrum mainnet");

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
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Arbitrum Goerli");

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
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Arbitrum Sepolia");

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Arbitrum Sepolia");
  }
}

/// @notice #_getL1FeeUpperLimit Optimism
/// @dev Optimism gas formula = ((l2_base_fee + l2_priority_fee) * l2_gas_used) + L1 data fee
/// @dev where L1 data fee = l1_gas_price * ((count_zero_bytes(tx_data) * 4 + count_non_zero_bytes(tx_data) * 16) + fixed_overhead + noncalldata_gas) * dynamic_overhead
contract ChainSpecificUtil__getL1FeeUpperLimit_Optimism is FunctionsFulfillmentSetup {
  address private constant L1BLOCK_ADDR = address(0x4200000000000000000000000000000000000015);
  L1Block private constant L1BLOCK = L1Block(L1BLOCK_ADDR);
  uint256 private constant L1_BASE_FEE_WEI = 27_400_956_733;
  uint256 private constant L1_BASE_FEE_SCALAR = 1_368;
  uint256 private constant L1_BLOB_BASE_FEE_WEI = 1;
  uint256 private constant L1_BLOB_BASE_FEE_SCALAR = 810_949;
  uint256 private constant L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI = 26_196_960_380_817_618;
  uint256 private constant L1_FEE_WEI = 14_930_089_287_329_164;

  uint96 l1FeeEstimateJuels = uint96((1e18 * L1_FEE_ESTIMATE_WITH_OVERESTIMATION_WEI) / uint256(LINK_ETH_RATE));
  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.basefee.selector),
      abi.encode(L1_BASE_FEE_WEI)
    );
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.baseFeeScalar.selector),
      abi.encode(L1_BASE_FEE_SCALAR)
    );
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.blobBaseFee.selector),
      abi.encode(L1_BLOB_BASE_FEE_WEI)
    );
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.blobBaseFeeScalar.selector),
      abi.encode(L1_BLOB_BASE_FEE_SCALAR)
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
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Optimism mainnet");

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
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Optimism Goerli");

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
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Optimism Sepolia");

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Optimism Sepolia");
  }
}

/// @notice #_getL1FeeUpperLimit Base
/// @dev Base gas formula uses Optimism formula = ((l2_base_fee + l2_priority_fee) * l2_gas_used) + L1 data fee
/// @dev where L1 data fee = l1_gas_price * ((count_zero_bytes(tx_data) * 4 + count_non_zero_bytes(tx_data) * 16) + fixed_overhead + noncalldata_gas) * dynamic_overhead
contract ChainSpecificUtil__getL1FeeUpperLimit_Base is FunctionsFulfillmentSetup {
  address private constant L1BLOCK_ADDR = address(0x4200000000000000000000000000000000000015);
  L1Block private constant L1BLOCK = L1Block(L1BLOCK_ADDR);
  uint256 private constant L1_BASE_FEE_WEI = 49_576_649_300;
  uint256 private constant L1_BASE_FEE_SCALAR = 1_101;
  uint256 private constant L1_BLOB_BASE_FEE_WEI = 1;
  uint256 private constant L1_BLOB_BASE_FEE_SCALAR = 659_851;

  uint96 l1FeeJuels = uint96((1e18 * L1_BASE_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.basefee.selector),
      abi.encode(L1_BASE_FEE_WEI)
    );
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.baseFeeScalar.selector),
      abi.encode(L1_BASE_FEE_SCALAR)
    );
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.blobBaseFee.selector),
      abi.encode(L1_BLOB_BASE_FEE_WEI)
    );
    vm.mockCall(
      L1BLOCK_ADDR,
      abi.encodeWithSelector(L1BLOCK.blobBaseFeeScalar.selector),
      abi.encode(L1_BLOB_BASE_FEE_SCALAR)
    );
  }

  function test__getL1FeeUpperLimit_SuccessWhenBaseMainnet() public {
    // Set the chainID
    vm.chainId(8453);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Base mainnet");

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
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Base Goerli");

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
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels, "Estimated cost mismatch for Base Sepolia");

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels, "Response cost mismatch for Base Sepolia");
  }
}
