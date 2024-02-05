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
import {OVM_GasPriceOracle} from "../../../vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol";

/// @notice #_getCurrentTxL1GasFees Arbitrum
/// @dev Arbitrum gas formula = L2 Gas Price * (Gas used on L2 + Extra Buffer for L1 cost)
/// @dev where Extra Buffer for L1 cost = (L1 Estimated Cost / L2 Gas Price)
contract ChainSpecificUtil__getCurrentTxL1GasFees_Arbitrum is FunctionsFulfillmentSetup {
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  uint256 private constant L1_FEE_WEI = 15_818_209_764_247;

  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(ARBGAS_ADDR, abi.encodeWithSelector(ArbGasInfo.getCurrentTxL1GasFees.selector), abi.encode(L1_FEE_WEI));
  }

  function test__getCurrentTxL1GasFees_SuccessWhenArbitrumMainnet() public {
    // Set the chainID
    vm.chainId(42161);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }

  function test__getCurrentTxL1GasFees_SuccessWhenArbitrumGoerli() public {
    // Set the chainID
    vm.chainId(421613);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }

  function test__getCurrentTxL1GasFees_SuccessWhenArbitrumSepolia() public {
    // Set the chainID
    vm.chainId(421614);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }
}

/// @notice #_getCurrentTxL1GasFees Optimism
/// @dev Optimism gas formula = ((l2_base_fee + l2_priority_fee) * l2_gas_used) + L1 data fee
/// @dev where L1 data fee = l1_gas_price * ((count_zero_bytes(tx_data) * 4 + count_non_zero_bytes(tx_data) * 16) + fixed_overhead + noncalldata_gas) * dynamic_overhead
contract ChainSpecificUtil__getCurrentTxL1GasFees_Optimism is FunctionsFulfillmentSetup {
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  uint256 private constant L1_FEE_WEI = 15_818_209_764_247;

  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(OVM_GasPriceOracle.getL1Fee.selector),
      abi.encode(L1_FEE_WEI)
    );
  }

  function test__getCurrentTxL1GasFees_SuccessWhenOptimismMainnet() public {
    // Set the chainID
    vm.chainId(10);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }

  function test__getCurrentTxL1GasFees_SuccessWhenOptimismGoerli() public {
    // Set the chainID
    vm.chainId(420);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }

  function test__getCurrentTxL1GasFees_SuccessWhenOptimismSepolia() public {
    // Set the chainID
    vm.chainId(11155420);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }
}

/// @notice #_getCurrentTxL1GasFees Base
/// @dev Base gas formula uses Optimism formula = ((l2_base_fee + l2_priority_fee) * l2_gas_used) + L1 data fee
/// @dev where L1 data fee = l1_gas_price * ((count_zero_bytes(tx_data) * 4 + count_non_zero_bytes(tx_data) * 16) + fixed_overhead + noncalldata_gas) * dynamic_overhead
contract ChainSpecificUtil__getCurrentTxL1GasFees_Base is FunctionsFulfillmentSetup {
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  uint256 private constant L1_FEE_WEI = 15_818_209_764_247;

  uint96 l1FeeJuels = uint96((1e18 * L1_FEE_WEI) / uint256(LINK_ETH_RATE));

  function setUp() public virtual override {
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(OVM_GasPriceOracle.getL1Fee.selector),
      abi.encode(L1_FEE_WEI)
    );
  }

  function test__getCurrentTxL1GasFees_SuccessWhenBaseMainnet() public {
    // Set the chainID
    vm.chainId(8453);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }

  function test__getCurrentTxL1GasFees_SuccessWhenBaseGoerli() public {
    // Set the chainID
    vm.chainId(84531);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }

  function test__getCurrentTxL1GasFees_SuccessWhenBaseSepolia() public {
    // Set the chainID
    vm.chainId(84532);

    // Setup sends and fulfills request #1
    FunctionsFulfillmentSetup.setUp();

    // Check request cost estimate
    uint96 expectedEstimatedTotalCostJuels = _getExpectedCostEstimate(s_requests[1].requestData.callbackGasLimit) +
      l1FeeJuels;
    assertEq(s_requests[1].commitment.estimatedTotalCostJuels, expectedEstimatedTotalCostJuels);

    // Check response actual cost
    uint96 expectedTotalCostJuels = _getExpectedCost(5416) + l1FeeJuels;
    assertEq(s_responses[1].totalCostJuels, expectedTotalCostJuels);
  }
}
