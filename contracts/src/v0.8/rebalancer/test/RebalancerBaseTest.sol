// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";

import {ERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract RebalancerBaseTest is Test {
  // ERC20 events
  event Transfer(address indexed from, address indexed to, uint256 value);
  event Approval(address indexed owner, address indexed spender, uint256 value);

  IERC20 internal s_l1Token;
  IERC20 internal s_l2Token;

  uint64 internal immutable i_localChainSelector = 1234;
  uint64 internal immutable i_remoteChainSelector = 9876;

  address internal constant OWNER = address(0x00000078772732723782873283);

  function setUp() public virtual {
    s_l1Token = new ERC20("l1", "L1");
    s_l2Token = new ERC20("l2", "L2");

    vm.startPrank(OWNER);
    vm.label(OWNER, "Owner");
  }
}
