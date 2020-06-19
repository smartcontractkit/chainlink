import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
//import { ethers } from 'ethers'
import { FlagsFactory } from '../../ethers/v0.6/FlagsFactory'

const provider = setup.provider()
const flagsFactory = new FlagsFactory()
let personas: setup.Personas

beforeAll(async () => {
  personas = (await setup.users(provider)).personas
})

describe('Flags', () => {
  let flags: contract.Instance<FlagsFactory>
  let questionable: string
  const deployment = setup.snapshot(provider, async () => {
    flags = await flagsFactory.connect(personas.Nelly).deploy()
    await flags.connect(personas.Nelly).disableAccessCheck()
  })

  beforeEach(async () => {
    await deployment()
    questionable = personas.Norbert.address
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
        assert.equal(false, await flags.getFlag(questionable))

        await flags.connect(personas.Nelly).setFlagOn(questionable)

        assert.equal(true, await flags.getFlag(questionable))
      })

      it('emits an event log', async () => {
        const tx = await flags.connect(personas.Nelly).setFlagOn(questionable)
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagOn,
        )
        assert.equal(questionable, h.eventArgs(event).subject)
      })
      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).setFlagOn(questionable)
        })

        it('emits an event log', async () => {
          const tx = await flags.connect(personas.Nelly).setFlagOn(questionable)
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('updates the warning flag', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).setFlagOn(questionable),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#setFlagOff', () => {
    beforeEach(async () => {
      await flags.connect(personas.Nelly).setFlagOn(questionable)
    })

    describe('when called by the owner', () => {
      it('updates the warning flag', async () => {
        assert.equal(true, await flags.getFlag(questionable))

        await flags.connect(personas.Nelly).setFlagOff(questionable)

        assert.equal(false, await flags.getFlag(questionable))
      })

      it('emits an event log', async () => {
        const tx = await flags.connect(personas.Nelly).setFlagOff(questionable)
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagOff,
        )
        assert.equal(questionable, h.eventArgs(event).subject)
      })

      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).setFlagOff(questionable)
        })

        it('emits an event log', async () => {
          const tx = await flags
            .connect(personas.Nelly)
            .setFlagOff(questionable)
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('updates the warning flag', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).setFlagOff(questionable),
          'Only callable by owner',
        )
      })
    })
  })
})
