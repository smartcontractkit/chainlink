'use strict';

contract('RunLog', () => {
  let LinkToken = artifacts.require("LinkToken.sol");
  let Oracle = artifacts.require("examples/Oracle.sol");
  let RunLog = artifacts.require("examples/RunLog.sol");
  let link, logger, oc;

  beforeEach(async () => {
    link = await LinkToken.new();
    oc = await Oracle.new(link.address);
    logger = await RunLog.new(link.address, oc.address, "SOME_JOB_ID");
    await link.transfer(logger.address, web3.toWei(1));
  });

  it("has a limited public interface", async () => {
    let tx = await logger.request();
    assert.equal(4, tx.receipt.logs.length);
  });
});
