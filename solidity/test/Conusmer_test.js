'use strict';

require('./support/helpers.js')

contract('Consumer', () => {
  let ChainLink = artifacts.require("./contracts/ChainLink.sol");
  let Consumer = artifacts.require("./test/contracts/Consumer.sol");
  let oc, cc;

  beforeEach(async () => {
    oc = await ChainLink.new({from: oracleNode});
    cc = await Consumer.new(oc.address, {from: stranger});
  });

  describe("#requestEthereumPrice", () => {
    it("triggers a log event in the ChainLink contract", async () => {
      let tx = await cc.requestEthereumPrice();

      let events = await getEvents(oc);
      assert.equal(1, events.length)
    });
  });
});
