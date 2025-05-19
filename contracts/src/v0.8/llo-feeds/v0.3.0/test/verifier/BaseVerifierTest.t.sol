// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IVerifier} from "../../interfaces/IVerifier.sol";
import {ErroredVerifier} from "../mocks/ErroredVerifier.sol";
import {Verifier} from "../../Verifier.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {FeeManager} from "../../FeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import {ERC20Mock} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {WERC20Mock} from "../../../../shared/mocks/WERC20Mock.sol";
import {RewardManager} from "../../RewardManager.sol";

contract BaseTest is Test {
  uint256 internal constant MAX_ORACLES = 31;
  address internal constant ADMIN = address(1);
  address internal constant USER = address(2);
  address internal constant MOCK_VERIFIER_ADDRESS = address(100);
  address internal constant MOCK_VERIFIER_ADDRESS_TWO = address(200);
  address internal constant ACCESS_CONTROLLER_ADDRESS = address(300);

  bytes32 internal constant V_MASK = 0x0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
  bytes32 internal constant V1_BITMASK = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V2_BITMASK = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 internal constant V3_BITMASK = 0x0003000000000000000000000000000000000000000000000000000000000000;

  //version 0 feeds
  bytes32 internal constant FEED_ID = (keccak256("ETH-USD") & V_MASK) | V1_BITMASK;
  bytes32 internal constant FEED_ID_2 = (keccak256("LINK-USD") & V_MASK) | V1_BITMASK;
  bytes32 internal constant FEED_ID_3 = (keccak256("BTC-USD") & V_MASK) | V1_BITMASK;

  //version 3 feeds
  bytes32 internal constant FEED_ID_V3 = (keccak256("ETH-USD") & V_MASK) | V3_BITMASK;

  bytes32 internal constant INVALID_FEED = keccak256("INVALID");
  uint32 internal constant OBSERVATIONS_TIMESTAMP = 1000;
  uint64 internal constant BLOCKNUMBER_LOWER_BOUND = 1000;
  uint64 internal constant BLOCKNUMBER_UPPER_BOUND = BLOCKNUMBER_LOWER_BOUND + 5;
  int192 internal constant MEDIAN = 1 ether;
  int192 internal constant BID = 500000000 gwei;
  int192 internal constant ASK = 2 ether;

  bytes32 internal constant EMPTY_BYTES = bytes32("");

  uint8 internal constant FAULT_TOLERANCE = 10;
  uint64 internal constant VERIFIER_VERSION = 1;

  string internal constant SERVER_URL = "https://mercury.server/client/";
  uint8 internal constant MAX_COMMITMENT_DELAY = 5;

  VerifierProxy internal s_verifierProxy;
  Verifier internal s_verifier;
  Verifier internal s_verifier_2;
  ErroredVerifier internal s_erroredVerifier;

  struct Signer {
    uint256 mockPrivateKey;
    address signerAddress;
  }

  struct V1Report {
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
  }

  Signer[MAX_ORACLES] internal s_signers;
  bytes32[] internal s_offchaintransmitters;
  bool private s_baseTestInitialized;

  function setUp() public virtual {
    // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;

    vm.startPrank(ADMIN);
    vm.mockCall(
      MOCK_VERIFIER_ADDRESS,
      abi.encodeWithSelector(IERC165.supportsInterface.selector, IVerifier.verify.selector),
      abi.encode(true)
    );
    s_verifierProxy = new VerifierProxy(AccessControllerInterface(address(0)));

    s_verifier = new Verifier(address(s_verifierProxy));
    s_verifier_2 = new Verifier(address(s_verifierProxy));
    s_erroredVerifier = new ErroredVerifier();

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

  function _getSignerAddresses(Signer[] memory signers) internal view returns (address[] memory) {
    address[] memory signerAddrs = new address[](signers.length);
    for (uint256 i = 0; i < signerAddrs.length; i++) {
      signerAddrs[i] = s_signers[i].signerAddress;
    }
    return signerAddrs;
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

  function _encodeReport(V1Report memory report) internal pure returns (bytes memory) {
    return
      abi.encode(
        report.feedId,
        report.observationsTimestamp,
        report.median,
        report.bid,
        report.ask,
        report.blocknumberUpperBound,
        report.upperBlockhash,
        report.blocknumberLowerBound,
        report.currentBlockTimestamp
      );
  }

  function _generateV1EncodedBlob(
    V1Report memory report,
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

  function _configDigestFromConfigData(
    bytes32 feedId,
    uint256 chainId,
    address verifierAddr,
    uint64 configCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          feedId,
          chainId,
          verifierAddr,
          configCount,
          signers,
          offchainTransmitters,
          f,
          onchainConfig,
          offchainConfigVersion,
          offchainConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x0006 << (256 - 16); // 0x000600..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  function _createV1Report(
    bytes32 feedId,
    uint32 observationsTimestamp,
    int192 median,
    int192 bid,
    int192 ask,
    uint64 blocknumberUpperBound,
    bytes32 upperBlockhash,
    uint64 blocknumberLowerBound,
    uint32 currentBlockTimestamp
  ) internal pure returns (V1Report memory) {
    return
      V1Report({
        feedId: feedId,
        observationsTimestamp: observationsTimestamp,
        median: median,
        bid: bid,
        ask: ask,
        blocknumberUpperBound: blocknumberUpperBound,
        upperBlockhash: upperBlockhash,
        blocknumberLowerBound: blocknumberLowerBound,
        currentBlockTimestamp: currentBlockTimestamp
      });
  }

  function _ccipReadURL(bytes32 feedId, uint256 commitmentBlock) internal pure returns (string memory url) {
    return
      string(
        abi.encodePacked(
          SERVER_URL,
          "?feedIDHex=",
          Strings.toHexString(uint256(feedId)),
          "&L2Blocknumber=",
          Strings.toString(commitmentBlock)
        )
      );
  }
}

contract BaseTestWithConfiguredVerifierAndFeeManager is BaseTest {
  FeeManager internal feeManager;
  RewardManager internal rewardManager;
  ERC20Mock internal link;
  WERC20Mock internal native;

  uint256 internal constant DEFAULT_REPORT_LINK_FEE = 1e10;
  uint256 internal constant DEFAULT_REPORT_NATIVE_FEE = 1e12;

  bytes32 internal v1ConfigDigest;
  bytes32 internal v3ConfigDigest;

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

  function setUp() public virtual override {
    BaseTest.setUp();
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    s_verifierProxy.initializeVerifier(address(s_verifier));
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
    (, , v1ConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);

    s_verifier.setConfig(
      FEED_ID_V3,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
    (, , v3ConfigDigest) = s_verifier.latestConfigDetails(FEED_ID_V3);

    link = new ERC20Mock("LINK", "LINK", ADMIN, 0);
    native = new WERC20Mock();

    rewardManager = new RewardManager(address(link));
    feeManager = new FeeManager(address(link), address(native), address(s_verifierProxy), address(rewardManager));

    s_verifierProxy.setFeeManager(feeManager);
    rewardManager.setFeeManager(address(feeManager));
  }

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

  function _generateV1Report() internal view returns (V1Report memory) {
    return
      _createV1Report(
        FEED_ID,
        OBSERVATIONS_TIMESTAMP,
        MEDIAN,
        BID,
        ASK,
        BLOCKNUMBER_UPPER_BOUND,
        bytes32(blockhash(BLOCKNUMBER_UPPER_BOUND)),
        BLOCKNUMBER_LOWER_BOUND,
        uint32(block.timestamp)
      );
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

  function _generateReportContext(bytes32 configDigest) internal pure returns (bytes32[3] memory) {
    bytes32[3] memory reportContext;
    reportContext[0] = configDigest;
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

  function _verify(bytes memory payload, address feeAddress, uint256 wrappedNativeValue, address sender) internal {
    address originalAddr = msg.sender;
    changePrank(sender);

    s_verifierProxy.verify{value: wrappedNativeValue}(payload, abi.encode(feeAddress));

    changePrank(originalAddr);
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
}

contract BaseTestWithMultipleConfiguredDigests is BaseTestWithConfiguredVerifierAndFeeManager {
  bytes32 internal s_configDigestOne;
  bytes32 internal s_configDigestTwo;
  bytes32 internal s_configDigestThree;
  bytes32 internal s_configDigestFour;
  bytes32 internal s_configDigestFive;

  uint32 internal s_numConfigsSet;

  uint8 internal constant FAULT_TOLERANCE_TWO = 2;
  uint8 internal constant FAULT_TOLERANCE_THREE = 1;

  function setUp() public virtual override {
    BaseTestWithConfiguredVerifierAndFeeManager.setUp();
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    (, , s_configDigestOne) = s_verifier.latestConfigDetails(FEED_ID);

    // Verifier 1, Feed 1, Config 2
    Signer[] memory secondSetOfSigners = _getSigners(8);
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(secondSetOfSigners),
      s_offchaintransmitters,
      FAULT_TOLERANCE_TWO,
      bytes(""),
      2,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
    (, , s_configDigestTwo) = s_verifier.latestConfigDetails(FEED_ID);

    // Verifier 1, Feed 1, Config 3
    Signer[] memory thirdSetOfSigners = _getSigners(5);
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(thirdSetOfSigners),
      s_offchaintransmitters,
      FAULT_TOLERANCE_THREE,
      bytes(""),
      3,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
    (s_numConfigsSet, , s_configDigestThree) = s_verifier.latestConfigDetails(FEED_ID);

    // Verifier 1, Feed 2, Config 1
    s_verifier.setConfig(
      FEED_ID_2,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      4,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
    (, , s_configDigestFour) = s_verifier.latestConfigDetails(FEED_ID_2);

    // Verifier 2, Feed 3, Config 1
    s_verifierProxy.initializeVerifier(address(s_verifier_2));
    s_verifier_2.setConfig(
      FEED_ID_3,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
    (, , s_configDigestFive) = s_verifier_2.latestConfigDetails(FEED_ID_3);
  }
}
