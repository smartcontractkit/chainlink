'use strict';

require('./support/helpers.js')

contract('Oracle', () => {
  let Oracle = artifacts.require("Oracle.sol");
  let LinkToken = artifacts.require("LinkToken.sol");
  let GetterSetter = artifacts.require("examples/GetterSetter.sol");
  let fHash = "0x12345678";
  let jobId = "4c7b7ffb66b344fbaa64995af81e355a";
  let to = "0x80e29acb842498fe6591f020bd82766dce619d43";
  let link, oc;

  beforeEach(async () => {
    link = await LinkToken.new();
    oc = await Oracle.new(link.address, {from: oracleNode});
  });

  it("has a limited public interface", () => {
    checkPublicABI(Oracle, [
      "onTokenTransfer",
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
    it("logs an event", async () => {
      let tx = await oc.requestData(1, jobId, to, fHash, "id", "");
      assert.equal(1, tx.receipt.logs.length)

      let log1 = tx.receipt.logs[0];
      assert.equal(jobId, web3.toUtf8(log1.topics[2]));
    });

    it("uses the expected event signature", async () => {
      // If updating this test, be sure to update TestServices_RunLogTopic_ExpectedEventSignature.
      let tx = await oc.requestData(1, jobId, to, fHash, "id", "");
      assert.equal(1, tx.receipt.logs.length)

      let log = tx.receipt.logs[0];
      let eventSignature = "0xd27ce9cd40e3b9de8d013e1c32693550a6f543fec0191156dc826978fffb3f48";
      assert.equal(eventSignature, log.topics[0]);
    });
  });

  describe("#fulfillData", () => {
    let mock, requestId;
    let externalId = "XID";

    beforeEach(async () => {
      mock = await GetterSetter.new();
      let fHash = functionSelector("requestedBytes32(bytes32,bytes32)");
      let req = await oc.requestData(1, jobId, mock.address, fHash, externalId, "");
      requestId = req.logs[0].args.id;
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
          await oc.fulfillData(requestId + 10000, "Hello World!", {from: oracleNode});
        });
      });

      it("sets the value on the requested contract", async () => {
        await oc.fulfillData(requestId, "Hello World!", {from: oracleNode});

        let currentExternalId = await mock.requestId.call();
        assert.equal(externalId.toString(), web3.toUtf8(currentExternalId));

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
