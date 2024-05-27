// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Predeploys
/// @notice Contains constant addresses for contracts that are pre-deployed to the L2 system.
library Predeploys {
  /// @notice Address of the L2ToL1MessagePasser predeploy.
  address internal constant L2_TO_L1_MESSAGE_PASSER = 0x4200000000000000000000000000000000000016;

  /// @notice Address of the L2CrossDomainMessenger predeploy.
  address internal constant L2_CROSS_DOMAIN_MESSENGER = 0x4200000000000000000000000000000000000007;

  /// @notice Address of the L2StandardBridge predeploy.
  address internal constant L2_STANDARD_BRIDGE = 0x4200000000000000000000000000000000000010;

  /// @notice Address of the L2ERC721Bridge predeploy.
  address internal constant L2_ERC721_BRIDGE = 0x4200000000000000000000000000000000000014;

  //// @notice Address of the SequencerFeeWallet predeploy.
  address internal constant SEQUENCER_FEE_WALLET = 0x4200000000000000000000000000000000000011;

  /// @notice Address of the OptimismMintableERC20Factory predeploy.
  address internal constant OPTIMISM_MINTABLE_ERC20_FACTORY = 0x4200000000000000000000000000000000000012;

  /// @notice Address of the OptimismMintableERC721Factory predeploy.
  address internal constant OPTIMISM_MINTABLE_ERC721_FACTORY = 0x4200000000000000000000000000000000000017;

  /// @notice Address of the L1Block predeploy.
  address internal constant L1_BLOCK_ATTRIBUTES = 0x4200000000000000000000000000000000000015;

  /// @notice Address of the GasPriceOracle predeploy. Includes fee information
  ///         and helpers for computing the L1 portion of the transaction fee.
  address internal constant GAS_PRICE_ORACLE = 0x420000000000000000000000000000000000000F;

  /// @custom:legacy
  /// @notice Address of the L1MessageSender predeploy. Deprecated. Use L2CrossDomainMessenger
  ///         or access tx.origin (or msg.sender) in a L1 to L2 transaction instead.
  address internal constant L1_MESSAGE_SENDER = 0x4200000000000000000000000000000000000001;

  /// @custom:legacy
  /// @notice Address of the DeployerWhitelist predeploy. No longer active.
  address internal constant DEPLOYER_WHITELIST = 0x4200000000000000000000000000000000000002;

  /// @notice Address of the canonical WETH9 contract.
  address internal constant WETH9 = 0x4200000000000000000000000000000000000006;

  /// @custom:legacy
  /// @notice Address of the LegacyERC20ETH predeploy. Deprecated. Balances are migrated to the
  ///         state trie as of the Bedrock upgrade. Contract has been locked and write functions
  ///         can no longer be accessed.
  address internal constant LEGACY_ERC20_ETH = 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000;

  /// @custom:legacy
  /// @notice Address of the L1BlockNumber predeploy. Deprecated. Use the L1Block predeploy
  ///         instead, which exposes more information about the L1 state.
  address internal constant L1_BLOCK_NUMBER = 0x4200000000000000000000000000000000000013;

  /// @custom:legacy
  /// @notice Address of the LegacyMessagePasser predeploy. Deprecate. Use the updated
  ///         L2ToL1MessagePasser contract instead.
  address internal constant LEGACY_MESSAGE_PASSER = 0x4200000000000000000000000000000000000000;

  /// @notice Address of the ProxyAdmin predeploy.
  address internal constant PROXY_ADMIN = 0x4200000000000000000000000000000000000018;

  /// @notice Address of the BaseFeeVault predeploy.
  address internal constant BASE_FEE_VAULT = 0x4200000000000000000000000000000000000019;

  /// @notice Address of the L1FeeVault predeploy.
  address internal constant L1_FEE_VAULT = 0x420000000000000000000000000000000000001A;

  /// @notice Address of the GovernanceToken predeploy.
  address internal constant GOVERNANCE_TOKEN = 0x4200000000000000000000000000000000000042;

  /// @notice Address of the SchemaRegistry predeploy.
  address internal constant SCHEMA_REGISTRY = 0x4200000000000000000000000000000000000020;

  /// @notice Address of the EAS predeploy.
  address internal constant EAS = 0x4200000000000000000000000000000000000021;

  /// @notice Address of the MultiCall3 predeploy.
  address internal constant MultiCall3 = 0xcA11bde05977b3631167028862bE2a173976CA11;

  /// @notice Address of the Create2Deployer predeploy.
  address internal constant Create2Deployer = 0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2;

  /// @notice Address of the Safe_v130 predeploy.
  address internal constant Safe_v130 = 0x69f4D1788e39c87893C980c06EdF4b7f686e2938;

  /// @notice Address of the SafeL2_v130 predeploy.
  address internal constant SafeL2_v130 = 0xfb1bffC9d739B8D520DaF37dF666da4C687191EA;

  /// @notice Address of the MultiSendCallOnly_v130 predeploy.
  address internal constant MultiSendCallOnly_v130 = 0xA1dabEF33b3B82c7814B6D82A79e50F4AC44102B;

  /// @notice Address of the SafeSingletonFactory predeploy.
  address internal constant SafeSingletonFactory = 0x914d7Fec6aaC8cd542e72Bca78B30650d45643d7;

  /// @notice Address of the DeterministicDeploymentProxy predeploy.
  address internal constant DeterministicDeploymentProxy = 0x4e59b44847b379578588920cA78FbF26c0B4956C;

  /// @notice Address of the MultiSend_v130 predeploy.
  address internal constant MultiSend_v130 = 0x998739BFdAAdde7C933B942a68053933098f9EDa;

  /// @notice Address of the Permit2 predeploy.
  address internal constant Permit2 = 0x000000000022D473030F116dDEE9F6B43aC78BA3;

  /// @notice Address of the SenderCreator predeploy.
  address internal constant SenderCreator = 0x7fc98430eAEdbb6070B35B39D798725049088348;

  /// @notice Address of the EntryPoint predeploy.
  address internal constant EntryPoint = 0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789;
}
