import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { FlagsFactory } from '../../ethers/v0.7/FlagsFactory'
import { FlagsTestHelperFactory } from '../../ethers/v0.7/FlagsTestHelperFactory'
import { SimpleWriteAccessControllerFactory } from '../../ethers/v0.7/SimpleWriteAccessControllerFactory'

const provider = setup.provider()
const flagsFactory = new FlagsFactory()
const consumerFactory = new FlagsTestHelperFactory()
const accessControlFactory = new SimpleWriteAccessControllerFactory()
let personas: setup.Personas

beforeAll(async () => {
  personas = (await setup.users(provider)).personas
})

describe('Flags', () => {
  let controller: contract.Instance<SimpleWriteAccessControllerFactory>
  let flags: contract.Instance<FlagsFactory>
  let consumer: contract.Instance<FlagsTestHelperFactory>

  const deployment = setup.snapshot(provider, async () => {
    controller = await accessControlFactory.connect(personas.Nelly).deploy()
    flags = await flagsFactory
      .connect(personas.Nelly)
      .deploy(controller.address)
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
      'getFlags',
      'lowerFlags',
      'raiseFlag',
      'raiseFlags',
      'raisingAccessController',
      'setRaisingAccessController',
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

  describe('#raiseFlag', () => {
    describe('when called by the owner', () => {
      it('updates the warning flag', async () => {
        assert.equal(false, await flags.getFlag(consumer.address))

        await flags.connect(personas.Nelly).raiseFlag(consumer.address)

        assert.equal(true, await flags.getFlag(consumer.address))
      })

      it('emits an event log', async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .raiseFlag(consumer.address)
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagRaised,
        )
        assert.equal(consumer.address, h.eventArgs(event).subject)
      })

      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).raiseFlag(consumer.address)
        })

        it('emits an event log', async () => {
          const tx = await flags
            .connect(personas.Nelly)
            .raiseFlag(consumer.address)
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by an enabled setter', () => {
      beforeEach(async () => {
        await controller
          .connect(personas.Nelly)
          .addAccess(personas.Neil.address)
      })

      it('sets the flags', async () => {
        await flags.connect(personas.Neil).raiseFlag(consumer.address),
          assert.equal(true, await flags.getFlag(consumer.address))
      })
    })

    describe('when called by a non-enabled setter', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).raiseFlag(consumer.address),
          'Not allowed to raise flags',
        )
      })
    })

    describe('when called when there is no raisingAccessController', () => {
      beforeEach(async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .setRaisingAccessController(
            '0x0000000000000000000000000000000000000000',
          )
        const receipt = await tx.wait()
        const event = matchers.eventExists(
          receipt,
          flags.interface.events.RaisingAccessControllerUpdated,
        )
        assert.equal(
          '0x0000000000000000000000000000000000000000',
          h.eventArgs(event).current,
        )
        assert.equal(
          '0x0000000000000000000000000000000000000000',
          await flags.raisingAccessController(),
        )
      })

      it('succeeds for the owner', async () => {
        await flags.connect(personas.Nelly).raiseFlag(consumer.address)
        assert.equal(true, await flags.getFlag(consumer.address))
      })

      it('reverts for non-owner', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).raiseFlag(consumer.address),
        )
      })
    })
  })

  describe('#raiseFlags', () => {
    describe('when called by the owner', () => {
      it('updates the warning flag', async () => {
        assert.equal(false, await flags.getFlag(consumer.address))

        await flags.connect(personas.Nelly).raiseFlags([consumer.address])

        assert.equal(true, await flags.getFlag(consumer.address))
      })

      it('emits an event log', async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .raiseFlags([consumer.address])
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagRaised,
        )
        assert.equal(consumer.address, h.eventArgs(event).subject)
      })

      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).raiseFlags([consumer.address])
        })

        it('emits an event log', async () => {
          const tx = await flags
            .connect(personas.Nelly)
            .raiseFlags([consumer.address])
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by an enabled setter', () => {
      beforeEach(async () => {
        await controller
          .connect(personas.Nelly)
          .addAccess(personas.Neil.address)
      })

      it('sets the flags', async () => {
        await flags.connect(personas.Neil).raiseFlags([consumer.address]),
          assert.equal(true, await flags.getFlag(consumer.address))
      })
    })

    describe('when called by a non-enabled setter', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).raiseFlags([consumer.address]),
          'Not allowed to raise flags',
        )
      })
    })

    describe('when called when there is no raisingAccessController', () => {
      beforeEach(async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .setRaisingAccessController(
            '0x0000000000000000000000000000000000000000',
          )
        const receipt = await tx.wait()
        const event = matchers.eventExists(
          receipt,
          flags.interface.events.RaisingAccessControllerUpdated,
        )
        assert.equal(
          '0x0000000000000000000000000000000000000000',
          h.eventArgs(event).current,
        )
        assert.equal(
          '0x0000000000000000000000000000000000000000',
          await flags.raisingAccessController(),
        )
      })

      it('succeeds for the owner', async () => {
        await flags.connect(personas.Nelly).raiseFlags([consumer.address])
        assert.equal(true, await flags.getFlag(consumer.address))
      })

      it('reverts for non-owners', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).raiseFlags([consumer.address]),
        )
      })
    })
  })

  describe('#lowerFlags', () => {
    beforeEach(async () => {
      await flags.connect(personas.Nelly).raiseFlags([consumer.address])
    })

    describe('when called by the owner', () => {
      it('updates the warning flag', async () => {
        assert.equal(true, await flags.getFlag(consumer.address))

        await flags.connect(personas.Nelly).lowerFlags([consumer.address])

        assert.equal(false, await flags.getFlag(consumer.address))
      })

      it('emits an event log', async () => {
        const tx = await flags
          .connect(personas.Nelly)
          .lowerFlags([consumer.address])
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          flags.interface.events.FlagLowered,
        )
        assert.equal(consumer.address, h.eventArgs(event).subject)
      })

      describe('if a flag has already been raised', () => {
        beforeEach(async () => {
          await flags.connect(personas.Nelly).lowerFlags([consumer.address])
        })

        it('emits an event log', async () => {
          const tx = await flags
            .connect(personas.Nelly)
            .lowerFlags([consumer.address])
          const receipt = await tx.wait()
          assert.equal(0, receipt.events?.length)
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          flags.connect(personas.Neil).lowerFlags([consumer.address]),
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

  describe('#getFlags', () => {
    beforeEach(async () => {
      await flags.connect(personas.Nelly).disableAccessCheck()
      await flags
        .connect(personas.Nelly)
        .raiseFlags([personas.Neil.address, personas.Norbert.address])
    })

    it('respects the access controls of #getFlag', async () => {
      await flags.connect(personas.Nelly).enableAccessCheck()

      await matchers.evmRevert(consumer.getFlag(consumer.address), 'No access')

      await flags.connect(personas.Nelly).addAccess(consumer.address)

      await consumer.getFlag(consumer.address)
    })

    it('returns the flags in the order they are requested', async () => {
      const response = await consumer.getFlags([
        personas.Nelly.address,
        personas.Neil.address,
        personas.Ned.address,
        personas.Norbert.address,
      ])

      assert.deepEqual([false, true, false, true], response)
    })
  })

  describe('#setRaisingAccessController', () => {
    let controller2: any

    beforeEach(async () => {
      controller2 = await accessControlFactory.connect(personas.Nelly).deploy()
      await controller2.connect(personas.Nelly).enableAccessCheck()
    })

    it('updates access control rules', async () => {
      await controller.connect(personas.Nelly).addAccess(personas.Neil.address)
      await flags.connect(personas.Neil).raiseFlags([consumer.address]) // doesn't raise

      await flags
        .connect(personas.Nelly)
        .setRaisingAccessController(controller2.address)

      await matchers.evmRevert(
        flags.connect(personas.Neil).raiseFlags([consumer.address]),
        'Not allowed to raise flags',
      ) // raises with new controller
    })

    it('emits a log announcing the change', async () => {
      const tx = await flags
        .connect(personas.Nelly)
        .setRaisingAccessController(controller2.address)
      const receipt = await tx.wait()

      const event = matchers.eventExists(
        receipt,
        flags.interface.events.RaisingAccessControllerUpdated,
      )
      assert.equal(controller.address, h.eventArgs(event).previous)
      assert.equal(controller2.address, h.eventArgs(event).current)
    })

    it('does not emit a log when there is no change', async () => {
      await flags
        .connect(personas.Nelly)
        .setRaisingAccessController(controller2.address)

      const tx = await flags
        .connect(personas.Nelly)
        .setRaisingAccessController(controller2.address)
      const receipt = await tx.wait()

      matchers.eventDoesNotExist(
        receipt,
        flags.interface.events.RaisingAccessControllerUpdated,
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          flags
            .connect(personas.Neil)
            .setRaisingAccessController(controller2.address),
          'Only callable by owner',
        )
      })
    })
  })
})
