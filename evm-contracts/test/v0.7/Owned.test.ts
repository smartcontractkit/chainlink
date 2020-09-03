import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { OwnedTestHelperFactory } from '../../ethers/v0.7/OwnedTestHelperFactory'

const ownedTestHelperFactory = new OwnedTestHelperFactory()
const provider = setup.provider()

let personas: setup.Personas
let owner: ethers.Wallet
let nonOwner: ethers.Wallet
let newOwner: ethers.Wallet

beforeAll(async () => {
  const users = await setup.users(provider)
  personas = users.personas
  owner = personas.Carol
  nonOwner = personas.Neil
  newOwner = personas.Ned
})

describe('Owned', () => {
  let owned: contract.Instance<OwnedTestHelperFactory>
  const ownedEvents = ownedTestHelperFactory.interface.events

  beforeEach(async () => {
    owned = await ownedTestHelperFactory.connect(owner).deploy()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(ownedTestHelperFactory, [
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // test helper public methods
      'modifierOnlyOwner',
    ])
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

  describe('#onlyOwner modifier', () => {
    describe('when called by an owner', () => {
      it('successfully calls the method', async () => {
        const tx = await owned.connect(owner).modifierOnlyOwner()
        const receipt = await tx.wait()

        expect(h.findEventIn(receipt, ownedEvents.Here)).toBeDefined()
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', () =>
        matchers.evmRevert(owned.connect(nonOwner).modifierOnlyOwner()))
    })
  })

  describe('#transferOwnership', () => {
    describe('when called by an owner', () => {
      it('emits a log', async () => {
        const tx = await owned
          .connect(owner)
          .transferOwnership(newOwner.address)
        const receipt = await tx.wait()

        const event = h.findEventIn(
          receipt,
          ownedEvents.OwnershipTransferRequested,
        )
        expect(h.eventArgs(event).to).toEqual(newOwner.address)
        expect(h.eventArgs(event).from).toEqual(owner.address)
      })
    })
  })

  describe('when called by anyone but the owner', () => {
    it('successfully calls the method', () =>
      matchers.evmRevert(
        owned.connect(nonOwner).transferOwnership(newOwner.address),
      ))
  })

  describe('#acceptOwnership', () => {
    describe('after #transferOwnership has been called', () => {
      beforeEach(async () => {
        await owned.connect(owner).transferOwnership(newOwner.address)
      })

      it('allows the recipient to call it', async () => {
        const tx = await owned.connect(newOwner).acceptOwnership()
        const receipt = await tx.wait()

        const event = h.findEventIn(receipt, ownedEvents.OwnershipTransferred)
        expect(h.eventArgs(event).to).toEqual(newOwner.address)
        expect(h.eventArgs(event).from).toEqual(owner.address)
      })

      it('does not allow a non-recipient to call it', () =>
        matchers.evmRevert(owned.connect(nonOwner).acceptOwnership()))
    })
  })
})
