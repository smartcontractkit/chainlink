// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./KeeperRegistryBase.sol";

/**
 * @notice Logic contract, works in tandem with KeeperRegistry as a proxy.
 * It's own state is never used.
 */
contract KeeperRegistryLogic is KeeperRegistryBase {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  /**
   * @param link address of the LINK Token
   * @param linkEthFeed address of the LINK/ETH price feed
   * @param fastGasFeed address of the Fast Gas price feed
   */
  constructor(
    address link,
    address linkEthFeed,
    address fastGasFeed
  ) KeeperRegistryBase(link, linkEthFeed, fastGasFeed) {}

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function checkUpkeep(uint256 id, address from)
    external
    cannotExecute
    returns (
      bytes memory performData,
      uint256 maxLinkPayment,
      uint256 gasLimit,
      uint256 adjustedGasWei,
      uint256 linkEth
    )
  {
    Upkeep memory upkeep = s_upkeep[id];

    bytes memory callData = abi.encodeWithSelector(CHECK_SELECTOR, s_checkData[id]);
    (bool success, bytes memory result) = upkeep.target.call{gas: s_storage.checkGasLimit}(callData);

    if (!success) revert TargetCheckReverted(result);

    (success, performData) = abi.decode(result, (bool, bytes));
    if (!success) revert UpkeepNotNeeded();

    PerformParams memory params = _generatePerformParams(from, id, performData, false);
    _prePerformUpkeep(upkeep, params.from, params.maxLinkPayment);

    return (performData, params.maxLinkPayment, params.gasLimit, params.adjustedGasWei, params.linkEth);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getUpkeep(uint256 id)
    external
    view
    returns (
      address target,
      uint32 executeGas,
      bytes memory checkData,
      uint96 balance,
      address lastKeeper,
      address admin,
      uint64 maxValidBlocknumber,
      uint96 amountSpent
    )
  {
    Upkeep memory reg = s_upkeep[id];
    return (
      reg.target,
      reg.executeGas,
      s_checkData[id],
      reg.balance,
      reg.lastKeeper,
      reg.admin,
      reg.maxValidBlocknumber,
      reg.amountSpent
    );
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) external view returns (uint256[] memory) {
    uint256 maxIdx = s_upkeepIDs.length();
    if (startIndex >= maxIdx) revert IndexOutOfRange();
    if (maxCount == 0) {
      maxCount = maxIdx - startIndex;
    }
    uint256[] memory ids = new uint256[](maxCount);
    for (uint256 idx = 0; idx < maxCount; idx++) {
      ids[idx] = s_upkeepIDs.at(startIndex + idx);
    }
    return ids;
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getKeeperInfo(address query)
    external
    view
    returns (
      address payee,
      bool active,
      uint96 balance
    )
  {
    KeeperInfo memory keeper = s_keeperInfo[query];
    return (keeper.payee, keeper.active, keeper.balance);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getState()
    external
    view
    returns (
      State memory state,
      Config memory config,
      address[] memory keepers
    )
  {
    Storage memory store = s_storage;
    state.nonce = store.nonce;
    state.ownerLinkBalance = s_ownerLinkBalance;
    state.expectedLinkBalance = s_expectedLinkBalance;
    state.numUpkeeps = s_upkeepIDs.length();
    config.paymentPremiumPPB = store.paymentPremiumPPB;
    config.flatFeeMicroLink = store.flatFeeMicroLink;
    config.blockCountPerTurn = store.blockCountPerTurn;
    config.checkGasLimit = store.checkGasLimit;
    config.stalenessSeconds = store.stalenessSeconds;
    config.gasCeilingMultiplier = store.gasCeilingMultiplier;
    config.minUpkeepSpend = store.minUpkeepSpend;
    config.maxPerformGas = store.maxPerformGas;
    config.fallbackGasPrice = s_fallbackGasPrice;
    config.fallbackLinkPrice = s_fallbackLinkPrice;
    config.transcoder = s_transcoder;
    config.registrar = s_registrar;
    return (state, config, s_keeperList);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getMinBalanceForUpkeep(uint256 id) external view returns (uint96 minBalance) {
    return getMaxPaymentForGas(s_upkeep[id].executeGas);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getMaxPaymentForGas(uint256 gasLimit) public view returns (uint96 maxPayment) {
    (uint256 gasWei, uint256 linkEth) = _getFeedData();
    uint256 adjustedGasWei = _adjustGasPrice(gasWei, false);
    return _calculatePaymentAmount(gasLimit, adjustedGasWei, linkEth);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function getPeerRegistryMigrationPermission(address peer) external view returns (MigrationPermission) {
    return s_peerRegistryMigrationPermission[peer];
  }
}
