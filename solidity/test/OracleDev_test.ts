import * as h from './support/helpers'

contract('OracleDev', () => {
  const sourcePath = 'OracleDev.sol'
  const priceFeed = 'LinkEx.sol'
  let link: any, ocd: any, usdFeed: any, ethFeed: any

  beforeEach(async () => {
    link = await h.linkContract()
    usdFeed = await h.deploy(priceFeed)
    ethFeed = await h.deploy(priceFeed)
    ocd = await h.deploy(sourcePath, link.address, usdFeed.address, ethFeed.address)
  })

  it('extends the public interface of the Oracle contract', () => {
    h.checkPublicABI(artifacts.require(sourcePath), [
      'EXPIRY_TIME',
      'cancelOracleRequest',
      'fulfillOracleRequest',
      'getAuthorizationStatus',
      'getEthPriceFeed',
      'getUsdPriceFeed',
      'onTokenTransfer',
      'owner',
      'renounceOwnership',
      'oracleRequest',
      'setFulfillmentPermission',
      'transferOwnership',
      'withdraw',
      'withdrawable'
    ])
  })

  describe('getEthPriceFeed', () =>{
    it('returns the address of the ETH price feed contract', async () => {
      assert.equal(await ocd.getEthPriceFeed(), ethFeed.address)
    })
  })

  describe('getUsdPriceFeed', () =>{
    it('returns the address of the USD price feed contract', async () => {
      assert.equal(await ocd.getUsdPriceFeed(), usdFeed.address)
    })
  })
})
