'use strict';

const abi = require('ethereumjs-abi');
const util = require('ethereumjs-util');
const cbor = require("cbor");

contract('SpecAndRunLog', () => {
  let LinkToken = artifacts.require("LinkToken.sol");
  let Oracle = artifacts.require("examples/Oracle.sol");
  let SpecAndRunLog = artifacts.require("examples/SpecAndRunLog.sol");
  let link, logger, oc, tx;

  beforeEach(async () => {
    link = await LinkToken.new();
    oc = await Oracle.new(link.address);
    logger = await SpecAndRunLog.new(link.address, oc.address);
    await link.transfer(logger.address, web3.toWei(1));
    tx = await logger.request();
  });

  it("emits the correct number of event logs", async () => {
    assert.equal(4, tx.receipt.logs.length);
  });

  it("uses the expected event signature from the oracle", () => {
    let log = tx.receipt.logs[2];
    let eventSignature = "0x40a86f3bd301164dcd67d63d081ecb2db540ac73bafb27eea27d65b3a2694f39";
    assert.equal(eventSignature, log.topics[0]);
    assert.equal(log.address, oc.address);
  });

  it("emits the on-chain data", async () => {
    let log = tx.receipt.logs[2];
    let runABI = util.toBuffer(log.data);
    let types = ["uint256", "bytes"];
    let [version, cborData] = abi.rawDecode(types, runABI);
    let params = await cbor.decodeFirst(cborData);
    let expected = {
      "tasks": ["httppost"],
      "params": {
        "msg":"hello_chainlink",
        "url":"http://localhost:6690"
      }
    };

    assert.deepEqual(expected, params);
  });
});

