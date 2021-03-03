import {
  contract,
  matchers,
  helpers as h,
  setup,
} from '@chainlink/test-helpers'
import { ethers } from 'ethers'
import { assert } from 'chai'
import { Flags__factory } from '../../ethers/v0.6/factories/Flags__factory'
import { MockV3Aggregator__factory } from '../../ethers/v0.6/factories/MockV3Aggregator__factory'
import { MockCompoundOracle__factory } from '../../ethers/v0.7/factories/MockCompoundOracle__factory'
import { SimpleWriteAccessController__factory } from '../../ethers/v0.6/factories/SimpleWriteAccessController__factory'
import { CompoundPriceFlaggingValidator__factory } from '../../ethers/v0.7/factories/CompoundPriceFlaggingValidator__factory'
import { ContractReceipt } from 'ethers/contract'

let personas: setup.Personas
const provider = setup.provider()
const validatorFactory = new CompoundPriceFlaggingValidator__factory()
const acFactory = new SimpleWriteAccessController__factory()
const flagsFactory = new Flags__factory()
const aggregatorFactory = new MockV3Aggregator__factory()
const compoundOracleFactory = new MockCompoundOracle__factory()

beforeAll(async () => {
  personas = await setup.users(provider).then((x) => x.personas)
})

describe('CompoundPriceFlaggingVlidator', () => {
  let validator: contract.Instance<CompoundPriceFlaggingValidator__factory>
  let aggregator: contract.Instance<MockV3Aggregator__factory>
  let compoundOracle: contract.Instance<MockCompoundOracle__factory>
  let flags: contract.Instance<Flags__factory>
  let ac: contract.Instance<SimpleWriteAccessController__factory>

  const aggregatorDecimals = 18
  // 1000
  const initialAggregatorPrice = ethers.utils.bigNumberify(
    '1000000000000000000000',
  )

  const compoundSymbol = 'ETH'
  const compoundDecimals = 6
  // 1100 (10% deviation from aggregator price)
  const initialCompoundPrice = ethers.utils.bigNumberify('1100000000')

  // (50,000,000 / 1,000,000,000) = 0.05 = 5% deviation threshold
  const initialDeviationNumerator = 50_000_000

  const deployment = setup.snapshot(provider, async () => {
    ac = await acFactory.connect(personas.Carol).deploy()
    flags = await flagsFactory.connect(personas.Carol).deploy(ac.address)
    aggregator = await aggregatorFactory
      .connect(personas.Carol)
      .deploy(aggregatorDecimals, initialAggregatorPrice)
    compoundOracle = await compoundOracleFactory
      .connect(personas.Carol)
      .deploy()
    await compoundOracle.setPrice(
      compoundSymbol,
      initialCompoundPrice,
      compoundDecimals,
    )
    validator = await validatorFactory
      .connect(personas.Carol)
      .deploy(flags.address, compoundOracle.address)
    await validator
      .connect(personas.Carol)
      .setFeedDetails(
        aggregator.address,
        compoundSymbol,
        compoundDecimals,
        initialDeviationNumerator,
      )
    await ac.connect(personas.Carol).addAccess(validator.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(validatorFactory, [
      'update',
      'check',
      'setFeedDetails',
      'setFlagsAddress',
      'setCompoundOpenOracleAddress',
      'getFeedDetails',
      'flags',
      'compoundOpenOracle',
      // Upkeep methods:
      'checkUpkeep',
      'performUpkeep',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#constructor', () => {
    it('sets the owner', async () => {
      assert.equal(await validator.owner(), personas.Carol.address)
    })

    it('sets the arguments passed in', async () => {
      assert.equal(await validator.flags(), flags.address)
      assert.equal(await validator.compoundOpenOracle(), compoundOracle.address)
    })
  })

  describe('#setOpenOracleAddress', () => {
    let newCompoundOracle: contract.Instance<MockCompoundOracle__factory>
    let receipt: ContractReceipt

    beforeEach(async () => {
      newCompoundOracle = await compoundOracleFactory
        .connect(personas.Carol)
        .deploy()
      const tx = await validator
        .connect(personas.Carol)
        .setCompoundOpenOracleAddress(newCompoundOracle.address)
      receipt = await tx.wait()
    })

    it('changes the compound oracke address', async () => {
      assert.equal(
        await validator.compoundOpenOracle(),
        newCompoundOracle.address,
      )
    })

    it('emits a log event', async () => {
      const eventLog = matchers.eventExists(
        receipt,
        validator.interface.events.CompoundOpenOracleAddressUpdated,
      )

      assert.equal(h.eventArgs(eventLog).from, compoundOracle.address)
      assert.equal(h.eventArgs(eventLog).to, newCompoundOracle.address)
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          validator
            .connect(personas.Neil)
            .setCompoundOpenOracleAddress(newCompoundOracle.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFlagsAddress', () => {
    let newFlagsContract: contract.Instance<Flags__factory>
    let receipt: ContractReceipt

    beforeEach(async () => {
      newFlagsContract = await flagsFactory
        .connect(personas.Carol)
        .deploy(ac.address)
      const tx = await validator
        .connect(personas.Carol)
        .setFlagsAddress(newFlagsContract.address)
      receipt = await tx.wait()
    })

    it('changes the flags address', async () => {
      assert.equal(await validator.flags(), newFlagsContract.address)
    })

    it('emits a log event', async () => {
      const eventLog = matchers.eventExists(
        receipt,
        validator.interface.events.FlagsAddressUpdated,
      )

      assert.equal(h.eventArgs(eventLog).from, flags.address)
      assert.equal(h.eventArgs(eventLog).to, newFlagsContract.address)
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          validator
            .connect(personas.Neil)
            .setFlagsAddress(newFlagsContract.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFeedDetails', () => {
    let mockAggregator: contract.Instance<MockV3Aggregator__factory>
    let receipt: ContractReceipt
    const symbol = 'BTC'
    const decimals = 8
    const deviationNumerator = 50_000_000 // 5%

    beforeEach(async () => {
      await compoundOracle.connect(personas.Carol).setPrice('BTC', 1500000, 2)
      mockAggregator = await aggregatorFactory
        .connect(personas.Carol)
        .deploy(decimals, 4000000000000)
      const tx = await validator
        .connect(personas.Carol)
        .setFeedDetails(
          mockAggregator.address,
          symbol,
          decimals,
          deviationNumerator,
        )
      receipt = await tx.wait()
    })

    it('sets the correct state', async () => {
      const response = await validator
        .connect(personas.Carol)
        .getFeedDetails(mockAggregator.address)

      assert.equal(response[0], symbol)
      assert.equal(response[1], decimals)
      assert.equal(response[2].toString(), deviationNumerator.toString())
    })

    it('uses the existing symbol if one already exists', async () => {
      const newSymbol = 'LINK'

      await compoundOracle
        .connect(personas.Carol)
        .setPrice(newSymbol, 1500000, 2)

      const tx = await validator
        .connect(personas.Carol)
        .setFeedDetails(
          mockAggregator.address,
          newSymbol,
          decimals,
          deviationNumerator,
        )
      receipt = await tx.wait()

      // Check the event
      const eventLog = matchers.eventExists(
        receipt,
        validator.interface.events.FeedDetailsSet,
      )
      const eventArgs = h.eventArgs(eventLog)
      assert.equal(eventArgs.symbol, symbol)

      // Check the state
      const response = await validator
        .connect(personas.Carol)
        .getFeedDetails(mockAggregator.address)
      assert.equal(response[0], symbol)
    })

    it('emits an event', async () => {
      const eventLog = matchers.eventExists(
        receipt,
        validator.interface.events.FeedDetailsSet,
      )

      const eventArgs = h.eventArgs(eventLog)
      assert.equal(eventArgs.aggregator, mockAggregator.address)
      assert.equal(eventArgs.symbol, symbol)
      assert.equal(eventArgs.decimals, decimals)
      assert.equal(
        eventArgs.deviationThresholdNumerator.toString(),
        deviationNumerator.toString(),
      )
    })

    it('fails when given a 0 numerator', async () => {
      await matchers.evmRevert(
        validator
          .connect(personas.Carol)
          .setFeedDetails(mockAggregator.address, symbol, decimals, 0),
        'Invalid threshold numerator',
      )
    })

    it('fails when given a numerator above 1 billion', async () => {
      await matchers.evmRevert(
        validator
          .connect(personas.Carol)
          .setFeedDetails(
            mockAggregator.address,
            symbol,
            decimals,
            1_200_000_000,
          ),
        'Invalid threshold numerator',
      )
    })

    it('fails when the compound price is invalid', async () => {
      await matchers.evmRevert(
        validator
          .connect(personas.Carol)
          .setFeedDetails(
            mockAggregator.address,
            'TEST',
            decimals,
            deviationNumerator,
          ),
        'Invalid Compound price',
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          validator
            .connect(personas.Neil)
            .setFeedDetails(
              mockAggregator.address,
              symbol,
              decimals,
              deviationNumerator,
            ),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#check', () => {
    describe('with a single aggregator', () => {
      describe('with a deviated price exceding threshold', () => {
        it('returns the deviated aggregator', async () => {
          const aggregators = [aggregator.address]
          const response = await validator.check(aggregators)
          assert.equal(response.length, 1)
          assert.equal(response[0], aggregator.address)
        })
      })

      describe('with a price within the threshold', () => {
        const newCompoundPrice = ethers.utils.bigNumberify('1000000000')
        beforeEach(async () => {
          await compoundOracle.setPrice(
            'ETH',
            newCompoundPrice,
            compoundDecimals,
          )
        })

        it('returns an empty array', async () => {
          const aggregators = [aggregator.address]
          const response = await validator.check(aggregators)
          assert.equal(response.length, 0)
        })
      })
    })
  })

  describe('#update', () => {
    describe('with a single aggregator', () => {
      describe('with a deviated price exceding threshold', () => {
        it('raises a flag on the flags contract', async () => {
          const aggregators = [aggregator.address]
          const tx = await validator.connect(personas.Carol).update(aggregators)
          const logs = await h.getLogs(tx)
          assert.equal(logs.length, 1)
          assert.equal(
            h.evmWordToAddress(logs[0].topics[1]),
            aggregator.address,
          )
        })
      })

      describe('with a price within the threshold', () => {
        const newCompoundPrice = ethers.utils.bigNumberify('1000000000')
        beforeEach(async () => {
          await compoundOracle.setPrice(
            'ETH',
            newCompoundPrice,
            compoundDecimals,
          )
        })

        it('does nothing', async () => {
          const aggregators = [aggregator.address]
          const tx = await validator.connect(personas.Carol).update(aggregators)
          const logs = await h.getLogs(tx)
          assert.equal(logs.length, 0)
        })
      })
    })
  })

  describe('#checkUpkeep', () => {
    describe('with a single aggregator', () => {
      describe('with a deviated price exceding threshold', () => {
        it('returns the deviated aggregator', async () => {
          const aggregators = [aggregator.address]
          const encodedAggregators = ethers.utils.defaultAbiCoder.encode(
            ['address[]'],
            [aggregators],
          )
          const response = await validator
            .connect(personas.Carol)
            .checkUpkeep(encodedAggregators)

          const decodedResponse = ethers.utils.defaultAbiCoder.decode(
            ['address[]'],
            response?.[1],
          )
          assert.equal(decodedResponse?.[0]?.[0], aggregators[0])
        })
      })

      describe('with a price within the threshold', () => {
        const newCompoundPrice = ethers.utils.bigNumberify('1000000000')
        beforeEach(async () => {
          await compoundOracle.setPrice(
            'ETH',
            newCompoundPrice,
            compoundDecimals,
          )
        })

        it('returns an empty array', async () => {
          const aggregators = [aggregator.address]
          const encodedAggregators = ethers.utils.defaultAbiCoder.encode(
            ['address[]'],
            [aggregators],
          )
          const response = await validator
            .connect(personas.Carol)
            .checkUpkeep(encodedAggregators)
          const decodedResponse = ethers.utils.defaultAbiCoder.decode(
            ['address[]'],
            response?.[1],
          )
          assert.equal(decodedResponse?.[0]?.length, 0)
        })
      })
    })
  })

  describe('#performUpkeep', () => {
    describe('with a single aggregator', () => {
      describe('with a deviated price exceding threshold', () => {
        it('raises a flag on the flags contract', async () => {
          const aggregators = [aggregator.address]
          const encodedAggregators = ethers.utils.defaultAbiCoder.encode(
            ['address[]'],
            [aggregators],
          )
          const tx = await validator
            .connect(personas.Carol)
            .performUpkeep(encodedAggregators)
          const logs = await h.getLogs(tx)
          assert.equal(logs.length, 1)
          assert.equal(
            h.evmWordToAddress(logs[0].topics[1]),
            aggregator.address,
          )
        })
      })

      describe('with a price within the threshold', () => {
        const newCompoundPrice = ethers.utils.bigNumberify('1000000000')
        beforeEach(async () => {
          await compoundOracle.setPrice(
            'ETH',
            newCompoundPrice,
            compoundDecimals,
          )
        })

        it('does nothing', async () => {
          const aggregators = [aggregator.address]
          const encodedAggregators = ethers.utils.defaultAbiCoder.encode(
            ['address[]'],
            [aggregators],
          )
          const tx = await validator
            .connect(personas.Carol)
            .performUpkeep(encodedAggregators)
          const logs = await h.getLogs(tx)
          assert.equal(logs.length, 0)
        })
      })
    })
  })
})
