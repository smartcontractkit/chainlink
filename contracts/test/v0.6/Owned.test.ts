import { ethers } from 'hardhat'
import { publicAbi } from '../test-helpers/helpers'
import { assert, expect } from 'chai'
import { Signer, Contract, ContractFactory } from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'

let personas: Personas

let owner: Signer
let nonOwner: Signer
let newOwner: Signer

let ownedFactory: ContractFactory
let owned: Contract

before(async () => {
  personas = (await getUsers()).personas
  owner = personas.Carol
  nonOwner = personas.Neil
  newOwner = personas.Ned
  ownedFactory = await ethers.getContractFactory(
    'src/v0.6/Owned.sol:Owned',
    owner,
  )
})

describe('Owned', () => {
  beforeEach(async () => {
    owned = await ownedFactory.connect(owner).deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(owned, ['acceptOwnership', 'owner', 'transferOwnership'])
  })

  describe('#constructor', () => {
    it('assigns ownership to the deployer', async () => {
      const [actual, expected] = await Promise.all([
        owner.getAddress(),
        owned.owner(),
      ])

      assert.equal(actual, expected)
    })
  })

  describe('#transferOwnership', () => {
    describe('when called by an owner', () => {
      it('emits a log', async () => {
        await expect(
          owned.connect(owner).transferOwnership(await newOwner.getAddress()),
        )
          .to.emit(owned, 'OwnershipTransferRequested')
          .withArgs(await owner.getAddress(), await newOwner.getAddress())
      })
    })
  })

  describe('when called by anyone but the owner', () => {
    it('reverts', async () =>
      await expect(
        owned.connect(nonOwner).transferOwnership(await newOwner.getAddress()),
      ).to.be.reverted)
  })

  describe('#acceptOwnership', () => {
    describe('after #transferOwnership has been called', () => {
      beforeEach(async () => {
        await owned
          .connect(owner)
          .transferOwnership(await newOwner.getAddress())
      })

      it('allows the recipient to call it', async () => {
        await expect(owned.connect(newOwner).acceptOwnership())
          .to.emit(owned, 'OwnershipTransferred')
          .withArgs(await owner.getAddress(), await newOwner.getAddress())
      })

      it('does not allow a non-recipient to call it', async () =>
        await expect(owned.connect(nonOwner).acceptOwnership()).to.be.reverted)
    })
  })
})
