// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import "./BaseDestinationFeeManager.t.sol";

/**
 * @title BaseDestinationFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the setup functionality of the feemanager
 */
contract DestinationFeeManagerProcessFeeTest is BaseDestinationFeeManagerTest {
  function setUp() public override {
    super.setUp();
  }

  function test_WithdrawERC20() public {
    //simulate a fee
    mintLink(address(feeManager), DEFAULT_LINK_MINT_QUANTITY);

    //get the balances to ne used for comparison
    uint256 contractBalance = getLinkBalance(address(feeManager));
    uint256 adminBalance = getLinkBalance(ADMIN);

    //the amount to withdraw
    uint256 withdrawAmount = contractBalance / 2;

    //withdraw some balance
    withdraw(address(link), ADMIN, withdrawAmount, ADMIN);

    //check the balance has been reduced
    uint256 newContractBalance = getLinkBalance(address(feeManager));
    uint256 newAdminBalance = getLinkBalance(ADMIN);

    //check the balance is greater than zero
    assertGt(newContractBalance, 0);
    //check the balance has been reduced by the correct amount
    assertEq(newContractBalance, contractBalance - withdrawAmount);
    //check the admin balance has increased by the correct amount
    assertEq(newAdminBalance, adminBalance + withdrawAmount);
  }

  function test_WithdrawUnwrappedNative() public {
    //issue funds straight to the contract to bypass the lack of fallback function
    issueUnwrappedNative(address(feeManager), DEFAULT_NATIVE_MINT_QUANTITY);

    //get the balances to be used for comparison
    uint256 contractBalance = getNativeUnwrappedBalance(address(feeManager));
    uint256 adminBalance = getNativeUnwrappedBalance(ADMIN);

    //the amount to withdraw
    uint256 withdrawAmount = contractBalance / 2;

    //withdraw some balance
    withdraw(NATIVE_WITHDRAW_ADDRESS, ADMIN, withdrawAmount, ADMIN);

    //check the balance has been reduced
    uint256 newContractBalance = getNativeUnwrappedBalance(address(feeManager));
    uint256 newAdminBalance = getNativeUnwrappedBalance(ADMIN);

    //check the balance is greater than zero
    assertGt(newContractBalance, 0);
    //check the balance has been reduced by the correct amount
    assertEq(newContractBalance, contractBalance - withdrawAmount);
    //check the admin balance has increased by the correct amount
    assertEq(newAdminBalance, adminBalance + withdrawAmount);
  }

  function test_WithdrawNonAdminAddr() public {
    //simulate a fee
    mintLink(address(feeManager), DEFAULT_LINK_MINT_QUANTITY);

    //should revert if not admin
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //withdraw some balance
    withdraw(address(link), ADMIN, DEFAULT_LINK_MINT_QUANTITY, USER);
  }

  function test_eventIsEmittedAfterSurchargeIsSet() public {
    //native surcharge
    uint64 nativeSurcharge = FEE_SCALAR / 5;

    //expect an emit
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit NativeSurchargeUpdated(nativeSurcharge);

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);
  }

  function test_subscriberDiscountEventIsEmittedOnUpdate() public {
    //native surcharge
    uint64 discount = FEE_SCALAR / 3;

    //an event should be emitted
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit SubscriberDiscountUpdated(USER, DEFAULT_FEED_1_V3, address(native), discount);

    //set the surcharge
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), discount, ADMIN);
  }

  function test_eventIsEmittedUponWithdraw() public {
    //simulate a fee
    mintLink(address(feeManager), DEFAULT_LINK_MINT_QUANTITY);

    //the amount to withdraw
    uint192 withdrawAmount = 1;

    //expect an emit
    vm.expectEmit();

    //the event to be emitted
    emit Withdraw(ADMIN, ADMIN, address(link), withdrawAmount);

    //withdraw some balance
    withdraw(address(link), ADMIN, withdrawAmount, ADMIN);
  }

  function test_linkAvailableForPaymentReturnsLinkBalance() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //check there's a balance
    assertGt(getLinkBalance(address(feeManager)), 0);

    //check the link available for payment is the link balance
    assertEq(feeManager.linkAvailableForPayment(), getLinkBalance(address(feeManager)));
  }

  function test_payLinkDeficit() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3));

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //not enough funds in the reward pool should trigger an insufficient link event
    vm.expectEmit();

    IDestinationRewardManager.FeePayment[] memory contractFees = new IDestinationRewardManager.FeePayment[](1);
    contractFees[0] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    emit InsufficientLink(contractFees);

    //process the fee
    processFee(contractFees[0].poolId, payload, USER, address(native), 0);

    //double check the rewardManager balance is 0
    assertEq(getLinkBalance(address(rewardManager)), 0);

    //simulate a deposit of link to cover the deficit
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    vm.expectEmit();
    emit LinkDeficitCleared(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_LINK_FEE);

    //pay the deficit which will transfer link from the rewardManager to the rewardManager
    payLinkDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);

    //check the rewardManager received the link
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);
  }

  function test_payLinkDeficitTwice() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3));

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //not enough funds in the reward pool should trigger an insufficient link event
    vm.expectEmit();

    IDestinationRewardManager.FeePayment[] memory contractFees = new IDestinationRewardManager.FeePayment[](1);
    contractFees[0] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    //emit the event that is expected to be emitted
    emit InsufficientLink(contractFees);

    //process the fee
    processFee(contractFees[0].poolId, payload, USER, address(native), 0);

    //double check the rewardManager balance is 0
    assertEq(getLinkBalance(address(rewardManager)), 0);

    //simulate a deposit of link to cover the deficit
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    vm.expectEmit();
    emit LinkDeficitCleared(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_LINK_FEE);

    //pay the deficit which will transfer link from the rewardManager to the rewardManager
    payLinkDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);

    //check the rewardManager received the link
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //paying again should revert with 0
    vm.expectRevert(ZERO_DEFICIT);

    payLinkDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);
  }

  function test_payLinkDeficitPaysAllFeesProcessed() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3));

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 2, USER);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(DEFAULT_CONFIG_DIGEST, payload, USER, address(native), 0);
    processFee(DEFAULT_CONFIG_DIGEST, payload, USER, address(native), 0);

    //check the deficit has been increased twice
    assertEq(getLinkDeficit(DEFAULT_CONFIG_DIGEST), DEFAULT_REPORT_LINK_FEE * 2);

    //double check the rewardManager balance is 0
    assertEq(getLinkBalance(address(rewardManager)), 0);

    //simulate a deposit of link to cover the deficit
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * 2);

    vm.expectEmit();
    emit LinkDeficitCleared(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_LINK_FEE * 2);

    //pay the deficit which will transfer link from the rewardManager to the rewardManager
    payLinkDeficit(DEFAULT_CONFIG_DIGEST, ADMIN);

    //check the rewardManager received the link
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 2);
  }

  function test_payLinkDeficitOnlyCallableByAdmin() public {
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    payLinkDeficit(DEFAULT_CONFIG_DIGEST, USER);
  }

  function test_revertOnSettingAnAddressZeroVerifier() public {
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    feeManager.addVerifier(address(0));
  }

  function test_onlyCallableByOwnerReverts() public {
    address STRANGER = address(999);
    changePrank(STRANGER);
    vm.expectRevert(bytes("Only callable by owner"));
    feeManager.addVerifier(address(0));
  }

  function test_addVerifierExistingAddress() public {
    address dummyAddress = address(998);
    feeManager.addVerifier(dummyAddress);
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    feeManager.addVerifier(dummyAddress);
  }

  function test_addVerifier() public {
    address dummyAddress = address(998);
    feeManager.addVerifier(dummyAddress);
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    feeManager.addVerifier(dummyAddress);

    // check calls to setFeeRecipients it should not error unauthorized
    changePrank(dummyAddress);
    bytes32 dummyConfigDigest = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef;
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](1);
    recipients[0] = Common.AddressAndWeight(address(991), 1e18);
    feeManager.setFeeRecipients(dummyConfigDigest, recipients);

    // removing this verifier should result in unauthorized when calling  setFeeRecipients
    changePrank(ADMIN);
    feeManager.removeVerifier(dummyAddress);
    changePrank(dummyAddress);
    vm.expectRevert(UNAUTHORIZED_ERROR);
    feeManager.setFeeRecipients(dummyConfigDigest, recipients);
  }

  function test_removeVerifierZeroAaddress() public {
    address dummyAddress = address(0);
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    feeManager.removeVerifier(dummyAddress);
  }

  function test_removeVerifierNonExistentAddress() public {
    address dummyAddress = address(991);
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    feeManager.removeVerifier(dummyAddress);
  }

  function test_setRewardManagerZeroAddress() public {
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    feeManager.setRewardManager(address(0));
  }
}
