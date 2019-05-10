import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const personas = h.personas

contract('ConverstionRate', () => {
  const SOURCE_PATH: string = 'ConversionRate.sol'
  const jobId = '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  let link, oc, rate

  beforeEach(async () => {
    link = await h.linkContract()
    oc = await h.deploy('Oracle.sol', link.address)
    await oc.transferOwnership(personas.Neil, { from: personas.Default })
    rate = await h.deploy(SOURCE_PATH, link.address, oc.address, jobId)
    await rate.transferOwnership(personas.Carol, { from: personas.Default })
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(SOURCE_PATH), [
      'currentRate',
      'getJobId',
      'getOracle',
      'updateCallback',
      'updateJobId',
      'updateOracle',
      'update',
      // Owable
      'owner',
      'renounceOwnership',
      'transferOwnership'
    ])
  })

  describe('#update', () => {
    const response = 100

    beforeEach(async () => {
      await link.transfer(rate.address, h.toWei('1', 'ether'), {from: personas.Default})

      const current = await rate.currentRate.call()
      assertBigNum(h.bigNum(0), current)
    })

    it('triggeers a request to the oracle and accepts a response', async () => {
      const requestTx = await rate.update()

      const log = requestTx.receipt.rawLogs[3]
      assert.equal(oc.address, log.address)
      const request = h.decodeRunRequest(log)

      await h.fulfillOracleRequest(oc, request, response, {
         from: personas.Neil
       })

      const current = await rate.currentRate.call()
      assertBigNum(response, current)
    })
  })

  describe('#updateOracle', () => {
    const updatedAddress = '0x12512b4b88566e455bf3bfbf98752da51925871a'

    beforeEach(async () => {
      assert.equal(oc.address, await rate.getOracle.call())
    })

    context('when called by the contract owner', () => {
      it('updates the oracle address', async () => {
        await rate.updateOracle(updatedAddress, { from: personas.Carol })

        const current = await rate.getOracle.call()
        assert.equal(updatedAddress, current.toLowerCase())
      })
    })

    context('when called by a non-owner', () => {
      it('does not update the oracle address', async () => {
        h.assertActionThrows(async () => {
          await rate.updateOracle(updatedAddress, { from: personas.Eddy })
        })

        assert.equal(oc.address, await rate.getOracle.call())
      })
    })
  })

  describe('#updateJobId', () => {
    const newJobId = '0x5d7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'

    beforeEach(async () => {
      assert.equal(jobId, await rate.getJobId.call())
    })

    context('when called by the contract owner', () => {
      it('updates the job ID', async () => {
        await rate.updateJobId(newJobId, { from: personas.Carol })

        assert.equal(newJobId, await rate.getJobId.call())
      })
    })

    context('when called by a non-owner', () => {
      it('does not update the job ID', async () => {
        h.assertActionThrows(async () => {
          await rate.updateJobId(newJobId, { from: personas.Eddy })
        })

        assert.equal(jobId, await rate.getJobId.call())
      })
    })
  })
})
