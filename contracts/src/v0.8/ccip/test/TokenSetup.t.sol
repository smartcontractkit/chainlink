// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import "../pools/BurnMintTokenPool.sol";
import "../pools/LockReleaseTokenPool.sol";
import "../libraries/Client.sol";
import {BurnMintERC677} from "../../shared/token/ERC677/BurnMintERC677.sol";
import {MaybeRevertingBurnMintTokenPool} from "./helpers/MaybeRevertingBurnMintTokenPool.sol";

contract TokenSetup is BaseTest {
  address[] internal s_sourceTokens;
  address[] internal s_destTokens;

  address[] internal s_sourcePools;
  address[] internal s_destPools;

  address internal s_sourceFeeToken;
  address internal s_destFeeToken;

  IPool internal s_destFeeTokenPool;

  function setUp() public virtual override {
    BaseTest.setUp();

    // Source tokens & pools
    if (s_sourceTokens.length == 0 && s_sourcePools.length == 0) {
      BurnMintERC677 sourceLink = new BurnMintERC677("sLINK", "sLNK", 18, 0);
      deal(address(sourceLink), OWNER, type(uint256).max);
      s_sourceTokens.push(address(sourceLink));
      s_sourcePools.push(address(new LockReleaseTokenPool(sourceLink, new address[](0), address(s_mockARM), true)));

      BurnMintERC677 sourceETH = new BurnMintERC677("sETH", "sETH", 18, 0);
      deal(address(sourceETH), OWNER, 2 ** 128);
      s_sourceTokens.push(address(sourceETH));
      s_sourcePools.push(address(new BurnMintTokenPool(sourceETH, new address[](0), address(s_mockARM))));
      sourceETH.grantMintAndBurnRoles(s_sourcePools[1]);
    }

    s_sourceFeeToken = s_sourceTokens[0];

    // Destination tokens & pools
    if (s_destTokens.length == 0 && s_destPools.length == 0) {
      BurnMintERC677 destLink = new BurnMintERC677("dLINK", "dLNK", 18, 0);
      deal(address(destLink), OWNER, type(uint256).max);
      s_destTokens.push(address(destLink));
      s_destPools.push(address(new LockReleaseTokenPool(destLink, new address[](0), address(s_mockARM), true)));

      BurnMintERC677 destEth = new BurnMintERC677("dETH", "dETH", 18, 0);
      deal(address(destEth), OWNER, 2 ** 128);
      s_destTokens.push(address(destEth));

      s_destPools.push(address(new MaybeRevertingBurnMintTokenPool(destEth, new address[](0), address(s_mockARM))));
      destEth.grantMintAndBurnRoles(s_destPools[1]);

      // Float the lockRelease pool with funds
      IERC20(s_destTokens[0]).transfer(address(s_destPools[0]), POOL_BALANCE);
    }

    s_destFeeToken = s_destTokens[0];
    s_destFeeTokenPool = IPool(s_destPools[0]);
  }

  function getCastedSourceEVMTokenAmountsWithZeroAmounts()
    internal
    view
    returns (Client.EVMTokenAmount[] memory tokenAmounts)
  {
    tokenAmounts = new Client.EVMTokenAmount[](s_sourceTokens.length);
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].token = s_sourceTokens[i];
    }
  }

  function getCastedDestinationEVMTokenAmountsWithZeroAmounts()
    internal
    view
    returns (Client.EVMTokenAmount[] memory tokenAmounts)
  {
    tokenAmounts = new Client.EVMTokenAmount[](s_destTokens.length);
    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      tokenAmounts[i].token = s_destTokens[i];
    }
  }

  function getCastedSourceTokens() internal view returns (IERC20[] memory sourceTokens) {
    // Convert address array into IERC20 array in one line
    sourceTokens = abi.decode(abi.encode(s_sourceTokens), (IERC20[]));
  }

  function getCastedDestinationTokens() internal view returns (IERC20[] memory destTokens) {
    // Convert address array into IERC20 array in one line
    destTokens = abi.decode(abi.encode(s_destTokens), (IERC20[]));
  }

  function getCastedSourcePools() internal view returns (IPool[] memory sourcePools) {
    // Convert address array into IPool array in one line
    sourcePools = abi.decode(abi.encode(s_sourcePools), (IPool[]));
  }

  function getCastedDestinationPools() internal view returns (IPool[] memory destPools) {
    // Convert address array into IPool array in one line
    destPools = abi.decode(abi.encode(s_destPools), (IPool[]));
  }
}
