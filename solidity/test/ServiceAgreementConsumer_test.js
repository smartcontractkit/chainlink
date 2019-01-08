import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

contract('ServiceAgreementConsumer', () => {
  const sourcePath = 'examples/ServiceAgreementConsumer.sol'
  const currency = 'USD'
  let link, coord, cc, agreement

  beforeEach(async () => {
    agreement = await h.newServiceAgreement()
    link = await h.linkContract()
    coord = await h.deploy('Coordinator.sol', link.address)
    cc = await h.deploy(sourcePath, link.address, coord.address, agreement.id)
  })

  it('gas price of contract deployment is predictable', async () => {
    const rec = await h.eth.getTransactionReceipt(cc.transactionHash)
    assert.isBelow(rec.gasUsed, 1500000)
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
      const paymentAmount = h.toWei('1', 'h.ether')
      beforeEach(async () => {
        await link.transfer(cc.address, paymentAmount)
      })

      it('triggers a log event in the Coordinator contract', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        let log = tx.receipt.logs[3]
        assert.equal(log.address, coord.address)

        let [jId, requester, wei, _, ver, cborData] = h.decodeRunRequest(log) 
        let params = await h.decodeDietCBOR(cborData)
        assert.equal(agreement.id, jId)
        assertBigNum(paymentAmount, wei,
                     "Logged transfer amount differed from actual amount")
        assert.equal(cc.address.slice(2), requester.slice(26))
        assert.equal(1, ver)
        const url = 'https://min-api.cryptocompare.com/' +
              'data/price?fsym=ETH&tsyms=USD,EUR,JPY'
        assert.deepEqual(params, { 'path': 'USD', url })
      })

      it('has a reasonable gas cost', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        assert.isBelow(tx.receipt.gasUsed, 169000)
      })
    })
  })

  describe('#fulfillData', () => {
    let response = '1,000,000.00'
    let requestId

    beforeEach(async () => {
      await link.transfer(cc.address, h.toWei(1, 'ether'))
      await cc.requestEthereumPrice(currency)
      let event = await h.getLatestEvent(coord)
      requestId = event.args.requestId
    })

    it('records the data given to it by the oracle', async () => {
      await coord.fulfillData(requestId, response, { from: h.oracleNode })

      let currentPrice = await cc.currentPrice.call()
      assert.equal(h.toUtf8(currentPrice), response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherId

      beforeEach(async () => {
        let funcSig = h.functionSelector('fulfill(bytes32,bytes32)')
        let args = h.executeServiceAgreementBytes(agreement.id, cc.address, funcSig, 1, '')
        await h.requestDataFrom(coord, link, 0, args)
        let event = await h.getLatestEvent(coord)
        otherId = event.args.requestId
      })

      it('does not accept the data provided', async () => {
        await coord.fulfillData(otherId, response, { from: h.oracleNode })

        let received = await cc.currentPrice.call()
        assert.equal(h.toUtf8(received), '')
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await cc.fulfill(requestId, response, { from: h.oracleNode })
        })
        let received = await cc.currentPrice.call()
        assert.equal(h.toUtf8(received), '')
      })
    })
  })
})
