'use strict';

require('./support/helpers.js')

contract('ChainLink', () => {
  let ChainLink = artifacts.require("./contracts/ChainLink.sol");
  let oracle;

  beforeEach(async () => {
    oracle = await ChainLink.new();
  });


  describe("initialization", () => {
    it("returns the value it was initialized with", async () => {
      let value = await oracle.value.call();
      assert.equal("Hello World!", web3.toUtf8(value));
    });
  });

  describe("#setValue", () => {
    it("sets the value", async () => {
      await oracle.setValue("Hello Jupiter");
      let value = await oracle.value.call();
      assert.equal("Hello Jupiter", web3.toUtf8(value));
    });
  });

  describe("#transferOwnership", () => {
    it("can change the owner", async () => {
      let newOwner = "0x80e29acb842498fe6591f020bd82766dce619d43"
      await oracle.transferOwnership(newOwner);
      let owner = await oracle.owner.call();
      assert.isTrue(web3.isAddress(owner));
      assert.equal(newOwner, owner);
    });
  });

  describe("#requestData", () => {
    it("logs an event", async () => {
      let fID = "0x12345678";
      let to = "0x80e29acb842498fe6591f020bd82766dce619d43";
      let tx = await oracle.requestData(to, fID);

      assert.equal(1, tx.receipt.logs.length)

      let log = tx.receipt.logs[0]
      assert.equal(to, hexToAddress(log.topics[2]))
    });

    it("increments the nonce", async () => {
      let fID = "0x12345678";
      let to = "0x80e29acb842498fe6591f020bd82766dce619d43";
      let tx = await oracle.requestData(to, fID);
      let nonce = await oracle.nonce.call();
      assert.isTrue(bigNum(1).eq(nonce));
    });
  });
});
