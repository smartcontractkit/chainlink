import * as h from './support/helpers'
const LinkEx = artifacts.require('LinkEx.sol')
const OracleDev = artifacts.require('OracleDev.sol')

contract('OracleDev', () => {
  const ethRate = 370160
  const usdRate = 500000
  const ethSymbol = h.toHex('ETH')
  const usdSymbol = h.toHex('USD')
  let link: any
  let ocd: any
  let usdFeed: any
  let ethFeed: any

  // Need to mine some blocks so that the check in update doesn't
  // fail when subtracting 25 from block.number.
  before(async () => {
    h.mineBlocks(25)
  })

  beforeEach(async () => {
    link = await h.linkContract()
    usdFeed = await LinkEx.new()
    ethFeed = await LinkEx.new()
    ocd = await OracleDev.new(link.address)
  })

  it('extends the public interface of the Oracle contract', () => {
    h.checkPublicABI(OracleDev, [
      'EXPIRY_TIME',
      'cancelOracleRequest',
      'currentRate',
      'fulfillOracleRequest',
      'getAuthorizationStatus',
      'getChainlinkToken',
      'onTokenTransfer',
      'oracleRequest',
      'priceFeeds',
      'setFulfillmentPermission',
      'setPriceFeed',
      'withdraw',
      'withdrawable',
      // Ownable methods:
      'isOwner',
      'owner',
      'transferOwnership'
    ])
  })

  describe('currentRate', () => {
    beforeEach(async () => {
      await ethFeed.addOracle(h.oracleNode, { from: h.defaultAccount })
      await usdFeed.addOracle(h.oracleNode, { from: h.defaultAccount })
      await ethFeed.update(ethRate, { from: h.oracleNode })
      await usdFeed.update(usdRate, { from: h.oracleNode })
      await ocd.setPriceFeed(ethFeed.address, ethSymbol)
      await ocd.setPriceFeed(usdFeed.address, usdSymbol)
    })

    context('when requesting the ETH rate', () => {
      it('returns the current ETH rate', async () => {
        const currentRate = await ocd.currentRate(ethSymbol)
        assert.equal(currentRate.toString(), ethRate.toString())
      })
    })

    context('when requesting the USD rate', () => {
      it('returns the current USD rate', async () => {
        const currentRate = await ocd.currentRate(usdSymbol)
        assert.equal(currentRate.toString(), usdRate.toString())
      })
    })
  })

  describe('setPriceFeed', () => {
    context('if a stranger tries setting a price feed', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await ocd.setPriceFeed(ethFeed.address, ethSymbol, {
            from: h.stranger
          })
        })
      })
    })

    context('owner setting an ETH price feed', () => {
      beforeEach(async () => {
        await ocd.setPriceFeed(ethFeed.address, ethSymbol, {
          from: h.defaultAccount
        })
      })

      it('sets the address of a price feed for a given currency', async () => {
        assert.equal(await ocd.priceFeeds.call(ethSymbol), ethFeed.address)
      })
    })

    context('owner setting a USD price feed', () => {
      beforeEach(async () => {
        await ocd.setPriceFeed(usdFeed.address, usdSymbol)
      })

      it('sets the address of a price feed for a given currency', async () => {
        assert.equal(await ocd.priceFeeds.call(usdSymbol), usdFeed.address)
      })
    })
  })
})
