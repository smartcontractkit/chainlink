'use strict';

require('./support/helpers.js')

contract('SpecAndRunRequester', () => {
  const Link = artifacts.require("LinkToken.sol");
  const Oracle = artifacts.require("Oracle.sol");
  const SpecAndRunRequester = artifacts.require("examples/SpecAndRunRequester.sol");
  const currency = "USD";
  let link, oc, cc;

  beforeEach(async () => {
    link = await Link.new();
    oc = await Oracle.new(link.address, {from: oracleNode});
    cc = await SpecAndRunRequester.new(link.address, oc.address, {from: consumer});
  });

  it("has a predictable gas price", async () => {
    let rec = await eth.getTransactionReceipt(cc.transactionHash);
    assert.isBelow(rec.gasUsed, 1650000);
  });

  describe("#requestEthereumPrice", () => {
    context("without LINK", () => {
      it("reverts", async () => {
        await assertActionThrows(async () => {
          await cc.requestEthereumPrice(currency);
        });
      });
    });

    context("with LINK", () => {
      beforeEach(async () => {
        await link.transfer(cc.address, web3.toWei('1', 'ether'));
      });

      it("triggers a log event in the Oracle contract", async () => {
        let tx = await cc.requestEthereumPrice(currency);
        let log = tx.receipt.logs[2];
        assert.equal(log.address, oc.address);

        let [id, wei, ver, cborData] = decodeSpecAndRunRequest(log);
        let params = await cbor.decodeFirst(cborData);
        let expected = {
          "tasks": ["httpget", "jsonparse", "ethint256", "ethtx"],
          "params": {
            "path":["USD"],
            "url":"https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY"
          }
        };

        assert.equal(web3.toWei('1', 'ether'), hexToInt(wei));
        assert.equal(1, ver);
        assert.deepEqual(expected, params);
      });

      it("has a reasonable gas cost", async () => {
        let tx = await cc.requestEthereumPrice(currency);
        assert.isBelow(tx.receipt.gasUsed, 200000);
      });
    });
  });
});
