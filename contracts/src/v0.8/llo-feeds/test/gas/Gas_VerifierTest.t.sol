// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest, BaseTestWithConfiguredVerifierAndFeeManager} from "../verifier/BaseVerifierTest.t.sol";
import {SimpleWriteAccessController} from "../../../shared/access/SimpleWriteAccessController.sol";
import {Common} from "../../libraries/Common.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

contract Verifier_setConfig is BaseTest {
  address[] internal s_signerAddrs;

  function setUp() public override {
    BaseTest.setUp();
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_signerAddrs = _getSignerAddresses(signers);
    s_verifierProxy.initializeVerifier(address(s_verifier));
  }

  function testSetConfigSuccess_gas() public {
    s_verifier.setConfig(
      FEED_ID,
      s_signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }
}

contract Verifier_verifyWithFee is BaseTestWithConfiguredVerifierAndFeeManager {
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  function setUp() public virtual override {
    super.setUp();

    //mint some link and eth to warm the storage
    link.mint(address(rewardManager), DEFAULT_LINK_MINT_QUANTITY);
    native.mint(address(feeManager), DEFAULT_NATIVE_MINT_QUANTITY);

    //warm the rewardManager
    link.mint(address(this), DEFAULT_NATIVE_MINT_QUANTITY);
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, address(this));
    (, , bytes32 latestConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some link tokens to the feeManager pool
    link.mint(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //approve funds prior to test
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](1);
    payments[0] = IRewardManager.FeePayment(latestConfigDigest, uint192(DEFAULT_REPORT_LINK_FEE));

    changePrank(address(feeManager));
    rewardManager.onFeePaid(payments, address(this));

    changePrank(USER);
  }

  function testVerifyProxyWithLinkFeeSuccess_gas() public {
    bytes memory signedLinkPayload = _generateV3EncodedBlob(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1)
    );

    s_verifierProxy.verify(signedLinkPayload, abi.encode(link));
  }

  function testVerifyProxyWithNativeFeeSuccess_gas() public {
    bytes memory signedNativePayload = _generateV3EncodedBlob(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1)
    );

    s_verifierProxy.verify(signedNativePayload, abi.encode(native));
  }
}

contract Verifier_bulkVerifyWithFee is BaseTestWithConfiguredVerifierAndFeeManager {
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;
  uint256 internal constant NUMBER_OF_REPORTS_TO_VERIFY = 5;

  function setUp() public virtual override {
    super.setUp();

    //mint some link and eth to warm the storage
    link.mint(address(rewardManager), DEFAULT_LINK_MINT_QUANTITY);
    native.mint(address(feeManager), DEFAULT_NATIVE_MINT_QUANTITY);

    //warm the rewardManager
    link.mint(address(this), DEFAULT_NATIVE_MINT_QUANTITY);
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, address(this));
    (, , bytes32 latestConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some link tokens to the feeManager pool
    link.mint(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS_TO_VERIFY);

    //approve funds prior to test
    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS_TO_VERIFY, USER);
    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS_TO_VERIFY, USER);

    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](1);
    payments[0] = IRewardManager.FeePayment(latestConfigDigest, uint192(DEFAULT_REPORT_LINK_FEE));

    changePrank(address(feeManager));
    rewardManager.onFeePaid(payments, address(this));

    changePrank(USER);
  }

  function testBulkVerifyProxyWithLinkFeeSuccess_gas() public {
    bytes memory signedLinkPayload = _generateV3EncodedBlob(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedLinkPayloads = new bytes[](NUMBER_OF_REPORTS_TO_VERIFY);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS_TO_VERIFY; i++) {
      signedLinkPayloads[i] = signedLinkPayload;
    }

    s_verifierProxy.verifyBulk(signedLinkPayloads, abi.encode(link));
  }

  function testBulkVerifyProxyWithNativeFeeSuccess_gas() public {
    bytes memory signedNativePayload = _generateV3EncodedBlob(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedNativePayloads = new bytes[](NUMBER_OF_REPORTS_TO_VERIFY);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS_TO_VERIFY; i++) {
      signedNativePayloads[i] = signedNativePayload;
    }

    s_verifierProxy.verifyBulk(signedNativePayloads, abi.encode(native));
  }
}

contract Verifier_verify is BaseTestWithConfiguredVerifierAndFeeManager {
  bytes internal s_signedReport;
  bytes32 internal s_configDigest;

  function setUp() public override {
    BaseTestWithConfiguredVerifierAndFeeManager.setUp();
    BaseTest.V1Report memory s_testReportOne = _createV1Report(
      FEED_ID,
      OBSERVATIONS_TIMESTAMP,
      MEDIAN,
      BID,
      ASK,
      BLOCKNUMBER_UPPER_BOUND,
      blockhash(BLOCKNUMBER_UPPER_BOUND),
      BLOCKNUMBER_LOWER_BOUND,
      uint32(block.timestamp)
    );
    (, , s_configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    bytes32[3] memory reportContext;
    reportContext[0] = s_configDigest;
    reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_signedReport = _generateV1EncodedBlob(s_testReportOne, reportContext, _getSigners(FAULT_TOLERANCE + 1));
  }

  function testVerifySuccess_gas() public {
    changePrank(address(s_verifierProxy));

    s_verifier.verify(s_signedReport, msg.sender);
  }

  function testVerifyProxySuccess_gas() public {
    s_verifierProxy.verify(s_signedReport, abi.encode(native));
  }
}

contract Verifier_accessControlledVerify is BaseTestWithConfiguredVerifierAndFeeManager {
  bytes internal s_signedReport;
  bytes32 internal s_configDigest;
  SimpleWriteAccessController s_accessController;

  address internal constant CLIENT = address(9000);
  address internal constant ACCESS_CONTROLLER_ADDR = address(10000);

  function setUp() public override {
    BaseTestWithConfiguredVerifierAndFeeManager.setUp();
    BaseTest.V1Report memory s_testReportOne = _createV1Report(
      FEED_ID,
      OBSERVATIONS_TIMESTAMP,
      MEDIAN,
      BID,
      ASK,
      BLOCKNUMBER_UPPER_BOUND,
      blockhash(BLOCKNUMBER_UPPER_BOUND),
      BLOCKNUMBER_LOWER_BOUND,
      uint32(block.timestamp)
    );
    (, , s_configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    bytes32[3] memory reportContext;
    reportContext[0] = s_configDigest;
    reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_signedReport = _generateV1EncodedBlob(s_testReportOne, reportContext, _getSigners(FAULT_TOLERANCE + 1));
    s_accessController = new SimpleWriteAccessController();
    s_verifierProxy.setAccessController(s_accessController);
    s_accessController.addAccess(CLIENT);
  }

  function testVerifyWithAccessControl_gas() public {
    changePrank(CLIENT);
    s_verifierProxy.verify(s_signedReport, abi.encode(native));
  }
}
