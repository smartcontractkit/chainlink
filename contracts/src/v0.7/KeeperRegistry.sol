// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./interfaces/AggregatorV3Interface.sol";
import "./interfaces/LinkTokenInterface.sol";
import "./interfaces/KeeperCompatibleInterface.sol";
import "./interfaces/KeeperRegistryInterface.sol";
import "./vendor/SafeMathChainlink.sol";
import "./vendor/Address.sol";
import "./vendor/Pausable.sol";
import "./vendor/ReentrancyGuard.sol";
import "./vendor/SignedSafeMath.sol";
import "./vendor/SafeMath96.sol";
import "./KeeperBase.sol";
import "./ConfirmedOwner.sol";

/**
  * @notice Registry for adding work for Chainlink Keepers to perform on client
  * contracts. Clients must support the Upkeep interface.
  */
contract KeeperRegistry is
  ConfirmedOwner,
  KeeperBase,
  ReentrancyGuard,
  Pausable,
  KeeperRegistryExecutableInterface
{
  using Address for address;
  using SafeMathChainlink for uint256;
  using SafeMath96 for uint96;
  using SignedSafeMath for int256;

  address constant private ZERO_ADDRESS = address(0);
  address constant private IGNORE_ADDRESS = 0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF;
  bytes4 constant private CHECK_SELECTOR = KeeperCompatibleInterface.checkUpkeep.selector;
  bytes4 constant private PERFORM_SELECTOR = KeeperCompatibleInterface.performUpkeep.selector;
  uint256 constant private CALL_GAS_MAX = 5_000_000;
  uint256 constant private CALL_GAS_MIN = 2_300;
  uint256 constant private CANCELATION_DELAY = 50;
  uint256 constant private CUSHION = 5_000;
  uint256 constant private REGISTRY_GAS_OVERHEAD = 80_000;
  uint256 constant private PPB_BASE = 1_000_000_000;
  uint64 constant private UINT64_MAX = 2**64 - 1;
  uint96 constant private LINK_TOTAL_SUPPLY = 1e27;

  uint256 private s_upkeepCount;
  uint256[] private s_canceledUpkeepList;
  address[] private s_keeperList;
  mapping(uint256 => Upkeep) private s_upkeep;
  mapping(address => KeeperInfo) private s_keeperInfo;
  mapping(address => address) private s_proposedPayee;
  mapping(uint256 => bytes) private s_checkData;
  Config private s_config;
  uint256 private s_fallbackGasPrice;  // not in config object for gas savings
  uint256 private s_fallbackLinkPrice; // not in config object for gas savings
  uint256 private s_expectedLinkBalance;

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_ETH_FEED;
  AggregatorV3Interface public immutable FAST_GAS_FEED;

  address private s_registrar;

  struct Upkeep {
    address target;
    uint32 executeGas;
    uint96 balance;
    address admin;
    uint64 maxValidBlocknumber;
    address lastKeeper;
  }

  struct KeeperInfo {
    address payee;
    uint96 balance;
    bool active;
  }

  struct Config {
    uint32 paymentPremiumPPB;
    uint24 blockCountPerTurn;
    uint32 checkGasLimit;
    uint24 stalenessSeconds;
    uint16 gasCeilingMultiplier;
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

  event UpkeepRegistered(
    uint256 indexed id,
    uint32 executeGas,
    address admin
  );
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    address indexed from,
    uint96 payment,
    bytes performData
  );
  event UpkeepCanceled(
    uint256 indexed id,
    uint64 indexed atBlockHeight
  );
  event FundsAdded(
    uint256 indexed id,
    address indexed from,
    uint96 amount
  );
  event FundsWithdrawn(
    uint256 indexed id,
    uint256 amount,
    address to
  );
  event ConfigSet(
    uint32 paymentPremiumPPB,
    uint24 blockCountPerTurn,
    uint32 checkGasLimit,
    uint24 stalenessSeconds,
    uint16 gasCeilingMultiplier,
    uint256 fallbackGasPrice,
    uint256 fallbackLinkPrice
  );
  event KeepersUpdated(
    address[] keepers,
    address[] payees
  );
  event PaymentWithdrawn(
    address indexed keeper,
    uint256 indexed amount,
    address indexed to,
    address payee
  );
  event PayeeshipTransferRequested(
    address indexed keeper,
    address indexed from,
    address indexed to
  );
  event PayeeshipTransferred(
    address indexed keeper,
    address indexed from,
    address indexed to
  );
  event RegistrarChanged(
    address indexed from,
    address indexed to
  );
  /**
   * @param link address of the LINK Token
   * @param linkEthFeed address of the LINK/ETH price feed
   * @param fastGasFeed address of the Fast Gas price feed
   * @param paymentPremiumPPB payment premium rate oracles receive on top of
   * being reimbursed for gas, measured in parts per billion
   * @param blockCountPerTurn number of blocks each oracle has during their turn to
   * perform upkeep before it will be the next keeper's turn to submit
   * @param checkGasLimit gas limit when checking for upkeep
   * @param stalenessSeconds number of seconds that is allowed for feed data to
   * be stale before switching to the fallback pricing
   * @param gasCeilingMultiplier multiplier to apply to the fast gas feed price
   * when calculating the payment ceiling for keepers
   * @param fallbackGasPrice gas price used if the gas price feed is stale
   * @param fallbackLinkPrice LINK price used if the LINK price feed is stale
   */
  constructor(
    address link,
    address linkEthFeed,
    address fastGasFeed,
    uint32 paymentPremiumPPB,
    uint24 blockCountPerTurn,
    uint32 checkGasLimit,
    uint24 stalenessSeconds,
    uint16 gasCeilingMultiplier,
    uint256 fallbackGasPrice,
    uint256 fallbackLinkPrice
  )
    ConfirmedOwner(msg.sender)
  {
    LINK = LinkTokenInterface(link);
    LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
    FAST_GAS_FEED = AggregatorV3Interface(fastGasFeed);

    setConfig(
      paymentPremiumPPB,
      blockCountPerTurn,
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      fallbackGasPrice,
      fallbackLinkPrice
    );
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
  )
    external
    override
    onlyOwnerOrRegistrar()
    returns (
      uint256 id
    )
  {
    require(target.isContract(), "target is not a contract");
    require(gasLimit >= CALL_GAS_MIN, "min gas is 2300");
    require(gasLimit <= CALL_GAS_MAX, "max gas is 5000000");

    id = s_upkeepCount;
    s_upkeep[id] = Upkeep({
      target: target,
      executeGas: gasLimit,
      balance: 0,
      admin: admin,
      maxValidBlocknumber: UINT64_MAX,
      lastKeeper: address(0)
    });
    s_checkData[id] = checkData;
    s_upkeepCount++;

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
  function checkUpkeep(
    uint256 id,
    address from
  )
    external
    override
    whenNotPaused()
    cannotExecute()
    returns (
      bytes memory performData,
      uint256 maxLinkPayment,
      uint256 gasLimit,
      uint256 adjustedGasWei,
      uint256 linkEth
    )
  {
    bytes memory callData = abi.encodeWithSelector(CHECK_SELECTOR, s_checkData[id]);
    (
      bool success,
      bytes memory result
    ) = s_upkeep[id].target.call{gas: s_config.checkGasLimit}(callData);

    if (!success) {
      string memory upkeepRevertReason = getRevertMsg(result);
      string memory reason = string(abi.encodePacked("call to check target failed: ", upkeepRevertReason));
      revert(reason);
    }

    (
      success,
      performData
    ) = abi.decode(result, (bool, bytes));
    require(success, "upkeep not needed");

    PerformParams memory params = generatePerformParams(from, id, performData, false);
    success = performUpkeepWithParams(params);
    require(success, "call to perform upkeep failed");

    return (performData, params.maxLinkPayment, params.gasLimit, params.adjustedGasWei, params.linkEth);
  }


  /**
   * @notice executes the upkeep with the perform data returned from
   * checkUpkeep, validates the keeper's permissions, and pays the keeper.
   * @param id identifier of the upkeep to execute the data with.
   * @param performData calldata parameter to be passed to the target upkeep.
   */
  function performUpkeep(
    uint256 id,
    bytes calldata performData
  )
    external
    override
    returns (
      bool success
    )
  {

    return performUpkeepWithParams(generatePerformParams(
      msg.sender,
      id,
      performData,
      true
    ));
  }

  /**
   * @notice prevent an upkeep from being performed in the future
   * @param id upkeep to be canceled
   */
  function cancelUpkeep(
    uint256 id
  )
    external
    override
  {
    uint64 maxValid = s_upkeep[id].maxValidBlocknumber;
    bool notCanceled = maxValid == UINT64_MAX;
    bool isOwner = msg.sender == owner();
    require(notCanceled || (isOwner && maxValid > block.number), "too late to cancel upkeep");
    require(isOwner || msg.sender == s_upkeep[id].admin, "only owner or admin");

    uint256 height = block.number;
    if (!isOwner) {
      height = height.add(CANCELATION_DELAY);
    }
    s_upkeep[id].maxValidBlocknumber = uint64(height);
    if (notCanceled) {
      s_canceledUpkeepList.push(id);
    }

    emit UpkeepCanceled(id, uint64(height));
  }

  /**
   * @notice adds LINK funding for an upkeep by tranferring from the sender's
   * LINK balance
   * @param id upkeep to fund
   * @param amount number of LINK to transfer
   */
  function addFunds(
    uint256 id,
    uint96 amount
  )
    external
    override
  {
    require(s_upkeep[id].maxValidBlocknumber == UINT64_MAX, "upkeep must be active");
    s_upkeep[id].balance = s_upkeep[id].balance.add(amount);
    s_expectedLinkBalance = s_expectedLinkBalance.add(amount);
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
  )
    external
  {
    require(msg.sender == address(LINK), "only callable through LINK");
    require(data.length == 32, "data must be 32 bytes");
    uint256 id = abi.decode(data, (uint256));
    require(s_upkeep[id].maxValidBlocknumber == UINT64_MAX, "upkeep must be active");

    s_upkeep[id].balance = s_upkeep[id].balance.add(uint96(amount));
    s_expectedLinkBalance = s_expectedLinkBalance.add(amount);

    emit FundsAdded(id, sender, uint96(amount));
  }

  /**
   * @notice removes funding from a canceled upkeep
   * @param id upkeep to withdraw funds from
   * @param to destination address for sending remaining funds
   */
  function withdrawFunds(
    uint256 id,
    address to
  )
    external
    validateRecipient(to)
  {
    require(s_upkeep[id].admin == msg.sender, "only callable by admin");
    require(s_upkeep[id].maxValidBlocknumber <= block.number, "upkeep must be canceled");

    uint256 amount = s_upkeep[id].balance;
    s_upkeep[id].balance = 0;
    s_expectedLinkBalance = s_expectedLinkBalance.sub(amount);
    emit FundsWithdrawn(id, amount, to);

    LINK.transfer(to, amount);
  }

  /**
   * @notice recovers LINK funds improperly transfered to the registry
   * @dev In principle this function’s execution cost could exceed block
   * gaslimit. However, in our anticipated deployment, the number of upkeeps and
   * keepers will be low enough to avoid this problem.
   */
  function recoverFunds()
    external
    onlyOwner()
  {
    uint256 total = LINK.balanceOf(address(this));
    LINK.transfer(msg.sender, total.sub(s_expectedLinkBalance));
  }

  /**
   * @notice withdraws a keeper's payment, callable only by the keeper's payee
   * @param from keeper address
   * @param to address to send the payment to
   */
  function withdrawPayment(
    address from,
    address to
  )
    external
    validateRecipient(to)
  {
    KeeperInfo memory keeper = s_keeperInfo[from];
    require(keeper.payee == msg.sender, "only callable by payee");

    s_keeperInfo[from].balance = 0;
    s_expectedLinkBalance = s_expectedLinkBalance.sub(keeper.balance);
    emit PaymentWithdrawn(from, keeper.balance, to, msg.sender);

    LINK.transfer(to, keeper.balance);
  }

  /**
   * @notice proposes the safe transfer of a keeper's payee to another address
   * @param keeper address of the keeper to transfer payee role
   * @param proposed address to nominate for next payeeship
   */
  function transferPayeeship(
    address keeper,
    address proposed
  )
    external
  {
    require(s_keeperInfo[keeper].payee == msg.sender, "only callable by payee");
    require(proposed != msg.sender, "cannot transfer to self");

    if (s_proposedPayee[keeper] != proposed) {
      s_proposedPayee[keeper] = proposed;
      emit PayeeshipTransferRequested(keeper, msg.sender, proposed);
    }
  }

  /**
   * @notice accepts the safe transfer of payee role for a keeper
   * @param keeper address to accept the payee role for
   */
  function acceptPayeeship(
    address keeper
  )
    external
  {
    require(s_proposedPayee[keeper] == msg.sender, "only callable by proposed payee");
    address past = s_keeperInfo[keeper].payee;
    s_keeperInfo[keeper].payee = msg.sender;
    s_proposedPayee[keeper] = ZERO_ADDRESS;

    emit PayeeshipTransferred(keeper, past, msg.sender);
  }

  /**
   * @notice signals to keepers that they should not perform upkeeps until the
   * contract has been unpaused
   */
  function pause()
    external
    onlyOwner()
  {
    _pause();
  }

  /**
   * @notice signals to keepers that they can perform upkeeps once again after
   * having been paused
   */
  function unpause()
    external
    onlyOwner()
  {
    _unpause();
  }


  // SETTERS

  /**
   * @notice updates the configuration of the registry
   * @param paymentPremiumPPB payment premium rate oracles receive on top of
   * being reimbursed for gas, measured in parts per billion
   * @param blockCountPerTurn number of blocks an oracle should wait before
   * checking for upkeep
   * @param checkGasLimit gas limit when checking for upkeep
   * @param stalenessSeconds number of seconds that is allowed for feed data to
   * be stale before switching to the fallback pricing
   * @param fallbackGasPrice gas price used if the gas price feed is stale
   * @param fallbackLinkPrice LINK price used if the LINK price feed is stale
   */
  function setConfig(
    uint32 paymentPremiumPPB,
    uint24 blockCountPerTurn,
    uint32 checkGasLimit,
    uint24 stalenessSeconds,
    uint16 gasCeilingMultiplier,
    uint256 fallbackGasPrice,
    uint256 fallbackLinkPrice
  )
    onlyOwner()
    public
  {
    s_config = Config({
      paymentPremiumPPB: paymentPremiumPPB,
      blockCountPerTurn: blockCountPerTurn,
      checkGasLimit: checkGasLimit,
      stalenessSeconds: stalenessSeconds,
      gasCeilingMultiplier: gasCeilingMultiplier
    });
    s_fallbackGasPrice = fallbackGasPrice;
    s_fallbackLinkPrice = fallbackLinkPrice;

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

  /**
   * @notice update the list of keepers allowed to perform upkeep
   * @param keepers list of addresses allowed to perform upkeep
   * @param payees addreses corresponding to keepers who are allowed to
   * move payments which have been accrued
   */
  function setKeepers(
    address[] calldata keepers,
    address[] calldata payees
  )
    external
    onlyOwner()
  {
    require(keepers.length == payees.length, "address lists not the same length");
    require(keepers.length >= 2, "not enough keepers");
    for (uint256 i = 0; i < s_keeperList.length; i++) {
      address keeper = s_keeperList[i];
      s_keeperInfo[keeper].active = false;
    }
    for (uint256 i = 0; i < keepers.length; i++) {
      address keeper = keepers[i];
      KeeperInfo storage s_keeper = s_keeperInfo[keeper];
      address oldPayee = s_keeper.payee;
      address newPayee = payees[i];
      require(newPayee != address(0), "cannot set payee to the zero address");
      require(oldPayee == ZERO_ADDRESS || oldPayee == newPayee || newPayee == IGNORE_ADDRESS, "cannot change payee");
      require(!s_keeper.active, "cannot add keeper twice");
      s_keeper.active = true;
      if (newPayee != IGNORE_ADDRESS) {
        s_keeper.payee = newPayee;
      }
    }
    s_keeperList = keepers;
    emit KeepersUpdated(keepers, payees);
  }

  /**
   * @notice update registrar
   * @param registrar new registrar
   */
  function setRegistrar(
    address registrar
  )
    external
    onlyOwnerOrRegistrar()
  {
    address previous = s_registrar;
    require(registrar != previous, "Same registrar");
    s_registrar = registrar;
    emit RegistrarChanged(previous, registrar);
  }

  // GETTERS

  /**
   * @notice read all of the details about an upkeep
   */
  function getUpkeep(
    uint256 id
  )
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

  /**
   * @notice read the total number of upkeep's registered
   */
  function getUpkeepCount()
    external
    view
    override
    returns (
      uint256
    )
  {
    return s_upkeepCount;
  }

  /**
   * @notice read the current list canceled upkeep IDs
   */
  function getCanceledUpkeepList()
    external
    view
    override
    returns (
      uint256[] memory
    )
  {
    return s_canceledUpkeepList;
  }

  /**
   * @notice read the current list of addresses allowed to perform upkeep
   */
  function getKeeperList()
    external
    view
    override
    returns (
      address[] memory
    )
  {
    return s_keeperList;
  }

 /**
   * @notice read the current registrar
   */
  function getRegistrar()
    external
    view
    returns (
      address
    )
  {
    return s_registrar;
  }

  /**
   * @notice read the current info about any keeper address
   */
  function getKeeperInfo(
    address query
  )
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
   * @notice read the current configuration of the registry
   */
  function getConfig()
    external
    view
    override
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
    Config memory config = s_config;
    return (
      config.paymentPremiumPPB,
      config.blockCountPerTurn,
      config.checkGasLimit,
      config.stalenessSeconds,
      config.gasCeilingMultiplier,
      s_fallbackGasPrice,
      s_fallbackLinkPrice
    );
  }

  /**
   * @notice calculates the minimum balance required for an upkeep to remain eligible
   */
  function getMinBalanceForUpkeep(
    uint256 id
  )
    external
    view
    returns (
      uint96 minBalance
    )
  {
    return getMaxPaymentForGas(s_upkeep[id].executeGas);
  }

  /**
   * @notice calculates the maximum payment for a given gas limit
   */
  function getMaxPaymentForGas(
    uint256 gasLimit
  )
    public
    view
    returns (
      uint96 maxPayment
    )
  {
    (uint256 gasWei, uint256 linkEth) = getFeedData();
    uint256 adjustedGasWei = adjustGasPrice(gasWei, false);
    return calculatePaymentAmount(gasLimit, adjustedGasWei, linkEth);
  }


  // PRIVATE

  /**
   * @dev retrieves feed data for fast gas/eth and link/eth prices. if the feed
   * data is stale it uses the configured fallback price. Once a price is picked
   * for gas it takes the min of gas price in the transaction or the fast gas
   * price in order to reduce costs for the upkeep clients.
   */
  function getFeedData()
    private
    view
    returns (
      uint256 gasWei,
      uint256 linkEth
    )
  {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 feedValue;
    (,feedValue,,timestamp,) = FAST_GAS_FEED.latestRoundData();
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp || feedValue <=0) {
      gasWei = s_fallbackGasPrice;
    } else {
      gasWei = uint256(feedValue);
    }
    (,feedValue,,timestamp,) = LINK_ETH_FEED.latestRoundData();
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp || feedValue <=0) {
      linkEth = s_fallbackLinkPrice;
    } else {
      linkEth = uint256(feedValue);
    }
    return (gasWei, linkEth);
  }

  /**
   * @dev calculates LINK paid for gas spent plus a configure premium percentage
   */
  function calculatePaymentAmount(
    uint256 gasLimit,
    uint256 gasWei,
    uint256 linkEth
  )
    private
    view
    returns (
      uint96 payment
    )
  {
    uint256 weiForGas = gasWei.mul(gasLimit.add(REGISTRY_GAS_OVERHEAD));
    uint256 premium = PPB_BASE.add(s_config.paymentPremiumPPB);
    uint256 total = weiForGas.mul(1e9).mul(premium).div(linkEth);
    require(total <= LINK_TOTAL_SUPPLY, "payment greater than all LINK");
    return uint96(total); // LINK_TOTAL_SUPPLY < UINT96_MAX
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available
   */
  function callWithExactGas(
    uint256 gasAmount,
    address target,
    bytes memory data
  )
    private
    returns (
      bool success
    )
  {
    assembly{
      let g := gas()
      // Compute g -= CUSHION and check for underflow
      if lt(g, CUSHION) { revert(0, 0) }
      g := sub(g, CUSHION)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) { revert(0, 0) }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) { revert(0, 0) }
      // call and return whether we succeeded. ignore return data
      success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
    }
    return success;
  }

  /**
   * @dev calls the Upkeep target with the performData param passed in by the
   * keeper and the exact gas required by the Upkeep
   */
  function performUpkeepWithParams(
    PerformParams memory params
  )
    private
    nonReentrant()
    validUpkeep(params.id)
    returns (
      bool success
    )
  {
    require(s_keeperInfo[params.from].active, "only active keepers");
    Upkeep memory upkeep = s_upkeep[params.id];
    require(upkeep.balance >= params.maxLinkPayment, "insufficient funds");
    require(upkeep.lastKeeper != params.from, "keepers must take turns");

    uint256  gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(PERFORM_SELECTOR, params.performData);
    success = callWithExactGas(params.gasLimit, upkeep.target, callData);
    gasUsed = gasUsed - gasleft();

    uint96 payment = calculatePaymentAmount(gasUsed, params.adjustedGasWei, params.linkEth);
    upkeep.balance = upkeep.balance.sub(payment);
    upkeep.lastKeeper = params.from;
    s_upkeep[params.id] = upkeep;
    uint96 newBalance = s_keeperInfo[params.from].balance.add(payment);
    s_keeperInfo[params.from].balance = newBalance;

    emit UpkeepPerformed(
      params.id,
      success,
      params.from,
      payment,
      params.performData
    );
    return success;
  }

  /**
   * @dev ensures a upkeep is valid
   */
  function validateUpkeep(
    uint256 id
  )
    private
    view
  {
    require(s_upkeep[id].maxValidBlocknumber > block.number, "invalid upkeep id");
  }

  /**
   * @dev adjusts the gas price to min(ceiling, tx.gasprice) or just uses the ceiling if tx.gasprice is disabled
   */
  function adjustGasPrice(
    uint256 gasWei,
    bool useTxGasPrice
  )
    private
    view
    returns(uint256 adjustedPrice)
  {
    adjustedPrice = gasWei.mul(s_config.gasCeilingMultiplier);
    if (useTxGasPrice && tx.gasprice < adjustedPrice) {
      adjustedPrice = tx.gasprice;
    }
  }

  /**
   * @dev generates a PerformParams struct for use in performUpkeepWithParams()
   */
  function generatePerformParams(
    address from,
    uint256 id,
    bytes memory performData,
    bool useTxGasPrice
  )
    private
    view
    returns(PerformParams memory)
  {
    uint256 gasLimit = s_upkeep[id].executeGas;
    (uint256 gasWei, uint256 linkEth) = getFeedData();
    uint256 adjustedGasWei = adjustGasPrice(gasWei, useTxGasPrice);
    uint96 maxLinkPayment = calculatePaymentAmount(gasLimit, adjustedGasWei, linkEth);

    return PerformParams({
      from: from,
      id: id,
      performData: performData,
      maxLinkPayment: maxLinkPayment,
      gasLimit: gasLimit,
      adjustedGasWei: adjustedGasWei,
      linkEth: linkEth
    });
  }

  /**
   * @dev extracts a revert reason from a call result payload
   */
  function getRevertMsg(bytes memory _payload) private pure returns (string memory) {
    if (_payload.length < 68) return 'transaction reverted silently';
    assembly {
        _payload := add(_payload, 0x04)
    }
    return abi.decode(_payload, (string));
  }

  // MODIFIERS

  /**
   * @dev ensures a upkeep is valid
   */
  modifier validUpkeep(
    uint256 id
  ) {
    validateUpkeep(id);
    _;
  }

  /**
   * @dev ensures that burns don't accidentally happen by sending to the zero
   * address
   */
  modifier validateRecipient(
    address to
  ) {
    require(to != address(0), "cannot send to zero address");
    _;
  }

    /**
   * @dev Reverts if called by anyone other than the contract owner or registrar.
   */
  modifier onlyOwnerOrRegistrar() {
    require(msg.sender == owner() || msg.sender == s_registrar, "Only callable by owner or registrar");
    _;
  }

}
