// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {OCR2BaseNoChecks} from "../../ocr/OCR2BaseNoChecks.sol";

contract OCR2NoChecksHelper is OCR2BaseNoChecks {
  function configDigestFromConfigData(
    uint256 chainSelector,
    address contractAddress,
    uint64 configCount,
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) public pure returns (bytes32) {
    return _configDigestFromConfigData(
      chainSelector,
      contractAddress,
      configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function _report(bytes calldata report) internal override {}

  function typeAndVersion() public pure override returns (string memory) {
    return "OCR2BaseHelper 1.0.0";
  }

  function _beforeSetConfig(bytes memory _onchainConfig) internal override {}
}
