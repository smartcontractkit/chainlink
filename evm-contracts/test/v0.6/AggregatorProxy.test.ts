import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { BigNumber } from 'ethers/utils'
import { MockV2Aggregator__factory } from '../../ethers/v0.6/factories/MockV2Aggregator__factory'
import { MockV3Aggregator__factory } from '../../ethers/v0.6/factories/MockV3Aggregator__factory'
import { AggregatorProxy__factory } from '../../ethers/v0.6/factories/AggregatorProxy__factory'
import { AggregatorFacade__factory } from '../../ethers/v0.6/factories/AggregatorFacade__factory'
import { FluxAggregator__factory } from '../../ethers/v0.6/factories/FluxAggregator__factory'
import { Reverter__factory } from '../../ethers/v0.6/factories/Reverter__factory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkToken__factory()
const aggregatorFactory = new MockV3Aggregator__factory()
const historicAggregatorFactory = new MockV2Aggregator__factory()
const aggregatorFacadeFactory = new AggregatorFacade__factory()
const aggregatorProxyFactory = new AggregatorProxy__factory()
const fluxAggregatorFactory = new FluxAggregator__factory()
const reverterFactory = new Reverter__factory()

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('AggregatorProxy', () => {
  const deposit = h.toWei('100')
  const response = h.numToBytes32(54321)
  const response2 = h.numToBytes32(67890)
  const decimals = 18
  const phaseBase = h.bigNum(2).pow(64)

  let link: contract.Instance<contract.LinkToken__factory>
  let aggregator: contract.Instance<MockV3Aggregator__factory>
  let aggregator2: contract.Instance<MockV3Aggregator__factory>
  let historicAggregator: contract.Instance<MockV2Aggregator__factory>
  let proxy: contract.Instance<AggregatorProxy__factory>
  let flux: contract.Instance<FluxAggregator__factory>
  let reverter: contract.Instance<Reverter__factory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(decimals, response)
    await link.transfer(aggregator.address, deposit)
    proxy = await aggregatorProxyFactory
      .connect(defaultAccount)
      .deploy(aggregator.address)
    const emptyAddress = '0x0000000000000000000000000000000000000000'
    flux = await fluxAggregatorFactory
      .connect(personas.Carol)
      .deploy(link.address, 0, 0, emptyAddress, 0, 0, 18, 'TEST / LINK')
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(aggregatorProxyFactory, [
      'aggregator',
      'confirmAggregator',
      'decimals',
      'description',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'phaseAggregators',
      'phaseId',
      'proposeAggregator',
      'proposedAggregator',
      'proposedGetRoundData',
      'proposedLatestRoundData',
      'version',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('constructor', () => {
    it('sets the proxy phase and aggregator', async () => {
      matchers.bigNum(1, await proxy.phaseId())
      assert.equal(aggregator.address, await proxy.phaseAggregators(1))
    })
  })

  describe('#latestRound', () => {
    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(phaseBase.add(1), await proxy.latestRound())
    })
  })

  describe('#latestAnswer', () => {
    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(response, await proxy.latestAnswer())
      const latestRound = await proxy.latestRound()
      matchers.bigNum(response, await proxy.getAnswer(latestRound))
    })

    describe('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(decimals, response2)
        await link.transfer(aggregator2.address, deposit)
        matchers.bigNum(response2, await aggregator2.latestAnswer())

        await proxy.proposeAggregator(aggregator2.address)
        await proxy.confirmAggregator(aggregator2.address)
      })

      it('pulls the rate from the new aggregator', async () => {
        matchers.bigNum(response2, await proxy.latestAnswer())
        const latestRound = await proxy.latestRound()
        matchers.bigNum(response2, await proxy.getAnswer(latestRound))
      })
    })

    describe('when the relevant info is not available', () => {
      beforeEach(async () => {
        await proxy.proposeAggregator(flux.address)
        await proxy.confirmAggregator(flux.address)
      })

      it('does not revert when called with a non existent ID', async () => {
        const actual = await proxy.latestAnswer()
        matchers.bigNum(0, actual)
      })
    })
  })

  describe('#getAnswer', () => {
    describe('when the relevant round is not available', () => {
      beforeEach(async () => {
        await proxy.proposeAggregator(flux.address)
        await proxy.confirmAggregator(flux.address)
      })

      it('does not revert when called with a non existent ID', async () => {
        const proxyId = phaseBase.mul(await proxy.phaseId()).add(1)
        const actual = await proxy.getAnswer(proxyId)
        matchers.bigNum(0, actual)
      })
    })

    describe('when the answer reverts in a non-predicted way', () => {
      it('reverts', async () => {
        reverter = await reverterFactory.connect(defaultAccount).deploy()
        await proxy.proposeAggregator(reverter.address)
        await proxy.confirmAggregator(reverter.address)
        assert.equal(reverter.address, await proxy.aggregator())

        const proxyId = phaseBase.mul(await proxy.phaseId())

        await matchers.evmRevert(
          proxy.getAnswer(proxyId),
          'Raised by Reverter.sol',
        )
      })
    })

    describe('after being updated to another contract', () => {
      let preUpdateRoundId: BigNumber
      let preUpdateAnswer: BigNumber

      beforeEach(async () => {
        preUpdateRoundId = await proxy.latestRound()
        preUpdateAnswer = await proxy.latestAnswer()

        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(decimals, response2)
        await link.transfer(aggregator2.address, deposit)
        matchers.bigNum(response2, await aggregator2.latestAnswer())

        await proxy.proposeAggregator(aggregator2.address)
        await proxy.confirmAggregator(aggregator2.address)
      })

      it('reports answers for previous phases', async () => {
        const actualAnswer = await proxy.getAnswer(preUpdateRoundId)
        matchers.bigNum(preUpdateAnswer, actualAnswer)
      })
    })

    describe('when the relevant info is not available', () => {
      it('returns 0', async () => {
        const actual = await proxy.getAnswer(phaseBase.mul(777))
        matchers.bigNum(0, actual)
      })
    })

    describe('when the round ID is too large', () => {
      const overflowRoundId = h
        .bigNum(2)
        .pow(255)
        .add(phaseBase) // get the original phase
        .add(1) // get the original round
      it('returns 0', async () => {
        const actual = await proxy.getTimestamp(overflowRoundId)
        matchers.bigNum(0, actual)
      })
    })
  })

  describe('#getTimestamp', () => {
    describe('when the relevant round is not available', () => {
      beforeEach(async () => {
        await proxy.proposeAggregator(flux.address)
        await proxy.confirmAggregator(flux.address)
      })

      it('does not revert when called with a non existent ID', async () => {
        const proxyId = phaseBase.mul(await proxy.phaseId()).add(1)
        const actual = await proxy.getTimestamp(proxyId)
        matchers.bigNum(0, actual)
      })
    })

    describe('when the relevant info is not available', () => {
      it('returns 0', async () => {
        const actual = await proxy.getTimestamp(phaseBase.mul(777))
        matchers.bigNum(0, actual)
      })
    })

    describe('when the round ID is too large', () => {
      const overflowRoundId = h
        .bigNum(2)
        .pow(255)
        .add(phaseBase) // get the original phase
        .add(1) // get the original round

      it('returns 0', async () => {
        const actual = await proxy.getTimestamp(overflowRoundId)
        matchers.bigNum(0, actual)
      })
    })
  })

  describe('#latestTimestamp', () => {
    beforeEach(async () => {
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
        await h.increaseTimeBy(30, provider)
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(decimals, response2)

        const height2 = await aggregator2.latestTimestamp()
        assert.notEqual('0', height2.toString())

        const height1 = await aggregator.latestTimestamp()
        assert.notEqual(
          height1.toString(),
          height2.toString(),
          'Height1 and Height2 should not be equal',
        )

        await proxy.proposeAggregator(aggregator2.address)
        await proxy.confirmAggregator(aggregator2.address)
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
        historicAggregator = await historicAggregatorFactory
          .connect(defaultAccount)
          .deploy(response2)
        await proxy.proposeAggregator(historicAggregator.address)
        await proxy.confirmAggregator(historicAggregator.address)
      })

      it('reverts', async () => {
        const latestRoundId = await historicAggregator.latestRound()
        await matchers.evmRevert(proxy.getRoundData(latestRoundId))
      })

      describe('when pointed at an Aggregator Facade', () => {
        beforeEach(async () => {
          const facade = await aggregatorFacadeFactory
            .connect(defaultAccount)
            .deploy(aggregator.address, 18, 'LINK/USD: Aggregator Facade')
          await proxy.proposeAggregator(facade.address)
          await proxy.confirmAggregator(facade.address)
        })

        it('works for a valid roundId', async () => {
          const aggId = await aggregator.latestRound()
          const phaseId = phaseBase.mul(await proxy.phaseId())
          const proxyId = phaseId.add(aggId)

          const round = await proxy.getRoundData(proxyId)
          matchers.bigNum(proxyId, round.roundId)
          matchers.bigNum(response, round.answer)
          const nowSeconds = new Date().valueOf() / 1000
          assert.isAbove(round.updatedAt.toNumber(), nowSeconds - 120)
          matchers.bigNum(round.updatedAt, round.startedAt)
          matchers.bigNum(proxyId, round.answeredInRound)
        })
      })
    })

    describe('when pointed at a FluxAggregator', () => {
      beforeEach(async () => {
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(decimals, response2)

        await proxy.proposeAggregator(aggregator2.address)
        await proxy.confirmAggregator(aggregator2.address)
      })

      it('works for a valid round ID', async () => {
        const aggId = phaseBase.sub(2)
        await aggregator2
          .connect(personas.Carol)
          .updateRoundData(aggId, response2, 77, 42)

        const phaseId = phaseBase.mul(await proxy.phaseId())
        const proxyId = phaseId.add(aggId)

        const round = await proxy.getRoundData(proxyId)
        matchers.bigNum(proxyId, round.roundId)
        matchers.bigNum(response2, round.answer)
        matchers.bigNum(42, round.startedAt)
        matchers.bigNum(77, round.updatedAt)
        matchers.bigNum(proxyId, round.answeredInRound)
      })
    })

    it('reads round ID of a previous phase', async () => {
      const oldphaseId = phaseBase.mul(await proxy.phaseId())
      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, response2)

      await proxy.proposeAggregator(aggregator2.address)
      await proxy.confirmAggregator(aggregator2.address)

      const aggId = await aggregator.latestRound()
      const proxyId = oldphaseId.add(aggId)

      const round = await proxy.getRoundData(proxyId)
      matchers.bigNum(proxyId, round.roundId)
      matchers.bigNum(response, round.answer)
      const nowSeconds = new Date().valueOf() / 1000
      assert.isAbove(round.startedAt.toNumber(), nowSeconds - 120)
      assert.isBelow(round.startedAt.toNumber(), nowSeconds)
      matchers.bigNum(round.startedAt, round.updatedAt)
      matchers.bigNum(proxyId, round.answeredInRound)
    })
  })

  describe('#latestRoundData', () => {
    describe('when pointed at a Historic Aggregator', () => {
      beforeEach(async () => {
        historicAggregator = await historicAggregatorFactory
          .connect(defaultAccount)
          .deploy(response2)
        await proxy.proposeAggregator(historicAggregator.address)
        await proxy.confirmAggregator(historicAggregator.address)
      })

      it('reverts', async () => {
        await matchers.evmRevert(proxy.latestRoundData())
      })

      describe('when pointed at an Aggregator Facade', () => {
        beforeEach(async () => {
          const facade = await aggregatorFacadeFactory
            .connect(defaultAccount)
            .deploy(
              historicAggregator.address,
              17,
              'DOGE/ZWL: Aggregator Facade',
            )
          await proxy.proposeAggregator(facade.address)
          await proxy.confirmAggregator(facade.address)
        })

        it('does not revert', async () => {
          const aggId = await historicAggregator.latestRound()
          const phaseId = phaseBase.mul(await proxy.phaseId())
          const proxyId = phaseId.add(aggId)

          const round = await proxy.latestRoundData()
          matchers.bigNum(proxyId, round.roundId)
          matchers.bigNum(response2, round.answer)
          const nowSeconds = new Date().valueOf() / 1000
          assert.isAbove(round.updatedAt.toNumber(), nowSeconds - 120)
          matchers.bigNum(round.updatedAt, round.startedAt)
          matchers.bigNum(proxyId, round.answeredInRound)
        })

        it('uses the decimals set in the constructor', async () => {
          matchers.bigNum(17, await proxy.decimals())
        })

        it('uses the description set in the constructor', async () => {
          assert.equal('DOGE/ZWL: Aggregator Facade', await proxy.description())
        })

        it('sets the version to 2', async () => {
          matchers.bigNum(2, await proxy.version())
        })
      })
    })

    describe('when pointed at a FluxAggregator', () => {
      beforeEach(async () => {
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(decimals, response2)

        await proxy.proposeAggregator(aggregator2.address)
        await proxy.confirmAggregator(aggregator2.address)
      })

      it('does not revert', async () => {
        const aggId = phaseBase.sub(2)
        await aggregator2
          .connect(personas.Carol)
          .updateRoundData(aggId, response2, 77, 42)

        const phaseId = phaseBase.mul(await proxy.phaseId())
        const proxyId = phaseId.add(aggId)

        const round = await proxy.latestRoundData()
        matchers.bigNum(proxyId, round.roundId)
        matchers.bigNum(response2, round.answer)
        matchers.bigNum(42, round.startedAt)
        matchers.bigNum(77, round.updatedAt)
        matchers.bigNum(proxyId, round.answeredInRound)
      })

      it('uses the decimals of the aggregator', async () => {
        matchers.bigNum(18, await proxy.decimals())
      })

      it('uses the description of the aggregator', async () => {
        assert.equal(
          'v0.6/tests/MockV3Aggregator.sol',
          await proxy.description(),
        )
      })

      it('uses the version of the aggregator', async () => {
        matchers.bigNum(0, await proxy.version())
      })
    })
  })

  describe('#proposeAggregator', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol.address)
      await proxy.connect(personas.Carol).acceptOwnership()

      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, 1)

      assert.equal(aggregator.address, await proxy.aggregator())
    })

    describe('when called by the owner', () => {
      it('sets the address of the proposed aggregator', async () => {
        await proxy
          .connect(personas.Carol)
          .proposeAggregator(aggregator2.address)

        assert.equal(aggregator2.address, await proxy.proposedAggregator())
      })
    })

    describe('when called by a non-owner', () => {
      it('does not update', async () => {
        await matchers.evmRevert(
          proxy.connect(personas.Neil).proposeAggregator(aggregator2.address),
          'Only callable by owner',
        )

        assert.equal(aggregator.address, await proxy.aggregator())
      })
    })
  })

  describe('#confirmAggregator', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol.address)
      await proxy.connect(personas.Carol).acceptOwnership()

      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, 1)

      assert.equal(aggregator.address, await proxy.aggregator())
    })

    describe('when called by the owner', () => {
      beforeEach(async () => {
        await proxy
          .connect(personas.Carol)
          .proposeAggregator(aggregator2.address)
      })

      it('sets the address of the new aggregator', async () => {
        await proxy
          .connect(personas.Carol)
          .confirmAggregator(aggregator2.address)

        assert.equal(aggregator2.address, await proxy.aggregator())
      })

      it('increases the phase', async () => {
        matchers.bigNum(1, await proxy.phaseId())

        await proxy
          .connect(personas.Carol)
          .confirmAggregator(aggregator2.address)

        matchers.bigNum(2, await proxy.phaseId())
      })

      it('increases the round ID', async () => {
        matchers.bigNum(phaseBase.add(1), await proxy.latestRound())

        await proxy
          .connect(personas.Carol)
          .confirmAggregator(aggregator2.address)

        matchers.bigNum(phaseBase.mul(2).add(1), await proxy.latestRound())
      })

      it('sets the proxy phase and aggregator', async () => {
        assert.equal(
          '0x0000000000000000000000000000000000000000',
          await proxy.phaseAggregators(2),
        )

        await proxy
          .connect(personas.Carol)
          .confirmAggregator(aggregator2.address)

        assert.equal(aggregator2.address, await proxy.phaseAggregators(2))
      })
    })

    describe('when called by a non-owner', () => {
      beforeEach(async () => {
        await proxy
          .connect(personas.Carol)
          .proposeAggregator(aggregator2.address)
      })

      it('does not update', async () => {
        await matchers.evmRevert(
          proxy.connect(personas.Neil).confirmAggregator(aggregator2.address),
          'Only callable by owner',
        )

        assert.equal(aggregator.address, await proxy.aggregator())
      })
    })
  })

  describe('#proposedGetRoundData', () => {
    beforeEach(async () => {
      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, response2)
    })

    describe('when an aggregator has been proposed', () => {
      beforeEach(async () => {
        await proxy
          .connect(defaultAccount)
          .proposeAggregator(aggregator2.address)
        assert.equal(await proxy.proposedAggregator(), aggregator2.address)
      })

      it('returns the data for the proposed aggregator', async () => {
        const roundId = await aggregator2.latestRound()
        const round = await proxy.proposedGetRoundData(roundId)
        matchers.bigNum(roundId, round.roundId)
        matchers.bigNum(response2, round.answer)
      })

      describe('after the aggregator has been confirmed', () => {
        beforeEach(async () => {
          await proxy
            .connect(defaultAccount)
            .confirmAggregator(aggregator2.address)
          assert.equal(await proxy.aggregator(), aggregator2.address)
        })

        it('reverts', async () => {
          const roundId = await aggregator2.latestRound()
          await matchers.evmRevert(
            proxy.proposedGetRoundData(roundId),
            'No proposed aggregator present',
          )
        })
      })
    })
  })

  describe('#proposedLatestRoundData', () => {
    beforeEach(async () => {
      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, response2)
    })

    describe('when an aggregator has been proposed', () => {
      beforeEach(async () => {
        await proxy
          .connect(defaultAccount)
          .proposeAggregator(aggregator2.address)
        assert.equal(await proxy.proposedAggregator(), aggregator2.address)
      })

      it('returns the data for the proposed aggregator', async () => {
        const roundId = await aggregator2.latestRound()
        const round = await proxy.proposedLatestRoundData()
        matchers.bigNum(roundId, round.roundId)
        matchers.bigNum(response2, round.answer)
      })

      describe('after the aggregator has been confirmed', () => {
        beforeEach(async () => {
          await proxy
            .connect(defaultAccount)
            .confirmAggregator(aggregator2.address)
          assert.equal(await proxy.aggregator(), aggregator2.address)
        })

        it('reverts', async () => {
          await matchers.evmRevert(
            proxy.proposedLatestRoundData(),
            'No proposed aggregator present',
          )
        })
      })
    })
  })
})
