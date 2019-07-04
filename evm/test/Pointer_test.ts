import * as h from './support/helpers'
const Pointer = artifacts.require('Pointer.sol')

contract('Pointer', () => {
  let contract: any
  let link: any

  beforeEach(async () => {
    link = await h.linkContract()
    contract = await Pointer.new(link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Pointer, ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await contract.getAddress(), link.address)
    })
  })
})
