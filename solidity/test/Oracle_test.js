'use strict';

require('./support/helpers.js')

contract('Oracle', () => {
  let Oracle = artifacts.require("Oracle.sol");
  let GetterSetter = artifacts.require("examples/GetterSetter.sol");
  let fHash = "0x12345678";
  let jobId = "4c7b7ffb66b344fbaa64995af81e355a";
  let to = "0x80e29acb842498fe6591f020bd82766dce619d43";
  let oc;

  beforeEach(async () => {
    oc = await Oracle.new({from: oracleNode});
  });

  it("has a limited public interface", () => {
    checkPublicABI(Oracle, [
      "owner",
      "transferOwnership",
      "requestData",
      "fulfillData",
    ]);
  });

  describe("#transferOwnership", () => {
    context("when called by the owner", () => {
      beforeEach( async () => {
        await oc.transferOwnership(stranger, {from: oracleNode});
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
    it("returns the id", async () => {
      let requestId = await oc.requestData.call(jobId, to, fHash, "");
      assert.equal(1, requestId);
    });

    it("logs an event", async () => {
      let tx = await oc.requestData(jobId, to, fHash, "");
      assert.equal(1, tx.receipt.logs.length)

      let log = tx.receipt.logs[0];
      assert.equal(jobId, web3.toUtf8(log.topics[2]));
    });

    it("uses the expected event signature", async () => {
      // If updating this test, be sure to update TestServices_RunLogTopic_ExpectedEventSignature.
      let tx = await oc.requestData(jobId, to, fHash, "");
      assert.equal(1, tx.receipt.logs.length)

      let log = tx.receipt.logs[0];
      let eventSignature = "0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2";
      assert.equal(eventSignature, log.topics[0]);
    });

    it("increments the request ID", async () => {
      let tx1 = await oc.requestData(jobId, to, fHash, "");
      let requestId1 = web3.toDecimal(tx1.receipt.logs[0].topics[1]);
      let tx2 = await oc.requestData(jobId, to, fHash, "");
      let requestId2 = web3.toDecimal(tx2.receipt.logs[0].topics[1]);

      assert.notEqual(requestId1, requestId2);
    });
  });

  describe("#fulfillData", () => {
    let mock, requestId;

    beforeEach(async () => {
      mock = await GetterSetter.new();
      let fHash = functionSelector("requestedBytes32(uint256,bytes32)");
      let req = await oc.requestData(jobId, mock.address, fHash, "");
      requestId = web3.toDecimal(req.receipt.logs[0].topics[1]);
    });

    context("when the called by a non-owner", () => {
      it("raises an error", async () => {
        await assertActionThrows(async () => {
          await oc.fulfillData(requestId, "Hello World!", {from: stranger});
        });
      });
    });

    context("when called by an owner", () => {
      it("raises an error if the request ID does not exist", async () => {
        await assertActionThrows(async () => {
          await oc.fulfillData(requestId + 1, "Hello World!", {from: oracleNode});
        });
      });

      it("sets the value on the requested contract", async () => {
        await oc.fulfillData(requestId, "Hello World!", {from: oracleNode});

        let currentRequestId = await mock.requestId.call();
        assert.equal(requestId, web3.toDecimal(currentRequestId));

        let currentValue = await mock.getBytes32.call();
        assert.equal("Hello World!", web3.toUtf8(currentValue));
      });

      it("does not allow a request to be fulfilled twice", async () => {
        await oc.fulfillData(requestId, "First message!", {from: oracleNode});
        await assertActionThrows(async () => {
          await oc.fulfillData(requestId, "Second message!!", {from: oracleNode});
        });
      });
    });
  });
});
