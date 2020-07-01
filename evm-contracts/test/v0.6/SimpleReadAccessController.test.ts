import { contract, helpers, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { SimpleReadAccessControllerFactory } from '../../ethers/v0.6/SimpleReadAccessControllerFactory'
import { AccessControlTestHelperFactory } from '../../ethers/v0.6/AccessControlTestHelperFactory'
import { ethers } from 'ethers'

const controllerFactory = new AccessControlTestHelperFactory()
const provider = setup.provider()
let personas: setup.Personas
let tx: ethers.ContractTransaction
beforeAll(async () => {
  await setup.users(provider).then(u => (personas = u.personas))
})

describe('SimpleReadAccessController', () => {
  let controller: contract.Instance<SimpleReadAccessControllerFactory>
  const value = 17
  const deployment = setup.snapshot(provider, async () => {
    controller = await controllerFactory.connect(personas.Carol).deploy(value)
  })
  beforeEach(deployment)

  it('has a limited public interface', () => {
    matchers.publicAbi(new SimpleReadAccessControllerFactory(), [
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

  describe('#constructor', () => {
    it('defaults checkEnabled to true', async () => {
      assert(await controller.checkEnabled())
    })
  })

  describe('#hasAccess', () => {
    it('allows unauthorized calls originating from the same account', async () => {
      assert.isTrue(
        await controller
          .connect(personas.Eddy)
          .hasAccess(personas.Eddy.address, '0x00'),
      )
    })

    it('blocks unauthorized calls originating from different accounts', async () => {
      assert.isFalse(
        await controller
          .connect(personas.Carol)
          .hasAccess(personas.Eddy.address, '0x00'),
      )
      assert.isFalse(
        await controller
          .connect(personas.Eddy)
          .hasAccess(personas.Carol.address, '0x00'),
      )
    })
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
      beforeEach(async () => {
        assert.isFalse(
          await controller.hasAccess(personas.Eddy.address, '0x00'),
        )
        tx = await controller.addAccess(personas.Eddy.address)
      })

      it('adds the address to the controller', async () => {
        assert.isTrue(await controller.hasAccess(personas.Eddy.address, '0x00'))
      })

      it('announces the change via a log', async () => {
        const receipt = await tx.wait()
        const event = helpers.findEventIn(
          receipt,
          controller.interface.events.AddedAccess,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Eddy.address)
      })

      describe('when called twice', () => {
        it('does not emit a log', async () => {
          const tx2 = await controller.addAccess(personas.Eddy.address)
          const receipt = await tx2.wait()
          assert.equal(receipt.events?.length, 0)
        })
      })
    })
  })

  describe('#removeAccess', () => {
    beforeEach(async () => {
      await controller.addAccess(personas.Eddy.address)
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
      beforeEach(async () => {
        tx = await controller.removeAccess(personas.Eddy.address)
      })

      it('removes the address from the controller', async () => {
        assert.isFalse(
          await controller.hasAccess(personas.Eddy.address, '0x00'),
        )
      })

      it('announces the change via a log', async () => {
        const receipt = await tx.wait()
        const event = helpers.findEventIn(
          receipt,
          controller.interface.events.RemovedAccess,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Eddy.address)
      })

      describe('when called twice', () => {
        it('does not emit a log', async () => {
          const tx2 = await controller.removeAccess(personas.Eddy.address)
          const receipt = await tx2.wait()
          assert.equal(receipt.events?.length, 0)
        })
      })
    })
  })

  describe('#disableAccessCheck', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          controller.connect(personas.Eddy).disableAccessCheck(),
          'Only callable by owner',
        )
        assert.isTrue(await controller.checkEnabled())
      })
    })

    describe('when called by the owner', () => {
      beforeEach(async () => {
        await controller.addAccess(personas.Eddy.address)
        tx = await controller.disableAccessCheck()
      })

      it('sets checkEnabled to false', async () => {
        assert.isFalse(await controller.checkEnabled())
      })

      it('allows users with access', async () => {
        assert.isTrue(await controller.hasAccess(personas.Eddy.address, '0x00'))
      })

      it('allows users without access', async () => {
        assert.isTrue(await controller.hasAccess(personas.Ned.address, '0x00'))
      })

      it('announces the change via a log', async () => {
        const receipt = await tx.wait()
        assert(
          helpers.findEventIn(
            receipt,
            controller.interface.events.CheckAccessDisabled,
          ),
        )
      })

      describe('when called twice', () => {
        it('does not emit a log', async () => {
          const tx2 = await controller.disableAccessCheck()
          const receipt = await tx2.wait()
          assert.equal(receipt.events?.length, 0)
        })
      })
    })
  })

  describe('#enableAccessCheck', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          controller.connect(personas.Eddy).enableAccessCheck(),
          'Only callable by owner',
        )
      })
    })

    describe('when called by the owner', () => {
      beforeEach(async () => {
        await controller.disableAccessCheck()
        await controller.addAccess(personas.Eddy.address)
        tx = await controller.enableAccessCheck()
      })

      it('allows users with access', async () => {
        assert.isTrue(await controller.hasAccess(personas.Eddy.address, '0x00'))
      })

      it('does not allow users without access', async () => {
        assert.isFalse(await controller.hasAccess(personas.Ned.address, '0x00'))
      })

      it('announces the change via a log', async () => {
        const receipt = await tx.wait()
        assert(
          helpers.findEventIn(
            receipt,
            controller.interface.events.CheckAccessEnabled,
          ),
        )
      })

      describe('when called twice', () => {
        it('does not emit a log', async () => {
          const tx2 = await controller.enableAccessCheck()
          const receipt = await tx2.wait()
          assert.equal(receipt.events?.length, 0)
        })
      })
    })
  })
})
