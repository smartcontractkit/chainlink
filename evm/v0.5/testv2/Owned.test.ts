import * as h from '../src/helpers'
import { OwnedTestHelperFactory } from '../src/generated'
import { makeTestProvider } from '../src/provider'
import { Instance } from '../src/contract'
import { ethers } from 'ethers'
import { assert } from 'chai'

const ownedTestHelperFactory = new OwnedTestHelperFactory()
const provider = makeTestProvider()

let personas: h.Personas
let owner: ethers.Wallet
let nonOwner: ethers.Wallet
let newOwner: ethers.Wallet

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)
  personas = rolesAndPersonas.personas
  owner = personas.Carol
  nonOwner = personas.Neil
  newOwner = personas.Ned
})

describe('Owned', () => {
  let owned: Instance<OwnedTestHelperFactory>
  const ownedEvents = ownedTestHelperFactory.interface.events

  beforeEach(async () => {
    owned = await ownedTestHelperFactory.connect(owner).deploy()
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(ownedTestHelperFactory, [
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
        h.assertActionThrows(owned.connect(nonOwner).modifierOnlyOwner()))
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
      h.assertActionThrows(
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

        const event = h.findEventIn(receipt, ownedEvents.OwnershipTransfered)
        expect(h.eventArgs(event).to).toEqual(newOwner.address)
        expect(h.eventArgs(event).from).toEqual(owner.address)
      })

      it('does not allow a non-recipient to call it', () =>
        h.assertActionThrows(owned.connect(nonOwner).acceptOwnership()))
    })
  })
})
