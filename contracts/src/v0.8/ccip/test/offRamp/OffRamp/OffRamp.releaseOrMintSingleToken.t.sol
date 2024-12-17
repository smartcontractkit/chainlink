// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITokenAdminRegistry} from "../../../interfaces/ITokenAdminRegistry.sol";

import {Internal} from "../../../libraries/Internal.sol";
import {Pool} from "../../../libraries/Pool.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OffRamp_releaseOrMintSingleToken is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test__releaseOrMintSingleToken() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    IERC20 dstToken1 = IERC20(s_destTokenBySourceToken[token]);
    uint256 startingBalance = dstToken1.balanceOf(OWNER);

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.expectCall(
      s_destPoolBySourceToken[token],
      abi.encodeWithSelector(
        LockReleaseTokenPool.releaseOrMint.selector,
        Pool.ReleaseOrMintInV1({
          originalSender: originalSender,
          receiver: OWNER,
          amount: amount,
          localToken: s_destTokenBySourceToken[token],
          remoteChainSelector: SOURCE_CHAIN_SELECTOR_1,
          sourcePoolAddress: tokenAmount.sourcePoolAddress,
          sourcePoolData: tokenAmount.extraData,
          offchainTokenData: offchainTokenData
        })
      )
    );

    s_offRamp.releaseOrMintSingleToken(tokenAmount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData);

    assertEq(startingBalance + amount, dstToken1.balanceOf(OWNER));
  }

  function test_RevertWhen_releaseOrMintToken_InvalidDataLength() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // Mock the call so returns 2 slots of data
    vm.mockCall(
      s_destTokenBySourceToken[token], abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER), abi.encode(0, 0)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.InvalidDataLength.selector, Internal.MAX_BALANCE_OF_RET_BYTES, 64));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR, "");
  }

  function test_RevertWhen_releaseOrMintTokenWhen_TokenHandlingError_BalanceOf() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destTokenAddress = s_destTokenBySourceToken[token];

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: destTokenAddress,
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    bytes memory revertData = "failed to balanceOf";

    // Mock the call so returns 2 slots of data
    vm.mockCallRevert(destTokenAddress, abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER), revertData);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, destTokenAddress, revertData));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR, "");
  }

  function test_RevertWhen_releaseOrMintToken_ReleaseOrMintBalanceMismatch() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    uint256 mockedStaticBalance = 50000;

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    vm.mockCall(
      s_destTokenBySourceToken[token],
      abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER),
      abi.encode(mockedStaticBalance)
    );

    vm.expectRevert(
      abi.encodeWithSelector(
        OffRamp.ReleaseOrMintBalanceMismatch.selector, amount, mockedStaticBalance, mockedStaticBalance
      )
    );

    s_offRamp.releaseOrMintSingleToken(tokenAmount, abi.encode(OWNER), OWNER, SOURCE_CHAIN_SELECTOR, "");
  }

  function test_RevertWhen_releaseOrMintToken_skip_ReleaseOrMintBalanceMismatch_if_pool() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    uint256 mockedStaticBalance = 50000;

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: s_destTokenBySourceToken[token],
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // This should make the call fail if it does not skip the check
    vm.mockCall(
      s_destTokenBySourceToken[token],
      abi.encodeWithSelector(IERC20.balanceOf.selector, OWNER),
      abi.encode(mockedStaticBalance)
    );

    s_offRamp.releaseOrMintSingleToken(
      tokenAmount, abi.encode(OWNER), s_destPoolBySourceToken[token], SOURCE_CHAIN_SELECTOR, ""
    );
  }

  function test_RevertWhen__releaseOrMintSingleToken_NotACompatiblePool() public {
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    vm.label(destToken, "destToken");
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: destToken,
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    // Address(0) should always revert
    address returnedPool = address(0);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData);

    // A contract that doesn't support the interface should also revert
    returnedPool = address(s_offRamp);

    vm.mockCall(
      address(s_tokenAdminRegistry),
      abi.encodeWithSelector(ITokenAdminRegistry.getPool.selector, destToken),
      abi.encode(returnedPool)
    );

    vm.expectRevert(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, returnedPool));

    s_offRamp.releaseOrMintSingleToken(tokenAmount, originalSender, OWNER, SOURCE_CHAIN_SELECTOR_1, offchainTokenData);
  }

  function test_RevertWhen_releaseOrMintSingleTokenWhen_TokenHandlingError_transfer() public {
    address receiver = makeAddr("receiver");
    uint256 amount = 123123;
    address token = s_sourceTokens[0];
    address destToken = s_destTokenBySourceToken[token];
    bytes memory originalSender = abi.encode(OWNER);
    bytes memory offchainTokenData = abi.encode(keccak256("offchainTokenData"));
    address destPool = s_destPoolByToken[destToken];

    Internal.Any2EVMTokenTransfer memory tokenAmount = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(s_sourcePoolByToken[token]),
      destTokenAddress: destToken,
      extraData: "",
      amount: amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    bytes memory revertData = "call reverted :o";

    vm.mockCallRevert(destToken, abi.encodeWithSelector(IERC20.transfer.selector, receiver, amount), revertData);

    vm.expectRevert(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, destPool, revertData));
    s_offRamp.releaseOrMintSingleToken(
      tokenAmount, originalSender, receiver, SOURCE_CHAIN_SELECTOR_1, offchainTokenData
    );
  }
}
