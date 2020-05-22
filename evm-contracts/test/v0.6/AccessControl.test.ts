import { contract, helpers, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { AccessControlFactory } from '../../ethers/v0.6/AccessControlFactory'
import { AccessControlTestHelperFactory } from '../../ethers/v0.6/AccessControlTestHelperFactory'

const controllerFactory = new AccessControlTestHelperFactory()
const provider = setup.provider()
let personas: setup.Personas
beforeAll(async () => {
  await setup.users(provider).then(u => (personas = u.personas))
})

describe('AccessControl', () => {
  let controller: contract.Instance<AccessControlFactory>
  const value = 17
  const deployment = setup.snapshot(provider, async () => {
    controller = await controllerFactory.connect(personas.Carol).deploy(value)
  })
  beforeEach(deployment)

  it('has a limited public interface', () => {
    matchers.publicAbi(new AccessControlFactory(), [
      'hasAccess',
      'addAccess',
      'disableAccessCheck',
      'enableAccessCheck',
      'removeAccess',
      'checkEnabled',
      // Owned
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#addAccess', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          controller.connect(personas.Eddy).addAccess(personas.Eddy.address),
          'Only callable by owner',
        )
      })
    })

    describe('when called by the owner', () => {
      it('adds the address to the controller', async () => {
        const tx = await controller
          .connect(personas.Carol)
          .addAccess(personas.Eddy.address)
        const receipt = await tx.wait()

        assert.isTrue(await controller.hasAccess(personas.Eddy.address, '0x00'))

        const event = helpers.findEventIn(
          receipt,
          controller.interface.events.AddedAccess,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Eddy.address)
      })

      it('allows controller users', async () => {
        await controller
          .connect(personas.Carol)
          .addAccess(personas.Eddy.address)

        matchers.bigNum(
          value,
          await controller.connect(personas.Eddy).getValue(),
        )
      })
    })
  })

  describe('#removeAccess', () => {
    beforeEach(async () => {
      await controller.connect(personas.Carol).addAccess(personas.Eddy.address)
      assert.isTrue(await controller.hasAccess(personas.Eddy.address, '0x00'))
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          controller.connect(personas.Eddy).removeAccess(personas.Eddy.address),
          'Only callable by owner',
        )
      })
    })

    describe('when called by the owner', () => {
      it('removes the address from the controller', async () => {
        const tx = await controller
          .connect(personas.Carol)
          .removeAccess(personas.Eddy.address)
        const receipt = await tx.wait()

        assert.isFalse(
          await controller.hasAccess(personas.Eddy.address, '0x00'),
        )

        const event = helpers.findEventIn(
          receipt,
          controller.interface.events.RemovedAccess,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Eddy.address)
      })

      it('does not allow users without access', async () => {
        await controller
          .connect(personas.Carol)
          .removeAccess(personas.Eddy.address)

        await matchers.evmRevert(controller.connect(personas.Eddy).getValue())
      })
    })
  })

  describe('#checkEnabled', () => {
    it('defaults to true', async () => {
      assert(await controller.checkEnabled())
    })
  })

  describe('#enableAccessCheck', () => {
    beforeEach(async () => {
      await controller.connect(personas.Carol).addAccess(personas.Eddy.address)
    })

    it('allows users with access', async () => {
      await controller.connect(personas.Carol).enableAccessCheck()

      matchers.bigNum(value, await controller.connect(personas.Eddy).getValue())
    })

    it('does not allow users without access', async () => {
      await controller.connect(personas.Carol).enableAccessCheck()

      await matchers.evmRevert(controller.connect(personas.Ned).getValue())
    })

    it('announces the change via a log', async () => {
      const tx = await controller.connect(personas.Carol).enableAccessCheck()
      const receipt = await tx.wait()

      assert(
        helpers.findEventIn(
          receipt,
          controller.interface.events.CheckAccessEnabled,
        ),
      )
    })

    it('reverts when called by a non-owner', async () => {
      await matchers.evmRevert(
        controller.connect(personas.Eddy).enableAccessCheck(),
        'Only callable by owner',
      )
    })
  })

  describe('#disableAccessCheck', () => {
    beforeEach(async () => {
      await controller.connect(personas.Carol).addAccess(personas.Eddy.address)
    })

    it('allows users with access', async () => {
      await controller.connect(personas.Carol).disableAccessCheck()

      matchers.bigNum(value, await controller.connect(personas.Eddy).getValue())
    })

    it('allows users without access', async () => {
      await controller.connect(personas.Carol).disableAccessCheck()

      matchers.bigNum(value, await controller.connect(personas.Ned).getValue())
    })

    it('announces the change via a log', async () => {
      const tx = await controller.connect(personas.Carol).disableAccessCheck()
      const receipt = await tx.wait()

      assert(
        helpers.findEventIn(
          receipt,
          controller.interface.events.CheckAccessDisabled,
        ),
      )
    })

    it('reverts when called by a non-owner', async () => {
      await matchers.evmRevert(
        controller.connect(personas.Eddy).disableAccessCheck(),
        'Only callable by owner',
      )
    })
  })
})
