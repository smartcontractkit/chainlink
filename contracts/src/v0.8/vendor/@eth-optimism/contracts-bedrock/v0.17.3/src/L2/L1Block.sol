// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/L2/L1Block.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

import {ISemver} from "../universal/ISemver.sol";
import {Constants} from "../libraries/Constants.sol";
import {GasPayingToken, IGasToken} from "../libraries/GasPayingToken.sol";
import "../libraries/L1BlockErrors.sol";

/// @custom:proxied
/// @custom:predeploy 0x4200000000000000000000000000000000000015
/// @title L1Block
/// @notice The L1Block predeploy gives users access to information about the last known L1 block.
///         Values within this contract are updated once per epoch (every L1 block) and can only be
///         set by the "depositor" account, a special system address. Depositor account transactions
///         are created by the protocol whenever we move to a new epoch.
contract L1Block is ISemver, IGasToken {
  /// @notice Event emitted when the gas paying token is set.
  event GasPayingTokenSet(address indexed token, uint8 indexed decimals, bytes32 name, bytes32 symbol);

  /// @notice Address of the special depositor account.
  function DEPOSITOR_ACCOUNT() public pure returns (address addr_) {
    addr_ = Constants.DEPOSITOR_ACCOUNT;
  }

  /// @notice The latest L1 block number known by the L2 system.
  uint64 public number;

  /// @notice The latest L1 timestamp known by the L2 system.
  uint64 public timestamp;

  /// @notice The latest L1 base fee.
  uint256 public basefee;

  /// @notice The latest L1 blockhash.
  bytes32 public hash;

  /// @notice The number of L2 blocks in the same epoch.
  uint64 public sequenceNumber;

  /// @notice The scalar value applied to the L1 blob base fee portion of the blob-capable L1 cost func.
  uint32 public blobBaseFeeScalar;

  /// @notice The scalar value applied to the L1 base fee portion of the blob-capable L1 cost func.
  uint32 public baseFeeScalar;

  /// @notice The versioned hash to authenticate the batcher by.
  bytes32 public batcherHash;

  /// @notice The overhead value applied to the L1 portion of the transaction fee.
  /// @custom:legacy
  uint256 public l1FeeOverhead;

  /// @notice The scalar value applied to the L1 portion of the transaction fee.
  /// @custom:legacy
  uint256 public l1FeeScalar;

  /// @notice The latest L1 blob base fee.
  uint256 public blobBaseFee;

  /// @custom:semver 1.4.1-beta.1
  function version() public pure virtual returns (string memory) {
    return "1.4.1-beta.1";
  }

  /// @notice Returns the gas paying token, its decimals, name and symbol.
  ///         If nothing is set in state, then it means ether is used.
  function gasPayingToken() public view returns (address addr_, uint8 decimals_) {
    (addr_, decimals_) = GasPayingToken.getToken();
  }

  /// @notice Returns the gas paying token name.
  ///         If nothing is set in state, then it means ether is used.
  function gasPayingTokenName() public view returns (string memory name_) {
    name_ = GasPayingToken.getName();
  }

  /// @notice Returns the gas paying token symbol.
  ///         If nothing is set in state, then it means ether is used.
  function gasPayingTokenSymbol() public view returns (string memory symbol_) {
    symbol_ = GasPayingToken.getSymbol();
  }

  /// @notice Getter for custom gas token paying networks. Returns true if the
  ///         network uses a custom gas token.
  function isCustomGasToken() public view returns (bool) {
    (address token, ) = gasPayingToken();
    return token != Constants.ETHER;
  }

  /// @custom:legacy
  /// @notice Updates the L1 block values.
  /// @param _number         L1 blocknumber.
  /// @param _timestamp      L1 timestamp.
  /// @param _basefee        L1 basefee.
  /// @param _hash           L1 blockhash.
  /// @param _sequenceNumber Number of L2 blocks since epoch start.
  /// @param _batcherHash    Versioned hash to authenticate batcher by.
  /// @param _l1FeeOverhead  L1 fee overhead.
  /// @param _l1FeeScalar    L1 fee scalar.
  function setL1BlockValues(
    uint64 _number,
    uint64 _timestamp,
    uint256 _basefee,
    bytes32 _hash,
    uint64 _sequenceNumber,
    bytes32 _batcherHash,
    uint256 _l1FeeOverhead,
    uint256 _l1FeeScalar
  ) external {
    require(msg.sender == DEPOSITOR_ACCOUNT(), "L1Block: only the depositor account can set L1 block values");

    number = _number;
    timestamp = _timestamp;
    basefee = _basefee;
    hash = _hash;
    sequenceNumber = _sequenceNumber;
    batcherHash = _batcherHash;
    l1FeeOverhead = _l1FeeOverhead;
    l1FeeScalar = _l1FeeScalar;
  }

  /// @notice Updates the L1 block values for an Ecotone upgraded chain.
  /// Params are packed and passed in as raw msg.data instead of ABI to reduce calldata size.
  /// Params are expected to be in the following order:
  ///   1. _baseFeeScalar      L1 base fee scalar
  ///   2. _blobBaseFeeScalar  L1 blob base fee scalar
  ///   3. _sequenceNumber     Number of L2 blocks since epoch start.
  ///   4. _timestamp          L1 timestamp.
  ///   5. _number             L1 blocknumber.
  ///   6. _basefee            L1 base fee.
  ///   7. _blobBaseFee        L1 blob base fee.
  ///   8. _hash               L1 blockhash.
  ///   9. _batcherHash        Versioned hash to authenticate batcher by.
  function setL1BlockValuesEcotone() external {
    address depositor = DEPOSITOR_ACCOUNT();
    assembly {
      // Revert if the caller is not the depositor account.
      if xor(caller(), depositor) {
        mstore(0x00, 0x3cc50b45) // 0x3cc50b45 is the 4-byte selector of "NotDepositor()"
        revert(0x1C, 0x04) // returns the stored 4-byte selector from above
      }
      // sequencenum (uint64), blobBaseFeeScalar (uint32), baseFeeScalar (uint32)
      sstore(sequenceNumber.slot, shr(128, calldataload(4)))
      // number (uint64) and timestamp (uint64)
      sstore(number.slot, shr(128, calldataload(20)))
      sstore(basefee.slot, calldataload(36)) // uint256
      sstore(blobBaseFee.slot, calldataload(68)) // uint256
      sstore(hash.slot, calldataload(100)) // bytes32
      sstore(batcherHash.slot, calldataload(132)) // bytes32
    }
  }

  /// @notice Sets the gas paying token for the L2 system. Can only be called by the special
  ///         depositor account. This function is not called on every L2 block but instead
  ///         only called by specially crafted L1 deposit transactions.
  function setGasPayingToken(address _token, uint8 _decimals, bytes32 _name, bytes32 _symbol) external {
    if (msg.sender != DEPOSITOR_ACCOUNT()) revert NotDepositor();

    GasPayingToken.set({_token: _token, _decimals: _decimals, _name: _name, _symbol: _symbol});

    emit GasPayingTokenSet({token: _token, decimals: _decimals, name: _name, symbol: _symbol});
  }
}
