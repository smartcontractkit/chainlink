// SPDX-License-Identifier: MIT

pragma solidity 0.8.24;

import {IBootloaderUtilities} from "./interfaces/IBootloaderUtilities.sol";
import {Transaction, TransactionHelper, EIP_712_TX_TYPE, LEGACY_TX_TYPE, EIP_2930_TX_TYPE, EIP_1559_TX_TYPE} from "./libraries/TransactionHelper.sol";
import {RLPEncoder} from "./libraries/RLPEncoder.sol";
import {EfficientCall} from "./libraries/EfficientCall.sol";
import {UnsupportedTxType, InvalidSig, SigField} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice A contract that provides some utility methods for the bootloader
 * that is very hard to write in Yul.
 */
contract BootloaderUtilities is IBootloaderUtilities {
    using TransactionHelper for *;

    /// @notice Calculates the canonical transaction hash and the recommended transaction hash.
    /// @param _transaction The transaction.
    /// @return txHash and signedTxHash of the transaction, i.e. the transaction hash to be used in the explorer and commits to all
    /// the fields of the transaction and the recommended hash to be signed for this transaction.
    /// @dev txHash must be unique for all transactions.
    function getTransactionHashes(
        Transaction calldata _transaction
    ) external view override returns (bytes32 txHash, bytes32 signedTxHash) {
        signedTxHash = _transaction.encodeHash();
        if (_transaction.txType == EIP_712_TX_TYPE) {
            txHash = keccak256(bytes.concat(signedTxHash, EfficientCall.keccak(_transaction.signature)));
        } else if (_transaction.txType == LEGACY_TX_TYPE) {
            txHash = encodeLegacyTransactionHash(_transaction);
        } else if (_transaction.txType == EIP_1559_TX_TYPE) {
            txHash = encodeEIP1559TransactionHash(_transaction);
        } else if (_transaction.txType == EIP_2930_TX_TYPE) {
            txHash = encodeEIP2930TransactionHash(_transaction);
        } else {
            revert UnsupportedTxType(_transaction.txType);
        }
    }

    /// @notice Calculates the hash for a legacy transaction.
    /// @param _transaction The legacy transaction.
    /// @return txHash The hash of the transaction.
    function encodeLegacyTransactionHash(Transaction calldata _transaction) internal view returns (bytes32 txHash) {
        // Hash of legacy transactions are encoded as one of the:
        // - RLP(nonce, gasPrice, gasLimit, to, value, data, chainId, 0, 0)
        // - RLP(nonce, gasPrice, gasLimit, to, value, data)
        //
        // In this RLP encoding, only the first one above list appears, so we encode each element
        // inside list and then concatenate the length of all elements with them.

        bytes memory encodedNonce = RLPEncoder.encodeUint256(_transaction.nonce);
        // Encode `gasPrice` and `gasLimit` together to prevent "stack too deep error".
        bytes memory encodedGasParam;
        {
            bytes memory encodedGasPrice = RLPEncoder.encodeUint256(_transaction.maxFeePerGas);
            bytes memory encodedGasLimit = RLPEncoder.encodeUint256(_transaction.gasLimit);
            encodedGasParam = bytes.concat(encodedGasPrice, encodedGasLimit);
        }

        bytes memory encodedTo = RLPEncoder.encodeAddress(address(uint160(_transaction.to)));
        bytes memory encodedValue = RLPEncoder.encodeUint256(_transaction.value);
        // Encode only the length of the transaction data, and not the data itself,
        // so as not to copy to memory a potentially huge transaction data twice.
        bytes memory encodedDataLength;
        {
            // Safe cast, because the length of the transaction data can't be so large.
            uint64 txDataLen = uint64(_transaction.data.length);
            if (txDataLen != 1) {
                // If the length is not equal to one, then only using the length can it be encoded definitely.
                encodedDataLength = RLPEncoder.encodeNonSingleBytesLen(txDataLen);
            } else if (_transaction.data[0] >= 0x80) {
                // If input is a byte in [0x80, 0xff] range, RLP encoding will concatenates 0x81 with the byte.
                encodedDataLength = hex"81";
            }
            // Otherwise the length is not encoded at all.
        }

        bytes memory rEncoded;
        {
            uint256 rInt = uint256(bytes32(_transaction.signature[0:32]));
            rEncoded = RLPEncoder.encodeUint256(rInt);
        }
        bytes memory sEncoded;
        {
            uint256 sInt = uint256(bytes32(_transaction.signature[32:64]));
            sEncoded = RLPEncoder.encodeUint256(sInt);
        }
        bytes memory vEncoded;
        {
            uint256 vInt = uint256(uint8(_transaction.signature[64]));
            if (vInt != 27 && vInt != 28) {
                revert InvalidSig(SigField.V, vInt);
            }

            // If the `chainId` is specified in the transaction, then the `v` value is encoded as
            // `35 + y + 2 * chainId == vInt + 8 + 2 * chainId`, where y - parity bit (see EIP-155).
            if (_transaction.reserved[0] != 0) {
                vInt += 8 + block.chainid * 2;
            }

            vEncoded = RLPEncoder.encodeUint256(vInt);
        }

        bytes memory encodedListLength;
        unchecked {
            uint256 listLength = encodedNonce.length +
                encodedGasParam.length +
                encodedTo.length +
                encodedValue.length +
                encodedDataLength.length +
                _transaction.data.length +
                rEncoded.length +
                sEncoded.length +
                vEncoded.length;

            // Safe cast, because the length of the list can't be so large.
            encodedListLength = RLPEncoder.encodeListLen(uint64(listLength));
        }

        return
            keccak256(
                // solhint-disable-next-line func-named-parameters
                bytes.concat(
                    encodedListLength,
                    encodedNonce,
                    encodedGasParam,
                    encodedTo,
                    encodedValue,
                    encodedDataLength,
                    _transaction.data,
                    vEncoded,
                    rEncoded,
                    sEncoded
                )
            );
    }

    /// @notice Calculates the hash for an EIP2930 transaction.
    /// @param _transaction The EIP2930 transaction.
    /// @return txHash The hash of the transaction.
    function encodeEIP2930TransactionHash(Transaction calldata _transaction) internal view returns (bytes32) {
        // Encode all fixed-length params to avoid "stack too deep error"
        bytes memory encodedFixedLengthParams;
        {
            bytes memory encodedChainId = RLPEncoder.encodeUint256(block.chainid);
            bytes memory encodedNonce = RLPEncoder.encodeUint256(_transaction.nonce);
            bytes memory encodedGasPrice = RLPEncoder.encodeUint256(_transaction.maxFeePerGas);
            bytes memory encodedGasLimit = RLPEncoder.encodeUint256(_transaction.gasLimit);
            bytes memory encodedTo = RLPEncoder.encodeAddress(address(uint160(_transaction.to)));
            bytes memory encodedValue = RLPEncoder.encodeUint256(_transaction.value);
            // solhint-disable-next-line func-named-parameters
            encodedFixedLengthParams = bytes.concat(
                encodedChainId,
                encodedNonce,
                encodedGasPrice,
                encodedGasLimit,
                encodedTo,
                encodedValue
            );
        }

        // Encode only the length of the transaction data, and not the data itself,
        // so as not to copy to memory a potentially huge transaction data twice.
        bytes memory encodedDataLength;
        {
            // Safe cast, because the length of the transaction data can't be so large.
            uint64 txDataLen = uint64(_transaction.data.length);
            if (txDataLen != 1) {
                // If the length is not equal to one, then only using the length can it be encoded definitely.
                encodedDataLength = RLPEncoder.encodeNonSingleBytesLen(txDataLen);
            } else if (_transaction.data[0] >= 0x80) {
                // If input is a byte in [0x80, 0xff] range, RLP encoding will concatenates 0x81 with the byte.
                encodedDataLength = hex"81";
            }
            // Otherwise the length is not encoded at all.
        }

        // On ZKsync, access lists are always zero length (at least for now).
        bytes memory encodedAccessListLength = RLPEncoder.encodeListLen(0);

        bytes memory rEncoded;
        {
            uint256 rInt = uint256(bytes32(_transaction.signature[0:32]));
            rEncoded = RLPEncoder.encodeUint256(rInt);
        }
        bytes memory sEncoded;
        {
            uint256 sInt = uint256(bytes32(_transaction.signature[32:64]));
            sEncoded = RLPEncoder.encodeUint256(sInt);
        }
        bytes memory vEncoded;
        {
            uint256 vInt = uint256(uint8(_transaction.signature[64]));
            if (vInt != 27 && vInt != 28) {
                revert InvalidSig(SigField.V, vInt);
            }

            vEncoded = RLPEncoder.encodeUint256(vInt - 27);
        }

        bytes memory encodedListLength;
        unchecked {
            uint256 listLength = encodedFixedLengthParams.length +
                encodedDataLength.length +
                _transaction.data.length +
                encodedAccessListLength.length +
                rEncoded.length +
                sEncoded.length +
                vEncoded.length;

            // Safe cast, because the length of the list can't be so large.
            encodedListLength = RLPEncoder.encodeListLen(uint64(listLength));
        }

        return
            keccak256(
                // solhint-disable-next-line func-named-parameters
                bytes.concat(
                    "\x01",
                    encodedListLength,
                    encodedFixedLengthParams,
                    encodedDataLength,
                    _transaction.data,
                    encodedAccessListLength,
                    vEncoded,
                    rEncoded,
                    sEncoded
                )
            );
    }

    /// @notice Calculates the hash for an EIP1559 transaction.
    /// @param _transaction The legacy transaction.
    /// @return txHash The hash of the transaction.
    function encodeEIP1559TransactionHash(Transaction calldata _transaction) internal view returns (bytes32) {
        // The formula for hash of EIP1559 transaction in the original proposal:
        // https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md

        // Encode all fixed-length params to avoid "stack too deep error"
        bytes memory encodedFixedLengthParams;
        {
            bytes memory encodedChainId = RLPEncoder.encodeUint256(block.chainid);
            bytes memory encodedNonce = RLPEncoder.encodeUint256(_transaction.nonce);
            bytes memory encodedMaxPriorityFeePerGas = RLPEncoder.encodeUint256(_transaction.maxPriorityFeePerGas);
            bytes memory encodedMaxFeePerGas = RLPEncoder.encodeUint256(_transaction.maxFeePerGas);
            bytes memory encodedGasLimit = RLPEncoder.encodeUint256(_transaction.gasLimit);
            bytes memory encodedTo = RLPEncoder.encodeAddress(address(uint160(_transaction.to)));
            bytes memory encodedValue = RLPEncoder.encodeUint256(_transaction.value);
            // solhint-disable-next-line func-named-parameters
            encodedFixedLengthParams = bytes.concat(
                encodedChainId,
                encodedNonce,
                encodedMaxPriorityFeePerGas,
                encodedMaxFeePerGas,
                encodedGasLimit,
                encodedTo,
                encodedValue
            );
        }

        // Encode only the length of the transaction data, and not the data itself,
        // so as not to copy to memory a potentially huge transaction data twice.
        bytes memory encodedDataLength;
        {
            // Safe cast, because the length of the transaction data can't be so large.
            uint64 txDataLen = uint64(_transaction.data.length);
            if (txDataLen != 1) {
                // If the length is not equal to one, then only using the length can it be encoded definitely.
                encodedDataLength = RLPEncoder.encodeNonSingleBytesLen(txDataLen);
            } else if (_transaction.data[0] >= 0x80) {
                // If input is a byte in [0x80, 0xff] range, RLP encoding will concatenates 0x81 with the byte.
                encodedDataLength = hex"81";
            }
            // Otherwise the length is not encoded at all.
        }

        // On ZKsync, access lists are always zero length (at least for now).
        bytes memory encodedAccessListLength = RLPEncoder.encodeListLen(0);

        bytes memory rEncoded;
        {
            uint256 rInt = uint256(bytes32(_transaction.signature[0:32]));
            rEncoded = RLPEncoder.encodeUint256(rInt);
        }
        bytes memory sEncoded;
        {
            uint256 sInt = uint256(bytes32(_transaction.signature[32:64]));
            sEncoded = RLPEncoder.encodeUint256(sInt);
        }
        bytes memory vEncoded;
        {
            uint256 vInt = uint256(uint8(_transaction.signature[64]));
            if (vInt != 27 && vInt != 28) {
                revert InvalidSig(SigField.V, vInt);
            }

            vEncoded = RLPEncoder.encodeUint256(vInt - 27);
        }

        bytes memory encodedListLength;
        unchecked {
            uint256 listLength = encodedFixedLengthParams.length +
                encodedDataLength.length +
                _transaction.data.length +
                encodedAccessListLength.length +
                rEncoded.length +
                sEncoded.length +
                vEncoded.length;

            // Safe cast, because the length of the list can't be so large.
            encodedListLength = RLPEncoder.encodeListLen(uint64(listLength));
        }

        return
            keccak256(
                // solhint-disable-next-line func-named-parameters
                bytes.concat(
                    "\x02",
                    encodedListLength,
                    encodedFixedLengthParams,
                    encodedDataLength,
                    _transaction.data,
                    encodedAccessListLength,
                    vEncoded,
                    rEncoded,
                    sEncoded
                )
            );
    }
}
