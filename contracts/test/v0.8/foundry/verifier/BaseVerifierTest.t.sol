// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";
import {IERC165} from "@openzeppelin/contracts/interfaces/IERC165.sol";
import {IVerifier} from "../../../../src/v0.8/interfaces/IVerifier.sol";
import {ErroredVerifier} from "./mocks/ErroredVerifier.sol";
import {Verifier} from "../../../../src/v0.8/Verifier.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {AccessControllerInterface} from "../../../../src/v0.8/interfaces/AccessControllerInterface.sol";

contract BaseTest is Test {
  uint256 internal constant MAX_ORACLES = 31;
  address internal constant ADMIN = address(1);
  address internal constant USER = address(2);
  address internal constant MOCK_VERIFIER_ADDRESS = address(100);
  address internal constant MOCK_VERIFIER_ADDRESS_TWO = address(200);
  address internal constant ACCESS_CONTROLLER_ADDRESS = address(300);

  bytes32 internal constant FEED_ID = keccak256("ETH-USD");
  bytes32 internal constant FEED_ID_2 = keccak256("LINK-USD");
  bytes32 internal constant FEED_ID_3 = keccak256("BTC-USD");
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

  struct Report {
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
  }

  Signer[MAX_ORACLES] internal s_signers;
  bytes32[] internal s_offchaintransmitters;

  function setUp() public virtual {
    changePrank(ADMIN);
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
    Report memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  )
    internal
    pure
    returns (
      bytes32[] memory rawRs,
      bytes32[] memory rawSs,
      bytes32 rawVs
    )
  {
    bytes32[] memory rs = new bytes32[](signers.length);
    bytes32[] memory ss = new bytes32[](signers.length);
    bytes memory vs = new bytes(signers.length);

    bytes32 hash = keccak256(abi.encodePacked(keccak256(abi.encode(report)), reportContext));

    for (uint256 i = 0; i < signers.length; i++) {
      (uint8 v, bytes32 r, bytes32 s) = vm.sign(signers[i].mockPrivateKey, hash);
      rs[i] = r;
      ss[i] = s;
      vs[i] = bytes1(v - 27);
    }
    return (rs, ss, bytes32(vs));
  }

  function _generateEncodedBlob(
    Report memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  ) internal pure returns (bytes memory) {
    (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _generateSignerSignatures(
      report,
      reportContext,
      signers
    );
    return abi.encode(reportContext, abi.encode(report), rs, ss, rawVs);
  }

  function _configDigestFromConfigData(
    bytes32 feedId,
    address verifierAddr,
    uint64 configCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal view returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          feedId,
          block.chainid,
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
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  function _createReport(
    bytes32 feedId,
    uint32 observationsTimestamp,
    int192 median,
    int192 bid,
    int192 ask,
    uint64 blocknumberUpperBound,
    bytes32 upperBlockhash,
    uint64 blocknumberLowerBound
  ) internal pure returns (Report memory) {
    return
      Report({
        feedId: feedId,
        observationsTimestamp: observationsTimestamp,
        median: median,
        bid: bid,
        ask: ask,
        blocknumberUpperBound: blocknumberUpperBound,
        upperBlockhash: upperBlockhash,
        blocknumberLowerBound: blocknumberLowerBound
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

contract BaseTestWithConfiguredVerifier is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    // Verifier 1, Feed 1, Config 1
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
    (, , bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    s_verifierProxy.initializeVerifier(address(s_verifier));
    changePrank(address(s_verifier));
    s_verifierProxy.setVerifier(bytes32(""), configDigest);
    changePrank(ADMIN);
  }
}

contract BaseTestWithMultipleConfiguredDigests is BaseTestWithConfiguredVerifier {
  bytes32 internal s_configDigestOne;
  bytes32 internal s_configDigestTwo;
  bytes32 internal s_configDigestThree;
  bytes32 internal s_configDigestFour;
  bytes32 internal s_configDigestFive;

  uint32 internal s_numConfigsSet;

  uint8 internal constant FAULT_TOLERANCE_TWO = 2;
  uint8 internal constant FAULT_TOLERANCE_THREE = 1;

  function setUp() public virtual override {
    BaseTestWithConfiguredVerifier.setUp();
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
      bytes("")
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
      bytes("")
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
      bytes("")
    );
    (, , s_configDigestFour) = s_verifier.latestConfigDetails(FEED_ID_2);
    changePrank(address(s_verifier));
    s_verifierProxy.setVerifier(bytes32(""), s_configDigestFour);
    changePrank(ADMIN);

    // Verifier 2, Feed 3, Config 1
    s_verifier_2.setConfig(
      FEED_ID_3,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
    (, , s_configDigestFive) = s_verifier_2.latestConfigDetails(FEED_ID_3);
    s_verifierProxy.initializeVerifier(address(s_verifier_2));
    changePrank(address(s_verifier_2));
    s_verifierProxy.setVerifier(bytes32(""), s_configDigestFive);
    changePrank(ADMIN);
  }
}
