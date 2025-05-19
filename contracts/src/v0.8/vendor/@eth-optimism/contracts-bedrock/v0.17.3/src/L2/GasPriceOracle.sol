// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/L2/GasPriceOracle.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

import {ISemver} from "../universal/ISemver.sol";
import {Predeploys} from "../libraries/Predeploys.sol";
import {L1Block} from "../L2/L1Block.sol";
import {Constants} from "../libraries/Constants.sol";
import {LibZip} from "../deps/LibZip.sol";

/// @custom:proxied
/// @custom:predeploy 0x420000000000000000000000000000000000000F
/// @title GasPriceOracle
/// @notice This contract maintains the variables responsible for computing the L1 portion of the
///         total fee charged on L2. Before Bedrock, this contract held variables in state that were
///         read during the state transition function to compute the L1 portion of the transaction
///         fee. After Bedrock, this contract now simply proxies the L1Block contract, which has
///         the values used to compute the L1 portion of the fee in its state.
///
///         The contract exposes an API that is useful for knowing how large the L1 portion of the
///         transaction fee will be. The following events were deprecated with Bedrock:
///         - event OverheadUpdated(uint256 overhead);
///         - event ScalarUpdated(uint256 scalar);
///         - event DecimalsUpdated(uint256 decimals);
contract GasPriceOracle is ISemver {
  /// @notice Number of decimals used in the scalar.
  uint256 public constant DECIMALS = 6;

  /// @notice Semantic version.
  /// @custom:semver 1.3.0
  string public constant version = "1.3.0";

  /// @notice This is the intercept value for the linear regression used to estimate the final size of the
  ///         compressed transaction.
  int32 private constant COST_INTERCEPT = -42_585_600;

  /// @notice This is the coefficient value for the linear regression used to estimate the final size of the
  ///         compressed transaction.
  uint32 private constant COST_FASTLZ_COEF = 836_500;

  /// @notice This is the minimum bound for the fastlz to brotli size estimation. Any estimations below this
  ///         are set to this value.
  uint256 private constant MIN_TRANSACTION_SIZE = 100;

  /// @notice Indicates whether the network has gone through the Ecotone upgrade.
  bool public isEcotone;

  /// @notice Indicates whether the network has gone through the Fjord upgrade.
  bool public isFjord;

  /// @notice Computes the L1 portion of the fee based on the size of the rlp encoded input
  ///         transaction, the current L1 base fee, and the various dynamic parameters.
  /// @param _data Unsigned fully RLP-encoded transaction to get the L1 fee for.
  /// @return L1 fee that should be paid for the tx
  function getL1Fee(bytes memory _data) external view returns (uint256) {
    if (isFjord) {
      return _getL1FeeFjord(_data);
    } else if (isEcotone) {
      return _getL1FeeEcotone(_data);
    }
    return _getL1FeeBedrock(_data);
  }

  /// @notice returns an upper bound for the L1 fee for a given transaction size.
  /// It is provided for callers who wish to estimate L1 transaction costs in the
  /// write path, and is much more gas efficient than `getL1Fee`.
  /// It assumes the worst case of fastlz upper-bound which covers %99.99 txs.
  /// @param _unsignedTxSize Unsigned fully RLP-encoded transaction size to get the L1 fee for.
  /// @return L1 estimated upper-bound fee that should be paid for the tx
  function getL1FeeUpperBound(uint256 _unsignedTxSize) external view returns (uint256) {
    require(isFjord, "GasPriceOracle: getL1FeeUpperBound only supports Fjord");

    // Add 68 to the size to account for unsigned tx:
    uint256 txSize = _unsignedTxSize + 68;
    // txSize / 255 + 16 is the pratical fastlz upper-bound covers %99.99 txs.
    uint256 flzUpperBound = txSize + txSize / 255 + 16;

    return _fjordL1Cost(flzUpperBound);
  }

  /// @notice Set chain to be Ecotone chain (callable by depositor account)
  function setEcotone() external {
    require(
      msg.sender == Constants.DEPOSITOR_ACCOUNT,
      "GasPriceOracle: only the depositor account can set isEcotone flag"
    );
    require(isEcotone == false, "GasPriceOracle: Ecotone already active");
    isEcotone = true;
  }

  /// @notice Set chain to be Fjord chain (callable by depositor account)
  function setFjord() external {
    require(
      msg.sender == Constants.DEPOSITOR_ACCOUNT,
      "GasPriceOracle: only the depositor account can set isFjord flag"
    );
    require(isEcotone, "GasPriceOracle: Fjord can only be activated after Ecotone");
    require(isFjord == false, "GasPriceOracle: Fjord already active");
    isFjord = true;
  }

  /// @notice Retrieves the current gas price (base fee).
  /// @return Current L2 gas price (base fee).
  function gasPrice() public view returns (uint256) {
    return block.basefee;
  }

  /// @notice Retrieves the current base fee.
  /// @return Current L2 base fee.
  function baseFee() public view returns (uint256) {
    return block.basefee;
  }

  /// @custom:legacy
  /// @notice Retrieves the current fee overhead.
  /// @return Current fee overhead.
  function overhead() public view returns (uint256) {
    require(!isEcotone, "GasPriceOracle: overhead() is deprecated");
    return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).l1FeeOverhead();
  }

  /// @custom:legacy
  /// @notice Retrieves the current fee scalar.
  /// @return Current fee scalar.
  function scalar() public view returns (uint256) {
    require(!isEcotone, "GasPriceOracle: scalar() is deprecated");
    return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).l1FeeScalar();
  }

  /// @notice Retrieves the latest known L1 base fee.
  /// @return Latest known L1 base fee.
  function l1BaseFee() public view returns (uint256) {
    return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).basefee();
  }

  /// @notice Retrieves the current blob base fee.
  /// @return Current blob base fee.
  function blobBaseFee() public view returns (uint256) {
    return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).blobBaseFee();
  }

  /// @notice Retrieves the current base fee scalar.
  /// @return Current base fee scalar.
  function baseFeeScalar() public view returns (uint32) {
    return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).baseFeeScalar();
  }

  /// @notice Retrieves the current blob base fee scalar.
  /// @return Current blob base fee scalar.
  function blobBaseFeeScalar() public view returns (uint32) {
    return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).blobBaseFeeScalar();
  }

  /// @custom:legacy
  /// @notice Retrieves the number of decimals used in the scalar.
  /// @return Number of decimals used in the scalar.
  function decimals() public pure returns (uint256) {
    return DECIMALS;
  }

  /// @notice Computes the amount of L1 gas used for a transaction. Adds 68 bytes
  ///         of padding to account for the fact that the input does not have a signature.
  /// @param _data Unsigned fully RLP-encoded transaction to get the L1 gas for.
  /// @return Amount of L1 gas used to publish the transaction.
  /// @custom:deprecated This method does not accurately estimate the gas used for a transaction.
  ///                    If you are calculating fees use getL1Fee or getL1FeeUpperBound.
  function getL1GasUsed(bytes memory _data) public view returns (uint256) {
    if (isFjord) {
      // Add 68 to the size to account for unsigned tx
      // Assume the compressed data is mostly non-zero, and would pay 16 gas per calldata byte
      // Divide by 1e6 due to the scaling factor of the linear regression
      return (_fjordLinearRegression(LibZip.flzCompress(_data).length + 68) * 16) / 1e6;
    }
    uint256 l1GasUsed = _getCalldataGas(_data);
    if (isEcotone) {
      return l1GasUsed;
    }
    return l1GasUsed + L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).l1FeeOverhead();
  }

  /// @notice Computation of the L1 portion of the fee for Bedrock.
  /// @param _data Unsigned fully RLP-encoded transaction to get the L1 fee for.
  /// @return L1 fee that should be paid for the tx
  function _getL1FeeBedrock(bytes memory _data) internal view returns (uint256) {
    uint256 l1GasUsed = _getCalldataGas(_data);
    uint256 fee = (l1GasUsed + L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).l1FeeOverhead()) *
      l1BaseFee() *
      L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).l1FeeScalar();
    return fee / (10 ** DECIMALS);
  }

  /// @notice L1 portion of the fee after Ecotone.
  /// @param _data Unsigned fully RLP-encoded transaction to get the L1 fee for.
  /// @return L1 fee that should be paid for the tx
  function _getL1FeeEcotone(bytes memory _data) internal view returns (uint256) {
    uint256 l1GasUsed = _getCalldataGas(_data);
    uint256 scaledBaseFee = baseFeeScalar() * 16 * l1BaseFee();
    uint256 scaledBlobBaseFee = blobBaseFeeScalar() * blobBaseFee();
    uint256 fee = l1GasUsed * (scaledBaseFee + scaledBlobBaseFee);
    return fee / (16 * 10 ** DECIMALS);
  }

  /// @notice L1 portion of the fee after Fjord.
  /// @param _data Unsigned fully RLP-encoded transaction to get the L1 fee for.
  /// @return L1 fee that should be paid for the tx
  function _getL1FeeFjord(bytes memory _data) internal view returns (uint256) {
    return _fjordL1Cost(LibZip.flzCompress(_data).length + 68);
  }

  /// @notice L1 gas estimation calculation.
  /// @param _data Unsigned fully RLP-encoded transaction to get the L1 gas for.
  /// @return Amount of L1 gas used to publish the transaction.
  function _getCalldataGas(bytes memory _data) internal pure returns (uint256) {
    uint256 total = 0;
    uint256 length = _data.length;
    for (uint256 i = 0; i < length; i++) {
      if (_data[i] == 0) {
        total += 4;
      } else {
        total += 16;
      }
    }
    return total + (68 * 16);
  }

  /// @notice Fjord L1 cost based on the compressed and original tx size.
  /// @param _fastLzSize estimated compressed tx size.
  /// @return Fjord L1 fee that should be paid for the tx
  function _fjordL1Cost(uint256 _fastLzSize) internal view returns (uint256) {
    // Apply the linear regression to estimate the Brotli 10 size
    uint256 estimatedSize = _fjordLinearRegression(_fastLzSize);
    uint256 feeScaled = baseFeeScalar() * 16 * l1BaseFee() + blobBaseFeeScalar() * blobBaseFee();
    return (estimatedSize * feeScaled) / (10 ** (DECIMALS * 2));
  }

  /// @notice Takes the fastLz size compression and returns the estimated Brotli
  /// @param _fastLzSize fastlz compressed tx size.
  /// @return Number of bytes in the compressed transaction
  function _fjordLinearRegression(uint256 _fastLzSize) internal pure returns (uint256) {
    int256 estimatedSize = COST_INTERCEPT + int256(COST_FASTLZ_COEF * _fastLzSize);
    if (estimatedSize < int256(MIN_TRANSACTION_SIZE) * 1e6) {
      estimatedSize = int256(MIN_TRANSACTION_SIZE) * 1e6;
    }
    return uint256(estimatedSize);
  }
}
