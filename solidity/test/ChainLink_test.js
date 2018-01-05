'use strict';

require('./support/helpers.js')

contract('ChainLink', () => {
  let ChainLink = artifacts.require("./contracts/ChainLink.sol");
  let GetterSetter = artifacts.require("./test/contracts/GetterSetter.sol");
  let oc;
  let fID = "0x12345678";
  let to = "0x80e29acb842498fe6591f020bd82766dce619d43";

  beforeEach(async () => {
    oc = await ChainLink.new({from: oracle});
  });

  it("has a limited public interface", () => {
    checkPublicABI(ChainLink, [
      "owner",
      "transferOwnership",
      "requestData",
      "fulfillData",
    ]);
  });

  describe("#transferOwnership", () => {
    context("when called by the owner", () => {
      beforeEach( async () => {
        await oc.transferOwnership(stranger, {from: oracle});
      });

      it("can change the owner", async () => {
        let owner = await oc.owner.call();
        assert.isTrue(web3.isAddress(owner));
        assert.equal(stranger, owner);
      });
    });

    context("when called by a non-owner", () => {
      it("cannot change the owner", async () => {
        await assertActionThrows(async () => {
          await oc.transferOwnership(stranger, {from: stranger});
        });
      });
    });
  });

  describe("#requestData", () => {
    it("logs an event", async () => {
      let tx = await oc.requestData(to, fID);
      assert.equal(1, tx.receipt.logs.length)

      let log = tx.receipt.logs[0]
      assert.equal(to, hexToAddress(log.topics[2]))
    });

    it("increments the nonce", async () => {
      let tx1 = await oc.requestData(to, fID);
      let nonce1 = web3.toDecimal(tx1.receipt.logs[0].topics[1]);
      let tx2 = await oc.requestData(to, fID);
      let nonce2 = web3.toDecimal(tx2.receipt.logs[0].topics[1]);

      assert.notEqual(nonce1, nonce2);
    });
  });

  describe("#fulfillData", () => {
    let mock, nonce;

    beforeEach(async () => {
      mock = await GetterSetter.new();
      let req = await oc.requestData(mock.address, functionID("setValue(bytes32)"));
      nonce = web3.toDecimal(req.receipt.logs[0].topics[1]);
    });

    context("when the called by a non-owner", () => {
      it("raises an error", async () => {
        await assertActionThrows(async () => {
          await oc.fulfillData(nonce, "Hello World!", {from: stranger});
        });
      });
    });

    context("when called by an owner", () => {
      it("raises an error if the request ID does not exist", async () => {
        await assertActionThrows(async () => {
          await oc.fulfillData(nonce + 1, "Hello World!", {from: oracle});
        });
      });

      it("sets the value on the requested contract", async () => {
        await oc.fulfillData(nonce, "Hello World!", {from: oracle});
        let current = await mock.value.call();
        assert.equal("Hello World!", web3.toUtf8(current));
      });

      it("does not allow a request to be fulfilled twice", async () => {
        await oc.fulfillData(nonce, "First message!", {from: oracle});
        await assertActionThrows(async () => {
          await oc.fulfillData(nonce, "Second message!!", {from: oracle});
        });
      });
    });
  });
});
