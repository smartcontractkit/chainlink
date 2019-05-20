import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const personas = h.personas

contract('ConverstionRate', () => {
  const SOURCE_PATH = 'ConversionRate.sol'
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const jobId2 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000002'
  const jobId3 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000003'
  const basePayment = h.toWei('1')
  let link, rate, oc1, oc2, oc3, oracles

  beforeEach(async () => {
    link = await h.linkContract()
    oc1 = await h.deploy('Oracle.sol', link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(SOURCE_PATH), [
      'authorizedRequesters',
      'chainlinkCallback',
      'currentRate',
      'destroy',
      'jobIds',
      'latestCompletedAnswer',
      'minimumResponses',
      'oracles',
      'paymentAmount',
      'requestRateUpdate',
      'setAuthorization',
      'transferLINK',
      'updateRequestDetails',
      // Ownable
      'owner',
      'renounceOwnership',
      'transferOwnership'
    ])
  })

  describe('#requestRateUpdate', () => {
    const response = 100

    context('with one oracle', () => {
      beforeEach(async () => {
        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1]
        )

        await link.transfer(rate.address, h.toWei('1'))

        const current = await rate.currentRate.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('triggeers a request to the oracle and accepts a response', async () => {
        const requestTx = await rate.requestRateUpdate()

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

        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          oracles.length,
          oracles.map(o => o.address),
          [jobId1, jobId2, jobId3]
        )

        await link.transfer(rate.address, h.toWei('100'))

        const current = await rate.currentRate.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('triggeers a request to the oracle and averages the responses', async () => {
        const requestTx = await rate.requestRateUpdate()
        const responses = [101, 102, 103]

        for (let i = 0; i < oracles.length; i++) {
          const oracle = oracles[i]
          const log = requestTx.receipt.rawLogs[i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        const current = await rate.currentRate.call()
        assertBigNum(102, current)
      })

      it('does not accept old responses', async () => {
        const request1 = await rate.requestRateUpdate()
        const response1 = 100

        const requests = [
          h.decodeRunRequest(request1.receipt.rawLogs[3]),
          h.decodeRunRequest(request1.receipt.rawLogs[7]),
          h.decodeRunRequest(request1.receipt.rawLogs[11])
        ]

        const request2 = await rate.requestRateUpdate()
        const response2 = 200
        for (let i = 0; i < oracles.length; i++) {
          const log = request2.receipt.rawLogs[i * 4 + 3]
          const request = h.decodeRunRequest(log)
          await h.fulfillOracleRequest(oracles[i], request, response2)
        }
        assertBigNum(response2, await rate.currentRate.call())

        for (let i = 0; i < oracles.length; i++) {
          await h.fulfillOracleRequest(oracles[i], requests[i], response1)
        }

        assertBigNum(response2, await rate.currentRate.call())
      })
    })
  })

  describe('#updateRequestDetails', () => {
    beforeEach(async () => {
      rate = await h.deploy(
        SOURCE_PATH,
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1]
      )
      await rate.transferOwnership(personas.Carol)

      oc2 = await h.deploy('Oracle.sol', link.address)

      await link.transfer(rate.address, h.toWei('100'))

      const current = await rate.currentRate.call()
      assertBigNum(h.bigNum(0), current)
    })

    context('when called by the owner', () => {
      it('succeeds', async () => {
        await rate.updateRequestDetails(
          basePayment,
          1,
          [oc2.address],
          [jobId2],
          {
            from: personas.Carol
          }
        )
      })

      context('and the number of jobs does not match number of oracles', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate.updateRequestDetails(
              basePayment,
              2,
              [oc1.address, oc2.address],
              [jobId2],
              {
                from: personas.Carol
              }
            )
          })
        })
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.updateRequestDetails(
            basePayment,
            1,
            [oc2.address],
            [jobId2],
            {
              from: personas.Eddy
            }
          )
        })
      })
    })

    context('when it is called before an answer is fulfilled', () => {
      beforeEach(async () => {
        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1]
        )
        await link.transfer(rate.address, h.toWei('100'))

        oc2 = await h.deploy('Oracle.sol', link.address)
        oc3 = await h.deploy('Oracle.sol', link.address)
      })

      it('accepts answers from oracles at the time the request was made', async () => {
        // make request 1
        const request1Tx = await rate.requestRateUpdate()
        const request1 = h.decodeRunRequest(request1Tx.receipt.rawLogs[3])

        // change oracles
        await rate.updateRequestDetails(
          basePayment,
          2,
          [oc2.address, oc3.address],
          [jobId2, jobId3]
        )

        // make new request
        const request2Tx = await rate.requestRateUpdate()
        const request2 = h.decodeRunRequest(request2Tx.receipt.rawLogs[3])
        const request3 = h.decodeRunRequest(request2Tx.receipt.rawLogs[7])

        // fulfill request 1
        const response1 = 100
        await h.fulfillOracleRequest(oc1, request1, response1)
        assertBigNum(response1, await rate.currentRate.call())

        // fulfill request 2
        const response2 = 200
        await h.fulfillOracleRequest(oc2, request2, response2)
        await h.fulfillOracleRequest(oc3, request3, response2)
        assertBigNum(response2, await rate.currentRate.call())
      })
    })
  })

  describe('#transferLINK', () => {
    beforeEach(async () => {
      rate = await h.deploy(
        SOURCE_PATH,
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1]
      )
      await rate.transferOwnership(personas.Carol)
      await link.transfer(rate.address, h.toWei('100'))
      assertBigNum(h.toWei('100'), await link.balanceOf.call(rate.address))
    })

    context('when called by the owner', () => {
      it('succeeds', async () => {
        await rate.transferLINK(personas.Carol, h.toWei('100'), {
          from: personas.Carol
        })

        assertBigNum(h.toWei('0'), await link.balanceOf.call(rate.address))
        assertBigNum(h.toWei('100'), await link.balanceOf.call(personas.Carol))
      })

      context('with a number higher than the LINK balance', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate.transferLINK(personas.Carol, h.toWei('101'), {
              from: personas.Carol
            })
          })

          assertBigNum(h.toWei('100'), await link.balanceOf.call(rate.address))
        })
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.transferLINK(personas.Carol, h.toWei('100'), {
            from: personas.Eddy
          })
        })

        assertBigNum(h.toWei('100'), await link.balanceOf.call(rate.address))
      })
    })
  })

  describe('#destroy', () => {
    beforeEach(async () => {
      rate = await h.deploy(
        SOURCE_PATH,
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1]
      )
      await rate.transferOwnership(personas.Carol)
      await link.transfer(rate.address, h.toWei('100'))
      assertBigNum(h.toWei('100'), await link.balanceOf.call(rate.address))
    })

    context('when called by the owner', () => {
      it('succeeds', async () => {
        await rate.destroy({ from: personas.Carol })

        assertBigNum(h.toWei('0'), await link.balanceOf.call(rate.address))
        assertBigNum(h.toWei('100'), await link.balanceOf.call(personas.Carol))

        assert.equal('0x', await web3.eth.getCode(rate.address))
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.destroy({ from: personas.Eddy })
        })

        assertBigNum(h.toWei('100'), await link.balanceOf.call(rate.address))
        assert.notEqual('0x', await web3.eth.getCode(rate.address))
      })
    })
  })

  describe('#setAuthorization', async () => {
    beforeEach(async () => {
      rate = await h.deploy(
        SOURCE_PATH,
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1]
      )
      await link.transfer(rate.address, h.toWei('100'))
    })

    context('when called by an authorized address', () => {
      beforeEach(async () => {
        await rate.setAuthorization(personas.Eddy, true)
        assert.equal(true, await rate.authorizedRequesters.call(personas.Eddy))
      })

      it('succeeds', async () => {
        await rate.requestRateUpdate({ from: personas.Eddy })
      })

      it('can be unset', async () => {
        await rate.setAuthorization(personas.Eddy, false)
        assert.equal(false, await rate.authorizedRequesters.call(personas.Eddy))

        await h.assertActionThrows(async () => {
          await rate.requestRateUpdate({ from: personas.Eddy })
        })
      })
    })

    context('when called by a non-authorized address', () => {
      beforeEach(async () => {
        assert.equal(false, await rate.authorizedRequesters.call(personas.Eddy))
      })

      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.requestRateUpdate({ from: personas.Eddy })
        })
      })
    })
  })
})
