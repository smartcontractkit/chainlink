// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {LiquidityManager} from "../LiquidityManager.sol";
import {MockL1BridgeAdapter} from "./mocks/MockBridgeAdapter.sol";
import {LiquidityManagerBaseTest} from "./LiquidityManagerBaseTest.t.sol";
import {LiquidityManagerHelper} from "./helpers/LiquidityManagerHelper.sol";
import {MockLockReleaseTokenPool} from "./mocks/MockTokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

// FOUNDRY_PROFILE=liquiditymanager forge test --match-path src/v0.8/liquiditymanager/test/LiquidityManager.t.sol

contract LiquidityManagerSetup is LiquidityManagerBaseTest {
  event FinalizationStepCompleted(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData
  );
  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount,
    bytes bridgeSpecificPayload,
    bytes bridgeReturnData
  );
  event FinalizationFailed(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData,
    bytes reason
  );
  event FinanceRoleSet(address financeRole);
  event LiquidityAddedToContainer(address indexed provider, uint256 indexed amount);
  event LiquidityRemovedFromContainer(address indexed remover, uint256 indexed amount);
  // Liquidity container event
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed remover, uint256 indexed amount);

  error NonceAlreadyUsed(uint256 nonce);

  LiquidityManagerHelper internal s_liquidityManager;
  MockLockReleaseTokenPool internal s_lockReleaseTokenPool;
  MockL1BridgeAdapter internal s_bridgeAdapter;

  // LiquidityManager that rebalances weth.
  LiquidityManagerHelper internal s_wethRebalancer;
  MockLockReleaseTokenPool internal s_wethMockLockReleaseTokenPool;
  MockL1BridgeAdapter internal s_wethBridgeAdapter;

  function setUp() public virtual override {
    LiquidityManagerBaseTest.setUp();

    s_bridgeAdapter = new MockL1BridgeAdapter(s_l1Token, false);
    s_lockReleaseTokenPool = new MockLockReleaseTokenPool(s_l1Token);
    s_liquidityManager = new LiquidityManagerHelper(
      s_l1Token,
      i_localChainSelector,
      s_lockReleaseTokenPool,
      0,
      FINANCE
    );

    s_wethBridgeAdapter = new MockL1BridgeAdapter(IERC20(address(s_l1Weth)), true);
    s_wethMockLockReleaseTokenPool = new MockLockReleaseTokenPool(IERC20(address(s_l1Weth)));
    s_wethRebalancer = new LiquidityManagerHelper(
      IERC20(address(s_l1Weth)),
      i_localChainSelector,
      s_wethMockLockReleaseTokenPool,
      0,
      FINANCE
    );
  }
}
