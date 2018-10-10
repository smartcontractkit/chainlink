'use strict';

contract('RunLog', () => {
  let LinkToken = artifacts.require("LinkToken");
  let Oracle = artifacts.require("Oracle");
  let RunLog = artifacts.require("RunLog");
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
