// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./KeeperRegistryBase1_3.sol";
import "./KeeperRegistryLogic1_3.sol";
import {AutomationRegistryExecutableInterface} from "../../interfaces/automation/1_3/AutomationRegistryInterface1_3.sol";
import "../../interfaces/automation/MigratableKeeperRegistryInterface.sol";
import "../../interfaces/TypeAndVersionInterface.sol";
import "../../interfaces/ERC677ReceiverInterface.sol";

/**
 * @notice Registry for adding work for Chainlink Keepers to perform on client
 * contracts. Clients must support the Upkeep interface.
 */
contract KeeperRegistry1_3 is
  KeeperRegistryBase1_3,
  Proxy,
  TypeAndVersionInterface,
  AutomationRegistryExecutableInterface,
  MigratableKeeperRegistryInterface,
  ERC677ReceiverInterface
{
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  address public immutable KEEPER_REGISTRY_LOGIC;

  /**
   * @notice versions:
   * - KeeperRegistry 1.3.0: split contract into Proxy and Logic
   *                       : account for Arbitrum and Optimism L1 gas fee
   *                       : allow users to configure upkeeps
   * - KeeperRegistry 1.2.0: allow funding within performUpkeep
   *                       : allow configurable registry maxPerformGas
   *                       : add function to let admin change upkeep gas limit
   *                       : add minUpkeepSpend requirement
                           : upgrade to solidity v0.8
   * - KeeperRegistry 1.1.0: added flatFeeMicroLink
   * - KeeperRegistry 1.0.0: initial release
   */
  string public constant override typeAndVersion = "KeeperRegistry 1.3.0";

  /**
   * @param keeperRegistryLogic the address of keeper registry logic
   * @param config registry config settings
   */
  constructor(KeeperRegistryLogic1_3 keeperRegistryLogic, Config memory config)
    KeeperRegistryBase1_3(
      keeperRegistryLogic.PAYMENT_MODEL(),
      keeperRegistryLogic.REGISTRY_GAS_OVERHEAD(),
      address(keeperRegistryLogic.LINK()),
      address(keeperRegistryLogic.LINK_ETH_FEED()),
      address(keeperRegistryLogic.FAST_GAS_FEED())
    )
  {
    KEEPER_REGISTRY_LOGIC = address(keeperRegistryLogic);
    setConfig(config);
  }

  // ACTIONS

  /**
   * @notice adds a new upkeep
   * @param target address to perform upkeep on
   * @param gasLimit amount of gas to provide the target contract when
   * performing upkeep
   * @param admin address to cancel upkeep and withdraw remaining funds
   * @param checkData data passed to the contract when checking for upkeep
   */
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    bytes calldata checkData
  ) external override returns (uint256 id) {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice simulated by keepers via eth_call to see if the upkeep needs to be
   * performed. If upkeep is needed, the call then simulates performUpkeep
   * to make sure it succeeds. Finally, it returns the success status along with
   * payment information and the perform data payload.
   * @param id identifier of the upkeep to check
   * @param from the address to simulate performing the upkeep from
   */
  function checkUpkeep(uint256 id, address from)
    external
    override
    cannotExecute
    returns (
      bytes memory performData,
      uint256 maxLinkPayment,
      uint256 gasLimit,
      uint256 adjustedGasWei,
      uint256 linkEth
    )
  {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice executes the upkeep with the perform data returned from
   * checkUpkeep, validates the keeper's permissions, and pays the keeper.
   * @param id identifier of the upkeep to execute the data with.
   * @param performData calldata parameter to be passed to the target upkeep.
   */
  function performUpkeep(uint256 id, bytes calldata performData)
    external
    override
    whenNotPaused
    returns (bool success)
  {
    return _performUpkeepWithParams(_generatePerformParams(msg.sender, id, performData, true));
  }

  /**
   * @notice prevent an upkeep from being performed in the future
   * @param id upkeep to be canceled
   */
  function cancelUpkeep(uint256 id) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice pause an upkeep
   * @param id upkeep to be paused
   */
  function pauseUpkeep(uint256 id) external override {
    Upkeep memory upkeep = s_upkeep[id];
    requireAdminAndNotCancelled(upkeep);
    if (upkeep.paused) revert OnlyUnpausedUpkeep();
    s_upkeep[id].paused = true;
    s_upkeepIDs.remove(id);
    emit UpkeepPaused(id);
  }

  /**
   * @notice unpause an upkeep
   * @param id upkeep to be resumed
   */
  function unpauseUpkeep(uint256 id) external override {
    Upkeep memory upkeep = s_upkeep[id];
    requireAdminAndNotCancelled(upkeep);
    if (!upkeep.paused) revert OnlyPausedUpkeep();
    s_upkeep[id].paused = false;
    s_upkeepIDs.add(id);
    emit UpkeepUnpaused(id);
  }

  /**
   * @notice update the check data of an upkeep
   * @param id the id of the upkeep whose check data needs to be updated
   * @param newCheckData the new check data
   */
  function updateCheckData(uint256 id, bytes calldata newCheckData) external override {
    Upkeep memory upkeep = s_upkeep[id];
    requireAdminAndNotCancelled(upkeep);
    s_checkData[id] = newCheckData;
    emit UpkeepCheckDataUpdated(id, newCheckData);
  }

  /**
   * @notice adds LINK funding for an upkeep by transferring from the sender's
   * LINK balance
   * @param id upkeep to fund
   * @param amount number of LINK to transfer
   */
  function addFunds(uint256 id, uint96 amount) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice uses LINK's transferAndCall to LINK and add funding to an upkeep
   * @dev safe to cast uint256 to uint96 as total LINK supply is under UINT96MAX
   * @param sender the account which transferred the funds
   * @param amount number of LINK transfer
   */
  function onTokenTransfer(
    address sender,
    uint256 amount,
    bytes calldata data
  ) external override {
    if (msg.sender != address(LINK)) revert OnlyCallableByLINKToken();
    if (data.length != 32) revert InvalidDataLength();
    uint256 id = abi.decode(data, (uint256));
    if (s_upkeep[id].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();

    s_upkeep[id].balance = s_upkeep[id].balance + uint96(amount);
    s_expectedLinkBalance = s_expectedLinkBalance + amount;

    emit FundsAdded(id, sender, uint96(amount));
  }

  /**
   * @notice removes funding from a canceled upkeep
   * @param id upkeep to withdraw funds from
   * @param to destination address for sending remaining funds
   */
  function withdrawFunds(uint256 id, address to) external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice withdraws LINK funds collected through cancellation fees
   */
  function withdrawOwnerFunds() external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice allows the admin of an upkeep to modify gas limit
   * @param id upkeep to be change the gas limit for
   * @param gasLimit new gas limit for the upkeep
   */
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice recovers LINK funds improperly transferred to the registry
   * @dev In principle this function’s execution cost could exceed block
   * gas limit. However, in our anticipated deployment, the number of upkeeps and
   * keepers will be low enough to avoid this problem.
   */
  function recoverFunds() external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice withdraws a keeper's payment, callable only by the keeper's payee
   * @param from keeper address
   * @param to address to send the payment to
   */
  function withdrawPayment(address from, address to) external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice proposes the safe transfer of a keeper's payee to another address
   * @param keeper address of the keeper to transfer payee role
   * @param proposed address to nominate for next payeeship
   */
  function transferPayeeship(address keeper, address proposed) external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice accepts the safe transfer of payee role for a keeper
   * @param keeper address to accept the payee role for
   */
  function acceptPayeeship(address keeper) external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice proposes the safe transfer of an upkeep's admin role to another address
   * @param id the upkeep id to transfer admin
   * @param proposed address to nominate for the new upkeep admin
   */
  function transferUpkeepAdmin(uint256 id, address proposed) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice accepts the safe transfer of admin role for an upkeep
   * @param id the upkeep id
   */
  function acceptUpkeepAdmin(uint256 id) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice signals to keepers that they should not perform upkeeps until the
   * contract has been unpaused
   */
  function pause() external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice signals to keepers that they can perform upkeeps once again after
   * having been paused
   */
  function unpause() external {
    // Executed through logic contract
    _fallback();
  }

  // SETTERS

  /**
   * @notice updates the configuration of the registry
   * @param config registry config fields
   */
  function setConfig(Config memory config) public onlyOwner {
    if (config.maxPerformGas < s_storage.maxPerformGas) revert GasLimitCanOnlyIncrease();
    s_storage = Storage({
      paymentPremiumPPB: config.paymentPremiumPPB,
      flatFeeMicroLink: config.flatFeeMicroLink,
      blockCountPerTurn: config.blockCountPerTurn,
      checkGasLimit: config.checkGasLimit,
      stalenessSeconds: config.stalenessSeconds,
      gasCeilingMultiplier: config.gasCeilingMultiplier,
      minUpkeepSpend: config.minUpkeepSpend,
      maxPerformGas: config.maxPerformGas,
      nonce: s_storage.nonce
    });
    s_fallbackGasPrice = config.fallbackGasPrice;
    s_fallbackLinkPrice = config.fallbackLinkPrice;
    s_transcoder = config.transcoder;
    s_registrar = config.registrar;
    emit ConfigSet(config);
  }

  /**
   * @notice update the list of keepers allowed to perform upkeep
   * @param keepers list of addresses allowed to perform upkeep
   * @param payees addresses corresponding to keepers who are allowed to
   * move payments which have been accrued
   */
  function setKeepers(address[] calldata keepers, address[] calldata payees) external {
    // Executed through logic contract
    _fallback();
  }

  // GETTERS

  /**
   * @notice read all of the details about an upkeep
   */
  function getUpkeep(uint256 id)
    external
    view
    override
    returns (
      address target,
      uint32 executeGas,
      bytes memory checkData,
      uint96 balance,
      address lastKeeper,
      address admin,
      uint64 maxValidBlocknumber,
      uint96 amountSpent,
      bool paused
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
      reg.amountSpent,
      reg.paused
    );
  }

  /**
   * @notice retrieve active upkeep IDs. Active upkeep is defined as an upkeep which is not paused and not canceled.
   * @param startIndex starting index in list
   * @param maxCount max count to retrieve (0 = unlimited)
   * @dev the order of IDs in the list is **not guaranteed**, therefore, if making successive calls, one
   * should consider keeping the blockheight constant to ensure a holistic picture of the contract state
   */
  function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) external view override returns (uint256[] memory) {
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
   * @notice read the current info about any keeper address
   */
  function getKeeperInfo(address query)
    external
    view
    override
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
   * @notice read the current state of the registry
   */
  function getState()
    external
    view
    override
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
   * @notice calculates the minimum balance required for an upkeep to remain eligible
   * @param id the upkeep id to calculate minimum balance for
   */
  function getMinBalanceForUpkeep(uint256 id) external view returns (uint96 minBalance) {
    return getMaxPaymentForGas(s_upkeep[id].executeGas);
  }

  /**
   * @notice calculates the maximum payment for a given gas limit
   * @param gasLimit the gas to calculate payment for
   */
  function getMaxPaymentForGas(uint256 gasLimit) public view returns (uint96 maxPayment) {
    (uint256 fastGasWei, uint256 linkEth) = _getFeedData();
    return _calculatePaymentAmount(gasLimit, fastGasWei, linkEth, false);
  }

  /**
   * @notice retrieves the migration permission for a peer registry
   */
  function getPeerRegistryMigrationPermission(address peer) external view returns (MigrationPermission) {
    return s_peerRegistryMigrationPermission[peer];
  }

  /**
   * @notice sets the peer registry migration permission
   */
  function setPeerRegistryMigrationPermission(address peer, MigrationPermission permission) external {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @inheritdoc MigratableKeeperRegistryInterface
   */
  function migrateUpkeeps(uint256[] calldata ids, address destination) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @inheritdoc MigratableKeeperRegistryInterface
   */
  UpkeepFormat public constant override upkeepTranscoderVersion = UPKEEP_TRANSCODER_VERSION_BASE;

  /**
   * @inheritdoc MigratableKeeperRegistryInterface
   */
  function receiveUpkeeps(bytes calldata encodedUpkeeps) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @dev This is the address to which proxy functions are delegated to
   */
  function _implementation() internal view override returns (address) {
    return KEEPER_REGISTRY_LOGIC;
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available
   */
  function _callWithExactGas(
    uint256 gasAmount,
    address target,
    bytes memory data
  ) private returns (bool success) {
    assembly {
      let g := gas()
      // Compute g -= PERFORM_GAS_CUSHION and check for underflow
      if lt(g, PERFORM_GAS_CUSHION) {
        revert(0, 0)
      }
      g := sub(g, PERFORM_GAS_CUSHION)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) {
        revert(0, 0)
      }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) {
        revert(0, 0)
      }
      // call and return whether we succeeded. ignore return data
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
    return success;
  }

  /**
   * @dev calls the Upkeep target with the performData param passed in by the
   * keeper and the exact gas required by the Upkeep
   */
  function _performUpkeepWithParams(PerformParams memory params) private nonReentrant returns (bool success) {
    Upkeep memory upkeep = s_upkeep[params.id];
    if (upkeep.maxValidBlocknumber <= block.number) revert UpkeepCancelled();
    _prePerformUpkeep(upkeep, params.from, params.maxLinkPayment);

    uint256 gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(PERFORM_SELECTOR, params.performData);
    success = _callWithExactGas(params.gasLimit, upkeep.target, callData);
    gasUsed = gasUsed - gasleft();
    uint96 payment = _calculatePaymentAmount(gasUsed, params.fastGasWei, params.linkEth, true);

    s_upkeep[params.id].balance = s_upkeep[params.id].balance - payment;
    s_upkeep[params.id].amountSpent = s_upkeep[params.id].amountSpent + payment;
    s_upkeep[params.id].lastKeeper = params.from;
    s_keeperInfo[params.from].balance = s_keeperInfo[params.from].balance + payment;

    emit UpkeepPerformed(params.id, success, params.from, payment, params.performData);
    return success;
  }
}
