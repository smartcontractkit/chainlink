import { ethers } from 'hardhat'
import { publicAbi } from '../test-helpers/helpers'
import { expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'

let personas: Personas

let controllerFactory: ContractFactory
let counterFactory: ContractFactory
let controller: Contract

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
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(controller, [
      'forward',
      'addPermission',
      'removePermission',
      'forwardPermissionList',
      // Owned
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#addPermission', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Eddy)
            .addPermission(
              await personas.Carol.getAddress(),
              await personas.Eddy.getAddress(),
            ),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {

        beforeEach(async () => {
          // Reset permission before testing
          await controller
            .connect(personas.Carol)
            .removePermission(await personas.Carol.getAddress())
        })

        it('adds the permission to the proxy', async () => {
          await controller
            .connect(personas.Carol)
            .addPermission(
              await personas.Carol.getAddress(),
              await personas.Eddy.getAddress(),
            )

          await expect(
            await controller.forwardPermissionList(
              await personas.Carol.getAddress(),
            ),
          ).to.be.equal(await personas.Eddy.getAddress())
        })
    })
  })

  describe('#removePermission', () => {
    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Eddy)
            .removePermission(
              await personas.Carol.getAddress(),
            ),
        ).to.be.revertedWith('Only callable by owner')
      })
    })

    describe('when called by the owner', () => {

        beforeEach(async () => {
          // Add permission before testing
          await controller
            .connect(personas.Carol)
            .addPermission(
              await personas.Carol.getAddress(),
              await personas.Eddy.getAddress(),
            )
        })

        it('removes the permission to the proxy', async () => {
          await controller
            .connect(personas.Carol)
            .removePermission(await personas.Carol.getAddress())

          await expect(
            await controller.forwardPermissionList(
              await personas.Carol.getAddress(),
            ),
          ).to.be.equal('0x0000000000000000000000000000000000000000') // Defaults to 0x0 address
        })
    })
  })


  describe('#removePermission', () => {
    describe('when permission does not exist', () => {
        beforeEach(async () => {
          // Reset permission before testing
          await controller
            .connect(personas.Carol)
            .removePermission(await personas.Carol.getAddress())
        })
      it('reverts', async () => {
        await expect(
          controller
            .connect(personas.Carol)
            .forward(await personas.Eddy.getAddress(), "0x"),
        ).to.be.revertedWith('Forwarding permission not found')
      })
    })

    describe('when permission exists', () => {
        let counter: Contract
        beforeEach(async () => {
          // Deploy Counter contract to call
          counter = await counterFactory.connect(personas.Carol).deploy()

          // Add permission before testing
          await controller
            .connect(personas.Carol)
            .addPermission(
              await personas.Carol.getAddress(),
              counter.address,
            )
        })

        it('calls target successfully', async () => {
            let encoder = new ethers.utils.Interface( [
                "function increment()"
            ]);
            let handler = encoder.encodeFunctionData("increment")

            await controller
              .connect(personas.Carol)
              .forward(counter.address, handler)

            await expect(await counter.count()).to.be.equal(1)
        })

        it('reverts when target reverts', async () => {
            let encoder = new ethers.utils.Interface( [
                "function alwaysRevert()"
            ]);
            let handler = encoder.encodeFunctionData("alwaysRevert")
            await expect(
              await controller
                .connect(personas.Carol)
                .forward(counter.address, handler),
            ).to.be.revertedWith('always revert')
        })
      })
  })
})
