// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./KeystoneForwarderBaseTest.t.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {MaliciousReportReceiver} from "./mocks/MaliciousReportReceiver.sol";
import {MaliciousRevertingReceiver} from "./mocks/MaliciousRevertingReceiver.sol";
import {KeystoneForwarder} from "../KeystoneForwarder.sol";
import {console} from "forge-std/console.sol";

contract KeystoneForwarder_ReportTest is BaseTest {
  event MessageReceived(bytes metadata, bytes[] mercuryReports);
  event ReportProcessed(
    address indexed receiver,
    bytes32 indexed workflowExecutionId,
    bytes2 indexed reportId,
    bool result
  );

  uint8 internal version = 1;
  uint32 internal timestamp = 0;
  bytes32 internal workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  bytes10 internal workflowName = hex"000000000000DEADBEEF";
  address internal workflowOwner = address(51);
  bytes32 internal executionId = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";
  bytes2 internal reportId = hex"0001";
  bytes[] internal mercuryReports = new bytes[](2);
  bytes internal rawReports;
  bytes internal header;
  bytes internal metadata;
  bytes internal report;
  bytes internal reportContext = new bytes(96);
  uint256 internal requiredSignaturesNum = F + 1;
  bytes[] internal signatures = new bytes[](2);

  function setUp() public override {
    BaseTest.setUp();

    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, _getSignerAddresses());
    s_router.addForwarder(address(s_forwarder));

    mercuryReports[0] = hex"010203";
    mercuryReports[1] = hex"aabbccdd";

    rawReports = abi.encode(mercuryReports);
    metadata = abi.encodePacked(workflowId, workflowName, workflowOwner, reportId);
    header = abi.encodePacked(version, executionId, timestamp, DON_ID, CONFIG_VERSION, metadata);
    report = abi.encodePacked(header, rawReports);

    signatures = _signReport(report, reportContext, requiredSignaturesNum);

    vm.startPrank(TRANSMITTER);
  }

  function test_RevertWhen_ReportHasIncorrectDON() public {
    uint32 invalidDONId = 111;
    bytes memory reportWithInvalidDONId = abi.encodePacked(
      version,
      executionId,
      timestamp,
      invalidDONId,
      CONFIG_VERSION,
      workflowId,
      workflowName,
      workflowOwner,
      reportId,
      rawReports
    );

    uint64 configId = (uint64(invalidDONId) << 32) | CONFIG_VERSION;
    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidConfig.selector, configId));
    s_forwarder.report(address(s_receiver), reportWithInvalidDONId, reportContext, signatures);
  }

  function test_RevertWhen_ReportHasInexistentConfigVersion() public {
    bytes memory reportWithInvalidDONId = abi.encodePacked(
      version,
      executionId,
      timestamp,
      DON_ID,
      CONFIG_VERSION + 1,
      workflowId,
      workflowName,
      workflowOwner,
      reportId,
      rawReports
    );

    uint64 configId = (uint64(DON_ID) << 32) | (CONFIG_VERSION + 1);
    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidConfig.selector, configId));
    s_forwarder.report(address(s_receiver), reportWithInvalidDONId, reportContext, signatures);
  }

  function test_RevertWhen_ReportIsMalformed() public {
    bytes memory shortenedReport = abi.encode(bytes32(report));

    vm.expectRevert(KeystoneForwarder.InvalidReport.selector);
    s_forwarder.report(address(s_receiver), shortenedReport, reportContext, signatures);
  }

  function test_RevertWhen_TooFewSignatures() public {
    bytes[] memory fewerSignatures = new bytes[](F);

    vm.expectRevert(
      abi.encodeWithSelector(KeystoneForwarder.InvalidSignatureCount.selector, F + 1, fewerSignatures.length)
    );
    s_forwarder.report(address(s_receiver), report, reportContext, fewerSignatures);
  }

  function test_RevertWhen_TooManySignatures() public {
    bytes[] memory moreSignatures = new bytes[](F + 2);

    vm.expectRevert(
      abi.encodeWithSelector(KeystoneForwarder.InvalidSignatureCount.selector, F + 1, moreSignatures.length)
    );
    s_forwarder.report(address(s_receiver), report, reportContext, moreSignatures);
  }

  function test_RevertWhen_AnySignatureIsInvalid() public {
    signatures[1] = abi.encode(1234); // invalid signature

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidSignature.selector, signatures[1]));
    s_forwarder.report(address(s_receiver), report, reportContext, signatures);
  }

  function test_RevertWhen_AnySignerIsInvalid() public {
    uint256 mockPK = 999;

    Signer memory maliciousSigner = Signer({mockPrivateKey: mockPK, signerAddress: vm.addr(mockPK)});
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(
      maliciousSigner.mockPrivateKey,
      keccak256(abi.encodePacked(keccak256(report), reportContext))
    );
    signatures[1] = bytes.concat(r, s, bytes1(v - 27));

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidSigner.selector, maliciousSigner.signerAddress));
    s_forwarder.report(address(s_receiver), report, reportContext, signatures);
  }

  function test_RevertWhen_ReportHasDuplicateSignatures() public {
    signatures[1] = signatures[0]; // repeat a signature

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.DuplicateSigner.selector, s_signers[0].signerAddress));
    s_forwarder.report(address(s_receiver), report, reportContext, signatures);
  }

  function test_RevertWhen_RetryingSuccessfulTransmission() public {
    s_forwarder.report{gas: 400_000}(address(s_receiver), report, reportContext, signatures);

    bytes32 transmissionId = s_forwarder.getTransmissionId(address(s_receiver), executionId, reportId);
    vm.expectRevert(abi.encodeWithSelector(IRouter.AlreadyAttempted.selector, transmissionId));
    // Retyring with more gas
    s_forwarder.report{gas: 450_000}(address(s_receiver), report, reportContext, signatures);
  }

  function test_RevertWhen_RetryingInvalidContractTransmission() public {
    // Receiver is not a contract
    address receiver = address(404);
    s_forwarder.report{gas: 400_000}(receiver, report, reportContext, signatures);

    bytes32 transmissionId = s_forwarder.getTransmissionId(receiver, executionId, reportId);
    vm.expectRevert(abi.encodeWithSelector(IRouter.AlreadyAttempted.selector, transmissionId));
    // Retyring with more gas
    s_forwarder.report{gas: 450_000}(receiver, report, reportContext, signatures);
  }

  function test_RevertWhen_AttemptingTransmissionWithInsufficientGas() public {
    bytes32 transmissionId = s_forwarder.getTransmissionId(address(s_receiver), executionId, reportId);
    vm.expectRevert(abi.encodeWithSelector(IRouter.InsufficientGasForRouting.selector, transmissionId));
    s_forwarder.report{gas: 150_000}(address(s_receiver), report, reportContext, signatures);
  }

  function test_Report_SuccessfulDelivery() public {
    IRouter.TransmissionInfo memory transmissionInfo = s_forwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.NOT_ATTEMPTED), "state mismatch");

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(metadata, mercuryReports);

    vm.expectEmit(address(s_forwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);

    s_forwarder.report(address(s_receiver), report, reportContext, signatures);

    transmissionInfo = s_forwarder.getTransmissionInfo(address(s_receiver), executionId, reportId);

    assertEq(transmissionInfo.transmitter, TRANSMITTER, "transmitter mismatch");
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
    assertGt(transmissionInfo.gasLimit, 100_000, "gas limit mismatch");
  }

  function test_Report_SuccessfulRetryWithMoreGas() public {
    s_forwarder.report{gas: 200_000}(address(s_receiver), report, reportContext, signatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_forwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );
    // Expect to fail with the receiver running out of gas
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.FAILED), "state mismatch");
    assertGt(transmissionInfo.gasLimit, 100_000, "gas limit mismatch");

    // Should succeed with more gas
    s_forwarder.report{gas: 300_000}(address(s_receiver), report, reportContext, signatures);

    transmissionInfo = s_forwarder.getTransmissionInfo(address(s_receiver), executionId, reportId);
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
    assertGt(transmissionInfo.gasLimit, 200_000, "gas limit mismatch");
  }

  function test_Report_FailedDeliveryWhenReceiverNotContract() public {
    // Receiver is not a contract
    address receiver = address(404);

    vm.expectEmit(address(s_forwarder));
    emit ReportProcessed(receiver, executionId, reportId, false);

    s_forwarder.report(receiver, report, reportContext, signatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_forwarder.getTransmissionInfo(receiver, executionId, reportId);
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.INVALID_RECEIVER), "state mismatch");
  }

  function test_Report_FailedDeliveryWhenReceiverInterfaceNotSupported() public {
    // Receiver is a contract but doesn't implement the required interface
    address receiver = address(s_forwarder);

    vm.expectEmit(true, true, true, false);
    emit ReportProcessed(receiver, executionId, reportId, false);

    s_forwarder.report(receiver, report, reportContext, signatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_forwarder.getTransmissionInfo(receiver, executionId, reportId);
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.INVALID_RECEIVER), "state mismatch");
  }

  //  function test_Report_FailedDeliveryWhenReportReceiverConsumesAllGasAndInterfaceCheckUsesMax() public {
  //    MaliciousRevertingReceiver maliciousReceiver = new MaliciousRevertingReceiver();
  //    // This should not revert if gas tracking is effective
  //    // It may revert if it fails to reserve sufficient gas for routing
  //    // This POC requires pretty specific initial gas, so that 1/64 of gas passed to `onReport()` is insufficient to store the success
  //    s_forwarder.report{gas: 200_000}(address(maliciousReceiver), report, reportContext, signatures);
  //
  //    IRouter.TransmissionInfo memory transmissionInfo = s_forwarder.getTransmissionInfo(
  //      address(maliciousReceiver),
  //      executionId,
  //      reportId
  //    );
  //
  //    assertEq(transmissionInfo.transmitter, TRANSMITTER, "transmitter mismatch");
  //    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
  //  }

  function test_Report_FailedDelieryWhenReportReceiverConsumesAllGas() public {
    MaliciousReportReceiver s_maliciousReceiver = new MaliciousReportReceiver();
    s_forwarder.report{gas: 500_000}(address(s_maliciousReceiver), report, reportContext, signatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_forwarder.getTransmissionInfo(
      address(s_maliciousReceiver),
      executionId,
      reportId
    );

    assertEq(transmissionInfo.transmitter, TRANSMITTER, "transmitter mismatch");
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.FAILED), "state mismatch");
    assertGt(transmissionInfo.gasLimit, 410_000, "gas limit mismatch");
  }

  function test_Report_ConfigVersion() public {
    vm.stopPrank();
    // configure a new configVersion
    vm.prank(ADMIN);
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION + 1, F, _getSignerAddresses());

    // old version still works
    vm.expectEmit(address(s_receiver));
    emit MessageReceived(metadata, mercuryReports);

    vm.expectEmit(address(s_forwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);

    vm.prank(TRANSMITTER);
    s_forwarder.report(address(s_receiver), report, reportContext, signatures);

    // after clear the old version doesn't work anymore
    vm.prank(ADMIN);
    s_forwarder.clearConfig(DON_ID, CONFIG_VERSION);

    uint64 configId = (uint64(DON_ID) << 32) | CONFIG_VERSION;
    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidConfig.selector, configId));
    vm.prank(TRANSMITTER);
    s_forwarder.report(address(s_receiver), report, reportContext, signatures);

    // but new config does
    bytes32 newExecutionId = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000001";
    bytes memory newMetadata = abi.encodePacked(workflowId, workflowName, workflowOwner, reportId);
    bytes memory newHeader = abi.encodePacked(
      version,
      newExecutionId,
      timestamp,
      DON_ID,
      CONFIG_VERSION + 1,
      newMetadata
    );
    bytes memory newReport = abi.encodePacked(newHeader, rawReports);
    // resign the new report
    bytes[] memory newSignatures = _signReport(newReport, reportContext, requiredSignaturesNum);

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(newMetadata, mercuryReports);

    vm.expectEmit(address(s_forwarder));
    emit ReportProcessed(address(s_receiver), newExecutionId, reportId, true);

    vm.prank(TRANSMITTER);
    s_forwarder.report(address(s_receiver), newReport, reportContext, newSignatures);

    // validate transmitter was recorded
    address transmitter = s_forwarder.getTransmitter(address(s_receiver), newExecutionId, reportId);
    assertEq(transmitter, TRANSMITTER, "transmitter mismatch");
  }
}
