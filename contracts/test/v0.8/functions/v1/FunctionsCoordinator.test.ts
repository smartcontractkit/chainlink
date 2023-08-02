import { ethers } from 'hardhat'
import { expect } from 'chai'
import {
  getSetupFactory,
  FunctionsContracts,
  coordinatorConfig,
  ids,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts

beforeEach(async () => {
  ;({ contracts } = setup())
})

describe('Functions Coordinator', () => {
  describe('Config', () => {
    it('#fulfillmentGasPriceOverEstimationBP overestimates gas cost', async () => {
      const estimateWithNoOverestimaton =
        await contracts.coordinator.estimateCost(1, 0x0, 100_000, 20)

      await contracts.router.proposeConfigUpdate(
        ids.donId,
        ethers.utils.defaultAbiCoder.encode(
          [
            'uint32',
            'uint32',
            'uint32',
            'uint32',
            'uint32',
            'uint80',
            'uint16',
            'uint256',
            'int256',
          ],
          [
            ...Object.values({
              ...coordinatorConfig,
              fulfillmentGasPriceOverEstimationBP: 10_000,
            }),
          ],
        ),
      )
      await contracts.router.updateConfig(ids.donId)

      // Halve the gas price, which should be the same estimate because of fulfillmentGasPriceOverEstimationBP doubling the gas price
      const estimateWithOverestimaton =
        await contracts.coordinator.estimateCost(1, 0x0, 100_000, 10)

      expect(estimateWithNoOverestimaton).to.equal(estimateWithOverestimaton)
    })
  })
})
