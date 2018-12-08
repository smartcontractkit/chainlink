import cbor from 'cbor'
import {
  assertActionThrows,
  decodeRunRequest,
  deploy,
  eth,
  executeServiceAgreementBytes,
  functionSelector,
  getLatestEvent,
  hexToInt,
  newServiceAgreement,
  oracleNode,
  requestDataFrom
} from './support/helpers'

contract('ServiceAgreementConsumer', () => {
  const sourcePath = 'examples/ServiceAgreementConsumer.sol'
  const agreement = newServiceAgreement()
  const currency = 'USD'
  let link, coord, cc

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    coord = await deploy('Coordinator.sol', link.address)
    cc = await deploy(sourcePath, link.address, coord.address, agreement.id)
  })

  it('gas price of contract deployment is predictable', async () => {
    const rec = await eth.getTransactionReceipt(cc.transactionHash)
    assert.isBelow(rec.gasUsed, 1500000)
  })

  describe('#requestEthereumPrice', () => {
    context('without LINK', () => {
      it('reverts', async () => {
        await assertActionThrows(async () => {
          await cc.requestEthereumPrice(currency)
        })
      })
    })

    context('with LINK', () => {
      beforeEach(async () => {
        await link.transfer(cc.address, web3.toWei('1', 'ether'))
      })

      it('triggers a log event in the Coordinator contract', async () => {
        let tx = await cc.requestEthereumPrice(currency)
        let log = tx.receipt.logs[3]
        assert.equal(log.address, coord.address)

        let [jId, requester, wei, id, ver, cborData] = decodeRunRequest(log) // eslint-disable-line no-unused-vars
        let params = await cbor.decodeFirst(cborData)
        let expected = {
          'path': 'USD',
          'url': 'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY'
        }

        assert.equal(agreement.id, jId)
        assert.equal(web3.toWei('1', 'ether'), hexToInt(wei))
        assert.equal(cc.address.slice(2), requester.slice(26))
        assert.equal(1, ver)
        assert.deepEqual(expected, params)
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
      await link.transfer(cc.address, web3.toWei('1', 'ether'))
      await cc.requestEthereumPrice(currency)
      let event = await getLatestEvent(coord)
      requestId = event.args.requestId
    })

    it('records the data given to it by the oracle', async () => {
      await coord.fulfillData(requestId, response, { from: oracleNode })

      let currentPrice = await cc.currentPrice.call()
      assert.equal(web3.toUtf8(currentPrice), response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherId

      beforeEach(async () => {
        let funcSig = functionSelector('fulfill(bytes32,bytes32)')
        let args = executeServiceAgreementBytes(agreement.id, cc.address, funcSig, 1, '')
        await requestDataFrom(coord, link, 0, args)
        let event = await getLatestEvent(coord)
        otherId = event.args.requestId
      })

      it('does not accept the data provided', async () => {
        await coord.fulfillData(otherId, response, { from: oracleNode })

        let received = await cc.currentPrice.call()
        assert.equal(web3.toUtf8(received), '')
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await assertActionThrows(async () => {
          await cc.fulfill(requestId, response, { from: oracleNode })
        })
        let received = await cc.currentPrice.call()
        assert.equal(web3.toUtf8(received), '')
      })
    })
  })
})
