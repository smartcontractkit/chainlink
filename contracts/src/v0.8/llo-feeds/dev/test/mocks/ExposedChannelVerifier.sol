// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

// ExposedChannelVerifier exposes certain internal Verifier
// methods/structures so that golang code can access them, and we get
// reliable type checking on their usage
contract ExposedChannelVerifier {
  constructor() {}

  function _configDigestFromConfigData(
    uint256 chainId,
    address contractAddress,
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
          chainId,
          contractAddress,
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
    // 0x0009 corresponds to ConfigDigestPrefixLLO in libocr
    uint256 prefix = 0x0009 << (256 - 16); // 0x000900..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  function exposedConfigDigestFromConfigData(
    uint256 _chainId,
    address _contractAddress,
    uint64 _configCount,
    address[] memory _signers,
    bytes32[] memory _offchainTransmitters,
    uint8 _f,
    bytes calldata _onchainConfig,
    uint64 _encodedConfigVersion,
    bytes memory _encodedConfig
  ) public pure returns (bytes32) {
    return
      _configDigestFromConfigData(
        _chainId,
        _contractAddress,
        _configCount,
        _signers,
        _offchainTransmitters,
        _f,
        _onchainConfig,
        _encodedConfigVersion,
        _encodedConfig
      );
  }
}
