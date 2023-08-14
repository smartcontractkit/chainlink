// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../dev/FeeManager.sol";
import {IFeeManager} from "../../dev/interfaces/IFeeManager.sol";
import {RewardManager} from "../../dev/RewardManager.sol";
import {Common} from "../../../libraries/Common.sol";
import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/mocks/ERC20Mock.sol";
import {WERC20Mock} from "../../../shared/mocks/WERC20Mock.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice Base class for all feeManager tests
 * @dev This contract is intended to be inherited from and not used directly. It contains functionality to setup the feeManager
 */
contract BaseFeeManagerTest is Test {
  //contracts
  FeeManager internal feeManager;
  RewardManager internal rewardManager;

  ERC20Mock internal link;
  WERC20Mock internal native;

  //erc20 config
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  //contract owner
  address internal constant INVALID_ADDRESS = address(0);
  address internal constant ADMIN = address(uint160(uint256(keccak256("ADMIN"))));
  address internal constant USER = address(uint160(uint256(keccak256("USER"))));
  address internal constant PROXY = address(uint160(uint256(keccak256("PROXY"))));

  //version masks
  bytes32 internal constant V_MASK = 0x0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
  bytes32 internal constant V1_BITMASK = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V2_BITMASK = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V3_BITMASK = 0x0003000000000000000000000000000000000000000000000000000000000000;

  //feed ids & config digests
  bytes32 internal constant DEFAULT_FEED_1_V1 = (keccak256("ETH-USD") & V_MASK) | V1_BITMASK;
  bytes32 internal constant DEFAULT_FEED_1_V2 = (keccak256("ETH-USD") & V_MASK) | V2_BITMASK;
  bytes32 internal constant DEFAULT_FEED_1_V3 = (keccak256("ETH-USD") & V_MASK) | V3_BITMASK;

  bytes32 internal constant DEFAULT_FEED_2_V3 = (keccak256("LINK-USD") & V_MASK) | V3_BITMASK;
  bytes32 internal constant DEFAULT_CONFIG_DIGEST = keccak256("DEFAULT_CONFIG_DIGEST");

  //report
  uint256 internal constant DEFAULT_REPORT_LINK_FEE = 1e10;
  uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e12;

  //rewards
  uint256 internal constant FEE_SCALAR = 1e18;

  address internal constant NATIVE_WITHDRAW_ADDRESS = address(0);

  //the selector for each error
  bytes4 internal immutable INVALID_DISCOUNT_ERROR = FeeManager.InvalidDiscount.selector;
  bytes4 internal immutable INVALID_ADDRESS_ERROR = FeeManager.InvalidAddress.selector;
  bytes4 internal immutable INVALID_SURCHARGE_ERROR = FeeManager.InvalidSurcharge.selector;
  bytes4 internal immutable EXPIRED_REPORT_ERROR = FeeManager.ExpiredReport.selector;
  bytes4 internal immutable INVALID_DEPOSIT_ERROR = FeeManager.InvalidDeposit.selector;
  bytes4 internal immutable INVALID_QUOTE_ERROR = FeeManager.InvalidQuote.selector;
  bytes4 internal immutable UNAUTHORIZED_ERROR = FeeManager.Unauthorized.selector;
  bytes4 internal immutable INVALID_REPORT_VERSION_ERROR = FeeManager.InvalidReportVersion.selector;
  bytes internal constant ONLY_CALLABLE_BY_OWNER_ERROR = "Only callable by owner";
  bytes internal constant INSUFFICIENT_ALLOWANCE_ERROR = "ERC20: insufficient allowance";

  //events emitted
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint256 discount);
  event NativeSurchargeUpdated(uint256 newSurcharge);
  event InsufficientLink(bytes32 indexed configDigest, uint256 linkQuantity, uint256 nativeQuantity);
  event Withdraw(address adminAddress, address assetAddress, uint256 quantity);

  function setUp() public virtual {
    //change to admin user
    vm.startPrank(ADMIN);

    //init required contracts
    _initializeContracts();
  }

  function _initializeContracts() internal {
    link = new ERC20Mock("LINK", "LINK", ADMIN, 0);
    native = new WERC20Mock();

    rewardManager = new RewardManager(getLinkAddress());
    feeManager = new FeeManager(getLinkAddress(), getNativeAddress(), PROXY, address(rewardManager));

    //link the feeManager to the reward manager
    rewardManager.setFeeManager(address(feeManager));

    //mint some tokens to the admin
    link.mint(ADMIN, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(ADMIN, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(ADMIN, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some tokens to the proxy
    link.mint(PROXY, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(PROXY, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(PROXY, DEFAULT_NATIVE_MINT_QUANTITY);
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

  function setNativeSurcharge(uint256 surcharge, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the surcharge
    feeManager.setNativeSurcharge(surcharge);

    //change back to the original address
    changePrank(originalAddr);
  }

  // solium-disable-next-line no-unused-vars
  function getFee(
    bytes memory report,
    IFeeManager.Quote memory quote,
    address subscriber
  ) public view returns (Common.Asset memory) {
    //get the fee
    (Common.Asset memory fee, ) = feeManager.getFeeAndReward(subscriber, report, quote);

    return fee;
  }

  function getReward(
    bytes memory report,
    IFeeManager.Quote memory quote,
    address subscriber
  ) public view returns (Common.Asset memory) {
    //get the reward
    (, Common.Asset memory reward) = feeManager.getFeeAndReward(subscriber, report, quote);

    return reward;
  }

  function getV0Report(bytes32 feedId) public pure returns (bytes memory) {
    return abi.encode(feedId, uint32(0), int192(0), int192(0), int192(0), uint64(0), bytes32(0), uint64(0), uint64(0));
  }

  function getV1Report(bytes32 feedId) public view returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        int192(0),
        uint32(0),
        uint32(block.timestamp),
        DEFAULT_REPORT_LINK_FEE,
        DEFAULT_REPORT_NATIVE_FEE
      );
  }

  function getV1ReportWithExpiryAndFee(
    bytes32 feedId,
    uint256 expiry,
    uint256 linkFee,
    uint256 nativeFee
  ) public view returns (bytes memory) {
    return abi.encode(feedId, uint32(0), int192(0), uint32(0), uint32(expiry), linkFee, nativeFee);
  }

  function getV2Report(bytes32 feedId) public view returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        int192(0),
        int192(0),
        int192(0),
        uint32(0),
        uint32(block.timestamp),
        uint192(DEFAULT_REPORT_LINK_FEE),
        uint192(DEFAULT_REPORT_NATIVE_FEE)
      );
  }

  function getV2ReportWithCustomExpiryAndFee(
    bytes32 feedId,
    uint256 expiry,
    uint256 linkFee,
    uint256 nativeFee
  ) public pure returns (bytes memory) {
    return
      abi.encode(
        feedId,
        uint32(0),
        int192(0),
        int192(0),
        int192(0),
        uint32(0),
        uint32(expiry),
        uint192(linkFee),
        uint192(nativeFee)
      );
  }

  function getLinkQuote() public view returns (IFeeManager.Quote memory) {
    return IFeeManager.Quote(getLinkAddress());
  }

  function getNativeQuote() public view returns (IFeeManager.Quote memory) {
    return IFeeManager.Quote(getNativeAddress());
  }

  function withdraw(address assetAddress, uint256 amount, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //set the surcharge
    feeManager.withdraw(assetAddress, amount);

    //change back to the original address
    changePrank(originalAddr);
  }

  function getLinkBalance(address balanceAddress) public view returns (uint256) {
    return link.balanceOf(balanceAddress);
  }

  function getNativeBalance(address balanceAddress) public view returns (uint256) {
    return native.balanceOf(balanceAddress);
  }

  function getNativeUnwrappedBalance(address balanceAddress) public view returns (uint256) {
    return balanceAddress.balance;
  }

  function mintLink(address recipient, uint256 amount) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(ADMIN);

    //mint the link to the recipient
    link.mint(recipient, amount);

    //change back to the original address
    changePrank(originalAddr);
  }

  function mintNative(address recipient, uint256 amount, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //mint the native to the recipient
    native.mint(recipient, amount);

    //change back to the original address
    changePrank(originalAddr);
  }

  function issueUnwrappedNative(address recipient, uint256 quantity) public {
    vm.deal(recipient, quantity);
  }

  function processFee(bytes memory payload, address subscriber, uint256 wrappedNativeValue, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //process the fee
    feeManager.processFee{value: wrappedNativeValue}(payload, subscriber);

    //change back to the original address
    changePrank(originalAddr);
  }

  function getPayload(bytes memory reportPayload, bytes memory quotePayload) public pure returns (bytes memory) {
    return
      abi.encode(
        [DEFAULT_CONFIG_DIGEST, 0, 0],
        reportPayload,
        new bytes32[](1),
        new bytes32[](1),
        bytes32(""),
        quotePayload
      );
  }

  function getQuotePayload(address quoteAddress) public pure returns (bytes memory) {
    return abi.encode(quoteAddress);
  }

  function approveLink(address spender, uint256 quantity, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //approve the link to be transferred
    link.approve(spender, quantity);

    //change back to the original address
    changePrank(originalAddr);
  }

  function approveNative(address spender, uint256 quantity, address sender) public {
    //record the current address and switch to the recipient
    address originalAddr = msg.sender;
    changePrank(sender);

    //approve the link to be transferred
    native.approve(spender, quantity);

    //change back to the original address
    changePrank(originalAddr);
  }

  function getLinkAddress() public view returns (address) {
    return address(link);
  }

  function getNativeAddress() public view returns (address) {
    return address(native);
  }
}
