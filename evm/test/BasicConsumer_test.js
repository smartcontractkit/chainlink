import cbor from 'cbor'
import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const BasicConsumer = artifacts.require('BasicConsumer.sol')
const Oracle = artifacts.require('Oracle.sol')

contract('BasicConsumer', () => {
  const specId = h.newHash('0x4c7b7ffb66b344fbaa64995af81e355a')
  const currency = 'USD'
  let link, oc, cc

  beforeEach(async () => {
    link = await h.linkContract()
    oc = await Oracle.new(link.address, { from: h.oracleNode })
    cc = await BasicConsumer.new(link.address, oc.address, h.toHex(specId))
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
        const tx = await cc.requestEthereumPrice(currency)
        const log = tx.receipt.rawLogs[3]
        assert.equal(log.address.toLowerCase(), oc.address.toLowerCase())

        const request = h.decodeRunRequest(log)
        const expected = {
          path: ['USD'],
          get:
            'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY',
        }

        assert.equal(h.toHex(specId), request.jobId)
        assertBigNum(h.toWei('1', 'ether'), request.payment)
        assert.equal(cc.address.toLowerCase(), request.requester.toLowerCase())
        assert.equal(1, request.dataVersion)
        assert.deepEqual(expected, await cbor.decodeFirst(request.data))
      })

      it('has a reasonable gas cost', async () => {
        const tx = await cc.requestEthereumPrice(currency)
        assert.isBelow(tx.receipt.gasUsed, 120000)
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = '1,000,000.00'
    let request

    beforeEach(async () => {
      await link.transfer(cc.address, h.toWei('1', 'ether'))
      const tx = await cc.requestEthereumPrice(currency)
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
    })

    it('records the data given to it by the oracle', async () => {
      await h.fulfillOracleRequest(oc, request, response, {
        from: h.oracleNode,
      })

      const currentPrice = await cc.currentPrice.call()
      assert.equal(h.toUtf8(currentPrice), response)
    })

    it('logs the data given to it by the oracle', async () => {
      const tx = await h.fulfillOracleRequest(oc, request, response, {
        from: h.oracleNode,
      })
      assert.equal(2, tx.receipt.rawLogs.length)
      const log = tx.receipt.rawLogs[1]

      assert.equal(h.toUtf8(log.topics[2]), response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherRequest

      beforeEach(async () => {
        const funcSig = h.functionSelector('fulfill(bytes32,bytes32)')
        const args = h.requestDataBytes(
          h.toHex(specId),
          cc.address,
          funcSig,
          43,
          '',
        )
        const tx = await h.requestDataFrom(oc, link, 0, args)
        otherRequest = h.decodeRunRequest(tx.receipt.rawLogs[2])
      })

      it('does not accept the data provided', async () => {
        await h.fulfillOracleRequest(oc, otherRequest, response, {
          from: h.oracleNode,
        })

        const received = await cc.currentPrice.call()
        assert.equal(h.toUtf8(received), '')
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await cc.fulfill(request.id, h.toHex(response), {
            from: h.oracleNode,
          })
        })

        const received = await cc.currentPrice.call()
        assert.equal(h.toUtf8(received), '')
      })
    })
  })

  describe('#cancelRequest', () => {
    const depositAmount = h.toWei('1', 'ether')
    let request

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      const tx = await cc.requestEthereumPrice(currency)
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
    })

    context('before 5 minutes', () => {
      it('cant cancel the request', async () => {
        await h.assertActionThrows(async () => {
          await cc.cancelRequest(
            request.id,
            request.payment,
            request.callbackFunc,
            request.expiration,
            { from: h.consumer },
          )
        })
      })
    })

    context('after 5 minutes', () => {
      it('can cancel the request', async () => {
        await h.increaseTime5Minutes()

        await cc.cancelRequest(
          request.id,
          request.payment,
          request.callbackFunc,
          request.expiration,
          { from: h.consumer },
        )
      })
    })
  })

  describe('#withdrawLink', () => {
    const depositAmount = h.toWei('1', 'ether')

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      const balance = await link.balanceOf(cc.address)
      assertBigNum(balance, depositAmount)
    })

    it('transfers LINK out of the contract', async () => {
      await cc.withdrawLink({ from: h.consumer })
      const ccBalance = await link.balanceOf(cc.address)
      const consumerBalance = h.bigNum(await link.balanceOf(h.consumer))
      assertBigNum(ccBalance, 0)
      assertBigNum(consumerBalance, depositAmount)
    })
  })
})
