import cbor from 'cbor'
import * as h from './support/helpers'

contract('Consumer', () => {
  const sourcePath = 'examples/Consumer.sol'
  let specId = '4c7b7ffb66b344fbaa64995af81e355a'
  let currency = 'USD'
  let link, oc, cc

  beforeEach(async () => {
    link = await h.deploy('LinkToken.sol')
    oc = await h.deploy('Oracle.sol', link.address)
    await oc.transferOwnership(h.oracleNode, {from: h.defaultAccount})
    cc = await h.deploy(sourcePath, link.address, oc.address, specId)
    await cc.transferOwnership(h.consumer, {from: h.defaultAccount})
  })

  it('has a predictable gas price', async () => {
    const rec = await h.eth.getTransactionReceipt(cc.transactionHash)
    assert.isBelow(rec.gasUsed, 1600000)
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
        await link.transfer(cc.address, web3.toWei('1', 'ether'))
      })

      it('triggers a log event in the Oracle contract', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        let log = tx.receipt.logs[2]
        assert.equal(log.address, oc.address)

        let [id, jId, wei, ver, cborData] = h.decodeRunRequest(log)
        let params = await cbor.decodeFirst(cborData)
        let expected = {
          'path': ['USD'],
          'url': 'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY'
        }

        assert.equal(`0x${h.toHex(h.rPad(specId))}`, jId)
        assert.equal(web3.toWei('1', 'ether'), h.hexToInt(wei))
        assert.equal(1, ver)
        assert.deepEqual(expected, params)
      })

      it('has a reasonable gas cost', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        assert.isBelow(tx.receipt.gasUsed, 190000)
      })
    })
  })

  describe('#fulfillData', () => {
    let response = '1,000,000.00'
    let internalId

    beforeEach(async () => {
      await link.transfer(cc.address, web3.toWei('1', 'ether'))
      await cc.requestEthereumPrice(currency)
      let event = await h.getLatestEvent(oc)
      internalId = event.args.internalId
    })

    it('records the data given to it by the oracle', async () => {
      await oc.fulfillData(internalId, response, {from: h.oracleNode})

      let currentPrice = await cc.currentPrice.call()
      assert.equal(web3.toUtf8(currentPrice), response)
    })

    it('logs the data given to it by the oracle', async () => {
      let tx = await oc.fulfillData(internalId, response, {from: h.oracleNode})
      assert.equal(2, tx.receipt.logs.length)
      let log = tx.receipt.logs[0]

      assert.equal(web3.toUtf8(log.topics[2]), response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherId

      beforeEach(async () => {
        let funcSig = h.functionSelector('fulfill(bytes32,bytes32)')
        let args = h.requestDataBytes(specId, cc.address, funcSig, 42, '')
        await h.requestDataFrom(oc, link, 0, args)
        let event = await h.getLatestEvent(oc)
        otherId = event.args.internalId
      })

      it('does not accept the data provided', async () => {
        await oc.fulfillData(otherId, response, {from: h.oracleNode})

        let received = await cc.currentPrice.call()
        assert.equal(web3.toUtf8(received), '')
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await cc.fulfill(internalId, response, {from: h.oracleNode})
        })

        let received = await cc.currentPrice.call()
        assert.equal(web3.toUtf8(received), '')
      })
    })
  })

  describe('#cancelRequest', () => {
    let requestId

    beforeEach(async () => {
      await link.transfer(cc.address, web3.toWei('1', 'ether'))
      await cc.requestEthereumPrice(currency)
      requestId = (await h.getLatestEvent(cc)).args.id
    })

    context('when called by a non-owner', () => {
      it('cannot cancel a request', async () => {
        await h.assertActionThrows(async () => {
          await cc.cancelRequest(requestId, {from: h.stranger})
        })
      })
    })

    context('when called by the owner', () => {
      it('can cancel the request', async () => {
        await cc.cancelRequest(requestId, {from: h.consumer})
      })
    })
  })
})
