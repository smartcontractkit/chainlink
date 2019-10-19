import * as h from './support/helpers'
const Aggregator = artifacts.require('Aggregator2.sol')

const personas = h.personas
const defaultAccount = h.personas.Default
const jobId1 =
  '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
let aggregator, link

contract('Aggregator2', () => {

  beforeEach(async () => {
    link = await h.linkContract(defaultAccount)
    aggregator = await Aggregator.new(link.address, { from: personas.Carol })
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Aggregator, [
      'addOracle',
      'oracleCount',
      // Ownable methods:
      'isOwner',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#addOracle', async () => {
    it('increases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount.call()
      await aggregator.addOracle(personas.Neil, { from: personas.Carol })
      const currentCount = await aggregator.oracleCount.call()

      assert.isAbove(currentCount.toNumber(), pastCount.toNumber())
    })

    context('when the oracle has already been added', async () => {
      beforeEach(async () => {
        await aggregator.addOracle(personas.Neil, { from: personas.Carol })
      })

      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.addOracle(personas.Neil, { from: personas.Carol })
        })
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.addOracle(personas.Neil, { from: personas.Neil })
        })
      })
    })
  })
})
