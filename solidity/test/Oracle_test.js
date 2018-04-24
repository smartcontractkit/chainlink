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

  describe("#onTokenTransfer", () => {
    context("when called from the LINK token", () => {
      it("triggers the intended method", async () => {
        let callData = requestDataBytes(jobId, to, fHash, "id", "");

        let tx = await link.transferAndCall(oc.address, 0, callData);
        assert.equal(3, tx.receipt.logs.length)
      });
    });

    context("when called from any address but the LINK token", () => {
      it("triggers the intended method", async () => {
        let callData = requestDataBytes(jobId, to, fHash, "id", "");

        await assertActionThrows(async () => {
          let tx = await oc.onTokenTransfer(oracleNode, 0, callData);
        });
      });
    });
  });

  describe("#requestData", () => {
    context("when called through the LINK token", () => {
      let log, tx;
      beforeEach(async () => {
        let args = requestDataBytes(jobId, to, fHash, "id", "");
        tx = await requestDataFrom(oc, link, 0, args);
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2];
      });

      it("logs an event", async () => {
        assert.equal(jobId, web3.toUtf8(log.topics[2]));
      });

      it("uses the expected event signature", async () => {
        // If updating this test, be sure to update TestServices_RunLogTopic_ExpectedEventSignature.
        let eventSignature = "0xd27ce9cd40e3b9de8d013e1c32693550a6f543fec0191156dc826978fffb3f48";
        assert.equal(eventSignature, log.topics[0]);
      });
    });

    context("when not called through the LINK token", () => {
      it("logs an event", async () => {
        await assertActionThrows(async () => {
          let tx = await oc.requestData(1, jobId, to, fHash, "id", "", {from: oracleNode});
        });
      });
    });
  });

  describe("#fulfillData", () => {
    let mock, requestId;
    let externalId = "XID";

    beforeEach(async () => {
      mock = await GetterSetter.new();
      let fHash = functionSelector("requestedBytes32(bytes32,bytes32)");
      let args = requestDataBytes(jobId, mock.address, fHash, externalId, "");
      let req = await requestDataFrom(oc, link, 0, args);
      requestId = req.receipt.logs[2].topics[1];
    });

    context("when called by a non-owner", () => {
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
