// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

contract KeeperRegistry1_1Mock {
  event ConfigSet(
    uint32 paymentPremiumPPB,
    uint24 blockCountPerTurn,
    uint32 checkGasLimit,
    uint24 stalenessSeconds,
    uint16 gasCeilingMultiplier,
    uint256 fallbackGasPrice,
    uint256 fallbackLinkPrice
  );
  event FlatFeeSet(uint32 flatFeeMicroLink);
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);
  event KeepersUpdated(address[] keepers, address[] payees);
  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);
  event Paused(address account);
  event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to);
  event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to);
  event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee);
  event RegistrarChanged(address indexed from, address indexed to);
  event Unpaused(address account);
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    address indexed from,
    uint96 payment,
    bytes performData
  );
  event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin);

  function emitConfigSet(
    uint32 paymentPremiumPPB,
    uint24 blockCountPerTurn,
    uint32 checkGasLimit,
    uint24 stalenessSeconds,
    uint16 gasCeilingMultiplier,
    uint256 fallbackGasPrice,
    uint256 fallbackLinkPrice
  ) public {
    emit ConfigSet(
      paymentPremiumPPB,
      blockCountPerTurn,
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      fallbackGasPrice,
      fallbackLinkPrice
    );
  }

  function emitFlatFeeSet(uint32 flatFeeMicroLink) public {
    emit FlatFeeSet(flatFeeMicroLink);
  }

  function emitFundsAdded(uint256 id, address from, uint96 amount) public {
    emit FundsAdded(id, from, amount);
  }

  function emitFundsWithdrawn(uint256 id, uint256 amount, address to) public {
    emit FundsWithdrawn(id, amount, to);
  }

  function emitKeepersUpdated(address[] memory keepers, address[] memory payees) public {
    emit KeepersUpdated(keepers, payees);
  }

  function emitOwnershipTransferRequested(address from, address to) public {
    emit OwnershipTransferRequested(from, to);
  }

  function emitOwnershipTransferred(address from, address to) public {
    emit OwnershipTransferred(from, to);
  }

  function emitPaused(address account) public {
    emit Paused(account);
  }

  function emitPayeeshipTransferRequested(address keeper, address from, address to) public {
    emit PayeeshipTransferRequested(keeper, from, to);
  }

  function emitPayeeshipTransferred(address keeper, address from, address to) public {
    emit PayeeshipTransferred(keeper, from, to);
  }

  function emitPaymentWithdrawn(address keeper, uint256 amount, address to, address payee) public {
    emit PaymentWithdrawn(keeper, amount, to, payee);
  }

  function emitRegistrarChanged(address from, address to) public {
    emit RegistrarChanged(from, to);
  }

  function emitUnpaused(address account) public {
    emit Unpaused(account);
  }

  function emitUpkeepCanceled(uint256 id, uint64 atBlockHeight) public {
    emit UpkeepCanceled(id, atBlockHeight);
  }

  function emitUpkeepPerformed(
    uint256 id,
    bool success,
    address from,
    uint96 payment,
    bytes memory performData
  ) public {
    emit UpkeepPerformed(id, success, from, payment, performData);
  }

  function emitUpkeepRegistered(uint256 id, uint32 executeGas, address admin) public {
    emit UpkeepRegistered(id, executeGas, admin);
  }

  uint256 private s_upkeepCount;

  // Function to set the current number of registered upkeeps
  function setUpkeepCount(uint256 _upkeepCount) external {
    s_upkeepCount = _upkeepCount;
  }

  // Function to get the current number of registered upkeeps
  function getUpkeepCount() external view returns (uint256) {
    return s_upkeepCount;
  }

  uint256[] private s_canceledUpkeepList;

  // Function to set the current number of canceled upkeeps
  function setCanceledUpkeepList(uint256[] memory _canceledUpkeepList) external {
    s_canceledUpkeepList = _canceledUpkeepList;
  }

  // Function to set the current number of canceled upkeeps
  function getCanceledUpkeepList() external view returns (uint256[] memory) {
    return s_canceledUpkeepList;
  }

  address[] private s_keeperList;

  // Function to set the keeper list for testing purposes
  function setKeeperList(address[] memory _keepers) external {
    s_keeperList = _keepers;
  }

  // Function to get the keeper list
  function getKeeperList() external view returns (address[] memory) {
    return s_keeperList;
  }

  struct Config {
    uint32 paymentPremiumPPB;
    uint32 flatFeeMicroLink; // min 0.000001 LINK, max 4294 LINK
    uint24 blockCountPerTurn;
    uint32 checkGasLimit;
    uint24 stalenessSeconds;
    uint16 gasCeilingMultiplier;
  }

  Config private s_config;
  uint256 private s_fallbackGasPrice;
  uint256 private s_fallbackLinkPrice;

  // Function to set the configuration for testing purposes
  function setConfig(
    uint32 _paymentPremiumPPB,
    uint32 _flatFeeMicroLink,
    uint24 _blockCountPerTurn,
    uint32 _checkGasLimit,
    uint24 _stalenessSeconds,
    uint16 _gasCeilingMultiplier,
    uint256 _fallbackGasPrice,
    uint256 _fallbackLinkPrice
  ) external {
    s_config.paymentPremiumPPB = _paymentPremiumPPB;
    s_config.flatFeeMicroLink = _flatFeeMicroLink;
    s_config.blockCountPerTurn = _blockCountPerTurn;
    s_config.checkGasLimit = _checkGasLimit;
    s_config.stalenessSeconds = _stalenessSeconds;
    s_config.gasCeilingMultiplier = _gasCeilingMultiplier;
    s_fallbackGasPrice = _fallbackGasPrice;
    s_fallbackLinkPrice = _fallbackLinkPrice;
  }

  // Function to get the configuration
  function getConfig()
    external
    view
    returns (
      uint32 paymentPremiumPPB,
      uint24 blockCountPerTurn,
      uint32 checkGasLimit,
      uint24 stalenessSeconds,
      uint16 gasCeilingMultiplier,
      uint256 fallbackGasPrice,
      uint256 fallbackLinkPrice
    )
  {
    return (
      s_config.paymentPremiumPPB,
      s_config.blockCountPerTurn,
      s_config.checkGasLimit,
      s_config.stalenessSeconds,
      s_config.gasCeilingMultiplier,
      s_fallbackGasPrice,
      s_fallbackLinkPrice
    );
  }

  struct Upkeep {
    address target;
    uint32 executeGas;
    uint96 balance;
    address admin;
    uint64 maxValidBlocknumber;
    address lastKeeper;
  }

  mapping(uint256 => Upkeep) private s_upkeep;
  mapping(uint256 => bytes) private s_checkData;

  // Function to set the upkeep and checkData for testing purposes
  function setUpkeep(
    uint256 id,
    address _target,
    uint32 _executeGas,
    uint96 _balance,
    address _admin,
    uint64 _maxValidBlocknumber,
    address _lastKeeper,
    bytes memory _checkData
  ) external {
    Upkeep memory upkeep = Upkeep({
      target: _target,
      executeGas: _executeGas,
      balance: _balance,
      admin: _admin,
      maxValidBlocknumber: _maxValidBlocknumber,
      lastKeeper: _lastKeeper
    });

    s_upkeep[id] = upkeep;
    s_checkData[id] = _checkData;
  }

  // Function to get the upkeep and checkData
  function getUpkeep(
    uint256 id
  )
    external
    view
    returns (
      address target,
      uint32 executeGas,
      bytes memory checkData,
      uint96 balance,
      address lastKeeper,
      address admin,
      uint64 maxValidBlocknumber
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
      reg.maxValidBlocknumber
    );
  }

  mapping(uint256 => uint96) private s_minBalances;

  // Function to set the minimum balance for a specific upkeep id
  function setMinBalance(uint256 id, uint96 minBalance) external {
    s_minBalances[id] = minBalance;
  }

  // Function to get the minimum balance for a specific upkeep id
  function getMinBalanceForUpkeep(uint256 id) external view returns (uint96) {
    return s_minBalances[id];
  }

  struct UpkeepData {
    bytes performData;
    uint256 maxLinkPayment;
    uint256 gasLimit;
    uint256 adjustedGasWei;
    uint256 linkEth;
  }

  mapping(uint256 => UpkeepData) private s_upkeepData;

  // Function to set mock data for the checkUpkeep function
  function setCheckUpkeepData(
    uint256 id,
    bytes memory performData,
    uint256 maxLinkPayment,
    uint256 gasLimit,
    uint256 adjustedGasWei,
    uint256 linkEth
  ) external {
    s_upkeepData[id] = UpkeepData({
      performData: performData,
      maxLinkPayment: maxLinkPayment,
      gasLimit: gasLimit,
      adjustedGasWei: adjustedGasWei,
      linkEth: linkEth
    });
  }

  // Mock checkUpkeep function
  function checkUpkeep(
    uint256 id,
    address from
  )
    external
    view
    returns (
      bytes memory performData,
      uint256 maxLinkPayment,
      uint256 gasLimit,
      uint256 adjustedGasWei,
      uint256 linkEth
    )
  {
    UpkeepData storage data = s_upkeepData[id];
    return (data.performData, data.maxLinkPayment, data.gasLimit, data.adjustedGasWei, data.linkEth);
  }

  mapping(uint256 => bool) private s_upkeepSuccess;

  // Function to set mock return data for the performUpkeep function
  function setPerformUpkeepSuccess(uint256 id, bool success) external {
    s_upkeepSuccess[id] = success;
  }

  // Mock performUpkeep function
  function performUpkeep(uint256 id, bytes calldata performData) external returns (bool success) {
    return s_upkeepSuccess[id];
  }
}
