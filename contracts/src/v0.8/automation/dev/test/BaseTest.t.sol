// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "forge-std/Test.sol";

// import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {LinkToken} from "../../../shared/token/ERC677/LinkToken.sol";
import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {AutomationForwarderLogic} from "../../AutomationForwarderLogic.sol";
import {AutomationRegistry2_3} from "../v2_3/AutomationRegistry2_3.sol";
import {AutomationRegistryBase2_3} from "../v2_3/AutomationRegistryBase2_3.sol";
import {AutomationRegistryLogicA2_3} from "../v2_3/AutomationRegistryLogicA2_3.sol";
import {AutomationRegistryLogicB2_3} from "../v2_3/AutomationRegistryLogicB2_3.sol";
import {IAutomationRegistryMaster2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";

/**
 * @title BaseTest provides basic test setup procedures and dependancies for use by other
 * unit tests
 */
contract BaseTest is Test {
  // constants
  address internal constant ZERO_ADDRESS = address(0);

  // contracts
  LinkToken internal linkToken;
  ERC20Mock internal mockERC20;
  MockV3Aggregator internal LINK_USD_FEED;
  MockV3Aggregator internal NATIVE_USD_FEED;
  MockV3Aggregator internal FAST_GAS_FEED;

  // roles
  address internal constant OWNER = address(uint160(uint256(keccak256("OWNER"))));
  address internal constant UPKEEP_ADMIN = address(uint160(uint256(keccak256("UPKEEP_ADMIN"))));
  address internal constant FINANCE_ADMIN = address(uint160(uint256(keccak256("FINANCE_ADMIN"))));

  function setUp() public virtual {
    vm.startPrank(OWNER);
    linkToken = new LinkToken();
    linkToken.grantMintRole(OWNER);
    mockERC20 = new ERC20Mock("MOCK_ERC20", "MOCK_ERC20", OWNER, 0);

    LINK_USD_FEED = new MockV3Aggregator(8, 2_000_000_000); // $20
    NATIVE_USD_FEED = new MockV3Aggregator(8, 400_000_000_000); // $4,000
    FAST_GAS_FEED = new MockV3Aggregator(0, 1_000_000_000); // 1 gwei
    vm.stopPrank();
  }

  function deployRegistry() internal returns (IAutomationRegistryMaster2_3) {
    AutomationForwarderLogic forwarderLogic = new AutomationForwarderLogic();
    AutomationRegistryLogicB2_3 logicB2_3 = new AutomationRegistryLogicB2_3(
      address(linkToken),
      address(LINK_USD_FEED),
      address(NATIVE_USD_FEED),
      address(FAST_GAS_FEED),
      address(forwarderLogic),
      ZERO_ADDRESS,
      AutomationRegistryBase2_3.PayoutMode.ON_CHAIN
    );
    AutomationRegistryLogicA2_3 logicA2_3 = new AutomationRegistryLogicA2_3(logicB2_3);
    return
      IAutomationRegistryMaster2_3(address(new AutomationRegistry2_3(AutomationRegistryLogicB2_3(address(logicA2_3))))); // wow this line is hilarious
  }
}
