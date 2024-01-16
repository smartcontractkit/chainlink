// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ILiquidityContainer} from "../../../liquidity-manager/interfaces/ILiquidityContainer.sol";

import {OCR3Base} from "../../ocr/OCR3Base.sol";

interface ILiquidityManager {
  struct SendLiquidityParams {
    uint256 amount;
    uint64 remoteChainSelector;
  }

  struct ReceiveLiquidityParams {
    uint256 amount;
    uint64 remoteChainSelector;
    bytes bridgeData;
  }

  struct LiquidityInstructions {
    SendLiquidityParams[] sendLiquidityParams;
    ReceiveLiquidityParams[] receiveLiquidityParams;
  }

  struct CrossChainLiquidityManagerArgs {
    address remoteLiquidityManager;
    // IBridgeAdapter localBridge;
    // address remoteToken;
    uint64 remoteChainSelector;
    bool enabled;
  }

  function getAllCrossChainLiquidityMangers() external view returns (CrossChainLiquidityManagerArgs[] memory);
}

/// @notice Dummy liquidity manager.
/// @dev this is only ever used in tests, do not deploy this for real liquidity management.
contract DummyLiquidityManager is ILiquidityManager, OCR3Base {
  error ZeroAddress();
  error InvalidRemoteChain(uint64 chainSelector);
  error ZeroChainSelector();
  error InsufficientLiquidity(uint256 requested, uint256 available);

  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount
  );
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed remover, uint256 indexed amount);

  struct CrossChainLiquidityManager {
    address remoteLiquidityManager;
    bool enabled;
  }

  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "DummyLiquidityManager 1.0.0";

  /// @notice The chain selector belonging to the chain this pool is deployed on.
  uint64 internal immutable i_localChainSelector;

  /// @notice Mapping of chain selector to liquidity container on other chains
  mapping(uint64 chainSelector => CrossChainLiquidityManager) private s_crossChainLiquidityManager;

  uint64[] private s_supportedDestChains;

  constructor(uint64 localChainSelector) OCR3Base() {
    if (localChainSelector == 0) {
      revert ZeroChainSelector();
    }

    i_localChainSelector = localChainSelector;
  }

  function _report(bytes calldata report, uint64 ocrSeqNum) internal override {
    // do nothing, dummy
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  function getSupportedDestChains() external view returns (uint64[] memory) {
    return s_supportedDestChains;
  }

  /// @notice Gets the cross chain liquidity manager
  function getCrossChainLiquidityManager(
    uint64 chainSelector
  ) external view returns (CrossChainLiquidityManager memory) {
    return s_crossChainLiquidityManager[chainSelector];
  }

  /// @notice Gets all cross chain liquidity managers
  /// @dev We don't care too much about gas since this function is intended for offchain usage.
  function getAllCrossChainLiquidityMangers() external view returns (CrossChainLiquidityManagerArgs[] memory) {
    CrossChainLiquidityManagerArgs[] memory managers = new CrossChainLiquidityManagerArgs[](
      s_supportedDestChains.length
    );
    for (uint256 i = 0; i < s_supportedDestChains.length; ++i) {
      uint64 chainSelector = s_supportedDestChains[i];
      CrossChainLiquidityManager memory currentManager = s_crossChainLiquidityManager[chainSelector];
      managers[i] = CrossChainLiquidityManagerArgs({
        remoteLiquidityManager: currentManager.remoteLiquidityManager,
        remoteChainSelector: chainSelector,
        enabled: currentManager.enabled
      });
    }

    return managers;
  }

  /// @notice Sets a list of cross chain liquidity managers.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function setCrossChainLiquidityManager(
    CrossChainLiquidityManagerArgs[] calldata crossChainLiquidityManagers
  ) external onlyOwner {
    for (uint256 i = 0; i < crossChainLiquidityManagers.length; ++i) {
      setCrossChainLiquidityManager(crossChainLiquidityManagers[i]);
    }
  }

  /// @notice Sets a single cross chain liquidity manager.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function setCrossChainLiquidityManager(
    CrossChainLiquidityManagerArgs calldata crossChainLiqManager
  ) public onlyOwner {
    if (crossChainLiqManager.remoteChainSelector == 0) {
      revert ZeroChainSelector();
    }

    if (crossChainLiqManager.remoteLiquidityManager == address(0)) {
      revert ZeroAddress();
    }

    // If the destination chain is new, add it to the list of supported chains
    if (s_crossChainLiquidityManager[crossChainLiqManager.remoteChainSelector].remoteLiquidityManager == address(0)) {
      s_supportedDestChains.push(crossChainLiqManager.remoteChainSelector);
    }

    s_crossChainLiquidityManager[crossChainLiqManager.remoteChainSelector] = CrossChainLiquidityManager({
      remoteLiquidityManager: crossChainLiqManager.remoteLiquidityManager,
      enabled: crossChainLiqManager.enabled
    });
  }
}
