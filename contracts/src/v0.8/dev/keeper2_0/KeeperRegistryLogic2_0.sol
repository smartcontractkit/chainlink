// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./KeeperRegistryBase2_0.sol";
import "../interfaces/MigratableKeeperRegistryInterfaceDev.sol";
import "../interfaces/UpkeepTranscoderInterfaceDev.sol";

/**
 * @notice Logic contract, works in tandem with KeeperRegistry as a proxy
 */
contract KeeperRegistryLogic2_0 is KeeperRegistryBase2_0 {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  /**
   * @param paymentModel one of Default, Arbitrum, Optimism
   * @param registryGasOverhead the gas overhead used by registry in performUpkeep
   * @param link address of the LINK Token
   * @param linkEthFeed address of the LINK/ETH price feed
   * @param fastGasFeed address of the Fast Gas price feed
   */
  constructor(
    PaymentModel paymentModel,
    uint256 registryGasOverhead,
    address link,
    address linkEthFeed,
    address fastGasFeed
  ) KeeperRegistryBase2_0(paymentModel, registryGasOverhead, link, linkEthFeed, fastGasFeed) {}

  function checkUpkeep(uint256 id)
    external
    cannotExecute
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed
    )
  {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.paused) return (false, performData, UpkeepFailureReason.UPKEEP_PAUSED, gasUsed);

    gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(CHECK_SELECTOR, s_checkData[id]);
    (bool success, bytes memory result) = upkeep.target.call{gas: s_onChainConfig.checkGasLimit}(callData);
    gasUsed = gasUsed - gasleft();

    if (!success) return (false, performData, UpkeepFailureReason.TARGET_CHECK_REVERTED, gasUsed);

    (upkeepNeeded, performData) = abi.decode(result, (bool, bytes));
    if (!upkeepNeeded) return (false, performData, UpkeepFailureReason.UPKEEP_NOT_NEEDED, gasUsed);

    PerformParams memory params = _generatePerformParams(id, performData, false);
    if (upkeep.balance < params.maxLinkPayment)
      return (false, performData, UpkeepFailureReason.INSUFFICIENT_BALANCE, gasUsed);

    return (true, performData, UpkeepFailureReason.NONE, gasUsed);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    uint16 numOcrInstances,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override onlyOwner {
    if (signers.length > maxNumOracles) revert TooManyOracles();
    if (f == 0) revert IncorrectNumberOfFaultyOracles();
    if (signers.length != transmitters.length || signers.length <= 3 * f) revert IncorrectNumberOfSigners();

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
      numOcrInstances,
      abi.encode(s_onChainConfig),
      offchainConfigVersion,
      offchainConfig
    );
  }

  /**
   * @dev Unimplemented on logic contract, implementation lives on KeeperRegistry main contract
   */
  function latestConfigDetails()
    external
    view
    override
    returns (
      uint32 configCount,
      uint32 blockNumber,
      bytes32 configDigest,
      bytes32[] memory configDigests
    )
  {}

  /**
   * @dev Unimplemented on logic contract, implementation lives on KeeperRegistry main contract
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
  {}

  /**
   * @dev Unimplemented on logic contract, implementation lives on KeeperRegistry main contract
   */
  function transmit(
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external override {}

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function withdrawOwnerFunds() external onlyOwner {
    uint96 amount = s_ownerLinkBalance;

    s_expectedLinkBalance = s_expectedLinkBalance - amount;
    s_ownerLinkBalance = 0;

    emit OwnerFundsWithdrawn(amount);
    LINK.transfer(msg.sender, amount);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function recoverFunds() external onlyOwner {
    uint256 total = LINK.balanceOf(address(this));
    LINK.transfer(msg.sender, total - s_expectedLinkBalance);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function pause() external onlyOwner {
    _pause();
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function unpause() external onlyOwner {
    _unpause();
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setPeerRegistryMigrationPermission(address peer, MigrationPermission permission) external onlyOwner {
    s_peerRegistryMigrationPermission[peer] = permission;
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    bytes calldata checkData
  ) external returns (uint256 id) {
    if (msg.sender != owner() && msg.sender != s_onChainConfig.registrar) revert OnlyCallableByOwnerOrRegistrar();

    id = uint256(keccak256(abi.encodePacked(blockhash(block.number - 1), address(this), s_nonce)));
    _createUpkeep(id, target, gasLimit, admin, 0, checkData, false);
    s_nonce++;
    emit UpkeepRegistered(id, gasLimit, admin);
    return id;
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function cancelUpkeep(uint256 id) external {
    Upkeep memory upkeep = s_upkeep[id];
    bool canceled = upkeep.maxValidBlocknumber != UINT32_MAX;
    bool isOwner = msg.sender == owner();

    if (canceled && !(isOwner && upkeep.maxValidBlocknumber > block.number)) revert CannotCancel();
    if (!isOwner && msg.sender != upkeep.admin) revert OnlyCallableByOwnerOrAdmin();

    uint256 height = block.number;
    if (!isOwner) {
      height = height + CANCELLATION_DELAY;
    }
    s_upkeep[id].maxValidBlocknumber = uint32(height);
    s_upkeepIDs.remove(id);

    // charge the cancellation fee if the minUpkeepSpend is not met
    uint96 minUpkeepSpend = s_onChainConfig.minUpkeepSpend;
    uint96 cancellationFee = 0;
    // cancellationFee is supposed to be min(max(minUpkeepSpend - amountSpent,0), amountLeft)
    if (upkeep.amountSpent < minUpkeepSpend) {
      cancellationFee = minUpkeepSpend - upkeep.amountSpent;
      if (cancellationFee > upkeep.balance) {
        cancellationFee = upkeep.balance;
      }
    }
    s_upkeep[id].balance = upkeep.balance - cancellationFee;
    s_ownerLinkBalance = s_ownerLinkBalance + cancellationFee;

    emit UpkeepCanceled(id, uint64(height));
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function addFunds(uint256 id, uint96 amount) external {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();

    s_upkeep[id].balance = upkeep.balance + amount;
    s_expectedLinkBalance = s_expectedLinkBalance + amount;
    LINK.transferFrom(msg.sender, address(this), amount);
    emit FundsAdded(id, msg.sender, amount);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function withdrawFunds(uint256 id, address to) external {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.admin != msg.sender) revert OnlyCallableByAdmin();
    if (upkeep.maxValidBlocknumber > block.number) revert UpkeepNotCanceled();

    uint96 amountToWithdraw = s_upkeep[id].balance;
    s_expectedLinkBalance = s_expectedLinkBalance - amountToWithdraw;
    s_upkeep[id].balance = 0;
    emit FundsWithdrawn(id, amountToWithdraw, to);

    LINK.transfer(to, amountToWithdraw);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external {
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_onChainConfig.maxPerformGas) revert GasLimitOutsideRange();
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    if (upkeep.admin != msg.sender) revert OnlyCallableByAdmin();

    s_upkeep[id].executeGas = gasLimit;

    emit UpkeepGasLimitSet(id, gasLimit);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function withdrawPayment(address from, address to) external {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    Transmitter memory keeper = s_transmitters[from];
    if (keeper.payee != msg.sender) revert OnlyCallableByPayee();

    s_transmitters[from].balance = 0;
    s_expectedLinkBalance = s_expectedLinkBalance - keeper.balance;
    emit PaymentWithdrawn(from, keeper.balance, to, msg.sender);

    LINK.transfer(to, keeper.balance);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function transferPayeeship(address keeper, address proposed) external {
    if (s_transmitters[keeper].payee != msg.sender) revert OnlyCallableByPayee();
    if (proposed == msg.sender) revert ValueNotChanged();

    if (s_proposedPayee[keeper] != proposed) {
      s_proposedPayee[keeper] = proposed;
      emit PayeeshipTransferRequested(keeper, msg.sender, proposed);
    }
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function acceptPayeeship(address keeper) external {
    if (s_proposedPayee[keeper] != msg.sender) revert OnlyCallableByProposedPayee();
    address past = s_transmitters[keeper].payee;
    s_transmitters[keeper].payee = msg.sender;
    s_proposedPayee[keeper] = ZERO_ADDRESS;

    emit PayeeshipTransferred(keeper, past, msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function transferUpkeepAdmin(uint256 id, address proposed) external {
    Upkeep memory upkeep = s_upkeep[id];
    requireAdminAndNotCancelled(upkeep);
    if (proposed == msg.sender) revert ValueNotChanged();
    if (proposed == ZERO_ADDRESS) revert InvalidRecipient();

    if (s_proposedAdmin[id] != proposed) {
      s_proposedAdmin[id] = proposed;
      emit UpkeepAdminTransferRequested(id, msg.sender, proposed);
    }
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function acceptUpkeepAdmin(uint256 id) external {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    if (s_proposedAdmin[id] != msg.sender) revert OnlyCallableByProposedAdmin();
    address past = upkeep.admin;
    s_upkeep[id].admin = msg.sender;
    s_proposedAdmin[id] = ZERO_ADDRESS;

    emit UpkeepAdminTransferred(id, past, msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function migrateUpkeeps(uint256[] calldata ids, address destination) external {
    // TODO (sc-49440): update upkeep format and handle migration from non OCR to OCR
    if (
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.OUTGOING &&
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    if (s_onChainConfig.transcoder == ZERO_ADDRESS) revert TranscoderNotSet();
    if (ids.length == 0) revert ArrayHasNoEntries();
    uint256 id;
    Upkeep memory upkeep;
    uint256 totalBalanceRemaining;
    bytes[] memory checkDatas = new bytes[](ids.length);
    Upkeep[] memory upkeeps = new Upkeep[](ids.length);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      id = ids[idx];
      upkeep = s_upkeep[id];
      requireAdminAndNotCancelled(upkeep);
      upkeeps[idx] = upkeep;
      checkDatas[idx] = s_checkData[id];
      totalBalanceRemaining = totalBalanceRemaining + upkeep.balance;
      delete s_upkeep[id];
      delete s_checkData[id];
      // nullify existing proposed admin change if an upkeep is being migrated
      delete s_proposedAdmin[id];
      s_upkeepIDs.remove(id);
      emit UpkeepMigrated(id, upkeep.balance, destination);
    }
    s_expectedLinkBalance = s_expectedLinkBalance - totalBalanceRemaining;
    bytes memory encodedUpkeeps = abi.encode(ids, upkeeps, checkDatas);
    MigratableKeeperRegistryInterfaceDev(destination).receiveUpkeeps(
      UpkeepTranscoderInterfaceDev(s_onChainConfig.transcoder).transcodeUpkeeps(
        UPKEEP_TRANSCODER_VERSION_BASE,
        MigratableKeeperRegistryInterfaceDev(destination).upkeepTranscoderVersion(),
        encodedUpkeeps
      )
    );
    LINK.transfer(destination, totalBalanceRemaining);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function receiveUpkeeps(bytes calldata encodedUpkeeps) external {
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
        checkDatas[idx],
        upkeeps[idx].paused
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
   * @param paused if this upkeep is paused
   */
  function _createUpkeep(
    uint256 id,
    address target,
    uint32 gasLimit,
    address admin,
    uint96 balance,
    bytes memory checkData,
    bool paused
  ) internal whenNotPaused {
    if (!target.isContract()) revert NotAContract();
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_onChainConfig.maxPerformGas) revert GasLimitOutsideRange();
    s_upkeep[id] = Upkeep({
      target: target,
      executeGas: gasLimit,
      balance: balance,
      admin: admin,
      maxValidBlocknumber: UINT32_MAX,
      lastPerformBlockNumber: 0,
      amountSpent: 0,
      paused: paused
    });
    s_expectedLinkBalance = s_expectedLinkBalance + balance;
    s_checkData[id] = checkData;
    s_upkeepIDs.add(id);
  }
}
