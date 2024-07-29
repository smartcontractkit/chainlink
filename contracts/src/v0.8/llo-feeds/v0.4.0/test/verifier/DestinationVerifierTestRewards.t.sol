// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {VerifierWithFeeManager} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierBillingTests is VerifierWithFeeManager {
  uint8 MINIMAL_FAULT_TOLERANCE = 2;
  address internal constant DEFAULT_RECIPIENT_1 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_1"))));
  address internal constant DEFAULT_RECIPIENT_2 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_2"))));
  address internal constant DEFAULT_RECIPIENT_3 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_3"))));
  address internal constant DEFAULT_RECIPIENT_4 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_4"))));
  address internal constant DEFAULT_RECIPIENT_5 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_5"))));
  address internal constant DEFAULT_RECIPIENT_6 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_6"))));
  address internal constant DEFAULT_RECIPIENT_7 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_7"))));

  function payRecipients(bytes32 poolId, address[] memory recipients, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //pay the recipients
    rewardManager.payRecipients(poolId, recipients);

    //change back to the original address
    changePrank(originalAddr);
  }

  bytes32[3] internal s_reportContext;
  V3Report internal s_testReport;

  function setUp() public virtual override {
    VerifierWithFeeManager.setUp();
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReport = generateReportAtTimestamp(block.timestamp);
  }

  function generateReportAtTimestamp(uint256 timestamp) public pure returns (V3Report memory) {
    return
      V3Report({
        feedId: FEED_ID_V3,
        observationsTimestamp: OBSERVATIONS_TIMESTAMP,
        validFromTimestamp: uint32(timestamp),
        nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
        linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
        // ask michael about this expires at, is it usually set at what blocks
        expiresAt: uint32(timestamp) + 500,
        benchmarkPrice: MEDIAN,
        bid: BID,
        ask: ASK
      });
  }

  function getRecipientAndWeightsGroup2() public pure returns (Common.AddressAndWeight[] memory, address[] memory) {
    address[] memory recipients = new address[](4);
    recipients[0] = DEFAULT_RECIPIENT_4;
    recipients[1] = DEFAULT_RECIPIENT_5;
    recipients[2] = DEFAULT_RECIPIENT_6;
    recipients[3] = DEFAULT_RECIPIENT_7;

    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](4);
    //init each recipient with even weights. 2500 = 25% of pool
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, POOL_SCALAR / 4);
    weights[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, POOL_SCALAR / 4);
    weights[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, POOL_SCALAR / 4);
    weights[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, POOL_SCALAR / 4);
    return (weights, recipients);
  }

  function getRecipientAndWeightsGroup1() public pure returns (Common.AddressAndWeight[] memory, address[] memory) {
    address[] memory recipients = new address[](4);
    recipients[0] = DEFAULT_RECIPIENT_1;
    recipients[1] = DEFAULT_RECIPIENT_2;
    recipients[2] = DEFAULT_RECIPIENT_3;
    recipients[3] = DEFAULT_RECIPIENT_4;

    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](4);
    //init each recipient with even weights. 2500 = 25% of pool
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, POOL_SCALAR / 4);
    weights[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, POOL_SCALAR / 4);
    weights[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, POOL_SCALAR / 4);
    weights[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, POOL_SCALAR / 4);
    return (weights, recipients);
  }

  function test_rewardsAreDistributedAccordingToWeights() public {
    /*
          Simple test verifying that rewards are distributed according to address and weights 
          associated to the DonConfig used to verify the report
         */
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, ONE_PERCENT * 100);
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    bytes memory signedReport = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);
    bytes32 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(signedReport, address(link), 0, USER);
    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);

    // internal state checks
    assertEq(feeManager.s_linkDeficit(expectedDonConfigId), 0);
    assertEq(rewardManager.s_totalRewardRecipientFees(expectedDonConfigId), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    // check the recipients are paid according to weights
    address[] memory recipients = new address[](1);
    recipients[0] = DEFAULT_RECIPIENT_1;
    payRecipients(expectedDonConfigId, recipients, ADMIN);
    assertEq(link.balanceOf(recipients[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);
  }

  function test_rewardsAreDistributedAccordingToWeightsMultipleWeigths() public {
    /*
          Rewards are distributed according to AddressAndWeight's 
          associated to the DonConfig used to verify the report:
          - multiple recipients
          - multiple verifications
         */
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    (Common.AddressAndWeight[] memory weights, address[] memory recipients) = getRecipientAndWeightsGroup1();
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);

    bytes memory signedReport = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);
    bytes32 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

    uint256 number_of_reports_verified = 10;

    for (uint256 i = 0; i < number_of_reports_verified; i++) {
      _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
      _verify(signedReport, address(link), 0, USER);
    }

    uint256 expected_pool_amount = DEFAULT_REPORT_LINK_FEE * number_of_reports_verified;

    //each recipient should receive 1/4 of the pool
    uint256 expectedRecipientAmount = expected_pool_amount / 4;

    payRecipients(expectedDonConfigId, recipients, ADMIN);
    for (uint256 i = 0; i < recipients.length; i++) {
      // checking each recipient got rewards as set by the weights
      assertEq(link.balanceOf(recipients[i]), expectedRecipientAmount);
    }
    // checking nothing left in reward manager
    assertEq(link.balanceOf(address(rewardManager)), 0);
  }

  function test_rewardsAreDistributedAccordingToWeightsUsingHistoricalConfigs() public {
    /*
          Verifies that reports verified with historical give rewards according to the verifying config AddressAndWeight.
          - Sets two Configs: ConfigA and ConfigB, These two Configs have different Recipient and Weights 
          - Verifies a couple reports with each config
          - Pays recipients
          - Asserts expected rewards for each recipient 
         */

    Signer[] memory signers = _getSigners(10);
    address[] memory signerAddrs = _getSignerAddresses(signers);

    (Common.AddressAndWeight[] memory weights, address[] memory recipients) = getRecipientAndWeightsGroup1();

    // Create ConfigA
    s_verifier.setConfig(signerAddrs, MINIMAL_FAULT_TOLERANCE, weights);
    vm.warp(block.timestamp + 100);

    V3Report memory testReportAtT1 = generateReportAtTimestamp(block.timestamp);
    bytes memory signedReportT1 = _generateV3EncodedBlob(testReportAtT1, s_reportContext, signers);
    bytes32 expectedDonConfigIdA = _donConfigIdFromConfigData(signerAddrs, MINIMAL_FAULT_TOLERANCE);

    uint256 number_of_reports_verified = 2;

    // advancing the blocktimestamp so we can test verifying with configs
    vm.warp(block.timestamp + 100);

    Signer[] memory signers2 = _getSigners(12);
    address[] memory signerAddrs2 = _getSignerAddresses(signers2);
    (Common.AddressAndWeight[] memory weights2, address[] memory recipients2) = getRecipientAndWeightsGroup2();

    // Create ConfigB
    s_verifier.setConfig(signerAddrs2, MINIMAL_FAULT_TOLERANCE, weights2);
    bytes32 expectedDonConfigIdB = _donConfigIdFromConfigData(signerAddrs2, MINIMAL_FAULT_TOLERANCE);

    V3Report memory testReportAtT2 = generateReportAtTimestamp(block.timestamp);

    // verifiying using ConfigA (report with Old timestamp)
    for (uint256 i = 0; i < number_of_reports_verified; i++) {
      _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
      _verify(signedReportT1, address(link), 0, USER);
    }

    // verifying using ConfigB (report with new timestamp)
    for (uint256 i = 0; i < number_of_reports_verified; i++) {
      _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
      _verify(_generateV3EncodedBlob(testReportAtT2, s_reportContext, signers2), address(link), 0, USER);
    }

    uint256 expected_pool_amount = DEFAULT_REPORT_LINK_FEE * number_of_reports_verified;
    assertEq(rewardManager.s_totalRewardRecipientFees(expectedDonConfigIdA), expected_pool_amount);
    assertEq(rewardManager.s_totalRewardRecipientFees(expectedDonConfigIdB), expected_pool_amount);

    // check the recipients are paid according to weights
    payRecipients(expectedDonConfigIdA, recipients, ADMIN);

    for (uint256 i = 0; i < recipients.length; i++) {
      // //each recipient should receive 1/4 of the pool
      assertEq(link.balanceOf(recipients[i]), expected_pool_amount / 4);
    }

    payRecipients(expectedDonConfigIdB, recipients2, ADMIN);

    for (uint256 i = 1; i < recipients2.length; i++) {
      // //each recipient should receive 1/4 of the pool
      assertEq(link.balanceOf(recipients2[i]), expected_pool_amount / 4);
    }

    // this recipient was part of the two config weights
    assertEq(link.balanceOf(recipients2[0]), (expected_pool_amount / 4) * 2);
    assertEq(link.balanceOf(address(rewardManager)), 0);
  }
}
