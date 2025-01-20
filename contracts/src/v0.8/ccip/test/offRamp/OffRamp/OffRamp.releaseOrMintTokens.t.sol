// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CallWithExactGas} from "../../../../shared/call/CallWithExactGas.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {Pool} from "../../../libraries/Pool.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {MaybeRevertingBurnMintTokenPool} from "../../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OffRamp_releaseOrMintTokens is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test_releaseOrMintTokens() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destFeeToken);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, new uint32[](0)
    );

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintTokens_WithGasOverride() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    IERC20 dstToken1 = IERC20(s_destFeeToken);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );

    uint32[] memory gasOverrides = new uint32[](sourceTokenAmounts.length);
    for (uint256 i = 0; i < gasOverrides.length; i++) {
      gasOverrides[i] = DEFAULT_TOKEN_DEST_GAS_OVERHEAD + 1;
    }
    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, gasOverrides
    );

    assertEq(startingBalance + amount1, dstToken1.balanceOf(OWNER));
  }

  function test_releaseOrMintTokens_destDenominatedDecimals() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount = 100;
    uint256 destinationDenominationMultiplier = 1000;
    srcTokenAmounts[1].amount = amount;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    address pool = s_destPoolBySourceToken[srcTokenAmounts[1].token];
    address destToken = s_destTokenBySourceToken[srcTokenAmounts[1].token];

    MaybeRevertingBurnMintTokenPool(pool).setReleaseOrMintMultiplier(destinationDenominationMultiplier);

    Client.EVMTokenAmount[] memory destTokenAmounts = s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, new uint32[](0)
    );
    assertEq(destTokenAmounts[1].amount, amount * destinationDenominationMultiplier);
    assertEq(destTokenAmounts[1].token, destToken);
  }

  // Revert

  function test_RevertWhen_releaseOrMintTokensWhen_TokenHandlingError() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();

    bytes memory unknownError = bytes("unknown error");
    s_maybeRevertingPool.setShouldRevert(unknownError);

    vm.expectRevert(
      abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, address(s_maybeRevertingPool), unknownError)
    );

    s_offRamp.releaseOrMintTokens(
      _getDefaultSourceTokenData(srcTokenAmounts),
      abi.encode(OWNER),
      OWNER,
      SOURCE_CHAIN_SELECTOR_1,
      new bytes[](srcTokenAmounts.length),
      new uint32[](0)
    );
  }

  function test_RevertWhen_releaseOrMintTokensWhenInvalidDataLengthReturnData() public {
    uint256 amount = 100;
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    srcTokenAmounts[0].amount = amount;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.mockCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      ),
      // Includes the amount twice, this will revert due to the return data being to long
      abi.encode(amount, amount)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.InvalidDataLength.selector, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, 64));

    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData, new uint32[](0)
    );
  }

  function test_RevertWhen_releaseOrMintTokensWhen_PoolIsNotAPool() public {
    // The offRamp is a contract, but not a pool
    address fakePoolAddress = address(s_offRamp);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = new Internal.Any2EVMTokenTransfer[](1);
    sourceTokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(fakePoolAddress),
      destTokenAddress: address(s_offRamp),
      extraData: "",
      amount: 1,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.expectRevert(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, address(0)));
    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, new bytes[](1), new uint32[](0)
    );
  }

  function test_RevertWhen_releaseOrMintTokensWhenPoolDoesNotSupportDest() public {
    Client.EVMTokenAmount[] memory srcTokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    uint256 amount1 = 100;
    srcTokenAmounts[0].amount = amount1;

    bytes[] memory offchainTokenData = new bytes[](srcTokenAmounts.length);
    offchainTokenData[0] = abi.encode(0x12345678);

    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = _getDefaultSourceTokenData(srcTokenAmounts);

    vm.expectCall(
      s_destPoolBySourceToken[srcTokenAmounts[0].token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: abi.encode(OWNER),
          receiver: OWNER,
          amount: srcTokenAmounts[0].amount,
          localToken: s_destTokenBySourceToken[srcTokenAmounts[0].token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_3,
          sourcePoolAddress: sourceTokenAmounts[0].sourcePoolAddress,
          sourcePoolData: sourceTokenAmounts[0].extraData,
          offchainTokenData: offchainTokenData[0]
        })
      )
    );
    vm.expectRevert();
    s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_3, offchainTokenData, new uint32[](0)
    );
  }

  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 1024
  // Uint256 gives a good range of values to test, both inside and outside of the eth address space.
  function testFuzz_releaseOrMintTokens_AnyRevertIsCaught(
    address destPool
  ) public {
    // Input 447301751254033913445893214690834296930546521452, which is 0x4E59B44847B379578588920CA78FBF26C0B4956C
    // triggers some Create2Deployer and causes it to fail
    vm.assume(destPool != 0x4e59b44847b379578588920cA78FbF26c0B4956C);
    bytes memory unusedVar = abi.encode(makeAddr("unused"));
    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts = new Internal.Any2EVMTokenTransfer[](1);
    sourceTokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: unusedVar,
      destTokenAddress: destPool,
      extraData: unusedVar,
      amount: 1,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    try s_offRamp.releaseOrMintTokens(
      sourceTokenAmounts, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR_1, new bytes[](1), new uint32[](0)
    ) {} catch (bytes memory reason) {
      // Any revert should be a TokenHandlingError, InvalidEVMAddress, InvalidDataLength or NoContract as those are caught by the offramp
      assertTrue(
        bytes4(reason) == OffRamp.TokenHandlingError.selector || bytes4(reason) == Internal.InvalidEVMAddress.selector
          || bytes4(reason) == OffRamp.InvalidDataLength.selector
          || bytes4(reason) == CallWithExactGas.NoContract.selector
          || bytes4(reason) == OffRamp.NotACompatiblePool.selector,
        "Expected TokenHandlingError or InvalidEVMAddress"
      );

      if (uint160(destPool) > type(uint160).max) {
        assertEq(reason, abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, abi.encode(destPool)));
      }
    }
  }

  function _getDefaultSourceTokenData(
    Client.EVMTokenAmount[] memory srcTokenAmounts
  ) internal view returns (Internal.Any2EVMTokenTransfer[] memory) {
    Internal.Any2EVMTokenTransfer[] memory sourceTokenData = new Internal.Any2EVMTokenTransfer[](srcTokenAmounts.length);
    for (uint256 i = 0; i < srcTokenAmounts.length; ++i) {
      sourceTokenData[i] = Internal.Any2EVMTokenTransfer({
        sourcePoolAddress: abi.encode(s_sourcePoolByToken[srcTokenAmounts[i].token]),
        destTokenAddress: s_destTokenBySourceToken[srcTokenAmounts[i].token],
        extraData: "",
        amount: srcTokenAmounts[i].amount,
        destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
      });
    }
    return sourceTokenData;
  }
}
