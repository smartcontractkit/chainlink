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
  const jobId4 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000004'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  let link, rate, oc1, oc2, oc3, oc4, oracles
  let jobIds = []

  beforeEach(async () => {
    link = await h.linkContract()
    oc1 = await h.deploy('Oracle.sol', link.address)
    oc2 = await h.deploy('Oracle.sol', link.address)
    oc3 = await h.deploy('Oracle.sol', link.address)
    oc4 = await h.deploy('Oracle.sol', link.address)
    oracles = [oc1, oc2, oc3]
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(SOURCE_PATH), [
      'authorizedRequesters',
      'cancelRequest',
      'chainlinkCallback',
      'currentAnswer',
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
      'updatedHeight',
      // Ownable methods:
      'owner',
      'renounceOwnership',
      'transferOwnership'
    ])
  })

  describe('#requestRateUpdate', () => {
    const response = 100

    context('basic updates', () => {
      beforeEach(async () => {
        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1]
        )

        await link.transfer(rate.address, deposit)

        const current = await rate.currentAnswer.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('trigger a request to the oracle and accepts a response', async () => {
        const requestTx = await rate.requestRateUpdate()

        const log = requestTx.receipt.rawLogs[3]
        assert.equal(oc1.address, log.address)
        const request = h.decodeRunRequest(log)

        await h.fulfillOracleRequest(oc1, request, response)

        const current = await rate.currentAnswer.call()
        assertBigNum(response, current)
      })

      it('change the updatedAt record', async () => {
        let updatedAt = await rate.updatedHeight.call()
        assert.equal('0', updatedAt.toString())

        const requestTx = await rate.requestRateUpdate()
        const request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
        await h.fulfillOracleRequest(oc1, request, response)

        updatedAt = await rate.updatedHeight.call()
        assert.notEqual('0', updatedAt.toString())
      })
    })

    context('with multiple oracles', () => {
      beforeEach(async () => {
        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          oracles.length,
          oracles.map(o => o.address),
          [jobId1, jobId2, jobId3]
        )

        await link.transfer(rate.address, deposit)

        const current = await rate.currentAnswer.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('triggers requests to the oracles and the median of the responses', async () => {
        const requestTx = await rate.requestRateUpdate()
        const responses = [77, 66, 111]

        for (let i = 0; i < oracles.length; i++) {
          const oracle = oracles[i]
          const log = requestTx.receipt.rawLogs[i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        const current = await rate.currentAnswer.call()
        assertBigNum(77, current)
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
        assertBigNum(response2, await rate.currentAnswer.call())

        for (let i = 0; i < oracles.length; i++) {
          await h.fulfillOracleRequest(oracles[i], requests[i], response1)
        }

        assertBigNum(response2, await rate.currentAnswer.call())
      })
    })

    context('with an even number of oracles', () => {
      beforeEach(async () => {
        oracles = [oc1, oc2, oc3, oc4]
        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          oracles.length,
          oracles.map(o => o.address),
          [jobId1, jobId2, jobId3, jobId4]
        )

        await link.transfer(rate.address, deposit)

        const current = await rate.currentAnswer.call()
        assertBigNum(h.bigNum(0), current)
      })

      it('triggers requests to the oracles and the median of the responses', async () => {
        const requestTx = await rate.requestRateUpdate()
        const responses = [66, 76, 78, 111]

        for (let i = 0; i < oracles.length; i++) {
          const oracle = oracles[i]
          const log = requestTx.receipt.rawLogs[i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        const current = await rate.currentAnswer.call()
        assertBigNum(77, current)
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
      await link.transfer(rate.address, deposit)

      const current = await rate.currentAnswer.call()
      assertBigNum(h.bigNum(0), current)
    })

    context('when called by the owner', () => {
      it('changes the amout of LINK sent on a request', async () => {
        const uniquePayment = 7777777
        await rate.updateRequestDetails(
          uniquePayment,
          1,
          [oc2.address],
          [jobId2],
          {
            from: personas.Carol
          }
        )

        await rate.requestRateUpdate({ from: personas.Carol })

        assertBigNum(uniquePayment, await link.balanceOf.call(oc2.address))
      })

      it('can be configured to accept fewer responses than oracles', async () => {
        await rate.updateRequestDetails(
          basePayment,
          1,
          [oc1.address, oc2.address],
          [jobId1, jobId2],
          {
            from: personas.Carol
          }
        )

        const requestTx = await rate.requestRateUpdate({ from: personas.Carol })
        const request1 = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
        const request2 = h.decodeRunRequest(requestTx.receipt.rawLogs[7])

        const response1 = 100
        await h.fulfillOracleRequest(oc1, request1, response1)
        assertBigNum(response1, await rate.currentAnswer.call())

        const response2 = 200
        await h.fulfillOracleRequest(oc2, request2, response2)
        assertBigNum(
          (response1 + response2) / 2,
          await rate.currentAnswer.call()
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

      context('and the oracles required exceeds the available amount', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate.updateRequestDetails(
              basePayment,
              3,
              [oc1.address, oc2.address],
              [jobId1, jobId2],
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

    context('when called before a past answer is fulfilled', () => {
      beforeEach(async () => {
        rate = await h.deploy(
          SOURCE_PATH,
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1]
        )
        await link.transfer(rate.address, deposit)

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
        assertBigNum(response1, await rate.currentAnswer.call())

        // fulfill request 2
        const responses2 = [202, 222]
        await h.fulfillOracleRequest(oc2, request2, responses2[0])
        await h.fulfillOracleRequest(oc3, request3, responses2[1])
        assertBigNum(212, await rate.currentAnswer.call())
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
      await link.transfer(rate.address, deposit)
      assertBigNum(deposit, await link.balanceOf.call(rate.address))
    })

    context('when called by the owner', () => {
      it('succeeds', async () => {
        await rate.transferLINK(personas.Carol, deposit, {
          from: personas.Carol
        })

        assertBigNum(0, await link.balanceOf.call(rate.address))
        assertBigNum(deposit, await link.balanceOf.call(personas.Carol))
      })

      context('with a number higher than the LINK balance', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate.transferLINK(personas.Carol, deposit.add(basePayment), {
              from: personas.Carol
            })
          })

          assertBigNum(deposit, await link.balanceOf.call(rate.address))
        })
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.transferLINK(personas.Carol, deposit, {
            from: personas.Eddy
          })
        })

        assertBigNum(deposit, await link.balanceOf.call(rate.address))
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
      await link.transfer(rate.address, deposit)
      assertBigNum(deposit, await link.balanceOf.call(rate.address))
    })

    context('when called by the owner', () => {
      it('succeeds', async () => {
        await rate.destroy({ from: personas.Carol })

        assertBigNum(0, await link.balanceOf.call(rate.address))
        assertBigNum(deposit, await link.balanceOf.call(personas.Carol))

        assert.equal('0x', await web3.eth.getCode(rate.address))
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.destroy({ from: personas.Eddy })
        })

        assertBigNum(deposit, await link.balanceOf.call(rate.address))
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
      await link.transfer(rate.address, deposit)
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

  describe('#cancelRequest', () => {
    let request

    beforeEach(async () => {
      rate = await h.deploy(
        SOURCE_PATH,
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1]
      )

      await link.transfer(rate.address, basePayment)

      assertBigNum(basePayment, await link.balanceOf.call(rate.address))
      assertBigNum(0, await link.balanceOf.call(oc1.address))

      const requestTx = await rate.requestRateUpdate()
      request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])

      assertBigNum(0, await link.balanceOf.call(rate.address))
      assertBigNum(basePayment, await link.balanceOf.call(oc1.address))

      await h.increaseTime5Minutes() // wait for request to expire
    })

    context('when a later answer has been provided', () => {
      beforeEach(async () => {
        await link.transfer(rate.address, basePayment)
        const requestTx2 = await rate.requestRateUpdate()
        const request2 = h.decodeRunRequest(requestTx2.receipt.rawLogs[3])
        await h.fulfillOracleRequest(oc1, request2, 17)

        assertBigNum(basePayment * 2, await link.balanceOf.call(oc1.address))
      })

      it('gets the LINK deposited back from the oracle', async () => {
        await rate.cancelRequest(
          request.id,
          request.payment,
          request.expiration
        )

        assertBigNum(basePayment, await link.balanceOf.call(rate.address))
        assertBigNum(basePayment, await link.balanceOf.call(oc1.address))
      })
    })

    context('when a later answer has not been provided', () => {
      it('does not allow the request to be cancelled', async () => {
        h.assertActionThrows(async () => {
          await rate.cancelRequest(
            request.id,
            request.payment,
            request.expiration
          )
        })

        assertBigNum(0, await link.balanceOf.call(rate.address))
        assertBigNum(basePayment, await link.balanceOf.call(oc1.address))
      })
    })
  })

  context('testing various sets of inputs', () => {
    const tests = [
      {
        name: 'ordered ascending',
        responses: [0, 1, 2, 3, 4, 5, 6, 7],
        want: 3
      },
      {
        name: 'ordered descending',
        responses: [7, 6, 5, 4, 3, 2, 1, 0],
        want: 3
      },
      {
        name: 'unordered 1',
        responses: [1001, 1, 101, 10, 11, 0, 111],
        want: 11
      },
      {
        name: 'unordered 2',
        responses: [8, 8, 4, 5, 5, 7, 9, 5, 9],
        want: 7
      },
      {
        name: 'unordered 3',
        responses: [33, 44, 89, 101, 67, 7, 23, 55, 88, 324, 0, 88],
        want: 61 // 67 + 55 / 2
      },
      {
        name: 'long unordered',
        responses: [
          333121,
          323453,
          337654,
          345363,
          345363,
          333456,
          335477,
          333323,
          332352,
          354648,
          983260,
          333856,
          335468,
          376987,
          333253,
          388867,
          337879,
          333324,
          338678
        ],
        want: 335477
      }
    ]

    beforeEach(async () => {
      rate = await h.deploy(SOURCE_PATH, link.address, basePayment, 0, [], [])
      await link.transfer(rate.address, deposit)
    })

    for (let test of tests) {
      const responses = test.responses
      let oracles = []
      let jobIds = []

      it(test.name, async () => {
        for (let i = 0; i < responses.length; i++) {
          oracles[i] = await h.deploy('Oracle.sol', link.address)
          jobIds[i] = jobId1 // doesn't really matter in this test
        }

        await rate.updateRequestDetails(
          basePayment,
          oracles.length,
          oracles.map(o => o.address),
          jobIds
        )

        const requestTx = await rate.requestRateUpdate()

        for (let i = 0; i < responses.length; i++) {
          const oracle = oracles[i]
          const log = requestTx.receipt.rawLogs[i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        assertBigNum(test.want, await rate.currentAnswer.call())
      })
    }
  })
})
