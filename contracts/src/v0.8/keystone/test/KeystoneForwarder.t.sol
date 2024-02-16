// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Test.sol";

import "../KeystoneForwarder.sol";
import {Utils} from "../libraries/Utils.sol";

contract Receiver {
  event MessageReceived(bytes32 indexed workflowId, bytes32 indexed workflowExecutionId, bytes[] mercuryReports);

  constructor() {}

  function foo(bytes calldata rawReport) external {
    // decode metadata
    (bytes32 workflowId, bytes32 workflowExecutionId) = Utils._splitReport(rawReport);
    // parse actual report
    bytes[] memory mercuryReports = abi.decode(rawReport[64:], (bytes[]));
    emit MessageReceived(workflowId, workflowExecutionId, mercuryReports);
  }
}

contract KeystoneForwarderTest is Test {
  function setUp() public virtual {}

  function test_abi_partial_decoding_works() public {
    bytes memory report = hex"0102";
    uint256 amount = 1;
    bytes memory payload = abi.encode(report, amount);
    bytes memory decodedReport = abi.decode(payload, (bytes));
    assertEq(decodedReport, report, "not equal");
  }

  function test_it_works() public {
    KeystoneForwarder forwarder = new KeystoneForwarder();
    Receiver receiver = new Receiver();

    // taken from https://github.com/smartcontractkit/chainlink/blob/2390ec7f3c56de783ef4e15477e99729f188c524/core/services/relay/evm/cap_encoder_test.go#L42-L55
    bytes
      memory report = hex"6d795f69640000000000000000000000000000000000000000000000000000006d795f657865637574696f6e5f696400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000301020300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004aabbccdd00000000000000000000000000000000000000000000000000000000";
    bytes memory data = abi.encodeWithSignature("foo(bytes)", report);
    bytes[] memory signatures = new bytes[](0);

    vm.expectCall(address(receiver), data);
    vm.recordLogs();

    bool delivered1 = forwarder.report(address(receiver), data, signatures);
    assertTrue(delivered1, "report not delivered");

    Vm.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].emitter, address(receiver));
    // validate workflow id and workflow execution id
    bytes32 workflowId = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
    bytes32 executionId = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";
    assertEq(entries[0].topics[1], workflowId);
    assertEq(entries[0].topics[2], executionId);
    bytes[] memory mercuryReports = abi.decode(entries[0].data, (bytes[]));
    assertEq(mercuryReports.length, 2);
    assertEq(mercuryReports[0], hex"010203");
    assertEq(mercuryReports[1], hex"aabbccdd");

    // doesn't deliver the same report more than once
    bool delivered2 = forwarder.report(address(receiver), data, signatures);
    assertFalse(delivered2, "report redelivered");
  }
}
