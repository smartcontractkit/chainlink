import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const personas = h.personas

contract.only('ConverstionRate', () => {
  const SOURCE_PATH: string = 'ConversionRate.sol'
  const jobId = '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  let link, oc, rate

  beforeEach(async () => {
    link = await h.linkContract()
    oc = await h.deploy('Oracle.sol', link.address)
    await oc.transferOwnership(personas.Neil, { from: personas.Default })
    rate = await h.deploy(SOURCE_PATH, link.address, [oc.address], [jobId])
    await rate.transferOwnership(personas.Carol, { from: personas.Default })
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(SOURCE_PATH), [
      'chainlinkCallback',
      'currentRate',
      'jobIds',
      'oracles',
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
})
