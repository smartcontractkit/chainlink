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
    it('#config', async () => {
      const config = await contracts.router.getConfig()
      expect(config.maxConsumersPerSubscription).to.be.equal(
        functionsRouterConfig.maxConsumersPerSubscription,
      )
      expect(config.adminFee).to.be.equal(functionsRouterConfig.adminFee)
      expect(config.handleOracleFulfillmentSelector).to.be.equal(
        functionsRouterConfig.handleOracleFulfillmentSelector,
      )
      expect(config.maxCallbackGasLimits.toString()).to.be.equal(
        functionsRouterConfig.maxCallbackGasLimits.toString(),
      )
    })
  })
})
