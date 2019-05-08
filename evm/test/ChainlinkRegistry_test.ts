import * as h from './support/helpers'

contract('ChainlinkRegistry', () => {
  const sourcePath = 'ChainlinkRegistry.sol'
  let contract: any
  let link: any

  beforeEach(async () => {
    link = await h.linkContract()
    contract = await h.deploy(sourcePath, link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(sourcePath), [
      'getChainlinkTokenAddress'
    ])
  })

  describe('#getChainlinkTokenAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await contract.getChainlinkTokenAddress(), link.address)
    })
  })
})
