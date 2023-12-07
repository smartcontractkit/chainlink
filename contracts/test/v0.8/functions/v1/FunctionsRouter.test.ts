import { expect } from 'chai'
import { ethers } from 'hardhat'
import { stringToBytes } from '../../../test-helpers/helpers'
import {
  getSetupFactory,
  FunctionsContracts,
  functionsRouterConfig,
  FunctionsRoles,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let roles: FunctionsRoles

beforeEach(async () => {
  ;({ contracts, roles } = setup())
})

describe('Functions Router - Request lifecycle', () => {
  describe('Config', () => {
    it('#typeAndVersion', async () => {
      expect(await contracts.router.typeAndVersion()).to.be.equal(
        'Functions Router v2.0.0',
      )
    })
    it('non-owner is unable to update config', async () => {
      await expect(
        contracts.router
          .connect(roles.stranger)
          .updateConfig(functionsRouterConfig),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('owner can update config', async () => {
      const beforeConfig = await contracts.router.getConfig()
      await expect(
        contracts.router.updateConfig({
          ...functionsRouterConfig,
          adminFee: 10,
        }),
      ).to.emit(contracts.router, 'ConfigUpdated')
      const afterConfig = await contracts.router.getConfig()
      expect(beforeConfig).to.not.equal(afterConfig)
    })

    it('returns the config set', async () => {
      const config = await contracts.router.connect(roles.stranger).getConfig()
      await Promise.all(
        Object.keys(functionsRouterConfig).map((key) =>
          expect(config[key]).to.deep.equal(
            functionsRouterConfig[key as keyof typeof functionsRouterConfig],
          ),
        ),
      )
    })
  })
  describe('Allow List path', () => {
    it('non-owner is unable to set Allow List ID', async () => {
      await expect(
        contracts.router
          .connect(roles.stranger)
          .setAllowListId(ethers.utils.hexZeroPad(stringToBytes(''), 32)),
      ).to.be.revertedWith('Only callable by owner')
    })
  })
})
