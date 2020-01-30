import {
  contract,
  coordinator,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/eth-test-helpers'
import cbor from 'cbor'
import { assert } from 'chai'
import { ethers } from 'ethers'
import {
  CoordinatorFactory,
  MeanAggregatorFactory,
  ServiceAgreementConsumerFactory,
} from '../src/generated'

const coordinatorFactory = new CoordinatorFactory()
const meanAggregatorFactory = new MeanAggregatorFactory()
const serviceAgreementConsumerFactory = new ServiceAgreementConsumerFactory()
const linkTokenFactory = new contract.LinkTokenFactory()

// create ethers provider from that web3js instance
const provider = setup.provider()

let roles: setup.Roles

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('ServiceAgreementConsumer', () => {
  const currency = 'USD'

  let link: contract.Instance<contract.LinkTokenFactory>
  let coord: contract.Instance<CoordinatorFactory>
  let cc: contract.Instance<ServiceAgreementConsumerFactory>
  let agreement: coordinator.ServiceAgreement

  beforeEach(async () => {
    const meanAggregator = await meanAggregatorFactory
      .connect(roles.defaultAccount)
      .deploy()
    agreement = await coordinator.newServiceAgreement({
      aggregator: meanAggregator.address,
      oracles: [roles.oracleNode],
    })
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    coord = await coordinatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
    await coordinator.initiateServiceAgreement(coord, agreement)
    cc = await serviceAgreementConsumerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, coord.address, coordinator.generateSAID(agreement))
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
        await matchers.evmRevert(async () => {
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

        const request = oracle.decodeRunRequest(log)

        assert.equal(coordinator.generateSAID(agreement), request.jobId)
        matchers.bigNum(paymentAmount, request.payment)
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
      let request: oracle.RunRequest

      beforeEach(async () => {
        await link.transfer(cc.address, h.toWei('1'))
        const tx = await cc.requestEthereumPrice(currency)
        const receipt = await tx.wait()
        const log = receipt?.logs?.[3]
        assert.equal(log?.address.toLowerCase(), coord.address.toLowerCase())

        request = oracle.decodeRunRequest(log)
      })

      it('records the data given to it by the oracle', async () => {
        await coord
          .connect(roles.oracleNode)
          .fulfillOracleRequest(request.id, response)
        const currentPrice = await cc.currentPrice()
        assert.equal(currentPrice, response)
      })

      describe('when the consumer does not recognize the request ID', () => {
        let request2: oracle.RunRequest

        beforeEach(async () => {
          // Create a request directly via the oracle, rather than through the
          // chainlink client (consumer). The client should not respond to
          // fulfillment of this request, even though the oracle will faithfully
          // forward the fulfillment to it.
          const args = oracle.encodeOracleRequest(
            coordinator.generateSAID(agreement),
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

          request2 = oracle.decodeRunRequest(receipt?.logs?.[2])
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
          await matchers.evmRevert(async () => {
            await cc.connect(roles.oracleNode).fulfill(request.id, response)
          })
          const received = await cc.currentPrice()
          assert.equal(ethers.utils.parseBytes32String(received), '')
        })
      })
    })
  })
})
