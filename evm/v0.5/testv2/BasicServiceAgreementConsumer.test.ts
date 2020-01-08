import * as h from '../src/helpers'
import cbor from 'cbor'
import { assertBigNum } from '../src/matchers'
import { makeTestProvider } from '../src/provider'
import {
  CoordinatorFactory,
  MeanAggregatorFactory,
  ServiceAgreementConsumerFactory,
  LinkTokenFactory,
} from '../src/generated'
import { Instance } from '../src/contract'
import { ethers } from 'ethers'
import { assert } from 'chai'

const coordinatorFactory = new CoordinatorFactory()
const meanAggregatorFactory = new MeanAggregatorFactory()
const serviceAgreementConsumerFactory = new ServiceAgreementConsumerFactory()
const linkTokenFactory = new LinkTokenFactory()

// create ethers provider from that web3js instance
const provider = makeTestProvider()

let roles: h.Roles

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)

  roles = rolesAndPersonas.roles
})

describe('ServiceAgreementConsumer', () => {
  const currency = 'USD'

  let link: Instance<LinkTokenFactory>
  let coord: Instance<CoordinatorFactory>
  let cc: Instance<ServiceAgreementConsumerFactory>
  let agreement: h.ServiceAgreement

  beforeEach(async () => {
    const meanAggregator = await meanAggregatorFactory
      .connect(roles.defaultAccount)
      .deploy()

    agreement = {
      oracles: [roles.oracleNode.address],
      aggregator: meanAggregator.address,
      aggInitiateJobSelector:
        meanAggregator.interface.functions.initiateJob.sighash,
      aggFulfillSelector: meanAggregator.interface.functions.fulfill.sighash,
      endAt: h.sixMonthsFromNow(),
      expiration: ethers.utils.bigNumberify(300),
      payment: ethers.utils.bigNumberify('1000000000000000000'),
      requestDigest:
        '0xbadc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5badc0de5',
    }
    const oracleSignatures = await h.computeOracleSignature(
      agreement,
      roles.oracleNode,
    )

    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    coord = await coordinatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)

    await coord.initiateServiceAgreement(
      h.encodeServiceAgreement(agreement),
      h.encodeOracleSignatures(oracleSignatures),
    )
    cc = await serviceAgreementConsumerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, coord.address, h.generateSAID(agreement))
  })

  it('gas price of contract deployment is predictable', async () => {
    const rec = await provider.getTransactionReceipt(
      cc.deployTransaction.hash ?? '',
    )
    assert.isBelow(rec.gasUsed?.toNumber() ?? 0, 1500000)
  })

  describe('#requestEthereumPrice', () => {
    describe('without LINK', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await cc.requestEthereumPrice(currency)
        })
      })
    })

    describe('with LINK', () => {
      const paymentAmount = h.toWei('1')
      beforeEach(async () => {
        await link.transfer(cc.address, paymentAmount)
      })

      it('triggers a log event in the Coordinator contract', async () => {
        const tx = await cc.requestEthereumPrice(currency)
        const receipt = await tx.wait()
        const log = receipt?.logs?.[3]
        assert.equal(log?.address.toLowerCase(), coord.address.toLowerCase())

        const request = h.decodeRunRequest(log)

        assert.equal(h.generateSAID(agreement), request.jobId)
        assertBigNum(paymentAmount, request.payment)
        assert.equal(cc.address.toLowerCase(), request.requester.toLowerCase())
        assert.equal(1, request.dataVersion)

        const expected = {
          path: currency,
          get:
            'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY',
        }
        assert.deepEqual(expected, cbor.decodeFirstSync(request.data))
      })

      it('has a reasonable gas cost', async () => {
        const tx = await cc.requestEthereumPrice(currency)
        const receipt = await tx.wait()

        assert.isBelow(receipt?.gasUsed?.toNumber() ?? -1, 175000)
      })
    })

    describe('#fulfillOracleRequest', () => {
      const response = ethers.utils.formatBytes32String('1,000,000.00')
      let request: h.RunRequest

      beforeEach(async () => {
        await link.transfer(cc.address, h.toWei('1'))
        const tx = await cc.requestEthereumPrice(currency)
        const receipt = await tx.wait()
        const log = receipt?.logs?.[3]
        assert.equal(log?.address.toLowerCase(), coord.address.toLowerCase())

        request = h.decodeRunRequest(log)
      })

      it('records the data given to it by the oracle', async () => {
        await coord
          .connect(roles.oracleNode)
          .fulfillOracleRequest(request.id, response)
        const currentPrice = await cc.currentPrice()
        assert.equal(currentPrice, response)
      })

      describe('when the consumer does not recognize the request ID', () => {
        let request2: h.RunRequest

        beforeEach(async () => {
          // Create a request directly via the oracle, rather than through the
          // chainlink client (consumer). The client should not respond to
          // fulfillment of this request, even though the oracle will faithfully
          // forward the fulfillment to it.
          const args = h.requestDataBytes(
            h.generateSAID(agreement),
            cc.address,
            serviceAgreementConsumerFactory.interface.functions.fulfill.sighash,
            48,
            '0x0',
          )

          const tx = await link.transferAndCall(
            coord.address,
            agreement.payment,
            args,
          )
          const receipt = await tx.wait()

          request2 = h.decodeRunRequest(receipt?.logs?.[2])
        })

        it('does not accept the data provided', async () => {
          await coord
            .connect(roles.oracleNode)
            .fulfillOracleRequest(request2.id, response)

          const received = await cc.currentPrice()
          assert.equal(ethers.utils.parseBytes32String(received), '')
        })
      })

      describe('when called by anyone other than the oracle contract', () => {
        it('does not accept the data provided', async () => {
          await h.assertActionThrows(async () => {
            await cc.connect(roles.oracleNode).fulfill(request.id, response)
          })
          const received = await cc.currentPrice()
          assert.equal(ethers.utils.parseBytes32String(received), '')
        })
      })
    })
  })
})
