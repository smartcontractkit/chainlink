import { personas, checkPublicABI } from './support/helpers'
import { expectEvent, expectRevert } from 'openzeppelin-test-helpers'

contract.only('Owned', () => {
  const Owned = artifacts.require('OwnedTestHelper.sol')
  let owned, owner, nonOwner, newOwner

  before(async () => {
    owner = personas.Carol
    nonOwner = personas.Neil
    newOwner = personas.Ned
  })

  beforeEach(async () => {
    owned = await Owned.new({ from: owner })
  })

  it('has a limited public interface', () => {
    checkPublicABI(Owned, [
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // test helper public methods
      'modifierIfOwner',
      'modifierOnlyOwner',
    ])
  })

  describe('#constructor', () => {
    it('assigns ownership to the deployer', async () => {
      assert.equal(owner, await owned.owner.call())
    })
  })

  describe('#ifOwner modifier', () => {
    context('when called by an owner', () => {
      it('successfully calls the method', async () => {
        const { logs } = await owned.modifierIfOwner({ from: owner })
        expectEvent.inLogs(logs, 'Here')
      })
    })

    context('when called by anyone but the owner', () => {
      it('reverts', async () => {
        const { logs } = await owned.modifierIfOwner({ from: nonOwner })
        assert.equal(0, logs.length)
      })
    })
  })

  describe('#onlyOwner modifier', () => {
    context('when called by an owner', () => {
      it('successfully calls the method', async () => {
        const { logs } = await owned.modifierOnlyOwner({ from: owner })
        expectEvent.inLogs(logs, 'Here')
      })
    })

    context('when called by anyone but the owner', () => {
      it('reverts', async () => {
        await expectRevert(
          owned.modifierOnlyOwner({ from: nonOwner }),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#transferOwnership', () => {
    context('when called by an owner', () => {
      it('emits a log', async () => {
        const { logs } = await owned.transferOwnership(newOwner, { from: owner })
        expectEvent.inLogs(logs, 'OwnershipTransferRequested', {to: newOwner, from: owner})
      })
    })

    context('when called by anyone but the owner', () => {
      it('successfully calls the method', async () => {
        await expectRevert(
          owned.transferOwnership(newOwner, { from: nonOwner }),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#acceptOwnership', () => {
    context('after #transferOwnership has been called', () => {
      beforeEach(async () => {
        await owned.transferOwnership(newOwner, { from: owner })
      })

      it('allows the recipient to call it', async () => {
        const { logs } = await owned.acceptOwnership({ from: newOwner })
        expectEvent.inLogs(logs, 'OwnershipTransfered', {to: newOwner, from: owner})
      })

      it('does not allow a non-recipient to call it', async () => {
        await expectRevert(
          owned.acceptOwnership({ from: nonOwner }),
          'Must be requested to accept ownership',
        )
      })
    })
  })
})
