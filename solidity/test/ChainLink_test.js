'use strict';

contract('ChainLink', () => {
  let ChainLink = artifacts.require("./contracts/ChainLink.sol");
  let oracle;

  beforeEach(async () => {
    oracle = await ChainLink.new();
  });

  it("returns the value it was initialized with", async () => {
    let value = await oracle.value.call();
    assert.equal("Hello World!", web3.toUtf8(value));
  });

  it("sets the value", async () => {
    await oracle.setValue("Hello Jupiter");
    let value = await oracle.value.call();
    assert.equal("Hello Jupiter", web3.toUtf8(value));
  });
});
