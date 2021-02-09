import {
  contract,
  // matchers,
  // helpers as h,
  setup,
} from '@chainlink/test-helpers'
import { ethers } from 'ethers'
import { assert } from 'chai'
import { Flags__factory } from '../../ethers/v0.6/factories/Flags__factory'
import { MockV3Aggregator__factory } from '../../ethers/v0.6/factories/MockV3Aggregator__factory'
import { MockCompoundOracle__factory } from '../../ethers/v0.7/factories/MockCompoundOracle__factory'
import { SimpleWriteAccessController__factory } from '../../ethers/v0.6/factories/SimpleWriteAccessController__factory'
import { CompoundPriceFlaggingValidator__factory } from '../../ethers/v0.7/factories/CompoundPriceFlaggingValidator__factory'

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
  const initialAggregatorPrice = ethers.utils.bigNumberify("1000000000000000000000")

  const compoundSymbol = "ETH"
  const compoundDecimals = 6
  // 1200 (20% deviation from aggregator price)
  const initialCompoundPrice = ethers.utils.bigNumberify("1200000000")

  const initialDeviationDenominator = 10

  const deployment = setup.snapshot(provider, async () => {
    ac = await acFactory.connect(personas.Carol).deploy()
    flags = await flagsFactory.connect(personas.Carol).deploy(ac.address)
    aggregator = await aggregatorFactory
      .connect(personas.Carol)
      .deploy(aggregatorDecimals, initialAggregatorPrice)
    compoundOracle = await compoundOracleFactory
      .connect(personas.Carol)
      .deploy()
    await compoundOracle
      .setPrice(compoundSymbol, initialCompoundPrice, compoundDecimals)
    validator = await validatorFactory
      .connect(personas.Carol)
      .deploy(flags.address, compoundOracle.address)
    await validator.connect(personas.Carol)
      .setThreshold(
        aggregator.address,
        compoundSymbol,
        compoundDecimals,
        initialDeviationDenominator
      )
    await ac.connect(personas.Carol).addAccess(validator.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#check', () => {
    it('returns deviated aggregator', async () => {
      const aggregators = [aggregator.address]
      const response = await validator.check(aggregators)
      assert.equal(response.length, 1)
      assert.equal(response[0], aggregator.address)
    })

  })
})
