// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {VerifierWithFeeManager} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {IDestinationFeeManager} from "../../../v0.4.0/interfaces/IDestinationFeeManager.sol";
import {IDestinationRewardManager} from "../../../v0.4.0/interfaces/IDestinationRewardManager.sol";
import {IDestinationVerifierProxy} from "../../../v0.4.0/interfaces/IDestinationVerifierProxy.sol";
import {Common} from "../../../libraries/Common.sol";
import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

/*
This test checks the interfaces of destination verifier matches the expectations.
The code here comes from this example:

https://docs.chain.link/chainlink-automation/guides/streams-lookup

*/

// Custom interfaces for IVerifierProxy and IFeeManager
interface IVerifierProxy {
  /**
   * @notice Verifies that the data encoded has been signed.
   * correctly by routing to the correct verifier, and bills the user if applicable.
   * @param payload The encoded data to be verified, including the signed
   * report.
   * @param parameterPayload Fee metadata for billing. For the current implementation this is just the abi-encoded fee token ERC-20 address.
   * @return verifierResponse The encoded report from the verifier.
   */
  function verify(
    bytes calldata payload,
    bytes calldata parameterPayload
  ) external payable returns (bytes memory verifierResponse);

  function s_feeManager() external view returns (IDestinationFeeManager);
}

interface IFeeManager {
  /**
   * @notice Calculates the fee and reward associated with verifying a report, including discounts for subscribers.
   * This function assesses the fee and reward for report verification, applying a discount for recognized subscriber addresses.
   * @param subscriber The address attempting to verify the report. A discount is applied if this address
   * is recognized as a subscriber.
   * @param unverifiedReport The report data awaiting verification. The content of this report is used to
   * determine the base fee and reward, before considering subscriber discounts.
   * @param quoteAddress The payment token address used for quoting fees and rewards.
   * @return fee The fee assessed for verifying the report, with subscriber discounts applied where applicable.
   * @return reward The reward allocated to the caller for successfully verifying the report.
   * @return totalDiscount The total discount amount deducted from the fee for subscribers.
   */
  function getFeeAndReward(
    address subscriber,
    bytes memory unverifiedReport,
    address quoteAddress
  ) external returns (Common.Asset memory, Common.Asset memory, uint256);

  function i_linkAddress() external view returns (address);

  function i_nativeAddress() external view returns (address);

  function i_rewardManager() external view returns (address);
}

//Tests
// https://docs.chain.link/chainlink-automation/guides/streams-lookup
contract VerifierInterfacesTest is VerifierWithFeeManager {
  address internal constant DEFAULT_RECIPIENT_1 = address(uint160(uint256(keccak256("DEFAULT_RECIPIENT_1"))));

  IVerifierProxy public verifier;
  V3Report internal s_testReport;

  address public FEE_ADDRESS;
  string public constant DATASTREAMS_FEEDLABEL = "feedIDs";
  string public constant DATASTREAMS_QUERYLABEL = "timestamp";
  int192 public last_retrieved_price;
  bytes internal signedReport;
  bytes32[3] internal s_reportContext;
  uint8 MINIMAL_FAULT_TOLERANCE = 2;

  function setUp() public virtual override {
    VerifierWithFeeManager.setUp();
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    s_testReport = V3Report({
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
    address[] memory signerAddrs = _getSignerAddresses(signers);
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, ONE_PERCENT * 100);
    s_verifier.setConfig(signerAddrs, MINIMAL_FAULT_TOLERANCE, weights);
    signedReport = _generateV3EncodedBlob(s_testReport, s_reportContext, signers);

    verifier = IVerifierProxy(address(s_verifierProxy));
  }

  function test_DestinationContractInterfaces() public {
    bytes memory unverifiedReport = signedReport;

    (, bytes memory reportData) = abi.decode(unverifiedReport, (bytes32[3], bytes));

    // Report verification fees
    IFeeManager feeManager = IFeeManager(address(verifier.s_feeManager()));
    IDestinationRewardManager rewardManager = IDestinationRewardManager(address(feeManager.i_rewardManager()));

    address feeTokenAddress = feeManager.i_linkAddress();
    (Common.Asset memory fee, , ) = feeManager.getFeeAndReward(address(this), reportData, feeTokenAddress);

    // Approve rewardManager to spend this contract's balance in fees
    _approveLink(address(rewardManager), fee.amount, USER);
    _verify(unverifiedReport, address(feeTokenAddress), 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - fee.amount);
    assertEq(link.balanceOf(address(rewardManager)), fee.amount);
  }
}
