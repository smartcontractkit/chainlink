// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../FeeManager.sol";
import {RewardManager} from "../../RewardManager.sol";
import {Common} from "../../../libraries/internal/Common.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice Base class for all fee manager tests
 * @dev This contract is intended to be inherited from and not used directly. It contains functionality to setup the fee manager
 */
contract BaseFeeManagerTest is Test {
  //contracts
  FeeManager internal feeManager;
  RewardManager internal rewardManager;

  //contract owner
  address internal constant INVALID_ADDRESS = address(0);
  address internal constant ADMIN = address(1);
  address internal constant USER = address(2);

  address internal constant LINK_ADDRESS = address(3);
  address internal constant NATIVE_ADDRESS = address(4);

  //feed ids
  bytes32 internal constant DEFAULT_FEED_1 = keccak256("feed_id_1");
  bytes32 internal constant DEFAULT_FEED_2 = keccak256("feed_id_2");

  //report
  uint256 internal constant DEFAULT_REPORT_LINK_FEE = 1e10;
  uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e12;
  uint256 internal constant DEFAULT_REPORT_EXPIRY_OFFSET_SECONDS = 300;

  //rewards
  uint256 internal constant FEE_SCALAR = 1e18;

  //the selector for each error
  bytes4 internal constant INVALID_DISCOUNT_ERROR = bytes4(keccak256("InvalidDiscount()"));
  bytes4 internal constant INVALID_ADDRESS_ERROR = bytes4(keccak256("InvalidAddress()"));
  bytes4 internal constant INVALID_PREMIUM_ERROR = bytes4(keccak256("InvalidPremium()"));
  bytes4 internal constant EXPIRED_REPORT_ERROR = bytes4(keccak256("ExpiredReport()"));
  bytes internal constant ONLY_CALLABLE_BY_OWNER_ERROR = "Only callable by owner";

  //events emitted
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint256 discount);
  event NativePremiumSet(uint256 newPremium);

  function setUp() public virtual {
    //change to admin user
    vm.startPrank(ADMIN);

    //init required contracts
    _initializeContracts();
  }

  function _initializeContracts() internal {
    rewardManager = new RewardManager(LINK_ADDRESS);
    feeManager = new FeeManager(LINK_ADDRESS, NATIVE_ADDRESS, USER, address(rewardManager));
    rewardManager.setFeeManager(address(feeManager));
  }

  function setSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint256 discount,
    address sender
  ) internal {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the discount
    feeManager.updateSubscriberDiscount(subscriber, feedId, token, discount);

    //change back to the original address
    changePrank(originalAddr);
  }

  function setNativePremium(uint256 premium, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the premium
    feeManager.setNativePremium(premium);

    //change back to the original address
    changePrank(originalAddr);
  }

  // solium-disable-next-line no-unused-vars
  function getFee(
    bytes memory report,
    FeeManager.Quote memory quote,
    address subscriber
  ) public view returns (Common.Asset memory) {
    //set the discount
    (Common.Asset memory fee, ) = feeManager.getFeeAndReward(subscriber, report, quote);

    return fee;
  }

  function getReward(
    bytes memory report,
    FeeManager.Quote memory quote,
    address subscriber
  ) public view returns (Common.Asset memory) {
    //set the discount
    (, Common.Asset memory reward) = feeManager.getFeeAndReward(subscriber, report, quote);

    return reward;
  }

  function getReport(bytes32 feedId) public pure returns (bytes memory) {
    return abi.encode(feedId, uint32(0), int192(0), int192(0), int192(0), uint64(0), bytes32(0), uint64(0));
  }

  function getReportWithFee(bytes32 feedId) public view returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        int192(0),
        int192(0),
        int192(0),
        uint64(0),
        bytes32(0),
        uint64(0),
        DEFAULT_REPORT_LINK_FEE,
        DEFAULT_REPORT_NATIVE_FEE,
        uint32(block.timestamp)
      );
  }

  function getReportWithCustomExpiryAndFee(
    bytes32 feedId,
    uint256 expiry,
    uint256 linkFee,
    uint256 nativeFee
  ) public view returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        int192(0),
        int192(0),
        int192(0),
        uint64(0),
        bytes32(0),
        uint64(0),
        linkFee,
        nativeFee,
        uint32(expiry)
      );
  }

  function getLinkQuote() public pure returns (FeeManager.Quote memory) {
    return FeeManager.Quote(LINK_ADDRESS);
  }

  function getNativeQuote() public pure returns (FeeManager.Quote memory) {
    return FeeManager.Quote(NATIVE_ADDRESS);
  }
}
