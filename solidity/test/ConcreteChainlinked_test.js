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
    oc = await Oracle.new(link.address);
    gs = await GetterSetter.new();
    cc = await Chainlinked.new(link.address, oc.address);
  });

  describe("#newRun", () => {
    it("forwards the information to the oracle contract through the link token", async () => {
      let tx = await cc.publicNewRun(
        jobId,
        gs.address,
        "requestedBytes32(uint256,bytes32)");

      assert.equal(1, tx.receipt.logs.length);
      let [jId, cbAddr, cbFId, cborData] = decodeRunABI(tx.receipt.logs[0]);
      let params = await cbor.decodeFirst(cborData);

      assert.equal(jobId, jId);
      assert.equal(gs.address, `0x${cbAddr}`);
      assert.equal("d67ce1e1", toHex(cbFId));
      assert.deepEqual({}, params);
    });
  });
});
