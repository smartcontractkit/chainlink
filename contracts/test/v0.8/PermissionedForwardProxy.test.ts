import { ethers } from 'hardhat'
import { publicAbi } from '../test-helpers/helpers'
import { expect, assert } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'

const PERMISSION_NOT_SET = 'PermissionNotSet'

let personas: Personas

let controllerFactory: ContractFactory
let counterFactory: ContractFactory
let controller: Contract
let counter: Contract

before(async () => {
  personas = (await getUsers()).personas
  controllerFactory = await ethers.getContractFactory(
    'src/v0.8/PermissionedForwardProxy.sol:PermissionedForwardProxy',
    personas.Carol,
  )
  counterFactory = await ethers.getContractFactory(
    'src/v0.8/tests/Counter.sol:Counter',
    personas.Carol,
  )
})

describe('PermissionedForwardProxy', () => {
  beforeEach(async () => {
    controller = await controllerFactory.connect(personas.Carol).deploy()
    counter = await counterFactory.connect(personas.Carol).deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(controller, [
      'forward',
      'setPermission',
      'removePermission',
      'getPermission',
      // Owned
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#setPermission', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Eddy)
            .setPermission(
              await personas.Carol.getAddress(),
              await personas.Eddy.getAddress(),
            ),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {
      it('adds the permission to the proxy', async () => {
        const tx = await controller
          .connect(personas.Carol)
          .setPermission(
            await personas.Carol.getAddress(),
            await personas.Eddy.getAddress(),
          )
        const receipt = await tx.wait()
        const eventLog = receipt?.events

        assert.equal(eventLog?.length, 1)
        assert.equal(eventLog?.[0].event, 'PermissionSet')
        assert.equal(eventLog?.[0].args?.[0], await personas.Carol.getAddress())
        assert.equal(eventLog?.[0].args?.[1], await personas.Eddy.getAddress())

        expect(
          await controller.getPermission(await personas.Carol.getAddress()),
        ).to.be.equal(await personas.Eddy.getAddress())
      })
    })
  })

  describe('#removePermission', () => {
    beforeEach(async () => {
      // Add permission before testing
      await controller
        .connect(personas.Carol)
        .setPermission(
          await personas.Carol.getAddress(),
          await personas.Eddy.getAddress(),
        )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Eddy)
            .removePermission(await personas.Carol.getAddress()),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {
      it('removes the permission to the proxy', async () => {
        const tx = await controller
          .connect(personas.Carol)
          .removePermission(await personas.Carol.getAddress())

        const receipt = await tx.wait()
        const eventLog = receipt?.events

        assert.equal(eventLog?.length, 1)
        assert.equal(eventLog?.[0].event, 'PermissionRemoved')
        assert.equal(eventLog?.[0].args?.[0], await personas.Carol.getAddress())

        expect(
          await controller.getPermission(await personas.Carol.getAddress()),
        ).to.be.equal(ethers.constants.AddressZero)
      })
    })
  })

  describe('#forward', () => {
    describe('when permission does not exist', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Carol)
            .forward(await personas.Eddy.getAddress(), '0x'),
        ).to.be.revertedWith(PERMISSION_NOT_SET)
      })
    })

    describe('when permission exists', () => {
      beforeEach(async () => {
        // Add permission before testing
        await controller
          .connect(personas.Carol)
          .setPermission(await personas.Carol.getAddress(), counter.address)
      })

      it('calls target successfully', async () => {
        await controller
          .connect(personas.Carol)
          .forward(
            counter.address,
            counter.interface.encodeFunctionData('increment'),
          )

        expect(await counter.count()).to.be.equal(1)
      })

      it('reverts when target reverts and bubbles up error', async () => {
        await expect(
          controller
            .connect(personas.Carol)
            .forward(
              counter.address,
              counter.interface.encodeFunctionData('alwaysRevertWithString'),
            ),
        ).to.be.revertedWith('always revert') // Revert strings should be bubbled up

        await expect(
          controller
            .connect(personas.Carol)
            .forward(
              counter.address,
              counter.interface.encodeFunctionData('alwaysRevert'),
            ),
        ).to.be.reverted // Javascript VM not able to parse custom errors defined on another contract
      })
    })
  })
})
