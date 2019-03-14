import * as h from './support/helpers'

contract('OracleDev', () => {
  const sourcePath = 'OracleDev.sol'
  const priceFeed = 'LinkEx.sol'
  const ethRate = 370160
  const usdRate = 500000
  let link: any
  let ocd: any
  let usdFeed: any
  let ethFeed: any

  beforeEach(async () => {
    link = await h.linkContract()
    usdFeed = await h.deploy(priceFeed)
    ethFeed = await h.deploy(priceFeed)
    ocd = await h.deploy(sourcePath, link.address, {from: h.defaultAccount})
  })

  it('extends the public interface of the Oracle contract', () => {
    h.checkPublicABI(artifacts.require(sourcePath), [
      'EXPIRY_TIME',
      'cancelOracleRequest',
      'currentRate',
      'fulfillOracleRequest',
      'getAuthorizationStatus',
      'onTokenTransfer',
      'oracleRequest',
      'owner',
      'priceFeeds',
      'renounceOwnership',
      'setFulfillmentPermission',
      'setPriceFeed',
      'transferOwnership',
      'withdraw',
      'withdrawable'
    ])
  })

  describe('currentRate', () => {
    beforeEach(async () => {
      await ethFeed.addOracle(h.oracleNode, {from: h.defaultAccount})
      await usdFeed.addOracle(h.oracleNode, {from: h.defaultAccount})
      await ethFeed.update(ethRate, {from: h.oracleNode})
      await usdFeed.update(usdRate, {from: h.oracleNode})
      await ocd.setPriceFeed(ethFeed.address, 'ETH')
      await ocd.setPriceFeed(usdFeed.address, 'USD')
    })

    context('when requesting the ETH rate', () => {
      it('returns the current ETH rate', async () => {
        const currentRate = await ocd.currentRate('ETH')
        assert.equal(currentRate.toString(), ethRate.toString())
      })
    })

    context('when requesting the USD rate', () => {
      it('returns the current USD rate', async () => {
        const currentRate = await ocd.currentRate('USD')
        assert.equal(currentRate.toString(), usdRate.toString())
      })
    })
  })

  describe('setPriceFeed', () => {
    context('if a stranger tries setting a price feed', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await ocd.setPriceFeed(ethFeed.address, 'ETH', {from: h.stranger})
        })
      })
    })

    context('owner setting an ETH price feed', () => {
      beforeEach(async () => {
        await ocd.setPriceFeed(ethFeed.address, 'ETH', {from: h.defaultAccount})
      })

      it('sets the address of a price feed for a given currency', async () => {
        assert.equal(await ocd.priceFeeds.call('ETH'), ethFeed.address)
      })
    })

    context('owner setting a USD price feed', () => {
      beforeEach(async () => {
        await ocd.setPriceFeed(usdFeed.address, 'USD')
      })

      it('sets the address of a price feed for a given currency', async () => {
        assert.equal(await ocd.priceFeeds.call('USD'), usdFeed.address)
      })
    })
  })
})
