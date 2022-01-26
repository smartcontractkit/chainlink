import { ethers } from 'hardhat'
import { evmWordToAddress, getLogs, publicAbi } from '../test-helpers/helpers'
import { assert, expect } from 'chai'
import {
  BigNumber,
  Contract,
  ContractFactory,
  ContractTransaction,
} from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'
import { evmRevert } from '../test-helpers/matchers'

let personas: Personas
let validatorFactory: ContractFactory
let acFactory: ContractFactory
let flagsFactory: ContractFactory
let aggregatorFactory: ContractFactory
let compoundOracleFactory: ContractFactory

before(async () => {
  personas = (await getUsers()).personas

  validatorFactory = await ethers.getContractFactory(
    'src/v0.7/dev/CompoundPriceFlaggingValidator.sol:CompoundPriceFlaggingValidator',
    personas.Carol,
  )
  acFactory = await ethers.getContractFactory(
    'src/v0.6/SimpleWriteAccessController.sol:SimpleWriteAccessController',
    personas.Carol,
  )
  flagsFactory = await ethers.getContractFactory(
    'src/v0.6/Flags.sol:Flags',
    personas.Carol,
  )
  aggregatorFactory = await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
    personas.Carol,
  )
  compoundOracleFactory = await ethers.getContractFactory(
    'src/v0.7/tests/MockCompoundOracle.sol:MockCompoundOracle',
    personas.Carol,
  )
})

describe('CompoundPriceFlaggingVlidator', () => {
  let validator: Contract
  let aggregator: Contract
  let compoundOracle: Contract
  let flags: Contract
  let ac: Contract

  const aggregatorDecimals = 18
  // 1000
  const initialAggregatorPrice = BigNumber.from('1000000000000000000000')

  const compoundSymbol = 'ETH'
  const compoundDecimals = 6
  // 1100 (10% deviation from aggregator price)
  const initialCompoundPrice = BigNumber.from('1100000000')

  // (50,000,000 / 1,000,000,000) = 0.05 = 5% deviation threshold
  const initialDeviationNumerator = 50_000_000

  beforeEach(async () => {
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

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(validator, [
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
      assert.equal(await validator.owner(), await personas.Carol.getAddress())
    })

    it('sets the arguments passed in', async () => {
      assert.equal(await validator.flags(), flags.address)
      assert.equal(await validator.compoundOpenOracle(), compoundOracle.address)
    })
  })

  describe('#setOpenOracleAddress', () => {
    let newCompoundOracle: Contract
    let tx: ContractTransaction

    beforeEach(async () => {
      newCompoundOracle = await compoundOracleFactory
        .connect(personas.Carol)
        .deploy()
      tx = await validator
        .connect(personas.Carol)
        .setCompoundOpenOracleAddress(newCompoundOracle.address)
    })

    it('changes the compound oracke address', async () => {
      assert.equal(
        await validator.compoundOpenOracle(),
        newCompoundOracle.address,
      )
    })

    it('emits a log event', async () => {
      await expect(tx)
        .to.emit(validator, 'CompoundOpenOracleAddressUpdated')
        .withArgs(compoundOracle.address, newCompoundOracle.address)
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await evmRevert(
          validator
            .connect(personas.Neil)
            .setCompoundOpenOracleAddress(newCompoundOracle.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFlagsAddress', () => {
    let newFlagsContract: Contract
    let tx: ContractTransaction

    beforeEach(async () => {
      newFlagsContract = await flagsFactory
        .connect(personas.Carol)
        .deploy(ac.address)
      tx = await validator
        .connect(personas.Carol)
        .setFlagsAddress(newFlagsContract.address)
    })

    it('changes the flags address', async () => {
      assert.equal(await validator.flags(), newFlagsContract.address)
    })

    it('emits a log event', async () => {
      await expect(tx)
        .to.emit(validator, 'FlagsAddressUpdated')
        .withArgs(flags.address, newFlagsContract.address)
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await evmRevert(
          validator
            .connect(personas.Neil)
            .setFlagsAddress(newFlagsContract.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFeedDetails', () => {
    let mockAggregator: Contract
    let tx: ContractTransaction
    const symbol = 'BTC'
    const decimals = 8
    const deviationNumerator = 50_000_000 // 5%

    beforeEach(async () => {
      await compoundOracle.connect(personas.Carol).setPrice('BTC', 1500000, 2)
      mockAggregator = await aggregatorFactory
        .connect(personas.Carol)
        .deploy(decimals, 4000000000000)
      tx = await validator
        .connect(personas.Carol)
        .setFeedDetails(
          mockAggregator.address,
          symbol,
          decimals,
          deviationNumerator,
        )
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

      tx = await validator
        .connect(personas.Carol)
        .setFeedDetails(
          mockAggregator.address,
          newSymbol,
          decimals,
          deviationNumerator,
        )

      // Check the event
      await expect(tx)
        .to.emit(validator, 'FeedDetailsSet')
        .withArgs(mockAggregator.address, symbol, decimals, deviationNumerator)

      // Check the state
      const response = await validator
        .connect(personas.Carol)
        .getFeedDetails(mockAggregator.address)
      assert.equal(response[0], symbol)
    })

    it('emits an event', async () => {
      await expect(tx)
        .to.emit(validator, 'FeedDetailsSet')
        .withArgs(mockAggregator.address, symbol, decimals, deviationNumerator)
    })

    it('fails when given a 0 numerator', async () => {
      await evmRevert(
        validator
          .connect(personas.Carol)
          .setFeedDetails(mockAggregator.address, symbol, decimals, 0),
        'Invalid threshold numerator',
      )
    })

    it('fails when given a numerator above 1 billion', async () => {
      await evmRevert(
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
      await evmRevert(
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
        await evmRevert(
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
        const newCompoundPrice = BigNumber.from('1000000000')
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
          const logs = await getLogs(tx)
          assert.equal(logs.length, 1)
          assert.equal(evmWordToAddress(logs[0].topics[1]), aggregator.address)
        })
      })

      describe('with a price within the threshold', () => {
        const newCompoundPrice = BigNumber.from('1000000000')
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
          const logs = await getLogs(tx)
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
        const newCompoundPrice = BigNumber.from('1000000000')
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
          const logs = await getLogs(tx)
          assert.equal(logs.length, 1)
          assert.equal(evmWordToAddress(logs[0].topics[1]), aggregator.address)
        })
      })

      describe('with a price within the threshold', () => {
        const newCompoundPrice = BigNumber.from('1000000000')
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
          const logs = await getLogs(tx)
          assert.equal(logs.length, 0)
        })
      })
    })
  })
})
