'use strict';

require('./support/helpers.js');

contract('LinkToken', () => {
  let LinkToken = artifacts.require("./contracts/LinkToken.sol");
  let LinkReceiver = artifacts.require("./contracts/mocks/LinkReceiver.sol");
  let Token677ReceiverMock = artifacts.require("../contracts/mocks/Token677ReceiverMock.sol");
  let NotERC677Compatible = artifacts.require("../contracts/mocks/NotERC677Compatible.sol");
  let allowance, owner, recipient, token;

  beforeEach(async () => {
    owner = Accounts[0];
    recipient = Accounts[1];
    token = await LinkToken.new({from: owner});
  });

  it("has a limited public ABI", () => {
    let expectedABI = [
      //public attributes
      'decimals',
      'name',
      'symbol',
      'totalSupply',
      //public functions
      'allowance',
      'approve',
      'balanceOf',
      'decreaseApproval',
      'increaseApproval',
      'transfer',
      'transferAndCall',
      'transferFrom',
    ];

    checkPublicABI(LinkToken, expectedABI);
  });

  it("assigns all of the balance to the owner", async () => {
    let balance = await token.balanceOf.call(owner);

    assert.equal(balance.toString(), '1e+27');
  });

  describe("#transfer(address,uint256)", () => {
    let receiver, sender, transferAmount;

    beforeEach(async () => {
      receiver = await Token677ReceiverMock.new();
      sender = Accounts[1];
      transferAmount = 100;

      await token.transfer(sender, transferAmount, {from: owner});
      assert.equal(await receiver.sentValue(), 0);
    });

    it("does not let you transfer to the null address", async () => {
      await assertActionThrows(async () => {
        await token.transfer(emptyAddress, transferAmount, {from: sender});
      });
    });

    it("does not let you transfer to the contract itself", async () => {
      await assertActionThrows(async () => {
        await token.transfer(token.address, transferAmount, {from: sender});
      });
    });

    it("transfers the tokens", async () => {
      let balance = await token.balanceOf(receiver.address);
      assert.equal(balance, 0);

      await token.transfer(receiver.address, transferAmount, {from: sender});

      balance = await token.balanceOf(receiver.address);
      assert.equal(balance.toString(), transferAmount.toString());
    });

    it("does NOT call the fallback on transfer", async () => {
      await token.transfer(receiver.address, transferAmount, {from: sender});

      let calledFallback = await receiver.calledFallback();
      assert(!calledFallback);
    });

    it("returns true when the transfer succeeds", async () => {
      let success = await token.transfer(receiver.address, transferAmount, {from: sender});
      assert(success);
    });

    it("throws when the transfer fails", async () => {
      await assertActionThrows(async () => {
        await token.transfer(receiver.address, 100000, {from: sender});
      });
    });

    context("when sending to a contract that is not ERC677 compatible", () => {
      let nonERC677;

      beforeEach(async () => {
        nonERC677 = await NotERC677Compatible.new();
      });

      it("transfers the token", async () => {
        let balance = await token.balanceOf(nonERC677.address);
        assert.equal(balance, 0);

        await token.transfer(nonERC677.address, transferAmount, {from: sender});

        balance = await token.balanceOf(nonERC677.address);
        assert.equal(balance.toString(), transferAmount.toString());
      });
    });
  });

  describe("#transfer(address,uint256,bytes)", () => {
    let value = 1000;

    beforeEach(async () => {
      recipient = await LinkReceiver.new({from: owner});

      assert.equal(await token.allowance.call(owner, recipient.address), 0);
      assert.equal(await token.balanceOf.call(recipient.address), 0);
    });

    it("does not let you transfer to an empty address", async () => {
      await assertActionThrows(async () => {
        let data = functionID("transferAndCall(address,uint256,bytes)") +
          encodeAddress(token.address) +
          encodeUint256(value) +
          encodeUint256(96) +
          encodeBytes("");

        await sendTransaction({
          from: owner,
          to: token.address,
          data: data,
        });
      });
    });

    it("does not let you transfer to the contract itself", async () => {
      await assertActionThrows(async () => {
        let data = "be45fd62" + // transfer(address,uint256,bytes)
          encodeAddress(emptyAddress) +
          encodeUint256(value) +
          encodeUint256(96) +
          encodeBytes("");

        await sendTransaction({
          from: owner,
          to: token.address,
          data: data,
        });
      });
    });

    it("transfers the amount to the contract and calls the contract", async () => {
      let data = functionID("transferAndCall(address,uint256,bytes)") +
        encodeAddress(recipient.address) +
        encodeUint256(value) +
        encodeUint256(96) +
        encodeBytes("043e94bd"); // callbackWithoutWithdrawl()

      await sendTransaction({
        from: owner,
        to: token.address,
        data: data,
      });

      assert.equal(await token.balanceOf.call(recipient.address), value);
      assert.equal(await token.allowance.call(owner, recipient.address), 0);
      assert.equal(await recipient.fallbackCalled.call(), true);
      assert.equal(await recipient.callDataCalled.call(), true);
    });

    it("does not blow up if no data is passed", async () => {
      let data = functionID("transferAndCall(address,uint256,bytes)") +
        encodeAddress(recipient.address) +
        encodeUint256(value) +
        encodeUint256(96) +
        encodeBytes("");

      await sendTransaction({
        from: owner,
        to: token.address,
        data: data,
      });

      assert.equal(await recipient.fallbackCalled.call(), true);
      assert.equal(await recipient.callDataCalled.call(), false);
    });
  });

  describe("#approve", () => {
    let amount = 1000;

    it("allows token approval amounts to be updated without first resetting to zero", async () => {
      let originalApproval = bigNum(1000);
      await token.approve(recipient, originalApproval, {from: owner});
      let approvedAmount = await token.allowance.call(owner, recipient);
      assert.equal(approvedAmount.toString(), originalApproval.toString());

      let laterApproval = bigNum(2000);
      await token.approve(recipient, laterApproval, {from: owner});
      approvedAmount = await token.allowance.call(owner, recipient);
      assert.equal(approvedAmount.toString(), laterApproval.toString());
    });

    it("throws an error when approving the null address", async () => {
      await assertActionThrows(async () => {
        await token.approve(emptyAddress, amount, {from: owner});
      });
    });

    it("throws an error when approving the token itself", async () => {
      await assertActionThrows(async () => {
        await token.approve(token.address, amount, {from: owner});
      });
    });
  });

  describe("#transferFrom", () => {
    let amount = 1000;

    beforeEach(async () => {
        await token.transfer(recipient, amount, {from: owner});
        await token.approve(owner, amount, {from: recipient});
    });

    it("throws an error when transferring to the null address", async () => {
      await assertActionThrows(async () => {
        await token.transferFrom(recipient, emptyAddress, amount, {from: owner});
      });
    });

    it("throws an error when transferring to the token itself", async () => {
      await assertActionThrows(async () => {
        await token.transferFrom(recipient, token.address, amount, {from: owner});
      });
    });
  });
});
