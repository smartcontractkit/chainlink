pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5_Optimism} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5_Optimism.sol";
import {OptimismL1Fees} from "../dev/OptimismL1Fees.sol";
import {BlockhashStore} from "../dev/BlockhashStore.sol";
import {GasPriceOracle as OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";
import {VmSafe} from "forge-std/Vm.sol";

contract VRFV2CoordinatorV2_5_Optimism is BaseTest {
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  /// @dev L1_FEE_DATA_PADDING includes 71 bytes for L1 data padding for Optimism
  bytes internal constant L1_FEE_DATA_PADDING =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";

  address internal constant DEPLOYER = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;

  BlockhashStore s_bhs;
  ExposedVRFCoordinatorV2_5_Optimism s_testCoordinator;
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkNativeFeed;

  uint256 s_startGas = 0.0038 gwei;
  uint256 s_weiPerUnitGas = 0.003 gwei;

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
    changePrank(DEPLOYER);

    vm.txGasPrice(100 gwei);

    // Instantiate BHS.
    s_bhs = new BlockhashStore();

    // Deploy coordinator, LINK token and LINK/Native feed.
    s_testCoordinator = new ExposedVRFCoordinatorV2_5_Optimism(address(s_bhs));
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Configure the coordinator.
    s_testCoordinator.setLINKAndLINKNativeFeed(address(s_linkToken), address(s_linkNativeFeed));
    s_testCoordinator.setConfig(
      0, // minRequestConfirmations
      2_500_000, // maxGasLimit
      1, // stalenessSeconds
      50_000, // gasAfterPaymentCalculation
      50000000000000000, // fallbackWeiPerUnitLink
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );
  }

  function _encodeCalculatePaymentAmountNativeExternal(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) internal pure returns (bytes memory) {
    return
      abi.encodeWithSelector(
        ExposedVRFCoordinatorV2_5_Optimism.calculatePaymentAmountNativeExternal.selector,
        startGas,
        weiPerUnitGas,
        onlyPremium
      );
  }

  function _encodeCalculatePaymentAmountLinkExternal(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) internal pure returns (bytes memory) {
    return
      abi.encodeWithSelector(
        ExposedVRFCoordinatorV2_5_Optimism.calculatePaymentAmountLinkExternal.selector,
        startGas,
        weiPerUnitGas,
        onlyPremium
      );
  }

  function _mockGasOraclePriceGetL1FeeUpperBoundCall() internal {
    // 171 bytes is the size of tx.data we are sending in this test
    // this is not expected fulfillment tx size!
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(OVM_GasPriceOracle.getL1FeeUpperBound.selector, 171),
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

  function _mockGasOraclePriceGetL1FeeCall(bytes memory txMsgData) internal {
    vm.mockCall(
      OVM_GASPRICEORACLE_ADDR,
      abi.encodeWithSelector(OVM_GasPriceOracle.getL1Fee.selector, bytes.concat(txMsgData, L1_FEE_DATA_PADDING)),
      abi.encode(uint256(0.001 ether))
    );
  }

  function _checkL1GasFeeEmittedLogs(uint256 expectedL1GasFee) internal {
    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries.length, 1);
    assertEq(entries[0].topics.length, 1);
    assertEq(entries[0].topics[0], keccak256("L1GasFee(uint256)"));
    // 1e15 is less than 1 percent discrepancy
    uint256 actualL1GasFee = abi.decode(entries[0].data, (uint256));
    assertApproxEqAbs(expectedL1GasFee, actualL1GasFee, 1e15);
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

  function test_setL1FeePaymentMethod() public {
    // check default settings after contract deployment
    assertEq(uint256(L1_GAS_FEES_MODE), uint256(s_testCoordinator.s_l1FeeCalculationMode()));
    assertEq(100, uint256(s_testCoordinator.s_l1FeeCoefficient()));
    vm.recordLogs();

    s_testCoordinator.setL1FeeCalculation(L1_CALLDATA_GAS_COST_MODE, 70);

    _checkL1FeeCalculationSetEmittedLogs(L1_CALLDATA_GAS_COST_MODE, 70);
    assertEq(uint256(L1_CALLDATA_GAS_COST_MODE), uint256(s_testCoordinator.s_l1FeeCalculationMode()));
    assertEq(70, uint256(s_testCoordinator.s_l1FeeCoefficient()));

    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 30);

    _checkL1FeeCalculationSetEmittedLogs(L1_GAS_FEES_UPPER_BOUND_MODE, 30);
    assertEq(uint256(L1_GAS_FEES_UPPER_BOUND_MODE), uint256(s_testCoordinator.s_l1FeeCalculationMode()));
    assertEq(30, uint256(s_testCoordinator.s_l1FeeCoefficient()));

    // should revert if invalid L1 fee calculation mode is used
    vm.expectRevert(abi.encodeWithSelector(OptimismL1Fees.InvalidL1FeeCalculationMode.selector, 4));
    s_testCoordinator.setL1FeeCalculation(4, 100);

    // should revert if invalid coefficient is used (equal to zero, this would disable L1 fees completely)
    vm.expectRevert(abi.encodeWithSelector(OptimismL1Fees.InvalidL1FeeCoefficient.selector, 0));
    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 0);

    // should revert if invalid coefficient is used (larger than 100%)
    vm.expectRevert(abi.encodeWithSelector(OptimismL1Fees.InvalidL1FeeCoefficient.selector, 150));
    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 150);
  }

  function test_getBlockNumber() public {
    // sanity check that OP stack will use block.number directly
    vm.roll(1589);
    assertEq(1589, s_testCoordinator.getBlockNumberExternal());
  }

  function test_getBlockhash() public {
    // sanity check that OP stack will use blockhash() directly
    vm.roll(1589);
    bytes32 blockHash = blockhash(1589);
    assertEq(blockHash, s_testCoordinator.getBlockhashExternal(1589));
  }

  // payment calculation depends on the L1 gas cost for fulfillment transaction payload
  // both calculatePayment functions use msg.data passed down from the fulfillRandomWords function and
  // in that case, msg.data contains the fulfillment transaction payload (calldata)
  // it's not easy to simulate this in tests below plus we only want to concentrate on the payment part
  // since we don't have to test with the correct payload, we can use any kind of payload for msg.data
  // in the case of tests below, msg.data will carry calculatePayment function selectors and parameters

  function test_calculatePaymentAmountNativeUsingL1GasFeesMode() public {
    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_MODE, 100);

    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeCall(txMsgData);
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.001 ether));

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.0162937 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeCall(txMsgData);

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.001 ether));

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.0015168 * 1e17, 1e15);
  }

  function test_calculatePaymentAmountLinkUsingL1GasFeesMode() public {
    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_MODE, 100);

    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountLinkExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeCall(txMsgData);
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.001 ether));

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.0222475 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountLinkExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeCall(txMsgData);

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.001 ether));

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.0020225 * 1e17, 1e15);
  }

  function test_calculatePaymentAmountNativeUsingCalldataGasCostMode() public {
    // for this type of cost calculation we are applying coefficient to reduce the inflated gas price
    s_testCoordinator.setL1FeeCalculation(L1_CALLDATA_GAS_COST_MODE, 70);

    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceFeeMethods();
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.0002352 ether)); // 3.36e14 actual price times the coefficient (0.7)

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.002834 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceFeeMethods();

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.0002352 ether)); // 3.36e14 actual price times the coefficient (0.7)

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.00037 * 1e17, 1e15);
  }

  function test_calculatePaymentAmountLinkUsingCalldataGasCostMode() public {
    // for this type of cost calculation we are applying coefficient to reduce the inflated gas price
    s_testCoordinator.setL1FeeCalculation(L1_CALLDATA_GAS_COST_MODE, 70);

    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountLinkExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceFeeMethods();
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.0002352 ether)); // 3.36e14 actual price times the coefficient (0.7)

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.0054219 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountLinkExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceFeeMethods();

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.0002352 ether)); // 3.36e14 actual price times the coefficient (0.7)

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.0004929 * 1e17, 1e15);
  }

  function test_calculatePaymentAmountNativeUsingL1GasFeesUpperBoundMode() public {
    // for this type of cost calculation we are applying coefficient to reduce the inflated gas price
    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 50);

    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeUpperBoundCall();
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.01 ether)); // 2e16 actual price times the coefficient (0.5)

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.115129 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountNativeExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeUpperBoundCall();

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.01 ether)); // 2e16 actual price times the coefficient (0.5)

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.015017 * 1e17, 1e15);
  }

  function test_calculatePaymentAmountLinkUsingL1GasFeesUpperBoundMode() public {
    // for this type of cost calculation we are applying coefficient to reduce the inflated gas price
    s_testCoordinator.setL1FeeCalculation(L1_GAS_FEES_UPPER_BOUND_MODE, 50);

    // first we test premium and flat fee payment combined
    bool onlyPremium = false;
    bytes memory txMsgData = _encodeCalculatePaymentAmountLinkExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeUpperBoundCall();
    vm.recordLogs();

    uint256 gasLimit = 0.0001 gwei; // needed because gasleft() is used in the payment calculation
    (bool success, bytes memory returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.01 ether)); // 2e16 actual price times the coefficient (0.5)

    uint96 payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.2202475 * 1e17, 1e15);

    // now we test only premium payment
    onlyPremium = true;
    txMsgData = _encodeCalculatePaymentAmountLinkExternal(s_startGas, s_weiPerUnitGas, onlyPremium);
    _mockGasOraclePriceGetL1FeeUpperBoundCall();

    (success, returnData) = address(s_testCoordinator).call{gas: gasLimit}(txMsgData);
    assertTrue(success);
    _checkL1GasFeeEmittedLogs(uint256(0.01 ether)); // 2e16 actual price times the coefficient (0.5)

    payment = abi.decode(returnData, (uint96));
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.0200225 * 1e17, 1e15);
  }
}
