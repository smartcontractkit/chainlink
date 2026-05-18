// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {MultipleVerifierWithMultipleFeeManagers} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";

contract MultiVerifierBillingTests is MultipleVerifierWithMultipleFeeManagers {
  uint8 MINIMAL_FAULT_TOLERANCE = 2;
  address internal constant DEFAULT_RECIPIENT_1 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_1"))));
  address internal constant DEFAULT_RECIPIENT_2 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_2"))));
  address internal constant DEFAULT_RECIPIENT_3 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_3"))));
  address internal constant DEFAULT_RECIPIENT_4 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_4"))));
  address internal constant DEFAULT_RECIPIENT_5 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_5"))));
  address internal constant DEFAULT_RECIPIENT_6 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_6"))));
  address internal constant DEFAULT_RECIPIENT_7 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_7"))));

  bytes32[3] internal s_reportContext;
  V3Report internal s_testReport;

  function setUp() public virtual override {
    MultipleVerifierWithMultipleFeeManagers.setUp();
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReport = generateReportAtTimestamp(block.timestamp);
  }

  function _verify(
    DestinationVerifierProxy proxy,
    bytes memory payload,
    address feeAddress,
    uint256 wrappedNativeValue,
    address sender
  ) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    proxy.verify{value: wrappedNativeValue}(payload, abi.encode(feeAddress));

    changePrank(originalAddr);
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

  function payRecipients(bytes32 poolId, address[] memory recipients, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //pay the recipients
    rewardManager.payRecipients(poolId, recipients);

    //change back to the original address
    changePrank(originalAddr);
  }

  function test_multipleFeeManagersAndVerifiers() public {
    /*
       In this test we got:
        - three verifiers (verifier, verifier2, verifier3).
        - two fee managers (feeManager, feeManager2)
        - one reward manager
        
       we glue:
       - feeManager is used by verifier1 and verifier2
       - feeManager is used by verifier3
       - Rewardmanager is used by feeManager and feeManager2
      
      In this test we do verificatons via verifier1, verifier2 and verifier3 and check that rewards are set accordingly
   
    */
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, ONE_PERCENT * 100);

    Common.AddressAndWeight[] memory weights2 = new Common.AddressAndWeight[](1);
    weights2[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, ONE_PERCENT * 100);

    Common.AddressAndWeight[] memory weights3 = new Common.AddressAndWeight[](1);
    weights3[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, ONE_PERCENT * 100);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    s_verifier2.setConfig(signerAddrs, MINIMAL_FAULT_TOLERANCE, weights2);
    s_verifier3.setConfig(signerAddrs, MINIMAL_FAULT_TOLERANCE + 1, weights3);
    bytes memory signedReport = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);
    bytes32 expectedDonConfigID = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);
    bytes32 expectedDonConfigID2 = _donConfigIdFromConfigData(signerAddrs, MINIMAL_FAULT_TOLERANCE);
    bytes32 expectedDonConfigID3 = _donConfigIdFromConfigData(signerAddrs, MINIMAL_FAULT_TOLERANCE + 1);

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(s_verifierProxy, signedReport, address(link), 0, USER);
    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);

    // internal state checks
    assertEq(feeManager.s_linkDeficit(expectedDonConfigID), 0);
    assertEq(rewardManager.s_totalRewardRecipientFees(expectedDonConfigID), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    // check the recipients are paid according to weights
    // These rewards happened through verifier1 and feeManager1
    address[] memory recipients = new address[](1);
    recipients[0] = DEFAULT_RECIPIENT_1;
    payRecipients(expectedDonConfigID, recipients, ADMIN);
    assertEq(link.balanceOf(recipients[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);

    // these rewards happaned through verifier2 and feeManager1
    address[] memory recipients2 = new address[](1);
    recipients2[0] = DEFAULT_RECIPIENT_2;
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(s_verifierProxy2, signedReport, address(link), 0, USER);
    payRecipients(expectedDonConfigID2, recipients2, ADMIN);
    assertEq(link.balanceOf(recipients2[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);

    // these rewards happened through verifier3 and feeManager2
    address[] memory recipients3 = new address[](1);
    recipients3[0] = DEFAULT_RECIPIENT_3;
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _verify(s_verifierProxy3, signedReport, address(link), 0, USER);
    payRecipients(expectedDonConfigID3, recipients3, ADMIN);
    assertEq(link.balanceOf(recipients3[0]), DEFAULT_REPORT_LINK_FEE);
    assertEq(link.balanceOf(address(rewardManager)), 0);
  }
}
