// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {MerkleMultiProof} from "../libraries/MerkleMultiProof.sol";
import {Client} from "./Client.sol";

// Library for CCIP internal definitions common to multiple contracts.
library Internal {
  error InvalidEVMAddress(bytes encodedAddress);

  /// @dev The minimum amount of gas to perform the call with exact gas.
  /// We include this in the offramp so that we can redeploy to adjust it
  /// should a hardfork change the gas costs of relevant opcodes in callWithExactGas.
  uint16 internal constant GAS_FOR_CALL_EXACT_CHECK = 5_000;
  // @dev We limit return data to a selector plus 4 words. This is to avoid
  // malicious contracts from returning large amounts of data and causing
  // repeated out-of-gas scenarios.
  uint16 internal constant MAX_RET_BYTES = 4 + 4 * 32;

  /// @notice A collection of token price and gas price updates.
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct PriceUpdates {
    TokenPriceUpdate[] tokenPriceUpdates;
    GasPriceUpdate[] gasPriceUpdates;
  }

  /// @notice Token price in USD.
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct TokenPriceUpdate {
    address sourceToken; // Source token
    uint224 usdPerToken; // 1e18 USD per 1e18 of the smallest token denomination.
  }

  /// @notice Gas price for a given chain in USD, its value may contain tightly packed fields.
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct GasPriceUpdate {
    uint64 destChainSelector; // Destination chain selector
    uint224 usdPerUnitGas; // 1e18 USD per smallest unit (e.g. wei) of destination chain gas
  }

  /// @notice A timestamped uint224 value that can contain several tightly packed fields.
  struct TimestampedPackedUint224 {
    uint224 value; // ───────╮ Value in uint224, packed.
    uint32 timestamp; // ────╯ Timestamp of the most recent price update.
  }

  /// @dev Gas price is stored in 112-bit unsigned int. uint224 can pack 2 prices.
  /// When packing L1 and L2 gas prices, L1 gas price is left-shifted to the higher-order bits.
  /// Using uint8 type, which cannot be higher than other bit shift operands, to avoid shift operand type warning.
  uint8 public constant GAS_PRICE_BITS = 112;

  struct PoolUpdate {
    address token; // The IERC20 token address
    address pool; // The token pool address
  }

  struct SourceTokenData {
    // The source pool address, abi encoded. This value is trusted as it was obtained through the onRamp. It can be
    // relied upon by the destination pool to validate the source pool.
    bytes sourcePoolAddress;
    // The address of the destination token, abi encoded in the case of EVM chains
    // This value is UNTRUSTED as any pool owner can return whatever value they want.
    bytes destTokenAddress;
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint32 destGasAmount; // The amount of gas available for the releaseOrMint and transfer calls on the offRamp
  }

  /// @notice Report that is submitted by the execution DON at the execution phase. (including chain selector data)
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct ExecutionReportSingleChain {
    uint64 sourceChainSelector; // Source chain selector for which the report is submitted
    Any2EVMRampMessage[] messages;
    // Contains a bytes array for each message, each inner bytes array contains bytes per transferred token
    bytes[][] offchainTokenData;
    bytes32[] proofs;
    uint256 proofFlagBits;
  }

  /// @notice Report that is submitted by the execution DON at the execution phase.
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct ExecutionReport {
    EVM2EVMMessage[] messages;
    // Contains a bytes array for each message, each inner bytes array contains bytes per transferred token
    bytes[][] offchainTokenData;
    bytes32[] proofs;
    uint256 proofFlagBits;
  }

  /// @notice The cross chain message that gets committed to EVM chains.
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct EVM2EVMMessage {
    uint64 sourceChainSelector; // ────────╮ the chain selector of the source chain, note: not chainId
    address sender; // ────────────────────╯ sender address on the source chain
    address receiver; // ──────────────────╮ receiver address on the destination chain
    uint64 sequenceNumber; // ─────────────╯ sequence number, not unique across lanes
    uint256 gasLimit; //                     user supplied maximum gas amount available for dest chain execution
    bool strict; // ───────────────────────╮ DEPRECATED
    uint64 nonce; //                       │ nonce for this lane for this sender, not unique across senders/lanes
    address feeToken; // ──────────────────╯ fee token
    uint256 feeTokenAmount; //               fee token amount
    bytes data; //                           arbitrary data payload supplied by the message sender
    Client.EVMTokenAmount[] tokenAmounts; // array of tokens and amounts to transfer
    bytes[] sourceTokenData; //              array of token data, one per token
    bytes32 messageId; //                    a hash of the message data
  }

  /// @dev EVM2EVMMessage struct has 13 fields, including 3 variable arrays.
  /// Each variable array takes 1 more slot to store its length.
  /// When abi encoded, excluding array contents,
  /// EVM2EVMMessage takes up a fixed number of 16 lots, 32 bytes each.
  /// For structs that contain arrays, 1 more slot is added to the front, reaching a total of 17.
  uint256 public constant MESSAGE_FIXED_BYTES = 32 * 17;

  /// @dev Each token transfer adds 1 EVMTokenAmount and 3 bytes at 3 slots each and one slot for the destGasAmount.
  /// When abi encoded, each EVMTokenAmount takes 2 slots, each bytes takes 1 slot for length, one slot of data and one
  /// slot for the offset. This results in effectively 3*3 slots per SourceTokenData.
  /// 0x20
  /// destGasAmount
  /// sourcePoolAddress_offset
  /// destTokenAddress_offset
  /// extraData_offset
  /// sourcePoolAddress_length
  /// sourcePoolAddress_content // assume 1 slot
  /// destTokenAddress_length
  /// destTokenAddress_content // assume 1 slot
  /// extraData_length // contents billed separately
  uint256 public constant MESSAGE_FIXED_BYTES_PER_TOKEN = 32 * ((1 + 3 * 3) + 2);

  /// @dev Any2EVMRampMessage struct has 10 fields, including 3 variable unnested arrays (data, receiver and tokenAmounts).
  /// Each variable array takes 1 more slot to store its length.
  /// When abi encoded, excluding array contents,
  /// Any2EVMMessage takes up a fixed number of 13 slots, 32 bytes each.
  /// For structs that contain arrays, 1 more slot is added to the front, reaching a total of 14.
  /// The fixed bytes does not cover struct data (this is represented by ANY_2_EVM_MESSAGE_FIXED_BYTES_PER_TOKEN)
  uint256 public constant ANY_2_EVM_MESSAGE_FIXED_BYTES = 32 * 14;

  /// @dev Each token transfer adds 1 RampTokenAmount
  /// RampTokenAmount has 4 fields, including 3 bytes.
  /// Each bytes takes 1 more slot to store its length, and one slot to store the offset.
  /// When abi encoded, each token transfer takes up 10 slots, excl bytes contents.
  uint256 public constant ANY_2_EVM_MESSAGE_FIXED_BYTES_PER_TOKEN = 32 * 10;

  bytes32 internal constant EVM_2_EVM_MESSAGE_HASH = keccak256("EVM2EVMMessageHashV2");

  /// @dev Used to hash messages for single-lane ramps.
  /// OnRamp hash(EVM2EVMMessage) = OffRamp hash(EVM2EVMMessage)
  /// The EVM2EVMMessage's messageId is expected to be the output of this hash function
  /// @param original Message to hash
  /// @param metadataHash Immutable metadata hash representing a lane with a fixed OnRamp
  /// @return hashedMessage hashed message as a keccak256
  function _hash(EVM2EVMMessage memory original, bytes32 metadataHash) internal pure returns (bytes32) {
    // Fixed-size message fields are included in nested hash to reduce stack pressure.
    // This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
    return keccak256(
      abi.encode(
        MerkleMultiProof.LEAF_DOMAIN_SEPARATOR,
        metadataHash,
        keccak256(
          abi.encode(
            original.sender,
            original.receiver,
            original.sequenceNumber,
            original.gasLimit,
            original.strict,
            original.nonce,
            original.feeToken,
            original.feeTokenAmount
          )
        ),
        keccak256(original.data),
        keccak256(abi.encode(original.tokenAmounts)),
        keccak256(abi.encode(original.sourceTokenData))
      )
    );
  }

  bytes32 internal constant ANY_2_EVM_MESSAGE_HASH = keccak256("Any2EVMMessageHashV1");
  bytes32 internal constant EVM_2_ANY_MESSAGE_HASH = keccak256("EVM2AnyMessageHashV1");

  /// @dev Used to hash messages for multi-lane family-agnostic OffRamps.
  /// OnRamp hash(EVM2AnyMessage) != Any2EVMRampMessage.messageId
  /// OnRamp hash(EVM2AnyMessage) != OffRamp hash(Any2EVMRampMessage)
  /// @param original OffRamp message to hash
  /// @param onRamp OnRamp to hash the message with - used to compute the metadataHash
  /// @return hashedMessage hashed message as a keccak256
  function _hash(Any2EVMRampMessage memory original, bytes memory onRamp) internal pure returns (bytes32) {
    // Fixed-size message fields are included in nested hash to reduce stack pressure.
    // This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
    return keccak256(
      abi.encode(
        MerkleMultiProof.LEAF_DOMAIN_SEPARATOR,
        // Implicit metadata hash
        keccak256(
          abi.encode(
            ANY_2_EVM_MESSAGE_HASH, original.header.sourceChainSelector, original.header.destChainSelector, onRamp
          )
        ),
        keccak256(
          abi.encode(
            original.header.messageId,
            original.sender,
            original.receiver,
            original.header.sequenceNumber,
            original.gasLimit,
            original.header.nonce
          )
        ),
        keccak256(original.data),
        keccak256(abi.encode(original.tokenAmounts))
      )
    );
  }

  function _hash(EVM2AnyRampMessage memory original, bytes32 metadataHash) internal pure returns (bytes32) {
    // Fixed-size message fields are included in nested hash to reduce stack pressure.
    // This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
    return keccak256(
      abi.encode(
        MerkleMultiProof.LEAF_DOMAIN_SEPARATOR,
        metadataHash,
        keccak256(
          abi.encode(
            original.sender,
            original.receiver,
            original.header.sequenceNumber,
            original.header.nonce,
            original.feeToken,
            original.feeTokenAmount
          )
        ),
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

  /// @notice This methods provides validation for parsing abi encoded addresses by ensuring the
  /// address is within the EVM address space. If it isn't it will revert with an InvalidEVMAddress error, which
  /// we can catch and handle more gracefully than a revert from abi.decode.
  /// @return The address if it is valid, the function will revert otherwise.
  function _validateEVMAddress(bytes memory encodedAddress) internal pure returns (address) {
    if (encodedAddress.length != 32) revert InvalidEVMAddress(encodedAddress);
    uint256 encodedAddressUint = abi.decode(encodedAddress, (uint256));
    if (encodedAddressUint > type(uint160).max || encodedAddressUint < PRECOMPILE_SPACE) {
      revert InvalidEVMAddress(encodedAddress);
    }
    return address(uint160(encodedAddressUint));
  }

  /// @notice Enum listing the possible message execution states within
  /// the offRamp contract.
  /// UNTOUCHED never executed
  /// IN_PROGRESS currently being executed, used a replay protection
  /// SUCCESS successfully executed. End state
  /// FAILURE unsuccessfully executed, manual execution is now enabled.
  /// @dev RMN depends on this enum, if changing, please notify the RMN maintainers.
  enum MessageExecutionState {
    UNTOUCHED,
    IN_PROGRESS,
    SUCCESS,
    FAILURE
  }

  /// @notice CCIP OCR plugin type, used to separate execution & commit transmissions and configs
  enum OCRPluginType {
    Commit,
    Execution
  }

  /// @notice Family-agnostic token amounts used for both OnRamp & OffRamp messages
  struct RampTokenAmount {
    // The source pool address, abi encoded. This value is trusted as it was obtained through the onRamp. It can be
    // relied upon by the destination pool to validate the source pool.
    bytes sourcePoolAddress;
    // The address of the destination token, abi encoded in the case of EVM chains
    // This value is UNTRUSTED as any pool owner can return whatever value they want.
    bytes destTokenAddress;
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint256 amount; // Amount of tokens.
  }

  /// @notice Family-agnostic header for OnRamp & OffRamp messages.
  /// The messageId is not expected to match hash(message), since it may originate from another ramp family
  struct RampMessageHeader {
    bytes32 messageId; // Unique identifier for the message, generated with the source chain's encoding scheme (i.e. not necessarily abi.encoded)
    uint64 sourceChainSelector; // ──╮ the chain selector of the source chain, note: not chainId
    uint64 destChainSelector; //     | the chain selector of the destination chain, note: not chainId
    uint64 sequenceNumber; //        │ sequence number, not unique across lanes
    uint64 nonce; // ────────────────╯ nonce for this lane for this sender, not unique across senders/lanes
  }

  /// @notice Family-agnostic message routed to an OffRamp
  /// Note: hash(Any2EVMRampMessage) != hash(EVM2AnyRampMessage), hash(Any2EVMRampMessage) != messageId
  /// due to encoding & parameter differences
  struct Any2EVMRampMessage {
    RampMessageHeader header; // Message header
    bytes sender; // sender address on the source chain
    bytes data; // arbitrary data payload supplied by the message sender
    address receiver; // receiver address on the destination chain
    uint256 gasLimit; // user supplied maximum gas amount available for dest chain execution
    RampTokenAmount[] tokenAmounts; // array of tokens and amounts to transfer
  }

  /// @notice Family-agnostic message emitted from the OnRamp
  /// Note: hash(Any2EVMRampMessage) != hash(EVM2AnyRampMessage) due to encoding & parameter differences
  /// messageId = hash(EVM2AnyRampMessage) using the source EVM chain's encoding format
  struct EVM2AnyRampMessage {
    RampMessageHeader header; // Message header
    address sender; // sender address on the source chain
    bytes data; // arbitrary data payload supplied by the message sender
    bytes receiver; // receiver address on the destination chain
    bytes extraArgs; // destination-chain specific extra args, such as the gasLimit for EVM chains
    address feeToken; // fee token
    uint256 feeTokenAmount; // fee token amount
    RampTokenAmount[] tokenAmounts; // array of tokens and amounts to transfer
  }

  // bytes4(keccak256("CCIP ChainFamilySelector EVM"))
  bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
}
