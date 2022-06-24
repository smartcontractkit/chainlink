// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./KeeperBase.sol";
import "./ConfirmedOwner.sol";
import "./interfaces/TypeAndVersionInterface.sol";
import "./interfaces/AggregatorV3Interface.sol";
import "./interfaces/LinkTokenInterface.sol";
import "./interfaces/KeeperCompatibleInterface.sol";
import "./interfaces/KeeperRegistryInterface.sol";
import "./interfaces/MigratableKeeperRegistryInterface.sol";
import "./interfaces/UpkeepTranscoderInterface.sol";
import "./interfaces/ERC677ReceiverInterface.sol";

/**
 * @notice Registry for adding work for Chainlink Keepers to perform on client
 * contracts. Clients must support the Upkeep interface.
 */
contract KeeperRegistry is
  TypeAndVersionInterface,
  ConfirmedOwner,
  KeeperBase,
  ReentrancyGuard,
  Pausable,
  KeeperRegistryExecutableInterface,
  MigratableKeeperRegistryInterface,
  ERC677ReceiverInterface
{
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  address private constant ZERO_ADDRESS = address(0);
  address private constant IGNORE_ADDRESS = 0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF;
  bytes4 private constant CHECK_SELECTOR = KeeperCompatibleInterface.checkUpkeep.selector;
  bytes4 private constant PERFORM_SELECTOR = KeeperCompatibleInterface.performUpkeep.selector;
  uint256 private constant PERFORM_GAS_MIN = 2_300;
  uint256 private constant CANCELATION_DELAY = 50;
  uint256 private constant PERFORM_GAS_CUSHION = 5_000;
  uint256 private constant REGISTRY_GAS_OVERHEAD = 80_000;
  uint256 private constant PPB_BASE = 1_000_000_000;
  uint64 private constant UINT64_MAX = 2**64 - 1;
  uint96 private constant LINK_TOTAL_SUPPLY = 1e27;

  address[] private s_keeperList;
  EnumerableSet.UintSet private s_upkeepIDs;
  mapping(uint256 => Upkeep) private s_upkeep;
  mapping(address => KeeperInfo) private s_keeperInfo;
  mapping(address => address) private s_proposedPayee;
  mapping(uint256 => bytes) private s_checkData;
  mapping(address => MigrationPermission) private s_peerRegistryMigrationPermission;
  Storage private s_storage;
  uint256 private s_fallbackGasPrice; // not in config object for gas savings
  uint256 private s_fallbackLinkPrice; // not in config object for gas savings
  uint96 private s_ownerLinkBalance;
  uint256 private s_expectedLinkBalance;
  address private s_transcoder;
  address private s_registrar;

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;
  AggregatorV3Interface public immutable FAST_GAS_FEED;

  /**
   * @notice versions:
   * - KeeperRegistry 1.2.0: allow funding within performUpkeep
   *                       : allow configurable registry maxPerformGas
   *                       : add function to let admin change upkeep gas limit
   *                       : add minUpkeepSpend requirement
                           : upgrade to solidity v0.8
   * - KeeperRegistry 1.1.0: added flatFeeMicroLink
   * - KeeperRegistry 1.0.0: initial release
   */
  string public constant override typeAndVersion = "KeeperRegistry 1.2.0";

  error CannotCancel();
  error UpkeepNotActive();
  error MigrationNotPermitted();
  error UpkeepNotCanceled();
  error UpkeepNotNeeded();
  error NotAContract();
  error PaymentGreaterThanAllLINK();
  error OnlyActiveKeepers();
  error InsufficientFunds();
  error KeepersMustTakeTurns();
  error ParameterLengthError();
  error OnlyCallableByOwnerOrAdmin();
  error OnlyCallableByLINKToken();
  error InvalidPayee();
  error DuplicateEntry();
  error ValueNotChanged();
  error IndexOutOfRange();
  error TranscoderNotSet();
  error ArrayHasNoEntries();
  error GasLimitOutsideRange();
  error OnlyCallableByPayee();
  error OnlyCallableByProposedPayee();
  error GasLimitCanOnlyIncrease();
  error OnlyCallableByAdmin();
  error OnlyCallableByOwnerOrRegistrar();
  error InvalidRecipient();
  error InvalidDataLength();
  error TargetCheckReverted(bytes reason);

  enum MigrationPermission {
    NONE,
    OUTGOING,
    INCOMING,
    BIDIRECTIONAL
  }

  /**
   * @notice storage of the registry, contains a mix of config and state data
   */
  struct Storage {
    uint32 paymentPremiumPPB;
    uint32 flatFeeMicroLink;
    uint24 blockCountPerTurn;
    uint32 checkGasLimit;
    uint24 stalenessSeconds;
    uint16 gasCeilingMultiplier;
    uint96 minUpkeepSpend; // 1 evm word
    uint32 maxPerformGas;
    uint32 nonce; // 2 evm words
  }

  struct Upkeep {
    uint96 balance;
    address lastKeeper; // 1 storage slot full
    uint32 executeGas;
    uint64 maxValidBlocknumber;
    address target; // 2 storage slots full
    uint96 amountSpent;
    address admin; // 3 storage slots full
  }

  struct KeeperInfo {
    address payee;
    uint96 balance;
    bool active;
  }

  struct PerformParams {
    address from;
    uint256 id;
    bytes performData;
    uint256 maxLinkPayment;
    uint256 gasLimit;
    uint256 adjustedGasWei;
    uint256 linkEth;
  }

  event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin);
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    address indexed from,
    uint96 payment,
    bytes performData
  );
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);
  event OwnerFundsWithdrawn(uint96 amount);
  event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination);
  event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom);
  event ConfigSet(Config config);
  event KeepersUpdated(address[] keepers, address[] payees);
  event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee);
  event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to);
  event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to);
  event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit);

  /**
   * @param link address of the LINK Token
   * @param linkEthFeed address of the LINK/ETH price feed
   * @param fastGasFeed address of the Fast Gas price feed
   * @param config registry config settings
   */
  constructor(
    address link,
    address linkEthFeed,
    address fastGasFeed,
    Config memory config
  ) ConfirmedOwner(msg.sender) {
    LINK = LinkTokenInterface(link);
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
    FAST_GAS_FEED = AggregatorV3Interface(fastGasFeed);
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
  ) external override onlyOwnerOrRegistrar returns (uint256 id) {
    id = uint256(keccak256(abi.encodePacked(blockhash(block.number - 1), address(this), s_storage.nonce)));
    _createUpkeep(id, target, gasLimit, admin, 0, checkData);
    s_storage.nonce++;
    emit UpkeepRegistered(id, gasLimit, admin);
    return id;
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
    uint64 maxValid = s_upkeep[id].maxValidBlocknumber;
    bool canceled = maxValid != UINT64_MAX;
    bool isOwner = msg.sender == owner();

    if (canceled && !(isOwner && maxValid > block.number)) revert CannotCancel();
    if (!isOwner && msg.sender != s_upkeep[id].admin) revert OnlyCallableByOwnerOrAdmin();

    uint256 height = block.number;
    if (!isOwner) {
      height = height + CANCELATION_DELAY;
    }
    s_upkeep[id].maxValidBlocknumber = uint64(height);
    s_upkeepIDs.remove(id);

    emit UpkeepCanceled(id, uint64(height));
  }

  /**
   * @notice adds LINK funding for an upkeep by transferring from the sender's
   * LINK balance
   * @param id upkeep to fund
   * @param amount number of LINK to transfer
   */
  function addFunds(uint256 id, uint96 amount) external override onlyActiveUpkeep(id) {
    s_upkeep[id].balance = s_upkeep[id].balance + amount;
    s_expectedLinkBalance = s_expectedLinkBalance + amount;
    LINK.transferFrom(msg.sender, address(this), amount);
    emit FundsAdded(id, msg.sender, amount);
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
    if (s_upkeep[id].maxValidBlocknumber != UINT64_MAX) revert UpkeepNotActive();

    s_upkeep[id].balance = s_upkeep[id].balance + uint96(amount);
    s_expectedLinkBalance = s_expectedLinkBalance + amount;

    emit FundsAdded(id, sender, uint96(amount));
  }

  /**
   * @notice removes funding from a canceled upkeep
   * @param id upkeep to withdraw funds from
   * @param to destination address for sending remaining funds
   */
  function withdrawFunds(uint256 id, address to) external validRecipient(to) onlyUpkeepAdmin(id) {
    if (s_upkeep[id].maxValidBlocknumber > block.number) revert UpkeepNotCanceled();

    uint96 minUpkeepSpend = s_storage.minUpkeepSpend;
    uint96 amountLeft = s_upkeep[id].balance;
    uint96 amountSpent = s_upkeep[id].amountSpent;

    uint96 cancellationFee = 0;
    // cancellationFee is supposed to be min(max(minUpkeepSpend - amountSpent,0), amountLeft)
    if (amountSpent < minUpkeepSpend) {
      cancellationFee = minUpkeepSpend - amountSpent;
      if (cancellationFee > amountLeft) {
        cancellationFee = amountLeft;
      }
    }
    uint96 amountToWithdraw = amountLeft - cancellationFee;

    s_upkeep[id].balance = 0;
    s_ownerLinkBalance = s_ownerLinkBalance + cancellationFee;

    s_expectedLinkBalance = s_expectedLinkBalance - amountToWithdraw;
    emit FundsWithdrawn(id, amountToWithdraw, to);

    LINK.transfer(to, amountToWithdraw);
  }

  /**
   * @notice withdraws LINK funds collected through cancellation fees
   */
  function withdrawOwnerFunds() external onlyOwner {
    uint96 amount = s_ownerLinkBalance;

    s_expectedLinkBalance = s_expectedLinkBalance - amount;
    s_ownerLinkBalance = 0;

    emit OwnerFundsWithdrawn(amount);
    LINK.transfer(msg.sender, amount);
  }

  /**
   * @notice allows the admin of an upkeep to modify gas limit
   * @param id upkeep to be change the gas limit for
   * @param gasLimit new gas limit for the upkeep
   */
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external override onlyActiveUpkeep(id) onlyUpkeepAdmin(id) {
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_storage.maxPerformGas) revert GasLimitOutsideRange();

    s_upkeep[id].executeGas = gasLimit;

    emit UpkeepGasLimitSet(id, gasLimit);
  }

  /**
   * @notice recovers LINK funds improperly transferred to the registry
   * @dev In principle this function’s execution cost could exceed block
   * gas limit. However, in our anticipated deployment, the number of upkeeps and
   * keepers will be low enough to avoid this problem.
   */
  function recoverFunds() external onlyOwner {
    uint256 total = LINK.balanceOf(address(this));
    LINK.transfer(msg.sender, total - s_expectedLinkBalance);
  }

  /**
   * @notice withdraws a keeper's payment, callable only by the keeper's payee
   * @param from keeper address
   * @param to address to send the payment to
   */
  function withdrawPayment(address from, address to) external validRecipient(to) {
    KeeperInfo memory keeper = s_keeperInfo[from];
    if (keeper.payee != msg.sender) revert OnlyCallableByPayee();

    s_keeperInfo[from].balance = 0;
    s_expectedLinkBalance = s_expectedLinkBalance - keeper.balance;
    emit PaymentWithdrawn(from, keeper.balance, to, msg.sender);

    LINK.transfer(to, keeper.balance);
  }

  /**
   * @notice proposes the safe transfer of a keeper's payee to another address
   * @param keeper address of the keeper to transfer payee role
   * @param proposed address to nominate for next payeeship
   */
  function transferPayeeship(address keeper, address proposed) external {
    if (s_keeperInfo[keeper].payee != msg.sender) revert OnlyCallableByPayee();
    if (proposed == msg.sender) revert ValueNotChanged();

    if (s_proposedPayee[keeper] != proposed) {
      s_proposedPayee[keeper] = proposed;
      emit PayeeshipTransferRequested(keeper, msg.sender, proposed);
    }
  }

  /**
   * @notice accepts the safe transfer of payee role for a keeper
   * @param keeper address to accept the payee role for
   */
  function acceptPayeeship(address keeper) external {
    if (s_proposedPayee[keeper] != msg.sender) revert OnlyCallableByProposedPayee();
    address past = s_keeperInfo[keeper].payee;
    s_keeperInfo[keeper].payee = msg.sender;
    s_proposedPayee[keeper] = ZERO_ADDRESS;

    emit PayeeshipTransferred(keeper, past, msg.sender);
  }

  /**
   * @notice signals to keepers that they should not perform upkeeps until the
   * contract has been unpaused
   */
  function pause() external onlyOwner {
    _pause();
  }

  /**
   * @notice signals to keepers that they can perform upkeeps once again after
   * having been paused
   */
  function unpause() external onlyOwner {
    _unpause();
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
  function setKeepers(address[] calldata keepers, address[] calldata payees) external onlyOwner {
    if (keepers.length != payees.length || keepers.length < 2) revert ParameterLengthError();
    for (uint256 i = 0; i < s_keeperList.length; i++) {
      address keeper = s_keeperList[i];
      s_keeperInfo[keeper].active = false;
    }
    for (uint256 i = 0; i < keepers.length; i++) {
      address keeper = keepers[i];
      KeeperInfo storage s_keeper = s_keeperInfo[keeper];
      address oldPayee = s_keeper.payee;
      address newPayee = payees[i];
      if (
        (newPayee == ZERO_ADDRESS) || (oldPayee != ZERO_ADDRESS && oldPayee != newPayee && newPayee != IGNORE_ADDRESS)
      ) revert InvalidPayee();
      if (s_keeper.active) revert DuplicateEntry();
      s_keeper.active = true;
      if (newPayee != IGNORE_ADDRESS) {
        s_keeper.payee = newPayee;
      }
    }
    s_keeperList = keepers;
    emit KeepersUpdated(keepers, payees);
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
   * @notice retrieve active upkeep IDs
   * @param startIndex starting index in list
   * @param maxCount max count to retrieve (0 = unlimited)
   * @dev the order of IDs in the list is **not guaranteed**, therefore, if making successive calls, one
   * should consider keeping the blockheight constant to ensure a wholistic picture of the contract state
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
    (uint256 gasWei, uint256 linkEth) = _getFeedData();
    uint256 adjustedGasWei = _adjustGasPrice(gasWei, false);
    return _calculatePaymentAmount(gasLimit, adjustedGasWei, linkEth);
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
  function setPeerRegistryMigrationPermission(address peer, MigrationPermission permission) external onlyOwner {
    s_peerRegistryMigrationPermission[peer] = permission;
  }

  /**
   * @inheritdoc MigratableKeeperRegistryInterface
   */
  function migrateUpkeeps(uint256[] calldata ids, address destination) external override {
    if (
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.OUTGOING &&
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    if (s_transcoder == ZERO_ADDRESS) revert TranscoderNotSet();
    if (ids.length == 0) revert ArrayHasNoEntries();
    uint256 id;
    Upkeep memory upkeep;
    uint256 totalBalanceRemaining;
    bytes[] memory checkDatas = new bytes[](ids.length);
    Upkeep[] memory upkeeps = new Upkeep[](ids.length);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      id = ids[idx];
      upkeep = s_upkeep[id];
      if (upkeep.admin != msg.sender) revert OnlyCallableByAdmin();
      if (upkeep.maxValidBlocknumber != UINT64_MAX) revert UpkeepNotActive();
      upkeeps[idx] = upkeep;
      checkDatas[idx] = s_checkData[id];
      totalBalanceRemaining = totalBalanceRemaining + upkeep.balance;
      delete s_upkeep[id];
      delete s_checkData[id];
      s_upkeepIDs.remove(id);
      emit UpkeepMigrated(id, upkeep.balance, destination);
    }
    s_expectedLinkBalance = s_expectedLinkBalance - totalBalanceRemaining;
    bytes memory encodedUpkeeps = abi.encode(ids, upkeeps, checkDatas);
    MigratableKeeperRegistryInterface(destination).receiveUpkeeps(
      UpkeepTranscoderInterface(s_transcoder).transcodeUpkeeps(
        UpkeepFormat.V1,
        MigratableKeeperRegistryInterface(destination).upkeepTranscoderVersion(),
        encodedUpkeeps
      )
    );
    LINK.transfer(destination, totalBalanceRemaining);
  }

  /**
   * @inheritdoc MigratableKeeperRegistryInterface
   */
  UpkeepFormat public constant override upkeepTranscoderVersion = UpkeepFormat.V1;

  /**
   * @inheritdoc MigratableKeeperRegistryInterface
   */
  function receiveUpkeeps(bytes calldata encodedUpkeeps) external override {
    if (
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.INCOMING &&
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    (uint256[] memory ids, Upkeep[] memory upkeeps, bytes[] memory checkDatas) = abi.decode(
      encodedUpkeeps,
      (uint256[], Upkeep[], bytes[])
    );
    for (uint256 idx = 0; idx < ids.length; idx++) {
      _createUpkeep(
        ids[idx],
        upkeeps[idx].target,
        upkeeps[idx].executeGas,
        upkeeps[idx].admin,
        upkeeps[idx].balance,
        checkDatas[idx]
      );
      emit UpkeepReceived(ids[idx], upkeeps[idx].balance, msg.sender);
    }
  }

  /**
   * @notice creates a new upkeep with the given fields
   * @param target address to perform upkeep on
   * @param gasLimit amount of gas to provide the target contract when
   * performing upkeep
   * @param admin address to cancel upkeep and withdraw remaining funds
   * @param checkData data passed to the contract when checking for upkeep
   */
  function _createUpkeep(
    uint256 id,
    address target,
    uint32 gasLimit,
    address admin,
    uint96 balance,
    bytes memory checkData
  ) internal whenNotPaused {
    if (!target.isContract()) revert NotAContract();
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_storage.maxPerformGas) revert GasLimitOutsideRange();
    s_upkeep[id] = Upkeep({
      target: target,
      executeGas: gasLimit,
      balance: balance,
      admin: admin,
      maxValidBlocknumber: UINT64_MAX,
      lastKeeper: ZERO_ADDRESS,
      amountSpent: 0
    });
    s_expectedLinkBalance = s_expectedLinkBalance + balance;
    s_checkData[id] = checkData;
    s_upkeepIDs.add(id);
  }

  /**
   * @dev retrieves feed data for fast gas/eth and link/eth prices. if the feed
   * data is stale it uses the configured fallback price. Once a price is picked
   * for gas it takes the min of gas price in the transaction or the fast gas
   * price in order to reduce costs for the upkeep clients.
   */
  function _getFeedData() private view returns (uint256 gasWei, uint256 linkEth) {
    uint32 stalenessSeconds = s_storage.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 feedValue;
    (, feedValue, , timestamp, ) = FAST_GAS_FEED.latestRoundData();
    if ((staleFallback && stalenessSeconds < block.timestamp - timestamp) || feedValue <= 0) {
      gasWei = s_fallbackGasPrice;
    } else {
      gasWei = uint256(feedValue);
    }
    (, feedValue, , timestamp, ) = LINK_ETH_FEED.latestRoundData();
    if ((staleFallback && stalenessSeconds < block.timestamp - timestamp) || feedValue <= 0) {
      linkEth = s_fallbackLinkPrice;
    } else {
      linkEth = uint256(feedValue);
    }
    return (gasWei, linkEth);
  }

  /**
   * @dev calculates LINK paid for gas spent plus a configure premium percentage
   */
  function _calculatePaymentAmount(
    uint256 gasLimit,
    uint256 gasWei,
    uint256 linkEth
  ) private view returns (uint96 payment) {
    uint256 weiForGas = gasWei * (gasLimit + REGISTRY_GAS_OVERHEAD);
    uint256 premium = PPB_BASE + s_storage.paymentPremiumPPB;
    uint256 total = ((weiForGas * (1e9) * (premium)) / (linkEth)) + (uint256(s_storage.flatFeeMicroLink) * (1e12));
    if (total > LINK_TOTAL_SUPPLY) revert PaymentGreaterThanAllLINK();
    return uint96(total); // LINK_TOTAL_SUPPLY < UINT96_MAX
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
    returns (bool success)
  {
    Upkeep memory upkeep = s_upkeep[params.id];
    _prePerformUpkeep(upkeep, params.from, params.maxLinkPayment);

    uint256 gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(PERFORM_SELECTOR, params.performData);
    success = _callWithExactGas(params.gasLimit, upkeep.target, callData);
    gasUsed = gasUsed - gasleft();

    uint96 payment = _calculatePaymentAmount(gasUsed, params.adjustedGasWei, params.linkEth);

    s_upkeep[params.id].balance = s_upkeep[params.id].balance - payment;
    s_upkeep[params.id].amountSpent = s_upkeep[params.id].amountSpent + payment;
    s_upkeep[params.id].lastKeeper = params.from;
    s_keeperInfo[params.from].balance = s_keeperInfo[params.from].balance + payment;

    emit UpkeepPerformed(params.id, success, params.from, payment, params.performData);
    return success;
  }

  /**
   * @dev ensures all required checks are passed before an upkeep is performed
   */
  function _prePerformUpkeep(
    Upkeep memory upkeep,
    address from,
    uint256 maxLinkPayment
  ) private view {
    if (!s_keeperInfo[from].active) revert OnlyActiveKeepers();
    if (upkeep.balance < maxLinkPayment) revert InsufficientFunds();
    if (upkeep.lastKeeper == from) revert KeepersMustTakeTurns();
  }

  /**
   * @dev adjusts the gas price to min(ceiling, tx.gasprice) or just uses the ceiling if tx.gasprice is disabled
   */
  function _adjustGasPrice(uint256 gasWei, bool useTxGasPrice) private view returns (uint256 adjustedPrice) {
    adjustedPrice = gasWei * s_storage.gasCeilingMultiplier;
    if (useTxGasPrice && tx.gasprice < adjustedPrice) {
      adjustedPrice = tx.gasprice;
    }
  }

  /**
   * @dev generates a PerformParams struct for use in _performUpkeepWithParams()
   */
  function _generatePerformParams(
    address from,
    uint256 id,
    bytes memory performData,
    bool useTxGasPrice
  ) private view returns (PerformParams memory) {
    uint256 gasLimit = s_upkeep[id].executeGas;
    (uint256 gasWei, uint256 linkEth) = _getFeedData();
    uint256 adjustedGasWei = _adjustGasPrice(gasWei, useTxGasPrice);
    uint96 maxLinkPayment = _calculatePaymentAmount(gasLimit, adjustedGasWei, linkEth);

    return
      PerformParams({
        from: from,
        id: id,
        performData: performData,
        maxLinkPayment: maxLinkPayment,
        gasLimit: gasLimit,
        adjustedGasWei: adjustedGasWei,
        linkEth: linkEth
      });
  }

  // MODIFIERS

  /**
   * @dev ensures a upkeep is valid
   */
  modifier validUpkeep(uint256 id) {
    if (s_upkeep[id].maxValidBlocknumber <= block.number) revert UpkeepNotActive();
    _;
  }

  /**
   * @dev Reverts if called by anyone other than the admin of upkeep #id
   */
  modifier onlyUpkeepAdmin(uint256 id) {
    if (msg.sender != s_upkeep[id].admin) revert OnlyCallableByAdmin();
    _;
  }

  /**
   * @dev Reverts if called on a cancelled upkeep
   */
  modifier onlyActiveUpkeep(uint256 id) {
    if (s_upkeep[id].maxValidBlocknumber != UINT64_MAX) revert UpkeepNotActive();
    _;
  }

  /**
   * @dev ensures that burns don't accidentally happen by sending to the zero
   * address
   */
  modifier validRecipient(address to) {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    _;
  }

  /**
   * @dev Reverts if called by anyone other than the contract owner or registrar.
   */
  modifier onlyOwnerOrRegistrar() {
    if (msg.sender != owner() && msg.sender != s_registrar) revert OnlyCallableByOwnerOrRegistrar();
    _;
  }
}
