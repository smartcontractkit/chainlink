import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ConfirmedOwnerTestHelperFactory } from '../../ethers/v0.7/ConfirmedOwnerTestHelperFactory'

const confirmedOwnerTestHelperFactory = new ConfirmedOwnerTestHelperFactory()
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

describe('ConfirmedOwner', () => {
  let confirmedOwner: contract.Instance<ConfirmedOwnerTestHelperFactory>
  const confirmedOwnerEvents = confirmedOwnerTestHelperFactory.interface.events

  beforeEach(async () => {
    confirmedOwner = await confirmedOwnerTestHelperFactory
      .connect(owner)
      .deploy()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(confirmedOwnerTestHelperFactory, [
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
        confirmedOwner.owner(),
      ])

      assert.equal(actual, expected)
    })
  })

  describe('#onlyOwner modifier', () => {
    describe('when called by an owner', () => {
      it('successfully calls the method', async () => {
        const tx = await confirmedOwner.connect(owner).modifierOnlyOwner()
        const receipt = await tx.wait()

        expect(h.findEventIn(receipt, confirmedOwnerEvents.Here)).toBeDefined()
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', () =>
        matchers.evmRevert(
          confirmedOwner.connect(nonOwner).modifierOnlyOwner(),
        ))
    })
  })

  describe('#transferOwnership', () => {
    describe('when called by an owner', () => {
      it('emits a log', async () => {
        const tx = await confirmedOwner
          .connect(owner)
          .transferOwnership(newOwner.address)
        const receipt = await tx.wait()

        const event = h.findEventIn(
          receipt,
          confirmedOwnerEvents.OwnershipTransferRequested,
        )
        expect(h.eventArgs(event).to).toEqual(newOwner.address)
        expect(h.eventArgs(event).from).toEqual(owner.address)
      })

      it('does not allow ownership transfer to self', async () => {
        await matchers.evmRevert(
          confirmedOwner.connect(owner).transferOwnership(owner.address),
          'Cannot transfer to self',
        )
      })
    })
  })

  describe('when called by anyone but the owner', () => {
    it('successfully calls the method', () =>
      matchers.evmRevert(
        confirmedOwner.connect(nonOwner).transferOwnership(newOwner.address),
      ))
  })

  describe('#acceptOwnership', () => {
    describe('after #transferOwnership has been called', () => {
      beforeEach(async () => {
        await confirmedOwner.connect(owner).transferOwnership(newOwner.address)
      })

      it('allows the recipient to call it', async () => {
        const tx = await confirmedOwner.connect(newOwner).acceptOwnership()
        const receipt = await tx.wait()

        const event = h.findEventIn(
          receipt,
          confirmedOwnerEvents.OwnershipTransferred,
        )
        expect(h.eventArgs(event).to).toEqual(newOwner.address)
        expect(h.eventArgs(event).from).toEqual(owner.address)
      })

      it('does not allow a non-recipient to call it', () =>
        matchers.evmRevert(confirmedOwner.connect(nonOwner).acceptOwnership()))
    })
  })
})
