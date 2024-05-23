// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./KeystoneForwarderBaseTest.t.sol";
import {KeystoneForwarder} from "../KeystoneForwarder.sol";

contract KeystoneForwarder_ReportTest is BaseTest {
  event MessageReceived(
    bytes32 workflowId,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bytes32 reportName,
    bytes rawReport
  );

  event ReportProcessed(
    address indexed receiver,
    bytes32 indexed workflowExecutionId,
    bytes32 indexed reportName,
    bytes32 workflowOwner,
    bytes32 workflowName,
    bool result
  );

  bytes32 internal EXECUTION_ID = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";
  bytes32 internal workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  bytes32 internal WORKFLOW_OWNER = bytes32(abi.encodePacked(address(51)));
  bytes32 internal WORKFLOW_NAME = bytes32(bytes("my_workflow"));
  bytes32 internal REPORT_NAME = bytes32(bytes("my_report"));
  bytes[] internal mercuryReports = new bytes[](2);
  bytes internal rawReports;
  bytes internal report;
  uint256 internal requiredSignaturesNum = F + 1;
  bytes[] internal signatures = new bytes[](2);

  function setUp() public override {
    BaseTest.setUp();

    s_forwarder.setConfig(DON_ID, F, _getSignerAddresses());

    mercuryReports[0] = hex"010203";
    mercuryReports[1] = hex"aabbccdd";

    rawReports = abi.encode(mercuryReports);
    report = abi.encodePacked(workflowId, DON_ID, EXECUTION_ID, WORKFLOW_OWNER, WORKFLOW_NAME, REPORT_NAME, rawReports);

    for (uint256 i = 0; i < requiredSignaturesNum; i++) {
      (uint8 v, bytes32 r, bytes32 s) = vm.sign(s_signers[i].mockPrivateKey, keccak256(report));
      signatures[i] = bytes.concat(r, s, bytes1(v));
    }

    vm.startPrank(TRANSMITTER);
  }

  function test_RevertWhen_ReportHasIncorrectDON() public {
    uint32 invalidDONId = 111;
    bytes memory reportWithInvalidDONId = abi.encodePacked(
      workflowId,
      invalidDONId,
      EXECUTION_ID,
      WORKFLOW_OWNER,
      rawReports
    );

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidDonId.selector, invalidDONId));
    s_forwarder.report(address(s_receiver), reportWithInvalidDONId, signatures);
  }

  function test_RevertWhen_ReportIsMalformed() public {
    bytes memory shortenedReport = abi.encode(bytes32(report));

    vm.expectRevert(KeystoneForwarder.InvalidReport.selector);
    s_forwarder.report(address(s_receiver), shortenedReport, signatures);
  }

  function test_RevertWhen_TooFewSignatures() public {
    bytes[] memory fewerSignatures = new bytes[](F);

    vm.expectRevert(
      abi.encodeWithSelector(KeystoneForwarder.InvalidSignatureCount.selector, F + 1, fewerSignatures.length)
    );
    s_forwarder.report(address(s_receiver), report, fewerSignatures);
  }

  function test_RevertWhen_TooManySignatures() public {
    bytes[] memory moreSignatures = new bytes[](F + 2);

    vm.expectRevert(
      abi.encodeWithSelector(KeystoneForwarder.InvalidSignatureCount.selector, F + 1, moreSignatures.length)
    );
    s_forwarder.report(address(s_receiver), report, moreSignatures);
  }

  function test_RevertWhen_AnySignatureIsInvalid() public {
    signatures[1] = abi.encode(1234); // invalid signature

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidSignature.selector, signatures[1]));
    s_forwarder.report(address(s_receiver), report, signatures);
  }

  function test_RevertWhen_AnySignerIsInvalid() public {
    uint256 mockPK = 999;

    Signer memory maliciousSigner = Signer({mockPrivateKey: mockPK, signerAddress: vm.addr(mockPK)});
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(maliciousSigner.mockPrivateKey, keccak256(report));
    signatures[1] = bytes.concat(r, s, bytes1(v));

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InvalidSigner.selector, maliciousSigner.signerAddress));
    s_forwarder.report(address(s_receiver), report, signatures);
  }

  function test_RevertWhen_ReportHasDuplicateSignatures() public {
    signatures[1] = signatures[0]; // repeat a signature

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.DuplicateSigner.selector, s_signers[0].signerAddress));
    s_forwarder.report(address(s_receiver), report, signatures);
  }

  function test_RevertWhen_ReportAlreadyProcessed() public {
    // First report should go through
    s_forwarder.report(address(s_receiver), report, signatures);

    bytes32 reportId = keccak256(abi.encode(address(s_receiver), EXECUTION_ID, REPORT_NAME));
    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.ReportAlreadyProcessed.selector, reportId));
    s_forwarder.report(address(s_receiver), report, signatures);
  }

  function test_Report_SuccessfulDelivery() public {
    // s_receiver.setAllowedOwnerReport(address(s_forwarder), WORKFLOW_OWNER, WORKFLOW_NAME, REPORT_NAME);

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(workflowId, WORKFLOW_OWNER, WORKFLOW_NAME, REPORT_NAME, rawReports);

    vm.expectEmit(address(s_forwarder));
    emit ReportProcessed(address(s_receiver), EXECUTION_ID, REPORT_NAME, WORKFLOW_OWNER, WORKFLOW_NAME, true);

    s_forwarder.report(address(s_receiver), report, signatures);

    // validate transmitter was recorded
    address transmitter = s_forwarder.getTransmitter(address(s_receiver), EXECUTION_ID, REPORT_NAME);
    assertEq(transmitter, TRANSMITTER, "transmitter mismatch");
  }
}
