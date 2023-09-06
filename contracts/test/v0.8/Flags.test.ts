import { ethers } from 'hardhat'
import { publicAbi } from '../test-helpers/helpers'
import { assert, expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'

let personas: Personas

let controllerFactory: ContractFactory
let flagsFactory: ContractFactory
let consumerFactory: ContractFactory

let controller: Contract
let flags: Contract
let consumer: Contract

before(async () => {
  personas = (await getUsers()).personas
  controllerFactory = await ethers.getContractFactory(
    'src/v0.8/shared/access/SimpleWriteAccessController.sol:SimpleWriteAccessController',
    personas.Nelly,
  )
  consumerFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FlagsTestHelper.sol:FlagsTestHelper',
    personas.Nelly,
  )
  flagsFactory = await ethers.getContractFactory(
    'src/v0.8/Flags.sol:Flags',
    personas.Nelly,
  )
})

describe('Flags', () => {
  beforeEach(async () => {
    controller = await controllerFactory.deploy()
    flags = await flagsFactory.deploy(controller.address)
    await flags.disableAccessCheck()
    consumer = await consumerFactory.deploy(flags.address)
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(flags, [
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
        await expect(flags.connect(personas.Nelly).raiseFlag(consumer.address))
          .to.emit(flags, 'FlagRaised')
          .withArgs(consumer.address)
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
          .addAccess(await personas.Neil.getAddress())
      })

      it('sets the flags', async () => {
        await flags.connect(personas.Neil).raiseFlag(consumer.address),
          assert.equal(true, await flags.getFlag(consumer.address))
      })
    })

    describe('when called by a non-enabled setter', () => {
      it('reverts', async () => {
        await expect(
          flags.connect(personas.Neil).raiseFlag(consumer.address),
        ).to.be.revertedWith('Not allowed to raise flags')
      })
    })

    describe('when called when there is no raisingAccessController', () => {
      beforeEach(async () => {
        await expect(
          flags
            .connect(personas.Nelly)
            .setRaisingAccessController(
              '0x0000000000000000000000000000000000000000',
            ),
        ).to.emit(flags, 'RaisingAccessControllerUpdated')
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
        await expect(flags.connect(personas.Neil).raiseFlag(consumer.address))
          .to.be.reverted
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
        await expect(
          flags.connect(personas.Nelly).raiseFlags([consumer.address]),
        )
          .to.emit(flags, 'FlagRaised')
          .withArgs(consumer.address)
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
          .addAccess(await personas.Neil.getAddress())
      })

      it('sets the flags', async () => {
        await flags.connect(personas.Neil).raiseFlags([consumer.address]),
          assert.equal(true, await flags.getFlag(consumer.address))
      })
    })

    describe('when called by a non-enabled setter', () => {
      it('reverts', async () => {
        await expect(
          flags.connect(personas.Neil).raiseFlags([consumer.address]),
        ).to.be.revertedWith('Not allowed to raise flags')
      })
    })

    describe('when called when there is no raisingAccessController', () => {
      beforeEach(async () => {
        await expect(
          flags
            .connect(personas.Nelly)
            .setRaisingAccessController(
              '0x0000000000000000000000000000000000000000',
            ),
        ).to.emit(flags, 'RaisingAccessControllerUpdated')

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
        await expect(
          flags.connect(personas.Neil).raiseFlags([consumer.address]),
        ).to.be.reverted
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
        await expect(
          flags.connect(personas.Nelly).lowerFlags([consumer.address]),
        )
          .to.emit(flags, 'FlagLowered')
          .withArgs(consumer.address)
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
        await expect(
          flags.connect(personas.Neil).lowerFlags([consumer.address]),
        ).to.be.revertedWith('Only callable by owner')
      })
    })
  })

  describe('#getFlag', () => {
    describe('if the access control is turned on', () => {
      beforeEach(async () => {
        await flags.connect(personas.Nelly).enableAccessCheck()
      })

      it('reverts', async () => {
        await expect(consumer.getFlag(consumer.address)).to.be.revertedWith(
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
        .raiseFlags([
          await personas.Neil.getAddress(),
          await personas.Norbert.getAddress(),
        ])
    })

    it('respects the access controls of #getFlag', async () => {
      await flags.connect(personas.Nelly).enableAccessCheck()

      await expect(consumer.getFlag(consumer.address)).to.be.revertedWith(
        'No access',
      )

      await flags.connect(personas.Nelly).addAccess(consumer.address)

      await consumer.getFlag(consumer.address)
    })

    it('returns the flags in the order they are requested', async () => {
      const response = await consumer.getFlags([
        await personas.Nelly.getAddress(),
        await personas.Neil.getAddress(),
        await personas.Ned.getAddress(),
        await personas.Norbert.getAddress(),
      ])

      assert.deepEqual([false, true, false, true], response)
    })
  })

  describe('#setRaisingAccessController', () => {
    let controller2: Contract

    beforeEach(async () => {
      controller2 = await controllerFactory.connect(personas.Nelly).deploy()
      await controller2.connect(personas.Nelly).enableAccessCheck()
    })

    it('updates access control rules', async () => {
      const neilAddress = await personas.Neil.getAddress()
      await controller.connect(personas.Nelly).addAccess(neilAddress)
      await flags.connect(personas.Neil).raiseFlags([consumer.address]) // doesn't raise

      await flags
        .connect(personas.Nelly)
        .setRaisingAccessController(controller2.address)

      await expect(
        flags.connect(personas.Neil).raiseFlags([consumer.address]),
      ).to.be.revertedWith('Not allowed to raise flags')
    })

    it('emits a log announcing the change', async () => {
      await expect(
        flags
          .connect(personas.Nelly)
          .setRaisingAccessController(controller2.address),
      )
        .to.emit(flags, 'RaisingAccessControllerUpdated')
        .withArgs(controller.address, controller2.address)
    })

    it('does not emit a log when there is no change', async () => {
      await flags
        .connect(personas.Nelly)
        .setRaisingAccessController(controller2.address)

      await expect(
        flags
          .connect(personas.Nelly)
          .setRaisingAccessController(controller2.address),
      ).to.not.emit(flags, 'RaisingAccessControllerUpdated')
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          flags
            .connect(personas.Neil)
            .setRaisingAccessController(controller2.address),
        ).to.be.revertedWith('Only callable by owner')
      })
    })
  })
})
