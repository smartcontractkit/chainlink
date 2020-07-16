import { contract, helpers, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { WhitelistedFactory } from '../../ethers/v0.6/WhitelistedFactory'
import { WhitelistedTestHelperFactory } from '../../ethers/v0.6/WhitelistedTestHelperFactory'

const whitelistedFactory = new WhitelistedTestHelperFactory()
const provider = setup.provider()
let personas: setup.Personas
beforeAll(async () => {
  await setup.users(provider).then(u => (personas = u.personas))
})

describe('Whitelisted', () => {
  let whitelisted: contract.Instance<WhitelistedFactory>
  const value = 17
  const deployment = setup.snapshot(provider, async () => {
    whitelisted = await whitelistedFactory.connect(personas.Carol).deploy(value)
  })
  beforeEach(deployment)

  it('has a limited public interface', () => {
    matchers.publicAbi(new WhitelistedFactory(), [
      'addToWhitelist',
      'disableWhitelist',
      'enableWhitelist',
      'removeFromWhitelist',
      'whitelisted',
      'whitelistEnabled',
      // Owned
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#addToWhitelist', () => {
    describe('when called by a non-owner', () => {
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

      it('allows whitelisted users', async () => {
        await whitelisted
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)

        matchers.bigNum(
          value,
          await whitelisted.connect(personas.Eddy).getValue(),
        )
      })
    })
  })

  describe('#removeFromWhitelist', () => {
    beforeEach(async () => {
      await whitelisted
        .connect(personas.Carol)
        .addToWhitelist(personas.Eddy.address)
      assert.isTrue(await whitelisted.whitelisted(personas.Eddy.address))
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          whitelisted
            .connect(personas.Eddy)
            .removeFromWhitelist(personas.Eddy.address),
          'Only callable by owner',
        )
      })
    })

    describe('when called by the owner', () => {
      it('removes the address from the whitelist', async () => {
        const tx = await whitelisted
          .connect(personas.Carol)
          .removeFromWhitelist(personas.Eddy.address)
        const receipt = await tx.wait()

        assert.isFalse(await whitelisted.whitelisted(personas.Eddy.address))

        const event = helpers.findEventIn(
          receipt,
          whitelisted.interface.events.RemovedFromWhitelist,
        )
        expect(helpers.eventArgs(event).user).toEqual(personas.Eddy.address)
      })

      it('does not allow non-whitelisted users', async () => {
        await whitelisted
          .connect(personas.Carol)
          .removeFromWhitelist(personas.Eddy.address)

        await matchers.evmRevert(whitelisted.connect(personas.Eddy).getValue())
      })
    })
  })

  describe('#whitelistEnabled', () => {
    it('defaults to true', async () => {
      assert(await whitelisted.whitelistEnabled())
    })
  })

  describe('#enableWhitelist', () => {
    beforeEach(async () => {
      await whitelisted
        .connect(personas.Carol)
        .addToWhitelist(personas.Eddy.address)
    })

    it('allows whitelisted users', async () => {
      await whitelisted.connect(personas.Carol).enableWhitelist()

      matchers.bigNum(
        value,
        await whitelisted.connect(personas.Eddy).getValue(),
      )
    })

    it('does not allow non-whitelisted users', async () => {
      await whitelisted.connect(personas.Carol).enableWhitelist()

      await matchers.evmRevert(whitelisted.connect(personas.Ned).getValue())
    })

    it('announces the change via a log', async () => {
      const tx = await whitelisted.connect(personas.Carol).enableWhitelist()
      const receipt = await tx.wait()

      assert(
        helpers.findEventIn(
          receipt,
          whitelisted.interface.events.WhitelistEnabled,
        ),
      )
    })

    it('reverts when called by a non-owner', async () => {
      await matchers.evmRevert(
        whitelisted.connect(personas.Eddy).enableWhitelist(),
        'Only callable by owner',
      )
    })
  })

  describe('#disableWhitelist', () => {
    beforeEach(async () => {
      await whitelisted
        .connect(personas.Carol)
        .addToWhitelist(personas.Eddy.address)
    })

    it('allows whitelisted users', async () => {
      await whitelisted.connect(personas.Carol).disableWhitelist()

      matchers.bigNum(
        value,
        await whitelisted.connect(personas.Eddy).getValue(),
      )
    })

    it('allows non-whitelisted users', async () => {
      await whitelisted.connect(personas.Carol).disableWhitelist()

      matchers.bigNum(value, await whitelisted.connect(personas.Ned).getValue())
    })

    it('announces the change via a log', async () => {
      const tx = await whitelisted.connect(personas.Carol).disableWhitelist()
      const receipt = await tx.wait()

      assert(
        helpers.findEventIn(
          receipt,
          whitelisted.interface.events.WhitelistDisabled,
        ),
      )
    })

    it('reverts when called by a non-owner', async () => {
      await matchers.evmRevert(
        whitelisted.connect(personas.Eddy).disableWhitelist(),
        'Only callable by owner',
      )
    })
  })
})
