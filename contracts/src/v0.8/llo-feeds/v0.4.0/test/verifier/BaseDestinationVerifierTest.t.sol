// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {DestinationVerifierProxy} from "../../DestinationVerifierProxy.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IDestinationVerifier} from "../../interfaces/IDestinationVerifier.sol";
import {IDestinationVerifierProxy} from "../../interfaces/IDestinationVerifierProxy.sol";
import {DestinationVerifier} from "../../DestinationVerifier.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {DestinationFeeManager} from "../../DestinationFeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import {ERC20Mock} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {WERC20Mock} from "../../../../shared/mocks/WERC20Mock.sol";
import {DestinationRewardManager} from "../../DestinationRewardManager.sol";
import {IDestinationRewardManager} from "../../interfaces/IDestinationRewardManager.sol";

contract BaseTest is Test {
  uint64 internal constant POOL_SCALAR = 1e18;
  uint64 internal constant ONE_PERCENT = POOL_SCALAR / 100;
  uint256 internal constant MAX_ORACLES = 31;
  address internal constant ADMIN = address(1);
  address internal constant USER = address(2);

  address internal constant MOCK_VERIFIER_ADDRESS = address(100);
  address internal constant ACCESS_CONTROLLER_ADDRESS = address(300);

  uint256 internal constant DEFAULT_REPORT_LINK_FEE = 1e10;
  uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e12;

  uint64 internal constant VERIFIER_VERSION = 1;

  uint8 internal constant FAULT_TOLERANCE = 10;

  DestinationVerifierProxy internal s_verifierProxy;
  DestinationVerifier internal s_verifier;
  DestinationFeeManager internal feeManager;
  DestinationRewardManager internal rewardManager;
  ERC20Mock internal link;
  WERC20Mock internal native;

  struct Signer {
    uint256 mockPrivateKey;
    address signerAddress;
  }

  Signer[MAX_ORACLES] internal s_signers;
  bytes32[] internal s_offchaintransmitters;
  bool private s_baseTestInitialized;

  struct V3Report {
    // The feed ID the report has data for
    bytes32 feedId;
    // The time the median value was observed on
    uint32 observationsTimestamp;
    // The timestamp the report is valid from
    uint32 validFromTimestamp;
    // The link fee
    uint192 linkFee;
    // The native fee
    uint192 nativeFee;
    // The expiry of the report
    uint32 expiresAt;
    // The median value agreed in an OCR round
    int192 benchmarkPrice;
    // The best bid value agreed in an OCR round
    int192 bid;
    // The best ask value agreed in an OCR round
    int192 ask;
  }

  bytes32 internal constant V_MASK = 0x0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
  bytes32 internal constant V1_BITMASK = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V2_BITMASK = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V3_BITMASK = 0x0003000000000000000000000000000000000000000000000000000000000000;

  bytes32 internal constant INVALID_FEED = keccak256("INVALID");
  uint32 internal constant OBSERVATIONS_TIMESTAMP = 1000;
  uint64 internal constant BLOCKNUMBER_LOWER_BOUND = 1000;
  uint64 internal constant BLOCKNUMBER_UPPER_BOUND = BLOCKNUMBER_LOWER_BOUND + 5;
  int192 internal constant MEDIAN = 1 ether;
  int192 internal constant BID = 500000000 gwei;
  int192 internal constant ASK = 2 ether;

  //version 0 feeds
  bytes32 internal constant FEED_ID = (keccak256("ETH-USD") & V_MASK) | V1_BITMASK;
  bytes32 internal constant FEED_ID_2 = (keccak256("LINK-USD") & V_MASK) | V1_BITMASK;
  bytes32 internal constant FEED_ID_3 = (keccak256("BTC-USD") & V_MASK) | V1_BITMASK;

  //version 3 feeds
  bytes32 internal constant FEED_ID_V3 = (keccak256("ETH-USD") & V_MASK) | V3_BITMASK;

  function _encodeReport(V3Report memory report) internal pure returns (bytes memory) {
    return
      abi.encode(
        report.feedId,
        report.observationsTimestamp,
        report.validFromTimestamp,
        report.nativeFee,
        report.linkFee,
        report.expiresAt,
        report.benchmarkPrice,
        report.bid,
        report.ask
      );
  }

  function _generateSignerSignatures(
    bytes memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  ) internal pure returns (bytes32[] memory rawRs, bytes32[] memory rawSs, bytes32 rawVs) {
    bytes32[] memory rs = new bytes32[](signers.length);
    bytes32[] memory ss = new bytes32[](signers.length);
    bytes memory vs = new bytes(signers.length);

    bytes32 hash = keccak256(abi.encodePacked(keccak256(report), reportContext));

    for (uint256 i = 0; i < signers.length; i++) {
      (uint8 v, bytes32 r, bytes32 s) = vm.sign(signers[i].mockPrivateKey, hash);
      rs[i] = r;
      ss[i] = s;
      vs[i] = bytes1(v - 27);
    }
    return (rs, ss, bytes32(vs));
  }

  function _generateV3EncodedBlob(
    V3Report memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  ) internal pure returns (bytes memory) {
    bytes memory reportBytes = _encodeReport(report);
    (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _generateSignerSignatures(
      reportBytes,
      reportContext,
      signers
    );
    return abi.encode(reportContext, reportBytes, rs, ss, rawVs);
  }

  function _verify(bytes memory payload, address feeAddress, uint256 wrappedNativeValue, address sender) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    s_verifierProxy.verify{value: wrappedNativeValue}(payload, abi.encode(feeAddress));

    changePrank(originalAddr);
  }

  function _generateV3Report() internal view returns (V3Report memory) {
    return
      V3Report({
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
  }

  function _verifyBulk(
    bytes[] memory payload,
    address feeAddress,
    uint256 wrappedNativeValue,
    address sender
  ) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    s_verifierProxy.verifyBulk{value: wrappedNativeValue}(payload, abi.encode(feeAddress));

    changePrank(originalAddr);
  }

  function _approveLink(address spender, uint256 quantity, address sender) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    link.approve(spender, quantity);
    changePrank(originalAddr);
  }

  function setUp() public virtual {
    // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;
    vm.startPrank(ADMIN);

    s_verifierProxy = new DestinationVerifierProxy();
    s_verifier = new DestinationVerifier(address(s_verifierProxy));
    s_verifierProxy.setVerifier(address(s_verifier));

    // setting up FeeManager and RewardManager
    native = new WERC20Mock();
    link = new ERC20Mock("LINK", "LINK", ADMIN, 0);
    rewardManager = new DestinationRewardManager(address(link));
    feeManager = new DestinationFeeManager(address(link), address(native), address(s_verifier), address(rewardManager));

    for (uint256 i; i < MAX_ORACLES; i++) {
      uint256 mockPK = i + 1;
      s_signers[i].mockPrivateKey = mockPK;
      s_signers[i].signerAddress = vm.addr(mockPK);
    }
  }

  function _getSigners(uint256 numSigners) internal view returns (Signer[] memory) {
    Signer[] memory signers = new Signer[](numSigners);
    for (uint256 i; i < numSigners; i++) {
      signers[i] = s_signers[i];
    }
    return signers;
  }

  function _getSignerAddresses(Signer[] memory signers) internal pure returns (address[] memory) {
    address[] memory signerAddrs = new address[](signers.length);
    for (uint256 i = 0; i < signerAddrs.length; i++) {
      signerAddrs[i] = signers[i].signerAddress;
    }
    return signerAddrs;
  }

  function _signerAddressAndDonConfigKey(address signer, bytes24 donConfigId) internal pure returns (bytes32) {
    return keccak256(abi.encodePacked(signer, donConfigId));
  }

  function _donConfigIdFromConfigData(address[] memory signers, uint8 f) internal pure returns (bytes24) {
    Common._quickSort(signers, 0, int256(signers.length - 1));
    bytes24 donConfigId = bytes24(keccak256(abi.encodePacked(signers, f)));
    return donConfigId;
  }

  function assertReportsEqual(bytes memory response, V3Report memory testReport) public pure {
    (
      bytes32 feedId,
      uint32 observationsTimestamp,
      uint32 validFromTimestamp,
      uint192 nativeFee,
      uint192 linkFee,
      uint32 expiresAt,
      int192 benchmarkPrice,
      int192 bid,
      int192 ask
    ) = abi.decode(response, (bytes32, uint32, uint32, uint192, uint192, uint32, int192, int192, int192));
    assertEq(feedId, testReport.feedId);
    assertEq(observationsTimestamp, testReport.observationsTimestamp);
    assertEq(validFromTimestamp, testReport.validFromTimestamp);
    assertEq(expiresAt, testReport.expiresAt);
    assertEq(benchmarkPrice, testReport.benchmarkPrice);
    assertEq(bid, testReport.bid);
    assertEq(ask, testReport.ask);
    assertEq(linkFee, testReport.linkFee);
    assertEq(nativeFee, testReport.nativeFee);
  }

  function _approveNative(address spender, uint256 quantity, address sender) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    native.approve(spender, quantity);
    changePrank(originalAddr);
  }
}

contract VerifierWithFeeManager is BaseTest {
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  function setUp() public virtual override {
    BaseTest.setUp();

    s_verifierProxy.setVerifier(address(s_verifier));
    s_verifier.setFeeManager(address(feeManager));
    rewardManager.addFeeManager(address(feeManager));

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some link tokens to the feeManager pool
    link.mint(address(feeManager), DEFAULT_REPORT_LINK_FEE);
  }
}

contract MultipleVerifierWithMultipleFeeManagers is BaseTest {
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  DestinationVerifier internal s_verifier2;
  DestinationVerifier internal s_verifier3;

  DestinationVerifierProxy internal s_verifierProxy2;
  DestinationVerifierProxy internal s_verifierProxy3;

  DestinationFeeManager internal feeManager2;

  function setUp() public virtual override {
    /*
      - Sets up 3 verifiers
      - Sets up 2 Fee managers, wire the fee managers and verifiers
      - Sets up a Reward Manager which can be used by both fee managers
     */
    BaseTest.setUp();

    s_verifierProxy2 = new DestinationVerifierProxy();
    s_verifierProxy3 = new DestinationVerifierProxy();

    s_verifier2 = new DestinationVerifier(address(s_verifierProxy2));
    s_verifier3 = new DestinationVerifier(address(s_verifierProxy3));

    s_verifierProxy2.setVerifier(address(s_verifier2));
    s_verifierProxy3.setVerifier(address(s_verifier3));

    feeManager2 = new DestinationFeeManager(
      address(link),
      address(native),
      address(s_verifier),
      address(rewardManager)
    );

    s_verifier.setFeeManager(address(feeManager));
    s_verifier2.setFeeManager(address(feeManager));
    s_verifier3.setFeeManager(address(feeManager2));

    // this is already set in the base contract
    // feeManager.addVerifier(address(s_verifier));
    feeManager.addVerifier(address(s_verifier2));
    feeManager2.addVerifier(address(s_verifier3));

    rewardManager.addFeeManager(address(feeManager));
    rewardManager.addFeeManager(address(feeManager2));

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);

    //mint some link tokens to the feeManager pool
    link.mint(address(feeManager), DEFAULT_REPORT_LINK_FEE);
  }
}
