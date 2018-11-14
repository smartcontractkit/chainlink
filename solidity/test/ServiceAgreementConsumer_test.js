import cbor from 'cbor'
import {
  assertActionThrows,
  consumer,
  decodeRunRequest,
  deploy,
  eth,
  executeServiceAgreementBytes,
  functionSelector,
  getLatestEvent,
  hexToInt,
  newHash,
  oracleNode,
  requestDataFrom,
  toHex
} from './support/helpers'

contract('ServiceAgreementConsumer', () => {
  const sourcePath = 'examples/ServiceAgreementConsumer.sol'
  let specId = newHash('0x4c7b7ffb66b344fbaa64995af81e355a')
  let currency = 'USD'
  let link, coord, cc

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    coord = await deploy('Coordinator.sol', link.address)
    cc = await deploy(sourcePath, link.address, coord.address, toHex(specId))
  })

  it('gas price of contract deployment is predictable', async () => {
    const rec = await eth.getTransactionReceipt(cc.transactionHash)
    assert.isBelow(rec.gasUsed, 1700000)
  })

  describe('#requestEthereumPrice works', () => {
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
          'path': ['USD'],
          'url': 'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY'
        }

        assert.equal(toHex(specId), jId)
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
    let internalId

    beforeEach(async () => {
      await link.transfer(cc.address, web3.toWei('1', 'ether'))
      await cc.requestEthereumPrice(currency)
      let event = await getLatestEvent(coord)
      internalId = event.args.internalId
    })

    it('records the data given to it by the oracle', async () => {
      await coord.fulfillData(internalId, response, { from: oracleNode })

      let currentPrice = await cc.currentPrice.call()
      assert.equal(web3.toUtf8(currentPrice), response)
    })

    it('logs the data given to it by the oracle', async () => {
      let tx = await coord.fulfillData(internalId, response, { from: oracleNode })
      assert.equal(2, tx.receipt.logs.length)
      let log = tx.receipt.logs[0]

      assert.equal(web3.toUtf8(log.topics[2]), response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherId

      beforeEach(async () => {
        let funcSig = functionSelector('fulfill(bytes32,bytes32)')
        let args = executeServiceAgreementBytes(toHex(specId), cc.address, funcSig, 42, '')
        await requestDataFrom(coord, link, 0, args)
        let event = await getLatestEvent(coord)
        otherId = event.args.internalId
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
          await cc.fulfill(internalId, response, { from: oracleNode })
        })
        let received = await cc.currentPrice.call()
        assert.equal(web3.toUtf8(received), '')
      })
    })
  })

  describe('#withdrawLink', () => {
    beforeEach(async () => {
      await link.transfer(cc.address, web3.toWei('1', 'ether'))
      const balance = await link.balanceOf(cc.address)
      assert.equal(balance.toString(), web3.toWei('1', 'ether'))
    })

    it('transfers LINK out of the contract', async () => {
      await cc.withdrawLink({ from: consumer })
      const ccBalance = await link.balanceOf(cc.address)
      const consumerBalance = await link.balanceOf(consumer)
      assert.equal(ccBalance.toString(), '0')
      assert.equal(consumerBalance.toString(), web3.toWei('1', 'ether'))
    })
  })
})
