import * as h from '../src/helpersV2'
import { assertBigNum } from '../src/matchersV2'
import { ethers } from 'ethers'
import { Instance } from '../src/contract'
import { OracleFactory } from '../src/generated/OracleFactory'
import { LinkTokenFactory } from '../src/generated/LinkTokenFactory'
import { AggregatorFactory } from '../src/generated/AggregatorFactory'
import { assert } from 'chai'
import ganache from 'ganache-core'

const aggregatorFactory = new AggregatorFactory()
const oracleFactory = new OracleFactory()
const linkTokenFactory = new LinkTokenFactory()

let personas: h.Personas
let defaultAccount: ethers.Wallet
const provider = new ethers.providers.Web3Provider(ganache.provider() as any)

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)

  personas = rolesAndPersonas.personas
  defaultAccount = rolesAndPersonas.roles.defaultAccount
})

describe('Aggregator', () => {
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
  let link: Instance<LinkTokenFactory>
  let rate: Instance<AggregatorFactory>
  let oc1: Instance<OracleFactory>
  let oc2: Instance<OracleFactory>
  let oc3: Instance<OracleFactory>
  let oc4: Instance<OracleFactory>
  let oracles: Instance<OracleFactory>[]
  let jobIds: string[] = []
  const deployment = h.useSnapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    oc2 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    oc3 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    oc4 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    oracles = [oc1, oc2, oc3]
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(aggregatorFactory, [
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
      'transferOwnership',
    ])
  })

  describe('#requestRateUpdate', () => {
    const response = h.numToBytes32(100)

    describe('basic updates', () => {
      beforeEach(async () => {
        rate = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])

        await link.transfer(rate.address, deposit)

        const current = await rate.currentAnswer()
        assertBigNum(ethers.constants.Zero, current)
      })

      it('trigger a request to the oracle and accepts a response', async () => {
        const requestTx = await rate.requestRateUpdate()
        const receipt = await requestTx.wait()

        const log = receipt.logs![3]
        assert.equal(oc1.address, log.address)
        const request = h.decodeRunRequest(log)

        await h.fulfillOracleRequest(oc1, request, response)

        const current = await rate.currentAnswer()

        assertBigNum(response, current)
      })

      it('change the updatedAt record', async () => {
        let updatedAt = await rate.updatedHeight()
        assert.equal('0', updatedAt.toString())

        const requestTx = await rate.requestRateUpdate()
        const receipt = await requestTx.wait()
        const request = h.decodeRunRequest(receipt.logs![3])
        await h.fulfillOracleRequest(oc1, request, response)

        updatedAt = await rate.updatedHeight()
        assert.notEqual('0', updatedAt.toString())
      })

      it('emits a log with the response, answer ID, and sender', async () => {
        const requestTx = await rate.requestRateUpdate()
        const requestTxreceipt = await requestTx.wait()

        const request = h.decodeRunRequest(requestTxreceipt.logs![3])
        const fulfillOracleRequest = await h.fulfillOracleRequest(
          oc1,
          request,
          response,
        )
        const fulfillOracleRequestReceipt = await fulfillOracleRequest.wait()
        const answerId = h.numToBytes32(1)

        const receivedLog = fulfillOracleRequestReceipt.logs![1]
        assert.equal(response, receivedLog.topics[1])
        assert.equal(answerId, receivedLog.topics[2])
        assert.equal(
          oc1.address,
          ethers.utils.getAddress(receivedLog.topics[3].slice(26, 66)),
        )
      })

      it('emits a log with the new answer', async () => {
        const requestTx = await rate.requestRateUpdate()
        const requestReceipt = await requestTx.wait()

        const request = h.decodeRunRequest(requestReceipt.logs![3])
        const fulfillOracleRequest = await h.fulfillOracleRequest(
          oc1,
          request,
          response,
        )
        const fulfillOracleRequestReceipt = await fulfillOracleRequest.wait()

        const answerId = h.numToBytes32(1)
        const answerUpdatedLog = fulfillOracleRequestReceipt.logs![2]
        assert.equal(response, answerUpdatedLog.topics[1])

        assert.equal(answerId, answerUpdatedLog.topics[2])
      })
    })

    describe('with multiple oracles', () => {
      beforeEach(async () => {
        rate = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(
            link.address,
            basePayment,
            oracles.length,
            oracles.map(o => o.address),
            [jobId1, jobId2, jobId3],
          )

        await link.transfer(rate.address, deposit)

        const current = await rate.currentAnswer()
        assertBigNum(ethers.constants.Zero, current)
      })

      it('triggers requests to the oracles and the median of the responses', async () => {
        const requestTx = await rate.requestRateUpdate()
        const receipt = await requestTx.wait()

        const responses = [77, 66, 111].map(h.numToBytes32)

        for (let i = 0; i < oracles.length; i++) {
          const oracle = oracles[i]
          const log = receipt.logs![i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        const current = await rate.currentAnswer()
        assertBigNum(h.numToBytes32(77), current)
      })

      it('does not accept old responses', async () => {
        const request1 = await rate.requestRateUpdate()
        const receipt1 = await request1.wait()

        const response1 = h.numToBytes32(100)

        const requests = [
          h.decodeRunRequest(receipt1.logs![3]),
          h.decodeRunRequest(receipt1.logs![7]),
          h.decodeRunRequest(receipt1.logs![11]),
        ]

        const request2 = await rate.requestRateUpdate()
        const receipt2 = await request2.wait()
        const response2 = h.numToBytes32(200)

        for (let i = 0; i < oracles.length; i++) {
          const log = receipt2.logs![i * 4 + 3]
          const request = h.decodeRunRequest(log)
          await h.fulfillOracleRequest(oracles[i], request, response2)
        }
        assertBigNum(response2, await rate.currentAnswer())

        for (let i = 0; i < oracles.length; i++) {
          await h.fulfillOracleRequest(oracles[i], requests[i], response1)
        }

        assertBigNum(response2, await rate.currentAnswer())
      })
    })

    describe('with an even number of oracles', () => {
      beforeEach(async () => {
        oracles = [oc1, oc2, oc3, oc4]
        rate = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(
            link.address,
            basePayment,
            oracles.length,
            oracles.map(o => o.address),
            [jobId1, jobId2, jobId3, jobId4],
          )

        await link.transfer(rate.address, deposit)

        const current = await rate.currentAnswer()
        assertBigNum(ethers.constants.Zero, current)
      })

      it('triggers requests to the oracles and the median of the responses', async () => {
        const requestTx = await rate.requestRateUpdate()
        const receipt = await requestTx.wait()

        const responses = [66, 76, 78, 111].map(h.numToBytes32)

        for (let i = 0; i < oracles.length; i++) {
          const oracle = oracles[i]
          const log = receipt.logs![i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        const current = await rate.currentAnswer()
        assertBigNum(77, current)
      })
    })
  })

  describe('#updateRequestDetails', () => {
    beforeEach(async () => {
      rate = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
      await rate.transferOwnership(personas.Carol.address)
      oc2 = await oracleFactory.connect(defaultAccount).deploy(link.address)
      await link.transfer(rate.address, deposit)

      const current = await rate.currentAnswer()
      assertBigNum(ethers.constants.Zero, current)
    })

    describe('when called by the owner', () => {
      it('changes the amout of LINK sent on a request', async () => {
        const uniquePayment = 7777777
        await rate
          .connect(personas.Carol)
          .updateRequestDetails(uniquePayment, 1, [oc2.address], [jobId2])

        await rate.connect(personas.Carol).requestRateUpdate()

        assertBigNum(uniquePayment, await link.balanceOf(oc2.address))
      })

      it('can be configured to accept fewer responses than oracles', async () => {
        await rate
          .connect(personas.Carol)
          .updateRequestDetails(
            basePayment,
            1,
            [oc1.address, oc2.address],
            [jobId1, jobId2],
          )

        const requestTx = await rate.connect(personas.Carol).requestRateUpdate()
        const requestTxReceipt = await requestTx.wait()
        const request1 = h.decodeRunRequest(requestTxReceipt.logs![3])
        const request2 = h.decodeRunRequest(requestTxReceipt.logs![7])

        const response1 = h.numToBytes32(100)
        await h.fulfillOracleRequest(oc1, request1, response1)
        assertBigNum(response1, await rate.currentAnswer())

        const response2 = h.numToBytes32(200)
        await h.fulfillOracleRequest(oc2, request2, response2)

        const response1Bn = ethers.utils.bigNumberify(response1)
        const response2Bn = ethers.utils.bigNumberify(response2)
        const expected = response1Bn.add(response2Bn).div(2)

        assert.isTrue(expected.eq(await rate.currentAnswer()))
      })

      describe('and the number of jobs does not match number of oracles', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate
              .connect(personas.Carol)
              .updateRequestDetails(
                basePayment,
                2,
                [oc1.address, oc2.address],
                [jobId2],
              )
          })
        })
      })

      describe('and the oracles required exceeds the available amount', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate
              .connect(personas.Carol)
              .updateRequestDetails(
                basePayment,
                3,
                [oc1.address, oc2.address],
                [jobId1, jobId2],
              )
          })
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate
            .connect(personas.Eddy)
            .updateRequestDetails(basePayment, 1, [oc2.address], [jobId2])
        })
      })
    })

    describe('when called before a past answer is fulfilled', () => {
      beforeEach(async () => {
        rate = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
        await link.transfer(rate.address, deposit)

        oc2 = await oracleFactory.connect(defaultAccount).deploy(link.address)
        oc3 = await oracleFactory.connect(defaultAccount).deploy(link.address)
      })

      it('accepts answers from oracles at the time the request was made', async () => {
        // make request 1
        const request1Tx = await rate.requestRateUpdate()
        const request1Receipt = await request1Tx.wait()
        const request1 = h.decodeRunRequest(request1Receipt.logs![3])

        // change oracles
        await rate.updateRequestDetails(
          basePayment,
          2,
          [oc2.address, oc3.address],
          [jobId2, jobId3],
        )

        // make new request
        const request2Tx = await rate.requestRateUpdate()
        const request2Receipt = await request2Tx.wait()
        const request2 = h.decodeRunRequest(request2Receipt.logs![3])
        const request3 = h.decodeRunRequest(request2Receipt.logs![7])

        // fulfill request 1
        const response1 = h.numToBytes32(100)
        await h.fulfillOracleRequest(oc1, request1, response1)
        assertBigNum(response1, await rate.currentAnswer())

        // fulfill request 2
        const responses2 = [202, 222].map(h.numToBytes32)
        await h.fulfillOracleRequest(oc2, request2, responses2[0])
        await h.fulfillOracleRequest(oc3, request3, responses2[1])
        assertBigNum(212, await rate.currentAnswer())
      })
    })

    describe('when calling with a large number of oracles', () => {
      const maxOracleCount = 45

      beforeEach(() => {
        oracles = []
        jobIds = []
      })

      it(`does not revert with up to ${maxOracleCount} oracles`, async () => {
        for (let i = 0; i < maxOracleCount; i++) {
          oracles.push(oc1)
          jobIds.push(jobId1)
        }
        assert.equal(maxOracleCount, oracles.length)
        assert.equal(maxOracleCount, jobIds.length)

        await rate
          .connect(personas.Carol)
          .updateRequestDetails(
            basePayment,
            maxOracleCount,
            oracles.map(o => o.address),
            jobIds,
          )
      })

      it(`reverts with more than ${maxOracleCount} oracles`, async () => {
        const overMaxOracles = maxOracleCount + 1

        for (let i = 0; i < overMaxOracles; i++) {
          oracles.push(oc1)
          jobIds.push(jobId1)
        }
        assert.equal(overMaxOracles, oracles.length)
        assert.equal(overMaxOracles, jobIds.length)

        await h.assertActionThrows(async () => {
          await rate
            .connect(personas.Carol)
            .updateRequestDetails(
              basePayment,
              overMaxOracles,
              oracles.map(o => o.address),
              jobIds,
            )
        })
      })
    })
  })

  describe('#transferLINK', () => {
    beforeEach(async () => {
      rate = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
      await rate.transferOwnership(personas.Carol.address)
      await link.transfer(rate.address, deposit)
      assertBigNum(deposit, await link.balanceOf(rate.address))
    })

    describe('when called by the owner', () => {
      it('succeeds', async () => {
        await rate
          .connect(personas.Carol)
          .transferLINK(personas.Carol.address, deposit)

        assertBigNum(0, await link.balanceOf(rate.address))
        assertBigNum(deposit, await link.balanceOf(personas.Carol.address))
      })

      describe('with a number higher than the LINK balance', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await rate
              .connect(personas.Carol)
              .transferLINK(personas.Carol.address, deposit.add(basePayment))
          })

          assertBigNum(deposit, await link.balanceOf(rate.address))
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate
            .connect(personas.Eddy)
            .transferLINK(personas.Carol.address, deposit)
        })

        assertBigNum(deposit, await link.balanceOf(rate.address))
      })
    })
  })

  describe('#destroy', () => {
    beforeEach(async () => {
      rate = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
      await rate.transferOwnership(personas.Carol.address)
      await link.transfer(rate.address, deposit)
      assertBigNum(deposit, await link.balanceOf(rate.address))
    })

    describe('when called by the owner', () => {
      it('succeeds', async () => {
        await rate.connect(personas.Carol).destroy()

        assertBigNum(0, await link.balanceOf(rate.address))
        assertBigNum(deposit, await link.balanceOf(personas.Carol.address))

        assert.equal('0x', await provider.getCode(rate.address))
      })
    })

    describe('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.connect(personas.Eddy).destroy()
        })

        assertBigNum(deposit, await link.balanceOf(rate.address))
        assert.notEqual('0x', await provider.getCode(rate.address))
      })
    })
  })

  describe('#setAuthorization', () => {
    beforeEach(async () => {
      rate = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
      await link.transfer(rate.address, deposit)
    })

    describe('when called by an authorized address', () => {
      beforeEach(async () => {
        await rate.setAuthorization(personas.Eddy.address, true)
        assert.equal(
          true,
          await rate.authorizedRequesters(personas.Eddy.address),
        )
      })

      it('succeeds', async () => {
        await rate.connect(personas.Eddy).requestRateUpdate()
      })

      it('can be unset', async () => {
        await rate.setAuthorization(personas.Eddy.address, false)
        assert.equal(
          false,
          await rate.authorizedRequesters(personas.Eddy.address),
        )

        await h.assertActionThrows(async () => {
          await rate.connect(personas.Eddy).requestRateUpdate()
        })
      })
    })

    describe('when called by a non-authorized address', () => {
      beforeEach(async () => {
        assert.equal(
          false,
          await rate.authorizedRequesters(personas.Eddy.address),
        )
      })

      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await rate.connect(personas.Eddy).requestRateUpdate()
        })
      })
    })
  })

  describe('#cancelRequest', () => {
    let request: h.RunRequest

    beforeEach(async () => {
      rate = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])

      await link.transfer(rate.address, basePayment)

      assertBigNum(basePayment, await link.balanceOf(rate.address))
      assertBigNum(0, await link.balanceOf(oc1.address))

      const requestTx = await rate.requestRateUpdate()
      const receipt = await requestTx.wait()
      request = h.decodeRunRequest(receipt.logs![3])

      assertBigNum(0, await link.balanceOf(rate.address))
      assertBigNum(basePayment, await link.balanceOf(oc1.address))

      await h.increaseTime5Minutes(provider) // wait for request to expire
    })

    describe('when a later answer has been provided', () => {
      beforeEach(async () => {
        await link.transfer(rate.address, basePayment)
        const requestTx2 = await rate.requestRateUpdate()
        const receipt = await requestTx2.wait()
        const request2 = h.decodeRunRequest(receipt.logs![3])
        await h.fulfillOracleRequest(oc1, request2, '17')

        assertBigNum(basePayment.mul(2), await link.balanceOf(oc1.address))
      })

      it('gets the LINK deposited back from the oracle', async () => {
        await rate.cancelRequest(
          request.id,
          request.payment,
          request.expiration,
        )

        assertBigNum(basePayment, await link.balanceOf(rate.address))
        assertBigNum(basePayment, await link.balanceOf(oc1.address))
      })
    })

    describe('when a later answer has not been provided', () => {
      it('does not allow the request to be cancelled', async () => {
        h.assertActionThrows(async () => {
          await rate.cancelRequest(
            request.id,
            request.payment,
            request.expiration,
          )
        })

        assertBigNum(0, await link.balanceOf(rate.address))
        assertBigNum(basePayment, await link.balanceOf(oc1.address))
      })
    })
  })

  describe('testing various sets of inputs', () => {
    const tests = [
      {
        name: 'ordered ascending',
        responses: [0, 1, 2, 3, 4, 5, 6, 7].map(h.numToBytes32),
        want: h.numToBytes32(3),
      },
      {
        name: 'ordered descending',
        responses: [7, 6, 5, 4, 3, 2, 1, 0].map(h.numToBytes32),
        want: h.numToBytes32(3),
      },
      {
        name: 'unordered 1',
        responses: [1001, 1, 101, 10, 11, 0, 111].map(h.numToBytes32),
        want: h.numToBytes32(11),
      },
      {
        name: 'unordered 2',
        responses: [8, 8, 4, 5, 5, 7, 9, 5, 9].map(h.numToBytes32),
        want: h.numToBytes32(7),
      },
      {
        name: 'unordered 3',
        responses: [33, 44, 89, 101, 67, 7, 23, 55, 88, 324, 0, 88].map(
          h.numToBytes32,
        ),
        want: h.numToBytes32(61), // 67 + 55 / 2
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
          338678,
        ].map(h.numToBytes32),
        want: h.numToBytes32(335477),
      },
    ]

    beforeEach(async () => {
      rate = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 0, [], [])
      await link.transfer(rate.address, deposit)
    })

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const test of tests) {
      const responses = test.responses
      const oracles: Instance<OracleFactory>[] = []
      const jobIds: string[] = []

      it(test.name, async () => {
        for (let i = 0; i < responses.length; i++) {
          oracles[i] = await oracleFactory
            .connect(defaultAccount)
            .deploy(link.address)
          jobIds[i] = jobId1 // doesn't really matter in this test
        }

        await rate.updateRequestDetails(
          basePayment,
          oracles.length,
          oracles.map(o => o.address),
          jobIds,
        )

        const requestTx = await rate.requestRateUpdate()

        for (let i = 0; i < responses.length; i++) {
          const oracle = oracles[i]
          const receipt = await requestTx.wait()
          const log = receipt.logs![i * 4 + 3]
          assert.equal(oracle.address, log.address)
          const request = h.decodeRunRequest(log)

          await h.fulfillOracleRequest(oracle, request, responses[i])
        }

        assertBigNum(test.want, await rate.currentAnswer())
      })
    }
  })
})
