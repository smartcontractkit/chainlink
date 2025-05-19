// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

// Derived from v0.5/Configurator.sol
contract MockConfigurator {
  struct ConfigurationState {
    uint64 configCount;
    uint32 latestConfigBlockNumber;
    bytes32 configDigest;
  }

  mapping(bytes32 => ConfigurationState) public s_configurationStates;

  function setStagingConfig(
    bytes32 configId,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external {
    ConfigurationState storage configurationState = s_configurationStates[configId];

    uint64 newConfigCount = ++configurationState.configCount;

    bytes32 configDigest = _configDigestFromConfigData(
      configId,
      block.chainid,
      address(this),
      newConfigCount,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    s_configurationStates[configId].configDigest = configDigest;
    configurationState.latestConfigBlockNumber = uint32(block.number);
  }

  function _configDigestFromConfigData(
    bytes32 configId,
    uint256 sourceChainId,
    address sourceAddress,
    uint64 configCount,
    bytes[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          configId,
          sourceChainId,
          sourceAddress,
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
    uint256 prefixMask = type(uint256).max << (256 - 16);
    uint256 prefix = 0x0009 << (256 - 16);
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }
}
