// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./KeystoneForwarderBaseTest.t.sol";
import {KeystoneForwarder} from "../KeystoneForwarder.sol";

contract KeystoneForwarder_ReportTest is BaseTest {
  event MessageReceived(bytes32 indexed workflowId, address indexed workflowOwner, bytes[] mercuryReports);
  event ReportProcessed(
    address indexed receiver,
    address indexed workflowOwner,
    bytes32 indexed workflowExecutionId,
    bool result
  );

  bytes32 internal workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  address internal workflowOwner = address(51);
  bytes32 internal executionId = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";
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
    report = abi.encodePacked(workflowId, DON_ID, executionId, workflowOwner, rawReports);

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
      executionId,
      workflowOwner,
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
    s_forwarder.report(address(s_receiver), report, signatures);
    bytes32 reportId = keccak256(bytes.concat(bytes20(uint160(address(s_receiver))), executionId));

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.ReportAlreadyProcessed.selector, reportId));
    s_forwarder.report(address(s_receiver), report, signatures);
  }

  function test_Report_SuccessfulDelivery() public {
    // taken from https://github.com/smartcontractkit/chainlink/blob/2390ec7f3c56de783ef4e15477e99729f188c524/core/services/relay/evm/cap_encoder_test.go#L42-L55
    // bytes memory report = hex"6d795f6964000000000000000000000000000000000000000000000000000000010203046d795f657865637574696f6e5f696400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000301020300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004aabbccdd00000000000000000000000000000000000000000000000000000000";

    vm.expectEmit(address(s_receiver));
    emit MessageReceived(workflowId, workflowOwner, mercuryReports);

    vm.expectEmit(address(s_forwarder));
    emit ReportProcessed(address(s_receiver), workflowOwner, executionId, true);

    s_forwarder.report(address(s_receiver), report, signatures);

    // validate transmitter was recorded
    address transmitter = s_forwarder.getTransmitter(address(s_receiver), executionId);
    assertEq(transmitter, TRANSMITTER, "transmitter mismatch");
  }
}
