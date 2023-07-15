import { ethers } from 'hardhat'
import { expect } from 'chai'
import {
  getSetupFactory,
  FunctionsContracts,
  FunctionsRoles,
  functionsRouterConfig,
  stringToHex,
  anyValue,
  createSubscription,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let roles: FunctionsRoles

const donLabel = ethers.utils.formatBytes32String('1')

beforeEach(async () => {
  ;({ contracts, roles } = setup())
})

describe('Functions Router - Request lifecycle', () => {
  describe('Getters', () => {
    it('#typeAndVersion', async () => {
      expect(await contracts.router.typeAndVersion()).to.be.equal(
        'Functions Router v1',
      )
    })
    it('#adminFee', async () => {
      expect(await contracts.router.getAdminFee()).to.be.equal(
        functionsRouterConfig.adminFee,
      )
    })
  })
})
