// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {MockLinkToken} from "../../functions/tests/v1_X/testhelpers/MockLinkToken.sol";
import {MockV3Aggregator} from "../../shared/mocks/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5_Optimism} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5_Optimism.sol";
import {VRFV2PlusWrapper_Optimism} from "../dev/VRFV2PlusWrapper_Optimism.sol";
import {OptimismL1Fees} from "../dev/OptimismL1Fees.sol";
import {GasPriceOracle as OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";
import {VmSafe} from "forge-std/Vm.sol";

contract VRFV2PlusWrapperOptimismAndBaseTest is BaseTest {
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  address internal constant DEPLOYER = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  bytes32 private vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint256 private s_wrapperSubscriptionId;

  ExposedVRFCoordinatorV2_5_Optimism private s_testCoordinator;
  MockLinkToken private s_linkToken;
  MockV3Aggregator private s_linkNativeFeed;
  VRFV2PlusWrapper_Optimism private s_wrapper;

  /// @dev Option 1: getL1Fee() function from predeploy GasPriceOracle contract with the fulfillment calldata payload
  /// @dev This option is only available for the Coordinator contract
  uint8 internal constant L1_GAS_FEES_MODE = 0;
  /// @dev Option 2: our own implementation of getL1Fee() function (Ecotone version) with projected
  /// @dev fulfillment calldata payload (number of non-zero bytes estimated based on historical data)
  /// @dev This option is available for the Coordinator and the Wrapper contract
  uint8 internal constant L1_CALLDATA_GAS_COST_MODE = 1;
  /// @dev Option 3: getL1FeeUpperBound() function from predeploy GasPriceOracle contract (available after Fjord upgrade)
  /// @dev This option is available for the Coordinator and the Wrapper contract
  uint8 internal constant L1_GAS_FEES_UPPER_BOUND_MODE = 2;

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(DEPLOYER, 10_000 ether);
    vm.stopPrank();
    vm.startPrank(DEPLOYER);

    // Deploy link token and link/native feed.
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Deploy coordinator.
    s_testCoordinator = new ExposedVRFCoordinatorV2_5_Optimism(address(0));

    // Create subscription for all future wrapper contracts.
    s_wrapperSubscriptionId = s_testCoordinator.createSubscription();

    // Deploy wrapper.
    s_wrapper = new VRFV2PlusWrapper_Optimism(
      address(s_linkToken),
      address(s_linkNativeFeed),
      address(s_testCoordinator),
      uint256(s_wrapperSubscriptionId)
    );

    // Configure the wrapper.
    s_wrapper.setConfig(
      100_000, // wrapper gas overhead
      200_000, // coordinator gas overhead native
      220_000, // coordinator gas overhead link
      500, // coordinator gas overhead per word
      15, // native premium percentage,
      10, // link premium percentage
      vrfKeyHash, // keyHash
      10, // max number of words,
      1, // stalenessSeconds
      50000000000000000, // fallbackWeiPerUnitLink
      500_000, // fulfillmentFlatFeeNativePPM
      100_000 // fulfillmentFlatFeeLinkDiscountPPM
    );

    // Add wrapper as a consumer to the wrapper's subscription.
    s_testCoordinator.addConsumer(uint256(s_wrapperSubscriptionId), address(s_wrapper));
  }

  function _mockGasOraclePriceGetL1FeeUpperBoundCall() internal {
    // fullfillment tx calldata size = 772 bytes
    // RLP-encoded unsigned tx headers (approx padding size) = 71 bytes
    // total = 843 bytes
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(OVM_GasPriceOracle.getL1FeeUpperBound.selector, 843),
      abi.encode(uint256(0.02 ether))
    );
  }

  function _mockGasOraclePriceFeeMethods() internal {
    // these values are taken from an example transaction on Base Sepolia
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(bytes4(keccak256("l1BaseFee()"))),
      abi.encode(64273426165)
    );
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(bytes4(keccak256("baseFeeScalar()"))),
      abi.encode(1101)
    );
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(bytes4(keccak256("blobBaseFeeScalar()"))),
      abi.encode(659851)
    );
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(bytes4(keccak256("blobBaseFee()"))),
      abi.encode(2126959908362)
    );
    vm.mockCall(OVM_GASPRICEORACLE_ADDR, abi.encodeWithSelector(bytes4(keccak256("decimals()"))), abi.encode(6));
  }

  function _checkL1FeeCalculationSetEmittedLogs(uint8 expectedMode, uint8 expectedCoefficient) internal {
    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries.length, 1);
    assertEq(entries[0].topics.length, 1);
    assertEq(entries[0].topics[0], keccak256("L1FeeCalculationSet(uint8,uint8)"));
    (uint8 actualMode, uint8 actualCoefficient) = abi.decode(entries[0].data, (uint8, uint8));
    assertEq(expectedMode, actualMode);
    assertEq(expectedCoefficient, actualCoefficient);
  }

  function test_setL1FeePaymentMethodOnOptimismWrapper() public {
    // check default settings after contract deployment
    assertEq(uint256(L1_CALLDATA_GAS_COST_MODE), uint256(s_wrapper.s_l1FeeCalculationMode()));
    assertEq(100, uint256(s_wrapper.s_l1FeeCoefficient()));

    vm.recordLogs();
    s_wrapper.setL1FeeCalculation(L1_CALLDATA_GAS_COST_MODE, 70);

    _checkL1FeeCalculationSetEmittedLogs(L1_CALLDATA_GAS_COST_MODE, 70);
    assertEq(uint256(L1_CALLDATA_GAS_COST_MODE), uint256(s_wrapper.s_l1FeeCalculationMode()));
    assertEq(70, uint256(s_wrapper.s_l1FeeCoefficient()));

    s_wrapper.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 30);

    _checkL1FeeCalculationSetEmittedLogs(L1_GAS_FEES_UPPER_BOUND_MODE, 30);
    assertEq(uint256(L1_GAS_FEES_UPPER_BOUND_MODE), uint256(s_wrapper.s_l1FeeCalculationMode()));
    assertEq(30, uint256(s_wrapper.s_l1FeeCoefficient()));

    // VRFWrapper doesn't support this mode
    vm.expectRevert(
      abi.encodeWithSelector(VRFV2PlusWrapper_Optimism.UnsupportedL1FeeCalculationMode.selector, L1_GAS_FEES_MODE)
    );
    s_wrapper.setL1FeeCalculation(L1_GAS_FEES_MODE, 100);

    // should revert if invalid L1 fee calculation mode is used
    vm.expectRevert(abi.encodeWithSelector(OptimismL1Fees.InvalidL1FeeCalculationMode.selector, 6));
    s_wrapper.setL1FeeCalculation(6, 100);

    // should revert if invalid coefficient is used (equal to zero, this would disable L1 fees completely)
    vm.expectRevert(abi.encodeWithSelector(OptimismL1Fees.InvalidL1FeeCoefficient.selector, 0));
    s_wrapper.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 0);

    // should revert if invalid coefficient is used (larger than 100%)
    vm.expectRevert(abi.encodeWithSelector(OptimismL1Fees.InvalidL1FeeCoefficient.selector, 101));
    s_wrapper.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 101);
  }

  function test_calculateRequestPriceNativeOnOptimismWrapper_UsingCalldataCostCall() public {
    s_wrapper.setL1FeeCalculation(L1_CALLDATA_GAS_COST_MODE, 80);

    vm.txGasPrice(1 gwei);
    _mockGasOraclePriceFeeMethods();
    uint32 callbackGasLimit = 1_000_000;
    uint32 numWords = 5;
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPriceNative(callbackGasLimit, numWords);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(wrapperCostCalculation, 5.02575 * 1e17, 1e15);
  }

  function test_calculateRequestPriceLinkOnOptimismWrapper_UsingCalldataCostCall() public {
    s_wrapper.setL1FeeCalculation(L1_CALLDATA_GAS_COST_MODE, 80);

    vm.txGasPrice(1 gwei);
    _mockGasOraclePriceFeeMethods();
    uint32 callbackGasLimit = 1_000_000;
    uint32 numWords = 5;
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPrice(callbackGasLimit, numWords);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(wrapperCostCalculation, 8.04934 * 1e17, 1e15);
  }

  function test_calculateRequestPriceNativeOnOptimismWrapper_UsingGetL1FeeUpperBoundCall() public {
    s_wrapper.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 60);

    vm.txGasPrice(1 gwei);
    _mockGasOraclePriceGetL1FeeUpperBoundCall();
    uint32 callbackGasLimit = 1_000_000;
    uint32 numWords = 5;
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPriceNative(callbackGasLimit, numWords);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(wrapperCostCalculation, 5.15283 * 1e17, 1e15);
  }

  function test_calculateRequestPriceLinkOnOptimismWrapper_UsingGetL1FeeUpperBoundCall() public {
    s_wrapper.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 60);

    vm.txGasPrice(1 gwei);
    _mockGasOraclePriceGetL1FeeUpperBoundCall();
    uint32 callbackGasLimit = 1_000_000;
    uint32 numWords = 5;
    uint256 wrapperCostCalculation = s_wrapper.calculateRequestPrice(callbackGasLimit, numWords);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(wrapperCostCalculation, 8.29246 * 1e17, 1e15);
  }
}
