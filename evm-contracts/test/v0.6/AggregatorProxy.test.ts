import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { AggregatorFactory } from '../../ethers/v0.4/AggregatorFactory'
import { AggregatorProxyFactory } from '../../ethers/v0.6/AggregatorProxyFactory'
import { AggregatorFacadeFactory } from '../../ethers/v0.6/AggregatorFacadeFactory'
import { OracleFactory } from '../../ethers/v0.6/OracleFactory'
import { FluxAggregatorFactory } from '../../ethers/v0.6/FluxAggregatorFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const aggregatorFactory = new AggregatorFactory()
const aggregatorFacadeFactory = new AggregatorFacadeFactory()
const oracleFactory = new OracleFactory()
const aggregatorProxyFactory = new AggregatorProxyFactory()
const fluxAggregatorFactory = new FluxAggregatorFactory()

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('AggregatorProxy', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  const response = h.numToBytes32(54321)
  const response2 = h.numToBytes32(67890)

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.CallableOverrideInstance<AggregatorFactory>
  let aggregator2: contract.CallableOverrideInstance<AggregatorFactory>
  let oc1: contract.Instance<OracleFactory>
  let proxy: contract.CallableOverrideInstance<AggregatorProxyFactory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = contract.callableAggregator(
      await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
    )
    await link.transfer(aggregator.address, deposit)
    proxy = contract.callableAggregator(
      await aggregatorProxyFactory
        .connect(defaultAccount)
        .deploy(aggregator.address),
    )
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(aggregatorProxyFactory, [
      'aggregator',
      'decimals',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'setAggregator',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#latestAnswer', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const receipt = await requestTx.wait()

      const request = oracle.decodeRunRequest(receipt.logs?.[3])
      await oc1.fulfillOracleRequest(
        ...oracle.convertFufillParams(request, response),
      )
      matchers.bigNum(
        ethers.utils.bigNumberify(response),
        await aggregator.latestAnswer(),
      )
    })

    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(response, await proxy.latestAnswer())
      const latestRound = await proxy.latestRound()
      matchers.bigNum(response, await proxy.getAnswer(latestRound))
    })

    describe('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = contract.callableAggregator(
          await aggregatorFactory
            .connect(defaultAccount)
            .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
        )
        await link.transfer(aggregator2.address, deposit)
        const requestTx = await aggregator2.requestRateUpdate()
        const receipt = await requestTx.wait()
        const request = oracle.decodeRunRequest(receipt.logs?.[3])

        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, response2),
        )
        matchers.bigNum(response2, await aggregator2.latestAnswer())

        await proxy.setAggregator(aggregator2.address)
      })

      it('pulls the rate from the new aggregator', async () => {
        matchers.bigNum(response2, await proxy.latestAnswer())
        const latestRound = await proxy.latestRound()
        matchers.bigNum(response2, await proxy.getAnswer(latestRound))
      })
    })
  })

  describe('#latestTimestamp', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const receipt = await requestTx.wait()
      const request = oracle.decodeRunRequest(receipt.logs?.[3])

      await oc1.fulfillOracleRequest(
        ...oracle.convertFufillParams(request, response),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the timestamp from the aggregator', async () => {
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.latestTimestamp(),
      )
      const latestRound = await proxy.latestRound()
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.getTimestamp(latestRound),
      )
    })

    describe('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = contract.callableAggregator(
          await aggregatorFactory
            .connect(defaultAccount)
            .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
        )
        await link.transfer(aggregator2.address, deposit)

        const requestTx = await aggregator2.requestRateUpdate()
        const receipt = await requestTx.wait()
        const request = oracle.decodeRunRequest(receipt.logs?.[3])

        await h.increaseTimeBy(30, provider)
        await h.mineBlock(provider)

        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, response2),
        )
        const height2 = await aggregator2.latestTimestamp()
        assert.notEqual('0', height2.toString())

        const height1 = await aggregator.latestTimestamp()
        assert.notEqual(
          height1.toString(),
          height2.toString(),
          'Height1 and Height2 should not be equal',
        )

        await proxy.setAggregator(aggregator2.address)
      })

      it('pulls the timestamp from the new aggregator', async () => {
        matchers.bigNum(
          await aggregator2.latestTimestamp(),
          await proxy.latestTimestamp(),
        )
        const latestRound = await proxy.latestRound()
        matchers.bigNum(
          await aggregator2.latestTimestamp(),
          await proxy.getTimestamp(latestRound),
        )
      })
    })
  })

  describe('#getRoundData', () => {
    describe('when pointed at a Historic Aggregator', () => {
      beforeEach(async () => {
        const requestTx = await aggregator.requestRateUpdate()
        const receipt = await requestTx.wait()

        const request = oracle.decodeRunRequest(receipt.logs?.[3])
        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, response),
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(response),
          await aggregator.latestAnswer(),
        )
      })

      it('reverts', async () => {
        const latestRoundId = await proxy.latestRound()
        await matchers.evmRevert(async () => {
          await proxy.getRoundData(latestRoundId)
        })
      })

      describe('when pointed at an Aggregator Facade', () => {
        beforeEach(async () => {
          const facade = await aggregatorFacadeFactory
            .connect(defaultAccount)
            .deploy(aggregator.address, 18)
          await proxy.setAggregator(facade.address)
        })

        it('works for a valid roundId', async () => {
          const roundId = await aggregator.latestRound()
          const round = await proxy.getRoundData(roundId)
          matchers.bigNum(roundId, round.roundId)
          matchers.bigNum(response, round.answer)
          const nowSeconds = new Date().valueOf() / 1000
          assert.isAbove(round.updatedAt.toNumber(), nowSeconds - 120)
          matchers.bigNum(round.updatedAt, round.startedAt)
          matchers.bigNum(roundId, round.answeredInRound)
        })
      })
    })

    describe('when pointed at a FluxAggregator', () => {
      const roundId = 1
      const submission = 42
      beforeEach(async () => {
        const fluxAggregator = await fluxAggregatorFactory
          .connect(defaultAccount)
          .deploy(
            link.address,
            basePayment,
            3600,
            18,
            ethers.utils.formatBytes32String('DOGE/ZWL'),
          )
        await link.transferAndCall(fluxAggregator.address, deposit, [])
        await fluxAggregator.addOracles(
          [defaultAccount.address],
          [defaultAccount.address],
          1,
          1,
          0,
        )
        await fluxAggregator.submit(roundId, submission)

        await proxy.setAggregator(fluxAggregator.address)
      })

      it('works for a valid round ID', async () => {
        const round = await proxy.getRoundData(roundId)
        matchers.bigNum(roundId, round.roundId)
        matchers.bigNum(submission, round.answer)
        const nowSeconds = new Date().valueOf() / 1000
        assert.isAbove(round.startedAt.toNumber(), nowSeconds - 120)
        assert.isBelow(round.startedAt.toNumber(), nowSeconds)
        matchers.bigNum(round.startedAt, round.updatedAt)
        matchers.bigNum(roundId, round.answeredInRound)
      })
    })
  })

  describe('#latestRoundData', () => {
    describe('when pointed at a Historic Aggregator', () => {
      beforeEach(async () => {
        const requestTx = await aggregator.requestRateUpdate()
        const receipt = await requestTx.wait()

        const request = oracle.decodeRunRequest(receipt.logs?.[3])
        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, response),
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(response),
          await aggregator.latestAnswer(),
        )
      })

      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await proxy.latestRoundData()
        })
      })

      describe('when pointed at an Aggregator Facade', () => {
        beforeEach(async () => {
          const facade = await aggregatorFacadeFactory
            .connect(defaultAccount)
            .deploy(aggregator.address, 18)
          await proxy.setAggregator(facade.address)
        })

        it('does not revert', async () => {
          const roundId = await aggregator.latestRound()
          const round = await proxy.latestRoundData()
          matchers.bigNum(roundId, round.roundId)
          matchers.bigNum(response, round.answer)
          const nowSeconds = new Date().valueOf() / 1000
          assert.isAbove(round.updatedAt.toNumber(), nowSeconds - 120)
          matchers.bigNum(round.updatedAt, round.startedAt)
          matchers.bigNum(roundId, round.answeredInRound)
        })
      })
    })

    describe('when pointed at a FluxAggregator', () => {
      const roundId = 1
      const submission = 42
      beforeEach(async () => {
        const fluxAggregator = await fluxAggregatorFactory
          .connect(defaultAccount)
          .deploy(
            link.address,
            basePayment,
            3600,
            18,
            ethers.utils.formatBytes32String('DOGE/ZWL'),
          )
        await link.transferAndCall(fluxAggregator.address, deposit, [])
        await fluxAggregator.addOracles(
          [defaultAccount.address],
          [defaultAccount.address],
          1,
          1,
          0,
        )
        await fluxAggregator.submit(roundId, submission)

        await proxy.setAggregator(fluxAggregator.address)
      })

      it('does not revert', async () => {
        const round = await proxy.latestRoundData()
        matchers.bigNum(roundId, round.roundId)
        matchers.bigNum(submission, round.answer)
        const nowSeconds = new Date().valueOf() / 1000
        assert.isAbove(round.startedAt.toNumber(), nowSeconds - 120)
        assert.isBelow(round.startedAt.toNumber(), nowSeconds)
        matchers.bigNum(round.startedAt, round.updatedAt)
        matchers.bigNum(roundId, round.answeredInRound)
      })
    })
  })

  describe('#setAggregator', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol.address)
      await proxy.connect(personas.Carol).acceptOwnership()

      aggregator2 = contract.callableAggregator(
        await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
      )

      assert.equal(aggregator.address, await proxy.aggregator())
    })

    describe('when called by the owner', () => {
      it('sets the address of the new aggregator', async () => {
        await proxy.connect(personas.Carol).setAggregator(aggregator2.address)

        assert.equal(aggregator2.address, await proxy.aggregator())
      })
    })

    describe('when called by a non-owner', () => {
      it('does not update', async () => {
        await matchers.evmRevert(async () => {
          await proxy.connect(personas.Neil).setAggregator(aggregator2.address)
        })

        assert.equal(aggregator.address, await proxy.aggregator())
      })
    })
  })
})
