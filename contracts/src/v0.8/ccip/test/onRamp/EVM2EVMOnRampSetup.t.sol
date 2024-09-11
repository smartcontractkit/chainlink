// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {TokenSetup} from "../TokenSetup.t.sol";

import {FeeQuoterSetup} from "../feeQuoter/FeeQuoterSetup.t.sol";
import {EVM2EVMOnRampHelper} from "../helpers/EVM2EVMOnRampHelper.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMOnRampSetup is TokenSetup, FeeQuoterSetup {
  uint256 internal immutable i_tokenAmount0 = 9;
  uint256 internal immutable i_tokenAmount1 = 7;

  bytes32 internal s_metadataHash;

  EVM2EVMOnRampHelper internal s_onRamp;
  address[] internal s_offRamps;

  address internal s_destTokenPool = makeAddr("destTokenPool");
  address internal s_destToken = makeAddr("destToken");

  EVM2EVMOnRamp.FeeTokenConfigArgs[] internal s_feeTokenConfigArgs;
  EVM2EVMOnRamp.TokenTransferFeeConfigArgs[] internal s_tokenTransferFeeConfigArgs;

  function setUp() public virtual override(TokenSetup, FeeQuoterSetup) {
    TokenSetup.setUp();
    FeeQuoterSetup.setUp();

    s_feeQuoter.updatePrices(_getSingleTokenPriceUpdateStruct(CUSTOM_TOKEN, CUSTOM_TOKEN_PRICE));

    address WETH = s_sourceRouter.getWrappedNative();

    s_feeTokenConfigArgs.push(
      EVM2EVMOnRamp.FeeTokenConfigArgs({
        token: s_sourceFeeToken,
        networkFeeUSDCents: 1_00, // 1 USD
        gasMultiplierWeiPerEth: 1e18, // 1x
        premiumMultiplierWeiPerEth: 5e17, // 0.5x
        enabled: true
      })
    );
    s_feeTokenConfigArgs.push(
      EVM2EVMOnRamp.FeeTokenConfigArgs({
        token: WETH,
        networkFeeUSDCents: 5_00, // 5 USD
        gasMultiplierWeiPerEth: 2e18, // 2x
        premiumMultiplierWeiPerEth: 2e18, // 2x
        enabled: true
      })
    );

    s_tokenTransferFeeConfigArgs.push(
      EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
        token: s_sourceFeeToken,
        minFeeUSDCents: 1_00, // 1 USD
        maxFeeUSDCents: 1000_00, // 1,000 USD
        deciBps: 2_5, // 2.5 bps, or 0.025%
        destGasOverhead: 84_000,
        destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
        aggregateRateLimitEnabled: true
      })
    );
    s_tokenTransferFeeConfigArgs.push(
      EVM2EVMOnRamp.TokenTransferFeeConfigArgs({
        token: CUSTOM_TOKEN,
        minFeeUSDCents: 2_00, // 1 USD
        maxFeeUSDCents: 500_00, // 500 USD
        deciBps: 10_0, // 10 bps, or 0.1%
        destGasOverhead: 83_000,
        destBytesOverhead: 200,
        aggregateRateLimitEnabled: true
      })
    );

    s_onRamp = new EVM2EVMOnRampHelper(
      EVM2EVMOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        destChainSelector: DEST_CHAIN_SELECTOR,
        defaultTxGasLimit: GAS_LIMIT,
        maxNopFeesJuels: MAX_NOP_FEES_JUELS,
        prevOnRamp: address(0),
        rmnProxy: address(s_mockRMN),
        tokenAdminRegistry: address(s_tokenAdminRegistry)
      }),
      generateDynamicOnRampConfig(address(s_sourceRouter), address(s_feeQuoter)),
      _getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getNopsAndWeights()
    );
    s_onRamp.setAdmin(ADMIN);

    s_metadataHash = keccak256(
      abi.encode(Internal.EVM_2_EVM_MESSAGE_HASH, SOURCE_CHAIN_SELECTOR, DEST_CHAIN_SELECTOR, address(s_onRamp))
    );

    s_offRamps = new address[](2);
    s_offRamps[0] = address(10);
    s_offRamps[1] = address(11);
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](2);
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(s_onRamp)});
    offRampUpdates[0] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: s_offRamps[0]});
    offRampUpdates[1] = Router.OffRamp({sourceChainSelector: SOURCE_CHAIN_SELECTOR, offRamp: s_offRamps[1]});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);

    // Pre approve the first token so the gas estimates of the tests
    // only cover actual gas usage from the ramps
    IERC20(s_sourceTokens[0]).approve(address(s_sourceRouter), 2 ** 128);
    IERC20(s_sourceTokens[1]).approve(address(s_sourceRouter), 2 ** 128);
  }

  function getNopsAndWeights() internal pure returns (EVM2EVMOnRamp.NopAndWeight[] memory) {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](3);
    nopsAndWeights[0] = EVM2EVMOnRamp.NopAndWeight({nop: USER_1, weight: 19284});
    nopsAndWeights[1] = EVM2EVMOnRamp.NopAndWeight({nop: USER_2, weight: 52935});
    nopsAndWeights[2] = EVM2EVMOnRamp.NopAndWeight({nop: USER_3, weight: 8});
    return nopsAndWeights;
  }

  function generateDynamicOnRampConfig(
    address router,
    address priceRegistry
  ) internal pure returns (EVM2EVMOnRamp.DynamicConfig memory) {
    return EVM2EVMOnRamp.DynamicConfig({
      router: router,
      maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
      destGasOverhead: DEST_GAS_OVERHEAD,
      destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
      destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
      destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
      destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
      priceRegistry: priceRegistry,
      maxDataBytes: MAX_DATA_SIZE,
      maxPerMsgGasLimit: MAX_GAS_LIMIT,
      defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
      defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
      enforceOutOfOrder: false
    });
  }

  function _generateTokenMessage() public view returns (Client.EVM2AnyMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = _getCastedSourceEVMTokenAmountsWithZeroAmounts();
    tokenAmounts[0].amount = i_tokenAmount0;
    tokenAmounts[1].amount = i_tokenAmount1;
    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }

  function _generateSingleTokenMessage(
    address token,
    uint256 amount
  ) public view returns (Client.EVM2AnyMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }

  function _generateEmptyMessage() public view returns (Client.EVM2AnyMessage memory) {
    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }

  function _messageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    address originalSender
  ) public view returns (Internal.EVM2EVMMessage memory) {
    // Slicing is only available for calldata. So we have to build a new bytes array.
    bytes memory args = new bytes(message.extraArgs.length - 4);
    for (uint256 i = 4; i < message.extraArgs.length; ++i) {
      args[i - 4] = message.extraArgs[i];
    }
    uint256 numberOfTokens = message.tokenAmounts.length;
    Client.EVMExtraArgsV2 memory extraArgs = _extraArgsFromBytes(bytes4(message.extraArgs), args);
    Internal.EVM2EVMMessage memory messageEvent = Internal.EVM2EVMMessage({
      sequenceNumber: seqNum,
      feeTokenAmount: feeTokenAmount,
      sender: originalSender,
      nonce: extraArgs.allowOutOfOrderExecution ? 0 : nonce,
      gasLimit: extraArgs.gasLimit,
      strict: false,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      receiver: abi.decode(message.receiver, (address)),
      data: message.data,
      tokenAmounts: message.tokenAmounts,
      sourceTokenData: new bytes[](numberOfTokens),
      feeToken: message.feeToken,
      messageId: ""
    });

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      EVM2EVMOnRamp.TokenTransferFeeConfig memory tokenTransferFeeConfig =
        s_onRamp.getTokenTransferFeeConfig(message.tokenAmounts[i].token);

      messageEvent.sourceTokenData[i] = abi.encode(
        Internal.SourceTokenData({
          sourcePoolAddress: abi.encode(s_sourcePoolByToken[message.tokenAmounts[i].token]),
          destTokenAddress: abi.encode(s_destTokenBySourceToken[message.tokenAmounts[i].token]),
          extraData: "",
          destGasAmount: tokenTransferFeeConfig.isEnabled
            ? tokenTransferFeeConfig.destGasOverhead
            : DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        })
      );
    }

    messageEvent.messageId = Internal._hash(messageEvent, s_metadataHash);
    return messageEvent;
  }

  function _extraArgsFromBytes(
    bytes4 sig,
    bytes memory extraArgData
  ) public pure returns (Client.EVMExtraArgsV2 memory) {
    if (sig == Client.EVM_EXTRA_ARGS_V1_TAG) {
      Client.EVMExtraArgsV1 memory extraArgsV1 = abi.decode(extraArgData, (Client.EVMExtraArgsV1));
      return Client.EVMExtraArgsV2({gasLimit: extraArgsV1.gasLimit, allowOutOfOrderExecution: false});
    } else if (sig == Client.EVM_EXTRA_ARGS_V2_TAG) {
      return abi.decode(extraArgData, (Client.EVMExtraArgsV2));
    } else {
      revert("Invalid extraArgs tag");
    }
  }
}
