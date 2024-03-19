// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {AutomationRegistryLogicC2_3} from "./AutomationRegistryLogicC2_3.sol";
import {Chainable} from "../../Chainable.sol";

contract AutomationRegistryLogicB2_3 is AutomationRegistryBase2_3, Chainable {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @param logicC the address of the third logic contract
   */
  constructor(
    AutomationRegistryLogicC2_3 logicC
  )
    AutomationRegistryBase2_3(
      logicC.getLinkAddress(),
      logicC.getLinkUSDFeedAddress(),
      logicC.getNativeUSDFeedAddress(),
      logicC.getFastGasFeedAddress(),
      logicC.getAutomationForwarderLogic(),
      logicC.getAllowedReadOnlyAddress(),
      logicC.getPayoutMode()
    )
    Chainable(address(logicC))
  {}

  /**
   * @notice settles NOPs' LINK payment offchain
   */
  function settleNOPsOffchain() external {
    _onlyFinanceAdminAllowed();
    if (s_payoutMode == PayoutMode.ON_CHAIN) revert MustSettleOnchain();

    uint256 length = s_transmittersList.length;
    uint256[] memory balances = new uint256[](length);
    for (uint256 i = 0; i < length; i++) {
      address transmitterAddr = s_transmittersList[i];
      uint96 balance = _updateTransmitterBalanceFromPool(transmitterAddr, s_hotVars.totalPremium, uint96(length));
      balances[i] = balance;
      s_transmitters[transmitterAddr].balance = 0;
    }

    emit NOPsSettledOffchain(s_transmittersList, balances);
  }

  /**
   * @notice disables offchain payment for NOPs
   */
  function disableOffchainPayments() external onlyOwner {
    s_payoutMode = PayoutMode.ON_CHAIN;
  }
}
