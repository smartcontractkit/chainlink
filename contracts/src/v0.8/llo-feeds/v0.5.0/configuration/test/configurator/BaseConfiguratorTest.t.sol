// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {IConfigurator} from "../../interfaces/IConfigurator.sol";
import {Test} from "forge-std/Test.sol";
import {Configurator} from "../../Configurator.sol";
import {ExposedConfigurator} from "../mocks/ExposedConfigurator.sol";
import {ExposedChannelConfigStore} from "../mocks/ExposedChannelConfigStore.sol";

/**
 * @title ConfiguratorTest
 * @author samsondav
 * @notice Base class for Configurator tests
 */
contract BaseTest is Test {
  uint256 internal constant MAX_ORACLES = 31;
  address internal constant USER = address(2);
  bytes32 internal constant CONFIG_ID_1 = (keccak256("CONFIG_ID_1"));
  uint8 internal constant FAULT_TOLERANCE = 10;
  uint64 internal constant OFFCHAIN_CONFIG_VERSION = 1;

  bytes32[] internal s_offchaintransmitters;
  bool private s_baseTestInitialized;

  Configurator internal s_configurator;
  ExposedConfigurator internal s_exposedConfigurator;

  event ProductionConfigSet(
    bytes32 indexed configId,
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    bytes[] signers,
    bytes32[] offchainTransmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig,
    bool isGreenProduction
  );
  event StagingConfigSet(
    bytes32 indexed configId,
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    bytes[] signers,
    bytes32[] offchainTransmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig,
    bool isGreenProduction
  );
  event PromoteStagingConfig(bytes32 indexed configId, bytes32 indexed retiredConfigDigest, bool isGreenProduction);

  bytes[MAX_ORACLES] internal s_signers;

  function setUp() public virtual {
    // BaseTest.setUp may be called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;

    s_configurator = new Configurator();
    s_exposedConfigurator = new ExposedConfigurator();

    for (uint256 i; i < MAX_ORACLES; i++) {
      bytes memory mockSigner = abi.encodePacked(i + 1);
      s_signers[i] = mockSigner;
    }

    for (uint256 i; i < MAX_ORACLES; i++) {
      s_offchaintransmitters.push(bytes32(i + 1));
    }
  }

  function _getSigners(uint256 numSigners) internal view returns (bytes[] memory) {
    bytes[] memory signers = new bytes[](numSigners);
    for (uint256 i; i < numSigners; i++) {
      signers[i] = s_signers[i];
    }
    return signers;
  }

  function _getOffchainTransmitters(uint256 numTransmitters) internal pure returns (bytes32[] memory) {
    bytes32[] memory transmitters = new bytes32[](numTransmitters);
    for (uint256 i; i < numTransmitters; i++) {
      transmitters[i] = bytes32(101 + i);
    }
    return transmitters;
  }
}
