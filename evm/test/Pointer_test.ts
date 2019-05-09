import * as h from './support/helpers'

contract('Pointer', () => {
  const sourcePath = 'Pointer.sol'
  let contract: any
  let link: any

  beforeEach(async () => {
    link = await h.linkContract()
    contract = await h.deploy(sourcePath, link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(sourcePath), ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await contract.getAddress(), link.address)
    })
  })
})
