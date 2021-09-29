import { ethers } from 'hardhat'
import { publicAbi } from '../test-helpers/helpers'
import { assert, expect } from 'chai'
import { Contract, ContractFactory, Transaction } from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'

let personas: Personas

let controllerFactory: ContractFactory
let controller: Contract

before(async () => {
  personas = (await getUsers()).personas
  controllerFactory = await ethers.getContractFactory(
    'src/v0.6/SimpleReadAccessController.sol:SimpleReadAccessController',
    personas.Carol,
  )
})

describe('SimpleReadAccessController', () => {
  beforeEach(async () => {
    controller = await controllerFactory.connect(personas.Carol).deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(controller, [
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
          .hasAccess(await personas.Eddy.getAddress(), '0x00'),
      )
    })

    it('blocks unauthorized calls originating from different accounts', async () => {
      assert.isFalse(
        await controller
          .connect(personas.Carol)
          .hasAccess(await personas.Eddy.getAddress(), '0x00'),
      )
      assert.isFalse(
        await controller
          .connect(personas.Eddy)
          .hasAccess(await personas.Carol.getAddress(), '0x00'),
      )
    })
  })

  describe('#addAccess', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Eddy)
            .addAccess(await personas.Eddy.getAddress()),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {
      let tx: Transaction
      beforeEach(async () => {
        assert.isFalse(
          await controller.hasAccess(await personas.Eddy.getAddress(), '0x00'),
        )
        tx = await controller.addAccess(await personas.Eddy.getAddress())
      })

      it('adds the address to the controller', async () => {
        assert.isTrue(
          await controller.hasAccess(await personas.Eddy.getAddress(), '0x00'),
        )
      })

      it('announces the change via a log', async () => {
        await expect(tx)
          .to.emit(controller, 'AddedAccess')
          .withArgs(await personas.Eddy.getAddress())
      })

      describe('when called twice', () => {
        it('does not emit a log', async () => {
          const tx2 = await controller.addAccess(
            await personas.Eddy.getAddress(),
          )
          const receipt = await tx2.wait()
          assert.equal(receipt.events?.length, 0)
        })
      })
    })
  })

  describe('#removeAccess', () => {
    beforeEach(async () => {
      await controller.addAccess(await personas.Eddy.getAddress())
      assert.isTrue(
        await controller.hasAccess(await personas.Eddy.getAddress(), '0x00'),
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Eddy)
            .removeAccess(await personas.Eddy.getAddress()),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {
      let tx: Transaction
      beforeEach(async () => {
        tx = await controller.removeAccess(await personas.Eddy.getAddress())
      })

      it('removes the address from the controller', async () => {
        assert.isFalse(
          await controller.hasAccess(await personas.Eddy.getAddress(), '0x00'),
        )
      })

      it('announces the change via a log', async () => {
        await expect(tx)
          .to.emit(controller, 'RemovedAccess')
          .withArgs(await personas.Eddy.getAddress())
      })

      describe('when called twice', () => {
        it('does not emit a log', async () => {
          const tx2 = await controller.removeAccess(
            await personas.Eddy.getAddress(),
          )
          const receipt = await tx2.wait()
          assert.equal(receipt.events?.length, 0)
        })
      })
    })
  })

  describe('#disableAccessCheck', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller.connect(personas.Eddy).disableAccessCheck(),
        ).to.be.revertedWith('Only callable by owner')
        assert.isTrue(await controller.checkEnabled())
      })
    })

    describe('when called by the owner', () => {
      let tx: Transaction
      beforeEach(async () => {
        await controller.addAccess(await personas.Eddy.getAddress())
        tx = await controller.disableAccessCheck()
      })

      it('sets checkEnabled to false', async () => {
        assert.isFalse(await controller.checkEnabled())
      })

      it('allows users with access', async () => {
        assert.isTrue(
          await controller.hasAccess(await personas.Eddy.getAddress(), '0x00'),
        )
      })

      it('allows users without access', async () => {
        assert.isTrue(
          await controller.hasAccess(await personas.Ned.getAddress(), '0x00'),
        )
      })

      it('announces the change via a log', async () => {
        await expect(tx).to.emit(controller, 'CheckAccessDisabled')
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
        await expect(
          controller.connect(personas.Eddy).enableAccessCheck(),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {
      let tx: Transaction
      beforeEach(async () => {
        await controller.disableAccessCheck()
        await controller.addAccess(await personas.Eddy.getAddress())
        tx = await controller.enableAccessCheck()
      })

      it('allows users with access', async () => {
        assert.isTrue(
          await controller.hasAccess(await personas.Eddy.getAddress(), '0x00'),
        )
      })

      it('does not allow users without access', async () => {
        assert.isFalse(
          await controller.hasAccess(await personas.Ned.getAddress(), '0x00'),
        )
      })

      it('announces the change via a log', async () => {
        expect(tx).to.emit(controller, 'CheckAccessEnabled')
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
