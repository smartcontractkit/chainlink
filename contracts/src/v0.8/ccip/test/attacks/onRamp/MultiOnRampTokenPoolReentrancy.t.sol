// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {OnRampSetup} from "../../onRamp/OnRampSetup.t.sol";
import {FacadeClient} from "./FacadeClient.sol";
import {ReentrantMaliciousTokenPool} from "./ReentrantMaliciousTokenPool.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/// @title MultiOnRampTokenPoolReentrancy
/// Attempts to perform a reentrancy exploit on Onramp with a malicious TokenPool
contract MultiOnRampTokenPoolReentrancy is OnRampSetup {
  FacadeClient internal s_facadeClient;
  ReentrantMaliciousTokenPool internal s_maliciousTokenPool;
  IERC20 internal s_sourceToken;
  IERC20 internal s_feeToken;
  address internal immutable i_receiver = makeAddr("receiver");

  function setUp() public virtual override {
    OnRampSetup.setUp();

    s_sourceToken = IERC20(s_sourceTokens[0]);
    s_feeToken = IERC20(s_sourceTokens[0]);

    s_facadeClient =
      new FacadeClient(address(s_sourceRouter), DEST_CHAIN_SELECTOR, s_sourceToken, s_feeToken, i_receiver);

    s_maliciousTokenPool = new ReentrantMaliciousTokenPool(
      address(s_facadeClient), s_sourceToken, address(s_mockRMN), address(s_sourceRouter)
    );

    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](1);
    chainUpdates[0] = TokenPool.ChainUpdate({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      remotePoolAddress: abi.encode(s_destPoolBySourceToken[s_sourceTokens[0]]),
      remoteTokenAddress: abi.encode(s_destTokens[0]),
      allowed: true,
      outboundRateLimiterConfig: _getOutboundRateLimiterConfig(),
      inboundRateLimiterConfig: _getInboundRateLimiterConfig()
    });
    s_maliciousTokenPool.applyChainUpdates(chainUpdates);
    s_sourcePoolByToken[address(s_sourceToken)] = address(s_maliciousTokenPool);

    Internal.PoolUpdate[] memory removes = new Internal.PoolUpdate[](1);
    removes[0].token = address(s_sourceToken);
    removes[0].pool = address(s_sourcePoolByToken[address(s_sourceToken)]);
    Internal.PoolUpdate[] memory adds = new Internal.PoolUpdate[](1);
    adds[0].token = address(s_sourceToken);
    adds[0].pool = address(s_maliciousTokenPool);

    s_tokenAdminRegistry.setPool(address(s_sourceToken), address(s_maliciousTokenPool));

    s_sourceToken.transfer(address(s_facadeClient), 1e18);
    s_feeToken.transfer(address(s_facadeClient), 1e18);
  }

  /// @dev This test was used to showcase a reentrancy exploit on OnRamp with malicious TokenPool.
  /// How it worked: OnRamp used to construct EVM2Any messages after calling TokenPool's lockOrBurn.
  /// This allowed the malicious TokenPool to break message sequencing expectations as follows:
  ///   Any user -> Facade -> 1st call to ccipSend -> pool’s lockOrBurn —>
  ///   (reenter)-> Facade -> 2nd call to ccipSend
  /// In this case, Facade's second call would produce an EVM2Any msg with a lower sequence number.
  /// The issue was fixed by implementing a reentrancy guard in OnRamp.
  function test_OnRampTokenPoolReentrancy_Success() public {
    uint256 amount = 1;

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0].token = address(s_sourceToken);
    tokenAmounts[0].amount = amount;

    Client.EVM2AnyMessage memory message1 = Client.EVM2AnyMessage({
      receiver: abi.encode(i_receiver),
      data: abi.encodePacked(uint256(1)), // message 1 contains data 1
      tokenAmounts: tokenAmounts,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000})),
      feeToken: address(s_feeToken)
    });

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message1);
    assertGt(expectedFee, 0);

    // Outcome of a successful exploit:
    // Message 1 event from OnRamp contains sequence/nonce 2, message 2 contains sequence/nonce 1
    // Internal.EVM2EVMMessage memory msgEvent1 = _messageToEvent(message1, 2, 2, expectedFee, address(s_facadeClient));
    // Internal.EVM2EVMMessage memory msgEvent2 = _messageToEvent(message2, 1, 1, expectedFee, address(s_facadeClient));

    vm.expectRevert(OnRamp.ReentrancyGuardReentrantCall.selector);
    s_facadeClient.send(amount);
  }
}
