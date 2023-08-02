// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifierAndFeeManager} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {ERC20Mock} from "../../../shared/vendor/ERC20Mock.sol";
import {WERC20Mock} from "../../../shared/vendor/WERC20Mock.sol";
import {Common} from "../../../libraries/internal/Common.sol";

contract VerifierTestWithConfiguredVerifierAndFeeManager is BaseTestWithConfiguredVerifierAndFeeManager {

    uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
    uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

    uint256 internal constant DEFAULT_REPORT_LINK_FEE = 1e10;
    uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e12;

    struct BillingReport {
        // The feed ID the report has data for
        bytes32 feedId;
        // The time the median value was observed on
        uint32 observationsTimestamp;
        // The median value agreed in an OCR round
        int192 median;
        // The best bid value agreed in an OCR round
        int192 bid;
        // The best ask value agreed in an OCR round
        int192 ask;
        // The upper bound of the block range the median value was observed within
        uint64 blocknumberUpperBound;
        // The blockhash for the upper bound of block range (ensures correct blockchain)
        bytes32 upperBlockhash;
        // The lower bound of the block range the median value was observed within
        uint64 blocknumberLowerBound;
        // The current block timestamp
        uint64 currentBlockTimestamp;
        // The link fee
        uint192 linkFee;
        // The native fee
        uint192 nativeFee;
        // The expiry of the report
        uint32 expiresAt;
    }

    function setUp() public virtual override {
        super.setUp();

        //mint some tokens to the user
        link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
        native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
        vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    }

    function _encodeReport(BillingReport memory report) internal pure returns (bytes memory) {
        return abi.encode(
            report.feedId,
            report.observationsTimestamp,
            report.median,
            report.bid,
            report.ask,
            report.blocknumberUpperBound,
            report.upperBlockhash,
            report.blocknumberLowerBound,
            report.currentBlockTimestamp,
            report.linkFee,
            report.nativeFee,
            report.expiresAt
        );
    }

    function _generateEncodedBlobWithFeesAndQuote(
        BillingReport memory report,
        bytes32[3] memory reportContext,
        Signer[] memory signers,
        bytes memory quote
    ) internal pure returns (bytes memory) {
        bytes memory reportBytes = _encodeReport(report);
        (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _generateSignerSignatures(
            reportBytes,
            reportContext,
            signers
        );

        return abi.encode(reportContext, reportBytes, rs, ss, rawVs, quote);
    }

    function _generateQuote(address billingAddress) internal returns (bytes memory) {
        return abi.encode(billingAddress);
    }

    function _generateBillingReport() internal returns (BillingReport memory) {
        return BillingReport({
            feedId: FEED_ID,
            observationsTimestamp: OBSERVATIONS_TIMESTAMP,
            median: MEDIAN,
            bid: BID,
            ask: ASK,
            blocknumberUpperBound: BLOCKNUMBER_UPPER_BOUND,
            upperBlockhash: blockhash(BLOCKNUMBER_UPPER_BOUND),
            blocknumberLowerBound: BLOCKNUMBER_LOWER_BOUND,
            currentBlockTimestamp: uint64(block.timestamp),
            linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
            nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
            expiresAt: uint32(block.timestamp)
        });
    }


    function _generateReportContext() internal returns (bytes32[3] memory) {
        (, , bytes32 latestConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
        bytes32[3] memory reportContext;
        reportContext[0] = latestConfigDigest;
        reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
        return reportContext;
    }

    function _approveLink(address spender, uint256 quantity, address sender) internal {
        address originalAddr = msg.sender;
        changePrank(sender);

        link.approve(spender, quantity);
        changePrank(originalAddr);
    }

    function _approveNative(address spender, uint256 quantity, address sender) internal {
        address originalAddr = msg.sender;
        changePrank(sender);

        native.approve(spender, quantity);
        changePrank(originalAddr);
    }

    function _verify(bytes memory payload, uint256 wrappedNativeValue, address sender) internal {
        address originalAddr = msg.sender;
        changePrank(sender);

        s_verifierProxy.verify{value: wrappedNativeValue}(payload);

        changePrank(originalAddr);
    }
}

contract VerifierTestBillingReport is VerifierTestWithConfiguredVerifierAndFeeManager {

    function test_verifyWithLink() public {
        bytes memory signedReport = _generateEncodedBlobWithFeesAndQuote(
            _generateBillingReport(),
            _generateReportContext(),
            _getSigners(FAULT_TOLERANCE + 1),
            _generateQuote(address(link))
        );

        _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

        _verify(signedReport, 0, USER);

        assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
    }

    function test_verifyWithNative() public {
        bytes memory signedReport = _generateEncodedBlobWithFeesAndQuote(
            _generateBillingReport(),
            _generateReportContext(),
            _getSigners(FAULT_TOLERANCE + 1),
            _generateQuote(address(native))
        );

        _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

        _verify(signedReport, 0, USER);

        assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    }

    function test_verifyWithNativeUnwrapped() public {
        bytes memory signedReport = _generateEncodedBlobWithFeesAndQuote(
            _generateBillingReport(),
            _generateReportContext(),
            _getSigners(FAULT_TOLERANCE + 1),
            _generateQuote(address(native))
        );

        _verify(signedReport, DEFAULT_REPORT_NATIVE_FEE, USER);

        assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
        assertEq(address(feeManager).balance, 0);
    }

    function test_verifyWithNativeUnwrappedReturnsChange() public {
        bytes memory signedReport = _generateEncodedBlobWithFeesAndQuote(
            _generateBillingReport(),
            _generateReportContext(),
            _getSigners(FAULT_TOLERANCE + 1),
            _generateQuote(address(native))
        );

        _verify(signedReport, DEFAULT_REPORT_NATIVE_FEE * 2, USER);

        assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
        assertEq(address(feeManager).balance, 0);
    }
}


