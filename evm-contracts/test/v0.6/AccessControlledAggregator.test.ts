import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { AccessControlledAggregator__factory } from '../../ethers/v0.6/factories/AccessControlledAggregator__factory'
import { FluxAggregatorTestHelper__factory } from '../../ethers/v0.6/factories/FluxAggregatorTestHelper__factory'

const aggregatorFactory = new AccessControlledAggregator__factory()
const linkTokenFactory = new contract.LinkToken__factory()
const testHelperFactory = new FluxAggregatorTestHelper__factory()
const provider = setup.provider()
let personas: setup.Personas

beforeAll(async () => {
  await setup.users(provider).then((u) => (personas = u.personas))
})

describe('AccessControlledAggregator', () => {
  const paymentAmount = h.toWei('3')
  const deposit = h.toWei('100')
  const answer = 100
  const minAns = 1
  const maxAns = 1
  const rrDelay = 0
  const timeout = 1800
  const decimals = 18
  const description = 'LINK/USD'
  const minSubmissionValue = h.bigNum('1')
  const maxSubmissionValue = h.bigNum('100000000000000000000')
  const emptyAddress = '0x0000000000000000000000000000000000000000'

  let link: contract.Instance<contract.LinkToken__factory>
  let aggregator: contract.Instance<AccessControlledAggregator__factory>
  let testHelper: contract.Instance<FluxAggregatorTestHelper__factory>
  let nextRound: number

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(personas.Default).deploy()
    aggregator = await (aggregatorFactory as any)
      .connect(personas.Carol)
      .deploy(
        link.address,
        paymentAmount,
        timeout,
        emptyAddress,
        minSubmissionValue,
        maxSubmissionValue,
        decimals,
        h.toBytes32String(description),
        // Remove when this PR gets merged:
        // https://github.com/ethereum-ts/TypeChain/pull/218
        { gasLimit: 8_000_000 },
      )
    await link.transfer(aggregator.address, deposit)
    await aggregator.updateAvailableFunds()
    matchers.bigNum(deposit, await link.balanceOf(aggregator.address))
    testHelper = await testHelperFactory.connect(personas.Carol).deploy()
  })

  beforeEach(async () => {
    await deployment()
    nextRound = 1
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(aggregatorFactory, [
      'acceptAdmin',
      'allocatedFunds',
      'availableFunds',
      'changeOracles',
      'decimals',
      'description',
      'getAdmin',
      'getAnswer',
      'getOracles',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'linkToken',
      'maxSubmissionCount',
      'maxSubmissionValue',
      'minSubmissionCount',
      'minSubmissionValue',
      'onTokenTransfer',
      'oracleCount',
      'oracleRoundState',
      'paymentAmount',
      'requestNewRound',
      'restartDelay',
      'setRequesterPermissions',
      'setValidator',
      'submit',
      'timeout',
      'transferAdmin',
      'updateAvailableFunds',
      'updateFutureRounds',
      'withdrawFunds',
      'withdrawPayment',
      'withdrawablePayment',
      'validator',
      'version',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // AccessControl methods:
      'addAccess',
      'disableAccessCheck',
      'enableAccessCheck',
      'removeAccess',
      'checkEnabled',
      'hasAccess',
    ])
  })

  describe('#constructor', () => {
    it('sets the paymentAmount', async () => {
      matchers.bigNum(h.bigNum(paymentAmount), await aggregator.paymentAmount())
    })

    it('sets the timeout', async () => {
      matchers.bigNum(h.bigNum(timeout), await aggregator.timeout())
    })

    it('sets the decimals', async () => {
      matchers.bigNum(h.bigNum(decimals), await aggregator.decimals())
    })

    it('sets the description', async () => {
      assert.equal(
        description,
        h.parseBytes32String(await aggregator.description()),
      )
    })
  })

  describe('#getAnswer', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when read by a contract', () => {
      describe('without explicit access', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            testHelper.readGetAnswer(aggregator.address, 0),
            'No access',
          )
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Carol).addAccess(testHelper.address)
          await testHelper.readGetAnswer(aggregator.address, 0)
        })
      })
    })

    describe('when read by a regular account', () => {
      describe('without explicit access', () => {
        it('succeeds', async () => {
          const round = await aggregator.latestRound()
          await aggregator.connect(personas.Eddy).getAnswer(round)
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator
            .connect(personas.Carol)
            .addAccess(personas.Eddy.address)
          const round = await aggregator.latestRound()
          await aggregator.connect(personas.Eddy).getAnswer(round)
        })
      })
    })
  })

  describe('#getTimestamp', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when read by a contract', () => {
      describe('without explicit access', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            testHelper.readGetTimestamp(aggregator.address, 0),
            'No access',
          )
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Carol).addAccess(testHelper.address)
          await testHelper.readGetTimestamp(aggregator.address, 0)
        })
      })
    })

    describe('when read by a regular account', () => {
      describe('without explicit access', () => {
        it('succeeds', async () => {
          const round = await aggregator.latestRound()
          const currentTimestamp = await aggregator
            .connect(personas.Eddy)
            .getTimestamp(round)
          assert.isAbove(currentTimestamp.toNumber(), 0)
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator
            .connect(personas.Carol)
            .addAccess(personas.Eddy.address)
          const round = await aggregator.latestRound()
          const currentTimestamp = await aggregator
            .connect(personas.Eddy)
            .getTimestamp(round)
          assert.isAbove(currentTimestamp.toNumber(), 0)
        })
      })
    })
  })

  describe('#latestAnswer', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when read by a contract', () => {
      describe('without explicit access', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            testHelper.readLatestAnswer(aggregator.address),
            'No access',
          )
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Carol).addAccess(testHelper.address)
          await testHelper.readLatestAnswer(aggregator.address)
        })
      })
    })

    describe('when read by a regular account', () => {
      describe('without explicit access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Eddy).latestAnswer()
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator
            .connect(personas.Carol)
            .addAccess(personas.Eddy.address)
          await aggregator.connect(personas.Eddy).latestAnswer()
        })
      })
    })
  })

  describe('#latestTimestamp', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when read by a contract', () => {
      describe('without explicit access', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            testHelper.readLatestTimestamp(aggregator.address),
            'No access',
          )
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Carol).addAccess(testHelper.address)
          await testHelper.readLatestTimestamp(aggregator.address)
        })
      })
    })

    describe('when read by a regular account', () => {
      describe('without explicit access', () => {
        it('succeeds', async () => {
          const currentTimestamp = await aggregator
            .connect(personas.Eddy)
            .latestTimestamp()
          assert.isAbove(currentTimestamp.toNumber(), 0)
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator
            .connect(personas.Carol)
            .addAccess(personas.Eddy.address)
          const currentTimestamp = await aggregator
            .connect(personas.Eddy)
            .latestTimestamp()
          assert.isAbove(currentTimestamp.toNumber(), 0)
        })
      })
    })
  })

  describe('#latestAnswer', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .changeOracles(
          [],
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when read by a contract', () => {
      describe('without explicit access', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            testHelper.readLatestAnswer(aggregator.address),
            'No access',
          )
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Carol).addAccess(testHelper.address)
          await testHelper.readLatestAnswer(aggregator.address)
        })
      })
    })

    describe('when read by a regular account', () => {
      describe('without explicit access', () => {
        it('succeeds', async () => {
          await aggregator.connect(personas.Eddy).latestAnswer()
        })
      })

      describe('with access', () => {
        it('succeeds', async () => {
          await aggregator
            .connect(personas.Carol)
            .addAccess(personas.Eddy.address)
          await aggregator.connect(personas.Eddy).latestAnswer()
        })
      })
    })
  })
})
