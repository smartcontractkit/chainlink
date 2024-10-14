// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRouter} from "../../interfaces/IRouter.sol";

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {NonceManager} from "../../NonceManager.sol";
import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";
import {FeeQuoterFeeSetup} from "../feeQuoter/FeeQuoterSetup.t.sol";
import {MessageInterceptorHelper} from "../helpers/MessageInterceptorHelper.sol";
import {OnRampHelper} from "../helpers/OnRampHelper.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OnRampSetup is FeeQuoterFeeSetup {
  uint256 internal immutable i_tokenAmount0 = 9;
  uint256 internal immutable i_tokenAmount1 = 7;

  bytes32 internal s_metadataHash;

  OnRampHelper internal s_onRamp;
  MessageInterceptorHelper internal s_outboundmessageInterceptor;
  address[] internal s_offRamps;
  NonceManager internal s_outboundNonceManager;

  function setUp() public virtual override {
    super.setUp();

    s_outboundmessageInterceptor = new MessageInterceptorHelper();
    s_outboundNonceManager = new NonceManager(new address[](0));
    (s_onRamp, s_metadataHash) = _deployOnRamp(
      SOURCE_CHAIN_SELECTOR, s_sourceRouter, address(s_outboundNonceManager), address(s_tokenAdminRegistry)
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

  /// @dev a helper function to compose EVM2AnyRampMessage messages
  /// @dev it is assummed that LINK is the payment token because feeTokenAmount == feeValueJuels
  function _messageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    address originalSender
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    return _messageToEvent(
      message,
      seqNum,
      nonce,
      feeTokenAmount, // fee paid
      feeTokenAmount, // converstion to jules is the same
      originalSender
    );
  }

  function _messageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    uint256 feeValueJuels,
    address originalSender
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    return _messageToEvent(
      message,
      SOURCE_CHAIN_SELECTOR,
      DEST_CHAIN_SELECTOR,
      seqNum,
      nonce,
      feeTokenAmount,
      feeValueJuels,
      originalSender,
      s_metadataHash,
      s_tokenAdminRegistry
    );
  }

  function _generateDynamicOnRampConfig(address feeQuoter) internal pure returns (OnRamp.DynamicConfig memory) {
    return OnRamp.DynamicConfig({
      feeQuoter: feeQuoter,
      reentrancyGuardEntered: false,
      messageInterceptor: address(0),
      feeAggregator: FEE_AGGREGATOR,
      allowListAdmin: address(0)
    });
  }

  // Slicing is only available for calldata. So we have to build a new bytes array.
  function _removeFirst4Bytes(bytes memory data) internal pure returns (bytes memory) {
    bytes memory result = new bytes(data.length - 4);
    for (uint256 i = 4; i < data.length; ++i) {
      result[i - 4] = data[i];
    }
    return result;
  }

  function _generateDestChainConfigArgs(IRouter router) internal pure returns (OnRamp.DestChainConfigArgs[] memory) {
    OnRamp.DestChainConfigArgs[] memory destChainConfigs = new OnRamp.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      OnRamp.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, router: router, allowListEnabled: false});
    return destChainConfigs;
  }

  function _deployOnRamp(
    uint64 sourceChainSelector,
    IRouter router,
    address nonceManager,
    address tokenAdminRegistry
  ) internal returns (OnRampHelper, bytes32 metadataHash) {
    OnRampHelper onRamp = new OnRampHelper(
      OnRamp.StaticConfig({
        chainSelector: sourceChainSelector,
        rmnRemote: s_mockRMNRemote,
        nonceManager: nonceManager,
        tokenAdminRegistry: tokenAdminRegistry
      }),
      _generateDynamicOnRampConfig(address(s_feeQuoter)),
      _generateDestChainConfigArgs(router)
    );

    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(onRamp);

    NonceManager(nonceManager).applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );

    return (
      onRamp,
      keccak256(abi.encode(Internal.EVM_2_ANY_MESSAGE_HASH, sourceChainSelector, DEST_CHAIN_SELECTOR, address(onRamp)))
    );
  }

  function _enableOutboundMessageInterceptor() internal {
    (, address msgSender,) = vm.readCallers();

    bool resetPrank = false;

    if (msgSender != OWNER) {
      vm.stopPrank();
      vm.startPrank(OWNER);
      resetPrank = true;
    }

    OnRamp.DynamicConfig memory dynamicConfig = s_onRamp.getDynamicConfig();
    dynamicConfig.messageInterceptor = address(s_outboundmessageInterceptor);
    s_onRamp.setDynamicConfig(dynamicConfig);

    if (resetPrank) {
      vm.stopPrank();
      vm.startPrank(msgSender);
    }
  }

  function _assertStaticConfigsEqual(OnRamp.StaticConfig memory a, OnRamp.StaticConfig memory b) internal pure {
    assertEq(a.chainSelector, b.chainSelector);
    assertEq(address(a.rmnRemote), address(b.rmnRemote));
    assertEq(a.tokenAdminRegistry, b.tokenAdminRegistry);
  }

  function _assertDynamicConfigsEqual(OnRamp.DynamicConfig memory a, OnRamp.DynamicConfig memory b) internal pure {
    assertEq(a.feeQuoter, b.feeQuoter);
  }
}
