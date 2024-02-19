// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {ITokenMessenger} from "./ITokenMessenger.sol";
import {IMessageTransmitter} from "./IMessageTransmitter.sol";

import {TokenPool} from "../TokenPool.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";

/// @notice This pool mints and burns USDC tokens through the Cross Chain Transfer
/// Protocol (CCTP).
contract USDCTokenPool is TokenPool, ITypeAndVersion {
  using SafeERC20 for IERC20;

  event DomainsSet(DomainUpdate[]);
  event ConfigSet(address tokenMessenger);

  error UnknownDomain(uint64 domain);
  error UnlockingUSDCFailed();
  error InvalidConfig();
  error InvalidDomain(DomainUpdate domain);
  error InvalidMessageVersion(uint32 version);
  error InvalidTokenMessengerVersion(uint32 version);
  error InvalidNonce(uint64 expected, uint64 got);
  error InvalidSourceDomain(uint32 expected, uint32 got);
  error InvalidDestinationDomain(uint32 expected, uint32 got);

  // This data is supplied from offchain and contains everything needed
  // to receive the USDC tokens.
  struct MessageAndAttestation {
    bytes message;
    bytes attestation;
  }

  // A domain is a USDC representation of a chain.
  struct DomainUpdate {
    bytes32 allowedCaller; //       Address allowed to mint on the domain
    uint32 domainIdentifier; // ──╮ Unique domain ID
    uint64 destChainSelector; //  │ The destination chain for this domain
    bool enabled; // ─────────────╯ Whether the domain is enabled
  }

  struct SourceTokenDataPayload {
    uint64 nonce;
    uint32 sourceDomain;
  }

  string public constant override typeAndVersion = "USDCTokenPool 1.4.0";

  // We restrict to the first version. New pool may be required for subsequent versions.
  uint32 public constant SUPPORTED_USDC_VERSION = 0;

  // The local USDC config
  ITokenMessenger public immutable i_tokenMessenger;
  IMessageTransmitter public immutable i_messageTransmitter;
  uint32 public immutable i_localDomainIdentifier;

  // The unique USDC pool flag to signal through EIP 165 that this is a USDC token pool.
  bytes4 private constant USDC_INTERFACE_ID = bytes4(keccak256("USDC"));

  /// A domain is a USDC representation of a destination chain.
  /// @dev Zero is a valid domain identifier.
  /// @dev The address to mint on the destination chain is the corresponding USDC pool.
  struct Domain {
    bytes32 allowedCaller; //      Address allowed to mint on the domain
    uint32 domainIdentifier; // ─╮ Unique domain ID
    bool enabled; // ────────────╯ Whether the domain is enabled
  }

  // A mapping of CCIP chain identifiers to destination domains
  mapping(uint64 chainSelector => Domain CCTPDomain) private s_chainToDomain;

  constructor(
    ITokenMessenger tokenMessenger,
    IERC20 token,
    address[] memory allowlist,
    address armProxy,
    address router
  ) TokenPool(token, allowlist, armProxy, router) {
    if (address(tokenMessenger) == address(0)) revert InvalidConfig();
    IMessageTransmitter transmitter = IMessageTransmitter(tokenMessenger.localMessageTransmitter());
    uint32 transmitterVersion = transmitter.version();
    if (transmitterVersion != SUPPORTED_USDC_VERSION) revert InvalidMessageVersion(transmitterVersion);
    uint32 tokenMessengerVersion = tokenMessenger.messageBodyVersion();
    if (tokenMessengerVersion != SUPPORTED_USDC_VERSION) revert InvalidTokenMessengerVersion(tokenMessengerVersion);

    i_tokenMessenger = tokenMessenger;
    i_messageTransmitter = transmitter;
    i_localDomainIdentifier = transmitter.localDomain();
    i_token.safeIncreaseAllowance(address(i_tokenMessenger), type(uint256).max);
    emit ConfigSet(address(tokenMessenger));
  }

  /// @notice returns the USDC interface flag used for EIP165 identification.
  function getUSDCInterfaceId() public pure returns (bytes4) {
    return USDC_INTERFACE_ID;
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == USDC_INTERFACE_ID || super.supportsInterface(interfaceId);
  }

  /// @notice Burn the token in the pool
  /// @dev Burn is not rate limited at per-pool level. Burn does not contribute to honey pot risk.
  /// Benefits of rate limiting here does not justify the extra gas cost.
  /// @param amount Amount to burn
  /// @dev emits ITokenMessenger.DepositForBurn
  /// @dev Assumes caller has validated destinationReceiver
  function lockOrBurn(
    address originalSender,
    bytes calldata destinationReceiver,
    uint256 amount,
    uint64 remoteChainSelector,
    bytes calldata
  ) external override onlyOnRamp(remoteChainSelector) checkAllowList(originalSender) returns (bytes memory) {
    Domain memory domain = s_chainToDomain[remoteChainSelector];
    if (!domain.enabled) revert UnknownDomain(remoteChainSelector);
    _consumeOutboundRateLimit(remoteChainSelector, amount);
    bytes32 receiver = bytes32(destinationReceiver[0:32]);
    // Since this pool is the msg sender of the CCTP transaction, only this contract
    // is able to call replaceDepositForBurn. Since this contract does not implement
    // replaceDepositForBurn, the tokens cannot be maliciously re-routed to another address.
    uint64 nonce = i_tokenMessenger.depositForBurnWithCaller(
      amount,
      domain.domainIdentifier,
      receiver,
      address(i_token),
      domain.allowedCaller
    );
    emit Burned(msg.sender, amount);
    return abi.encode(SourceTokenDataPayload({nonce: nonce, sourceDomain: i_localDomainIdentifier}));
  }

  /// @notice Mint tokens from the pool to the recipient
  /// @param receiver Recipient address
  /// @param amount Amount to mint
  /// @param extraData Encoded return data from `lockOrBurn` and offchain attestation data
  /// @dev sourceTokenData is part of the verified message and passed directly from
  /// the offramp so it is guaranteed to be what the lockOrBurn pool released on the
  /// source chain. It contains (nonce, sourceDomain) which is guaranteed by CCTP
  /// to be unique.
  /// offchainTokenData is untrusted (can be supplied by manual execution), but we assert
  /// that (nonce, sourceDomain) is equal to the message's (nonce, sourceDomain) and
  /// receiveMessage will assert that Attestation contains a valid attestation signature
  /// for that message, including its (nonce, sourceDomain). This way, the only
  /// non-reverting offchainTokenData that can be supplied is a valid attestation for the
  /// specific message that was sent on source.
  function releaseOrMint(
    bytes memory,
    address receiver,
    uint256 amount,
    uint64 remoteChainSelector,
    bytes memory extraData
  ) external override onlyOffRamp(remoteChainSelector) {
    _consumeInboundRateLimit(remoteChainSelector, amount);
    (bytes memory sourceData, bytes memory offchainTokenData) = abi.decode(extraData, (bytes, bytes));
    SourceTokenDataPayload memory sourceTokenData = abi.decode(sourceData, (SourceTokenDataPayload));
    MessageAndAttestation memory msgAndAttestation = abi.decode(offchainTokenData, (MessageAndAttestation));

    _validateMessage(msgAndAttestation.message, sourceTokenData);

    if (!i_messageTransmitter.receiveMessage(msgAndAttestation.message, msgAndAttestation.attestation))
      revert UnlockingUSDCFailed();
    emit Minted(msg.sender, receiver, amount);
  }

  /// @notice Validates the USDC encoded message against the given parameters.
  /// @param usdcMessage The USDC encoded message
  /// @param sourceTokenData The expected source chain token data to check against
  /// @dev Only supports version SUPPORTED_USDC_VERSION of the CCTP message format
  /// @dev Message format for USDC:
  ///     * Field                 Bytes      Type       Index
  ///     * version               4          uint32     0
  ///     * sourceDomain          4          uint32     4
  ///     * destinationDomain     4          uint32     8
  ///     * nonce                 8          uint64     12
  ///     * sender                32         bytes32    20
  ///     * recipient             32         bytes32    52
  ///     * destinationCaller     32         bytes32    84
  ///     * messageBody           dynamic    bytes      116
  function _validateMessage(bytes memory usdcMessage, SourceTokenDataPayload memory sourceTokenData) internal view {
    uint32 version;
    // solhint-disable-next-line no-inline-assembly
    assembly {
      // We truncate using the datatype of the version variable, meaning
      // we will only be left with the first 4 bytes of the message.
      version := mload(add(usdcMessage, 4)) // 0 + 4 = 4
    }
    // This token pool only supports version 0 of the CCTP message format
    // We check the version prior to loading the rest of the message
    // to avoid unexpected reverts due to out-of-bounds reads.
    if (version != SUPPORTED_USDC_VERSION) revert InvalidMessageVersion(version);

    uint32 sourceDomain;
    uint32 destinationDomain;
    uint64 nonce;

    // solhint-disable-next-line no-inline-assembly
    assembly {
      sourceDomain := mload(add(usdcMessage, 8)) // 4 + 4 = 8
      destinationDomain := mload(add(usdcMessage, 12)) // 8 + 4 = 12
      nonce := mload(add(usdcMessage, 20)) // 12 + 8 = 20
    }

    if (sourceDomain != sourceTokenData.sourceDomain)
      revert InvalidSourceDomain(sourceTokenData.sourceDomain, sourceDomain);
    if (destinationDomain != i_localDomainIdentifier)
      revert InvalidDestinationDomain(i_localDomainIdentifier, destinationDomain);
    if (nonce != sourceTokenData.nonce) revert InvalidNonce(sourceTokenData.nonce, nonce);
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Gets the CCTP domain for a given CCIP chain selector.
  function getDomain(uint64 chainSelector) external view returns (Domain memory) {
    return s_chainToDomain[chainSelector];
  }

  /// @notice Sets the CCTP domain for a CCIP chain selector.
  /// @dev Must verify mapping of selectors -> (domain, caller) offchain.
  function setDomains(DomainUpdate[] calldata domains) external onlyOwner {
    for (uint256 i = 0; i < domains.length; ++i) {
      DomainUpdate memory domain = domains[i];
      if (domain.allowedCaller == bytes32(0) || domain.destChainSelector == 0) revert InvalidDomain(domain);

      s_chainToDomain[domain.destChainSelector] = Domain({
        domainIdentifier: domain.domainIdentifier,
        allowedCaller: domain.allowedCaller,
        enabled: domain.enabled
      });
    }
    emit DomainsSet(domains);
  }
}
