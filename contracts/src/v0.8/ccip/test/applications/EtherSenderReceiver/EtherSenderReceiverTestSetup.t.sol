// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

import {ICCIPRouter} from "../../../applications/EtherSenderReceiver.sol";

import {EtherSenderReceiverHelper} from "../../helpers/EtherSenderReceiverHelper.sol";

import {WETH9} from "../../../../vendor/canonical-weth/WETH9.sol";
import {ERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";

contract EtherSenderReceiverTestSetup is Test {
  EtherSenderReceiverHelper internal s_etherSenderReceiver;
  WETH9 internal s_weth;
  WETH9 internal s_someOtherWeth;
  ERC20 internal s_linkToken;

  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant ROUTER = 0x0F3779ee3a832D10158073ae2F5e61ac7FBBF880;
  address internal constant XCHAIN_RECEIVER = 0xBd91b2073218AF872BF73b65e2e5950ea356d147;
  uint256 internal constant AMOUNT = 100;

  function setUp() public {
    vm.startPrank(OWNER);

    s_linkToken = new ERC20("Chainlink Token", "LINK");
    s_someOtherWeth = new WETH9();
    s_weth = new WETH9();
    vm.mockCall(ROUTER, abi.encodeWithSelector(ICCIPRouter.getWrappedNative.selector), abi.encode(address(s_weth)));
    s_etherSenderReceiver = new EtherSenderReceiverHelper(ROUTER);

    deal(OWNER, 1_000_000 ether);
    deal(address(s_linkToken), OWNER, 1_000_000 ether);

    // deposit some eth into the weth contract.
    s_weth.deposit{value: 10 ether}();
    uint256 wethSupply = s_weth.totalSupply();
    assertEq(wethSupply, 10 ether, "total weth supply must be 10 ether");
  }
}
