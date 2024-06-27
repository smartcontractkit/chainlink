// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifierAndFeeManager} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierTestWithConfiguredVerifierAndFeeManager is BaseTestWithConfiguredVerifierAndFeeManager {
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  function setUp() public virtual override {
    super.setUp();

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);
  }
}

contract VerifierTestBillingReport is VerifierTestWithConfiguredVerifierAndFeeManager {
  function test_verifyWithLink() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    _verify(signedReport, 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_verifyWithNative() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    _verify(signedReport, 0, USER);

    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_verifyWithNativeUnwrapped() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    _verify(signedReport, DEFAULT_REPORT_NATIVE_FEE, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyWithNativeUnwrappedReturnsChange() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    _verify(signedReport, DEFAULT_REPORT_NATIVE_FEE * 2, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(address(feeManager).balance, 0);
  }
}
