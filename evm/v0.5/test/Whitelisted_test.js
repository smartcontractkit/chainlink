import * as h from './support/helpers'
import { expectEvent, expectRevert } from 'openzeppelin-test-helpers'

contract('Whitelisted', () => {
  const Whitelisted = artifacts.require('Whitelisted.sol')
  const personas = h.personas

  let whitelisted

  beforeEach(async () => {
    whitelisted = await Whitelisted.new({ from: personas.Carol })
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Whitelisted, [
      'acceptOwnership',
      'addToWhitelist',
      'owner',
      'removeFromWhitelist',
      'transferOwnership',
      'whitelisted',
    ])
  })

  describe('#addToWhitelist', () => {
    context('when called by a stranger', () => {
      it('reverts', async () => {
        await expectRevert(
          whitelisted.addToWhitelist(personas.Eddy, {
            from: personas.Eddy,
          }),
          'Only callable by owner',
        )
      })
    })

    context('when called by the owner', () => {
      it('adds the address to the whitelist', async () => {
        const { logs } = await whitelisted.addToWhitelist(personas.Eddy, {
          from: personas.Carol,
        })
        assert.isTrue(await whitelisted.whitelisted.call(personas.Eddy))
        expectEvent.inLogs(logs, 'AddedToWhitelist', {
          user: personas.Eddy,
        })
      })
    })
  })

  describe('#removeFromWhitelist', () => {
    beforeEach(async () => {
      await whitelisted.addToWhitelist(personas.Neil, {
        from: personas.Carol,
      })
      assert.isTrue(await whitelisted.whitelisted.call(personas.Neil))
    })

    context('when called by a stranger', () => {
      it('reverts', async () => {
        await expectRevert(
          whitelisted.removeFromWhitelist(personas.Neil, {
            from: personas.Eddy,
          }),
          'Only callable by owner',
        )
      })
    })

    context('when called by the owner', () => {
      it('removes the address from the whitelist', async () => {
        const { logs } = await whitelisted.removeFromWhitelist(personas.Neil, {
          from: personas.Carol,
        })
        assert.isFalse(await whitelisted.whitelisted.call(personas.Neil))
        await expectEvent.inLogs(logs, 'RemovedFromWhitelist', {
          user: personas.Neil,
        })
      })
    })
  })
})
