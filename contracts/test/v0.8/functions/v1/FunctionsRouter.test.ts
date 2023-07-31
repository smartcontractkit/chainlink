import { ethers } from 'hardhat'
import { expect } from 'chai'
import {
  getSetupFactory,
  FunctionsContracts,
  functionsRouterConfig,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts

beforeEach(async () => {
  ;({ contracts } = setup())
})

describe('Functions Router - Request lifecycle', () => {
  describe('Getters', () => {
    it('#typeAndVersion', async () => {
      expect(await contracts.router.typeAndVersion()).to.be.equal(
        'Functions Router v1.0.0',
      )
    })
    it('#subscriptionConfig', async () => {
      const subscriptionConfig = await contracts.router.getSubscriptionConfig()
      expect(subscriptionConfig.maxConsumers).to.be.equal(
        functionsRouterConfig.maxConsumers,
      )
      expect(subscriptionConfig.adminFee).to.be.equal(
        functionsRouterConfig.adminFee,
      )
      expect(subscriptionConfig.handleOracleFulfillmentSelector).to.be.equal(
        functionsRouterConfig.handleOracleFulfillmentSelector,
      )
      expect(subscriptionConfig.maxCallbackGasLimits.toString()).to.be.equal(
        functionsRouterConfig.maxCallbackGasLimits.toString(),
      )
    })
  })
})
