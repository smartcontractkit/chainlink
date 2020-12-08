import {
  contract,
  matchers,
  // helpers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { StalenessFlaggingValidatorFactory } from '../../ethers/v0.7/StalenessFlaggingValidatorFactory'
import { FlagsFactory } from '../../ethers/v0.6/FlagsFactory'
import { SimpleWriteAccessControllerFactory } from '../../ethers/v0.6/SimpleWriteAccessControllerFactory'


let personas: setup.Personas
const provider = setup.provider()
const validatorFactory = new StalenessFlaggingValidatorFactory()
const flagsFactory = new FlagsFactory()
const acFactory = new SimpleWriteAccessControllerFactory()


beforeAll(async () => {
  personas = await setup.users(provider).then((x) => x.personas)
})

describe('StalenessFlaggingValidator', () => {
  let validator: contract.Instance<StalenessFlaggingValidatorFactory>
  let flags: contract.Instance<FlagsFactory>
  let ac: contract.Instance<SimpleWriteAccessControllerFactory>

  const flaggingThreshold = 10000

  const deployment = setup.snapshot(provider, async () => {
    ac = await acFactory.connect(personas.Carol).deploy()
    flags = await flagsFactory.connect(personas.Carol).deploy(ac.address)
    validator = await validatorFactory
      .connect(personas.Carol)
      .deploy(flags.address, flaggingThreshold)
    await ac.connect(personas.Carol).addAccess(validator.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(validatorFactory, [
      'update',
      'check',
      'setThreshold',
      'setFlagsAddress',
      'threshold',
      'flags',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#constructor', () => {
    it('sets the arguments passed in', async () => {
      assert.equal(flags.address, await validator.flags())
      matchers.bigNum(flaggingThreshold, await validator.threshold())
    })
  })

  // TODO
})