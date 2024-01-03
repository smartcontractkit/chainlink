// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Metis_GasPriceOracle} from "./../../vendor/@metis/contracts/L2/predeploys/Metis_GasPriceOracle.sol";
import {IScrollL1GasPriceOracle} from "./../../vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol";

contract ChainCostTester {
    event BlockLog(uint256 block, uint256 timestamp, bytes32 hash);
    event GasDetails(uint256 nativeGasUsed, uint256 l1Cost, uint256 l1DataLength, uint256 l1BaseFee, uint256 nativeGasPriceFromOracle, uint256 txGasPrice);

    Metis_GasPriceOracle public immutable METIS_ORACLE = Metis_GasPriceOracle(0x420000000000000000000000000000000000000F);
    IScrollL1GasPriceOracle public immutable SCROLL_ORACLE = IScrollL1GasPriceOracle(0x5300000000000000000000000000000000000002);
    bytes public data;
    uint256 public sum;
    bytes public scrollPadding;
    bytes public metisPadding;
    bytes public gnosisPadding;
    bytes public zkSyncPadding;
    bytes public zkEVMPadding;
    bytes public celoPadding;

    function setScrollPadding(bytes calldata _padding) public {
        scrollPadding = _padding;
    }

    function setMetisPadding(bytes calldata _padding) public {
        metisPadding = _padding;
    }

    function setGnosisPadding(bytes calldata _padding) public {
        gnosisPadding = _padding;
    }

    function setZKSyncPadding(bytes calldata _padding) public {
        zkSyncPadding = _padding;
    }

    function setZKEVMPadding(bytes calldata _padding) public {
        zkEVMPadding = _padding;
    }

    function setCeloPadding(bytes calldata _padding) public {
        celoPadding = _padding;
    }

    function computeAndUpdateData(bytes calldata _data, uint256 _n) public {
        uint256 g1 = gasleft();
        emit BlockLog(block.number, block.timestamp, blockhash(block.number - 1));

        // some computation
        uint256 _sum = 0;
        for (uint256 i = 0; i < _n; i++) {
            _sum += i;
        }
        sum = _sum;

        // store bytes
        data = _data;

        uint256 l1Cost;
        uint256 l1BaseFee;
        uint256 nativeGasPriceFromOracle;
        bytes memory totalBytes;
        if (isMetis()) {
            totalBytes = bytes.concat(_data, metisPadding);
            l1Cost = METIS_ORACLE.getL1Fee(totalBytes);
            l1BaseFee = METIS_ORACLE.l1BaseFee();
            nativeGasPriceFromOracle = METIS_ORACLE.gasPrice();
        } else if (isScroll()) {
            totalBytes = bytes.concat(_data, scrollPadding);
            l1Cost = SCROLL_ORACLE.getL1Fee(totalBytes);
            l1BaseFee = SCROLL_ORACLE.l1BaseFee();
            nativeGasPriceFromOracle = 0; // the oracle does not expose this
        } else if (isGnosis()) {
            totalBytes = bytes.concat(_data, gnosisPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isPolygonZKEVM()) {
            totalBytes = bytes.concat(_data, zkEVMPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isZKSync()) {
            totalBytes = bytes.concat(_data, zkSyncPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        } else if (isCelo()) {
            totalBytes = bytes.concat(_data, celoPadding);
            l1Cost = 0;
            l1BaseFee = 0;
            nativeGasPriceFromOracle = 0;
        }

        uint256 g2 = gasleft();
        emit GasDetails(g1 - g2, l1Cost, totalBytes.length, l1BaseFee, nativeGasPriceFromOracle, tx.gasprice);
    }

//    function updateData(bytes calldata _data) public {
//        uint256 g1 = gasLeft();
//        emit BlockLog(block.number, block.timestamp, blockhash(block.number));
//
//        // store bytes
//        data = _data;
//
//        uint256 g2 = gasLeft();
//        emit GasDetails(g2 - g1, );
//    }
//
//    function compute(uint256 _n) public {
//        uint256 g1 = gasLeft();
//        emit BlockLog(block.number, block.timestamp, blockhash(block.number));
//
//        // some computation
//        uint256 _sum = 0;
//        for (uint256 i = 0; i < _n; i++) {
//            _sum += i;
//        }
//        sum = _sum;
//
//        uint256 g2 = gasLeft();
//        emit GasDetails(g2 - g1, );
//    }

    function isScroll() public view returns (bool) {
        // Scroll Sepolia or Scroll mainnet
        return block.chainid == 534351 || block.chainid == 534352;
    }

    function isMetis() public view returns (bool) {
        // Metis Goerli or Metis Andromeda mainnet
        return block.chainid == 599 || block.chainid == 1088;
    }

    function isGnosis() public view returns (bool) {
        // Gnosis Chiado or Gnosis mainnet
        return block.chainid == 10200 || block.chainid == 100;
    }

    function isZKSync() public view returns (bool) {
        // zkSync Sepolia or zkSync mainnet
        return block.chainid == 300 || block.chainid == 324;
    }

    function isPolygonZKEVM() public view returns (bool) {
        // zkEVM testnet or zkEVM mainnet
        return block.chainid == 1442 || block.chainid == 1101;
    }

    function isCelo() public view returns (bool) {
        // celo Alfajores testnet or celo mainnet
        return block.chainid == 44787 || block.chainid == 42220;
    }
}