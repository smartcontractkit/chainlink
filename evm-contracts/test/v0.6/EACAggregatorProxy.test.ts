import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { EACAggregatorProxyFactory } from '../../ethers/v0.6/EACAggregatorProxyFactory'
import { AccessControlledAggregatorFactory } from '../../ethers/v0.6/AccessControlledAggregatorFactory'
import { SimpleReadAccessControllerFactory } from '../../ethers/v0.6/SimpleReadAccessControllerFactory'
import { MockV3AggregatorFactory } from '../../ethers/v0.6/MockV3AggregatorFactory'
import { FluxAggregatorTestHelperFactory } from '../../ethers/v0.6/FluxAggregatorTestHelperFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const accessControlFactory = new SimpleReadAccessControllerFactory()
const aggregatorFactory = new MockV3AggregatorFactory()
const acAggregatorFactory = new AccessControlledAggregatorFactory()
const testHelperFactory = new FluxAggregatorTestHelperFactory()
const proxyFactory = new EACAggregatorProxyFactory()
const emptyAddress = '0x0000000000000000000000000000000000000000'

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('EACAggregatorProxy', () => {
  const deposit = h.toWei('100')
  const answer = h.numToBytes32(54321)
  const answer2 = h.numToBytes32(54320)
  const roundId = 17
  const decimals = 18
  const timestamp = 678
  const startedAt = 677

  let link: contract.Instance<contract.LinkTokenFactory>
  let controller: contract.Instance<SimpleReadAccessControllerFactory>
  let aggregator: contract.Instance<MockV3AggregatorFactory>
  let aggregator2: contract.Instance<MockV3AggregatorFactory>
  let proxy: contract.Instance<EACAggregatorProxyFactory>
  let testHelper: contract.Instance<FluxAggregatorTestHelperFactory>
  const phaseBase = h.bigNum(2).pow(64)

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(decimals, 0)
    controller = await accessControlFactory.connect(defaultAccount).deploy()
    await aggregator.updateRoundData(roundId, answer, timestamp, startedAt)
    await link.transfer(aggregator.address, deposit)
    proxy = await proxyFactory
      .connect(defaultAccount)
      .deploy(aggregator.address, controller.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(proxyFactory, [
      'accessController',
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
      'setController',
      'version',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('callers can call view functions without explicit access', () => {
    it('#latestAnswer', async () => {
      await proxy.connect(personas.Carol).latestAnswer()
    })

    it('#latestTimestamp', async () => {
      await proxy.connect(personas.Carol).latestTimestamp()
    })

    it('#getAnswer', async () => {
      await proxy.connect(personas.Carol).getAnswer(phaseBase.add(1))
    })

    it('#getTimestamp', async () => {
      await proxy.connect(personas.Carol).getTimestamp(phaseBase.add(1))
    })

    it('#latestRound', async () => {
      await proxy.connect(personas.Carol).latestRound()
    })

    it('#getRoundData', async () => {
      await proxy.connect(personas.Carol).getRoundData(phaseBase.add(1))
    })

    it('#proposedGetRoundData', async () => {
      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, answer2)
      await proxy.proposeAggregator(aggregator2.address)
      const latestRound = await aggregator2.latestRound()
      await proxy.connect(personas.Carol).proposedGetRoundData(latestRound)
    })

    it('#proposedLatestRoundData', async () => {
      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(decimals, answer2)
      await proxy.proposeAggregator(aggregator2.address)
      await proxy.connect(personas.Carol).proposedLatestRoundData()
    })
  })

  describe('if the caller is granted access', () => {
    beforeEach(async () => {
      await controller.addAccess(defaultAccount.address)

      matchers.bigNum(
        ethers.utils.bigNumberify(answer),
        await aggregator.latestAnswer(),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(answer, await proxy.latestAnswer())
      const latestRound = await proxy.latestRound()
      matchers.bigNum(answer, await proxy.getAnswer(latestRound))
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

    it('getRoundData works', async () => {
      const proxyRoundId = await proxy.latestRound()
      const round = await proxy.getRoundData(proxyRoundId)
      matchers.bigNum(proxyRoundId, round.roundId)
      matchers.bigNum(answer, round.answer)
      matchers.bigNum(startedAt, round.startedAt)
      matchers.bigNum(timestamp, round.updatedAt)
    })

    describe('and an aggregator has been proposed', () => {
      beforeEach(async () => {
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(decimals, answer2)
        await proxy.proposeAggregator(aggregator2.address)
      })

      it('proposedGetRoundData works', async () => {
        const latestRound = await aggregator2.latestRound()
        const round = await proxy.proposedGetRoundData(latestRound)
        matchers.bigNum(latestRound, round.roundId)
        matchers.bigNum(answer2, round.answer)
      })

      it('proposedLatestRoundData works', async () => {
        const latestRound = await aggregator2.latestRound()
        const round = await proxy.proposedLatestRoundData()
        matchers.bigNum(latestRound, round.roundId)
        matchers.bigNum(answer2, round.answer)
      })
    })

    describe('without a proposed aggregator', () => {
      it('proposedGetRoundData reverts', async () => {
        await matchers.evmRevert(async () => {
          await proxy.proposedGetRoundData(1)
        }, 'No proposed aggregator present')
      })

      it('proposedLatestRoundData reverts', async () => {
        await matchers.evmRevert(async () => {
          await proxy.proposedLatestRoundData()
        }, 'No proposed aggregator present')
      })
    })

    describe('when read from a contract that is not permissioned', () => {
      beforeEach(async () => {
        testHelper = await testHelperFactory.connect(personas.Carol).deploy()
      })

      it('does not allow reading', async () => {
        await matchers.evmRevert(
          testHelper.readLatestRoundData(proxy.address),
          'No access',
        )
      })
    })
  })

  describe('#setController', () => {
    let newController: contract.Instance<SimpleReadAccessControllerFactory>

    beforeEach(async () => {
      newController = await accessControlFactory
        .connect(defaultAccount)
        .deploy()
    })

    describe('when called by a stranger', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await proxy
            .connect(personas.Carol)
            .setController(newController.address)
        }, 'Only callable by owner')
      })
    })

    describe('when called by the owner', () => {
      it('updates the controller contract', async () => {
        await proxy.connect(defaultAccount).setController(newController.address)
        assert.equal(await proxy.accessController(), newController.address)
      })
    })

    describe('when set to the zero address', () => {
      beforeEach(async () => {
        testHelper = await testHelperFactory.connect(personas.Carol).deploy()
      })

      it('allows anyone to read', async () => {
        await matchers.evmRevert(
          testHelper.readLatestRoundData(proxy.address),
          'No access',
        )

        await proxy.connect(defaultAccount).setController(emptyAddress)

        await testHelper.readLatestRoundData(proxy.address)
      })
    })
  })

  describe('#latestAnswer', () => {
    it('adds a small gas overhead on top of reading directly from the aggregator', async () => {
      testHelper = await testHelperFactory.connect(personas.Default).deploy()
      const link = await linkTokenFactory.connect(personas.Default).deploy()
      const aggregator = await (acAggregatorFactory as any)
        .connect(personas.Default)
        .deploy(
          link.address,
          0,
          0,
          emptyAddress,
          0,
          h.bigNum(2).pow(254),
          decimals,
          h.toBytes32String('TEST/LINK'),
          { gasLimit: 8_000_000 },
        )
      await proxy.proposeAggregator(aggregator.address)
      await proxy.confirmAggregator(aggregator.address)

      await aggregator.changeOracles(
        [],
        [personas.Neil.address],
        [personas.Neil.address],
        1,
        1,
        0,
      )
      await aggregator.connect(personas.Neil).submit(1, 100)

      await proxy.connect(personas.Default).setController(emptyAddress)
      await aggregator.disableAccessCheck()
      await aggregator.addAccess(proxy.address)

      const tx1 = await testHelper.readLatestAnswer(aggregator.address)
      const receipt1 = await tx1.wait()
      const tx2 = await testHelper.readLatestAnswer(proxy.address)
      const receipt2 = await tx2.wait()
      const diff = receipt2.gasUsed?.sub(receipt1.gasUsed || 0)
      assert.isAbove(3000, diff?.toNumber() || 3001)
    })
  })
})
