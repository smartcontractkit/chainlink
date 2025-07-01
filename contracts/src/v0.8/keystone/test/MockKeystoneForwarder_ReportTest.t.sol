// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockBaseTest} from "./MockKeystoneForwarderBaseTest.t.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {MockKeystoneForwarder} from "../MockKeystoneForwarder.sol";

contract MockKeystoneForwarder_ReportTest is MockBaseTest {
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
  bytes[] internal signatures = new bytes[](0); // Empty signatures since mock doesn't validate

  function setUp() public override {
    MockBaseTest.setUp();

    mercuryReports[0] = hex"010203";
    mercuryReports[1] = hex"aabbccdd";

    rawReports = abi.encode(mercuryReports);
    metadata = abi.encodePacked(workflowId, workflowName, workflowOwner, reportId);
    header = abi.encodePacked(version, executionId, timestamp, DON_ID, CONFIG_VERSION, metadata);
    report = abi.encodePacked(header, rawReports);

    vm.startPrank(TRANSMITTER);
  }

  function test_Report_SuccessfulDelivery() public {
    IRouter.TransmissionInfo memory transmissionInfo = s_mockForwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.NOT_ATTEMPTED), "state mismatch");

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(metadata, mercuryReports);

    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);

    s_mockForwarder.report(address(s_receiver), report, reportContext, signatures);

    transmissionInfo = s_mockForwarder.getTransmissionInfo(address(s_receiver), executionId, reportId);

    assertEq(transmissionInfo.transmitter, TRANSMITTER, "transmitter mismatch");
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
  }

  function test_Report_PermissionlessExecution() public {
    // Test that anyone can call report (no signature validation)
    vm.stopPrank();
    address randomCaller = address(123);
    vm.prank(randomCaller);

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(metadata, mercuryReports);

    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);

    s_mockForwarder.report(address(s_receiver), report, reportContext, signatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_mockForwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );

    assertEq(transmissionInfo.transmitter, randomCaller, "transmitter mismatch");
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
  }

  function test_RevertWhen_ReportIsMalformed() public {
    bytes memory shortenedReport = abi.encode(bytes32(report));

    vm.expectRevert(MockKeystoneForwarder.InvalidReport.selector);
    s_mockForwarder.report(address(s_receiver), shortenedReport, reportContext, signatures);
  }

  function test_Report_FailedDeliveryWhenReceiverNotContract() public {
    // Receiver is not a contract
    address receiver = address(404);

    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(receiver, executionId, reportId, false);

    s_mockForwarder.report(receiver, report, reportContext, signatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_mockForwarder.getTransmissionInfo(receiver, executionId, reportId);
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.FAILED), "state mismatch");
  }

  function test_Report_MultipleReportsWithDifferentIds() public {
    bytes2 reportId2 = hex"0002";
    bytes32 executionId2 = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000001";
    
    bytes memory metadata2 = abi.encodePacked(workflowId, workflowName, workflowOwner, reportId2);
    bytes memory header2 = abi.encodePacked(version, executionId2, timestamp, DON_ID, CONFIG_VERSION, metadata2);
    bytes memory report2 = abi.encodePacked(header2, rawReports);

    // First report
    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);
    s_mockForwarder.report(address(s_receiver), report, reportContext, signatures);

    // Second report with different ID
    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(address(s_receiver), executionId2, reportId2, true);
    s_mockForwarder.report(address(s_receiver), report2, reportContext, signatures);

    // Both should be successful
    IRouter.TransmissionInfo memory transmissionInfo1 = s_mockForwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );
    IRouter.TransmissionInfo memory transmissionInfo2 = s_mockForwarder.getTransmissionInfo(
      address(s_receiver),
      executionId2,
      reportId2
    );

    assertEq(uint8(transmissionInfo1.state), uint8(IRouter.TransmissionState.SUCCEEDED), "first report state mismatch");
    assertEq(uint8(transmissionInfo2.state), uint8(IRouter.TransmissionState.SUCCEEDED), "second report state mismatch");
  }

  function test_Report_EmptySignaturesArray() public {
    // Mock doesn't validate signatures, so empty array should work
    bytes[] memory emptySignatures = new bytes[](0);

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(metadata, mercuryReports);

    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);

    s_mockForwarder.report(address(s_receiver), report, reportContext, emptySignatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_mockForwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
  }

  function test_Report_ArbitrarySignaturesIgnored() public {
    // Mock doesn't validate signatures, so invalid signatures should still work
    bytes[] memory invalidSignatures = new bytes[](3);
    invalidSignatures[0] = hex"deadbeef";
    invalidSignatures[1] = hex"cafebabe";
    invalidSignatures[2] = hex"1234567890abcdef";

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(metadata, mercuryReports);

    vm.expectEmit(address(s_mockForwarder));
    emit ReportProcessed(address(s_receiver), executionId, reportId, true);

    s_mockForwarder.report(address(s_receiver), report, reportContext, invalidSignatures);

    IRouter.TransmissionInfo memory transmissionInfo = s_mockForwarder.getTransmissionInfo(
      address(s_receiver),
      executionId,
      reportId
    );
    assertEq(uint8(transmissionInfo.state), uint8(IRouter.TransmissionState.SUCCEEDED), "state mismatch");
  }

  function test_GetTransmissionId() public view {
    bytes32 transmissionId = s_mockForwarder.getTransmissionId(address(s_receiver), executionId, reportId);
    bytes32 expectedId = keccak256(bytes.concat(bytes20(uint160(address(s_receiver))), executionId, reportId));
    assertEq(transmissionId, expectedId, "transmission ID mismatch");
  }

  function test_GetTransmitter() public {
    // Initially should be zero address
    address transmitter = s_mockForwarder.getTransmitter(address(s_receiver), executionId, reportId);
    assertEq(transmitter, address(0), "initial transmitter should be zero");

    // After report, should be the caller
    s_mockForwarder.report(address(s_receiver), report, reportContext, signatures);
    
    transmitter = s_mockForwarder.getTransmitter(address(s_receiver), executionId, reportId);
    assertEq(transmitter, TRANSMITTER, "transmitter mismatch after report");
  }
}