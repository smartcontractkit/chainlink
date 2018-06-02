'use strict'

require('./support/helpers.js')

contract('SpecAndRunRequester', () => {
  const sourcePath = 'examples/SpecAndRunRequester.sol'
  const currency = 'USD'
  let link, oc, cc

  beforeEach(async () => {
    link = await deploy('linkToken/contracts/LinkToken.sol')
    oc = await deploy('Oracle.sol', link.address)
    await oc.transferOwnership(oracleNode, {from: defaultAccount})
    cc = await deploy(sourcePath, link.address, oc.address, {from: consumer});
    await cc.transferOwnership(oracleNode, {from: defaultAccount})
  })

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
