// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {TermsOfServiceAllowList} from "../../dev/v1_0_0/accessControl/TermsOfServiceAllowList.sol";
import {FunctionsClientTestHelper} from "./testhelpers/FunctionsClientTestHelper.sol";

import {FunctionsRoutesSetup, FunctionsOwnerAcceptTermsOfServiceSetup} from "./Setup.t.sol";

/// @notice #constructor
contract FunctionsTermsOfServiceAllowList_Constructor is FunctionsRoutesSetup {
  function test_Constructor_Success() public {
    assertEq(s_termsOfServiceAllowList.typeAndVersion(), "Functions Terms of Service Allow List v1.0.0");
    assertEq(s_termsOfServiceAllowList.owner(), OWNER_ADDRESS);
  }
}

/// @notice #getConfig
contract FunctionsTermsOfServiceAllowList_GetConfig is FunctionsRoutesSetup {
  function test_GetConfig_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    TermsOfServiceAllowList.Config memory config = s_termsOfServiceAllowList.getConfig();
    assertEq(config.enabled, getTermsOfServiceConfig().enabled);
    assertEq(config.signerPublicKey, getTermsOfServiceConfig().signerPublicKey);
  }
}

/// @notice #updateConfig
contract FunctionsTermsOfServiceAllowList_UpdateConfig is FunctionsRoutesSetup {
  function test_UpdateConfig_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_termsOfServiceAllowList.updateConfig(
      TermsOfServiceAllowList.Config({enabled: true, signerPublicKey: STRANGER_ADDRESS})
    );
  }

  event ConfigUpdated(TermsOfServiceAllowList.Config config);

  function test_UpdateConfig_Success() public {
    TermsOfServiceAllowList.Config memory configToSet = TermsOfServiceAllowList.Config({
      enabled: false,
      signerPublicKey: TOS_SIGNER
    });

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ConfigUpdated(configToSet);

    s_termsOfServiceAllowList.updateConfig(configToSet);

    TermsOfServiceAllowList.Config memory config = s_termsOfServiceAllowList.getConfig();
    assertEq(config.enabled, configToSet.enabled);
    assertEq(config.signerPublicKey, configToSet.signerPublicKey);
  }
}

/// @notice #getMessage
contract FunctionsTermsOfServiceAllowList_GetMessage is FunctionsRoutesSetup {
  function test_GetMessage_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);

    assertEq(message, keccak256(abi.encodePacked(STRANGER_ADDRESS, STRANGER_ADDRESS)));
  }
}

/// @notice #acceptTermsOfService
contract FunctionsTermsOfServiceAllowList_AcceptTermsOfService is FunctionsRoutesSetup {
  function test_AcceptTermsOfService_RevertIfBlockedSender() public {
    s_termsOfServiceAllowList.blockSender(STRANGER_ADDRESS);

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);

    vm.expectRevert(TermsOfServiceAllowList.RecipientIsBlocked.selector);

    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);
  }

  function test_AcceptTermsOfService_RevertIfInvalidSigner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(STRANGER_PRIVATE_KEY, prefixedMessage);

    vm.expectRevert(TermsOfServiceAllowList.InvalidSignature.selector);

    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);
  }

  function test_AcceptTermsOfService_RevertIfRecipientIsNotSender() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(OWNER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);

    vm.expectRevert(TermsOfServiceAllowList.InvalidUsage.selector);

    s_termsOfServiceAllowList.acceptTermsOfService(OWNER_ADDRESS, STRANGER_ADDRESS, r, s, v);
  }

  function test_AcceptTermsOfService_RevertIfAcceptorIsNotSender() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, OWNER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);

    vm.expectRevert(TermsOfServiceAllowList.InvalidUsage.selector);

    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, OWNER_ADDRESS, r, s, v);
  }

  function test_AcceptTermsOfService_RevertIfRecipientContractIsNotSender() public {
    FunctionsClientTestHelper s_functionsClientHelper = new FunctionsClientTestHelper(address(s_functionsRouter));

    // Send as externally owned account
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    // Attempt to accept for a contract account
    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, address(s_functionsClientHelper));
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);

    vm.expectRevert(TermsOfServiceAllowList.InvalidUsage.selector);

    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, address(s_functionsClientHelper), r, s, v);
  }

  event AddedAccess(address user);

  function test_AcceptTermsOfService_SuccessIfAcceptingForSelf() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit AddedAccess(STRANGER_ADDRESS);

    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);

    assertEq(s_termsOfServiceAllowList.hasAccess(STRANGER_ADDRESS, new bytes(0)), true);
  }

  function test_AcceptTermsOfService_SuccessIfAcceptingForContract() public {
    FunctionsClientTestHelper s_functionsClientHelper = new FunctionsClientTestHelper(address(s_functionsRouter));

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, address(s_functionsClientHelper));
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit AddedAccess(address(s_functionsClientHelper));

    s_functionsClientHelper.acceptTermsOfService(STRANGER_ADDRESS, address(s_functionsClientHelper), r, s, v);

    assertEq(s_termsOfServiceAllowList.hasAccess(address(s_functionsClientHelper), new bytes(0)), true);
  }
}

/// @notice #getAllAllowedSenders
contract FunctionsTermsOfServiceAllowList_GetAllAllowedSenders is FunctionsOwnerAcceptTermsOfServiceSetup {
  function test_GetAllAllowedSenders_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    address[] memory expectedSenders = new address[](1);
    expectedSenders[0] = OWNER_ADDRESS;

    assertEq(s_termsOfServiceAllowList.getAllAllowedSenders(), expectedSenders);
  }
}

/// @notice #hasAccess
contract FunctionsTermsOfServiceAllowList_HasAccess is FunctionsRoutesSetup {
  function test_HasAccess_FalseWhenEnabled() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    // Check access of account that is not on the allow list
    assertEq(s_termsOfServiceAllowList.hasAccess(STRANGER_ADDRESS, new bytes(0)), false);
  }

  function test_HasAccess_TrueWhenDisabled() public {
    // Disable allow list, which opens all access
    s_termsOfServiceAllowList.updateConfig(
      TermsOfServiceAllowList.Config({enabled: false, signerPublicKey: TOS_SIGNER})
    );

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    // Check access of account that is not on the allow list
    assertEq(s_termsOfServiceAllowList.hasAccess(STRANGER_ADDRESS, new bytes(0)), true);
  }
}

/// @notice #isBlockedSender
contract FunctionsTermsOfServiceAllowList_IsBlockedSender is FunctionsRoutesSetup {
  function test_IsBlockedSender_SuccessFalse() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    assertEq(s_termsOfServiceAllowList.isBlockedSender(STRANGER_ADDRESS), false);
  }

  function test_IsBlockedSender_SuccessTrue() public {
    // Block sender
    s_termsOfServiceAllowList.blockSender(STRANGER_ADDRESS);

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    assertEq(s_termsOfServiceAllowList.isBlockedSender(STRANGER_ADDRESS), true);
  }
}

/// @notice #blockSender
contract FunctionsTermsOfServiceAllowList_BlockSender is FunctionsRoutesSetup {
  function test_BlockSender_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_termsOfServiceAllowList.blockSender(OWNER_ADDRESS);
  }

  event BlockedAccess(address user);

  function test_BlockSender_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit BlockedAccess(STRANGER_ADDRESS);

    s_termsOfServiceAllowList.blockSender(STRANGER_ADDRESS);
    assertEq(s_termsOfServiceAllowList.hasAccess(STRANGER_ADDRESS, new bytes(0)), false);
    assertEq(s_termsOfServiceAllowList.isBlockedSender(STRANGER_ADDRESS), true);

    // Account can no longer accept Terms of Service
    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    vm.expectRevert(TermsOfServiceAllowList.RecipientIsBlocked.selector);
    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);
  }
}

/// @notice #unblockSender
contract FunctionsTermsOfServiceAllowList_UnblockSender is FunctionsRoutesSetup {
  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    s_termsOfServiceAllowList.blockSender(STRANGER_ADDRESS);
  }

  function test_UnblockSender_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_termsOfServiceAllowList.unblockSender(STRANGER_ADDRESS);
  }

  event UnblockedAccess(address user);

  function test_UnblockSender_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit UnblockedAccess(STRANGER_ADDRESS);

    s_termsOfServiceAllowList.unblockSender(STRANGER_ADDRESS);

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    // Account can now accept the Terms of Service
    bytes32 message = s_termsOfServiceAllowList.getMessage(STRANGER_ADDRESS, STRANGER_ADDRESS);
    bytes32 prefixedMessage = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message));
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(TOS_SIGNER_PRIVATE_KEY, prefixedMessage);
    s_termsOfServiceAllowList.acceptTermsOfService(STRANGER_ADDRESS, STRANGER_ADDRESS, r, s, v);
  }
}
