// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {MultiOCR3Helper} from "../../helpers/MultiOCR3Helper.sol";

contract MultiOCR3BaseSetup is BaseTest {
  // Signer private keys used for these test
  uint256 internal constant PRIVATE0 = 0x7b2e97fe057e6de99d6872a2ef2abf52c9b4469bc848c2465ac3fcd8d336e81d;

  address[] internal s_validSigners;
  address[] internal s_validTransmitters;
  uint256[] internal s_validSignerKeys;

  address[] internal s_partialSigners;
  address[] internal s_partialTransmitters;
  uint256[] internal s_partialSignerKeys;

  address[] internal s_emptySigners;

  bytes internal constant REPORT = abi.encode("testReport");
  MultiOCR3Helper internal s_multiOCR3;

  function setUp() public virtual override {
    BaseTest.setUp();

    uint160 numSigners = 7;
    s_validSignerKeys = new uint256[](numSigners);
    s_validSigners = new address[](numSigners);
    s_validTransmitters = new address[](numSigners);

    for (uint160 i; i < numSigners; ++i) {
      s_validTransmitters[i] = address(4 + i);
      s_validSignerKeys[i] = PRIVATE0 + i;
      s_validSigners[i] = vm.addr(s_validSignerKeys[i]);
    }

    s_partialSigners = new address[](4);
    s_partialSignerKeys = new uint256[](4);
    s_partialTransmitters = new address[](4);
    for (uint256 i; i < s_partialSigners.length; ++i) {
      s_partialSigners[i] = s_validSigners[i];
      s_partialSignerKeys[i] = s_validSignerKeys[i];
      s_partialTransmitters[i] = s_validTransmitters[i];
    }

    s_emptySigners = new address[](0);

    s_multiOCR3 = new MultiOCR3Helper();
  }

  /// @dev returns a mock config digest with config digest computation logic similar to OCR2Base
  function _getBasicConfigDigest(
    uint8 F,
    address[] memory signers,
    address[] memory transmitters
  ) internal view returns (bytes32) {
    bytes memory configBytes = abi.encode("");
    uint256 configVersion = 1;

    uint256 h = uint256(
      keccak256(
        abi.encode(
          block.chainid, address(s_multiOCR3), signers, transmitters, F, configBytes, configVersion, configBytes
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  function _assertOCRConfigEquality(
    MultiOCR3Base.OCRConfig memory configA,
    MultiOCR3Base.OCRConfig memory configB
  ) internal pure {
    vm.assertEq(configA.configInfo.configDigest, configB.configInfo.configDigest);
    vm.assertEq(configA.configInfo.F, configB.configInfo.F);
    vm.assertEq(configA.configInfo.n, configB.configInfo.n);
    vm.assertEq(configA.configInfo.isSignatureVerificationEnabled, configB.configInfo.isSignatureVerificationEnabled);

    vm.assertEq(configA.signers, configB.signers);
    vm.assertEq(configA.transmitters, configB.transmitters);
  }

  function _assertOCRConfigUnconfigured(
    MultiOCR3Base.OCRConfig memory config
  ) internal pure {
    assertEq(config.configInfo.configDigest, bytes32(""));
    assertEq(config.signers.length, 0);
    assertEq(config.transmitters.length, 0);
  }

  function _getSignaturesForDigest(
    uint256[] memory signerPrivateKeys,
    bytes memory report,
    bytes32[2] memory reportContext,
    uint8 signatureCount
  ) internal pure returns (bytes32[] memory rs, bytes32[] memory ss, uint8[] memory vs, bytes32 rawVs) {
    rs = new bytes32[](signatureCount);
    ss = new bytes32[](signatureCount);
    vs = new uint8[](signatureCount);

    bytes32 reportDigest = keccak256(abi.encodePacked(keccak256(report), reportContext));

    // Calculate signatures
    for (uint256 i; i < signatureCount; ++i) {
      (vs[i], rs[i], ss[i]) = vm.sign(signerPrivateKeys[i], reportDigest);
      rawVs = rawVs | (bytes32(bytes1(vs[i] - 27)) >> (8 * i));
    }

    return (rs, ss, vs, rawVs);
  }
}
