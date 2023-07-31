// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @notice config of the registry
 * @dev only used in params and return values
 * @member paymentPremiumPPB payment premium rate oracles receive on top of
 * being reimbursed for gas, measured in parts per billion
 * @member flatFeeMicroLink flat fee paid to oracles for performing upkeeps,
 * priced in MicroLink; can be used in conjunction with or independently of
 * paymentPremiumPPB
 * @member blockCountPerTurn number of blocks each oracle has during their turn to
 * perform upkeep before it will be the next keeper's turn to submit
 * @member checkGasLimit gas limit when checking for upkeep
 * @member stalenessSeconds number of seconds that is allowed for feed data to
 * be stale before switching to the fallback pricing
 * @member gasCeilingMultiplier multiplier to apply to the fast gas feed price
 * when calculating the payment ceiling for keepers
 * @member minUpkeepSpend minimum LINK that an upkeep must spend before cancelling
 * @member maxPerformGas max executeGas allowed for an upkeep on this registry
 * @member fallbackGasPrice gas price used if the gas price feed is stale
 * @member fallbackLinkPrice LINK price used if the LINK price feed is stale
 * @member transcoder address of the transcoder contract
 * @member registrar address of the registrar contract
 */
struct Config {
  uint32 paymentPremiumPPB;
  uint32 flatFeeMicroLink; // min 0.000001 LINK, max 4294 LINK
  uint24 blockCountPerTurn;
  uint32 checkGasLimit;
  uint24 stalenessSeconds;
  uint16 gasCeilingMultiplier;
  uint96 minUpkeepSpend;
  uint32 maxPerformGas;
  uint256 fallbackGasPrice;
  uint256 fallbackLinkPrice;
  address transcoder;
  address registrar;
}

/**
 * @notice state of the registry
 * @dev only used in params and return values
 * @member nonce used for ID generation
 * @member ownerLinkBalance withdrawable balance of LINK by contract owner
 * @member expectedLinkBalance the expected balance of LINK of the registry
 * @member numUpkeeps total number of upkeeps on the registry
 */
struct State {
  uint32 nonce;
  uint96 ownerLinkBalance;
  uint256 expectedLinkBalance;
  uint256 numUpkeeps;
}

interface AutomationRegistryBaseInterface {
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    bytes calldata checkData
  ) external returns (uint256 id);

  function performUpkeep(uint256 id, bytes calldata performData) external returns (bool success);

  function cancelUpkeep(uint256 id) external;

  function addFunds(uint256 id, uint96 amount) external;

  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external;

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
      uint64 maxValidBlocknumber,
      uint96 amountSpent
    );

  function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) external view returns (uint256[] memory);

  function getKeeperInfo(address query) external view returns (address payee, bool active, uint96 balance);

  function getState() external view returns (State memory, Config memory, address[] memory);
}

/**
 * @dev The view methods are not actually marked as view in the implementation
 * but we want them to be easily queried off-chain. Solidity will not compile
 * if we actually inherit from this interface, so we document it here.
 */
interface AutomationRegistryInterface is AutomationRegistryBaseInterface {
  function checkUpkeep(
    uint256 upkeepId,
    address from
  )
    external
    view
    returns (bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, int256 gasWei, int256 linkEth);
}

interface AutomationRegistryExecutableInterface is AutomationRegistryBaseInterface {
  function checkUpkeep(
    uint256 upkeepId,
    address from
  )
    external
    returns (
      bytes memory performData,
      uint256 maxLinkPayment,
      uint256 gasLimit,
      uint256 adjustedGasWei,
      uint256 linkEth
    );
}
