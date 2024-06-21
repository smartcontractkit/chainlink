// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/libraries/Predeploys.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

/// @title Predeploys
/// @notice Contains constant addresses for protocol contracts that are pre-deployed to the L2 system.
//          This excludes the preinstalls (non-protocol contracts).
library Predeploys {
  /// @notice Number of predeploy-namespace addresses reserved for protocol usage.
  uint256 internal constant PREDEPLOY_COUNT = 2048;

  /// @custom:legacy
  /// @notice Address of the LegacyMessagePasser predeploy. Deprecate. Use the updated
  ///         L2ToL1MessagePasser contract instead.
  address internal constant LEGACY_MESSAGE_PASSER = 0x4200000000000000000000000000000000000000;

  /// @custom:legacy
  /// @notice Address of the L1MessageSender predeploy. Deprecated. Use L2CrossDomainMessenger
  ///         or access tx.origin (or msg.sender) in a L1 to L2 transaction instead.
  ///         Not embedded into new OP-Stack chains.
  address internal constant L1_MESSAGE_SENDER = 0x4200000000000000000000000000000000000001;

  /// @custom:legacy
  /// @notice Address of the DeployerWhitelist predeploy. No longer active.
  address internal constant DEPLOYER_WHITELIST = 0x4200000000000000000000000000000000000002;

  /// @notice Address of the canonical WETH contract.
  address internal constant WETH = 0x4200000000000000000000000000000000000006;

  /// @notice Address of the L2CrossDomainMessenger predeploy.
  address internal constant L2_CROSS_DOMAIN_MESSENGER = 0x4200000000000000000000000000000000000007;

  /// @notice Address of the GasPriceOracle predeploy. Includes fee information
  ///         and helpers for computing the L1 portion of the transaction fee.
  address internal constant GAS_PRICE_ORACLE = 0x420000000000000000000000000000000000000F;

  /// @notice Address of the L2StandardBridge predeploy.
  address internal constant L2_STANDARD_BRIDGE = 0x4200000000000000000000000000000000000010;

  //// @notice Address of the SequencerFeeWallet predeploy.
  address internal constant SEQUENCER_FEE_WALLET = 0x4200000000000000000000000000000000000011;

  /// @notice Address of the OptimismMintableERC20Factory predeploy.
  address internal constant OPTIMISM_MINTABLE_ERC20_FACTORY = 0x4200000000000000000000000000000000000012;

  /// @custom:legacy
  /// @notice Address of the L1BlockNumber predeploy. Deprecated. Use the L1Block predeploy
  ///         instead, which exposes more information about the L1 state.
  address internal constant L1_BLOCK_NUMBER = 0x4200000000000000000000000000000000000013;

  /// @notice Address of the L2ERC721Bridge predeploy.
  address internal constant L2_ERC721_BRIDGE = 0x4200000000000000000000000000000000000014;

  /// @notice Address of the L1Block predeploy.
  address internal constant L1_BLOCK_ATTRIBUTES = 0x4200000000000000000000000000000000000015;

  /// @notice Address of the L2ToL1MessagePasser predeploy.
  address internal constant L2_TO_L1_MESSAGE_PASSER = 0x4200000000000000000000000000000000000016;

  /// @notice Address of the OptimismMintableERC721Factory predeploy.
  address internal constant OPTIMISM_MINTABLE_ERC721_FACTORY = 0x4200000000000000000000000000000000000017;

  /// @notice Address of the ProxyAdmin predeploy.
  address internal constant PROXY_ADMIN = 0x4200000000000000000000000000000000000018;

  /// @notice Address of the BaseFeeVault predeploy.
  address internal constant BASE_FEE_VAULT = 0x4200000000000000000000000000000000000019;

  /// @notice Address of the L1FeeVault predeploy.
  address internal constant L1_FEE_VAULT = 0x420000000000000000000000000000000000001A;

  /// @notice Address of the SchemaRegistry predeploy.
  address internal constant SCHEMA_REGISTRY = 0x4200000000000000000000000000000000000020;

  /// @notice Address of the EAS predeploy.
  address internal constant EAS = 0x4200000000000000000000000000000000000021;

  /// @notice Address of the GovernanceToken predeploy.
  address internal constant GOVERNANCE_TOKEN = 0x4200000000000000000000000000000000000042;

  /// @custom:legacy
  /// @notice Address of the LegacyERC20ETH predeploy. Deprecated. Balances are migrated to the
  ///         state trie as of the Bedrock upgrade. Contract has been locked and write functions
  ///         can no longer be accessed.
  address internal constant LEGACY_ERC20_ETH = 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000;

  /// @notice Address of the CrossL2Inbox predeploy.
  address internal constant CROSS_L2_INBOX = 0x4200000000000000000000000000000000000022;

  /// @notice Address of the L2ToL2CrossDomainMessenger predeploy.
  address internal constant L2_TO_L2_CROSS_DOMAIN_MESSENGER = 0x4200000000000000000000000000000000000023;

  /// @notice Returns the name of the predeploy at the given address.
  function getName(address _addr) internal pure returns (string memory out_) {
    require(isPredeployNamespace(_addr), "Predeploys: address must be a predeploy");
    if (_addr == LEGACY_MESSAGE_PASSER) return "LegacyMessagePasser";
    if (_addr == L1_MESSAGE_SENDER) return "L1MessageSender";
    if (_addr == DEPLOYER_WHITELIST) return "DeployerWhitelist";
    if (_addr == WETH) return "WETH";
    if (_addr == L2_CROSS_DOMAIN_MESSENGER) return "L2CrossDomainMessenger";
    if (_addr == GAS_PRICE_ORACLE) return "GasPriceOracle";
    if (_addr == L2_STANDARD_BRIDGE) return "L2StandardBridge";
    if (_addr == SEQUENCER_FEE_WALLET) return "SequencerFeeVault";
    if (_addr == OPTIMISM_MINTABLE_ERC20_FACTORY) return "OptimismMintableERC20Factory";
    if (_addr == L1_BLOCK_NUMBER) return "L1BlockNumber";
    if (_addr == L2_ERC721_BRIDGE) return "L2ERC721Bridge";
    if (_addr == L1_BLOCK_ATTRIBUTES) return "L1Block";
    if (_addr == L2_TO_L1_MESSAGE_PASSER) return "L2ToL1MessagePasser";
    if (_addr == OPTIMISM_MINTABLE_ERC721_FACTORY) return "OptimismMintableERC721Factory";
    if (_addr == PROXY_ADMIN) return "ProxyAdmin";
    if (_addr == BASE_FEE_VAULT) return "BaseFeeVault";
    if (_addr == L1_FEE_VAULT) return "L1FeeVault";
    if (_addr == SCHEMA_REGISTRY) return "SchemaRegistry";
    if (_addr == EAS) return "EAS";
    if (_addr == GOVERNANCE_TOKEN) return "GovernanceToken";
    if (_addr == LEGACY_ERC20_ETH) return "LegacyERC20ETH";
    if (_addr == CROSS_L2_INBOX) return "CrossL2Inbox";
    if (_addr == L2_TO_L2_CROSS_DOMAIN_MESSENGER) return "L2ToL2CrossDomainMessenger";
    revert("Predeploys: unnamed predeploy");
  }

  /// @notice Returns true if the predeploy is not proxied.
  function notProxied(address _addr) internal pure returns (bool) {
    return _addr == GOVERNANCE_TOKEN || _addr == WETH;
  }

  /// @notice Returns true if the address is a defined predeploy that is embedded into new OP-Stack chains.
  function isSupportedPredeploy(address _addr, bool _useInterop) internal pure returns (bool) {
    return
      _addr == LEGACY_MESSAGE_PASSER ||
      _addr == DEPLOYER_WHITELIST ||
      _addr == WETH ||
      _addr == L2_CROSS_DOMAIN_MESSENGER ||
      _addr == GAS_PRICE_ORACLE ||
      _addr == L2_STANDARD_BRIDGE ||
      _addr == SEQUENCER_FEE_WALLET ||
      _addr == OPTIMISM_MINTABLE_ERC20_FACTORY ||
      _addr == L1_BLOCK_NUMBER ||
      _addr == L2_ERC721_BRIDGE ||
      _addr == L1_BLOCK_ATTRIBUTES ||
      _addr == L2_TO_L1_MESSAGE_PASSER ||
      _addr == OPTIMISM_MINTABLE_ERC721_FACTORY ||
      _addr == PROXY_ADMIN ||
      _addr == BASE_FEE_VAULT ||
      _addr == L1_FEE_VAULT ||
      _addr == SCHEMA_REGISTRY ||
      _addr == EAS ||
      _addr == GOVERNANCE_TOKEN ||
      (_useInterop && _addr == CROSS_L2_INBOX) ||
      (_useInterop && _addr == L2_TO_L2_CROSS_DOMAIN_MESSENGER);
  }

  function isPredeployNamespace(address _addr) internal pure returns (bool) {
    return uint160(_addr) >> 11 == uint160(0x4200000000000000000000000000000000000000) >> 11;
  }

  /// @notice Function to compute the expected address of the predeploy implementation
  ///         in the genesis state.
  function predeployToCodeNamespace(address _addr) internal pure returns (address) {
    require(isPredeployNamespace(_addr), "Predeploys: can only derive code-namespace address for predeploy addresses");
    return
      address(
        uint160((uint256(uint160(_addr)) & 0xffff) | uint256(uint160(0xc0D3C0d3C0d3C0D3c0d3C0d3c0D3C0d3c0d30000)))
      );
  }
}
