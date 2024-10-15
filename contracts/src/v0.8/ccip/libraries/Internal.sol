// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {MerkleMultiProof} from "../libraries/MerkleMultiProof.sol";

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
  /// @dev The expected number of bytes returned by the balanceOf function.
  uint256 internal constant MAX_BALANCE_OF_RET_BYTES = 32;

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
    uint224 value; // ──────╮ Value in uint224, packed.
    uint32 timestamp; // ───╯ Timestamp of the most recent price update.
  }

  /// @dev Gas price is stored in 112-bit unsigned int. uint224 can pack 2 prices.
  /// When packing L1 and L2 gas prices, L1 gas price is left-shifted to the higher-order bits.
  /// Using uint8 type, which cannot be higher than other bit shift operands, to avoid shift operand type warning.
  uint8 public constant GAS_PRICE_BITS = 112;

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
    uint32 destGasAmount; // The amount of gas available for the releaseOrMint and balanceOf calls on the offRamp
  }

  /// @notice Report that is submitted by the execution DON at the execution phase. (including chain selector data)
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct ExecutionReport {
    uint64 sourceChainSelector; // Source chain selector for which the report is submitted
    Any2EVMRampMessage[] messages;
    // Contains a bytes array for each message, each inner bytes array contains bytes per transferred token
    bytes[][] offchainTokenData;
    bytes32[] proofs;
    uint256 proofFlagBits;
  }

  /// @dev Any2EVMRampMessage struct has 10 fields, including 3 variable unnested arrays (data, receiver and tokenAmounts).
  /// Each variable array takes 1 more slot to store its length.
  /// When abi encoded, excluding array contents,
  /// Any2EVMMessage takes up a fixed number of 13 slots, 32 bytes each.
  /// For structs that contain arrays, 1 more slot is added to the front, reaching a total of 14.
  /// The fixed bytes does not cover struct data (this is represented by MESSAGE_FIXED_BYTES_PER_TOKEN)
  uint256 public constant MESSAGE_FIXED_BYTES = 32 * 14;

  /// @dev Each token transfer adds 1 RampTokenAmount
  /// RampTokenAmount has 5 fields, 2 of which are bytes type, 1 Address, 1 uint256 and 1 uint32.
  /// Each bytes type takes 1 slot for length, 1 slot for data and 1 slot for the offset.
  /// address
  /// uint256 amount takes 1 slot.
  /// uint32 destGasAmount takes 1 slot.
  uint256 public constant MESSAGE_FIXED_BYTES_PER_TOKEN = 32 * ((2 * 3) + 3);

  bytes32 internal constant ANY_2_EVM_MESSAGE_HASH = keccak256("Any2EVMMessageHashV1");
  bytes32 internal constant EVM_2_ANY_MESSAGE_HASH = keccak256("EVM2AnyMessageHashV1");

  /// @dev Used to hash messages for multi-lane family-agnostic OffRamps.
  /// OnRamp hash(EVM2AnyMessage) != Any2EVMRampMessage.messageId
  /// OnRamp hash(EVM2AnyMessage) != OffRamp hash(Any2EVMRampMessage)
  /// @param original OffRamp message to hash
  /// @param metadataHash Hash preimage to ensure global uniqueness
  /// @return hashedMessage hashed message as a keccak256
  function _hash(Any2EVMRampMessage memory original, bytes32 metadataHash) internal pure returns (bytes32) {
    // Fixed-size message fields are included in nested hash to reduce stack pressure.
    // This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
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
    // This hashing scheme is also used by RMN. If changing it, please notify the RMN maintainers.
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

  /// @notice This methods provides validation for parsing abi encoded addresses by ensuring the
  /// address is within the EVM address space. If it isn't it will revert with an InvalidEVMAddress error, which
  /// we can catch and handle more gracefully than a revert from abi.decode.
  /// @return The address if it is valid, the function will revert otherwise.
  function _validateEVMAddress(
    bytes memory encodedAddress
  ) internal pure returns (address) {
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

  /// @notice Family-agnostic header for OnRamp & OffRamp messages.
  /// The messageId is not expected to match hash(message), since it may originate from another ramp family
  struct RampMessageHeader {
    bytes32 messageId; // Unique identifier for the message, generated with the source chain's encoding scheme (i.e. not necessarily abi.encoded)
    uint64 sourceChainSelector; // ──╮ the chain selector of the source chain, note: not chainId
    uint64 destChainSelector; //     │ the chain selector of the destination chain, note: not chainId
    uint64 sequenceNumber; //        │ sequence number, not unique across lanes
    uint64 nonce; // ────────────────╯ nonce for this lane for this sender, not unique across senders/lanes
  }

  struct EVM2AnyTokenTransfer {
    // The source pool EVM address. This value is trusted as it was obtained through the onRamp. It can be
    // relied upon by the destination pool to validate the source pool.
    address sourcePoolAddress;
    // The EVM address of the destination token
    // This value is UNTRUSTED as any pool owner can return whatever value they want.
    bytes destTokenAddress;
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint256 amount; // Amount of tokens.
    // Destination chain specific execution data encoded in bytes
    // for an EVM destination, it consists of the amount of gas available for the releaseOrMint
    // and transfer calls made by the offRamp
    bytes destExecData;
  }

  struct Any2EVMTokenTransfer {
    // The source pool EVM address encoded to bytes. This value is trusted as it is obtained through the onRamp. It can be
    // relied upon by the destination pool to validate the source pool.
    bytes sourcePoolAddress;
    address destTokenAddress; // ───╮ Address of destination token
    uint32 destGasAmount; //────────╯ The amount of gas available for the releaseOrMint and transfer calls on the offRamp.
    // Optional pool data to be transferred to the destination chain. Be default this is capped at
    // CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
    // has to be set for the specific token.
    bytes extraData;
    uint256 amount; // Amount of tokens.
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
    Any2EVMTokenTransfer[] tokenAmounts; // array of tokens and amounts to transfer
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
    uint256 feeValueJuels; // fee amount in Juels
    EVM2AnyTokenTransfer[] tokenAmounts; // array of tokens and amounts to transfer
  }

  // bytes4(keccak256("CCIP ChainFamilySelector EVM"))
  bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;

  /// @dev Struct to hold a merkle root and an interval for a source chain so that an array of these can be passed in the CommitReport.
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  /// @dev inefficient struct packing intentionally chosen to maintain order of specificity. Not a storage struct so impact is minimal.
  // solhint-disable-next-line gas-struct-packing
  struct MerkleRoot {
    uint64 sourceChainSelector; //  Remote source chain selector that the Merkle Root is scoped to
    bytes onRampAddress; //         Generic onramp address, to support arbitrary sources; for EVM, use abi.encode
    uint64 minSeqNr; // ──────────╮ Minimum sequence number, inclusive
    uint64 maxSeqNr; // ──────────╯ Maximum sequence number, inclusive
    bytes32 merkleRoot; //          Merkle root covering the interval & source chain messages
  }
}
