'use strict';

require('./support/helpers.js')

contract('DynamicConsumer', () => {
  let Link = artifacts.require("LinkToken.sol");
  let Oracle = artifacts.require("Oracle.sol");
  let Consumer = artifacts.require("examples/DynamicConsumer1.sol");
  let jobId = "4c7b7ffb66b344fbaa64995af81e355a";
  let link, oc, cc;

  beforeEach(async () => {
    link = await Link.new();
    oc = await Oracle.new({from: oracleNode});
    cc = await Consumer.new(link.address, oc.address, {from: stranger});
    await cc.DynamicConsumer(link.address, oc.address);
  });

  it("has a predictable gas price", async () => {
    let rec = await eth.getTransactionReceipt(cc.transactionHash);
    assert.isBelow(rec.gasUsed, 2250000);
  });

  describe("#requestEthereumPrice", () => {
    it("triggers a log event in the Oracle contract", async () => {
      let tx = await cc.requestEthereumPrice("usd");
      let log = tx.receipt.logs[2];
      assert.equal(log.address, oc.address);

      let names = lPadHex("09") + toHex(rPad("url,path,"));
      let types = lPadHex("11") + toHex(rPad("string,bytes32[],"));
      let values = lPadHex("a0") +
        lPadHex("1a") + toHex(rPad("https://etherprice.com/api")) +
        lPadHex("02") + toHex(rPad("recent") + rPad("usd"));
      let expected = names + types + values;

      let params = abi.rawDecode(["uint256", "bytes"], util.toBuffer(log.data));
      let [version, logData] = params;
      assert.equal(version, 1);
      assert.equal(toHex(logData), expected);

      let event = await getLatestEvent(oc);
      assert.equal(web3.toUtf8(event.args.jobId), "someJobId");

    });

    it("has a reasonable gas cost", async () => {
      let tx = await cc.requestEthereumPrice("usd");
      assert.isBelow(tx.receipt.gasUsed, 200000);
    });
  });

  describe("#fulfillData", () => {
    let response = "1,000,000.00";
    let requestId;

    beforeEach(async () => {
      await cc.requestEthereumPrice("usd");
      let event = await getLatestEvent(oc);
      requestId = event.args.id
    });

    it("records the data given to it by the oracle", async () => {
      await oc.fulfillData(requestId, response, {from: oracleNode})

      let currentPrice = await cc.currentPrice.call();
      assert.equal(web3.toUtf8(currentPrice), response);
    });

    context("when the consumer does not recognize the request ID", () => {
      beforeEach(async () => {
        let funcSig = functionSelector("fulfill(bytes32,bytes32)");
        let reqId = "~weird~Request~ID~";
        await oc.requestData(1, jobId, cc.address, funcSig, reqId, "");
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
