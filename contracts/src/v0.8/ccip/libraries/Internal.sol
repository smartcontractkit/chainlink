// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {MerkleMultiProof} from "../libraries/MerkleMultiProof.sol";

/// @notice Library for CCIP internal definitions common to multiple contracts.
/// @dev The following is a non-exhaustive list of "known issues" for CCIP:
/// - We could implement yield claiming for Blast. This is not worth the custom code path on non-blast chains.
/// - uint32 is used for timestamps, which will overflow in 2106. This is not a concern for the current use case, as we
/// expect to have migrated to a new version by then.
library Internal {
  error InvalidEVMAddress(bytes encodedAddress);
  error Invalid32ByteAddress(bytes encodedAddress);

  /// @dev We limit return data to a selector plus 4 words. This is to avoid malicious contracts from returning
  /// large amounts of data and causing repeated out-of-gas scenarios.
  uint16 internal constant MAX_RET_BYTES = 4 + 4 * 32;
  /// @dev The expected number of bytes returned by the balanceOf function.
  uint256 internal constant MAX_BALANCE_OF_RET_BYTES = 32;

  /// @dev The address used to send calls for gas estimation.
  /// You only need to use this address if the minimum gas limit specified by the user is not actually enough to execute the
  /// given message and you're attempting to estimate the actual necessary gas limit
  address public constant GAS_ESTIMATION_SENDER = address(0xC11C11C11C11C11C11C11C11C11C11C11C11C1);

  /// @notice A collection of token price and gas price updates.
  struct PriceUpdates {
    TokenPriceUpdate[] tokenPriceUpdates;
    GasPriceUpdate[] gasPriceUpdates;
  }

  /// @notice Token price in USD.
  struct TokenPriceUpdate {
    address sourceToken; // Source token.
    uint224 usdPerToken; // 1e18 USD per 1e18 of the smallest token denomination.
  }

  /// @notice Gas price for a given chain in USD, its value may contain tightly packed fields.
  struct GasPriceUpdate {
    uint64 destChainSelector; // Destination chain selector.
    uint224 usdPerUnitGas; // 1e18 USD per smallest unit (e.g. wei) of destination chain gas.
  }

  /// @notice A timestamped uint224 value that can contain several tightly packed fields.
  struct TimestampedPackedUint224 {
    uint224 value; // ────╮ Value in uint224, packed.
    uint32 timestamp; // ─╯ Timestamp of the most recent price update.
  }

  /// @dev Gas price is stored in 112-bit unsigned int. uint224 can pack 2 prices.
  /// When packing L1 and L2 gas prices, L1 gas price is left-shifted to the higher-order bits.
  /// Using uint8 type, which cannot be higher than other bit shift operands, to avoid shift operand type warning.
  uint8 public constant GAS_PRICE_BITS = 112;

  struct SourceTokenData {
    // The source pool address, abi encoded. This value is trusted as it was obtained through the onRamp. It can be
    // relied upon by the destination pool to validate the source pool.
    bytes sourcePoolAddress;
    // The address of the destination token, abi encoded in the case of EVM chains.
    // This value is UNTRUSTED as any pool owner can return whatever value they want.
    bytes destTokenAddress;
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint32 destGasAmount; // The amount of gas available for the releaseOrMint and balanceOf calls on the offRamp
  }

  /// @notice Report that is submitted by the execution DON at the execution phase, including chain selector data.
  struct ExecutionReport {
    uint64 sourceChainSelector; // Source chain selector for which the report is submitted.
    Any2EVMRampMessage[] messages;
    // Contains a bytes array for each message, each inner bytes array contains bytes per transferred token.
    bytes[][] offchainTokenData;
    bytes32[] proofs;
    uint256 proofFlagBits;
  }

  /// @dev Any2EVMRampMessage struct has 10 fields, including 3 variable unnested arrays, sender, data and tokenAmounts.
  /// Each variable array takes 1 more slot to store its length.
  /// When abi encoded, excluding array contents, Any2EVMMessage takes up a fixed number of 13 slots, 32 bytes each.
  /// Assume 1 slot for sender
  /// For structs that contain arrays, 1 more slot is added to the front, reaching a total of 14.
  /// The fixed bytes does not cover struct data (this is represented by MESSAGE_FIXED_BYTES_PER_TOKEN)
  uint256 public constant MESSAGE_FIXED_BYTES = 32 * 15;

  /// @dev Any2EVMTokensTransfer struct bytes length
  /// 0x20
  /// sourcePoolAddress_offset
  /// destTokenAddress
  /// destGasAmount
  /// extraData_offset
  /// amount
  /// sourcePoolAddress_length
  /// sourcePoolAddress_content // assume 1 slot
  /// extraData_length // contents billed separately
  uint256 public constant MESSAGE_FIXED_BYTES_PER_TOKEN = 32 * (4 + (3 + 2));

  bytes32 internal constant ANY_2_EVM_MESSAGE_HASH = keccak256("Any2EVMMessageHashV1");
  bytes32 internal constant EVM_2_ANY_MESSAGE_HASH = keccak256("EVM2AnyMessageHashV1");

  /// @dev Used to hash messages for multi-lane family-agnostic OffRamps.
  /// OnRamp hash(EVM2AnyMessage) != Any2EVMRampMessage.messageId.
  /// OnRamp hash(EVM2AnyMessage) != OffRamp hash(Any2EVMRampMessage).
  /// @param original OffRamp message to hash.
  /// @param metadataHash Hash preimage to ensure global uniqueness.
  /// @return hashedMessage hashed message as a keccak256.
  function _hash(Any2EVMRampMessage memory original, bytes32 metadataHash) internal pure returns (bytes32) {
    // Fixed-size message fields are included in nested hash to reduce stack pressure.
    return keccak256(
      abi.encode(
        MerkleMultiProof.LEAF_DOMAIN_SEPARATOR,
        metadataHash,
        keccak256(
          abi.encode(
            original.header.messageId,
            original.receiver,
            original.header.sequenceNumber,
            original.gasLimit,
            original.header.nonce
          )
        ),
        keccak256(original.sender),
        keccak256(original.data),
        keccak256(abi.encode(original.tokenAmounts))
      )
    );
  }

  function _hash(EVM2AnyRampMessage memory original, bytes32 metadataHash) internal pure returns (bytes32) {
    // Fixed-size message fields are included in nested hash to reduce stack pressure.
    return keccak256(
      abi.encode(
        MerkleMultiProof.LEAF_DOMAIN_SEPARATOR,
        metadataHash,
        keccak256(
          abi.encode(
            original.sender,
            original.header.sequenceNumber,
            original.header.nonce,
            original.feeToken,
            original.feeTokenAmount
          )
        ),
        keccak256(original.receiver),
        keccak256(original.data),
        keccak256(abi.encode(original.tokenAmounts)),
        keccak256(original.extraArgs)
      )
    );
  }

  /// @dev We disallow the first 1024 addresses to avoid calling into a range known for hosting precompiles. Calling
  /// into precompiles probably won't cause any issues, but to be safe we can disallow this range. It is extremely
  /// unlikely that anyone would ever be able to generate an address in this range. There is no official range of
  /// precompiles, but EIP-7587 proposes to reserve the range 0x100 to 0x1ff. Our range is more conservative, even
  /// though it might not be exhaustive for all chains, which is OK. We also disallow the zero address, which is a
  /// common practice.
  uint256 public constant PRECOMPILE_SPACE = 1024;

  /// @notice This methods provides validation for parsing abi encoded addresses by ensuring the address is within the
  /// EVM address space. If it isn't it will revert with an InvalidEVMAddress error, which we can catch and handle
  /// more gracefully than a revert from abi.decode.
  function _validateEVMAddress(
    bytes memory encodedAddress
  ) internal pure {
    if (encodedAddress.length != 32) revert InvalidEVMAddress(encodedAddress);
    uint256 encodedAddressUint = abi.decode(encodedAddress, (uint256));
    if (encodedAddressUint > type(uint160).max || encodedAddressUint < PRECOMPILE_SPACE) {
      revert InvalidEVMAddress(encodedAddress);
    }
  }

  function _validate32ByteAddress(bytes memory encodedAddress, bool mustBeNonZero) internal pure {
    if (encodedAddress.length != 32) revert Invalid32ByteAddress(encodedAddress);
    if (mustBeNonZero) {
      if (abi.decode(encodedAddress, (bytes32)) == bytes32(0)) {
        revert Invalid32ByteAddress(encodedAddress);
      }
    }
  }

  /// @notice Enum listing the possible message execution states within the offRamp contract.
  /// UNTOUCHED never executed.
  /// IN_PROGRESS currently being executed, used a replay protection.
  /// SUCCESS successfully executed. End state.
  /// FAILURE unsuccessfully executed, manual execution is now enabled.
  enum MessageExecutionState {
    UNTOUCHED,
    IN_PROGRESS,
    SUCCESS,
    FAILURE
  }

  /// @notice CCIP OCR plugin type, used to separate execution & commit transmissions and configs.
  enum OCRPluginType {
    Commit,
    Execution
  }

  /// @notice Family-agnostic header for OnRamp & OffRamp messages.
  /// The messageId is not expected to match hash(message), since it may originate from another ramp family.
  /// Note RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct RampMessageHeader {
    bytes32 messageId; // Unique identifier for the message, generated with the source chain's encoding scheme (i.e. not necessarily abi.encoded).
    uint64 sourceChainSelector; // ─╮ the chain selector of the source chain, note: not chainId.
    uint64 destChainSelector; //    │ the chain selector of the destination chain, note: not chainId.
    uint64 sequenceNumber; //       │ sequence number, not unique across lanes.
    uint64 nonce; // ───────────────╯ nonce for this lane for this sender, not unique across senders/lanes.
  }

  /// Note RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct EVM2AnyTokenTransfer {
    // The source pool EVM address. This value is trusted as it was obtained through the onRamp. It can be relied
    // upon by the destination pool to validate the source pool.
    address sourcePoolAddress;
    // The EVM address of the destination token.
    // This value is UNTRUSTED as any pool owner can return whatever value they want.
    bytes destTokenAddress;
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint256 amount; // Amount of tokens.
    // Destination chain data used to execute the token transfer on the destination chain. For an EVM destination, it
    // consists of the amount of gas available for the releaseOrMint and transfer calls made by the offRamp.
    bytes destExecData;
  }

  /// Note RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct Any2EVMTokenTransfer {
    // The source pool EVM address encoded to bytes. This value is trusted as it is obtained through the onRamp. It can
    // be relied upon by the destination pool to validate the source pool.
    bytes sourcePoolAddress;
    address destTokenAddress; // ─╮ Address of destination token
    uint32 destGasAmount; // ─────╯ The amount of gas available for the releaseOrMint and transfer calls on the offRamp.
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint256 amount; // Amount of tokens.
  }

  /// @notice Family-agnostic message routed to an OffRamp.
  /// Note: hash(Any2EVMRampMessage) != hash(EVM2AnyRampMessage), hash(Any2EVMRampMessage) != messageId due to encoding
  /// and parameter differences.
  /// Note RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct Any2EVMRampMessage {
    RampMessageHeader header; // Message header.
    bytes sender; // sender address on the source chain.
    bytes data; // arbitrary data payload supplied by the message sender.
    address receiver; // receiver address on the destination chain.
    uint256 gasLimit; // user supplied maximum gas amount available for dest chain execution.
    Any2EVMTokenTransfer[] tokenAmounts; // array of tokens and amounts to transfer.
  }

  /// @notice Family-agnostic message emitted from the OnRamp.
  /// Note: hash(Any2EVMRampMessage) != hash(EVM2AnyRampMessage) due to encoding & parameter differences.
  /// messageId = hash(EVM2AnyRampMessage) using the source EVM chain's encoding format.
  /// Note RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct EVM2AnyRampMessage {
    RampMessageHeader header; // Message header.
    address sender; // sender address on the source chain.
    bytes data; // arbitrary data payload supplied by the message sender.
    bytes receiver; // receiver address on the destination chain.
    bytes extraArgs; // destination-chain specific extra args, such as the gasLimit for EVM chains.
    address feeToken; // fee token.
    uint256 feeTokenAmount; // fee token amount.
    uint256 feeValueJuels; // fee amount in Juels.
    EVM2AnyTokenTransfer[] tokenAmounts; // array of tokens and amounts to transfer.
  }

  // bytes4(keccak256("CCIP ChainFamilySelector EVM"));
  bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;

  // bytes4(keccak256("CCIP ChainFamilySelector SVM"));
  bytes4 public constant CHAIN_FAMILY_SELECTOR_SVM = 0x1e10bdc4;

  // bytes4(keccak256("CCIP ChainFamilySelector APTOS"));
  bytes4 public constant CHAIN_FAMILY_SELECTOR_APTOS = 0xac77ffec;

  /// @dev Holds a merkle root and interval for a source chain so that an array of these can be passed in the CommitReport.
  /// @dev inefficient struct packing intentionally chosen to maintain order of specificity. Not a storage struct so impact is minimal.
  // solhint-disable-next-line gas-struct-packing
  struct MerkleRoot {
    uint64 sourceChainSelector; // Remote source chain selector that the Merkle Root is scoped to
    bytes onRampAddress; //        Generic onRamp address, to support arbitrary sources; for EVM, use abi.encode
    uint64 minSeqNr; // ─────────╮ Minimum sequence number, inclusive
    uint64 maxSeqNr; // ─────────╯ Maximum sequence number, inclusive
    bytes32 merkleRoot; //         Merkle root covering the interval & source chain messages
  }
}
