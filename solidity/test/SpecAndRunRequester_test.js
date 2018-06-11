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
    await cc.transferOwnership(consumer, {from: defaultAccount})
  })

  it("has a predictable gas price", async () => {
    let rec = await eth.getTransactionReceipt(cc.transactionHash);
    assert.isBelow(rec.gasUsed, 1900000);
  });

  describe("#requestEthereumPrice", () => {
    let tx, log, requestId;

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
        tx = await cc.requestEthereumPrice(currency);
        log = tx.receipt.logs[2];
        let event = await getLatestEvent(oc);
        requestId = event.args.internalId;
      });

      it("triggers a log event in the Oracle contract", async () => {
        assert.equal(log.address, oc.address);

        let [id, wei, ver, cborData] = decodeSpecAndRunRequest(log);
        let params = await cbor.decodeFirst(cborData);
        let expected = {
          "tasks": ["httpget", "jsonparse", "multiply", "ethuint256", "ethtx"],
          "params": {
            "path":["USD"],
            "times": 100,
            "url":"https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY"
          }
        };

        assert.equal(web3.toWei('1', 'ether'), hexToInt(wei));
        assert.equal(1, ver);
        assert.deepEqual(expected, params);
      });

      it("has a reasonable gas cost", async () => {
        assert.isBelow(tx.receipt.gasUsed, 200000);
      });

      it("records the data given to it by the oracle", async () => {
        let expected = 60729;
        let response = "0x" + encodeUint256(expected);
        await oc.fulfillData(requestId, response, {from: oracleNode});
  
        let currentPrice = await cc.currentPrice();
        let decoded = await abi.rawDecode(["uint256"], new Buffer(intToHexNoPrefix(currentPrice), "hex"));
        assert.equal(decoded.toString(), expected);
      });
    });
  });

  describe("#withdrawLink", () => {
    beforeEach(async () => {
      await link.transfer(cc.address, web3.toWei('1', 'ether'));
    });

    context("as a non-owner", () => {
      it("reverts", async () => {
        let startingBalance = await link.balanceOf(cc.address);
        assert.equal(startingBalance, web3.toWei('1', 'ether'));

        await assertActionThrows(async () => {
          await cc.withdrawLink({from: stranger});
        });

        let endingBalance = await link.balanceOf(cc.address);
        assert.equal(endingBalance.toString(), startingBalance.toString());

        let strangerBalance = await link.balanceOf(stranger);
        assert.equal(strangerBalance, 0);
      });
    });
    
    context("as the owner", () => {
      it("returns contract LINK to the owner", async () => {
        let startingBalance = await link.balanceOf(cc.address);
        assert.equal(startingBalance, web3.toWei('1', 'ether'));

        await cc.withdrawLink({from: consumer});

        let endingBalance = await link.balanceOf(cc.address);
        assert.equal(endingBalance, 0);

        let ownerBalance = await link.balanceOf(consumer);
        assert.equal(ownerBalance.toString(), startingBalance.toString());
      });
    });
  });
});
