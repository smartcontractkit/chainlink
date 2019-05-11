import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const personas = h.personas

contract('ConverstionRate', () => {
  const SOURCE_PATH: string = 'ConversionRate.sol'
  const jobId1 = '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const jobId2 = '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000002'
  const jobId3 = '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000003'
  let link, rate, oc1, oc2, oc3, oracles

  beforeEach(async () => {
    link = await h.linkContract()
    oc1 = await h.deploy('Oracle.sol', link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(SOURCE_PATH), [
      'chainlinkCallback',
      'currentRate',
      'jobIds',
      'oracles',
      'update',
    ])
  })

  describe('#update', () => {
    const response = 100

    context('with one oracle', () => {
      beforeEach(async () => {
        rate = await h.deploy(SOURCE_PATH, link.address, [oc1.address], [jobId1])

        await link.transfer(rate.address, h.toWei('1', 'ether'))

        const current = await rate.currentRate.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('triggeers a request to the oracle and accepts a response', async () => {
        const requestTx = await rate.update()

        const log = requestTx.receipt.rawLogs[3]
        assert.equal(oc1.address, log.address)
        const request = h.decodeRunRequest(log)

        await h.fulfillOracleRequest(oc1, request, response)

        const current = await rate.currentRate.call()
        assertBigNum(response, current)
      })
    })

    context('with multiple oracles', () => {
      beforeEach(async () => {
        oc2 = await h.deploy('Oracle.sol', link.address)
        oc3 = await h.deploy('Oracle.sol', link.address)
        oracles = [oc1, oc2, oc3]

        rate = await h.deploy(SOURCE_PATH, link.address, oracles.map(o => o.address), [jobId1, jobId2, jobId3])

        await link.transfer(rate.address, h.toWei('3', 'ether'))

        const current = await rate.currentRate.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('triggeers a request to the oracle and accepts a response', async () => {
        const requestTx = await rate.update()
        const responses = [101, 102, 103]

        for (let i = 0; i < oracles.length; i++) {
            const oracle = oracles[i]
            const log = requestTx.receipt.rawLogs[(i * 4) + 3]
            assert.equal(oracle.address, log.address)
            const request = h.decodeRunRequest(log)

            await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        const current = await rate.currentRate.call()
        assertBigNum(102, current)
      })
    })
  })
})
