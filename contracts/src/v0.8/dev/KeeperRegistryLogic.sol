// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./KeeperRegistryBase.sol";

/**
 * @notice Logic contract, works in tandem with KeeperRegistry as a proxy
 */
contract KeeperRegistryLogic is KeeperRegistryBase {
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
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override onlyOwner {
    require(signers.length <= maxNumOracles, "too many oracles");
    require(signers.length == transmitters.length, "oracle length mismatch");
    require(3 * f < signers.length, "faulty-oracle f too high");
    require(0 < f, "f must be positive");
    require(onchainConfig.length == 0, "onchainConfig must be empty");

    // remove any old signer/transmitter addresses
    uint256 oldLength = s_signersList.length;
    for (uint256 i = 0; i < oldLength; i++) {
      address signer = s_signersList[i];
      address transmitter = s_transmittersList[i];
      delete s_signers[signer];
      delete s_transmitters[transmitter];
    }
    delete s_signersList;
    delete s_transmittersList;

    // add new signer/transmitter addresses
    for (uint256 i = 0; i < signers.length; i++) {
      require(!s_signers[signers[i]].active, "repeated signer address");
      s_signers[signers[i]] = Signer({active: true, index: uint8(i)});
      require(!s_transmitters[transmitters[i]].active, "repeated transmitter address");
      s_transmitters[transmitters[i]] = Transmitter({active: true, index: uint8(i), paymentJuels: 0});
    }
    s_signersList = signers;
    s_transmittersList = transmitters;

    s_hotVars.latestEpochAndRound = 0;
    s_hotVars.f = f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    s_latestConfigDigest = _configDigestFromConfigData(
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
      s_latestConfigDigest,
      s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    // TODO: understand if this is needed
    /*uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    for (uint256 i = 0; i < signers.length; i++) {
      s_rewardFromAggregatorRoundId[i] = latestAggregatorRoundId;
    }*/
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
      bytes32 configDigest
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
   * @dev Unimplemented on logic contract, implementation lives on KeeperRegistry main contract
   */
  function typeAndVersion() external pure override returns (string memory) {}
}
