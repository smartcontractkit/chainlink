'use strict';

require('./support/helpers.js')

contract('SimpleConsumer', () => {
  let Oracle = artifacts.require("Oracle.sol");
  let Consumer = artifacts.require("examples/SimpleConsumer.sol");
  let oc, cc;

  beforeEach(async () => {
    oc = await Oracle.new({from: oracleNode});
    cc = await Consumer.new(oc.address, {from: stranger});
  });

  it("has a predictable gas price", async () => {
    let rec = await eth.getTransactionReceipt(cc.transactionHash);
    assert.isBelow(rec.gasUsed, 400000);
  });

  describe("#requestEthereumPrice", () => {
    it("triggers a log event in the Oracle contract", async () => {
      let tx = await cc.requestEthereumPrice();

      let events = await getEvents(oc);
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.args.data, `{"url":"https://etherprice.com/api","path":["recent","usd"]}`)
    });

    it("has a reasonable gas cost", async () => {
      let tx = await cc.requestEthereumPrice();
      assert.isBelow(tx.receipt.gasUsed, 100000);
    });
  });

  describe("#fulfillData", () => {
    let response = "1,000,000.00";
    let requestId;

    beforeEach(async () => {
      await cc.requestEthereumPrice();
      let event = await getLatestEvent(oc);
      requestId = event.args.id;
    });

    it("records the data given to it by the oracle", async () => {
      await oc.fulfillData(requestId, response, {from: oracleNode})

      let received = await cc.currentPrice.call();
      assert.equal(web3.toUtf8(received), response);
    });
  });
});
