import {
  contract,
  matchers,
  helpers as h,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { DeviationFlaggingValidatorFactory } from '../../ethers/v0.6-ovm/DeviationFlaggingValidatorFactory'
import { FlagsFactory } from '../../ethers/v0.6-ovm/FlagsFactory'
import { SimpleWriteAccessControllerFactory } from '../../ethers/v0.6-ovm/SimpleWriteAccessControllerFactory'

let personas: setup.Personas
const provider = setup.provider()
const validatorFactory = new DeviationFlaggingValidatorFactory()
const flagsFactory = new FlagsFactory()
const acFactory = new SimpleWriteAccessControllerFactory()

beforeAll(async () => {
  personas = await setup.users(provider).then((x) => x.personas)
})

describe('DeviationFlaggingValidator', () => {
  let validator: contract.Instance<DeviationFlaggingValidatorFactory>
  let flags: contract.Instance<FlagsFactory>
  let ac: contract.Instance<SimpleWriteAccessControllerFactory>
  const flaggingThreshold = 10000 // 10%
  const previousRoundId = 2
  const previousValue = 1000000
  const currentRoundId = 3
  const currentValue = 1000000

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
      'THRESHOLD_MULTIPLIER',
      'flaggingThreshold',
      'flags',
      'isValid',
      'setFlagsAddress',
      'setFlaggingThreshold',
      'validate',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#constructor', () => {
    it('sets the arguments passed in', async () => {
      assert.equal(flags.address, await validator.flags())
      matchers.bigNum(flaggingThreshold, await validator.flaggingThreshold())
    })
  })

  describe('#validate', () => {
    describe('when the deviation is greater than the threshold', () => {
      const currentValue = 1100010

      it('does raises a flag for the calling address', async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(
            previousRoundId,
            previousValue,
            currentRoundId,
            currentValue,
          )
        const receipt = await tx.wait()
        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagRaised,
        )

        assert.equal(flags.address, event.address)
        assert.equal(
          personas.Nelly.address,
          h.evmWordToAddress(event.topics[1]),
        )
      })

      // OVM CHANGE: gas metering on OVM not implemented correctly
      it.skip('uses less than the gas allotted by the aggregator', async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(
            previousRoundId,
            previousValue,
            currentRoundId,
            currentValue,
          )
        const receipt = await tx.wait()
        assert(receipt)
        if (receipt && receipt.gasUsed) {
          assert.isAbove(60000, receipt.gasUsed.toNumber())
        }
      })
    })

    describe('when the deviation is less than or equal to the threshold', () => {
      const currentValue = 1100009

      it('does raises a flag for the calling address', async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(
            previousRoundId,
            previousValue,
            currentRoundId,
            currentValue,
          )
        const receipt = await tx.wait()
        matchers.eventDoesNotExist(receipt, flags.interface.events.FlagRaised)
      })

      // OVM CHANGE: gas metering on OVM not implemented correctly
      it.skip('uses less than the gas allotted by the aggregator', async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(
            previousRoundId,
            previousValue,
            currentRoundId,
            currentValue,
          )
        const receipt = await tx.wait()
        assert(receipt)
        if (receipt && receipt.gasUsed) {
          assert.isAbove(24000, receipt.gasUsed.toNumber())
        }
      })
    })

    describe('when called with a previous value of zero', () => {
      const previousValue = 0

      it('does not raise any flags', async () => {
        const tx = await validator
          .connect(personas.Nelly)
          .validate(
            previousRoundId,
            previousValue,
            currentRoundId,
            currentValue,
          )
        const receipt = await tx.wait()
        assert.equal(0, receipt.events?.length)
      })
    })
  })

  describe('#isValid', () => {
    const previousValue = 1000000

    describe('with a validation larger than the deviation', () => {
      const currentValue = 1100010
      it('is not valid', async () => {
        assert.isFalse(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })

    describe('with a validation smaller than the deviation', () => {
      const currentValue = 1100009
      it('is valid', async () => {
        assert.isTrue(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })

    describe('with positive previous and negative current', () => {
      const previousValue = 1000000
      const currentValue = -900000
      it('correctly detects the difference', async () => {
        assert.isFalse(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })

    describe('with negative previous and positive current', () => {
      const previousValue = -900000
      const currentValue = 1000000
      it('correctly detects the difference', async () => {
        assert.isFalse(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })

    describe('when the difference overflows', () => {
      const previousValue = h.bigNum(2).pow(255).sub(1)
      const currentValue = h.bigNum(-1)

      it('does not revert and returns false', async () => {
        assert.isFalse(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })

    describe('when the rounding overflows', () => {
      const previousValue = h.bigNum(2).pow(255).div(10000)
      const currentValue = h.bigNum(1)

      it('does not revert and returns false', async () => {
        assert.isFalse(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })

    describe('when the division overflows', () => {
      const previousValue = h.bigNum(2).pow(255).sub(1)
      const currentValue = h.bigNum(-1)

      it('does not revert and returns false', async () => {
        assert.isFalse(
          await validator.isValid(0, previousValue, 1, currentValue),
        )
      })
    })
  })

  describe('#setFlaggingThreshold', () => {
    const newThreshold = 777

    it('changes the flagging thresold', async () => {
      assert.equal(flaggingThreshold, await validator.flaggingThreshold())

      await validator.connect(personas.Carol).setFlaggingThreshold(newThreshold)

      assert.equal(newThreshold, await validator.flaggingThreshold())
    })

    it('emits a log event only when actually changed', async () => {
      const tx = await validator
        .connect(personas.Carol)
        .setFlaggingThreshold(newThreshold)
      const receipt = await tx.wait()
      const eventLog = matchers.eventExists(
        receipt,
        validator.interface.events.FlaggingThresholdUpdated,
      )

      assert.equal(flaggingThreshold, h.eventArgs(eventLog).previous)
      assert.equal(newThreshold, h.eventArgs(eventLog).current)

      const sameChangeTx = await validator
        .connect(personas.Carol)
        .setFlaggingThreshold(newThreshold)
      const sameChangeReceipt = await sameChangeTx.wait()
      assert.equal(0, sameChangeReceipt.events?.length)
      matchers.eventDoesNotExist(
        sameChangeReceipt,
        validator.interface.events.FlaggingThresholdUpdated,
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          validator.connect(personas.Neil).setFlaggingThreshold(newThreshold),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFlagsAddress', () => {
    const newFlagsAddress = '0x0123456789012345678901234567890123456789'

    it('changes the flags address', async () => {
      assert.equal(flags.address, await validator.flags())

      await validator.connect(personas.Carol).setFlagsAddress(newFlagsAddress)

      assert.equal(newFlagsAddress, await validator.flags())
    })

    it('emits a log event only when actually changed', async () => {
      const tx = await validator
        .connect(personas.Carol)
        .setFlagsAddress(newFlagsAddress)
      const receipt = await tx.wait()
      const eventLog = matchers.eventExists(
        receipt,
        validator.interface.events.FlagsAddressUpdated,
      )

      assert.equal(flags.address, h.eventArgs(eventLog).previous)
      assert.equal(newFlagsAddress, h.eventArgs(eventLog).current)

      const sameChangeTx = await validator
        .connect(personas.Carol)
        .setFlagsAddress(newFlagsAddress)
      const sameChangeReceipt = await sameChangeTx.wait()
      assert.equal(0, sameChangeReceipt.events?.length)
      matchers.eventDoesNotExist(
        sameChangeReceipt,
        validator.interface.events.FlagsAddressUpdated,
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          validator.connect(personas.Neil).setFlagsAddress(newFlagsAddress),
          'Only callable by owner',
        )
      })
    })
  })
})
