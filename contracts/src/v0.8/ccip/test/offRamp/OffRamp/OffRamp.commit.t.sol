// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IFeeQuoter} from "../../../interfaces/IFeeQuoter.sol";
import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";

import {FeeQuoter} from "../../../FeeQuoter.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_commit is OffRampSetup {
  uint64 internal s_maxInterval = 12;

  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();

    s_latestSequenceNumber = uint64(uint256(s_configDigestCommit));
  }

  function test_ReportAndPriceUpdate() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_maxInterval + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_ReportOnlyRootSuccess_gas() public {
    uint64 max1 = 931;
    bytes32 root = "Only a single root";

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: max1,
      merkleRoot: root
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(max1 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(block.timestamp, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR_1, root));
  }

  function test_RootWithRMNDisabled() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    // but ☝️ doesn't matter because RMN verification is disabled
    OffRamp.SourceChainConfig memory sourceChainConfig = s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1);

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigUpdates = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigUpdates[0] = OffRamp.SourceChainConfigArgs({
      router: sourceChainConfig.router,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      isEnabled: sourceChainConfig.isEnabled,
      isRMNVerificationDisabled: true,
      onRamp: sourceChainConfig.onRamp
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigUpdates);

    uint64 max1 = 931;
    bytes32 root = "Only a single root";

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: max1,
      merkleRoot: root
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: new Internal.MerkleRoot[](0),
      unblessedMerkleRoots: roots,
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(max1 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(block.timestamp, s_offRamp.getMerkleRoot(SOURCE_CHAIN_SELECTOR_1, root));
  }

  function test_StaleReportWithRoot() public {
    uint64 maxSeq = 12;
    uint224 tokenStartPrice = IFeeQuoter(s_offRamp.getDynamicConfig().feeQuoter).getTokenPrice(s_sourceFeeToken).value;

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: maxSeq,
      merkleRoot: "stale report 1"
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    commitReport.blessedMerkleRoots[0].minSeqNr = maxSeq + 1;
    commitReport.blessedMerkleRoots[0].maxSeqNr = maxSeq * 2;
    commitReport.blessedMerkleRoots[0].merkleRoot = "stale report 2";

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq * 2 + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());
    assertEq(tokenStartPrice, IFeeQuoter(s_offRamp.getDynamicConfig().feeQuoter).getTokenPrice(s_sourceFeeToken).value);
  }

  function test_OnlyTokenPriceUpdates() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_OnlyGasPriceUpdates() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  function test_PriceSequenceNumberCleared() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    _commit(commitReport, s_latestSequenceNumber);

    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    vm.startPrank(OWNER);
    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    // Execution plugin OCR config should not clear latest epoch and round
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    // Commit plugin config should clear latest epoch & round
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: F,
      isSignatureVerificationEnabled: true,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    assertEq(0, s_offRamp.getLatestPriceSequenceNumber());

    // The same sequence number can be reported again
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_ValidPriceUpdateThenStaleReportWithRoot() public {
    uint64 maxSeq = 12;
    uint224 tokenPrice1 = 4e18;
    uint224 tokenPrice2 = 5e18;
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice1),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, tokenPrice1, block.timestamp);

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());

    roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: maxSeq,
      merkleRoot: "stale report"
    });
    commitReport.priceUpdates = _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice2);
    commitReport.blessedMerkleRoots = roots;

    vm.expectEmit();
    emit OffRamp.CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    vm.expectEmit();
    emit MultiOCR3Base.Transmitted(uint8(Internal.OCRPluginType.Commit), s_configDigestCommit, s_latestSequenceNumber);

    _commit(commitReport, s_latestSequenceNumber);

    assertEq(maxSeq + 1, s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(tokenPrice1, IFeeQuoter(s_offRamp.getDynamicConfig().feeQuoter).getTokenPrice(s_sourceFeeToken).value);
    assertEq(s_latestSequenceNumber, s_offRamp.getLatestPriceSequenceNumber());
  }

  // Reverts

  function test_RevertWhen_UnauthorizedTransmitter() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[2] memory reportContext = [s_configDigestCommit, bytes32(uint256(s_latestSequenceNumber))];

    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, F + 1);

    vm.expectRevert(MultiOCR3Base.UnauthorizedTransmitter.selector);
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_RevertWhen_NoConfig() public {
    _redeployOffRampWithNoOCRConfigs();

    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[2] memory reportContext = [bytes32(""), s_configDigestCommit];
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, F + 1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_RevertWhen_NoConfigWithOtherConfigPresent() public {
    _redeployOffRampWithNoOCRConfigs();

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Execution),
      configDigest: s_configDigestExec,
      F: F,
      isSignatureVerificationEnabled: false,
      signers: s_emptySigners,
      transmitters: s_validTransmitters
    });
    s_offRamp.setOCR3Configs(ocrConfigs);

    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    bytes32[2] memory reportContext = [bytes32(""), s_configDigestCommit];
    (bytes32[] memory rs, bytes32[] memory ss,, bytes32 rawVs) =
      _getSignaturesForDigest(s_validSignerKeys, abi.encode(commitReport), reportContext, F + 1);

    vm.startPrank(s_validTransmitters[0]);
    vm.expectRevert();
    s_offRamp.commit(reportContext, abi.encode(commitReport), rs, ss, rawVs);
  }

  function test_RevertWhen_FailedRMNVerifications() public {
    // force RMN verification to fail
    vm.mockCallRevert(address(s_mockRMNRemote), abi.encodeWithSelector(IRMNRemote.verify.selector), bytes(""));

    OffRamp.CommitReport memory commitReport = _constructCommitReport();
    vm.expectRevert();
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_Unhealthy() public {
    _setMockRMNChainCurse(SOURCE_CHAIN_SELECTOR_1, true);
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "Only a single root",
      onRampAddress: abi.encode(ON_RAMP_ADDRESS_1)
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectRevert(abi.encodeWithSelector(OffRamp.CursedByRMN.selector, roots[0].sourceChainSelector));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_InvalidRoot() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 4,
      merkleRoot: bytes32(0)
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectRevert(OffRamp.InvalidRoot.selector);
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_InvalidInterval() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 2,
      maxSeqNr: 2,
      merkleRoot: bytes32(0)
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.InvalidInterval.selector, roots[0].sourceChainSelector, roots[0].minSeqNr, roots[0].maxSeqNr
      )
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_InvalidIntervalMinLargerThanMax() public {
    s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR);
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 0,
      merkleRoot: bytes32(0)
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.InvalidInterval.selector, roots[0].sourceChainSelector, roots[0].minSeqNr, roots[0].maxSeqNr
      )
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_ZeroEpochAndRound() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectRevert(OffRamp.StaleCommitReport.selector);
    _commit(commitReport, 0);
  }

  function test_RevertWhen_OnlyPriceUpdateStaleReport() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](0);
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    _commit(commitReport, s_latestSequenceNumber);

    vm.expectRevert(OffRamp.StaleCommitReport.selector);
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_SourceChainNotEnabled() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: 0,
      onRampAddress: abi.encode(ON_RAMP_ADDRESS_1),
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "Only a single root"
    });

    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    vm.expectRevert(abi.encodeWithSelector(OffRamp.SourceChainNotEnabled.selector, 0));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_RootAlreadyCommitted() public {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "Only a single root"
    });
    OffRamp.CommitReport memory commitReport = OffRamp.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });

    _commit(commitReport, s_latestSequenceNumber);
    commitReport.blessedMerkleRoots[0].minSeqNr = 3;
    commitReport.blessedMerkleRoots[0].maxSeqNr = 3;

    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.RootAlreadyCommitted.selector, roots[0].sourceChainSelector, roots[0].merkleRoot)
    );
    _commit(commitReport, ++s_latestSequenceNumber);
  }

  function test_RevertWhen_CommitOnRampMismatch() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    commitReport.blessedMerkleRoots[0].onRampAddress = ON_RAMP_ADDRESS_2;

    vm.expectRevert(abi.encodeWithSelector(OffRamp.CommitOnRampMismatch.selector, ON_RAMP_ADDRESS_2, ON_RAMP_ADDRESS_1));
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_RootBlessingMismatch_blessedButShouldNot() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();
    uint64 sourceChainSelector = commitReport.blessedMerkleRoots[0].sourceChainSelector;

    OffRamp.SourceChainConfig memory sourceChainConfig = s_offRamp.getSourceChainConfig(sourceChainSelector);

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigUpdates = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigUpdates[0] = OffRamp.SourceChainConfigArgs({
      router: sourceChainConfig.router,
      sourceChainSelector: sourceChainSelector,
      isEnabled: sourceChainConfig.isEnabled,
      isRMNVerificationDisabled: true,
      onRamp: sourceChainConfig.onRamp
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.RootBlessingMismatch.selector,
        commitReport.blessedMerkleRoots[0].sourceChainSelector,
        commitReport.blessedMerkleRoots[0].merkleRoot,
        true
      )
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function test_RevertWhen_RootBlessingMismatch_unblessedButShouldBeBlessed() public {
    OffRamp.CommitReport memory commitReport = _constructCommitReport();

    commitReport.unblessedMerkleRoots = commitReport.blessedMerkleRoots;
    commitReport.blessedMerkleRoots = new Internal.MerkleRoot[](0);

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.RootBlessingMismatch.selector,
        commitReport.unblessedMerkleRoots[0].sourceChainSelector,
        commitReport.unblessedMerkleRoots[0].merkleRoot,
        false
      )
    );
    _commit(commitReport, s_latestSequenceNumber);
  }

  function _constructCommitReport() internal view returns (OffRamp.CommitReport memory) {
    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: s_maxInterval,
      merkleRoot: "test #2"
    });

    return OffRamp.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      blessedMerkleRoots: roots,
      unblessedMerkleRoots: new Internal.MerkleRoot[](0),
      rmnSignatures: s_rmnSignatures
    });
  }
}
