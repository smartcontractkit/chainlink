// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IPool} from "../../interfaces/IPool.sol";

import {PriceRegistry} from "../../PriceRegistry.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMMultiOnRamp} from "../../onRamp/EVM2EVMMultiOnRamp.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {EVM2EVMMultiOnRampHelper} from "../helpers/EVM2EVMMultiOnRampHelper.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract EVM2EVMMultiOnRampSetup is TokenSetup, PriceRegistrySetup {
  // Duplicate event of the CCIPSendRequested in the IOnRamp
  event CCIPSendRequested(Internal.EVM2EVMMessage message);

  address internal constant CUSTOM_TOKEN = address(12345);
  uint224 internal constant CUSTOM_TOKEN_PRICE = 1e17; // $0.1 CUSTOM

  uint256 internal immutable i_tokenAmount0 = 9;
  uint256 internal immutable i_tokenAmount1 = 7;

  bytes32 internal s_metadataHash;

  EVM2EVMMultiOnRampHelper internal s_onRamp;
  address[] internal s_offRamps;

  address internal s_destTokenPool = makeAddr("destTokenPool");

  EVM2EVMMultiOnRamp.FeeTokenConfigArgs[] internal s_feeTokenConfigArgs;
  EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs[] internal s_tokenTransferFeeConfigArgs;

  function setUp() public virtual override(TokenSetup, PriceRegistrySetup) {
    TokenSetup.setUp();
    PriceRegistrySetup.setUp();

    s_priceRegistry.updatePrices(getSingleTokenPriceUpdateStruct(CUSTOM_TOKEN, CUSTOM_TOKEN_PRICE));

    address WETH = s_sourceRouter.getWrappedNative();

    s_feeTokenConfigArgs.push(
      EVM2EVMMultiOnRamp.FeeTokenConfigArgs({
        token: s_sourceFeeToken,
        networkFeeUSDCents: 1_00, // 1 USD
        gasMultiplierWeiPerEth: 1e18, // 1x
        premiumMultiplierWeiPerEth: 5e17, // 0.5x
        enabled: true
      })
    );
    s_feeTokenConfigArgs.push(
      EVM2EVMMultiOnRamp.FeeTokenConfigArgs({
        token: WETH,
        networkFeeUSDCents: 5_00, // 5 USD
        gasMultiplierWeiPerEth: 2e18, // 2x
        premiumMultiplierWeiPerEth: 2e18, // 2x
        enabled: true
      })
    );

    s_tokenTransferFeeConfigArgs.push(
      EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs({
        token: s_sourceFeeToken,
        minFeeUSDCents: 1_00, // 1 USD
        maxFeeUSDCents: 1000_00, // 1,000 USD
        deciBps: 2_5, // 2.5 bps, or 0.025%
        destGasOverhead: 40_000,
        destBytesOverhead: 0,
        aggregateRateLimitEnabled: true
      })
    );
    s_tokenTransferFeeConfigArgs.push(
      EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs({
        token: s_sourceRouter.getWrappedNative(),
        minFeeUSDCents: 50, // 0.5 USD
        maxFeeUSDCents: 500_00, // 500 USD
        deciBps: 5_0, // 5 bps, or 0.05%
        destGasOverhead: 10_000,
        destBytesOverhead: 100,
        aggregateRateLimitEnabled: true
      })
    );
    s_tokenTransferFeeConfigArgs.push(
      EVM2EVMMultiOnRamp.TokenTransferFeeConfigArgs({
        token: CUSTOM_TOKEN,
        minFeeUSDCents: 2_00, // 1 USD
        maxFeeUSDCents: 2000_00, // 1,000 USD
        deciBps: 10_0, // 10 bps, or 0.1%
        destGasOverhead: 1,
        destBytesOverhead: 200,
        aggregateRateLimitEnabled: true
      })
    );

    s_onRamp = new EVM2EVMMultiOnRampHelper(
      EVM2EVMMultiOnRamp.StaticConfig({
        linkToken: s_sourceTokens[0],
        chainSelector: SOURCE_CHAIN_SELECTOR,
        destChainSelector: DEST_CHAIN_SELECTOR,
        defaultTxGasLimit: GAS_LIMIT,
        maxNopFeesJuels: MAX_NOP_FEES_JUELS,
        prevOnRamp: address(0),
        rmnProxy: address(s_mockRMN)
      }),
      generateDynamicMultiOnRampConfig(address(s_sourceRouter), address(s_priceRegistry), address(s_tokenAdminRegistry)),
      generateDestChainDynamicConfigArgs(),
      getOutboundRateLimiterConfig(),
      s_feeTokenConfigArgs,
      s_tokenTransferFeeConfigArgs,
      getMultiOnRampNopsAndWeights()
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

  function _generateTokenMessage() public view returns (Client.EVM2AnyMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = getCastedSourceEVMTokenAmountsWithZeroAmounts();
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
    Internal.EVM2EVMMessage memory messageEvent = Internal.EVM2EVMMessage({
      sequenceNumber: seqNum,
      feeTokenAmount: feeTokenAmount,
      sender: originalSender,
      nonce: nonce,
      gasLimit: abi.decode(args, (Client.EVMExtraArgsV1)).gasLimit,
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
      address sourcePool = s_sourcePoolByToken[message.tokenAmounts[i].token];
      address destPool = s_destPoolBySourceToken[message.tokenAmounts[i].token];
      messageEvent.sourceTokenData[i] = abi.encode(
        Internal.SourceTokenData({
          sourcePoolAddress: abi.encode(sourcePool),
          destPoolAddress: abi.encode(destPool),
          extraData: ""
        })
      );
    }

    messageEvent.messageId = Internal._hash(messageEvent, s_metadataHash);
    return messageEvent;
  }

  function generateDynamicMultiOnRampConfig(
    address router,
    address priceRegistry,
    address tokenAdminRegistry
  ) internal pure returns (EVM2EVMMultiOnRamp.DynamicConfig memory) {
    return EVM2EVMMultiOnRamp.DynamicConfig({
      router: router,
      priceRegistry: priceRegistry,
      tokenAdminRegistry: tokenAdminRegistry
    });
  }

  function generateDestChainDynamicConfigArgs()
    internal
    pure
    returns (EVM2EVMMultiOnRamp.DestChainDynamicConfigArgs[] memory)
  {
    EVM2EVMMultiOnRamp.DestChainDynamicConfigArgs[] memory destChainConfigs =
      new EVM2EVMMultiOnRamp.DestChainDynamicConfigArgs[](1);
    destChainConfigs[0] = EVM2EVMMultiOnRamp.DestChainDynamicConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      dynamicConfig: EVM2EVMMultiOnRamp.DestChainDynamicConfig({
        isEnabled: true,
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        destGasOverhead: DEST_GAS_OVERHEAD,
        destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
        destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
        destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
        destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
        maxDataBytes: MAX_DATA_SIZE,
        maxPerMsgGasLimit: MAX_GAS_LIMIT,
        defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
        defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        defaultTokenDestBytesOverhead: DEFAULT_TOKEN_BYTES_OVERHEAD
      })
    });
    return destChainConfigs;
  }

  function getMultiOnRampNopsAndWeights() internal pure returns (EVM2EVMMultiOnRamp.NopAndWeight[] memory) {
    EVM2EVMMultiOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMMultiOnRamp.NopAndWeight[](3);
    nopsAndWeights[0] = EVM2EVMMultiOnRamp.NopAndWeight({nop: USER_1, weight: 19284});
    nopsAndWeights[1] = EVM2EVMMultiOnRamp.NopAndWeight({nop: USER_2, weight: 52935});
    nopsAndWeights[2] = EVM2EVMMultiOnRamp.NopAndWeight({nop: USER_3, weight: 8});
    return nopsAndWeights;
  }
}
