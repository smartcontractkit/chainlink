// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../FeeManager.sol";
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
  uint256 internal constant DEFAULT_REPORT_SIZE_EXCLUDING_FEED_ID = 4 + 24 + 24 + 24 + 8 + 32 + 8;
  uint256 internal constant DEFAULT_REPORT_LINK_FEE = 1e10;
  uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e10;

  //rewards
  uint16 internal constant FEE_SCALAR = 10000;

  //the selector for each error
  bytes4 internal constant INVALID_DISCOUNT_ERROR = bytes4(keccak256("InvalidDiscount()"));
  bytes4 internal constant INVALID_ADDRESS_ERROR = bytes4(keccak256("InvalidAddress()"));
  bytes4 internal constant INVALID_PREMIUM_ERROR = bytes4(keccak256("InvalidPremium()"));
  bytes internal constant ONLY_CALLABLE_BY_OWNER_ERROR = "Only callable by owner";

  //events emitted
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint16 discount);
  event NativePremiumSet(uint16 newPremium);

  function setUp() public virtual {
    //change to admin user
    vm.startPrank(ADMIN);

    //init required contracts
    _initializeFeeManager();
  }

  function _initializeFeeManager() internal {
    //create the contract
    feeManager = new FeeManager(LINK_ADDRESS, NATIVE_ADDRESS);
  }

  function setSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint16 discount,
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

  function setNativePremium(uint16 premium, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the premium
    feeManager.setNativePremium(premium);

    //change back to the original address
    changePrank(originalAddr);
  }

  function getFee(bytes memory report, bytes memory quote, address subscriber) public returns (Common.Asset memory) {
    //set the discount
    return feeManager.getFee(subscriber, report, quote);
  }

  function getReport(bytes32 feedId) public pure returns (bytes memory) {
    return abi.encodePacked(feedId, new bytes(DEFAULT_REPORT_SIZE_EXCLUDING_FEED_ID));
  }

  function getReportWithFee(bytes32 feedId) public view returns (bytes memory) {
    return abi.encodePacked(getReport(feedId), bytes32(DEFAULT_REPORT_LINK_FEE), bytes32(DEFAULT_REPORT_NATIVE_FEE));
  }

  function getReportWithCustomFee(
    bytes32 feedId,
    uint256 linkFee,
    uint256 nativeFee
  ) public view returns (bytes memory) {
    return abi.encodePacked(getReport(feedId), bytes32(linkFee), bytes32(nativeFee));
  }

  function getLinkQuote() public pure returns (bytes memory) {
    return abi.encodePacked(LINK_ADDRESS);
  }

  function getNativeQuote() public pure returns (bytes memory) {
    return abi.encodePacked(NATIVE_ADDRESS);
  }
}
