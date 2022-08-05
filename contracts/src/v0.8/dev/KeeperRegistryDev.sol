// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./KeeperRegistryBase.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import {KeeperRegistryExecutableInterface} from "./interfaces/KeeperRegistryInterfaceDev.sol";
import "../interfaces/MigratableKeeperRegistryInterface.sol";
import "../interfaces/ERC677ReceiverInterface.sol";

/**
 * @notice Registry for adding work for Chainlink Keepers to perform on client
 * contracts. Clients must support the Upkeep interface.
 */
contract KeeperRegistryDev is
  KeeperRegistryBase,
  Proxy,
  KeeperRegistryExecutableInterface,
  MigratableKeeperRegistryInterface,
  ERC677ReceiverInterface
{
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  address public immutable KEEPER_REGISTRY_LOGIC;

  /**
   * @notice versions:
   * - KeeperRegistry 2.0.0: Split contract into Proxy and Logic
   * - KeeperRegistry 1.2.0: allow funding within performUpkeep
   *                       : allow configurable registry maxPerformGas
   *                       : add function to let admin change upkeep gas limit
   *                       : add minUpkeepSpend requirement
                           : upgrade to solidity v0.8
   * - KeeperRegistry 1.1.0: added flatFeeMicroLink
   * - KeeperRegistry 1.0.0: initial release
   */
  string public constant override typeAndVersion = "KeeperRegistry 2.0.0";

  /**
   * @param paymentModel one of Default, Arbitrum, and Optimism
   * @param registryGasOverhead the gas overhead used by registry in performUpkeep
   * @param link address of the LINK Token
   * @param linkEthFeed address of the LINK/ETH price feed
   * @param fastGasFeed address of the Fast Gas price feed
   * @param params registry parameters settings
   */
  constructor(
    PaymentModel paymentModel,
    uint256 registryGasOverhead,
    address link,
    address linkEthFeed,
    address fastGasFeed,
    address keeperRegistryLogic,
    RegistryParams memory params
  ) KeeperRegistryBase(link, linkEthFeed, fastGasFeed) {
    KEEPER_REGISTRY_LOGIC = keeperRegistryLogic;
    setRegistryParams(params);
  }

  /**
   * @inheritdoc OCR2Abstract
   */
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @inheritdoc OCR2Abstract
   */
  function latestConfigDetails()
    external
    view
    override
    returns (
      uint32 configCount,
      uint32 blockNumber,
      bytes32 configDigest
    )
  {
    return (s_configCount, s_latestConfigBlockNumber, s_latestConfigDigest);
  }

  /**
   * @inheritdoc OCR2Abstract
   */
  function latestConfigDigestAndEpoch()
    external
    view
    override
    returns (
      bool scanLogs,
      bytes32 configDigest,
      uint32 epoch
    )
  {
    return (true, configDigest, epoch);
  }

  /**
   * @inheritdoc OCR2Abstract
   */
  function transmit(
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs // signatures
  ) external override whenNotPaused {
    uint256 initialGas = gasleft(); // This line must come first

    HotVars memory hotVars = s_hotVars;

    uint40 epochAndRound = uint40(uint256(reportContext[1]));
    // TODO: Convert revret strings to errros

    // TODO: Do we need this since we'll check last block of perform
    require(hotVars.latestEpochAndRound < epochAndRound, "stale report");

    require(s_transmitters[msg.sender].active, "unauthorized transmitter");

    require(s_latestConfigDigest == reportContext[0], "configDigest mismatch");

    // TODO: Maybe don't need this as payment calculation is separate in keepers
    _requireExpectedMsgDataLength(report, rs, ss);

    require(rs.length == hotVars.f + 1, "wrong number of signatures");
    require(rs.length == ss.length, "signatures out of registration");

    // Verify signatures attached to report
    {
      bytes32 h = keccak256(abi.encodePacked(keccak256(report), reportContext));

      // i-th byte counts number of sigs made by i-th signer
      uint256 signedCount = 0;

      Signer memory signer;
      for (uint256 i = 0; i < rs.length; i++) {
        address signerAddress = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
        signer = s_signers[signerAddress];
        require(signer.active, "signature error");
        unchecked {
          signedCount += 1 << (8 * signer.index);
        }
      }

      // The first byte of the mask can be 0, because we only ever have 31 oracles
      require(
        signedCount & 0x0001010101010101010101010101010101010101010101010101010101010101 == signedCount,
        "duplicate signer"
      );
    }

    int192 juelsPerFeeCoin = _report(hotVars, reportContext[0], epochAndRound, report);
  }

  function _report(
    HotVars memory hotVars,
    bytes32 configDigest,
    uint40 epochAndRound,
    bytes memory rawReport
  ) internal returns (int192 juelsPerFeeCoin) {
    Report memory report = _decodeReport(rawReport);

    hotVars.latestEpochAndRound = epochAndRound;

    hotVars.latestAggregatorRoundId++;

    // performUpkeep here

    // persist updates to hotVars
    s_hotVars = hotVars;
  }

  struct Report {
    uint256 upkeepId; // Id of upkeep
    uint32 checkBlockNum; // Block number at which checkUpkeep was true
    bytes performData; // Perform Data for the upkeep
    uint256 linkNativePrice; // Price of link to native token (18 decimals)
  }

  // _decodeReport decodes a serialized report into a Report struct
  function _decodeReport(bytes memory rawReport) internal pure returns (Report memory) {
    uint256 upkeepId;
    uint32 checkBlockNum;
    bytes performData;
    uint256 linkNativePrice;
    (upkeepId, checkBlockNum, performData, linkNativePrice) = abi.decode(rawReport, (uint256, uint32, bytes, uint256));

    return
      Report({
        upkeepId: upkeepId,
        checkBlockNum: checkBlockNum,
        performData: performData,
        linkNativePrice: linkNativePrice
      });
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
   * TODO: update docs
   * @notice executes the upkeep with the perform data returned from
   * checkUpkeep, validates the keeper's permissions, and pays the keeper.
   * @param id identifier of the upkeep to execute the data with.
   * @param performData calldata parameter to be passed to the target upkeep.
   */
  function simulatePerformUpkeep(uint256 id, bytes calldata performData)
    external
    cannotExecute
    whenNotPaused
    returns (bool success, uint256 gasUsed)
  {
    return _performUpkeepWithParams(_generatePerformParams(id, performData, false, 0));
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
    if (s_upkeep[id].maxValidBlocknumber != UINT64_MAX) revert UpkeepCancelled();

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
   * @param params registry parameter fields
   */
  // TODO: Try to combine with onChain config
  function setRegistryParams(RegistryParams memory params) public onlyOwner {
    if (params.maxPerformGas < s_storage.maxPerformGas) revert GasLimitCanOnlyIncrease();
    s_storage = Storage({
      paymentPremiumPPB: params.paymentPremiumPPB,
      flatFeeMicroLink: params.flatFeeMicroLink,
      blockCountPerTurn: params.blockCountPerTurn,
      checkGasLimit: params.checkGasLimit,
      stalenessSeconds: params.stalenessSeconds,
      gasCeilingMultiplier: params.gasCeilingMultiplier,
      minUpkeepSpend: params.minUpkeepSpend,
      maxPerformGas: params.maxPerformGas,
      nonce: s_storage.nonce
    });
    s_fallbackGasPrice = params.fallbackGasPrice;
    s_fallbackLinkPrice = params.fallbackLinkPrice;
    s_transcoder = params.transcoder;
    s_registrar = params.registrar;
    emit RegistryParamsSet(params);
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
      RegistryParams memory params,
      address[] memory keepers
    )
  {
    Storage memory store = s_storage;
    state.nonce = store.nonce;
    state.ownerLinkBalance = s_ownerLinkBalance;
    state.expectedLinkBalance = s_expectedLinkBalance;
    state.numUpkeeps = s_upkeepIDs.length();
    params.paymentPremiumPPB = store.paymentPremiumPPB;
    params.flatFeeMicroLink = store.flatFeeMicroLink;
    params.blockCountPerTurn = store.blockCountPerTurn;
    params.checkGasLimit = store.checkGasLimit;
    params.stalenessSeconds = store.stalenessSeconds;
    params.gasCeilingMultiplier = store.gasCeilingMultiplier;
    params.minUpkeepSpend = store.minUpkeepSpend;
    params.maxPerformGas = store.maxPerformGas;
    params.fallbackGasPrice = s_fallbackGasPrice;
    params.fallbackLinkPrice = s_fallbackLinkPrice;
    params.transcoder = s_transcoder;
    params.registrar = s_registrar;
    return (state, params, s_keeperList);
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
  UpkeepFormat public constant override upkeepTranscoderVersion = UPKEEP_TRANSCODER_VESION_BASE;

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
  function _performUpkeepWithParams(PerformParams memory params)
    private
    nonReentrant
    validUpkeep(params.id)
    returns (bool success, uint256 gasUsed)
  {
    if (params.upkeep.paused) revert OnlyUnpausedUpkeep();
    if (params.upkeep.balance < maxLinkPayment) revert InsufficientFunds();

    uint256 gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(PERFORM_SELECTOR, params.performData);
    success = _callWithExactGas(params.gasLimit, upkeep.target, callData);
    gasUsed = gasUsed - gasleft();

    return (success, gasUsed);
  }

  // TODO: Move this to transmit
  //emit UpkeepPerformed(params.id, success, params.from, payment, params.performData);

  // TODO: Call this in transmit
  function _processUpkeepPayment(PerformParams memory params, uint256 gasUsed) private {
    uint96 payment = _calculatePaymentAmount(gasUsed, params.adjustedGasWei, params.linkEth);
    s_upkeep[params.id].balance = s_upkeep[params.id].balance - payment;
    s_upkeep[params.id].amountSpent = s_upkeep[params.id].amountSpent + payment;
    s_upkeep[params.id].lastKeeper = params.from;
    s_keeperInfo[params.from].balance = s_keeperInfo[params.from].balance + payment;
  }
}
