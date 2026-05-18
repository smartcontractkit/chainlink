// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {TokenProxy} from "../../applications/TokenProxy.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {EVM2EVMOnRampSetup} from "../onRamp/EVM2EVMOnRampSetup.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract TokenProxySetup is EVM2EVMOnRampSetup {
  TokenProxy internal s_tokenProxy;
  IERC20 internal s_feeToken;
  IERC20 internal s_transferToken;

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    s_feeToken = IERC20(s_sourceTokens[0]);
    s_transferToken = IERC20(s_sourceTokens[1]);
    s_tokenProxy = new TokenProxy(address(s_sourceRouter), address(s_transferToken));

    s_transferToken.approve(address(s_tokenProxy), type(uint256).max);
    s_feeToken.approve(address(s_tokenProxy), type(uint256).max);
  }
}

contract TokenProxy_constructor is TokenProxySetup {
  function test_Constructor() public view {
    assertEq(address(s_tokenProxy.getRouter()), address(s_sourceRouter));
    assertEq(address(s_tokenProxy.getToken()), address(s_transferToken));
  }
}

contract TokenProxy_getFee is TokenProxySetup {
  function test_GetFee_Success() public view {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_tokenProxy),
      data: "",
      tokenAmounts: tokens,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}))
    });

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    uint256 actualFee = s_tokenProxy.getFee(DEST_CHAIN_SELECTOR, message);
    assertEq(expectedFee, actualFee);
  }

  // Reverts

  function test_GetFeeInvalidToken_Revert() public {
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_tokenProxy),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}))
    });

    vm.expectRevert(TokenProxy.InvalidToken.selector);

    s_tokenProxy.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_GetFeeNoDataAllowed_Revert() public {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_tokenProxy),
      data: "not empty",
      tokenAmounts: tokens,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}))
    });

    vm.expectRevert(TokenProxy.NoDataAllowed.selector);

    s_tokenProxy.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_GetFeeGasShouldBeZero_Revert() public {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_tokenProxy),
      data: "",
      tokenAmounts: tokens,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 10}))
    });

    vm.expectRevert(TokenProxy.GasShouldBeZero.selector);

    s_tokenProxy.getFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract TokenProxy_ccipSend is TokenProxySetup {
  function test_CcipSend_Success() public {
    vm.pauseGasMetering();
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = tokens;
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}));

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);

    s_feeToken.approve(address(s_tokenProxy), expectedFee);

    Internal.EVM2EVMMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee, OWNER);
    msgEvent.sender = address(s_tokenProxy);
    msgEvent.messageId = Internal._hash(msgEvent, s_metadataHash);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(msgEvent);

    vm.resumeGasMetering();
    s_tokenProxy.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_CcipSendNative_Success() public {
    vm.pauseGasMetering();
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = tokens;
    message.feeToken = address(0);
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}));

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);

    Internal.EVM2EVMMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee, OWNER);
    msgEvent.sender = address(s_tokenProxy);
    msgEvent.feeToken = s_sourceRouter.getWrappedNative();
    msgEvent.messageId = Internal._hash(msgEvent, s_metadataHash);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(msgEvent);

    vm.resumeGasMetering();
    s_tokenProxy.ccipSend{value: expectedFee}(DEST_CHAIN_SELECTOR, message);
  }

  // Reverts

  function test_CcipSendInsufficientAllowance_Revert() public {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = tokens;
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}));

    // Revoke allowance
    s_transferToken.approve(address(s_tokenProxy), 0);

    vm.expectRevert("ERC20: insufficient allowance");

    s_tokenProxy.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_CcipSendInvalidToken_Revert() public {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_feeToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = tokens;
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}));

    vm.expectRevert(TokenProxy.InvalidToken.selector);

    s_tokenProxy.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_CcipSendNoDataAllowed_Revert() public {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = tokens;
    message.data = "not empty";
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 0}));

    vm.expectRevert(TokenProxy.NoDataAllowed.selector);

    s_tokenProxy.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_CcipSendGasShouldBeZero_Revert() public {
    Client.EVMTokenAmount[] memory tokens = new Client.EVMTokenAmount[](1);
    tokens[0] = Client.EVMTokenAmount({token: address(s_transferToken), amount: 1e18});

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.tokenAmounts = tokens;
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 1}));

    vm.expectRevert(TokenProxy.GasShouldBeZero.selector);

    s_tokenProxy.ccipSend(DEST_CHAIN_SELECTOR, message);
  }
}
