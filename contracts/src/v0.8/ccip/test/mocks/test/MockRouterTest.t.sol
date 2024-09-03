pragma solidity ^0.8.0;

import {Client} from "../../../libraries/Client.sol";

import {TokenSetup} from "../../TokenSetup.t.sol";
import {IRouterClient, MockCCIPRouter} from "../MockRouter.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract MockRouterTest is TokenSetup {
  using SafeERC20 for IERC20;

  MockCCIPRouter public mockRouter;

  uint64 public constant mockChainSelector = 123456;

  Client.EVM2AnyMessage public message;

  function setUp() public override {
    mockRouter = new MockCCIPRouter();

    //Configure the Fee to 0.1 ether for native token fees
    mockRouter.setFee(0.1 ether);

    deal(address(this), 100 ether);

    message.receiver = abi.encode(address(0x12345));
    message.data = abi.encode("Hello World");

    s_sourceFeeToken = _deploySourceToken("sLINK", type(uint256).max, 18);
  }

  function test_ccipSendWithInsufficientNativeTokens_Revert() public {
    //Should revert because did not include sufficient eth to pay for fees
    vm.expectRevert(IRouterClient.InsufficientFeeTokenAmount.selector);
    mockRouter.ccipSend(mockChainSelector, message);
  }

  function test_ccipSendWithSufficientNativeFeeTokens_Success() public {
    //ccipSend with sufficient native tokens for fees
    mockRouter.ccipSend{value: 0.1 ether}(mockChainSelector, message);
  }

  function test_ccipSendWithInvalidMsgValue_Revert() public {
    message.feeToken = address(1); //Set to non native-token fees

    vm.expectRevert(IRouterClient.InvalidMsgValue.selector);
    mockRouter.ccipSend{value: 0.1 ether}(mockChainSelector, message);
  }

  function test_ccipSendWithLinkFeeTokenbutInsufficientAllowance_Revert() public {
    message.feeToken = s_sourceFeeToken;

    vm.expectRevert(bytes("ERC20: insufficient allowance"));
    mockRouter.ccipSend(mockChainSelector, message);
  }

  function test_ccipSendWithLinkFeeTokenAndValidMsgValue_Success() public {
    message.feeToken = s_sourceFeeToken;

    vm.startPrank(OWNER, OWNER);

    IERC20(s_sourceFeeToken).safeApprove(address(mockRouter), type(uint256).max);

    mockRouter.ccipSend(mockChainSelector, message);
  }
}
