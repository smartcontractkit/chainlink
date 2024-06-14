import { expect } from 'chai'
import {
  getSetupFactory,
  FunctionsContracts,
  coordinatorConfig,
  FunctionsRoles,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let roles: FunctionsRoles

beforeEach(async () => {
  ;({ contracts, roles } = setup())
})

describe('Functions Coordinator', () => {
  describe('Config', () => {
    it('non-owner is unable to update config', async () => {
      await expect(
        contracts.coordinator
          .connect(roles.stranger)
          .updateConfig(coordinatorConfig),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('Owner can update config', async () => {
      const beforeConfig = await contracts.coordinator.getConfig()
      await expect(
        contracts.coordinator.updateConfig({
          ...coordinatorConfig,
          donFee: 10,
        }),
      ).to.emit(contracts.coordinator, 'ConfigUpdated')
      const afterConfig = await contracts.coordinator.getConfig()
      expect(beforeConfig).to.not.equal(afterConfig)
    })

    it('returns the config set', async () => {
      const config = await contracts.coordinator
        .connect(roles.stranger)
        .getConfig()
      await Promise.all(
        Object.keys(coordinatorConfig).map((key) =>
          expect(config[key]).to.equal(
            coordinatorConfig[key as keyof typeof coordinatorConfig],
          ),
        ),
      )
    })

    it('#fulfillmentGasPriceOverEstimationBP overestimates gas cost', async () => {
      const estimateWithNoOverestimaton =
        await contracts.coordinator.estimateCost(1, 0x0, 100_000, 20)

      await contracts.coordinator.updateConfig({
        ...coordinatorConfig,
        fulfillmentGasPriceOverEstimationBP: 10_000,
      })

      // Halve the gas price, which should be the same estimate because of fulfillmentGasPriceOverEstimationBP doubling the gas price
      const estimateWithOverestimaton =
        await contracts.coordinator.estimateCost(1, 0x0, 100_000, 10)

      expect(estimateWithNoOverestimaton).to.equal(estimateWithOverestimaton)
    })
  })
})
