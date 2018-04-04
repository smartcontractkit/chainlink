require('./support/helpers.js')

contract('ERC677Token', (accounts) => {
  let Token677 = artifacts.require("../contracts/mocks/Token677.sol");
  let Token677ReceiverMock = artifacts.require("../contracts/mocks/Token677ReceiverMock.sol");
  let NotERC677Compatible = artifacts.require("../contracts/mocks/NotERC677Compatible.sol");

  let receiver, sender, token, transferAmount;

  beforeEach(async () => {
    receiver = await Token677ReceiverMock.new();
    sender = Accounts[0];
    token = await Token677.new(1000);
    transferAmount = 100;

    await token.transfer(sender, transferAmount);
    assert.equal(await receiver.sentValue(), 0);
  });

  describe("#transferAndCall(address, uint, bytes)", () => {
    let params;

    beforeEach(() => {
      let data = functionID("transferAndCall(address,uint256,bytes)") +
        encodeAddress(receiver.address) +
        encodeUint256(transferAmount) +
        encodeUint256(96) +
        encodeBytes("deadbeef");

      params = {
        from: sender,
        to: token.address,
        data: data,
        gas: 1000000
      }
    });

    it("transfers the tokens", async () => {
      let balance = await token.balanceOf(receiver.address);
      assert.equal(balance, 0);

      await sendTransaction(params);

      balance = await token.balanceOf(receiver.address);
      assert.equal(balance.toString(), transferAmount.toString());
    });

    it("calls the token fallback function on transfer", async () => {
      await sendTransaction(params);

      let calledFallback = await receiver.calledFallback();
      assert(calledFallback);

      let tokenSender = await receiver.tokenSender();
      assert.equal(tokenSender, sender);

      let sentValue = await receiver.sentValue();
      assert.equal(sentValue, transferAmount);
    });

    it("returns true when the transfer succeeds", async () => {
      let success = await sendTransaction(params);

      assert(success);
    });

    it("throws when the transfer fails", async () => {
      let data = "be45fd62" + // transfer(address,uint256,bytes)
        encodeAddress(receiver.address) +
        encodeUint256(100000) +
        encodeUint256(96) +
        encodeBytes("deadbeef");
      params = {
        from: sender,
        to: token.address,
        data: data,
        gas: 1000000
      }

      await assertActionThrows(async () => {
        await sendTransaction(params);
      });
    });

    context("when sending to a contract that is not ERC677 compatible", () => {
      let nonERC677;

      beforeEach(async () => {
        nonERC677 = await NotERC677Compatible.new();

        let data = functionID("transferAndCall(address,uint256,bytes)") +
          encodeAddress(nonERC677.address) +
          encodeUint256(100000) +
          encodeUint256(96) +
          encodeBytes("deadbeef");

        params = {
          from: sender,
          to: token.address,
          data: data,
          gas: 1000000
        }
      });

      it("throws an error", async () => {
        await assertActionThrows(async () => {
          await sendTransaction(params);
        });

        let balance = await token.balanceOf(nonERC677.address);
        assert.equal(balance.toString(), '0');
      });
    });
  });
});
