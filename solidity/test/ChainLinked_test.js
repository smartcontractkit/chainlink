'use strict';

require('./support/helpers.js')

contract('ChainLinked', () => {
  let Oracle = artifacts.require("./contracts/Oracle.sol");
  let ChainLinked = artifacts.require("./contracts/ChainLinked.sol");

  it("has a limited public interface", () => {
    checkPublicABI(ChainLinked, []);
  });

  it("does not cost too much gas", async () => {
    let cl = await ChainLinked.new();
    let rec = await eth.getTransactionReceipt(cl.transactionHash);
    assert.equal(68653, rec.gasUsed);
  });
});
