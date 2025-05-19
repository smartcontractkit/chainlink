pragma solidity 0.8.6;

import "./BaseTest.t.sol";
import {ChainSpecificUtil} from "../ChainSpecificUtil_v0_8_6.sol";

import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import {ArbGasInfo} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts/v0.8.6/contracts/L2/predeploys/OVM_GasPriceOracle.sol";

contract ChainSpecificUtilTest is BaseTest {
  // ------------ Start Arbitrum Constants ------------

  /// @dev ARBSYS_ADDR is the address of the ArbSys precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbSys.sol#L10
  address private constant ARBSYS_ADDR = address(0x0000000000000000000000000000000000000064);
  ArbSys private constant ARBSYS = ArbSys(ARBSYS_ADDR);

  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  uint256 private constant ARB_MAINNET_CHAIN_ID = 42161;
  uint256 private constant ARB_GOERLI_TESTNET_CHAIN_ID = 421613;
  uint256 private constant ARB_SEPOLIA_TESTNET_CHAIN_ID = 421614;

  // ------------ End Arbitrum Constants ------------

  // ------------ Start Optimism Constants ------------
  /// @dev L1_FEE_DATA_PADDING includes 35 bytes for L1 data padding for Optimism
  bytes internal constant L1_FEE_DATA_PADDING =
    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  uint256 private constant OP_MAINNET_CHAIN_ID = 10;
  uint256 private constant OP_GOERLI_CHAIN_ID = 420;
  uint256 private constant OP_SEPOLIA_CHAIN_ID = 11155420;

  /// @dev Base is a OP stack based rollup and follows the same L1 pricing logic as Optimism.
  uint256 private constant BASE_MAINNET_CHAIN_ID = 8453;
  uint256 private constant BASE_GOERLI_CHAIN_ID = 84531;

  // ------------ End Optimism Constants ------------

  function setUp() public override {
    BaseTest.setUp();
    vm.clearMockedCalls();
  }

  function testGetBlockhashArbitrum() public {
    uint256[3] memory chainIds = [ARB_MAINNET_CHAIN_ID, ARB_GOERLI_TESTNET_CHAIN_ID, ARB_SEPOLIA_TESTNET_CHAIN_ID];
    bytes32[3] memory expectedBlockHashes = [keccak256("mainnet"), keccak256("goerli"), keccak256("sepolia")];
    uint256[3] memory expectedBlockNumbers = [uint256(10), 11, 12];
    for (uint256 i = 0; i < chainIds.length; i++) {
      vm.chainId(chainIds[i]);
      bytes32 expectedBlockHash = expectedBlockHashes[i];
      uint256 expectedBlockNumber = expectedBlockNumbers[i];
      vm.mockCall(
        ARBSYS_ADDR,
        abi.encodeWithSelector(ArbSys.arbBlockNumber.selector),
        abi.encode(expectedBlockNumber + 1)
      );
      vm.mockCall(
        ARBSYS_ADDR,
        abi.encodeWithSelector(ArbSys.arbBlockHash.selector, expectedBlockNumber),
        abi.encodePacked(expectedBlockHash)
      );
      bytes32 actualBlockHash = ChainSpecificUtil._getBlockhash(uint64(expectedBlockNumber));
      assertEq(expectedBlockHash, actualBlockHash, "incorrect blockhash");
    }
  }

  function testGetBlockhashOptimism() public {
    // Optimism L2 block hash is simply blockhash()
    bytes32 actualBlockhash = ChainSpecificUtil._getBlockhash(uint64(block.number - 1));
    assertEq(blockhash(block.number - 1), actualBlockhash);
  }

  function testGetBlockNumberArbitrum() public {
    uint256[2] memory chainIds = [ARB_MAINNET_CHAIN_ID, ARB_GOERLI_TESTNET_CHAIN_ID];
    uint256[3] memory expectedBlockNumbers = [uint256(10), 11, 12];
    for (uint256 i = 0; i < chainIds.length; i++) {
      vm.chainId(chainIds[i]);
      uint256 expectedBlockNumber = expectedBlockNumbers[i];
      vm.mockCall(ARBSYS_ADDR, abi.encodeWithSelector(ArbSys.arbBlockNumber.selector), abi.encode(expectedBlockNumber));
      uint256 actualBlockNumber = ChainSpecificUtil._getBlockNumber();
      assertEq(expectedBlockNumber, actualBlockNumber, "incorrect block number");
    }
  }

  function testGetBlockNumberOptimism() public {
    // Optimism L2 block number is simply block.number
    uint256 actualBlockNumber = ChainSpecificUtil._getBlockNumber();
    assertEq(block.number, actualBlockNumber);
  }

  function testGetCurrentTxL1GasFeesArbitrum() public {
    uint256[3] memory chainIds = [ARB_MAINNET_CHAIN_ID, ARB_GOERLI_TESTNET_CHAIN_ID, ARB_SEPOLIA_TESTNET_CHAIN_ID];
    uint256[3] memory expectedGasFees = [uint256(10 gwei), 12 gwei, 14 gwei];
    for (uint256 i = 0; i < chainIds.length; i++) {
      vm.chainId(chainIds[i]);
      uint256 expectedGasFee = expectedGasFees[i];
      vm.mockCall(
        ARBGAS_ADDR,
        abi.encodeWithSelector(ArbGasInfo.getCurrentTxL1GasFees.selector),
        abi.encode(expectedGasFee)
      );
      uint256 actualGasFee = ChainSpecificUtil._getCurrentTxL1GasFees("");
      assertEq(expectedGasFee, actualGasFee, "incorrect gas fees");
    }
  }

  function testGetCurrentTxL1GasFeesOptimism() public {
    // set optimism chain id
    uint256[5] memory chainIds = [
      OP_MAINNET_CHAIN_ID,
      OP_GOERLI_CHAIN_ID,
      OP_SEPOLIA_CHAIN_ID,
      BASE_MAINNET_CHAIN_ID,
      BASE_GOERLI_CHAIN_ID
    ];
    uint256[5] memory expectedGasFees = [uint256(10 gwei), 12 gwei, 14 gwei, 16 gwei, 18 gwei];
    for (uint256 i = 0; i < chainIds.length; i++) {
      vm.chainId(chainIds[i]);
      uint256 expectedL1Fee = expectedGasFees[i];
      bytes memory someCalldata = abi.encode(address(0), "blah", uint256(1));
      vm.mockCall(
        OVM_GASPRICEORACLE_ADDR,
        abi.encodeWithSelector(OVM_GasPriceOracle.getL1Fee.selector, bytes.concat(someCalldata, L1_FEE_DATA_PADDING)),
        abi.encode(expectedL1Fee)
      );
      uint256 actualL1Fee = ChainSpecificUtil._getCurrentTxL1GasFees(someCalldata);
      assertEq(expectedL1Fee, actualL1Fee, "incorrect gas fees");
    }
  }

  function testGetL1CalldataGasCostArbitrum() public {
    uint256[3] memory chainIds = [ARB_MAINNET_CHAIN_ID, ARB_GOERLI_TESTNET_CHAIN_ID, ARB_SEPOLIA_TESTNET_CHAIN_ID];
    for (uint256 i = 0; i < chainIds.length; i++) {
      vm.chainId(chainIds[i]);
      vm.mockCall(
        ARBGAS_ADDR,
        abi.encodeWithSelector(ArbGasInfo.getPricesInWei.selector),
        abi.encode(0, 10, 0, 0, 0, 0)
      );

      // fee = l1PricePerByte * (calldataSizeBytes + 140)
      // fee = 10 * (10 + 140) = 1500
      uint256 dataFee = ChainSpecificUtil._getL1CalldataGasCost(10);
      assertEq(dataFee, 1500);
    }
  }

  function testGetL1CalldataGasCostOptimism() public {
    uint256[5] memory chainIds = [
      OP_MAINNET_CHAIN_ID,
      OP_GOERLI_CHAIN_ID,
      OP_SEPOLIA_CHAIN_ID,
      BASE_MAINNET_CHAIN_ID,
      BASE_GOERLI_CHAIN_ID
    ];
    for (uint256 i = 0; i < chainIds.length; i++) {
      vm.chainId(chainIds[i]);
      vm.mockCall(
        OVM_GASPRICEORACLE_ADDR,
        abi.encodeWithSelector(bytes4(hex"519b4bd3")), // l1BaseFee()
        abi.encode(10)
      );
      vm.mockCall(
        OVM_GASPRICEORACLE_ADDR,
        abi.encodeWithSelector(bytes4(hex"0c18c162")), // overhead()
        abi.encode(160)
      );
      vm.mockCall(
        OVM_GASPRICEORACLE_ADDR,
        abi.encodeWithSelector(bytes4(hex"f45e65d8")), // scalar()
        abi.encode(500_000)
      );
      vm.mockCall(
        OVM_GASPRICEORACLE_ADDR,
        abi.encodeWithSelector(bytes4(hex"313ce567")), // decimals()
        abi.encode(6)
      );

      // tx_data_gas = count_zero_bytes(tx_data) * 4 + count_non_zero_bytes(tx_data) * 16
      // tx_data_gas = 0 * 4 + 10 * 16 = 160
      // l1_data_fee = l1_gas_price * (tx_data_gas + fixed_overhead) * dynamic_overhead
      // l1_data_fee = 10 * (160 + 160) * 500_000 / 1_000_000 = 1600
      uint256 dataFee = ChainSpecificUtil._getL1CalldataGasCost(10);
      assertEq(dataFee, 1600);
    }
  }
}
