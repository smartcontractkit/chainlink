import { contract, helpers, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { WhitelistedFactory } from '../../ethers/v0.5/WhitelistedFactory'

const whitelistedFactory = new WhitelistedFactory()
const provider = setup.provider()
let personas: setup.Personas
beforeAll(async () => {
  await setup.users(provider).then(u => (personas = u.personas))
})

describe('Whitelisted', () => {
  let whitelisted: contract.Instance<WhitelistedFactory>
  const deployment = setup.snapshot(provider, async () => {
    whitelisted = await whitelistedFactory.connect(personas.Carol).deploy()
  })
  beforeEach(deployment)

  it('has a limited public interface', () => {
    matchers.publicAbi(whitelistedFactory, [
      'acceptOwnership',
      'addToWhitelist',
      'owner',
      'removeFromWhitelist',
      'transferOwnership',
      'whitelisted',
    ])
  })

  describe('#addToWhitelist', () => {
    describe('when called by a stranger', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          whitelisted
            .connect(personas.Eddy)
            .addToWhitelist(personas.Eddy.address),
          'Only callable by owner',
        )
      })
    })

    describe('when called by the owner', () => {
      it('adds the address to the whitelist', async () => {
        const tx = await whitelisted
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
        const receipt = await tx.wait()

        assert.isTrue(await whitelisted.whitelisted(personas.Eddy.address))

        const event = helpers.findEventIn(
          receipt,
          whitelisted.interface.events.AddedToWhitelist,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Eddy.address)
      })
    })
  })

  describe('#removeFromWhitelist', () => {
    beforeEach(async () => {
      await whitelisted
        .connect(personas.Carol)
        .addToWhitelist(personas.Neil.address)
      assert.isTrue(await whitelisted.whitelisted(personas.Neil.address))
    })

    describe('when called by a stranger', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          whitelisted
            .connect(personas.Eddy)
            .removeFromWhitelist(personas.Neil.address),
          'Only callable by owner',
        )
      })
    })

    describe('when called by the owner', () => {
      it('removes the address from the whitelist', async () => {
        const tx = await whitelisted
          .connect(personas.Carol)
          .removeFromWhitelist(personas.Neil.address)
        const receipt = await tx.wait()

        assert.isFalse(await whitelisted.whitelisted(personas.Neil.address))

        const event = helpers.findEventIn(
          receipt,
          whitelisted.interface.events.RemovedFromWhitelist,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Neil.address)
      })
    })
  })
})
