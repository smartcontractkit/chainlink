// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./KeeperRegistryBase2_0.sol";
import {KeeperRegistryExecutableInterface} from "./interfaces/KeeperRegistryInterface2_0.sol";
import "../../interfaces/MigratableKeeperRegistryInterface.sol";
import "../../interfaces/ERC677ReceiverInterface.sol";
import "./interfaces/OCR2Abstract.sol";

/**
 * @notice Registry for adding work for Chainlink Keepers to perform on client
 * contracts. Clients must support the Upkeep interface.
 */
contract KeeperRegistry2_0 is
  KeeperRegistryBase2_0,
  Proxy,
  OCR2Abstract,
  KeeperRegistryExecutableInterface,
  MigratableKeeperRegistryInterface,
  ERC677ReceiverInterface
{
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  address public immutable KEEPER_REGISTRY_LOGIC;

  /**
   * @notice versions:
   * - KeeperRegistry 2.0.0: implement OCR interface
   * - KeeperRegistry 1.3.0: split contract into Proxy and Logic
   *                       : account for Arbitrum and Optimism L1 gas fee
   *                       : allow users to configure upkeeps
   * - KeeperRegistry 1.2.0: allow funding within performUpkeep
   *                       : allow configurable registry maxPerformGas
   *                       : add function to let admin change upkeep gas limit
   *                       : add minUpkeepSpend requirement
   *                       : upgrade to solidity v0.8
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
   * @param onChainConfig registry on chain config settings
   */
  constructor(
    PaymentModel paymentModel,
    uint256 registryGasOverhead,
    address link,
    address linkEthFeed,
    address fastGasFeed,
    address keeperRegistryLogic,
    OnChainConfig memory onChainConfig
  ) KeeperRegistryBase2_0(paymentModel, registryGasOverhead, link, linkEthFeed, fastGasFeed) {
    KEEPER_REGISTRY_LOGIC = keeperRegistryLogic;
    setOnChainConfig(onChainConfig);
  }

  // ACTIONS

  /**
   * @inheritdoc OCR2Abstract
   */
  function transmit(
    bytes32[4] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs // signatures
  ) external override whenNotPaused {
    if (!s_transmitters[msg.sender].active) revert OnlyActiveKeepers();
    // reportContext consists of:
    // reportContext[0]: OCR instance index
    // reportContext[1]: ConfigDigest
    // reportContext[2]: 27 byte padding, 4-byte epoch and 1-byte round
    // reportContext[3]: ExtraHash
    if (s_latestRootConfigDigest ^ reportContext[0] != reportContext[1]) revert ConfigDisgestMismatch();
    if (rs.length != s_f + 1 || rs.length != ss.length) revert IncorrectNumberOfSignatures();

    uint8[] memory signerIndices = new uint8[](rs.length);
    // Verify signatures attached to report
    {
      bytes32 h = keccak256(abi.encodePacked(keccak256(report), reportContext));
      // i-th byte counts number of sigs made by i-th signer
      uint256 signedCount = 0;

      Signer memory signer;
      for (uint256 i = 0; i < rs.length; i++) {
        address signerAddress = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
        signer = s_signers[signerAddress];
        if (!signer.active) revert OnlyActiveSigners();
        unchecked {
          signedCount += 1 << (8 * signer.index);
        }
        signerIndices[i] = signer.index;
      }
      // The first byte of the mask can be 0, because we only ever have 31 oracles
      if (signedCount & 0x0001010101010101010101010101010101010101010101010101010101010101 != signedCount)
        revert DuplicateSigners();
    }

    // Deocde the report and performUpkeep
    Report memory parsedReport = _decodeReport(report);
    if (parsedReport.checkBlockNumber <= s_upkeep[parsedReport.upkeepId].lastPerformBlockNumber) revert StaleReport();

    // Perform target upkeep
    PerformParams memory params = _generatePerformParams(parsedReport.upkeepId, parsedReport.performData, true);
    (bool success, uint256 gasUsed) = _performUpkeepWithParams(params);

    // Calculate actual payment amount
    (uint96 gasPayment, uint96 premium) = _calculatePaymentAmount(gasUsed, params.fastGasWei, params.linkEth, true);
    uint96 premiumPerSigner = premium / uint96(rs.length);
    uint96 totalPayment = gasPayment + premiumPerSigner * uint96(rs.length);

    s_upkeep[parsedReport.upkeepId].balance = s_upkeep[params.id].balance - totalPayment;
    s_upkeep[parsedReport.upkeepId].amountSpent = s_upkeep[params.id].amountSpent + totalPayment;
    s_upkeep[parsedReport.upkeepId].lastPerformBlockNumber = uint32(block.number);

    s_transmitters[msg.sender].balance = s_transmitters[msg.sender].balance + gasPayment;
    for (uint256 i = 0; i < rs.length; i++) {
      address transmitterToPay = s_transmittersList[signerIndices[i]];
      s_transmitters[transmitterToPay].balance += premiumPerSigner;
    }

    emit UpkeepPerformed(parsedReport.upkeepId, success, gasUsed, parsedReport.checkBlockNumber, totalPayment);
  }

  // TODO(sc-50641): Evaluate if we need link/eth in the report and finalize the fields
  struct Report {
    uint256 upkeepId; // Id of upkeep
    bytes performData; // Perform Data for the upkeep
    uint32 checkBlockNumber; // Block number at which checkUpkeep was true
  }

  // _decodeReport decodes a serialized report into a Report struct
  function _decodeReport(bytes memory rawReport) internal pure returns (Report memory) {
    uint256 upkeepId;

    bytes memory performData;
    uint32 checkBlockNumber;
    (upkeepId, performData, checkBlockNumber) = abi.decode(rawReport, (uint256, bytes, uint32));

    return Report({upkeepId: upkeepId, performData: performData, checkBlockNumber: checkBlockNumber});
  }

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
   * performed. It returns the success status / failure reason along with the perform data payload.
   * @param id identifier of the upkeep to check
   */
  function checkUpkeep(uint256 id)
    external
    override
    cannotExecute
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed
    )
  {
    // Executed through logic contract
    _fallback();
  }

  /**
   * @notice simulates the upkeep with the perform data returned from
   * checkUpkeep
   * @param id identifier of the upkeep to execute the data with.
   * @param performData calldata parameter to be passed to the target upkeep.
   */
  function simulatePerformUpkeep(uint256 id, bytes calldata performData)
    external
    cannotExecute
    whenNotPaused
    returns (bool success, uint256 gasUsed)
  {
    return _performUpkeepWithParams(_generatePerformParams(id, performData, false));
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
   * @inheritdoc OCR2Abstract
   */
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override onlyOwner {
    if (signers.length > maxNumOracles) revert TooManyOracles();
    if (f == 0) revert IncorrectNumberOfFaultyOracles();
    if (signers.length != transmitters.length || signers.length <= 3 * f)
      if (onchainConfig.length != 0) revert OnchainConfigNonEmpty();

    // remove any old signer/transmitter addresses
    uint256 oldLength = s_signersList.length;
    address signer;
    address transmitter;
    for (uint256 i = 0; i < oldLength; i++) {
      signer = s_signersList[i];
      transmitter = s_transmittersList[i];
      delete s_signers[signer];
      // Do not delete the whole transmitter struct as it has balance information stored
      s_transmitters[transmitter].active = false;
    }
    delete s_signersList;
    delete s_transmittersList;

    // add new signer/transmitter addresses
    for (uint256 i = 0; i < signers.length; i++) {
      if (s_signers[signers[i]].active) revert RepeatedSigner();
      s_signers[signers[i]] = Signer({active: true, index: uint8(i)});

      if (s_transmitters[transmitters[i]].active) revert RepeatedTransmitter();
      s_transmitters[transmitters[i]].active = true;
      s_transmitters[transmitters[i]].index = uint8(i);
    }
    s_signersList = signers;
    s_transmittersList = transmitters;
    s_f = f;
    s_offchainConfigVersion = offchainConfigVersion;
    s_offchainConfig = offchainConfig;

    _computeAndStoreConfigDigest(
      signers,
      transmitters,
      f,
      abi.encode(s_onChainConfig),
      offchainConfigVersion,
      offchainConfig
    );
  }

  /**
   * @notice updates the configuration of the registry
   * @param onChainConfig registry config fields
   */
  function setOnChainConfig(OnChainConfig memory onChainConfig) public onlyOwner {
    if (onChainConfig.maxPerformGas < s_onChainConfig.maxPerformGas) revert GasLimitCanOnlyIncrease();
    s_onChainConfig = onChainConfig;

    _computeAndStoreConfigDigest(
      s_signersList,
      s_transmittersList,
      s_f,
      abi.encode(s_onChainConfig),
      s_offchainConfigVersion,
      s_offchainConfig
    );
    emit OnChainConfigSet(onChainConfig);
  }

  /**
   * @dev Should be called on every config change, either OCR or onChainConfig
   * Recomputed the config digest and stores it
   */
  function _computeAndStoreConfigDigest(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal {
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;

    s_latestRootConfigDigest = _configDigestFromConfigData(
      block.chainid,
      address(this),
      s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    emit ConfigSet(
      previousConfigBlockNumber,
      s_latestRootConfigDigest,
      s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  /**
   * @notice update the list of payees corresponding to the transmitters
   * @param payees addresses corresponding to transmitters who are allowed to
   * move payments which have been accrued
   */
  function setPayees(address[] calldata payees) external {
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
      address admin,
      uint64 maxValidBlocknumber,
      uint32 lastPerformBlockNumber,
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
      reg.admin,
      reg.maxValidBlocknumber,
      reg.lastPerformBlockNumber,
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
      bool active,
      uint8 index,
      uint96 balance,
      address payee
    )
  {
    Transmitter memory keeper = s_transmitters[query];
    return (keeper.active, keeper.index, keeper.balance, keeper.payee);
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
      OnChainConfig memory config,
      address[] memory signers,
      address[] memory transmitters,
      uint8 f,
      uint64 offchainConfigVersion,
      bytes memory offchainConfig
    )
  {
    state.nonce = s_nonce;
    state.ownerLinkBalance = s_ownerLinkBalance;
    state.expectedLinkBalance = s_expectedLinkBalance;
    state.numUpkeeps = s_upkeepIDs.length();
    return (state, s_onChainConfig, s_signersList, s_transmittersList, s_f, s_offchainConfigVersion, s_offchainConfig);
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
    (uint96 gasPayment, uint96 premium) = _calculatePaymentAmount(gasLimit, fastGasWei, linkEth, false);
    return gasPayment + premium;
  }

  /**
   * @notice retrieves the migration permission for a peer registry
   */
  function getPeerRegistryMigrationPermission(address peer) external view returns (MigrationPermission) {
    return s_peerRegistryMigrationPermission[peer];
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
      bytes32 rootConfigDigest
    )
  {
    return (s_configCount, s_latestConfigBlockNumber, s_latestRootConfigDigest);
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

  // MIGRATION

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
  function _performUpkeepWithParams(PerformParams memory params)
    private
    nonReentrant
    returns (bool success, uint256 gasUsed)
  {
    Upkeep memory upkeep = s_upkeep[params.id];
    // TODO (sc-50783) Move these reverts to transmit before sig verification to optimise gas
    if (upkeep.maxValidBlocknumber <= block.number) revert UpkeepCancelled();
    if (upkeep.paused) revert OnlyUnpausedUpkeep();
    if (upkeep.balance < params.maxLinkPayment) revert InsufficientFunds();

    uint256 gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(PERFORM_SELECTOR, params.performData);
    success = _callWithExactGas(upkeep.executeGas, upkeep.target, callData);
    gasUsed = gasUsed - gasleft();

    return (success, gasUsed);
  }
}
