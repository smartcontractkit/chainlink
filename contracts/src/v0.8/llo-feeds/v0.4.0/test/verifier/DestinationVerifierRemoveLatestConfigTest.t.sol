// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationRewardManager} from "../../../v0.4.0/DestinationRewardManager.sol";
import {Common} from "../../../libraries/Common.sol";

contract DestinationVerifierSetConfigTest is BaseTest {
  bytes32[3] internal s_reportContext;
  V3Report internal s_testReport;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
  }

  function test_removeLatestConfigWhenNoConfigShouldFail() public {
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.DonConfigDoesNotExist.selector));
    s_verifier.removeLatestConfig();
  }

  function test_removeLatestConfig() public {
    /*
       This test sets two Configs: Config A and Config B.
       - it removes and readds config B multiple times while trying Config A verifications
      */
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersA = new BaseTest.Signer[](7);
    signersA[0] = signers[0];
    signersA[1] = signers[1];
    signersA[2] = signers[2];
    signersA[3] = signers[3];
    signersA[4] = signers[4];
    signersA[5] = signers[5];
    signersA[6] = signers[6];

    // ConfigA
    address[] memory signersAddrA = _getSignerAddresses(signersA);
    s_verifier.setConfig(signersAddrA, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    vm.warp(block.timestamp + 10);
    V3Report memory s_testReportA = V3Report({
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

    vm.warp(block.timestamp + 100);
    // Config B
    BaseTest.Signer[] memory signersB = new BaseTest.Signer[](7);
    // signers in ConfigA
    signersB[0] = signers[8];
    signersB[1] = signers[9];
    signersB[2] = signers[10];
    signersB[3] = signers[11];
    signersB[4] = signers[12];
    signersB[5] = signers[13];
    signersB[6] = signers[14];
    address[] memory signersAddrsB = _getSignerAddresses(signersB);
    s_verifier.setConfig(signersAddrsB, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    V3Report memory s_testReportB = V3Report({
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

    BaseTest.Signer[] memory reportSignersA = new BaseTest.Signer[](3);
    reportSignersA[0] = signers[0];
    reportSignersA[1] = signers[1];
    reportSignersA[2] = signers[2];

    BaseTest.Signer[] memory reportSignersB = new BaseTest.Signer[](3);
    reportSignersB[0] = signers[8];
    reportSignersB[1] = signers[9];
    reportSignersB[2] = signers[10];

    bytes memory signedReportA = _generateV3EncodedBlob(s_testReportA, s_reportContext, reportSignersA);
    bytes memory signedReportB = _generateV3EncodedBlob(s_testReportB, s_reportContext, reportSignersB);

    // verifying should work
    s_verifierProxy.verify(signedReportA, abi.encode(native));
    s_verifierProxy.verify(signedReportB, abi.encode(native));

    s_verifier.removeLatestConfig();

    // this should remove the latest config, so ConfigA should be able to verify reports still
    s_verifierProxy.verify(signedReportA, abi.encode(native));
    // this report cannot be verified any longer because ConfigB is not there
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReportB, abi.encode(native));

    // since ConfigB is removed we should be able to set it again with no errors
    s_verifier.setConfig(signersAddrsB, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    // we should be able to remove ConfigB
    s_verifier.removeLatestConfig();
    // removing configA
    s_verifier.removeLatestConfig();

    // verifigny should fail
    // verifying should work
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReportA, abi.encode(native));
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReportB, abi.encode(native));

    // removing again should fail. no other configs exist
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.DonConfigDoesNotExist.selector));
    s_verifier.removeLatestConfig();
  }
}
