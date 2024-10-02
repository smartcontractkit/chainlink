// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {DataFeedsBase} from "../dev/DataFeedsBase.sol";
import {BaseTest} from "../../shared/test/BaseTest.t.sol";
import {ERC20Mock} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";

contract DataFeedsTestSetup is BaseTest {
  uint256 internal constant PERCENTAGE_SCALAR = 1e18;

  bytes32 internal constant V_MASK = 0x0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
  bytes32 private constant BLOCK_PREMIUM_SCHEMA = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant BASIC_SCHEMA = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant PREMIUM_SCHEMA = 0x0003000000000000000000000000000000000000000000000000000000000000;

  string constant BLOCK_PREMIUM_FEED_ID_STR = "0x0001e44ae0c3807eaf5db446f48222cb7b23f1a944abc299b1b24fccfe4ef5f7";
  string constant BASIC_FEED_ID_STR = "0x0002144b359ba1147d6134b23f204c0ffb561826c5e3bfe7d7e1c08e0210d108";
  string constant PREMIUM_FEED_ID_STR = "0x0003fc78f4eb041431937a8ab0a648c5d284bf8c382891324a25cae47d19e3ba";

  ERC20Mock internal link;

  uint256 internal constant ORACLE_NUM = 4;

  struct Signer {
    uint256 mockPrivateKey;
    address signerAddress;
  }

  DataFeedsBase.BlockPremiumSchema internal reportStructBlockPremium =
    DataFeedsBase.BlockPremiumSchema({
      feedId: (keccak256("FEED_ID_BLOCK_PREMIUM") & V_MASK) | BLOCK_PREMIUM_SCHEMA,
      observationsTimestamp: uint32(BLOCK_TIME),
      benchmarkPrice: 23456,
      bid: 23457,
      ask: 23455,
      currentBlockNum: 10,
      currentBlockHash: blockhash(10),
      validFromBlockNum: 9,
      currentBlockTimestamp: uint32(BLOCK_TIME + 5)
    });

  DataFeedsBase.BasicSchema internal reportStructBasic =
    DataFeedsBase.BasicSchema({
      feedId: (keccak256("FEED_ID_BASIC") & V_MASK) | BASIC_SCHEMA,
      validFromTimestamp: uint32(BLOCK_TIME),
      observationsTimestamp: uint32(BLOCK_TIME + 5),
      nativeFee: 10000000,
      linkFee: 20000000,
      expiresAt: uint32(BLOCK_TIME + 10),
      benchmarkPrice: 34567
    });

  DataFeedsBase.PremiumSchema internal reportStructPremium =
    DataFeedsBase.PremiumSchema({
      feedId: (keccak256("FEED_ID_PREMIUM") & V_MASK) | PREMIUM_SCHEMA,
      validFromTimestamp: uint32(BLOCK_TIME + 15),
      observationsTimestamp: uint32(BLOCK_TIME + 20),
      nativeFee: 30000000,
      linkFee: 40000000,
      expiresAt: uint32(BLOCK_TIME + 25),
      benchmarkPrice: 45678,
      bid: 45679,
      ask: 45677
    });

  bytes reportBlockPremium;
  bytes reportBasic;
  bytes reportPremium;

  function setUp() public virtual override {
    BaseTest.setUp();

    link = new ERC20Mock("LINK", "LINK", OWNER, 0);

    Signer[] memory signers = new Signer[](ORACLE_NUM);
    for (uint256 i = 0; i < ORACLE_NUM; i++) {
      uint256 mockPK = i + 1000;
      signers[i].mockPrivateKey = mockPK;
      signers[i].signerAddress = vm.addr(mockPK);
    }

    bytes32[3] memory reportContextBlockPremium;
    reportContextBlockPremium[0] = keccak256("CONFIG_DIGEST_BLOCK_PREMIUM");
    reportContextBlockPremium[1] = bytes32(abi.encode(uint32(1), uint8(2)));

    reportBlockPremium = _generateBlockPremiumReport(reportStructBlockPremium, reportContextBlockPremium, signers);

    bytes32[3] memory reportContextBasic;
    reportContextBasic[0] = keccak256("CONFIG_DIGEST_BASIC");
    reportContextBasic[1] = bytes32(abi.encode(uint32(3), uint8(4)));

    reportBasic = _generateBasicReport(reportStructBasic, reportContextBasic, signers);

    bytes32[3] memory reportContextPremium;
    reportContextPremium[0] = keccak256("CONFIG_DIGEST_PREMIUM");
    reportContextPremium[1] = bytes32(abi.encode(uint32(5), uint8(6)));

    reportPremium = _generatePremiumReport(reportStructPremium, reportContextPremium, signers);
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

  function _generateBlockPremiumReport(
    DataFeedsBase.BlockPremiumSchema memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  ) internal pure returns (bytes memory) {
    bytes memory reportData = abi.encode(
      report.feedId,
      report.observationsTimestamp,
      report.benchmarkPrice,
      report.bid,
      report.ask,
      report.currentBlockNum,
      report.currentBlockHash,
      report.validFromBlockNum,
      report.currentBlockTimestamp
    );

    (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _generateSignerSignatures(
      reportData,
      reportContext,
      signers
    );
    return abi.encode(reportContext, reportData, rs, ss, rawVs);
  }

  function _generateBasicReport(
    DataFeedsBase.BasicSchema memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  ) internal pure returns (bytes memory) {
    bytes memory reportData = abi.encode(
      report.feedId,
      report.validFromTimestamp,
      report.observationsTimestamp,
      report.nativeFee,
      report.linkFee,
      report.expiresAt,
      report.benchmarkPrice
    );

    (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _generateSignerSignatures(
      reportData,
      reportContext,
      signers
    );
    return abi.encode(reportContext, reportData, rs, ss, rawVs);
  }

  function _generatePremiumReport(
    DataFeedsBase.PremiumSchema memory report,
    bytes32[3] memory reportContext,
    Signer[] memory signers
  ) internal pure returns (bytes memory) {
    bytes memory reportData = abi.encode(
      report.feedId,
      report.validFromTimestamp,
      report.observationsTimestamp,
      report.nativeFee,
      report.linkFee,
      report.expiresAt,
      report.benchmarkPrice,
      report.bid,
      report.ask
    );

    (bytes32[] memory rs, bytes32[] memory ss, bytes32 rawVs) = _generateSignerSignatures(
      reportData,
      reportContext,
      signers
    );
    return abi.encode(reportContext, reportData, rs, ss, rawVs);
  }
}
