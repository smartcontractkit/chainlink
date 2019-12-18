import { personas, checkPublicABI } from './support/helpers'
import { expectEvent, expectRevert } from 'openzeppelin-test-helpers'

contract('Owned', () => {
  const Owned = artifacts.require('OwnedTestHelper.sol')
  let owned, owner, nonOwner

  beforeEach(async () => {
    owner = personas.Carol
    nonOwner = personas.Neil
    owned = await Owned.new({ from: owner })
  })

  it('has a limited public interface', () => {
    checkPublicABI(Owned, [
      'owner',
      // test helper public methods
      'modifierOnlyOwner',
    ])
  })

  describe('#constructor', () => {
    it('assigns ownership to the deployer', async () => {
      assert.equal(owner, await owned.owner.call())
    })
  })

  describe('#modifierOnlyOwner', () => {
    context('when called by an owner', () => {
      it('successfully calls the method', async () => {
        const { logs } = await owned.modifierOnlyOwner({ from: owner })
        expectEvent.inLogs(logs, 'Here')
      })
    })

    context('when called by an owner', () => {
      it('successfully calls the method', async () => {
        expectRevert(
          owned.modifierOnlyOwner({ from: nonOwner }),
          'Only callable by owner',
        )
      })
    })
  })
})
