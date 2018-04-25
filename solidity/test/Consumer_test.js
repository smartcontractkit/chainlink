'use strict';

require('./support/helpers.js')

contract('Consumer', () => {
  let Link = artifacts.require("LinkToken.sol");
  let Oracle = artifacts.require("Oracle.sol");
  let Consumer = artifacts.require("examples/Consumer.sol");
  let jobId = "4c7b7ffb66b344fbaa64995af81e355a";
  let link, oc, cc;

  beforeEach(async () => {
    link = await Link.new();
    oc = await Oracle.new(link.address, {from: oracleNode});
    cc = await Consumer.new(link.address, oc.address, {from: stranger});
  });

  it("has a predictable gas price", async () => {
    let rec = await eth.getTransactionReceipt(cc.transactionHash);
    assert.isBelow(rec.gasUsed, 1500000);
  });

  describe("#requestEthereumPrice", () => {
    context("without LINK", () => {
      it("reverts", async () => {
        await assertActionThrows(async () => {
          await cc.requestEthereumPrice("usd");
        });
      });
    });

    context("with LINK", () => {
      beforeEach(async () => {
        await link.transfer(cc.address, web3.toWei('1', 'ether'));
      });

      it("triggers a log event in the Oracle contract", async () => {
        let tx = await cc.requestEthereumPrice("usd");
        let log = tx.receipt.logs[2];
        assert.equal(log.address, oc.address);

        let [id, jId, wei, ver, cborData] = decodeRunRequest(log);
        let params = await cbor.decodeFirst(cborData);
        let expected = {
          "path":["recent", "usd"],
          "url":"https://etherprice.com/api"
        };

        assert.equal(`0x${toHex(rPad("someJobId"))}`, jId);
        assert.equal(web3.toWei('1', 'szabo'), hexToInt(wei));
        assert.equal(1, ver);
        assert.deepEqual(expected, params);
      });

      it("has a reasonable gas cost", async () => {
        let tx = await cc.requestEthereumPrice("usd");
        assert.isBelow(tx.receipt.gasUsed, 190000);
      });
    });
  });

  describe("#fulfillData", () => {
    let response = "1,000,000.00";
    let requestId;

    beforeEach(async () => {
      await link.transfer(cc.address, web3.toWei('1', 'ether'));
      await cc.requestEthereumPrice("usd");
      let event = await getLatestEvent(oc);
      requestId = event.args.id;
    });

    it("records the data given to it by the oracle", async () => {
      await oc.fulfillData(requestId, response, {from: oracleNode})

      let currentPrice = await cc.currentPrice.call();
      assert.equal(web3.toUtf8(currentPrice), response);
    });

    context("when the consumer does not recognize the request ID", () => {
      beforeEach(async () => {
        let funcSig = functionSelector("fulfill(bytes32,bytes32)");
        let args = requestDataBytes(jobId, cc.address, funcSig, 42, "");
        await requestDataFrom(oc, link, 0, args);
        let event = await getLatestEvent(oc);
        requestId = event.args.id;
      });

      it("does not accept the data provided", async () => {
        let tx = await cc.requestEthereumPrice("usd");

        await assertActionThrows(async () => {
          await oc.fulfillData(requestId, response, {from: oracleNode})
        });

        let received = await cc.currentPrice.call();
        assert.equal(web3.toUtf8(received), "");
      });
    });

    context("when called by anyone other than the oracle contract", () => {
      it("does not accept the data provided", async () => {
        await assertActionThrows(async () => {
          await cc.fulfill(requestId, response, {from: oracleNode})
        });

        let received = await cc.currentPrice.call();
        assert.equal(web3.toUtf8(received), "");
      });
    });
  });
});
