import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
//import { ethers } from 'ethers'
import { FlagsFactory } from '../../ethers/v0.6/FlagsFactory'
import { FlagsTestHelperFactory } from '../../ethers/v0.6/FlagsTestHelperFactory'

const provider = setup.provider()
const flagsFactory = new FlagsFactory()
const consumerFactory = new FlagsTestHelperFactory()
let personas: setup.Personas

beforeAll(async () => {
  personas = (await setup.users(provider)).personas
})

describe('Flags', () => {
  let flags: contract.Instance<FlagsFactory>
  let consumer: contract.Instance<FlagsTestHelperFactory>
  const deployment = setup.snapshot(provider, async () => {
    flags = await flagsFactory.connect(personas.Nelly).deploy()
    await flags.connect(personas.Nelly).disableAccessCheck()
    consumer = await consumerFactory
      .connect(personas.Nelly)
      .deploy(flags.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(flags, [
      'getFlag',
      'setFlagOff',
      'setFlagOn',
      // Ownable methods:
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

  describe('#setFlagOn', () => {
    describe('when called by the owner', () => {
      it('updates the warning flag', async () => {
        assert.equal(false, await flags.getFlag(consumer.address))

        await flags.connect(personas.Nelly).setFlagOn(consumer.address)

        assert.equal(true, await flags.getFlag(consumer.address))
      })

      it('emits an event log', async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .setFlagOn(consumer.address)
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagOn,
        )
        assert.equal(consumer.address, h.eventArgs(event).subject)
      })
      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).setFlagOn(consumer.address)
        })

        it('emits an event log', async () => {
          const tx = await flags
            .connect(personas.Nelly)
            .setFlagOn(consumer.address)
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('updates the warning flag', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).setFlagOn(consumer.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFlagOff', () => {
    beforeEach(async () => {
      await flags.connect(personas.Nelly).setFlagOn(consumer.address)
    })

    describe('when called by the owner', () => {
      it('updates the warning flag', async () => {
        assert.equal(true, await flags.getFlag(consumer.address))

        await flags.connect(personas.Nelly).setFlagOff(consumer.address)

        assert.equal(false, await flags.getFlag(consumer.address))
      })

      it('emits an event log', async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .setFlagOff(consumer.address)
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagOff,
        )
        assert.equal(consumer.address, h.eventArgs(event).subject)
      })

      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).setFlagOff(consumer.address)
        })

        it('emits an event log', async () => {
          const tx = await flags
            .connect(personas.Nelly)
            .setFlagOff(consumer.address)
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('updates the warning flag', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).setFlagOff(consumer.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#getFlag', () => {
    describe('if the access control is turned on', () => {
      beforeEach(async () => {
        await flags.connect(personas.Nelly).enableAccessCheck()
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          consumer.getFlag(consumer.address),
          'No access',
        )
      })

      describe('if access is granted to the address', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).addAccess(consumer.address)
        })

        it('does not revert', async () => {
          await consumer.getFlag(consumer.address)
        })
      })
    })

    describe('if the access control is turned off', () => {
      beforeEach(async () => {
        await flags.connect(personas.Nelly).disableAccessCheck()
      })

      it('does not revert', async () => {
        await consumer.getFlag(consumer.address)
      })

      describe('if access is granted to the address', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).addAccess(consumer.address)
        })

        it('does not revert', async () => {
          await consumer.getFlag(consumer.address)
        })
      })
    })
  })
})
