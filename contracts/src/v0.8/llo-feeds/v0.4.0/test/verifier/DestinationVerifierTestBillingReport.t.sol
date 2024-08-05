// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {VerifierWithFeeManager} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierBillingTests is VerifierWithFeeManager {
  bytes32[3] internal s_reportContext;
  V3Report internal s_testReportThree;

  function setUp() public virtual override {
    VerifierWithFeeManager.setUp();
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReportThree = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(block.timestamp),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });
  }

  function test_verifyWithLinkV3Report() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);
    bytes32 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(signedReport, address(link), 0, USER);
    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);

    // internal state checks
    assertEq(feeManager.s_linkDeficit(expectedDonConfigId), 0);
    assertEq(rewardManager.s_totalRewardRecipientFees(expectedDonConfigId), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);
  }

  function test_verifyWithNativeERC20() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(signerAddrs[0], ONE_PERCENT * 100);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    bytes memory signedReport = _generateV3EncodedBlob(
      s_testReportThree,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);
    _verify(signedReport, address(native), 0, USER);
    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);

    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);
  }

  function test_verifyWithNativeUnwrapped() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    bytes memory signedReport = _generateV3EncodedBlob(
      s_testReportThree,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    _verify(signedReport, address(native), DEFAULT_REPORT_NATIVE_FEE, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyWithNativeUnwrappedReturnsChange() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    bytes memory signedReport = _generateV3EncodedBlob(
      s_testReportThree,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    _verify(signedReport, address(native), DEFAULT_REPORT_NATIVE_FEE * 2, USER);
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(address(feeManager).balance, 0);
  }
}

contract DestinationVerifierBulkVerifyBillingReport is VerifierWithFeeManager {
  uint256 internal constant NUMBERS_OF_REPORTS = 5;

  bytes32[3] internal s_reportContext;

  function setUp() public virtual override {
    VerifierWithFeeManager.setUp();
    // setting a DonConfig we can reuse in the rest of tests
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](0);
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
  }

  function test_verifyWithBulkLink() public {
    bytes memory signedReport = _generateV3EncodedBlob(
      _generateV3Report(),
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i = 0; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBERS_OF_REPORTS, USER);

    _verifyBulk(signedReports, address(link), 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * NUMBERS_OF_REPORTS);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBERS_OF_REPORTS);
  }

  function test_verifyWithBulkNative() public {
    bytes memory signedReport = _generateV3EncodedBlob(
      _generateV3Report(),
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i = 0; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS, USER);
    _verifyBulk(signedReports, address(native), 0, USER);
    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS);
  }

  function test_verifyWithBulkNativeUnwrapped() public {
    bytes memory signedReport = _generateV3EncodedBlob(
      _generateV3Report(),
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _verifyBulk(signedReports, address(native), 200 * DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyWithBulkNativeUnwrappedReturnsChange() public {
    bytes memory signedReport = _generateV3EncodedBlob(
      _generateV3Report(),
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i = 0; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _verifyBulk(signedReports, address(native), DEFAULT_REPORT_NATIVE_FEE * (NUMBERS_OF_REPORTS * 2), USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS);
    assertEq(address(feeManager).balance, 0);
  }
}
