'use strict';

require('./support/helpers.js')

contract('ConcreteChainlinked', () => {
  let Link = artifacts.require("LinkToken.sol");
  let Oracle = artifacts.require("Oracle.sol");
  let Chainlinked = artifacts.require("examples/ConcreteChainlinked.sol");
  let GetterSetter = artifacts.require("examples/GetterSetter.sol");
  let jobId = "4c7b7ffb66b344fbaa64995af81e355a";
  let cc, gs, oc, link;

  beforeEach(async () => {
    link = await Link.new();
    oc = await Oracle.new();
    gs = await GetterSetter.new();
    cc = await Chainlinked.new(link.address, oc.address);
  });

  describe("#newRun", () => {
    it("forwards the information to the oracle contract through the link token", async () => {
      let tx = await cc.publicNewRun(jobId, gs.address, "requestedBytes32(uint256,bytes32)");

      assert.equal(3, tx.receipt.logs.length);
      let transferLog = tx.receipt.logs[0];
      let transferAndCallLog = tx.receipt.logs[1];
      let oracleLog = tx.receipt.logs[2];

      let expected = "0x" + lPadHex("1") + // version number
        lPadHex("40") + // payload offset
        lPadHex("60") + // total payload length
        lPadHex("0") + // payload prefix
        lPadHex("0") + // payload internal length
        lPadHex("0"); // payload internal value
      assert.equal(expected, oracleLog.data);
    });
  });
});
