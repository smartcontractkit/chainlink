import cbor from 'cbor'
import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

contract('BasicConsumer', () => {
  const sourcePath = 'examples/BasicConsumer.sol'
  let specId = h.newHash('0x4c7b7ffb66b344fbaa64995af81e355a')
  let currency = 'USD'
  let link, oc, cc

  beforeEach(async () => {
    link = await h.linkContract()
    oc = await h.deploy('Oracle.sol', link.address)
    await oc.transferOwnership(h.oracleNode, { from: h.defaultAccount })
    cc = await h.deploy(sourcePath, link.address, oc.address, h.toHex(specId))
  })

  it('has a predictable gas price', async () => {
    const rec = await h.eth.getTransactionReceipt(cc.transactionHash)
    assert.isBelow(rec.gasUsed, 1700000)
  })

  describe('#requestEthereumPrice', () => {
    context('without LINK', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await cc.requestEthereumPrice(currency)
        })
      })
    })

    context('with LINK', () => {
      beforeEach(async () => {
        await link.transfer(cc.address, h.toWei('1', 'ether'))
      })

      it('triggers a log event in the Oracle contract', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        let log = tx.receipt.logs[3]
        assert.equal(log.address, oc.address)

        let [jId, requester, wei, id, ver, addr, func, exp, cborData] = h.decodeRunRequest(log)
        let params = await cbor.decodeFirst(cborData)
        let expected = {
          'path': ['USD'],
          'url': 'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY'
        }

        assert.equal(h.toHex(specId), jId)
        assertBigNum(h.toWei('1', 'ether'), wei)
        assert.equal(cc.address.slice(2), requester.slice(26))
        assert.equal(1, ver)
        assert.deepEqual(expected, params)
      })

      it('has a reasonable gas cost', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        assert.isBelow(tx.receipt.gasUsed, 167500)
      })
    })
  })

  describe('#fulfillData', () => {
    let response = '1,000,000.00'
    let requestId

    beforeEach(async () => {
      await link.transfer(cc.address, h.toWei('1', 'ether'))
      await cc.requestEthereumPrice(currency)
      let event = await h.getLatestEvent(oc)
      requestId = event.args.requestId
    })

    it('records the data given to it by the oracle', async () => {
      await oc.fulfillData(requestId, response, {from: h.oracleNode})

      let currentPrice = await cc.currentPrice.call()
      assert.equal(h.toUtf8(currentPrice), response)
    })

    it('logs the data given to it by the oracle', async () => {
      let tx = await oc.fulfillData(requestId, response, {from: h.oracleNode})
      assert.equal(2, tx.receipt.logs.length)
      let log = tx.receipt.logs[1]

      assert.equal(h.toUtf8(log.topics[2]), response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherId

      beforeEach(async () => {
        let funcSig = h.functionSelector('fulfill(bytes32,bytes32)')
        let args = h.requestDataBytes(h.toHex(specId), cc.address, funcSig, 43, '')
        await h.requestDataFrom(oc, link, 0, args)
        let event = await h.getLatestEvent(oc)
        otherId = event.args.requestId
      })

      it('does not accept the data provided', async () => {
        await oc.fulfillData(otherId, response, {from: h.oracleNode})

        let received = await cc.currentPrice.call()
        assert.equal(h.toUtf8(received), '')
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await cc.fulfill(requestId, response, {from: h.oracleNode})
        })

        let received = await cc.currentPrice.call()
        assert.equal(h.toUtf8(received), '')
      })
    })
  })

  describe('#cancelRequest', () => {
    const depositAmount = h.toWei('1', 'ether')
    let requestId

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      await cc.requestEthereumPrice(currency)
      let event = await h.getLatestEvent(oc)
      requestId = event.args.requestId
    })

    context("before 5 minutes", () => {
      it('cant cancel the request', async () => {
        await h.assertActionThrows(async () => {
          await cc.cancelRequest(requestId, {from: h.consumer})
        })
      })
    })

    context("after 5 minutes", () => {
      it('can cancel the request', async () => {
        await h.increaseTime5Minutes();
        await cc.cancelRequest(requestId, {from: h.consumer})
      })
    })
  })

  describe('#withdrawLink', () => {
    const depositAmount = h.toWei('1', 'ether')
    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      const balance = await link.balanceOf(cc.address);
      assertBigNum(balance, depositAmount);
    })

    it('transfers LINK out of the contract', async () => {
      await cc.withdrawLink({from: h.consumer});
      const ccBalance = await link.balanceOf(cc.address);
      const consumerBalance = h.bigNum(await link.balanceOf(h.consumer));
      assertBigNum(ccBalance, 0);
      assertBigNum(consumerBalance, depositAmount);
    })
  })
})
