import { deploy } from './support/helpers'

contract('ConcreteChainlinked', () => {
  let contract: any

  beforeEach(async () => {
    contract = await deploy('LinkEx.sol')
  })

  describe('#currentRate', () => {
    it('returns 0', async () => {
      const rate = await contract.currentRate()
      assert.equal(rate, 0)
    })
  })

  describe('#updateRate', () => {
    it('returns last set value', async () => {
      await contract.update(8616460799)
      const rate = await contract.currentRate()
      assert.equal(rate, 8616460799)
    })

    it('returns the historic rate', async () => {
      await contract.update(8616460799)
      await contract.update(9970646168)
      const rate = await contract.currentRate()
      assert.equal(rate, 8616460799)
    })
  })
})
